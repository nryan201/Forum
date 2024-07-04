package back

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
)

var (
	discordConfig    *oauth2.Config
	oauthStateString = "random"
)

type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	discordConfig = &oauth2.Config{
		ClientID:     "1258399519654023301",
		ClientSecret: "CfOJqTTak9wL8fCI5Ary2oRVY79isDsA",
		RedirectURL:  "http://localhost:8080/auth/discord/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		Scopes: []string{"identify", "email"},
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html><body><a href="/auth/discord">Login with Discord</a></body></html>`
	fmt.Fprint(w, htmlIndex)
}

func handleDiscordLogin(w http.ResponseWriter, r *http.Request) {
	url := discordConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleDiscordCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := discordConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Println("Code exchange failed:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := discordConfig.Client(context.Background(), token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		log.Println("Failed to get user info:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var user DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Println("Failed to decode user info:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Insert user into database
	_, err = db.Exec(`
		INSERT OR REPLACE INTO users (id, username, password, email, role) VALUES (?, ?, ?, ?, ?)`,
		user.ID, user.Username+"#"+user.Discriminator, "", user.Email, "user")
	if err != nil {
		log.Println("Failed to insert user into database:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Display user information
	fmt.Fprintf(w, "User Info: %+v\n", user)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleMain)
	r.HandleFunc("/auth/discord", handleDiscordLogin)
	r.HandleFunc("/auth/discord/callback", handleDiscordCallback)

	fmt.Println("Started running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
