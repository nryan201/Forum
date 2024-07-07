package back

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var (
	facebookOauthConfig = &oauth2.Config{
		RedirectURL:  "https://localhost:443/callbackFacebook",
		ClientID:     "1515434692706809",
		ClientSecret: "3d109972564526228111d7206c651932",
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
	facebookOauthStateString = "randomstringFacebook"
)

func handleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	url := facebookOauthConfig.AuthCodeURL(facebookOauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != facebookOauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := facebookOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Printf("Could not get token: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + token.AccessToken)
	if err != nil {
		log.Printf("Could not create request: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()

	var facebookUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(response.Body).Decode(&facebookUser); err != nil {
		log.Printf("Could not decode response: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	facebookUser.ID = strings.TrimSpace(facebookUser.ID)
	facebookUser.Name = strings.TrimSpace(facebookUser.Name)
	facebookUser.Email = strings.TrimSpace(facebookUser.Email)

	if !isValidString(facebookUser.ID) || !isValidString(facebookUser.Name) || !isValidString(facebookUser.Email) {
		log.Printf("Invalid data found in user details: ID: %s, Name: %s, Email: %s\n", facebookUser.ID, facebookUser.Name, facebookUser.Email)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	db := dbConn()
	defer db.Close()

	var userID, username, birthday sql.NullString
	err = db.QueryRow("SELECT id, username, birthday FROM users WHERE id = ?", facebookUser.ID).Scan(&userID, &username, &birthday)
	if err == sql.ErrNoRows {
		// Insert new user
		statement, err := db.Prepare("INSERT INTO users (id, name, email, role) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Fatal("Failed to prepare statement: ", err)
		}
		defer statement.Close()

		_, err = statement.Exec(facebookUser.ID, facebookUser.Name, facebookUser.Email, "user")
		if err != nil {
			log.Printf("Failed to insert new user: %s", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		log.Println("New user inserted successfully")
		// Set cookie with the Facebook ID
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    facebookUser.ID,
			Path:     "/",
			HttpOnly: true,
		})
		// Redirect to the complete profile page only if user data is incomplete
		http.Redirect(w, r, "/completeProfile", http.StatusSeeOther)
	} else if err != nil {
		log.Fatal("Failed to query existing user: ", err)
	} else {
		log.Printf("User found with ID: %s, Username Valid: %t, Birthday Valid: %t", userID.String, username.Valid, birthday.Valid)
		if !username.Valid || !birthday.Valid {
			// Set cookie with the Facebook ID
			http.SetCookie(w, &http.Cookie{
				Name:     "user_id",
				Value:    facebookUser.ID,
				Path:     "/",
				HttpOnly: true,
			})
			http.Redirect(w, r, "/completeProfile", http.StatusSeeOther)
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     "user_id",
				Value:    facebookUser.ID,
				Path:     "/",
				HttpOnly: true,
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}
