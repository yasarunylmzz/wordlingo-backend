package routes

import (
	"github.com/labstack/echo/v4"
	user_controller "github.com/yasarunylmzz/wordlingo-backend/controllers/user"
)

func RegisterUserRoutes(e *echo.Echo) {
	e.POST("/createuser", user_controller.CreateUser)
}

func LoginUserRoutes(e *echo.Echo){
	e.POST("/login", user_controller.LoginUser)
}

func VerificationUserRouters(e *echo.Echo){
	e.POST("/verification", user_controller.UserVerification)
}

func CreateDesk(e *echo.Echo){
	e.POST("/create-desk",desk_controller.createDesk)
}