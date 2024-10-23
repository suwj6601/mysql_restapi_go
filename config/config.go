package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDsnDB() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	return dbUser + ":" + dbPassword + "@tcp(" + dbHost + ")/" + dbName
}
