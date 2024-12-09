package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/yasarunylmzz/wordlingo-backend/controllers"
)

func RegisterUserRoutes(e *echo.Echo) {
	e.POST("/createuser", func(c echo.Context) error {
		return controllers.CreateUser(c) 
	})

}

func LoginUserRoutes(e *echo.Echo){
	e.GET("/login", func (c echo.Context) error{
		return controllers.LoginUser(c)
	})
}