// services/jwt_services.go
package jwt_services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

// createAccessToken creates an access token (short-lived).
func CreateAccessToken(username,name, email, surname string, id string) (string, error) {
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
func CreateRefreshToken(username,name, email, surname string, id string) (string, error) {
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
// VerifyRefreshToken kontrol fonksiyonu
func VerifyRefreshToken(tokenString string) (bool, error) {
    token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Signature algoritması doğru mu?
        if token.Method != jwt.SigningMethodHS256 {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })

    if err != nil || !token.Valid {
        return false, fmt.Errorf("token geçersiz: %v", err)
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return false, fmt.Errorf("claims parse edilemedi")
    }

    // Expiration claimini kontrol edelim
    if exp, ok := claims["exp"].(float64); ok {
        if int64(exp) < time.Now().Unix() {
            return false, fmt.Errorf("token süresi dolmuş")
        }
    } else {
        return false, fmt.Errorf("exp claimi bulunamadı")
    }

    return true, nil // her şey yolunda, süresi dolmamış
}