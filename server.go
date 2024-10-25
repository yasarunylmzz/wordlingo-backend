package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"wordlingo-backend/internal/db"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq" // PostgreSQL sürücüsü
	// queries.sql.go dosyasının modül yolu
)

// Veritabanı bağlantısını `Queries` yapıtaşına aktarmak için global bir değişken tanımlayın
var q *db.Queries

func main() {
	// Veritabanı bağlantısını aç
	conn, err := sql.Open("postgres", "postgres://username:password@localhost:5432/dbname?sslmode=disable")
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer conn.Close()

	// `q` nesnesini veritabanı bağlantısıyla başlat
	q = db.New(conn)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// `createUser` handler fonksiyonunu rota ile bağlayın
	e.POST("/create-user", createUser)
	e.Logger.Fatal(e.Start(":1323"))
}

func createUser(c echo.Context) error {
	// Bağlam oluşturma
	ctx := context.Background()

	// İstek gövdesini (`body`) `CreateUserParams` yapısına çevirme
	var params db.CreateUserParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Veritabanı sorgusunu çağırma
	if err := q.CreateUser(ctx, params); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	// Başarı yanıtı
	return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}
