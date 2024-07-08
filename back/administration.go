package back

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

var tmplAdmin = template.Must(template.ParseFiles("./template/html/admin.html"))
var tmplModerator = template.Must(template.ParseFiles("./template/html/moderator.html"))

func AdminHandle(w http.ResponseWriter, r *http.Request) {
	// Lire le cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	currentUserID := cookie.Value

	db := dbConn()
	defer db.Close()

	// Vérifier le rôle de l'utilisateur
	var role string
	err = db.QueryRow("SELECT role FROM users WHERE id = ?", currentUserID).Scan(&role)
	if err != nil {
		log.Printf("Error fetching user role: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if role != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Récupérer tous les utilisateurs
	userRows, err := db.Query("SELECT id, username, name, email, role FROM users")
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer userRows.Close()

	var users []User
	for userRows.Next() {
		var user User
		if err := userRows.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Role); err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		users = append(users, user)
	}

	// Récupérer tous les topics
	topicRows, err := db.Query("SELECT id, user_id, title, description, created_at FROM topics")
	if err != nil {
		log.Printf("Error fetching topics: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer topicRows.Close()

	var topics []Topic
	for topicRows.Next() {
		var topic Topic
		if err := topicRows.Scan(&topic.ID, &topic.UserID, &topic.Title, &topic.Description, &topic.CreatedAt); err != nil {
			log.Printf("Error scanning topic: %v", err)
			continue
		}
		topics = append(topics, topic)
	}

	// Récupérer tous les commentaires
	commentRows, err := db.Query("SELECT id, topic_id, user_id, content, created_at FROM comments")
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer commentRows.Close()

	var comments []Comment
	for commentRows.Next() {
		var comment Comment
		if err := commentRows.Scan(&comment.ID, &comment.TopicID, &comment.UserID, &comment.Content, &comment.CreatedAt); err != nil {
			log.Printf("Error scanning comment: %v", err)
			continue
		}
		comments = append(comments, comment)
	}

	// Récupérer toutes les catégories
	categoryRows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer categoryRows.Close()

	var categories []Category
	for categoryRows.Next() {
		var category Category
		if err := categoryRows.Scan(&category.ID, &category.Name); err != nil {
			log.Printf("Error scanning category: %v", err)
			continue
		}
		categories = append(categories, category)
	}

	// Récupérer tous les hashtags
	hashtagRows, err := db.Query("SELECT id, name FROM hashtags")
	if err != nil {
		log.Printf("Error fetching hashtags: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer hashtagRows.Close()

	var hashtags []Hashtag
	for hashtagRows.Next() {
		var hashtag Hashtag
		if err := hashtagRows.Scan(&hashtag.ID, &hashtag.Name); err != nil {
			log.Printf("Error scanning hashtag: %v", err)
			continue
		}
		hashtags = append(hashtags, hashtag)
	}

	// Récupérer tous les rapports
	reportRows, err := db.Query("SELECT id, topic_id, comment_id, user_id, reason, created_at, status FROM reports")
	if err != nil {
		log.Printf("Error fetching reports: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer reportRows.Close()

	var reports []Report
	for reportRows.Next() {
		var report Report
		if err := reportRows.Scan(&report.ID, &report.TopicID, &report.CommentID, &report.UserID, &report.Reason, &report.CreatedAt, &report.Status); err != nil {
			log.Printf("Error scanning report: %v", err)
			continue
		}
		reports = append(reports, report)
	}

	data := AdminData{
		CurrentUserID: currentUserID,
		Users:         users,
		Topics:        topics,
		Comments:      comments,
		Categories:    categories,
		Hashtags:      hashtags,
		Reports:       reports,
	}

	err = tmplAdmin.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func DeleteUserHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.FormValue("user_id")

	// Retrieve the current user's ID from the cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	currentUserID := cookie.Value

	if currentUserID == userID {
		http.Error(w, "You cannot delete your own account", http.StatusForbidden)
		return
	}

	db := dbConn()
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeleteTopicHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	topicID := r.FormValue("topic_id")

	db := dbConn()
	defer db.Close()

	_, err := db.Exec("DELETE FROM topics WHERE id = ?", topicID)
	if err != nil {
		log.Printf("Error deleting topic: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeleteCommentHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	commentID := r.FormValue("comment_id")

	db := dbConn()
	defer db.Close()

	_, err := db.Exec("DELETE FROM comments WHERE id = ?", commentID)
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeleteCategoryHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categoryID := r.FormValue("category_id")

	db := dbConn()
	defer db.Close()

	_, err := db.Exec("DELETE FROM categories WHERE id = ?", categoryID)
	if err != nil {
		log.Printf("Error deleting category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeleteHashtagHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hashtagID := r.FormValue("hashtag_id")

	db := dbConn()
	defer db.Close()

	_, err := db.Exec("DELETE FROM hashtags WHERE id = ?", hashtagID)
	if err != nil {
		log.Printf("Error deleting hashtag: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func HandleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reportID := r.FormValue("report_id")
	action := r.FormValue("action")

	db := dbConn()
	defer db.Close()

	switch action {
	case "ignore":
		_, err := db.Exec("UPDATE reports SET status = 'ignored' WHERE id = ?", reportID)
		if err != nil {
			log.Printf("Error ignoring report: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	case "delete":
		var topicID string
		var commentID sql.NullInt64
		err := db.QueryRow("SELECT topic_id, comment_id FROM reports WHERE id = ?", reportID).Scan(&topicID, &commentID)
		if err != nil {
			log.Printf("Error fetching report details: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if topicID != "" {
			_, err = db.Exec("DELETE FROM topics WHERE id = ?", topicID)
			if err != nil {
				log.Printf("Error deleting topic: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		_, err = db.Exec("UPDATE reports SET status = 'resolved' WHERE id = ?", reportID)
		if err != nil {
			log.Printf("Error resolving report: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	case "deleteReport":
		_, err := db.Exec("DELETE FROM reports WHERE id = ?", reportID)
		if err != nil {
			log.Printf("Error deleting report: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func PromoteUserHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.FormValue("user_id")
	role := r.FormValue("role")

	if role != "user" && role != "moderator" && role != "admin" {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	// Retrieve the current user's ID from the cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	currentUserID := cookie.Value

	// Prevent the current user from promoting themselves
	if currentUserID == userID {
		http.Error(w, "You cannot change your own role", http.StatusForbidden)
		return
	}

	db := dbConn()
	defer db.Close()

	_, err = db.Exec("UPDATE users SET role = ? WHERE id = ?", role, userID)
	if err != nil {
		log.Printf("Error promoting user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
