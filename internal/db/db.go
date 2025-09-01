package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Database connection successful!")

	return db, nil
}
