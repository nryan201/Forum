package back

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

var tmpl = template.Must(template.ParseFiles("./template/html/connexion.html"))

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

		// Chiffrer le mot de passe
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		db := dbConn()
		defer db.Close()

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		var existingUsername, existingEmail string

		// Vérifier si le nom d'utilisateur est déjà pris
		err = tx.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existingUsername)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			log.Fatal(err)
		}
		if existingUsername != "" {
			tx.Rollback()
			fmt.Fprintf(w, "Nom d'utilisateur déjà pris. Veuillez en choisir un autre.")
			return
		}

		// Vérifier si l'adresse e-mail est déjà prise
		err = tx.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existingEmail)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			log.Fatal(err)
		}
		if existingEmail != "" {
			tx.Rollback()
			fmt.Fprintf(w, "Adresse e-mail déjà prise. Veuillez en choisir une autre.")
			return
		}

		// Générer un identifiant séquentiel de type texte
		var maxID sql.NullString
		err = tx.QueryRow("SELECT MAX(CAST(id AS INTEGER)) FROM users").Scan(&maxID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		newID := "1"
		if maxID.Valid {
			maxIDInt, _ := strconv.Atoi(maxID.String)
			newID = strconv.Itoa(maxIDInt + 1)
		}

		// Insérer le nouvel utilisateur avec le rôle par défaut "user"
		_, err = tx.Exec("INSERT INTO users(id, username, password, email, role) VALUES(?, ?, ?, ?, ?)", newID, username, hashedPassword, email, "user")
		if err != nil {
			tx.Rollback()
			fmt.Fprintf(w, "Erreur lors de l'ajout de l'utilisateur : %v", err)
			return
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		fmt.Fprintf(w, "Utilisateur ajouté avec succès")
	} else {
		tmpl.Execute(w, nil)
	}
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		db := dbConn()
		defer db.Close()

		var dbUsername, dbPassword string
		err := db.QueryRow("SELECT username, password FROM users WHERE username = ?", username).Scan(&dbUsername, &dbPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
			} else {
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				log.Fatal(err)
			}
			return
		}

		// Vérifier le mot de passe
		err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
		if err != nil {
			http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: username,
			Path: "/",
			HttpOnly: true,
	})

		fmt.Fprintf(w, "Connexion réussie. Bienvenue %s!", username)
	} else {
		tmpl.Execute(w, nil)
	}
}


