package routes

import (
	"github.com/labstack/echo"
	"github.com/yasarunylmzz/wordlingo-backend/controllers"
)

func RegisterUserRoutes() {
	e.POST("/createuser", func(c echo.Context) error {
		return controllers.CreateUser(c) 
	})
}