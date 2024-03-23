package model

import (
	"context"
	"database/sql"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// AddCategory inserts new category into database using provided DB connection
// Returns ID of newly inserted category and any error encountered
func AddCategory(ctx context.Context, db *sql.DB, category *Category) (int64, error) {
	// Execute INSERT query to add new category to database with provided name
	query := "INSERT INTO categories (name) VALUES (?)"
	result, err := db.Exec(query, category.Name)
	if err != nil {
		logger.Printf(ctx, "Error inserting category into database: %s", err)
		return 0, err
	}

	// Retrieve ID of newly inserted category
	categoryID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	logger.Printf(ctx, "Category added successfully")
	// Return ID of newly inserted category and nil error
	return categoryID, nil
}

// GetAllCategories retrieves all existing categories from database using provided DB connection
// Returns slice of strings containing names of all categories and any error encountered
func GetAllCategories(ctx context.Context, db *sql.DB) ([]string, error) {
	// Execute SELECT query to retrieve all category names from database
	query := "SELECT name FROM categories"
	rows, err := db.Query(query)
	if err != nil {
		logger.Printf(ctx, "Error querying categories: %s", err)
		return nil, err
	}
	defer rows.Close()

	// Initialize slice to store category names retrieved from database
	var categories []string
	// Iterate through each row of result set
	for rows.Next() {
		var category string
		// Scan category name from current row into 'category' variable
		if err = rows.Scan(&category); err != nil {
			logger.Printf(ctx, "Error scanning category row: %s", err)
			return nil, err
		}
		// Append category name to 'categories' slice
		categories = append(categories, category)
	}
	// Check for any error occurred during iteration
	if err = rows.Err(); err != nil {
		logger.Printf(ctx, "Error iterating through category rows: %s", err)
		return nil, err
	}

	logger.Printf(ctx, "Successfully got all categories")
	// Return slice containing all category names and nil error
	return categories, nil
}

// UpdateCategory edits existing category in database with specified name
// Takes database connection, current category name, and updated category information as parameters
// Returns error if any occurred during update process
func UpdateCategory(ctx context.Context, db *sql.DB, categoryName string, category *Category) error {
	// Execute UPDATE query to update category name in database
	query := "UPDATE categories SET name = ? WHERE name = ?"
	result, err := db.Exec(query, category.Name, categoryName)
	if err != nil {
		logger.Printf(ctx, "Error updating category in database: %s", err)
		return err
	}

	// Get number of rows affected by UPDATE operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Printf(ctx, "Error getting rows affected: %s", err)
		return err
	}

	// Check if no rows were affected by UPDATE operation
	if rowsAffected == 0 {
		logger.Printf(ctx, "Category %s not found in database", categoryName)
		return errors.New("no rows affected, category not found")
	}

	logger.Printf(ctx, "Category updated successfully")
	// Return error indicating that no rows were affected
	return nil
}

// DeleteCategory deletes specified category from database
// Takes database connection and name of category to be deleted as parameters
// Returns error if any occurred during deletion process
func DeleteCategory(ctx context.Context, db *sql.DB, categoryName string) error {
	// Execute DELETE query to delete category from database
	query := "DELETE FROM categories WHERE name = ?"
	result, err := db.Exec(query, categoryName)
	if err != nil {
		logger.Printf(ctx, "Error deleting category from database: %s", err)
		return err
	}

	// Get number of rows affected by DELETE operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Printf(ctx, "Error getting rows affected: %s", err)
		// Return error indicating that no rows were affected
		return err
	}

	// Check if no rows were affected by DELETE operation
	if rowsAffected == 0 {
		logger.Printf(ctx, "Category %s not found in database", categoryName)
		return errors.New("no rows affected, category not found")
	}

	logger.Printf(ctx, "Category deleted successfully")
	// Return nil to indicate success and no error
	return nil
}
