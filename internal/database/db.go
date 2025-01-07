package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // Untuk memuat file .env
	_ "github.com/go-sql-driver/mysql" // Import driver MySQL
)

var db *sql.DB

// InitDB menginisialisasi koneksi database
func InitDB() (*sql.DB, error) {
	// Memuat file .env
	err := godotenv.Load()
if err != nil {
    log.Fatal("Error loading .env file")
} else {
    log.Println("Successfully loaded .env file")
}


	

	// Ambil nilai variabel dari .env
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	log.Println("DB User:", dbUser)
	log.Println("DB Password:", dbPassword)


	// Bangun data source name untuk MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Membuka koneksi ke database
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
		return nil, err
	}

	// Cek koneksi ke database
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}

// GetDB mengembalikan koneksi database yang sudah terinisialisasi
func GetDB() *sql.DB {
	return db
}
