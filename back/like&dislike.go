package back

import (
	"encoding/json"
	"log"
	"net/http"
)

func addLikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	topicID := r.FormValue("topicID")
	userID := r.FormValue("userID")

	_, err := db.Exec("INSERT INTO likes (topic_id, user_id) VALUES (?, ?)", topicID, userID)
	if err != nil {
		log.Printf("Error inserting like: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "like added"})
}

func addDislikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	topicID := r.FormValue("topicID")
	userID := r.FormValue("userID")

	_, err := db.Exec("INSERT INTO dislikes (topic_id, user_id) VALUES (?, ?)", topicID, userID)
	if err != nil {
		log.Printf("Error inserting dislike: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "dislike added"})
}
