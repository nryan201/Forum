package back

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callbackGoogle",
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

	// Nettoyer les valeurs pour éliminer les caractères invisibles
	GoogleUser.ID = strings.TrimSpace(GoogleUser.ID)
	GoogleUser.Name = strings.TrimSpace(GoogleUser.Name)
	GoogleUser.Email = strings.TrimSpace(GoogleUser.Email)

	// Validation des données pour s'assurer qu'il n'y a pas de caractères non visibles
	if !isValidString(GoogleUser.ID) || !isValidString(GoogleUser.Name) || !isValidString(GoogleUser.Email) {
		log.Printf("Invalid data found in user details: ID: %s, Name: %s, Email: %s\n", GoogleUser.ID, GoogleUser.Name, GoogleUser.Email)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	log.Println("Database opened successfully")

	var userGoogleID string
	err = db.QueryRow("SELECT id FROM users WHERE id = ? OR email = ?", GoogleUser.ID, GoogleUser.Email).Scan(&userGoogleID)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Attempting to insert new user with ID: %s, Name: %s, Email: %s", GoogleUser.ID, GoogleUser.Name, GoogleUser.Email)
		statement, err := db.Prepare("INSERT INTO users (id, username, email, role) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Fatal("Failed to prepare statement: ", err)
		}
		defer statement.Close()

		// Ajout de journaux pour vérifier les types de données
		log.Printf("Inserting user with ID (type: %T): %s", GoogleUser.ID, GoogleUser.ID)
		log.Printf("Inserting user with Name (type: %T): %s", GoogleUser.Name, GoogleUser.Name)
		log.Printf("Inserting user with Email (type: %T): %s", GoogleUser.Email, GoogleUser.Email)

		_, err = statement.Exec(GoogleUser.ID, GoogleUser.Name, GoogleUser.Email, "user")
		if err != nil {
			log.Printf("Failed to insert new user: %s", err)
		} else {
			log.Println("New user inserted successfully")
			fmt.Fprintf(w, "Connexion réussie. Bienvenue %s! (ID utilisateur: %s)", GoogleUser.Name, GoogleUser.ID)
		}
	} else if err != nil {
		log.Fatal("Failed to query existing user: ", err)
	} else {
		log.Printf("User found with ID: %s", userGoogleID)
		fmt.Fprintf(w, "Welcome back %s! (User ID: %s)", GoogleUser.Name, userGoogleID)
	}
}

func isValidString(str string) bool {
	for _, r := range str {
		if r == '\uFFFD' {
			return false
		}
	}
	return true
}
