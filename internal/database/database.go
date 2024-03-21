package database

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Init initializes database connection using provided database path
// Returns  error if  connection cannot be established
func Init(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Println("Failed to open database connection:", err)
		return err
	}

	// Database connection check
	if err = db.Ping(); err != nil {
		logger.Println("Failed to ping database:", err)
		return err
	}

	logger.Println("Database connection established")

	return nil
}

// Close closes database connection
func Close() {
	if db != nil {
		if err := db.Close(); err != nil {
			logger.Println("Error closing database connection:", err)
		} else {
			logger.Println("Database connection closed")
		}
	}
}

// GetDB returns reference to database connection
func GetDB() *sql.DB {
	return db
}
