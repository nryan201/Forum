package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Topic struct {
	ID 			 int  		`json:"id"`
	Title		 string		`json:"title"`
	Description	 string		`json:"description"`
}

func OpenDB() (*sql.DB, error) {
    connStr := "user=youruser dbname=yourdbname password=yourpassword host=yourhost sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    return db, nil
}

