package routes

import (
	"github.com/labstack/echo/v4"
	desk_controller "github.com/yasarunylmzz/wordlingo-backend/controllers/desk"
	user_controller "github.com/yasarunylmzz/wordlingo-backend/controllers/user"
	jwt_middleware "github.com/yasarunylmzz/wordlingo-backend/middleware/jwt"
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
	e.POST("/create-desk", desk_controller.CreateDesk, jwt_middleware.RefreshAccessTokenMiddleware)
}

func UpdateDesk(e *echo.Echo){
	e.PATCH("/update-desk", desk_controller.UpdateDesk, jwt_middleware.RefreshAccessTokenMiddleware)
}