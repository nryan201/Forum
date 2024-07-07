package back

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

var tmplmessage = template.Must(template.ParseFiles("./template/html/messagerie/messages.html"))
var tmplconversationDetail = template.Must(template.ParseFiles("./template/html/messagerie/conversation_detail.html"))

func getUsers(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	rows, err := db.Query("SELECT id, username FROM users")
	if err != nil {
		log.Printf("Error querying users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		users = append(users, user)
	}

	for _, user := range users {
		log.Printf("User: ID=%s, Username=%s\n", user.ID, user.Username)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	var message struct {
		ReceiverUsername string `json:"receiverUsername"`
		Content          string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if message.ReceiverUsername == "" || message.Content == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Retrieve sender ID from the cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	senderID := cookie.Value

	// Retrieve sender username for additional checks or logs
	var currentUsername string
	db := dbConn()
	defer db.Close()
	err = db.QueryRow("SELECT username FROM users WHERE id = ?", senderID).Scan(&currentUsername)
	if err != nil {
		log.Printf("Error retrieving sender username: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Retrieve receiver ID from the database using receiver's username
	var receiverID string
	err = db.QueryRow("SELECT id FROM users WHERE username = ?", message.ReceiverUsername).Scan(&receiverID)
	if err != nil {
		log.Printf("Error querying receiver ID: %v", err)
		http.Error(w, "Receiver not found", http.StatusNotFound)
		return
	}

	// Check if sender and receiver are the same
	if senderID == receiverID {
		http.Error(w, "Cannot send messages to yourself", http.StatusBadRequest)
		return
	}

	log.Printf("Current Username: %s, Sender ID: %s, Receiver ID: %s, Content: %s", currentUsername, senderID, receiverID, message.Content)

	// Insert the message into the database
	_, err = db.Exec("INSERT INTO messenger (sender_id, receiver_id, content) VALUES (?, ?, ?)", senderID, receiverID, message.Content)
	if err != nil {
		log.Printf("Error inserting message: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"success": "Message sent"})
}

func messagePage(w http.ResponseWriter, r *http.Request) {
	tmplmessage.Execute(w, nil)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from the cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := cookie.Value

	db := dbConn()
	defer db.Close()

	rows, err := db.Query(`SELECT messenger.id, messenger.sender_id, messenger.receiver_id, messenger.content, messenger.created_at, users.username 
                           FROM messenger 
                           JOIN users ON messenger.sender_id = users.id 
                           WHERE messenger.receiver_id = ? 
                           ORDER BY messenger.created_at ASC`, userID)
	if err != nil {
		log.Printf("Error querying messages: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt, &message.Username); err != nil {
			log.Printf("Error scanning message: %v", err)
			continue
		}
		messages = append(messages, message)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func conversationsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/html/messagerie/conversations.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func getConversationMessages(w http.ResponseWriter, r *http.Request) {
	conversationWith := r.URL.Path[len("/conversation/messages/"):]

	// Retrieve user ID from the cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := cookie.Value

	db := dbConn()
	defer db.Close()

	rows, err := db.Query(`SELECT messenger.id, messenger.sender_id, messenger.receiver_id, messenger.content, messenger.created_at, users.username 
                           FROM messenger 
                           JOIN users ON messenger.sender_id = users.id 
                           WHERE (messenger.receiver_id = ? AND users.username = ?) OR (messenger.sender_id = ? AND messenger.receiver_id = (SELECT id FROM users WHERE username = ?))
                           ORDER BY messenger.created_at ASC`, userID, conversationWith, userID, conversationWith)
	if err != nil {
		log.Printf("Error querying messages: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt, &message.Username); err != nil {
			log.Printf("Error scanning message: %v", err)
			continue
		}
		messages = append(messages, message)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func conversationDetailPage(w http.ResponseWriter, r *http.Request) {
	// Extract receiverID from URL path
	receiverID := r.URL.Path[len("/conversation/"):]

	// Retrieve user ID from the cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		// Redirect or handle error if the cookie is not found
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userID := cookie.Value

	// Data to pass to the template
	data := struct {
		UserID     string
		ReceiverID string
	}{
		UserID:     userID,
		ReceiverID: receiverID,
	}

	// Execute the conversation detail template with userID and receiverID
	err = tmplconversationDetail.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getConversations(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := cookie.Value

	db := dbConn()
	defer db.Close()

	rows, err := db.Query(`
        SELECT DISTINCT users.username
        FROM messenger
        JOIN users ON (messenger.sender_id = users.id OR messenger.receiver_id = users.id)
        WHERE messenger.sender_id = ? OR messenger.receiver_id = ?`, userID, userID)
	if err != nil {
		log.Printf("Error querying conversations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var conversations []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			log.Printf("Error scanning conversation: %v", err)
			continue
		}
		conversations = append(conversations, username)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}
