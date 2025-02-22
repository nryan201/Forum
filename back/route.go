package back

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Topic represents a topic
type Topic struct {
	ID          int
	UserID      string
	Title       string
	Description string
	CreatedAt   string
	Categories  []string
	Hashtags    []string
}
type TopicDetail struct {
	ID              int
	Username        string
	Title           string
	Description     string
	IsOwner         bool
	IsAuthenticated bool
	Comments        []Comment
	Role            string
	Categories      []Category
	Hashtags        []Hashtag
}

// Comment represents a comment
type Comment struct {
	ID        int
	TopicID   int
	UserID    int
	Content   string
	CreatedAt time.Time
	Username  string
}

// User represents a user
type User struct {
	ID        string `json:"id"`
	Username  string `json:"Username"`
	Name      string
	Password  string
	Birthday  string
	ProfilImg string
	Email     string
	Role      string
}

type Message struct {
	ID              int    `json:"id"`
	SenderID        string `json:"sender_id"`
	CurrentUsername string
	ReceiverID      string    `json:"receiver_id"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
	Username        string    `json:"Username"`
}

// Category represents a category
type Category struct {
	ID   int
	Name string
}
type Hashtag struct {
	ID   int
	Name string
}
type Data struct {
	Categories []Category
	Hashtags   []Hashtag
	Topics     []Topic
	Username   string
}

type Report struct {
	ID        int
	TopicID   sql.NullInt64
	CommentID sql.NullInt64
	UserID    int
	Reason    string
	CreatedAt string
	Status    string
}
type AdminData struct {
	Users         []User
	Topics        []Topic
	Comments      []Comment
	Categories    []Category
	Hashtags      []Hashtag
	Reports       []Report
	CurrentUserID string
}

var db *sql.DB

// OpenDB initializes the database connection
func OpenDB() {
	var err error
	db, err = sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
}

// HomeHandle handles the home page
func HomeHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("template/html/connexion.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

// Handle for topic
func BlockHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/block" {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, "https://www.minecraft.net/fr-fr", http.StatusFound)

}
