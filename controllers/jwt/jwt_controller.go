package jwt_controller

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/yasarunylmzz/wordlingo-backend/helpers"
	jwt_services "github.com/yasarunylmzz/wordlingo-backend/services/jwt"
)

func RefreshAccessToken(c echo.Context) error {
	// 1. Database bağlantısı
	queries, dbConn, err := helpers.OpenDatabaseConnection()
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer dbConn.Close()

	// 2. Bearer token'ı al
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing Authorization header"})
	}

	// Bearer token'ı parse et
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Authorization header format"})
	}
	refreshToken := parts[1]

	// 3. Refresh token'ı verify et
	refreshTokenParsed, err := jwt_services.VerifyRefreshToken(refreshToken)
	if err != nil || !refreshTokenParsed.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid refresh token"})
	}

	// 4. Claims'leri al
	claims, ok := refreshTokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
	}

	// 5. UserID'yi al
	userID, ok := claims["id"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	// 6. Kullanıcıyı kontrol et
	user, err := queries.GetUserbyId(c.Request().Context(), int32(userID))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		log.Printf("User retrieval error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user"})
	}

	// 7. Refresh token süresini kontrol et
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Refresh token expired"})
		}
	}

	// 8. Access token süresini kontrol et (Client'tan gelen access token)
	var req struct {
		AccessToken string `json:"access_token"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	accessValid := false
	if req.AccessToken != "" {
		accessToken, _ := jwt_services.VerifyAccessToken(req.AccessToken)
		if accessToken != nil && accessToken.Valid {
			accessValid = true
		}
	}

	// 9. Access token hala geçerliyse yenileme yapma
	if accessValid {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"valid":         true,
			"access_token":  req.AccessToken,
			"refresh_token": refreshToken,
		})
	}

	// 10. Yeni token'ları üret
	newAccessToken, err := jwt_services.CreateAccessToken(user.Username, user.Name, user.Email, user.Surname, int(user.ID))
	if err != nil {
		log.Printf("Access token creation failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create access token"})
	}

	newRefreshToken, err := jwt_services.CreateRefreshToken(user.Username, user.Name, user.Email, user.Surname, int(user.ID))
	if err != nil {
		log.Printf("Refresh token creation failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create refresh token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"valid":         true,
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}