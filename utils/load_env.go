package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load_env() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func EnvGet(key string) string {
	return os.Getenv(key)
}
