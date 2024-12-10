package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yasarunylmzz/wordlingo-backend/routes"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	routes.RegisterUserRoutes(e)
	routes.LoginUserRoutes(e)
	
	e.Logger.Fatal(e.Start(":1323"))
}
