package back

import (
	"database/sql"
	"html/template"
	"log"
)

var tmpl = template.Must(template.ParseFiles("./template/html/connexion.html"))

func dbConn() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

/* Non fonctionel pour le moment attendre que Ryan et Alexis modifient ceci
func addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

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

		// Insérer le nouvel utilisateur
		_, err = tx.Exec("INSERT INTO users(username, password, email) VALUES(?, ?, ?)", username, password, email)
		if err != nil {
			tx.Rollback()
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				fmt.Fprintf(w, "Nom d'utilisateur ou adresse e-mail déjà pris. Veuillez en choisir un autre.")
				return
			}
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
				fmt.Fprintf(w, "Nom d'utilisateur ou mot de passe incorrect")
			} else {
				log.Fatal(err)
			}
			return
		}

		if dbPassword != password {
			fmt.Fprintf(w, "Nom d'utilisateur ou mot de passe incorrect")
			return
		}

		fmt.Fprintf(w, "Connexion réussie. Bienvenue %s!", username)
	} else {
		tmpl.Execute(w, nil)
	}
}
*/
