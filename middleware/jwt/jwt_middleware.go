package jwt_middleware

import (
	"bytes"
	"fmt"
	"io"
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

		// **Request body’yi oku ve sakla**
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to read request body"})
		}
		// **Request body’yi tekrar kullanılabilir yap**
		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		
		var request struct {
			AccessToken string `json:"access_token"`
		}

		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
		}

		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		//verify access token
		_, err = jwt_services.VerifyAccessToken(request.AccessToken)
		fmt.Print(err)
		if err != nil {
			newAccessToken, err := jwt_services.CreateAccessToken("username", "name", "email", "surname", 1)
			if err != nil{
				return c.JSON(http.StatusNotAcceptable, map[string]string{"message":err.Error()})
			}
			c.Response().Header().Set("New-Access-Token", newAccessToken)
			return c.JSON(http.StatusBadRequest, map[string]string{"message":"please take new access token in headers"})
		}

		return next(c)


	}
}
