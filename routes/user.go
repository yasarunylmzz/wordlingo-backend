package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/yasarunylmzz/wordlingo-backend/controllers"
)

func RegisterUserRoutes(e *echo.Echo) {
	e.POST("/createuser", controllers.CreateUser)

}

func LoginUserRoutes(e *echo.Echo){
	e.GET("/login", controllers.LoginUser)
}

func VerificationUserRouters(e *echo.Echo){
	e.POST("/verification", controllers.UserVerification)
}