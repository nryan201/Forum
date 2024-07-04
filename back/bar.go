package back

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Result struct {
	ID    int
	Title string
}


func SearchDatabase(query string) ([]Result, error) {
	// Open the database
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query("SELECT id, title FROM topics WHERE title LIKE ? UNION SELECT id, username AS title FROM users WHERE username LIKE ?", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Fetch the results
	var results []Result
	for rows.Next() {
		var result Result
		err = rows.Scan(&result.ID, &result.Title)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func SearchTopics (query string)([]Topic, error){

	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, title, description FROM topic WHERE title LIKE ? OR description LIKE ?")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + query + "%", "%" + query + "%")
	if err != nil {
		log.Printf("Error getting topics: %v", err)
		return nil, err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var t Topic
		err := rows.Scan(&t.ID, &t.Title, &t.Description)
		if err != nil {
			log.Printf("Error scanning topic: %v", err)
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}

func GetAllTopics() ([]Topic, error) {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	
	
	// 
	rows, err := db.Query("SELECT id, title, description FROM topics")
	if err != nil {
		log.Printf("Error getting topics: %v", err)
		return nil, err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var t Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.Description); err != nil {
			log.Printf("Error scanning topic: %v", err)
			return nil, err
		}
		topics = append(topics, t)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	return topics, nil
}
