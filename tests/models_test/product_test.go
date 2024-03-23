package models_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/MaximInnopolis/ProductCatalog/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

func TestAddProduct(t *testing.T) {
	ctx := context.WithValue(context.Background(), "endpoint", "products/new")
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	createTables(db)

	product := &model.Product{Name: "Test Product"}
	categories := []model.Category{{Name: "Test Category"}}

	err = model.AddProduct(ctx, db, product, categories)
	if err != nil {
		t.Fatalf("Error adding product: %v", err)
	}

	// Retrieve product ID after adding product
	productID, err := model.GetProductID(ctx, db, product.Name)
	if err != nil {
		t.Fatalf("Error getting product ID: %v", err)
	}

	// Check if categories were associated with product
	associatedCategories, err := model.GetCategoriesByProductID(ctx, db, productID)
	if err != nil {
		t.Fatalf("Error getting categories associated with product: %v", err)
	}
	if len(associatedCategories) != 1 || associatedCategories[0] != categories[0].Name {
		t.Errorf("Expected categories to be associated with product: %v, got: %v", categories, associatedCategories)
	}
}

func TestUpdateProduct(t *testing.T) {
	ctx := context.WithValue(context.Background(), "endpoint", "products")

	// Open in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create tables
	createTables(db)

	// Create test categories
	testCategories := []model.Category{
		{Name: "Category A"},
		{Name: "Category B"},
	}

	// Add test categories to database
	for _, category := range testCategories {
		_, err := model.AddCategory(ctx, db, &category)
		if err != nil {
			t.Fatalf("Error adding category: %v", err)
		}
	}

	// Create test product
	testProduct := &model.Product{Name: "Test Product"}

	// Add test product to database
	err = model.AddProduct(ctx, db, testProduct, testCategories)
	if err != nil {
		t.Fatalf("Error adding product: %v", err)
	}

	// Get product ID
	productID, err := model.GetProductID(ctx, db, testProduct.Name)
	if err != nil {
		t.Fatalf("Error getting product ID: %v", err)
	}

	// Update test product
	updatedProduct := &model.Product{ID: int(productID), Name: "Test Product"}
	err = model.UpdateProduct(ctx, db, updatedProduct, testCategories)
	if err != nil {
		t.Fatalf("Error updating product: %v", err)
	}

	if updatedProduct.Name != "Test Product" {
		t.Errorf("Expected product name to be 'Test Product', got '%s'", updatedProduct.Name)
	}
}

func TestDeleteProduct(t *testing.T) {
	ctx := context.WithValue(context.Background(), "endpoint", "products/")

	// Open in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create tables
	createTables(db)

	// Create test product
	product := &model.Product{Name: "Test Product"}
	categories := []model.Category{{Name: "Test Category"}}

	// Add product
	err = model.AddProduct(ctx, db, product, categories)
	if err != nil {
		t.Fatalf("Error adding product: %v", err)
	}

	// Delete product
	err = model.DeleteProduct(ctx, db, product.Name)
	if err != nil {
		t.Fatalf("Error deleting product: %v", err)
	}

	// Check if product was deleted successfully
	_, err = model.GetProductID(ctx, db, product.Name)
	if err == nil {
		t.Error("Expected product to be deleted, but it still exists in database")
	}
}

func TestContains(t *testing.T) {
	// Test case 1: string exists in slice
	slice1 := []string{"apple", "banana", "orange"}
	str1 := "banana"
	if !model.Contains(slice1, str1) {
		t.Errorf("Test case 1 failed: expected true, got false")
	}

	// Test case 2: string does not exist in slice
	slice2 := []string{"apple", "banana", "orange"}
	str2 := "grape"
	if model.Contains(slice2, str2) {
		t.Errorf("Test case 2 failed: expected false, got true")
	}

	// Test case 3: empty slice
	var slice3 []string
	str3 := "apple"
	if model.Contains(slice3, str3) {
		t.Errorf("Test case 3 failed: expected false, got true")
	}

	// Test case 4: empty string
	slice4 := []string{"apple", "banana", "orange"}
	str4 := ""
	if model.Contains(slice4, str4) {
		t.Errorf("Test case 4 failed: expected false, got true")
	}
}

func TestGetCategoryID(t *testing.T) {
	ctx := context.WithValue(context.Background(), "endpoint", "products/")
	// Open in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create tables
	createTables(db)

	// Insert category into database
	categoryName := "Test Category"
	categoryID, err := model.AddCategory(ctx, db, &model.Category{Name: categoryName})
	if err != nil {
		t.Fatalf("Error adding category: %v", err)
	}

	// Retrieve ID of inserted category
	retrievedCategoryID, err := model.GetCategoryID(ctx, db, categoryName)
	if err != nil {
		t.Fatalf("Error retrieving category ID: %v", err)
	}

	if categoryID != retrievedCategoryID {
		t.Errorf("Expected category ID: %d, got: %d", categoryID, retrievedCategoryID)
	}
}

func createTables(db *sql.DB) {
	createCategoriesTableQuery := `
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
		)
	`

	createProductsTableQuery := `
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
		)
	`

	createProductCategoriesTableQuery := `
		CREATE TABLE IF NOT EXISTS product_categories (
			product_id INTEGER,
			category_id INTEGER,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
		)
	`

	_, err := db.Exec(createCategoriesTableQuery)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createProductsTableQuery)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createProductCategoriesTableQuery)
	if err != nil {
		panic(err)
	}
}
