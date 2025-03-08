package jwt_middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	jwt_services "github.com/yasarunylmzz/wordlingo-backend/services/jwt"
)


func RefreshAccessTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == ""{
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Authorization header is required"})
		}
		
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { 
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid authorization header"})
		}

		//verify refresh token
		_, err := jwt_services.VerifyRefreshToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired refresh token"})
		}

		var request struct {
			AccessToken string `json:"access_token"`
		}

		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
		}

		//verify access token
		_, err = jwt_services.VerifyAccessToken(request.AccessToken)
		// fmt.Print(err)
		if err == nil {
			return c.JSON(http.StatusOK, map[string]string{"message": "Access token is still valid","access_token": request.AccessToken,"refresh_token":tokenString})
		}

		newAccessToken, err := jwt_services.CreateAccessToken("username", "name", "email", "surname", 1)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not generate new access token"})
		}

		// if access token is invalid return newAccessToken
		return c.JSON(http.StatusOK, map[string]string{"access_token": newAccessToken,"refresh_token": tokenString})


	}
}
