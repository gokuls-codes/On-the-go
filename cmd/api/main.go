package main

import (
	"log"

	"github.com/gokuls-codes/on-the-go/internal/config"
	"github.com/gokuls-codes/on-the-go/internal/db"
	"github.com/gokuls-codes/on-the-go/internal/server"
)

func main() {
	config := config.NewConfig()

	store, err := db.NewStore()
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(config.Port, store)
	server.Start()
}
