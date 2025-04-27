package jwt_middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	jwt_services "github.com/yasarunylmzz/wordlingo-backend/services/jwt"
)


func RefreshAccessTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		authHeaderRefresh := c.Request().Header.Get("Authorization")
		if authHeaderRefresh == ""{
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Authorization header is required"})
		}
		
		tokenString := strings.TrimPrefix(authHeaderRefresh, "Bearer ")
		if tokenString == authHeaderRefresh { 
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid authorization header"})
		}

		_, err := jwt_services.VerifyRefreshToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired refresh token"})
		}

		authHeaderAccess := c.Request().Header.Get("X-Access-Token")
		if authHeaderAccess == ""{
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "x-access-token header is required"})
		}

		tokenString2 := strings.TrimPrefix(authHeaderAccess, "Bearer ")
		if tokenString2 == authHeaderAccess { 
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid authorization header"})
		}

		_, err = jwt_services.VerifyAccessToken(tokenString2)
		fmt.Print(err)
		if err != nil {
			newAccessToken, err := jwt_services.CreateAccessToken("username", "name", "email", "surname", fmt.Sprintf("%d", 1))
			if err != nil{
				return c.JSON(http.StatusNotAcceptable, map[string]string{"message":err.Error()})
			}
			c.Response().Header().Set("New-Access-Token", newAccessToken)
			return c.JSON(http.StatusBadRequest, map[string]string{"message":"please take new access token in headers"})
		}

		return next(c)


	}
}
