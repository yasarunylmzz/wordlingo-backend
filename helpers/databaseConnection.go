package helpers

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yasarunylmzz/wordlingo-backend/internal/db"
)

func OpenDatabaseConnection() (*db.Queries, *sql.DB, error) {
	// GO_ENV ortam değişkenini kontrol ediyoruz, varsayılan development
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Ortama göre ilgili .env dosyasını belirle
	var envFile string
	if env == "production" {
		envFile = ".env.production"
	} else {
		envFile = ".env.development"
	}

	// İlgili .env dosyasını yükle
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Could not load %s file: %v", envFile, err)
	}

	// Ortam değişkenlerini oku
	dbHost := os.Getenv("DB_HOST")      // Yerelde 'localhost', prod'da 'db'
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Bağlantı string'ini oluştur
	connStr := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil, nil, err
	}

	if err := dbConn.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, nil, err
	}

	return db.New(dbConn), dbConn, nil
}
