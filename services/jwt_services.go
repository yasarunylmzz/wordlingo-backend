package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("secret-key")

// createAccessToken creates an access token (short-lived).
func createAccessToken(password string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256,
        jwt.MapClaims{
            "password": password,
            "exp":      time.Now().Add(time.Hour * 24).Unix(), // Access token expires in 1 day
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// createRefreshToken creates a refresh token (long-lived).
func createRefreshToken(password string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256,
        jwt.MapClaims{
            "password": password,
            "exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // Refresh token expires in 1 week
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// Verify Access Token
func verifyAccessToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if token.Method != jwt.SigningMethodHS256 {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })
    if err != nil || !token.Valid {
        return nil, fmt.Errorf("invalid access token")
    }
    return token, nil
}

// Verify Refresh Token
func verifyRefreshToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if token.Method != jwt.SigningMethodHS256 {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })
    if err != nil || !token.Valid {
        return nil, fmt.Errorf("invalid refresh token")
    }
    return token, nil
}
