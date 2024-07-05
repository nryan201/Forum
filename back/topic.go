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

	db := dbConn()
	defer db.Close()

	var topic struct {
		ID          int
		Username    string
		Title       string
		Description string
	}

	// Query to get the topic details along with the author's username
	query := `
		SELECT t.id, u.username, t.title, t.description
		FROM topics t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = ?
	`

	err := db.QueryRow(query, topicID).Scan(&topic.ID, &topic.Username, &topic.Title, &topic.Description)
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
