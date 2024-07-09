package back

import (
	_ "database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"
	_ "time"
)

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

func listTopics() ([]Topic, error) {
	rows, err := db.Query("SELECT id, user_id, title, description, created_at FROM topics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var topic Topic
		if err := rows.Scan(&topic.ID, &topic.UserID, &topic.Title, &topic.Description, &topic.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, nil
}

func addCategoryHandler(w http.ResponseWriter, r *http.Request) {
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

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := strings.TrimPrefix(r.URL.Path, "/category/")
	topics, categoryName, err := getTopicsByCategory(categoryID)
	if err != nil {
		log.Printf("Error retrieving topics for category %v: %v", categoryID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	hashtags, err := listHashtags() // Fetch the list of hashtags
	if err != nil {
		log.Printf("Error listing hashtags: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	categories, err := listCategories() // Fetch the list of categories
	if err != nil {
		log.Printf("Error listing categories: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Topics       []Topic
		CategoryName string
		Hashtags     []Hashtag
		Categories   []Category
	}{
		Topics:       topics,
		CategoryName: categoryName,
		Hashtags:     hashtags,
		Categories:   categories,
	}

	tmpl, err := template.ParseFiles("template/html/category.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func hashtagHandler(w http.ResponseWriter, r *http.Request) {
	hashtagID := strings.TrimPrefix(r.URL.Path, "/hashtag/")
	topics, hashtagName, err := getTopicsByHashtag(hashtagID)
	if err != nil {
		log.Printf("Error retrieving topics for hashtag %v: %v", hashtagID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	hashtags, err := listHashtags() // Fetch the list of hashtags
	if err != nil {
		log.Printf("Error listing hashtags: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	categories, err := listCategories() // Fetch the list of categories
	if err != nil {
		log.Printf("Error listing categories: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Topics      []Topic
		HashtagName string
		Hashtags    []Hashtag
		Categories  []Category
	}{
		Topics:      topics,
		HashtagName: hashtagName,
		Hashtags:    hashtags,
		Categories:  categories,
	}

	tmpl, err := template.ParseFiles("template/html/hashtag.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func getTopicsByCategory(categoryID string) ([]Topic, string, error) {
	var categoryName string
	err := db.QueryRow("SELECT name FROM categories WHERE id = ?", categoryID).Scan(&categoryName)
	if err != nil {
		return nil, "", err
	}

	rows, err := db.Query(`
		SELECT t.id, t.user_id, t.title, t.description, t.created_at 
		FROM topics t 
		JOIN topic_categories tc ON t.id = tc.topic_id 
		WHERE tc.category_id = ?`, categoryID)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var topic Topic
		if err := rows.Scan(&topic.ID, &topic.UserID, &topic.Title, &topic.Description, &topic.CreatedAt); err != nil {
			return nil, "", err
		}
		topics = append(topics, topic)
	}
	return topics, categoryName, nil
}

func getTopicsByHashtag(hashtagID string) ([]Topic, string, error) {
	rows, err := db.Query(`
		SELECT t.id, t.user_id, t.title, t.description, t.created_at 
		FROM topics t 
		JOIN topic_hashtags th ON t.id = th.topic_id 
		WHERE th.hashtag_id = ?`, hashtagID)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var topic Topic
		if err := rows.Scan(&topic.ID, &topic.UserID, &topic.Title, &topic.Description, &topic.CreatedAt); err != nil {
			return nil, "", err
		}
		topics = append(topics, topic)
	}

	var hashtagName string
	err = db.QueryRow(`SELECT name FROM hashtags WHERE id = ?`, hashtagID).Scan(&hashtagName)
	if err != nil {
		return nil, "", err
	}

	return topics, hashtagName, nil
}
