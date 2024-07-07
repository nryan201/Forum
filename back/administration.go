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
	db := dbConn()
	defer db.Close()

	// Fetch all users
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

	// Fetch all topics
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

	// Fetch all comments
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

	// Fetch all categories
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

	// Fetch all hashtags
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

	// Fetch all reports
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
		Users:      users,
		Topics:     topics,
		Comments:   comments,
		Categories: categories,
		Hashtags:   hashtags,
		Reports:    reports,
	}

	err = tmplAdmin.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
func ModeratorHandle(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	// Fetch all reports
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

	data := struct {
		Reports []Report
	}{
		Reports: reports,
	}

	err = tmplModerator.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Handlers for deleting items

func DeleteUserHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.FormValue("user_id")

	db := dbConn()
	defer db.Close()

	_, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
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

	if action == "ignore" {
		_, err := db.Exec("UPDATE reports SET status = 'ignored' WHERE id = ?", reportID)
		if err != nil {
			log.Printf("Error ignoring report: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else if action == "delete" {
		var topicID sql.NullInt64
		var commentID sql.NullInt64
		err := db.QueryRow("SELECT topic_id, comment_id FROM reports WHERE id = ?", reportID).Scan(&topicID, &commentID)
		if err != nil {
			log.Printf("Error fetching report details: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if topicID.Valid {
			_, err = db.Exec("DELETE FROM topics WHERE id = ?", topicID.Int64)
			if err != nil {
				log.Printf("Error deleting topic: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		if commentID.Valid {
			_, err = db.Exec("DELETE FROM comments WHERE id = ?", commentID.Int64)
			if err != nil {
				log.Printf("Error deleting comment: %v", err)
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
	}

	http.Redirect(w, r, "/moderator", http.StatusSeeOther)
}
func ReportTopicHandler(w http.ResponseWriter, r *http.Request) {
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

	// Retrieve user_id from cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login page if not authenticated
		return
	}
	userID := cookie.Value

	db := dbConn()
	defer db.Close()

	_, err = db.Exec("INSERT INTO reports (topic_id, user_id, reason) VALUES (?, ?, 'Inappropriate content')", topicID, userID)
	if err != nil {
		log.Printf("Error reporting topic: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+topicID, http.StatusSeeOther) // Redirect back to the post's detail page
}
