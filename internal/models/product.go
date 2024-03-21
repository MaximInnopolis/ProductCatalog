package models

import (
	"context"
	"database/sql"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
)

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// AddProduct inserts new product into 'products' table in database and associates with specified categories
// takes database connection, pointer to product information, and slice of categories as parameters
// returns error if any occurred during insertion process
func AddProduct(ctx context.Context, db *sql.DB, product *Product, categories []Category) error {
	// If product does not have any categories, return error
	if len(categories) == 0 {
		logger.Printf(ctx, "Error adding product: product must have at least one category")
		return errors.New("product must have at least one category")
	}

	// Insert product into 'products' table
	query := "INSERT INTO products (name) VALUES (?)"
	result, err := db.Exec(query, product.Name)
	if err != nil {
		logger.Printf(ctx, "Error inserting product into database: %s", err)
		return err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		logger.Printf(ctx, "Error getting last insert ID: %s", err)
		return err
	}

	// Insert associations into 'product_categories' table
	for _, category := range categories {
		categoryID, err := GetCategoryID(ctx, db, category.Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {

				// If category doesn't exist, add to 'categories' table
				categoryID, err = AddCategory(ctx, db, &category)
				if err != nil {
					logger.Printf(ctx, "Error adding category: %s", err)
					return err
				}
			} else {
				logger.Printf(ctx, "Error getting category ID: %s", err)
				return err
			}
		}
		// Insert association into 'product_categories' table
		query = "INSERT INTO product_categories (product_id, category_id) VALUES (?, ?)"
		_, err = db.Exec(query, productID, categoryID)
		if err != nil {
			logger.Printf(ctx, "Error inserting association into product_categories: %s", err)
			return err
		}
	}

	logger.Printf(ctx, "Product added successfully")
	return nil
}

// UpdateProduct edits existing product in database along with its associated categories
// takes database connection, pointer to updated product information, and slice of updated categories as parameters
// returns error if any occurred during update process
func UpdateProduct(ctx context.Context, db *sql.DB, product *Product, categories []Category) error {
	// If product does not have any categories, return error
	if len(categories) == 0 {
		logger.Printf(ctx, "Error adding product: product must have at least one category")
		return errors.New("product must have at least one category")
	}

	// Check if product exists in database
	productID, err := GetProductID(ctx, db, product.Name)
	if err != nil {
		return err
	}

	// Get current categories associated with product
	currentCategories, err := GetCategoriesByProductID(ctx, db, productID)
	if err != nil {
		return err
	}

	// Map to store category names to their IDs
	categoryIDMap := make(map[string]int64)

	// Insert new categories into 'categories' table if they don't exist yet
	for _, category := range categories {
		categoryID, err := GetCategoryID(ctx, db, category.Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// If category doesn't exist, add to 'categories' table
				categoryID, err = AddCategory(ctx, db, &category)
				if err != nil {
					return err
				}
			} else {
				logger.Printf(ctx, "Error getting category ID: %s", err)
				return err
			}
		}
		categoryIDMap[category.Name] = categoryID
	}

	// Remove categories that are no longer associated with product
	for _, currentCategory := range currentCategories {
		if _, exists := categoryIDMap[currentCategory]; !exists {
			if err := DeleteProductCategory(ctx, db, productID, currentCategory); err != nil {
				return err
			}
		}
	}

	// Insert new associations into 'product_categories' table
	for _, category := range categories {
		categoryID := categoryIDMap[category.Name]
		if !Contains(currentCategories, category.Name) {
			query := "INSERT INTO product_categories (product_id, category_id) VALUES (?, ?)"
			_, err = db.Exec(query, productID, categoryID)
			if err != nil {
				logger.Printf(ctx, "Error inserting association into product_categories: %s", err)
				return err
			}
		}
	}

	logger.Printf(ctx, "Product updated successfully")
	return nil
}

// DeleteProduct deletes specified product from database along with its associated records in 'product_categories' table
// takes database connection and name of product to be deleted as parameters
// returns error if any occurred during deletion process
func DeleteProduct(ctx context.Context, db *sql.DB, productName string) error {

	// Check if product exists in database
	productID, err := GetProductID(ctx, db, productName)
	if err != nil {
		return err
	}

	// Delete product from 'products' table
	query := "DELETE FROM products WHERE id = ?"
	result, err := db.Exec(query, productID)
	if err != nil {
		logger.Printf(ctx, "Error deleting product from database: %s", err)
		return err
	}

	// Get number of rows affected by delete operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Printf(ctx, "Error getting rows affected: %s", err)
		return err
	}

	// If no rows affected, log message indicating that product was not found in database
	if rowsAffected == 0 {
		logger.Printf(ctx, "Product %s not found in database", productName)
		return errors.New("no rows affected, product not found")
	}

	// Delete associated records from 'product_categories' table
	query = "DELETE FROM product_categories WHERE product_id = ?"
	_, err = db.Exec(query, productID)
	if err != nil {
		logger.Printf(ctx, "Error deleting product associations from database: %s", err)
		return err
	}

	logger.Printf(ctx, "Product deleted successfully")
	return nil
}

// GetProductsByCategory retrieves names of products belonging to specified category from database
// takes database connection and name of category as parameters
// returns slice of strings containing names of products in specified category and any error encountered
func GetProductsByCategory(ctx context.Context, db *sql.DB, categoryName string) ([]string, error) {
	query := `
        SELECT p.name
        FROM products p
        JOIN product_categories pc ON p.id = pc.product_id
        JOIN categories c ON pc.category_id = c.id
        WHERE c.name = ?
    `

	// Execute database query to retrieve products in specified category
	rows, err := db.Query(query, categoryName)
	if err != nil {
		logger.Printf(ctx, "Error executing database query: %s", err)
		return nil, err
	}
	defer rows.Close()

	// Initialize slice to store names of products in category
	var products []string
	// Iterate through query result rows
	for rows.Next() {
		var product string
		// Scan product name from current row
		if err := rows.Scan(&product); err != nil {
			logger.Printf(ctx, "Error scanning row from query result: %s", err)
			return nil, err
		}
		// Append product name to products slice
		products = append(products, product)
	}
	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		logger.Printf(ctx, "Error processing query result: %s", err)
		return nil, err
	}

	logger.Printf(ctx, "Successfully got products in category %s", categoryName)
	// Return slice of product names and nil error, indicating success
	return products, nil
}
