package card_controller

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yasarunylmzz/wordlingo-backend/helpers"
	"github.com/yasarunylmzz/wordlingo-backend/internal/db"
)

func CreateCard(c echo.Context) error {
	ctx := context.Background()
    var params db.CreateCardParams
	queries,dbConn, err := helpers.OpenDatabaseConnection()


	if err := c.Bind(&params); err != nil {
		log.Printf("İnvalid inputs")
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err != nil {
        log.Printf("Failed to open database connection: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
    }	
	return nil
}