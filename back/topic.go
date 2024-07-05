package back

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

var tmplPost = template.Must(template.ParseFiles("./template/html/postDetail.html"))

func addTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	userID := r.FormValue("userID") // Assume user ID is passed as a form value

	// SQL to insert the new topic into the database
	_, err := db.Exec("INSERT INTO topics (user_id, title, description) VALUES (?, ?, ?)", userID, title, description)
	if err != nil {
		log.Printf("Error inserting new topic: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
}

func getTopicsHandler(w http.ResponseWriter, r *http.Request) {
	dbQuery := `SELECT id, title, description FROM topics ORDER BY created_at DESC`
	rows, err := db.Query(dbQuery)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var topics []struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	for rows.Next() {
		var t struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
		}
		if err := rows.Scan(&t.ID, &t.Title, &t.Description); err != nil {
			log.Printf("Error scanning topics: %v", err)
			continue
		}
		topics = append(topics, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}
func postDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Get the topic ID from the query parameters
	topicID := r.URL.Query().Get("id")
	if topicID == "" {
		http.Error(w, "Topic ID is required", http.StatusBadRequest)
		return
	}

	// Retrieve user_id from cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login page if not authenticated
		return
	}
	userID := cookie.Value

	db := dbConn()
	defer db.Close()

	var topic struct {
		ID          int
		Username    string
		Title       string
		Description string
		IsOwner     bool
	}

	// Query to get the topic details along with the author's username and check ownership
	query := `
        SELECT t.id, u.username, t.title, t.description, (t.user_id = ?) AS is_owner
        FROM topics t
        JOIN users u ON t.user_id = u.id
        WHERE t.id = ?
    `

	err = db.QueryRow(query, userID, topicID).Scan(&topic.ID, &topic.Username, &topic.Title, &topic.Description, &topic.IsOwner)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Topic not found", http.StatusNotFound)
		} else {
			log.Printf("Error retrieving topic: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	err = tmplPost.Execute(w, topic)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func editPostHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/editpost" {
		http.NotFound(w, r)
		return
	}

	topicID := r.URL.Query().Get("id")
	if topicID == "" {
		http.Error(w, "Topic ID is required", http.StatusBadRequest)
		return
	}

	db := dbConn()
	defer db.Close()

	var topic struct {
		ID          int
		Title       string
		Description string
		Username    string
	}

	err := db.QueryRow(`
		SELECT t.id, t.title, t.description, u.username
		FROM topics t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = ?`, topicID).Scan(&topic.ID, &topic.Title, &topic.Description, &topic.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Topic not found", http.StatusNotFound)
		} else {
			log.Printf("Error retrieving topic: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	tmpl, err := template.ParseFiles("template/html/edit.html") // Path to your edit template
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		ID          int
		Title       string
		Description string
		Username    string
	}{
		ID:          topic.ID,
		Title:       topic.Title,
		Description: topic.Description,
		Username:    topic.Username,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
func editHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	topicID := r.FormValue("id")
	title := r.FormValue("title")
	description := r.FormValue("description")

	// Retrieve user_id from cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login page if not authenticated
		return
	}
	userID := cookie.Value

	db := dbConn()
	defer db.Close()

	// Check if the user is the owner of the topic
	var ownerID string
	err = db.QueryRow("SELECT user_id FROM topics WHERE id = ?", topicID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Topic not found", http.StatusNotFound)
		} else {
			log.Printf("Error retrieving topic owner: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if ownerID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Update the topic
	_, err = db.Exec("UPDATE topics SET title = ?, description = ? WHERE id = ?", title, description, topicID)
	if err != nil {
		log.Printf("Error updating topic: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+topicID, http.StatusSeeOther) // Redirect to the updated post's page
}
