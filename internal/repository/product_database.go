package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
)

type ProductDatabase struct {
	CategoryDatabase
}

func NewProductDatabase(db *sql.DB) *ProductDatabase {
	return &ProductDatabase{CategoryDatabase{db: db}}
}

// CreateProduct inserts new product into 'products' table in database and associates with specified categories
// takes database connection, pointer to product information, and slice of categories as parameters
// returns error if any occurred during insertion process
func (r *ProductDatabase) CreateProduct(ctx context.Context, product *model.Product, categories []model.Category) error {
	if len(categories) == 0 {
		logger.Printf(ctx, "Error adding product: product must have at least one category")
		return errors.New("product must have at least one category")
	}

	// Insert product into 'products' table
	query := "INSERT INTO products (name) VALUES (?)"
	result, err := r.db.Exec(query, product.Name)
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
		categoryID, err := r.GetCategoryID(ctx, category.Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {

				// If category doesn't exist, add to 'categories' table
				categoryID, err = r.CreateCategory(ctx, &category)
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
		_, err = r.db.Exec(query, productID, categoryID)
		if err != nil {
			logger.Printf(ctx, "Error inserting association into product_categories: %s", err)
			return err
		}
	}

	logger.Printf(ctx, "Product added successfully")
	return nil
}

// GetProductsByCategory retrieves names of products belonging to specified category from database
// takes database connection and name of category as parameters
// returns slice of strings containing names of products in specified category and any error encountered
func (r *ProductDatabase) GetProductsByCategory(ctx context.Context, categoryName string) ([]string, error) {
	query := `
       SELECT p.name
       FROM products p
       JOIN product_categories pc ON p.id = pc.product_id
       JOIN categories c ON pc.category_id = c.id
       WHERE c.name = ?
   `

	// Execute database query to retrieve products in specified category
	rows, err := r.db.Query(query, categoryName)
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

// UpdateProduct edits existing product in database along with its associated categories
// takes database connection, pointer to updated product information, and slice of updated categories as parameters
// returns error if any occurred during update process
func (r *ProductDatabase) UpdateProduct(ctx context.Context, product *model.Product, categories []model.Category) error {
	// If product does not have any categories, return error
	if len(categories) == 0 {
		logger.Printf(ctx, "Error adding product: product must have at least one category")
		return errors.New("product must have at least one category")
	}

	// Check if product exists in database
	productID, err := r.GetProductID(ctx, product.Name)
	if err != nil {
		return err
	}

	// Get current categories associated with product
	currentCategories, err := r.GetCategoriesByProductID(ctx, productID)
	if err != nil {
		return err
	}

	// Map to store category names to their IDs
	categoryIDMap := make(map[string]int64)

	// Insert new categories into 'categories' table if they don't exist yet
	for _, category := range categories {
		categoryID, err := r.GetCategoryID(ctx, category.Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// If category doesn't exist, add to 'categories' table
				categoryID, err = r.CreateCategory(ctx, &category)
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
			if err := r.DeleteProductCategory(ctx, productID, currentCategory); err != nil {
				return err
			}
		}
	}

	// Insert new associations into 'product_categories' table
	for _, category := range categories {
		categoryID := categoryIDMap[category.Name]
		if !contains(currentCategories, category.Name) {
			query := "INSERT INTO product_categories (product_id, category_id) VALUES (?, ?)"
			_, err = r.db.Exec(query, productID, categoryID)
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
func (r *ProductDatabase) DeleteProduct(ctx context.Context, productID int64) error {

	// Delete product from 'products' table
	query := "DELETE FROM products WHERE id = ?"
	result, err := r.db.Exec(query, productID)
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
		return errors.New("no rows affected, product not found")
	}

	// Delete associated records from 'product_categories' table
	query = "DELETE FROM product_categories WHERE product_id = ?"
	_, err = r.db.Exec(query, productID)
	if err != nil {
		logger.Printf(ctx, "Error deleting product associations from database: %s", err)
		return err
	}

	logger.Printf(ctx, "Product deleted successfully")
	return nil
}

// GetProductID retrieves ID of product with specified name from database
// takes database connection and product name as parameters
// returns product ID and any error encountered
func (r *ProductDatabase) GetProductID(ctx context.Context, productName string) (int64, error) {
	// Construct SQL query to select product ID based on product name
	query := "SELECT id FROM products WHERE name = ?"
	// Execute query and retrieve single row result
	row := r.db.QueryRow(query, productName)

	// Initialize variable to store product ID
	var productID int64
	// Scan product ID from result row into productID variable
	err := row.Scan(&productID)
	if err != nil {
		logger.Printf(ctx, "Product %v not found in database.", productName)
		return 0, err
	}

	// Return retrieved product ID and nil error, indicating success
	return productID, nil
}

// GetCategoriesByProductID retrieves categories associated with specified product from database
// takes database connection and product ID as parameters
// returns slice of category names and any error encountered
func (r *ProductDatabase) GetCategoriesByProductID(ctx context.Context, productID int64) ([]string, error) {
	// Construct SQL query to select categories associated with given product ID
	query := `
		SELECT c.name
		FROM categories c
		JOIN product_categories pc ON c.id = pc.category_id
		WHERE pc.product_id = ?
	`

	// Execute query to retrieve categories associated with product ID
	rows, err := r.db.Query(query, productID)
	if err != nil {
		logger.Printf(ctx, "Error executing database query: %s", err)
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
			logger.Printf(ctx, "Error scanning row from query result: %s", err)
			return nil, err
		}
		// Append category name to categories slice
		categories = append(categories, category)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		logger.Printf(ctx, "Error processing query result: %s", err)
		return nil, err
	}

	// Return slice of category names and nil error, indicating success
	return categories, nil
}

// DeleteProductCategory deletes association between specified product and category from database
// takes database connection, product ID, and category name as parameters
// returns error if any occurred during deletion process
func (r *ProductDatabase) DeleteProductCategory(ctx context.Context, productID int64, categoryName string) error {
	// Construct SQL query to delete association between product and category
	query := `
		DELETE FROM product_categories
		WHERE product_id = ? AND category_id = (SELECT id FROM categories WHERE name = ?)
	`

	// Execute delete query with provided product ID and category name
	_, err := r.db.Exec(query, productID, categoryName)
	if err != nil {
		logger.Printf(ctx, "Error deleting product category association from database: %s", err)
		return err
	}

	logger.Printf(ctx, "Product category association deleted successfully")
	return nil
}

// contains checks if string exists in slice of strings
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
