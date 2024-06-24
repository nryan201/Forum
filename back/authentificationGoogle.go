package back

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	url := googleOauthConfig.AuthCodeURL(googleOauthStateString)
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
		log.Printf("Could not create request: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(response.Body).Decode(&googleUser); err != nil {
		log.Printf("Could not decode response: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	db, err := sql.Open("sqlite3", "./mydatabase.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE google_id = ? OR email = ?", googleUser.ID, googleUser.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		// User does not exist, create a new user
		statement, err := db.Prepare("INSERT INTO users (username, email, google_id) VALUES (?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer statement.Close()
		res, err := statement.Exec(googleUser.Name, googleUser.Email, googleUser.ID)
		if err != nil {
			log.Fatal(err)
		}
		userID64, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		userID = int(userID64)
	} else if err != nil {
		log.Fatal(err)
	}

	// Create session for user, etc.
	fmt.Fprintf(w, "User ID: %d\n", userID)
	fmt.Fprintf(w, "User Info: %s\n", googleUser.Name)
}

func main() {
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
