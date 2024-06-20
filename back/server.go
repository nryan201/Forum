package back

import (
	"log"
	"net/http"
)

func Server() {
	// Create a file server to serve static files
	htmlFs := http.FileServer(http.Dir("./template/html"))
	http.Handle("/html/", http.StripPrefix("/html/", htmlFs))

	cssFs := http.FileServer(http.Dir("./template/css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFs))

	jsFs := http.FileServer(http.Dir("./template/script"))
	http.Handle("/script/", http.StripPrefix("/script/", jsFs))

	// Handle the routes

	http.HandleFunc("/", HomeHandle) // Handle the home page

	// Start the server
	log.Println("Hello there !")
	log.Println("Server started on http://localhost:8080/")
	log.Println("Press Ctrl+C to stop the server")

	err := http.ListenAndServe(":8080", nil) // Start the server
	if err != nil {
		log.Fatalf("Could not start the server: %v", err)
	}

}
