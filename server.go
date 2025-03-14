package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yasarunylmzz/wordlingo-backend/routes"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:19006",
			"http://192.168.1.11:19006",
		},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))
	
	routes.RegisterUserRoutes(e)
	routes.LoginUserRoutes(e)
	routes.VerificationUserRouters(e)
	routes.CreateDesk(e)
	routes.UpdateDesk(e)
	routes.DeleteDesk(e)
	routes.GetAllDesk(e)
	routes.CreateCard(e)
	routes.DeleteCard(e)
	routes.UpdateCard(e)
	routes.GetAllCardByDeskId(e)
	e.Logger.Fatal(e.Start(":1323"))
}
