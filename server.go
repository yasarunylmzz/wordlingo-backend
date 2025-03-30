package main

import (
	"os"

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

	httpPort := os.Getenv("PORT")
    if httpPort == "" {
        httpPort = "8080"
    }

	e.Logger.Fatal(e.Start(":"+httpPort))
}
