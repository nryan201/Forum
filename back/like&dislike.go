package back

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func addLikeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		topicID := r.FormValue("topicID")
		userID := r.FormValue("userID")

		_, err := db.Exec("INSERT INTO likes (topic_id, user_id) VALUES (?, ?)", topicID, userID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "like added"})
	}
}
func addDislikeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		topicID := r.FormValue("topicID")
		userID := r.FormValue("userID")

		_, err := db.Exec("INSERT INTO dislikes (topic_id, user_id) VALUES (?, ?)", topicID, userID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "dislike added"})
	}
}
