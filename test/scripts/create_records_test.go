package scripts

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/scripts"
	"testing"
)

func TestCreateRecords(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY,
		name TEXT UNIQUE
	);`)
	if err != nil {
		t.Fatalf("Error creating products table': %v", err)
	}

	_, err = db.Exec(`CREATE TABLE categories (
		id INTEGER PRIMARY KEY,
		name TEXT UNIQUE
	);`)
	if err != nil {
		t.Fatalf("Error creating 'categories' table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE product_categories (
		product_id INTEGER,
		category_id INTEGER,
		FOREIGN KEY(product_id) REFERENCES products(id),
		FOREIGN KEY(category_id) REFERENCES categories(id),
		PRIMARY KEY(product_id, category_id)
	);`)
	if err != nil {
		t.Fatalf("Error creating 'product_categories' table: %v", err)
	}

	err = scripts.CreateRecords(db)
	if err != nil {
		t.Fatalf("CreateRecords returned error: %v", err)
	}

	var productCount, categoryCount, productCategoryCount int
	err = db.QueryRow("SELECT COUNT(*) FROM products").Scan(&productCount)
	if err != nil {
		t.Fatalf("Error getting product count: %v", err)
	}
	if productCount != 1 {
		t.Fatalf("Expected 1 product, got %d", productCount)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&categoryCount)
	if err != nil {
		t.Fatalf("Error getting category count: %v", err)
	}
	if categoryCount != 2 {
		t.Fatalf("Expected 2 categories, got %d", categoryCount)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM product_categories").Scan(&productCategoryCount)
	if err != nil {
		t.Fatalf("Error getting productCategoryCount count: %v", err)
	}
	if productCategoryCount != 2 {
		t.Fatalf("Expected 2 productCategoryCount, got %d", productCategoryCount)
	}
}
