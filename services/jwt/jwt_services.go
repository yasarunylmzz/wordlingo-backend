// services/jwt_services.go
package jwt_services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

// createAccessToken creates an access token (short-lived).
func CreateAccessToken(username,name, email, surname string, id int) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256,
        jwt.MapClaims{
            "username": username,
            "emaik":email,
            "name":name,
            "surname":surname,
            "id":id,
            "exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // Access token expires in 1 week
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// createRefreshToken creates a refresh token (long-lived).
func CreateRefreshToken(username,name, email, surname string, id int) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256,
        jwt.MapClaims{
            "username": username,
            "emaik":email,
            "name":name,
            "surname":surname,
            "id":id,
            "exp":      time.Now().Add(time.Hour * 24 * 30).Unix(), // Refresh token expires in 1 month
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// Verify Access Token
func VerifyAccessToken(tokenString string) (*jwt.Token, error) {
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
func VerifyRefreshToken(tokenString string) (*jwt.Token, error) {
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
