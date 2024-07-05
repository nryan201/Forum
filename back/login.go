package back

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var tmpl = template.Must(template.ParseFiles("./template/html/connexion.html"))
var tmplProfile = template.Must(template.ParseFiles("./template/html/profil.html"))

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
		name := r.FormValue("name")
		birthday := r.FormValue("birthday")

		log.Printf("Username: %s\n", username)
		log.Printf("Password: %s\n", password)
		log.Printf("Email: %s\n", email)
		log.Printf("Name: %s\n", name)
		log.Printf("Date de naissance (originale): %s\n", birthday)

		// Vérifier le format de la date
		_, err := time.Parse("2006-01-02", birthday)
		if err != nil {
			http.Error(w, "Format de date incorrect", http.StatusBadRequest)
			return
		}

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

		var existingUsername string

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

		// Générer un nouvel ID unique
		var count int
		err = tx.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		newID := strconv.Itoa(count + 1)

		// Insérer le nouvel utilisateur avec le rôle par défaut "user"
		_, err = tx.Exec("INSERT INTO users(id, username, name, birthday, password, email, role) VALUES(?, ?, ?, ?, ?, ?, ?)", newID, username, name, birthday, hashedPassword, email, "user")
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

		// Créer un cookie de session pour l'utilisateur nouvellement ajouté avec user_id
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    newID,
			Path:     "/",
			HttpOnly: true,
		})

		// Rediriger vers la page d'accueil
		log.Printf("Nouvel utilisateur ajouté : %s\n", username)
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
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

		var userID, dbPassword string
		err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&userID, &dbPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
			} else {
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				log.Println("Error querying database:", err)
			}
			return
		}

		// Vérifier le mot de passe
		err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
		if err != nil {
			http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// Set user_id cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    userID,
			Path:     "/",
			HttpOnly: true,
		})
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
	} else {
		tmpl.Execute(w, nil)
	}
}

func profilePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering profilePage function")

	// Lire le cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		if err == http.ErrNoCookie {
			// Pas de cookie, utilisateur non connecté
			log.Println("No cookie found, redirecting to login page")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Autre erreur
		log.Println("Error reading cookie:", err)
		http.Error(w, "Erreur lors de la lecture du cookie", http.StatusInternalServerError)
		return
	}

	userID := cookie.Value
	log.Println("Cookie found, user_id:", userID)

	// Récupérer les informations de l'utilisateur depuis la base de données
	db := dbConn()
	defer db.Close()

	var user struct {
		Username  string
		Name      string
		Firstname string
		Birthdate string
		Email     string
	}
	user.Username = "-"
	user.Name = "-"
	user.Firstname = "-"
	user.Birthdate = "-"
	user.Email = "-"

	err = db.QueryRow("SELECT username, name, strftime('%Y-%m-%d', birthday), email FROM users WHERE id = ?", userID).Scan(&user.Username, &user.Name, &user.Birthdate, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found with id:", userID)
			http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
			return
		}
		log.Println("Error retrieving user data:", err)
		http.Error(w, "Erreur lors de la récupération des données utilisateur", http.StatusInternalServerError)
		return
	}

	log.Printf("User data retrieved: %+v\n", user)

	// Rendre le template avec les données de l'utilisateur
	tmplProfile := template.Must(template.ParseFiles("./template/html/profil.html"))
	err = tmplProfile.Execute(w, user)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Erreur lors du rendu du template", http.StatusInternalServerError)
	}
}

func clearCookie(w http.ResponseWriter, name string) {
	log.Println("Clearing cookie", name)
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	})
}
func logout(w http.ResponseWriter, r *http.Request) {
	log.Println("Logging out user")
	clearCookie(w, "username")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
