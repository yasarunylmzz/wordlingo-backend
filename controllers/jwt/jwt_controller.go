package jwt_controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/yasarunylmzz/wordlingo-backend/helpers"
	jwt_services "github.com/yasarunylmzz/wordlingo-backend/services/jwt"
)

type JwtController struct{}

func (ac *JwtController) RefreshAccessToken(c echo.Context) error {
    // 1. Database bağlantısını aç ve kapatmayı unutma
    queries, dbConn, err := helpers.OpenDatabaseConnection()
    if err != nil {
        log.Printf("Database connection failed: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }
    defer dbConn.Close()

    // 2. Refresh token'ı BODY'den al (Header'dan değil!)
    var requestBody struct {
        RefreshToken string `json:"refresh_token"`
    }
    if err := c.Bind(&requestBody); err != nil || requestBody.RefreshToken == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing refresh token"})
    }

    // 3. Refresh token'ı verify et (Doğru servis fonksiyonunu kullan)
    token, err := jwt_services.VerifyRefreshToken(requestBody.RefreshToken)
    if err != nil || !token.Valid {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid refresh token"})
    }

    // 4. Claims'leri doğru şekilde parse et
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
    }

    // 5. UserID'yi güvenli şekilde al
    userID, ok := claims["id"].(float64)
    if !ok {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
    }

    // 6. Kullanıcıyı database'den kontrol et
    user, err := queries.GetUserbyId(c.Request().Context(), int32(userID)) 
    if err != nil {
        if err == sql.ErrNoRows {
            return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
        }
        log.Printf("User retrieval error: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user"})
    }

    // 7. Yeni Access Token üret
    newAccessToken, err := jwt_services.CreateAccessToken(user.Username, user.Name, user.Email, user.Surname, int(user.ID))
    if err != nil {
        log.Printf("Access token creation failed: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create access token"})
    }

    // 8. (Önemli!) Refresh token'ı rotate et (Yeni refresh token üret ve eskisini invalidate et)
    newRefreshToken, err := jwt_services.CreateRefreshToken(user.Username, user.Name, user.Email, user.Surname, int(user.ID)) // int32 → int dönüşümü
	    if err != nil {
        log.Printf("Refresh token rotation failed: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to rotate refresh token"})
    }

    return c.JSON(http.StatusOK, map[string]string{
        "access_token":  newAccessToken,
        "refresh_token": newRefreshToken, // Yeni refresh token'ı da dön
    })
}