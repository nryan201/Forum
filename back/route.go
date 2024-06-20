package back

import (
    "database/sql"
    "encoding/json"
    "html/template"
    "log"
    "net/http"
    "strconv"
	"main/database"

    "github.com/gorilla/mux"
    _ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Handle for home page
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

// Handle for  Create topic 
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	// Logic for creating a topic
	var topic database.Topic 
	err :=json.NewDecoder(r.Body).Decode(&topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
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
}


// GetTopic retrieves a single topic by ID from PortgresSql
func GetTopic(w http.ResponseWriter, r *http.Request) {
	// Logic for getting a topic
	vars := mux.Vars(r)
	topicID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var topic database.Topic
	err = db.QueryRow("SELECT id, title, description FROM topics WHERE id = ?", topicID).Scan(&topic.ID, &topic.Title, &topic.Description)
	if err != nil{
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

	var topic database.Topic
	err = json.NewDecoder(r.Body).Decode(&topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
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

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
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
	var comment database.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
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
	var comment database.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
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
	var comment database.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
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
	var comment database.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
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
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.OpenDB()

	if err != nil {
		log.Fatal(err)
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Logic for getting a user
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
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

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Logic for updating a user
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Logic for deleting a user
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
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
}

// Handle for logout
func Logout(w http.ResponseWriter, r *http.Request) {
	// Logic for logout
}

