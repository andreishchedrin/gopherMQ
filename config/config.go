package config

import (
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load("config/.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
