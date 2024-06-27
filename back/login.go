package back

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var tmpl = template.Must(template.ParseFiles("./template/html/connexion1.html"))

func dbConn() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		db := dbConn()
		defer db.Close()

		insForm, err := db.Prepare("INSERT INTO users(username, password, email) VALUES(?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		_, err = insForm.Exec(username, password, email)
		if err != nil {
			fmt.Fprintf(w, "Erreur lors de l'ajout de l'utilisateur : %v", err)
			return
		}
		fmt.Fprintf(w, "Utilisateur ajouté avec succès")
	} else {
		tmpl.Execute(w, nil)
	}
}
