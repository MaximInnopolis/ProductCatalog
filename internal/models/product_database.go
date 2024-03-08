package models

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
)

// GetCategoryID retrieves ID of category from database
func GetCategoryID(db *sql.DB, categoryName string) (int64, error) {
	query := "SELECT id FROM categories WHERE name = ?"
	row := db.QueryRow(query, categoryName)

	var categoryID int64
	err := row.Scan(&categoryID)
	if err != nil {
		logger.Printf("Category %v does not exist in database yet. Creating...", categoryName)
		return 0, err
	}

	return categoryID, nil
}

// GetProductID retrieves ID of product from database
func GetProductID(db *sql.DB, productName string) (int64, error) {
	query := "SELECT id FROM products WHERE name = ?"
	row := db.QueryRow(query, productName)

	var productID int64
	err := row.Scan(&productID)
	if err != nil {
		logger.Printf("Product %v does not exist in database.", productName)
		return 0, err
	}

	return productID, nil
}

func DeleteProductCategory(db *sql.DB, productID int64, categoryName string) error {
	query := `
		DELETE FROM product_categories
		WHERE product_id = ? AND category_id = (SELECT id FROM categories WHERE name = ?)
	`

	_, err := db.Exec(query, productID, categoryName)
	if err != nil {
		logger.Println("Error deleting product category association from database:", err)
		return err
	}

	logger.Println("Product category association deleted successfully")
	return nil
}

// GetCategoriesByProductID retrieves categories associated with product from database
func GetCategoriesByProductID(db *sql.DB, productID int64) ([]string, error) {
	query := `
		SELECT c.name
		FROM categories c
		JOIN product_categories pc ON c.id = pc.category_id
		WHERE pc.product_id = ?
	`

	rows, err := db.Query(query, productID)
	if err != nil {
		logger.Println("Error executing database query:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			logger.Println("Error scanning row from query result:", err)
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		logger.Println("Error processing query result:", err)
		return nil, err
	}

	return categories, nil
}
