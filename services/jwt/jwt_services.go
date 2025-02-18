package jwt_services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

// CreateAccessToken generates a short-lived access token (e.g., 15 minutes)
func CreateAccessToken(username, name, email, surname string, id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"email":    email,
		"name":     name,
		"surname":  surname,
		"id":       id,
		"exp":      time.Now().Add(time.Minute * 15).Unix(), // 15-minute expiry
	})

	return token.SignedString(secretKey)
}

// CreateRefreshToken generates a long-lived refresh token (e.g., 30 days)
func CreateRefreshToken(username, name, email, surname string, id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"email":    email,
		"name":     name,
		"surname":  surname,
		"id":       id,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days expiry
	})

	return token.SignedString(secretKey)
}

// VerifyAccessToken checks if the given access token is valid
func VerifyAccessToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}
	return token, nil
}

// VerifyRefreshToken checks if the given refresh token is valid
func VerifyRefreshToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}
	return token, nil
}
