// controllers/user_controller.go
package controllers

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/yasarunylmzz/wordlingo-backend/helpers"
	"github.com/yasarunylmzz/wordlingo-backend/internal/db"
	"github.com/yasarunylmzz/wordlingo-backend/mail"
	hash_services "github.com/yasarunylmzz/wordlingo-backend/services/hash"
	jwt_services "github.com/yasarunylmzz/wordlingo-backend/services/jwt"
)

func CreateUser(c echo.Context) error {
    ctx := context.Background()
    var params db.CreateUserParams
	queries, err := helpers.OpenDatabaseConnection()

    if err := c.Bind(&params); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
    }

    if err != nil {
        log.Printf("Failed to open database connection: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
    }

	encodedHash, err := hash_services.HashPassword(params.Password) 
	if err != nil {
		log.Printf("Şifre hash'lenemedi: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "password_hashing_failed"})
	}
	
	 accessToken, err := jwt_services.CreateAccessToken(params.Username) 
	 if err != nil {
		 log.Printf("Failed to create access token: %v", err)
		 return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create access token"})
	 }
 
	 refreshToken, err := jwt_services.CreateRefreshToken(params.Username) 
	 if err != nil {
		 log.Printf("Failed to create refresh token: %v", err)
		 return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create refresh token"})
	 }

	params.Password = encodedHash
	userID, err := queries.CreateUser(ctx, params)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

    verificationCode := mail.GenerateVerificationCode()
    verificationParams := db.VerificationCodeCreateParams{
        UserID: sql.NullInt32{Int32: userID, Valid: true},
        Code:   verificationCode,
    }

    if _, err := queries.VerificationCodeCreate(ctx, verificationParams); err != nil {
        log.Printf("Failed to create verification code: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create verification code"})
    }

    if err := mail.SendMail(params.Email, verificationCode); err != nil {
        log.Printf("Failed to send email: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
    }

    return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully","accessToken": accessToken,"refreshToken":refreshToken})
}

func LoginUser(c echo.Context) error {
	ctx := context.Background()

	// Veritabanı bağlantısını aç
	queries, err := helpers.OpenDatabaseConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database connection failed",
		})
	}

	// Giriş için gerekli parametreleri al (Artık LoginUserParams yok)
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid input",
		})
	}

	log.Printf("Attempting to login user with email: %s", loginRequest.Email)

	// Kullanıcının hash'lenmiş şifresini veritabanından al
	hashPass, err := queries.GetHashPass(ctx, loginRequest.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error":   "Invalid email or password",
				"details": "No account found with this email",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user",
		})
	}

	// Şifreyi doğrula
	if !hash_services.VerifyPassword(loginRequest.Password, hashPass) {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error":   "Invalid email or password",
			"details": "Password is incorrect",
		})
	}

	// Şifre doğrulandı, kullanıcı bilgilerini çek
	user, err := queries.GetUserByEmail(ctx, loginRequest.Email) 
	if err != nil {
		log.Printf("GetUserByEmail error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Login failed",
			"details": "Database error occurred",
		})
	}

	// Başarılı yanıt
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User logged in successfully",
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"surname": user.Surname,

		},
	})
}
