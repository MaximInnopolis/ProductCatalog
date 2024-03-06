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

// AddCategory create new category in database
func AddCategory(db *sql.DB, category *Category) error {
	query := "INSERT INTO categories (name) VALUES (?)"
	_, err := db.Exec(query, category.Name)
	if err != nil {
		logger.Println("Error inserting category into database:", err)
		return err
	}

	logger.Println("Category added successfully")
	return nil
}

// UpdateCategory edit existing category in database
func UpdateCategory(db *sql.DB, categoryName string, category *Category) error {
	query := "UPDATE categories SET name = ? WHERE name = ?"
	result, err := db.Exec(query, category.Name, categoryName)
	if err != nil {
		logger.Println("Error updating category in database:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Println("Error getting rows affected:", err)
		return err
	}

	if rowsAffected == 0 {
		logger.Printf("Category %s not found in database", categoryName)
		return errors.New("no rows affected, category not found")
	}

	logger.Println("Category updated successfully")
	return nil
}

// GetAllCategories returns all existing categories
func GetAllCategories(db *sql.DB) ([]string, error) {
	query := "SELECT name FROM categories"
	rows, err := db.Query(query)
	if err != nil {
		logger.Println("Error querying categories:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			logger.Println("Error scanning category row:", err)
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		logger.Println("Error iterating through category rows:", err)
		return nil, err
	}

	logger.Println("Successfully got all categories")
	return categories, nil
}

// DeleteCategory deletes the specified category from the database
func DeleteCategory(db *sql.DB, categoryName string) error {
	query := "DELETE FROM categories WHERE name = ?"
	result, err := db.Exec(query, categoryName)
	if err != nil {
		logger.Println("Error deleting category from database:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Println("Error getting rows affected:", err)
		return err
	}

	if rowsAffected == 0 {
		logger.Printf("Category %s not found in database", categoryName)
		return errors.New("no rows affected, category not found")
	}

	logger.Println("Category deleted successfully")
	return nil
}
