package routes

import (
	"github.com/labstack/echo/v4"
	auth_controller "github.com/yasarunylmzz/wordlingo-backend/controllers/auth"
	card_controller "github.com/yasarunylmzz/wordlingo-backend/controllers/card"
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

func DeleteDesk(e *echo.Echo) {
	e.DELETE("/delete-desk", desk_controller.DeleteDesk, jwt_middleware.RefreshAccessTokenMiddleware)
}

func GetAllDesk(e *echo.Echo) {
	e.GET("/get-desk", desk_controller.GetAllDesk, jwt_middleware.RefreshAccessTokenMiddleware)
}

func CreateCard(e *echo.Echo){
	e.POST("/create-card",card_controller.CreateCard,jwt_middleware.RefreshAccessTokenMiddleware)
}

func DeleteCard(e *echo.Echo){
	e.DELETE("/delete-card", card_controller.DeleteCard, jwt_middleware.RefreshAccessTokenMiddleware)
}

func UpdateCard(e *echo.Echo){
	e.PATCH("/update-card", card_controller.UpdateCard, jwt_middleware.RefreshAccessTokenMiddleware)
}

func GetAllCardByDeskId(e *echo.Echo){
	e.POST("/get-card", card_controller.GetAllCardByDeskId, jwt_middleware.RefreshAccessTokenMiddleware)
}

func VerifyRefreshTokenInSplashScreen(e *echo.Echo){
	e.POST("/verify-refresh", auth_controller.VerifyRefresh)
}
