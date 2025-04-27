// controllers/desk_controller.go
package desk_controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yasarunylmzz/wordlingo-backend/helpers"
	"github.com/yasarunylmzz/wordlingo-backend/internal/db"
)

type DeskRequest struct {
	UserID int32 `json:"user_id"`
}


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

	defer dbConn.Close()

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

	defer dbConn.Close()

	return c.JSON(http.StatusOK, map[string]string{"message":"ok"})
}

func DeleteDesk(c echo.Context) error {
	ctx := context.Background()
	var params db.DeleteDeskParams
	queries, dbConn, err := helpers.OpenDatabaseConnection()

	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"message":"internal server error"})
	}

	if c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message":"Invalid Input"})
	}

	err = queries.DeleteDesk(ctx, params)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message":err.Error(),
		})
	}

	defer dbConn.Close()

	return c.JSON(http.StatusOK, map[string]string{"message":"ok"})
}


func GetAllDesk(c echo.Context) error {
    // URL'den gelen query parametrelerini almak
    userID := c.QueryParam("user_id")  // user_id'yi query parametre olarak alıyoruz.

    if userID == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "User ID is required"})
    }

    ctx := context.Background()
    queries, dbConn, err := helpers.OpenDatabaseConnection()

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
    }

	userIDParsed, err := uuid.Parse(userID)
    desks, err := queries.GetAllDesksByUserId(ctx, userIDParsed)

    defer dbConn.Close()

    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, map[string]interface{}{"data": desks})
}
