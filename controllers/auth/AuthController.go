package auth_controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	jwt_services "github.com/yasarunylmzz/wordlingo-backend/services/jwt"
)

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func VerifyRefresh(c echo.Context) error{
	var req RefreshRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message":"Refresh Token missing or invalid"})
	}

	isValid, err := jwt_services.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error verifying refresh token"})
	}

	if isValid != false {
		return nil
	}
	


return c.JSON(http.StatusNotAcceptable, echo.Map{
	"message":"token is invalid",
})
	
}