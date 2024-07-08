package back

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

func Server() {
	OpenDB()
	defer db.Close()

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("template/css/"))))
	http.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("template/ressource/image/"))))
	http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("template/script/"))))

	http.HandleFunc("/", routeHandler)
	http.HandleFunc("/block", BlockHandler)
	http.HandleFunc("/login", loginUser)
	http.HandleFunc("/addUser", addUser)
	http.HandleFunc("/addTopic", addTopicHandler) // Add this for creating new topics
	http.HandleFunc("/topics", getTopicsHandler)  // Add this to fetch and display topics on the homepage

	//logout
	http.HandleFunc("/logout", logout)

	http.HandleFunc("/completeProfile", completeProfile)
	http.HandleFunc("/completeProfileGithub", completeProfileGithub)
	//Handle Social Media
	http.HandleFunc("/loginFacebook", handleFacebookLogin)
	http.HandleFunc("/callbackFacebook", handleFacebookCallback)
	http.HandleFunc("/loginGoogle", handleGoogleLogin)
	http.HandleFunc("/callbackGoogle", handleGoogleCallback)
	http.HandleFunc("/loginGithub", handleGithubLogin)
	http.HandleFunc("/callbackGithub", handleGithubCallback)

	//Handle Accueil, Contact, Profil, Post
	http.HandleFunc("/accueil", AccueilHandle)
	http.HandleFunc("/contact", ContactHandle)
	http.HandleFunc("/profil", profilePage)

	//Handle Post
	http.HandleFunc("/createpost", PostHandle) // celui sert a la creation de post
	http.HandleFunc("/submit-post", postHandler)
	http.HandleFunc("/post", postDetailHandler)        // celui sert a la visualisation de post
	http.HandleFunc("/add-comment", addCommentHandler) // celui sert a la creation de commentaire
	http.HandleFunc("/editpost", editPostHandle)       // for editing posts
	http.HandleFunc("/submit-edit", editHandler)

	//Handle Category and Hashtag
	http.HandleFunc("/addHashtag", addHashtagHandler(db))
	http.HandleFunc("/hashtags", listHashtagsHandler(db))
	http.HandleFunc("/addCategory", addCategoryHandler(db))
	http.HandleFunc("/categories", listCategoriesHandler(db))

	// Admin and moderation routes
	http.HandleFunc("/admin", AdminHandle)
	http.HandleFunc("/admin/delete-user", DeleteUserHandle)
	http.HandleFunc("/admin/promote-user", PromoteUserHandle)
	http.HandleFunc("/admin/delete-topic", DeleteTopicHandle)
	http.HandleFunc("/admin/delete-comment", DeleteCommentHandle)
	http.HandleFunc("/admin/delete-category", DeleteCategoryHandle)
	http.HandleFunc("/admin/delete-hashtag", DeleteHashtagHandle)
	http.HandleFunc("/admin/handle-report", HandleReport)
	http.HandleFunc("/moderator", ModeratorHandle)
	http.HandleFunc("/moderator/handle-report", HandleReport)
	http.HandleFunc("/report-topic", reportTopicDetailHandler)
	http.HandleFunc("/submit-report", submitReportHandler)

	// Path to your SSL certificate and key
	certPath := "./permsHttps/cert.pem"
	keyPath := "./permsHttps/key.pem"

	// Start the server
	log.Println("Hello there !")
	log.Println("Server started on https://localhost:443/")
	log.Println("Press Ctrl+C to stop the server")

	// Start the server with TLS
	err := http.ListenAndServeTLS(":443", certPath, keyPath, nil)
	if err != nil {
		log.Fatalf("ListenAndServeTLS: %v", err)
	}
}

// routeHandler is the main handler for the server
func routeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request URL path is /api
	path := r.URL.Path
	method := r.Method

	if path == "/" {
		switch method {
		case "GET":
			HomeHandle(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else if strings.HasPrefix(path, "/topic") {
		handleTopic(w, r)
	} else if strings.HasPrefix(path, "/comment") {
		handleComment(w, r)
	} else if strings.HasPrefix(path, "/user") {
		handleUser(w, r)
	} else if strings.HasPrefix(path, "/login") {
		handleLogin(w, r)
	} else if strings.HasPrefix(path, "/addUser") {
		addUser(w, r)
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func handleTopic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case path == "/topic":
		HomeHandle(w, r)
	case strings.HasPrefix(path, "/topic/"):
		handleTopic(w, r)
	case strings.HasPrefix(path, "/comment/"):
		handleComment(w, r)
	case strings.HasPrefix(path, "/user/"):
		handleUser(w, r)
	case strings.HasPrefix(path, "/login"):
		handleLogin(w, r)
	case strings.HasPrefix(path, "/addUser"):
		addUser(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func handleComment(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/comment/")
	switch r.Method {
	case "GET":
		GetComment(w, r, idStr)
	case "POST":
		//CreateComment(w, r)
	case "PUT":
		UpdateComment(w, r, idStr)
	case "DELETE":
		DeleteComment(w, r, idStr)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/user/")
	switch r.Method {
	case "GET":
		GetUser(w, r, idStr)
	case "POST":
		CreateUser(w, r)
	case "PUT":
		UpdateUser(w, r, idStr)
	case "DELETE":
		DeleteUser(w, r, idStr)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		loginUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func AccueilHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/accueil" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles("template/html/accueil.html")
	if err != nil {
		log.Printf("Error parsing template %v", err)
		http.Error(w, "internal server errror ", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
func ContactHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contact" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles("template/html/contact.html") // return to contact
	if err != nil {
		log.Printf("Error parsing template %v", err)
		http.Error(w, "internal server errror ", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
func ProfilHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profil" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles("template/html/profil.html") // return to profil
	if err != nil {
		log.Printf("Error parsing template %v", err)
		http.Error(w, "internal server errror ", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
func PostHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/createpost" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles("template/html/post.html") // Update to the path of your create post template
	if err != nil {
		log.Printf("Error parsing template %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	// Retrieve user_id from cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login page if not authenticated
		return
	}
	userID := cookie.Value

	db := dbConn()
	defer db.Close()

	// Generate the next topic ID by counting existing entries
	var topicID int
	err = db.QueryRow("SELECT COUNT(*) FROM topics").Scan(&topicID)
	if err != nil {
		log.Printf("Error counting topics: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	topicID++ // Increment to get the new topic ID

	// Insert the new topic
	_, err = db.Exec("INSERT INTO topics (id, user_id, title, description) VALUES (?, ?, ?, ?)", topicID, userID, title, description)
	if err != nil {
		log.Printf("Error inserting new topic: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to home or confirmation page
}
