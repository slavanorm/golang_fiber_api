package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config func to get env value
func Config(key string) string {
	// load .env file
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Print("Error loading .env file")
	}
	return os.Getenv(key)
}
