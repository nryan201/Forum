package back

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"main/data"

)

func HomeHandle (w http.ResponseWriter, r *http.Request) {
	
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

// Handle for topic 
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	// Logic for creating a topic
	if r.Method == "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Extract the data

	title := r.Form.Get("title")
	description := r.FormValue("description")

	//Validate the data
	if title == "" || description == ""{
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	
	}

	// Open the database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//prepare the statement
	stmt, err := db.Prepare("INSERT INTO topics (title, description) VALUES (?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	//Execute the statement
	_, err = stmt.Exec(title, description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/succesPage", http.StatusSeeOther)
}


// GetTopic retrieves a single topic by ID from SQLite.
func GetTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["id"]

	db := data.OpenDB()
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, title, description, created_at FROM topics WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var topic data.Topic
	err = stmt.QueryRow(topicID).Scan(&topic.ID, &topic.Title, &topic.Description, &topic.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "ID: %d\nTitle: %s\nDescription: %s\nCreated At: %s\n", topic.ID, topic.Title, topic.Description, topic.CreatedAt)
}




func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	
}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	// Logic for deleting a topic

}

// Handle for comment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	// Logic for creating a comment
}

func GetComment(w http.ResponseWriter, r *http.Request) {
	// Logic for getting a comment
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	// Logic for updating a comment
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	// Logic for deleting a comment
}

// Handle for user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Logic for creating a user
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Logic for getting a user
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Logic for updating a user
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Logic for deleting a user
}

// Handle for login
func Login(w http.ResponseWriter, r *http.Request) {
	// Logic for login
}

// Handle for logout
func Logout(w http.ResponseWriter, r *http.Request) {
	// Logic for logout
}

