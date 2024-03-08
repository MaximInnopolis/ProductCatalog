package utils

import (
	"github.com/joho/godotenv"
	"log"
)

// LoadEnv loads environment variables from a .env file
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
