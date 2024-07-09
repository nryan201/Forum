package back

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

var tmplPost = template.Must(template.ParseFiles("./template/html/postDetail.html"))
var tmplReport = template.Must(template.ParseFiles("./template/html/reportDetail.html"))

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
	log.Printf("New topic added: %s", title)
	log.Println("userID", userID)
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

	db := dbConn()
	defer db.Close()

	var topic TopicDetail

	// Query to get the topic details along with the author's username
	query := `
		SELECT t.id, u.username, t.title, t.description, CASE WHEN t.user_id = ? THEN 1 ELSE 0 END AS is_owner
		FROM topics t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = ?
	`

	cookie, err := r.Cookie("user_id")
	userID := ""
	isAuthenticated := false
	role := ""
	if err == nil {
		userID = cookie.Value
		isAuthenticated = true

		// Fetch user role
		err = db.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
		if err != nil {
			log.Printf("Error fetching user role: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

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

	// Set authentication status and role
	topic.IsAuthenticated = isAuthenticated
	topic.Role = role
	log.Println("role", role)

	// Query to get comments for the topic
	commentQuery := `
		SELECT c.id, u.username, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.topic_id = ?
		ORDER BY c.created_at ASC
	`

	rows, err := db.Query(commentQuery, topicID)
	if err != nil {
		log.Printf("Error retrieving comments: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Username, &comment.Content, &comment.CreatedAt)
		if err != nil {
			log.Printf("Error scanning comment: %v", err)
			continue
		}
		topic.Comments = append(topic.Comments, comment)
	}

	// Query to get categories for the topic
	categoryQuery := `
		SELECT c.id, c.name
		FROM categories c
		JOIN topic_categories tc ON c.id = tc.category_id
		WHERE tc.topic_id = ?
	`

	categoryRows, err := db.Query(categoryQuery, topicID)
	if err != nil {
		log.Printf("Error retrieving categories: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer categoryRows.Close()

	for categoryRows.Next() {
		var category Category
		err := categoryRows.Scan(&category.ID, &category.Name)
		if err != nil {
			log.Printf("Error scanning category: %v", err)
			continue
		}
		topic.Categories = append(topic.Categories, category)
	}

	// Query to get hashtags for the topic
	hashtagQuery := `
		SELECT h.id, h.name
		FROM hashtags h
		JOIN topic_hashtags th ON h.id = th.hashtag_id
		WHERE th.topic_id = ?
	`

	hashtagRows, err := db.Query(hashtagQuery, topicID)
	if err != nil {
		log.Printf("Error retrieving hashtags: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer hashtagRows.Close()

	for hashtagRows.Next() {
		var hashtag Hashtag
		err := hashtagRows.Scan(&hashtag.ID, &hashtag.Name)
		if err != nil {
			log.Printf("Error scanning hashtag: %v", err)
			continue
		}
		topic.Hashtags = append(topic.Hashtags, hashtag)
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

func addCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	topicID := r.FormValue("topic_id")
	content := r.FormValue("content")

	// Retrieve user_id from cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login page if not authenticated
		return
	}
	userID := cookie.Value

	db := dbConn()
	defer db.Close()

	_, err = db.Exec("INSERT INTO comments (topic_id, user_id, content) VALUES (?, ?, ?)", topicID, userID, content)
	if err != nil {
		log.Printf("Error inserting new comment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+topicID, http.StatusSeeOther) // Redirect back to the post's detail page
}
func reportTopicDetailHandler(w http.ResponseWriter, r *http.Request) {
	topicID := r.URL.Query().Get("topic_id")
	if topicID == "" {
		http.Error(w, "Topic ID is required", http.StatusBadRequest)
		return
	}

	data := struct {
		TopicID string
	}{
		TopicID: topicID,
	}

	err := tmplReport.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func submitReportHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("submitReportHandler reached")
	if r.Method != http.MethodPost {
		log.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	topicID := r.FormValue("topic_id")
	reason := r.FormValue("reason")

	// Retrieve user_id from cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login page if not authenticated
		return
	}
	userID := cookie.Value

	db := dbConn()
	defer db.Close()

	_, err = db.Exec("INSERT INTO reports (topic_id, user_id, reason, status) VALUES (?, ?, ?, 'pending')", topicID, userID, reason)
	if err != nil {
		log.Printf("Error inserting report: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+topicID, http.StatusSeeOther)
}
