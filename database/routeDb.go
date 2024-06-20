package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Topic struct
type Topic struct {
	ID 			 int  		`json:"id"`
	Title		 string		`json:"title"`
	Description	 string		`json:"description"`
}


func OpenDB() (*sql.DB, error) {
    return sql.Open("sqlite3", "file:forum.db?cache=shared&mode=rwc")
}

