package back

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  "https://localhost:443/callbackGithub",
		ClientID:     "Ov23ctETeI5Rk1nEsbPa",
		ClientSecret: "b12f971f88aa12cb9d6e57f4d24e101eeddc1dac",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	githubOauthStateString = "randomstringGithub"
)

func handleGithubLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL(githubOauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGithubCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != githubOauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := githubOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Printf("Could not get token: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := githubOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Printf("Could not create request: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()

	var githubUser struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Login string `json:"login"`
	}

	if err := json.NewDecoder(response.Body).Decode(&githubUser); err != nil {
		log.Printf("Could not decode response: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if githubUser.Email == "" {
		emailsResponse, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			log.Printf("Could not fetch emails: %s\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer emailsResponse.Body.Close()

		var emails []struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
		}
		if err := json.NewDecoder(emailsResponse.Body).Decode(&emails); err != nil {
			log.Printf("Could not decode emails response: %s\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		for _, email := range emails {
			if email.Primary {
				githubUser.Email = email.Email
				break
			}
		}
	}

	if githubUser.Name == "" {
		githubUser.Name = githubUser.Login
	}

	githubUser.Name = strings.TrimSpace(githubUser.Name)
	githubUser.Email = strings.TrimSpace(githubUser.Email)
	githubUser.Login = strings.TrimSpace(githubUser.Login)

	if !isValidString(githubUser.Name) || !isValidString(githubUser.Email) || !isValidString(githubUser.Login) {
		log.Printf("Invalid data found in user details: Name: %s, Email: %s, Login: %s\n", githubUser.Name, githubUser.Email, githubUser.Login)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	log.Println("Database opened successfully")

	var userID string
	err = db.QueryRow("SELECT id FROM users WHERE id = ?", fmt.Sprint(githubUser.ID)).Scan(&userID)
	if err == sql.ErrNoRows {
		log.Printf("Attempting to insert new user with ID: %d, Name: %s, Email: %s", githubUser.ID, githubUser.Name, githubUser.Email)
		statement, err := db.Prepare("INSERT INTO users (id, username, email, role) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Fatal("Failed to prepare statement: ", err)
		}
		defer statement.Close()

		log.Printf("Inserting user with ID (type: %T): %d", githubUser.ID, githubUser.ID)
		log.Printf("Inserting user with Name (type: %T): %s", githubUser.Name, githubUser.Name)
		log.Printf("Inserting user with Email (type: %T): %s", githubUser.Email, githubUser.Email)

		_, err = statement.Exec(fmt.Sprint(githubUser.ID), githubUser.Name, githubUser.Email, "user")
		if err != nil {
			log.Printf("Failed to insert new user: %s", err)
		} else {
			log.Println("New user inserted successfully")
			log.Printf("User ID: %d", githubUser.ID)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else if err != nil {
		log.Fatal("Failed to query existing user: ", err)
	} else {
		log.Printf("User found with ID: %s", userID)
		log.Printf("Welcome back %s!", githubUser.Name)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
