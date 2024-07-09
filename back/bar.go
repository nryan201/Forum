package back

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// searchTopicsHandler handles the AJAX request and returns JSON data.
func searchTopicsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("search")
	if query == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	dbQuery := `SELECT id, title, description FROM topics WHERE title LIKE ? OR description LIKE ?`
	rows, err := db.Query(dbQuery, "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%")
	if err != nil {
		log.Printf("Error querying topics: %v", err)
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
