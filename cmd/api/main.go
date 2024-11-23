package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gokuls-codes/on-the-go/internal/config"
	"github.com/gokuls-codes/on-the-go/internal/db"
	"github.com/gokuls-codes/on-the-go/internal/server"
)

func main() {
	config := config.NewConfig()
	conn, err := db.NewDatabase(config.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)

	server := server.NewServer(config.Port, conn)
	server.Start()
}