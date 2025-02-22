package back

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

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

		// Générer un nouvel UUID unique
		newID := uuid.New().String()

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
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    newID,
			Path:     "/",
			HttpOnly: true,
		})

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
				log.Println("Erreur lors de la vérification des informations de connexion:", err)
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				log.Println("Error querying database:", err)
			}
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Set user_id cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    userID,
			Path:     "/",
			HttpOnly: false,
			Secure:   true,
		})
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
	} else {
		tmpl.Execute(w, nil)
	}
}

func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie("user_id")
	if err != nil {
		if err == http.ErrNoCookie {
			json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
			return
		}

		log.Println("Error retrieving user_id cookie:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur lors de la récupération du cookie"})
		return
	}

	authenticated := cookie.Value != ""
	log.Printf("User authenticated: %v", authenticated)
	json.NewEncoder(w).Encode(map[string]bool{"authenticated": authenticated})
}

func profilePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering profilePage function")

	cookie, err := r.Cookie("user_id")
	if err != nil {
		log.Printf("Error reading cookie or no cookie found: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID := cookie.Value
	log.Println("Cookie found, user_id:", userID)

	db := dbConn()
	defer db.Close()

	var user struct {
		Username  string
		Name      string
		Firstname string
		Birthdate string
		Email     string
		ProfilImg string
		Role      string
		IsAdmin   bool
	}
	user.Username = "-"
	user.Name = "-"
	user.Firstname = "-"
	user.Birthdate = "-"
	user.Email = "-"
	user.Role = "-"
	user.IsAdmin = false
	user.ProfilImg = "/image/fondblanc.png"

	var profileImg sql.NullString

	err = db.QueryRow("SELECT username, name, strftime('%Y-%m-%d', birthday), profile_image, email, role FROM users WHERE id = ?", userID).Scan(&user.Username, &user.Name, &user.Birthdate, &profileImg, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found with id:", userID)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		log.Println("Error retrieving user data:", err)
		http.Error(w, "Erreur lors de la récupération des données utilisateur", http.StatusInternalServerError)
		return
	}

	// Gérer le cas où profile_image est NULL
	if profileImg.Valid {
		user.ProfilImg = profileImg.String
	} else {
		user.ProfilImg = "/image/fondblanc.png"
	}

	// Vérifier si l'utilisateur est un administrateur
	if user.Role == "admin" {
		user.IsAdmin = true
	}

	log.Printf("User data retrieved: %+v\n", user)

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
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

func ClearAllCookies(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, "user_id")
	clearCookie(w, "username")
	clearCookie(w, "role")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request) {

	ClearAllCookies(w, r)
	w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	w.Write([]byte("Logged out"))

}
