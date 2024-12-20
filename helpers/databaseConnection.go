package helpers

import "database/sql"

func openDatabaseConnection() (*sql.DB, error) {
    connStr := "postgres://postgres:abc123@localhost:5432/mydb?sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    // Veritabanına bağlantı sağlandı mı, kontrol edelim
    if err := db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}