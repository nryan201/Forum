package back

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "https://localhost:443/callbackGoogle",
		ClientID:     "767814867927-m0ou233go88bec08h0qobi3v50nn5qhg.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	googleOauthStateString = "randomstringGoogle"
)

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(googleOauthStateString, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != googleOauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Printf("Could not get token: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Printf("Failed to request user info: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()

	var GoogleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(response.Body).Decode(&GoogleUser); err != nil {
		log.Printf("Could not decode user data: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	GoogleUser.ID = strings.TrimSpace(GoogleUser.ID)
	GoogleUser.Name = strings.TrimSpace(GoogleUser.Name)
	GoogleUser.Email = strings.TrimSpace(GoogleUser.Email)

	if !isValidString(GoogleUser.ID) || !isValidString(GoogleUser.Name) || !isValidString(GoogleUser.Email) {
		log.Printf("Invalid data found in user details: ID: %s, Name: %s, Email: %s\n", GoogleUser.ID, GoogleUser.Name, GoogleUser.Email)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	db := dbConn()
	defer db.Close()

	log.Println("Database opened successfully")

	var userGoogleID, username, birthday sql.NullString
	err = db.QueryRow("SELECT id, username, birthday FROM users WHERE id = ?", GoogleUser.ID).Scan(&userGoogleID, &username, &birthday)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Attempting to insert new user with ID: %s, Name: %s, Email: %s", GoogleUser.ID, GoogleUser.Name, GoogleUser.Email)
		statement, err := db.Prepare("INSERT INTO users (id, email, name, role) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Fatal("Failed to prepare statement: ", err)
		}
		defer statement.Close()

		_, err = statement.Exec(GoogleUser.ID, GoogleUser.Email, GoogleUser.Name, "user")
		if err != nil {
			log.Printf("Failed to insert new user: %s", err)
		} else {
			log.Println("New user inserted successfully")
		}

		// Set cookie with the Google ID
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    GoogleUser.ID,
			Path:     "/",
			HttpOnly: true,
		})

		// Redirect to the complete profile page
		http.Redirect(w, r, "/completeProfile", http.StatusSeeOther)
	} else if err != nil {
		log.Fatal("Failed to query existing user: ", err)
	} else {
		log.Printf("User found with ID: %s", userGoogleID.String)
		if !username.Valid || !birthday.Valid {
			// Set cookie with the Google ID
			http.SetCookie(w, &http.Cookie{
				Name:     "user_id",
				Value:    GoogleUser.ID,
				Path:     "/",
				HttpOnly: true,
			})

			http.Redirect(w, r, "/completeProfile", http.StatusSeeOther)
		} else {
			// Set the cookie with the user ID
			http.SetCookie(w, &http.Cookie{
				Name:     "user_id",
				Value:    userGoogleID.String,
				Path:     "/",
				HttpOnly: true,
			})

			http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		}
	}
}

func completeProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		tmplCompleteProfile := template.Must(template.ParseFiles("./template/html/missingData.html"))
		tmplCompleteProfile.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	birthday := r.FormValue("birthday")

	_, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		http.Error(w, "Format de date incorrect", http.StatusBadRequest)
		return
	}

	var userID string
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du cookie", http.StatusInternalServerError)
		return
	}
	userID = cookie.Value

	db := dbConn()
	defer db.Close()

	var existingUsername string
	err = db.QueryRow("SELECT username FROM users WHERE username = ? AND id != ?", username, userID).Scan(&existingUsername)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking username: %s", err)
		http.Error(w, "Erreur lors de la vérification du nom d'utilisateur", http.StatusInternalServerError)
		return
	}
	if existingUsername != "" {
		http.Error(w, "Nom d'utilisateur déjà pris. Veuillez en choisir un autre.", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("UPDATE users SET username = ?, birthday = ? WHERE id = ?", username, birthday, userID)
	if err != nil {
		log.Printf("Failed to update user data: %s", err)
		http.Error(w, "Erreur lors de la mise à jour des données utilisateur", http.StatusInternalServerError)
		return
	}

	// Set the cookie with the updated user ID
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    userID,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
}

func isValidString(str string) bool {
	for _, r := range str {
		if r == '\uFFFD' {
			return false
		}
	}
	return true
}
