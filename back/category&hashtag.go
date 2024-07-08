package back

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func addHashtagHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")

		_, err := db.Exec("INSERT INTO hashtags (name) VALUES (?)", name)
		if err != nil {
			log.Printf("Error inserting new hashtag: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/hashtags", http.StatusSeeOther)
	}
}

func listHashtagsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM hashtags")
		if err != nil {
			log.Printf("Error querying hashtags: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var hashtags []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		for rows.Next() {
			var hashtag struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}
			if err := rows.Scan(&hashtag.ID, &hashtag.Name); err != nil {
				log.Printf("Error scanning hashtag: %v", err)
				continue
			}
			hashtags = append(hashtags, hashtag)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(hashtags)
	}
}
func addCategoryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")

		_, err := db.Exec("INSERT INTO categories (name) VALUES (?)", name)
		if err != nil {
			log.Printf("Error inserting new category: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/categories", http.StatusSeeOther)
	}
}

func listCategoriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM categories")
		if err != nil {
			log.Printf("Error querying categories: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var categories []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		for rows.Next() {
			var category struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}
			if err := rows.Scan(&category.ID, &category.Name); err != nil {
				log.Printf("Error scanning category: %v", err)
				continue
			}
			categories = append(categories, category)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
	}
}
