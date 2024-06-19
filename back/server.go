package back

import (
	"log"
	"net/http"
)

func Server(){
    // Create a file server to serve static files
	htmlFs := http.FileServer(http.Dir("./web/html")) 
	http.Handle("/html/",http.StripPrefix("/html/",htmlFs))

    cssFs := http.FileServer(http.Dir("./web/css")) 
    http.Handle("/css/",http.StripPrefix("/css/",cssFs)) 

    jsFs := http.FileServer(http.Dir("./web/js")) 
    http.Handle("/js/",http.StripPrefix("/js/",jsFs)) 

    // Handle the routes

    http.HandleFunc("/", HomeHandle) // Handle the home page

    // Handle the topic page
    http.HandleFunc("/topic", CreateTopic) // Handle the topic page
    http.HandleFunc("/topic/{id}", GetTopic) // Handle the view topic page
    http.HandleFunc("/topic/{id}", UpdateTopic) // Handle the update topic page
    http.HandleFunc("/topic/{id}", DeleteTopic) // Handle the delete topic page

    // Handle the comment page
    http.HandleFunc("/comment", CreateComment) // Handle the comment page
    http.HandleFunc("/comment/{id}", GetComment) // Handle the view comment page
    http.HandleFunc("/comment/{id}", UpdateComment) // Handle the update comment page
    http.HandleFunc("/comment/{id}", DeleteComment) // Handle the delete comment page

    // Handle the user page
    http.HandleFunc("/user", CreateUser) // Handle the user page
    http.HandleFunc("/user/{id}", GetUser) // Handle the view user page
    http.HandleFunc("/user/{id}", UpdateUser) // Handle the update user page
    http.HandleFunc("/user/{id}", DeleteUser) // Handle the delete user page

    // Handle the login page
    http.HandleFunc("/login", Login) // Handle the login page

    // Handle the logout page
    http.HandleFunc("/logout", Logout) // Handle the logout page
    


    // Start the server
    log.Println("Hello there !")
    log.Println("Server started on http://localhost:8080/")
    log.Println("Press Ctrl+C to stop the server")

    err := http.ListenAndServe(":8080", nil) // Start the server
    if err != nil {
        log.Fatalf("Could not start the server: %v", err)
    }


}