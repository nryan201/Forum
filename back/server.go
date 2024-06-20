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
	router.PathPrefix("/script/").Handler(http.StripPrefix("/script/", http.FileServer(http.Dir("./template/script"))))
    router.PathPrefix("/image/").Handler(http.StripPrefix("/image/", http.FileServer(http.Dir("./template/ressource/image"))))

    // Handle the routes
    router.HandleFunc("/", HomeHandle).Methods("Get") 

    // Handle the topic page
    router.HandleFunc("/topic", CreateTopic).Methods("POST") 
    router.HandleFunc("/topic/{id}", GetTopic).Methods("GET")
    router.HandleFunc("/topic/{id}", UpdateTopic).Methods("PUT") 
    router.HandleFunc("/topic/{id}", DeleteTopic).Methods("DELETE") 

    // Handle the comment page
    router.HandleFunc("/comment", CreateComment).Methods("POST")
    router.HandleFunc("/comment/{id}", GetComment).Methods("GET")
    router.HandleFunc("/comment/{id}", UpdateComment).Methods("PUT")
    router.HandleFunc("/comment/{id}", DeleteComment).Methods("DELETE")

    // Handle the user page
    router.HandleFunc("/user", CreateUser).Methods("POST") 
    router.HandleFunc("/user/{id}", GetUser).Methods("GET") 
    router.HandleFunc("/user/{id}", UpdateUser).Methods("PUT") 
    router.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE") 

    // Authentication routes
    router.HandleFunc("/login", Login).Methods("POST") 
    router.HandleFunc("/logout", Logout).Methods("POST") 



    // Start the server
    log.Println("Hello there !")
    log.Println("Server started on http://localhost:8080/")
    log.Println("Press Ctrl+C to stop the server")

    err := http.ListenAndServe(":8080", router)
    if err != nil {
        log.Fatalf("Could not start the server: %v", err)
    }


}