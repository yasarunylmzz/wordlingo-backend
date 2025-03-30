package helpers

import (
	"database/sql"
	"log"

	"github.com/yasarunylmzz/wordlingo-backend/internal/db"
)



func OpenDatabaseConnection() (*db.Queries,*sql.DB, error) {
	connStr := "postgres://postgres:abc123@db:5432/flashcards?sslmode=disable"
	dbConn, err := sql.Open("postgres", connStr)
    
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil, nil, err
	}
	if err := dbConn.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil,nil, err
	}

	return db.New(dbConn),dbConn, nil
}
