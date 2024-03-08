package database_test

import (
	"testing"

	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

func TestInitAndGetDB(t *testing.T) {
	// Open temporary in-memory SQLite database for testing
	dbPath := ":memory:"
	err := database.Init(dbPath)
	if err != nil {
		t.Fatalf("Error initializing database: %v", err)
	}
	defer database.Close()

	// Get database connection
	db := database.GetDB()

	// Check if database connection is not nil
	if db == nil {
		t.Fatal("Expected database connection, got nil")
	}

	// Perform simple query to check if database connection is valid
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("Error querying database: %v", err)
	}

	// Check if the result is expected
	if result != 1 {
		t.Fatalf("Expected result to be 1, got %d", result)
	}
}
