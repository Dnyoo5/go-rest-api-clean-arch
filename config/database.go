package config

import (
	"database/sql"
	"fmt"
	"log"
	"os" // Library bawaan Go untuk baca sistem environment

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv" // Import library yang baru diinstall
)

func ConnectDB() *sql.DB {
    err := godotenv.Load()
    if err != nil {
        log.Println("Note: Tidak menemukan file .env, menggunakan environment system")
    }

    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Gagal membuka koneksi:", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Gagal konek ke database:", err)
    }

    fmt.Println("Config: Database Connected (Secure Mode) 🔒")
    return db
}