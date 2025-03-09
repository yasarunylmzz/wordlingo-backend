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
            "email":email,
            "name":name,
            "surname":surname,
            "id":id,
            "exp":      time.Now().Add(time.Minute * 15).Unix(), // Access token expires in 15 minute
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
            "email":email,
            "name":name,
            "surname":surname,
            "id":id,
            "exp":      time.Now().Add(time.Hour * 24 * 30 * 12).Unix(), // Refresh token expires in 1 year
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// Verify Access Token
func VerifyAccessToken(tokenString string) (*jwt.Token, error) {
   token, err := jwt.ParseWithClaims(tokenString,jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        if token.Method != jwt.SigningMethodHS256 {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })

    // fmt.Print(token)
    if err != nil || !token.Valid {
        return nil, fmt.Errorf("invalid access token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, fmt.Errorf("invalid claims")
    }

    if exp, ok := claims["exp"].(float64); ok {
        if time.Now().Unix() > int64(exp) {
            return nil, fmt.Errorf("access token expired")
        }
    }

    return token, nil

}

// Verify Refresh Token
func VerifyRefreshToken(tokenString string) (*jwt.Token, error) {
    return jwt.ParseWithClaims(tokenString,jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        if token.Method != jwt.SigningMethodHS256 {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })

}
