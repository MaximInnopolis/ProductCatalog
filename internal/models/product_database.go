package models

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
)

// GetCategoryID retrieves ID of category with specified name from database
// takes database connection and category name as parameters
// returns category ID and any error encountered
func GetCategoryID(db *sql.DB, categoryName string) (int64, error) {
	// Construct SQL query to select category ID based on category name
	query := "SELECT id FROM categories WHERE name = ?"
	// Execute query and retrieve single row result
	row := db.QueryRow(query, categoryName)

	// Initialize variable to store category ID
	var categoryID int64
	// Scan category ID from result row into categoryID variable
	err := row.Scan(&categoryID)
	if err != nil {
		logger.Printf("Category %v does not exist in database yet. Creating...", categoryName)
		return 0, err
	}

	// Return retrieved category ID and nil error, indicating success
	return categoryID, nil
}

// GetProductID retrieves ID of product with specified name from database
// takes database connection and product name as parameters
// returns product ID and any error encountered
func GetProductID(db *sql.DB, productName string) (int64, error) {
	// Construct SQL query to select product ID based on product name
	query := "SELECT id FROM products WHERE name = ?"
	// Execute query and retrieve single row result
	row := db.QueryRow(query, productName)

	// Initialize variable to store product ID
	var productID int64
	// Scan product ID from result row into productID variable
	err := row.Scan(&productID)
	if err != nil {
		logger.Printf("Product %v does not exist in database.", productName)
		return 0, err
	}

	// Return retrieved product ID and nil error, indicating success
	return productID, nil
}

// DeleteProductCategory deletes association between specified product and category from database
// takes database connection, product ID, and category name as parameters
// returns error if any occurred during deletion process
func DeleteProductCategory(db *sql.DB, productID int64, categoryName string) error {
	// Construct SQL query to delete association between product and category
	query := `
		DELETE FROM product_categories
		WHERE product_id = ? AND category_id = (SELECT id FROM categories WHERE name = ?)
	`

	// Execute delete query with provided product ID and category name
	_, err := db.Exec(query, productID, categoryName)
	if err != nil {
		logger.Println("Error deleting product category association from database:", err)
		return err
	}

	logger.Println("Product category association deleted successfully")
	return nil
}

// GetCategoriesByProductID retrieves categories associated with specified product from database
// takes database connection and product ID as parameters
// returns slice of category names and any error encountered
func GetCategoriesByProductID(db *sql.DB, productID int64) ([]string, error) {
	// Construct SQL query to select categories associated with given product ID
	query := `
		SELECT c.name
		FROM categories c
		JOIN product_categories pc ON c.id = pc.category_id
		WHERE pc.product_id = ?
	`

	// Execute query to retrieve categories associated with product ID
	rows, err := db.Query(query, productID)
	if err != nil {
		logger.Println("Error executing database query:", err)
		return nil, err
	}
	defer rows.Close()

	// Initialize slice to store category names
	var categories []string
	// Iterate through query result rows
	for rows.Next() {
		var category string
		// Scan category name from current row
		if err := rows.Scan(&category); err != nil {
			logger.Println("Error scanning row from query result:", err)
			return nil, err
		}
		// Append category name to categories slice
		categories = append(categories, category)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		logger.Println("Error processing query result:", err)
		return nil, err
	}

	// Return slice of category names and nil error, indicating success
	return categories, nil
}
