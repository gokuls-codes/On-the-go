package main

import (
	"github.com/gokuls-codes/on-the-go/internal/config"
	"github.com/gokuls-codes/on-the-go/internal/server"
)

func main() {
	config := config.NewConfig()
	
	server := server.NewServer(config.Port)
	server.Start()
}