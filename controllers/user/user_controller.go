// controllers/user_controller.go
package user_controller

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
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
	queries,dbConn, err := helpers.OpenDatabaseConnection()

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

	params.Password = encodedHash
	userID, err := queries.CreateUser(ctx, params)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

    verificationCode := mail.GenerateVerificationCode()
    verificationParams := db.VerificationCodeCreateParams{
        UserID: uuid.NullUUID{
		 UUID: userID,
		 Valid: true,
		},
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
	defer dbConn.Close()

	return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully", "user_id": userID.String()})
}


func LoginUser(c echo.Context) error {
	ctx := context.Background()

	queries,dbConn, err := helpers.OpenDatabaseConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database connection failed",
		})
	}

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

	if !hash_services.VerifyPassword(loginRequest.Password, hashPass) {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error":   "Invalid email or password",
			"details": "Password is incorrect",
		})
	}

	user, err := queries.GetUserByEmail(ctx, loginRequest.Email) 
	if err != nil {
		log.Printf("GetUserByEmail error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Login failed",
			"details": "Database error occurred",
		})
	}
	accessToken, err := jwt_services.CreateAccessToken(user.Username, user.Name, user.Email, user.Surname, user.ID.String())
	if err != nil {
		 log.Printf("Failed to create access token: %v", err)
		 return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create access token"})
	 }
 
	 refreshToken, err := jwt_services.CreateRefreshToken(user.Username, user.Name, user.Email, user.Surname, user.ID.String())
	 if err != nil {
		 log.Printf("Failed to create refresh token: %v", err)
		 return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create refresh token"})
	 }

	if !user.IsVerified.Valid || (user.IsVerified.Valid && !user.IsVerified.Bool) {
		log.Printf("verification is false or not set")
		return c.JSON(http.StatusNonAuthoritativeInfo, map[string]interface{}{
			"error":   "verification required",
			"details": "please complete verification process",
			"user": map[string]interface{}{
				"id": user.ID,
				"email": user.Email,
				"name":  user.Name,
				"surname": user.Surname,
				"username":user.Username,
				"is_verified": user.IsVerified.Bool,
			},
		})
	}
	defer dbConn.Close()

	c.Response().Header().Set("access_token", accessToken)
	c.Response().Header().Set("refresh_token", refreshToken)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User logged in successfully",
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"surname": user.Surname,
			"username":user.Username,
			"is_verified": user.IsVerified.Bool,
			
		},
	})
}

func UserVerification(c echo.Context) error {
	ctx := context.Background()
	queries, dbConn, err := helpers.OpenDatabaseConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database connection failed",
		})
	}

	var params db.VerifyUserParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	code, err := queries.GetVerificationCodeById(ctx, uuid.NullUUID{
		UUID: params.ID,
		Valid: true,
	})


	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error retrieving verification code",
			"message": err.Error(),
		})
	}
	if code.Code != params.Code {
		return c.JSON(http.StatusNotAcceptable, map[string]string{
			"error":"code expired or not allowed",
		})
	}

	 if err := queries.VerifyUser(ctx, params); err != nil {
	 	return c.JSON(http.StatusInternalServerError, map[string]string{
	 		"error": "User verification failed",
	 	})
	 }
	defer dbConn.Close() 

	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "User verified successfully",
	})
}