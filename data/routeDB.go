package data

import (
	"database/sql"
	_"github.com/mattn/go-sqlite3"

	"log"
)

type Topic struct {
	ID 			 int
	Title		 string
	Description	 string
	CreatedAt	 string
}

func OpenDB()*sql.DB{
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
