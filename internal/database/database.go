package database

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Printf("Failed to open database connection: %v", err)
		return err
	}

	// Database connection check
	if err := db.Ping(); err != nil {
		logger.Printf("Failed to ping database: %v", err)
		return err
	}

	logger.Println("Database connection established")

	return nil
}

func Close() {
	if db != nil {
		if err := db.Close(); err != nil {
			logger.Println("Error closing database connection:", err)
		} else {
			logger.Println("Database connection closed")
		}
	}
}

func GetDB() *sql.DB {
	return db
}
