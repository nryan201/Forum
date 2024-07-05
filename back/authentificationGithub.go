package back

import (
	"context"
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		ID    int    `json:"id"` // Changed from string to int for correct JSON unmarshalling
		Email string `json:"email"`
		Name  string `json:"name"`
		Login string `json:"login"`
	}

	if err := json.NewDecoder(response.Body).Decode(&githubUser); err != nil {
		log.Printf("Could not decode response: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	githubUserID := strconv.Itoa(githubUser.ID) // Convert ID from int to string

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

	db := dbConn()
	defer db.Close()

	log.Println("Database opened successfully")

	var userID, username, birthday sql.NullString
	err = db.QueryRow("SELECT id, username, birthday FROM users WHERE id = ?", githubUserID).Scan(&userID, &username, &birthday)
	if err == sql.ErrNoRows {
		// Insert new user
		statement, err := db.Prepare("INSERT INTO users (id, username, email, role) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Fatal("Failed to prepare statement: ", err)
		}
		defer statement.Close()

		_, err = statement.Exec(githubUserID, githubUser.Name, githubUser.Email, "user")
		if err != nil {
			log.Printf("Failed to insert new user: %s", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		log.Println("New user inserted successfully")
		http.SetCookie(w, &http.Cookie{
			Name:     "github_id",
			Value:    githubUserID,
			Path:     "/",
			HttpOnly: true,
		})
		http.Redirect(w, r, "/completeProfileGithub", http.StatusSeeOther)
	} else if err != nil {
		log.Fatal("Failed to query existing user: ", err)
	} else {
		log.Printf("User found with ID: %s", userID.String)
		if !username.Valid || !birthday.Valid {
			http.SetCookie(w, &http.Cookie{
				Name:     "github_id",
				Value:    githubUserID,
				Path:     "/",
				HttpOnly: true,
			})
			http.Redirect(w, r, "/completeProfileGithub", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func completeProfileGithub(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		tmplCompleteProfile := template.Must(template.ParseFiles("./template/html/missingDataGithub.html"))
		tmplCompleteProfile.Execute(w, nil)
		return
	}

	name := r.FormValue("name")
	birthday := r.FormValue("birthday")

	_, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		http.Error(w, "Format de date incorrect", http.StatusBadRequest)
		return
	}

	cookieGithub, err := r.Cookie("github_id")
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du cookie GitHub", http.StatusInternalServerError)
		return
	}

	userID := cookieGithub.Value
	log.Printf("User ID: %s\n", userID)
	db := dbConn()
	defer db.Close()

	_, err = db.Exec("UPDATE users SET name = ?, birthday = ? WHERE id = ?", name, birthday, userID)

	if err != nil {
		log.Printf("Failed to update user data: %s", err)
		http.Error(w, "Erreur lors de la mise à jour des données utilisateur", http.StatusInternalServerError)
		return
	}

	// Supprimer le cookie github_id
	http.SetCookie(w, &http.Cookie{
		Name:    "github_id",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	})

	// Définir le cookie avec l'id de l'utilisateur
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    userID,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
