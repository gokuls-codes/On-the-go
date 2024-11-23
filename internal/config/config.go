package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Port string
	PostgresURL string
}


func NewConfig() *config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &config{
		Port: getEnv("PORT", "8080"),
		PostgresURL: getEnv("POSTGRES_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}