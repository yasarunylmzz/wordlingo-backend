// controllers/userController.go
package controllers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/yasarunylmzz/wordlingo-backend/internal/db"
)

func CreateUser(c echo.Context) error {
	ctx := context.Background()

	var params db.CreateUserParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	connStr := "postgres://postgres:abc123@localhost:5432/flashcards?sslmode=disable"
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to ping database"})
	}

	queries := db.New(dbConn)
	if err := queries.CreateUser(ctx, params); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}
