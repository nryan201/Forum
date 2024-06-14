package back

import (
	"fmt"
	"net/http"
)

func Server() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %q", r.URL.Path)
    })

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    fmt.Println("Server started on http://localhost:8080/")
}