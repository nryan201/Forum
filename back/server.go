package back

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Server() {
	r := mux.NewRouter()
	db := dbConn()
	defer db.Close()
	// Create a file server to serve static files
	r.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("./template/html"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./template/css"))))
	r.PathPrefix("/script/").Handler(http.StripPrefix("/script/", http.FileServer(http.Dir("./template/script"))))
	r.PathPrefix("/image/").Handler(http.StripPrefix("/image/", http.FileServer(http.Dir("./template/ressource/image"))))

	// Handle the routes
	r.HandleFunc("/", HomeHandle).Methods("Get")
	// Handle the topic page
	r.HandleFunc("/topic", CreateTopic).Methods("POST")
	r.HandleFunc("/topic/{id}", GetTopic).Methods("GET")
	r.HandleFunc("/topic/{id}", UpdateTopic).Methods("PUT")
	r.HandleFunc("/topic/{id}", DeleteTopic).Methods("DELETE")

	// Handle the comment page
	r.HandleFunc("/comment", CreateComment).Methods("POST")
	r.HandleFunc("/comment/{id}", GetComment).Methods("GET")
	r.HandleFunc("/comment/{id}", UpdateComment).Methods("PUT")
	r.HandleFunc("/comment/{id}", DeleteComment).Methods("DELETE")

	// Handle the user page
	r.HandleFunc("/user", CreateUser).Methods("POST")
	r.HandleFunc("/user/{id}", GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE")

	// Handle for Catergory
	r.HandleFunc("/category", CreateCategory).Methods("POST")
	r.HandleFunc("/category/{id}", CategoryHandler).Methods("GET")
	r.HandleFunc("/category/{id}", DeleteCategory).Methods("DELETE")

	/* non fonctionel pour le moment
	r.HandleFunc("/adduser", addUser).Methods("GET", "POST")
	r.HandleFunc("/login", loginUser).Methods("GET", "POST")
	*/

	// Start the server
	log.Println("Hello there !")
	log.Println("Server started on http://localhost:8080/")
	log.Println("Press Ctrl+C to stop the server")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Could not start the server: %v", err)
	}

}
