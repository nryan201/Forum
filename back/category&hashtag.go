package back

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func addCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("entering addCategoryHandler")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("categoryName")

	_, err := db.Exec("INSERT INTO categories (name) VALUES (?)", name)
	if err != nil {
		log.Printf("Error inserting new category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
}

func listCategories() ([]Category, error) {
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func addHashtagHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("hashtagsName")

	_, err := db.Exec("INSERT INTO hashtags (name) VALUES (?)", name)
	if err != nil {
		log.Printf("Error inserting new hashtag: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
}

func listHashtags() ([]Hashtag, error) {
	rows, err := db.Query("SELECT id, name FROM hashtags")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hashtags []Hashtag
	for rows.Next() {
		var hashtag Hashtag
		if err := rows.Scan(&hashtag.ID, &hashtag.Name); err != nil {
			return nil, err
		}
		hashtags = append(hashtags, hashtag)
	}
	return hashtags, nil
}
