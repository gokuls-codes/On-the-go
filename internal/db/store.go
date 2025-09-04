package db

import (
	"database/sql"
	"log"

	"github.com/gokuls-codes/on-the-go/internal/db/sqlc"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	*sqlc.Queries
	db *sql.DB
}

func NewStore() (*Store, error) {
	conn, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Database connection successful!")

	return &Store{
		db:      conn,
		Queries: sqlc.New(conn),
	}, nil
}
