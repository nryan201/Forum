package back

import (
	"database/sql"
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