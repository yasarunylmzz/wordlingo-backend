// controllers/desk_controller.go
package desk_controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yasarunylmzz/wordlingo-backend/helpers"
	"github.com/yasarunylmzz/wordlingo-backend/internal/db"
)



func CreateDesk(c echo.Context) error {

	ctx := context.Background()
	var params db.CreateDeskParams
	queries, dbConn, err := helpers.OpenDatabaseConnection()

	fmt.Print(params.Description)

	if err := c.Bind(&params); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
    }

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message":"internal server error"})
	}
	
	err = queries.CreateDesk(ctx, params)

	fmt.Print(err)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message":"error while creating desk"})
	}

	dbConn.Close()

	return c.JSON(http.StatusAccepted, map[string]string{"message":"ok","title":params.Title,"description":params.Description,"image_url":  func() string {
        if params.ImageLink.Valid {
            return params.ImageLink.String
        }
        return ""  // veya default bir değer
    }(),})
}

func UpdateDesk(c echo.Context) error {
	ctx := context.Background()
	var params db.UpdateDeskParams
	queries, dbConn, err := helpers.OpenDatabaseConnection()

	if err != nil {
		return c.JSON(http.StatusInternalServerError,map[string]string{"message":"internal server error"})
	}

	if c.Bind(&params); err != nil{
		return c.JSON(http.StatusBadRequest, map[string]string{"message":"Invalid Input"})
	}

	err = queries.UpdateDesk(ctx, params)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, map[string]string{"message":err.Error(),"message2":"erroorrr"})
	}

	dbConn.Close()

	return c.JSON(http.StatusOK, map[string]string{"message":"ok"})
}

func DeleteDesk(c echo.Context) error {

	return c.JSON(http.StatusOK, map[string]string{"message":"ok"})

}