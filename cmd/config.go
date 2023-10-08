package cmd

import (
	"log"
	"os"
	"transaction_system/app/lib/db"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load("development.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func SetupDBConnection() {
	dbURL := os.Getenv("DATABASE_URL")

	log.Println("Connecting to the database")
	if err := db.Connect(dbURL, 10, 10); err != nil {
		panic(err)
	}
}
