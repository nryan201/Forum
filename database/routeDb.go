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

// Comment struct
type Comment struct {
    ID 			 int  		`json:"id"`
    Content		 string		`json:"content"`
    TopicID		 int		`json:"topic_id"`
}

// User struct
type User struct {
    ID 			 int  		`json:"id"`
    Mail         string		`json:"mail"`
    Username	 string		`json:"username"`
    Password	 string		`json:"password"`
}

func OpenDB() (*sql.DB, error) {
    return sql.Open("sqlite3", "file:forum.db?cache=shared&mode=rwc")
}

