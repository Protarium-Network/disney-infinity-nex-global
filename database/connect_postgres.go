package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ConnectPostgres opens a connection pool and verifies the server answers.
func ConnectPostgres(uri string) (*sql.DB, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}
