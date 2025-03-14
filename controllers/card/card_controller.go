package card_controller

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yasarunylmzz/wordlingo-backend/helpers"
	"github.com/yasarunylmzz/wordlingo-backend/internal/db"
)

func CreateCard(c echo.Context) error {
	ctx := context.Background()
    var params db.CreateCardParams
	queries,dbConn, err := helpers.OpenDatabaseConnection()


	if err := c.Bind(&params); err != nil {
		log.Printf("İnvalid inputs", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err != nil {
        log.Printf("Failed to open database connection: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
    }	

	err = queries.CreateCard(ctx, params)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	defer dbConn.Close()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})

}

func UpdateCard(c echo.Context) error {

	ctx := context.Background()
	var params db.UpdateCardParams
	queries, dbConn, err := helpers.OpenDatabaseConnection()


	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}
	

	err = queries.UpdateCard(ctx, params)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	defer dbConn.Close()


	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"params":  params,
	})
}


func DeleteCard(c echo.Context) error {
	ctx := context.Background()
	var params db.DeleteCardParams
	queries, dbConn, err := helpers.OpenDatabaseConnection()


	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	err = queries.DeleteCard(ctx, params)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	defer dbConn.Close()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
		"params": params,
	})
}

func GetAllCardByDeskId(c echo.Context) error {
	ctx := context.Background()
	queries, dbConn, err := helpers.OpenDatabaseConnection()

	deskID := c.QueryParam("desk_id")  

	if deskID == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Desk ID is required"})
    }

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	deskIDstr, err := strconv.Atoi(deskID)

	data, err := queries.GetCardsByDeskId(ctx, int32(deskIDstr))

	
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":err.Error(),
		})
	}

	defer dbConn.Close()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
		"params": data,
	})
}