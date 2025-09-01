package main

import (
	"log"

	"github.com/gokuls-codes/on-the-go/internal/config"
	"github.com/gokuls-codes/on-the-go/internal/db"
	"github.com/gokuls-codes/on-the-go/internal/server"
)

func main() {
	config := config.NewConfig()

	db, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	
	server := server.NewServer(config.Port, db)
	server.Start()
}