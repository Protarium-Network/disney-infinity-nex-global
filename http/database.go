package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDatabase connects to PostgreSQL using HTTP_POSTGRES_URI. When the
// variable is unset the HTTP server runs without a database.
func InitDatabase() error {
	uri := os.Getenv("HTTP_POSTGRES_URI")
	if uri == "" {
		log.Println("HTTP_POSTGRES_URI not set, running without a database")
		return nil
	}

	db, err := sql.Open("postgres", uri)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return err
	}

	DB = db
	log.Println("Connected to PostgreSQL database")

	return nil
}

func CloseDatabase() {
	if DB != nil {
		DB.Close()
	}
}
