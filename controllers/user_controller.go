package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/yasarunylmzz/wordlingo-backend/internal/db"
)

func createUser(c echo.Context) error {
	// Bağlam oluşturma
	ctx := context.Background()

	// İstek gövdesini (`body`) `CreateUserParams` yapısına çevirme
	var params db.CreateUserParams // `db` paketindeki `CreateUserParams` yapısını kullanıyoruz
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Veritabanı bağlantı dizesi
	connStr := "postgres://postgres:abc123@localhost:5432/flashcards?sslmode=disable"
	dbConn, err := sql.Open("postgres", connStr)

	// Veritabanı bağlantısı
	if err != nil {
		fmt.Println("Database connection error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
	}
	defer dbConn.Close()

	// Bağlantıyı doğrulama
	if err := dbConn.Ping(); err != nil {
		fmt.Println("Database ping error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to ping database"})
	}

	// Queries nesnesi oluşturma
	queries := db.New(dbConn)

	// Veritabanı sorgusunu çağırma
	if err := queries.CreateUser(ctx, params); err != nil {
		fmt.Println("Failed to create user:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	// Başarı yanıtı
	return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}
