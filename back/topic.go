package back

import (
	"encoding/json"
	"log"
	"net/http"
)

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
