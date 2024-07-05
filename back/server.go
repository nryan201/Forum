package back

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func loadPrivateKey(path string) (*tls.Certificate, error) {
	// Read the certificate file
	certPEM, err := ioutil.ReadFile("./permsHttps/cert.pem")
	if err != nil {
		return nil, err
	}

	// Read the private key file
	keyPEM, err := ioutil.ReadFile("./permsHttps/key.pem")
	if err != nil {
		return nil, err
	}

	// Decode the PEM blocks
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, errors.New("failed to parse certificate PEM")
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, errors.New("failed to parse key PEM")
	}

	// Parse the certificate and key
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// Create the tls.Certificate
	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
	}

	return &tlsCert, nil
}

func Server() {
	OpenDB()
	defer db.Close()

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("template/css/"))))
	http.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("template/ressource/image/"))))
	http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("template/script/"))))

	http.HandleFunc("/logout", logout)
	http.HandleFunc("/block", BlockHandler)
	http.HandleFunc("/login", loginUser)
	http.HandleFunc("/addUser", addUser)
	http.HandleFunc("/loginFacebook", handleFacebookLogin)
	http.HandleFunc("/callbackFacebook", handleFacebookCallback)
	http.HandleFunc("/loginGoogle", handleGoogleLogin)
	http.HandleFunc("/callbackGoogle", handleGoogleCallback)
	http.HandleFunc("/loginGithub", handleGithubLogin)
	http.HandleFunc("/callbackGithub", handleGithubCallback)
	http.HandleFunc("/accueil", AccueilHandle)
	http.HandleFunc("/contact", ContactHandle)
	http.HandleFunc("/profil", profilePage)
	http.HandleFunc("/completeProfile", completeProfile)
	http.HandleFunc("/post", PostHandle)

	// Route pour effacer le cookie et rediriger vers l'accueil
	http.HandleFunc("/clearCookies", func(w http.ResponseWriter, r *http.Request) {
		clearCookie(w, "username")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Route initiale pour gérer la page d'accueil et les autres routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == "GET" {
			clearCookie(w, "username")
			AccueilHandle(w, r)
		} else {
			routeHandler(w, r)
		}
	})

	// Load the certificate and private key
	certPath := "./permsHttps/cert.pem"
	keyPath := "./permsHttps/key.pem"
	tconf := &tls.Config{}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}
	tconf.Certificates = append(tconf.Certificates, cert)

	// Lancer un serveur HTTP sur le port 8080
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Configuration du serveur HTTPS
	server := &http.Server{
		Addr:      ":443",
		Handler:   http.DefaultServeMux, // Utiliser le DefaultServeMux pour gérer les routes
		TLSConfig: tconf,
	}

	// Start the server
	log.Println("Hello there !")
	log.Println("Server started on https://localhost:443/")
	log.Println("Press Ctrl+C to stop the server")

	err = server.ListenAndServeTLS("", "")
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
			AccueilHandle(w, r)
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
		AccueilHandle(w, r)
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
		CreateComment(w, r)
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
	if r.URL.Path != "/" && r.URL.Path != "/accueil" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles("template/html/accueil.html") // return to accueil
	if err != nil {
		log.Printf("Error parsing template %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
