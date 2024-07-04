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
	http.HandleFunc("/loginFacebook", handleFacebookLogin)
	http.HandleFunc("/callbackFacebook", handleFacebookCallback)
	http.HandleFunc("/loginGoogle", handleGoogleLogin)
	http.HandleFunc("/callbackGoogle", handleGoogleCallback)
	http.HandleFunc("/accueil", AccueilHandle)
	http.HandleFunc("/contact", ContactHandle)
	http.HandleFunc("/profil", ProfilHandle)
	http.HandleFunc("/post", PostHandle)

	// Path to your SSL certificate and key
	certPath := "./permsHttps/cert.pem"
	keyPath := "./permsHttps/key.pem"

	log.Println("Hello there !")
	log.Println("Server started on http://localhost:8080/")
	log.Println("Press Ctrl+C to stop the server")

	err := http.ListenAndServeTLS(":8080", certPath, keyPath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: %v", err)
	}
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request URL path is /api
	path := r.URL.Path
	method := r.Method

	if path == "/" {
		switch method {
		case "GET":
			HomeHandle(w, r)
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
	idStr := strings.TrimPrefix(r.URL.Path, "/topic/")
	switch r.Method {
	case "GET":
		GetTopic(w, r, idStr)
		GetTopic(w, r, idStr)
	case "POST":
		CreateTopic(w, r, idStr)
	case "PUT":
		UpdateTopic(w, r, idStr)
	case "DELETE":
		DeleteTopic(w, r, idStr)
		DeleteTopic(w, r, idStr)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleComment(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/comment/")
	switch r.Method {
	case "GET":
		GetComment(w, r, idStr)
		GetComment(w, r, idStr)
	case "POST":
		CreateComment(w, r)
	case "PUT":
		UpdateComment(w, r, idStr)
		UpdateComment(w, r, idStr)
	case "DELETE":
		DeleteComment(w, r, idStr)
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
		GetUser(w, r, idStr)
	case "POST":
		CreateUser(w, r)
	case "PUT":
		UpdateUser(w, r, idStr)
	case "DELETE":
		DeleteUser(w, r, idStr)
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
	tmpl, err := template.ParseFiles("template/html/accueil.html") // return to accueil
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
	if r.URL.Path != "/post" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles("template/html/post.html") // return to post
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
