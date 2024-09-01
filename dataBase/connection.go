package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// ConectarBD initializes and returns a connection to the MySQL database
func ConectarBD() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println("Error opening database connection:", err)
		return nil, err
	}

	// Verify if the database connection is successful
	if err := db.Ping(); err != nil {
		log.Println("Failed to connect to the database:", err)
		return nil, err
	}

	log.Println("Database connection successful")
	return db, nil
}
