package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yasarunylmzz/wordlingo-backend/controllers"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/createuser", func(c echo.Context) error {
		return controllers.CreateUser(c) 
	})
	e.Logger.Fatal(e.Start(":1323"))
}
