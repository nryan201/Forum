package back

import (
	"log"
	"net/http"
	"html/template"

)

func HomeHandle (w http.ResponseWriter, r *http.Request) {
	
	tmp, err := template.ParseFiles("template/html/accueil.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)	
		return
	}

	err = tmp.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
