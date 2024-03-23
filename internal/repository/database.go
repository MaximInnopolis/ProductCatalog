package repository

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

// NewDB initializes database connection using provided database path
// Returns  error if  connection cannot be established
func NewDB(dbPath string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Println("Failed to open database connection:", err)
		return nil, err
	}

	// Database connection check
	if err = db.Ping(); err != nil {
		logger.Println("Failed to ping database:", err)
		return nil, err
	}

	logger.Println("Database connection established")

	return db, nil
}

// Close closes database connection
func Close(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			logger.Println("Error closing database connection:", err)
		} else {
			logger.Println("Database connection closed")
		}
	}
}
