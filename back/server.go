package back

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Server(){
    router := mux.NewRouter()

    // Create a file server to serve static files
	router.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("./template/html"))))
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./template/css"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./template/script"))))
    router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./template/ressource/image/img"))))

    // Handle the routes

    router.HandleFunc("/", HomeHandle).Methods("Get") // Handle the home page

    // Handle the topic page
   /* router.HandleFunc("/topic", CreateTopic).Methods("POST") // Handle the topic page
    router.HandleFunc("/topic/{id}", GetTopic).Methods("GET") // Handle the view topic page
    router.HandleFunc("/topic/{id}", UpdateTopic).Methods("PUT") // Handle the update topic page
    router.HandleFunc("/topic/{id}", DeleteTopic).Methods("DELETE") // Handle the delete topic page

    // Handle the comment page
    router.HandleFunc("/comment", CreateComment).Methods("POST") // Handle the comment page
    router.HandleFunc("/comment/{id}", GetComment).Methods("GET") // Handle the view comment page
    router.HandleFunc("/comment/{id}", UpdateComment).Methods("PUT") // Handle the update comment page
    router.HandleFunc("/comment/{id}", DeleteComment).Methods("DELETE") // Handle the delete comment page

    // Handle the user page
    router.HandleFunc("/user", CreateUser).Methods("POST") // Handle the user page
    router.HandleFunc("/user/{id}", GetUser).Methods("GET") // Handle the view user page
    router.HandleFunc("/user/{id}", UpdateUser).Methods("PUT") // Handle the update user page
    router.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE") // Handle the delete user page

    // Authentication routes
    router.HandleFunc("/login", Login).Methods("POST") 
    router.HandleFunc("/logout", Logout).Methods("POST") */



    // Start the server
    log.Println("Hello there !")
    log.Println("Server started on http://localhost:8080/")
    log.Println("Press Ctrl+C to stop the server")

    err := http.ListenAndServe(":8080", router)
    if err != nil {
        log.Fatalf("Could not start the server: %v", err)
    }


}