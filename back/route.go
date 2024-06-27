package back

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"golang.org/x/crypto/bcrypt"
)

// Topic struct
type Topic struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Comment struct
type Comment struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	TopicID int    `json:"topic_id"`
}

// User struct
type User struct {
	ID       int    `json:"id"`
	Mail     string `json:"mail"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var store = sessions.NewCookieStore([]byte("secret-key"))

// Open data base
func OpenDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "chemin de la bdd")
}

// Handle for home page
func HomeHandle(w http.ResponseWriter, r *http.Request) {

	tmp, err := template.ParseFiles("template/html/accueil.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmp.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

// Handle for  Create topic
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	// Logic for creating a topic
	var topic Topic
	err := json.NewDecoder(r.Body).Decode(&topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO topics (title, description) VALUES (?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(topic.Title, topic.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topic)
}

// GetTopic retrieves a single topic by ID
func GetTopic(w http.ResponseWriter, r *http.Request) {
	// Logic for getting a topic
	vars := mux.Vars(r)
	topicID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var topic Topic
	err = db.QueryRow("SELECT id, title, description FROM topics WHERE id = ?", topicID).Scan(&topic.ID, &topic.Title, &topic.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Topic not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topic)
}

// UpdateTopic updates a single topic by ID from PostgresSql
func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var topic Topic
	err = json.NewDecoder(r.Body).Decode(&topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", 500) // error 500 = internal server error
		return

	}

	defer db.Close()

	statement, err := db.Prepare("UPDATE topics SET title = ?, description = ? WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(topic.Title, topic.Description, topicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// DeleteTopic deletes a single topic by ID from PostgresSql
func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	// Logic for deleting a topic
	vars := mux.Vars(r)
	topicID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", 500) // error 500 = internal server error
		return
	}
	defer db.Close()

	statement, err := db.Prepare("DELETE FROM topics WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(topicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// Handle for create commment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	// Logic for creating a comment
	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", 500) // error 500 = internal server error
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO comments (content, topic_id) VALUES (?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(comment.Content, comment.TopicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetComment(w http.ResponseWriter, r *http.Request) {
	// Logic for getting a comment
	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", 500) // error 500 = internal server error
		return
	}
	defer db.Close()

	statement, err := db.Prepare("SELECT id, content, topic_id FROM comments WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	err = statement.QueryRow(comment.ID).Scan(&comment.ID, &comment.Content, &comment.TopicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	// Logic for updating a comment
	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", 500) // error 500 = internal server error
		return

	}
	defer db.Close()

	statement, err := db.Prepare("UPDATE comments SET content = ?, topic_id = ? WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(comment.Content, comment.TopicID, comment.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	// Logic for deleting a comment
	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", 500) // error 500 = internal server error
		return
	}
	defer db.Close()

	statement, err := db.Prepare("DELETE FROM comments WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(comment.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Handle for user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Logic for creating a user
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := OpenDB()

	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", 500) // error 500 = internal server error
		return
	}

	defer db.Close()

	statement, err := db.Prepare("INSERT INTO users (username, mail, password) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetUser retrieves a single user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Logic for getting a user
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", 500) // error 500 = internal server error
		return
	}
	defer db.Close()

	statement, err := db.Prepare("SELECT id, username, mail, password FROM users WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	err = statement.QueryRow(user.ID).Scan(&user.ID, &user.Username, &user.Password, &user.Mail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates a single user by ID
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Logic for updating a user
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", 500) // error 500 = internal server error
		return
	}
	defer db.Close()

	statement, err := db.Prepare("UPDATE users SET username = ?, password = ?, mail = ?, WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(user.Username, user.Password, user.ID, user.Mail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser deletes a single user by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Logic for deleting a user
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", 500) // error 500 = internal server error
		return
	}
	defer db.Close()

	statement, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Handle for login
func Login(w http.ResponseWriter, r *http.Request) {
	// Logic for login
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Erreur lors de la connexion " +err.Error(), http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Open the database
	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Erreur lors de la connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	
	var dbPassword, role string
	err = db.QueryRow("SELECT password, role FROM users WHERE username = ?", username).Scan(&dbPassword, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
			return
		}
		http.Error(w, "Erreur lors de la requête à la base de données"+err.Error(), http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Values["role"] = role
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Erreur lors de la sauvegarde de la session"+err.Error(), http.StatusInternalServerError)
		return
	}

	if role == "admin" {
		w.Write([]byte("Welcome home Sir !!"))
	} else {
		w.Write([]byte("Welcome Guys !!"))
	}
}

// Handle for logout
func Logout(w http.ResponseWriter, r *http.Request) {
	// Logic for logout
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["authentified"] = false
	
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Erreur lors de la sauvegarde de la session" +err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response to the client
	w.Write([]byte("You are logged out"))
}

// Handle for category
func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Logic for getting a category
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["id"]) // Get the category ID from the URL
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var category Category

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", 500) // error 500 = internal server error
		return
	}
	defer db.Close()

	statement, err := db.Prepare("SELECT id, name FROM categories WHERE id = ?") // Prepare the SQL statement
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer statement.Close()

	err = statement.QueryRow(categoryID).Scan(&category.ID, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error querring database"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)

}

// CreateCategory creates a single category
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Logic for creating a category
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return

	}

	categoryName := r.FormValue("categoryName") // Get the category name from the form

	// Open the database
	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err)
		http.Error(w, "Internal server error", 500) // error 500 = internal server error
		return
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO categories (name) VALUES (?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(categoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Category created: %s", categoryName)
}

// UpdateCategory updates a single category by ID
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Logic for deleting a category
	vars := mux.Vars(r)
	categoryID := vars["id"]

	db, err := OpenDB()
	if err != nil {
		log.Println("Failed to open the database", err) // print the error message without stopping the server
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	statement, err := db.Prepare("DELETE FROM categories WHERE id = ?")
	if err != nil {
		http.Error(w, "Failed to prepare statement"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(categoryID)
	if err != nil {
		http.Error(w, "Failed to execute statement"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
