package models

import (
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
func AddCategory(db *sql.DB, category *Category) (int64, error) {
	// Execute INSERT query to add new category to database with provided name
	query := "INSERT INTO categories (name) VALUES (?)"
	result, err := db.Exec(query, category.Name)
	if err != nil {
		logger.Println("Error inserting category into database:", err)
		return 0, err
	}

	// Retrieve ID of newly inserted category
	categoryID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	logger.Println("Category added successfully")
	// Return ID of newly inserted category and nil error
	return categoryID, nil
}

// GetAllCategories retrieves all existing categories from database using provided DB connection
// Returns slice of strings containing names of all categories and any error encountered
func GetAllCategories(db *sql.DB) ([]string, error) {
	// Execute SELECT query to retrieve all category names from database
	query := "SELECT name FROM categories"
	rows, err := db.Query(query)
	if err != nil {
		logger.Println("Error querying categories:", err)
		return nil, err
	}
	defer rows.Close()

	// Initialize slice to store category names retrieved from database
	var categories []string
	// Iterate through each row of result set
	for rows.Next() {
		var category string
		// Scan category name from current row into 'category' variable
		if err := rows.Scan(&category); err != nil {
			logger.Println("Error scanning category row:", err)
			return nil, err
		}
		// Append category name to 'categories' slice
		categories = append(categories, category)
	}
	// Check for any error occurred during iteration
	if err := rows.Err(); err != nil {
		logger.Println("Error iterating through category rows:", err)
		return nil, err
	}

	logger.Println("Successfully got all categories")
	// Return slice containing all category names and nil error
	return categories, nil
}

// UpdateCategory edits existing category in database with specified name
// Takes database connection, current category name, and updated category information as parameters
// Returns error if any occurred during update process
func UpdateCategory(db *sql.DB, categoryName string, category *Category) error {
	// Execute UPDATE query to update category name in database
	query := "UPDATE categories SET name = ? WHERE name = ?"
	result, err := db.Exec(query, category.Name, categoryName)
	if err != nil {
		logger.Println("Error updating category in database:", err)
		return err
	}

	// Get number of rows affected by UPDATE operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Println("Error getting rows affected:", err)
		return err
	}

	// Check if no rows were affected by UPDATE operation
	if rowsAffected == 0 {
		logger.Printf("Category %s not found in database", categoryName)
		return errors.New("no rows affected, category not found")
	}

	logger.Println("Category updated successfully")
	// Return error indicating that no rows were affected
	return nil
}

// DeleteCategory deletes specified category from database
// Takes database connection and name of category to be deleted as parameters
// Returns error if any occurred during deletion process
func DeleteCategory(db *sql.DB, categoryName string) error {
	// Execute DELETE query to delete category from database
	query := "DELETE FROM categories WHERE name = ?"
	result, err := db.Exec(query, categoryName)
	if err != nil {
		logger.Println("Error deleting category from database:", err)
		return err
	}

	// Get number of rows affected by DELETE operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Println("Error getting rows affected:", err)
		// Return error indicating that no rows were affected
		return err
	}

	// Check if no rows were affected by DELETE operation
	if rowsAffected == 0 {
		logger.Printf("Category %s not found in database", categoryName)
		return errors.New("no rows affected, category not found")
	}

	logger.Println("Category deleted successfully")
	// Return nil to indicate success and no error
	return nil
}
