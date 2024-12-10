// controllers/userController.go
package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/yasarunylmzz/wordlingo-backend/internal/db"
	_ "github.com/yasarunylmzz/wordlingo-backend/mail"
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
		log.Printf("Failed to create user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}
	// Kullanıcı oluşturulduktan sonra doğrulama kodu oluştur
	verificationCode := send_mail.generateVerificationCode() // mail paketinden fonksiyonu çağırın

	// Doğrulama kodunu gönder
	if err := send_mail.sendMail(params.Email, verificationCode); err != nil { // mail paketinden fonksiyonu çağırın
		log.Printf("Failed to send email: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func LoginUser(c echo.Context) error {
	ctx := context.Background()

	var params db.LoginUserParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Debug: Gelen inputu loglayın
	fmt.Printf("Params: %+v", params)

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
	log.Printf("Executing LoginUser query")
	user, err := queries.LoginUser(ctx, params)
	log.Printf("Query executed")
	if err != nil {
		log.Printf("LoginUser error: %v", err)
		if err == sql.ErrNoRows {
			log.Printf("No user found for email: %s", params.Email)
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Not Found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to login user"})
	}
	
	log.Printf("User logged in successfully: %+v", user)
	

	// Debug: Bulunan kullanıcıyı loglayın
	log.Printf("User found: %+v", user)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User logged in successfully",
		"user":    map[string]interface{}{
			"id": user.ID,
			"email": user.Email,
			"name": user.Name,
			"surname": user.Surname,
			"is_verified": user.IsVerified,
			"username": user.Username,
		},
	})
}
