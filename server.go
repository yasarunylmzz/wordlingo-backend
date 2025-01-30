package main

import (
	"github.com/labstack/echo/v4"
	"github.com/yasarunylmzz/wordlingo-backend/routes"
)

func main() {
	e := echo.New()
	routes.RegisterUserRoutes(e)
	routes.LoginUserRoutes(e)
	routes.VerificationUserRouters(e)
	e.Logger.Fatal(e.Start(":1323"))
}
