package models_test

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestAddCategory tests AddCategory function
func TestAddCategory(t *testing.T) {
	// Open in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create necessary tables
	createCategoryTable(db)

	// Create test category
	category := &models.Category{Name: "Test Category"}

	// Add category to database
	categoryID, err := models.AddCategory(db, category)
	if err != nil {
		t.Fatalf("Error adding category: %v", err)
	}

	// Check if category ID is greater than zero
	if categoryID <= 0 {
		t.Error("Expected category ID to be greater than zero")
	}

	// Retrieve all categories from database
	categories, err := models.GetAllCategories(db)
	if err != nil {
		t.Fatalf("Error retrieving categories: %v", err)
	}

	// Check if newly added category exists in list of categories
	found := false
	for _, c := range categories {
		if c == category.Name {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected newly added category to exist in list of categories")
	}
}

// TestUpdateCategory tests UpdateCategory function
func TestUpdateCategory(t *testing.T) {
	// Open in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create necessary tables
	createCategoryTable(db)

	// Create new category
	category := &models.Category{Name: "Test Category"}

	// Add category to database
	_, err = models.AddCategory(db, category)
	if err != nil {
		t.Fatalf("Error adding category: %v", err)
	}

	// Update category name
	newCategoryName := "Updated Category"
	err = models.UpdateCategory(db, category.Name, &models.Category{Name: newCategoryName})
	if err != nil {
		t.Fatalf("Error updating category: %v", err)
	}

	// Retrieve all categories from database
	categories, err := models.GetAllCategories(db)
	if err != nil {
		t.Fatalf("Error retrieving categories: %v", err)
	}

	// Check if updated category name exists in list of categories
	found := false
	for _, c := range categories {
		if c == newCategoryName {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected updated category name to exist in list of categories")
	}
}

// TestDeleteCategory tests DeleteCategory function
func TestDeleteCategory(t *testing.T) {
	// Open in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create necessary tables
	createCategoryTable(db)

	// Create new category
	category := &models.Category{Name: "Test Category"}

	// Add category to database
	_, err = models.AddCategory(db, category)
	if err != nil {
		t.Fatalf("Error adding category: %v", err)
	}

	// Delete category from database
	err = models.DeleteCategory(db, category.Name)
	if err != nil {
		t.Fatalf("Error deleting category: %v", err)
	}

	// Retrieve all categories from database
	categories, err := models.GetAllCategories(db)
	if err != nil {
		t.Fatalf("Error retrieving categories: %v", err)
	}

	// Check if deleted category name exists in list of categories
	found := false
	for _, c := range categories {
		if c == category.Name {
			found = true
			break
		}
	}
	if found {
		t.Error("Expected deleted category name not to exist in list of categories")
	}
}

// createCategoryTable creates necessary tables in database
func createCategoryTable(db *sql.DB) {
	createCategoriesTableQuery := `
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
		)
	`

	_, err := db.Exec(createCategoriesTableQuery)
	if err != nil {
		panic(err)
	}
}
