package helpers

import (
	"database/sql"
	"log"

	"github.com/yasarunylmzz/wordlingo-backend/internal/db"
)



func OpenDatabaseConnection() (*db.Queries, error) {
	connStr := "postgres://postgres:abc123@localhost:5432/flashcards?sslmode=disable"
	dbConn, err := sql.Open("postgres", connStr)
    
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil, err
	}
	if err := dbConn.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, err
	}

	return db.New(dbConn), nil
}
