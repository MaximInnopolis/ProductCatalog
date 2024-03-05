package models

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// func GetAllCategories(db *sql.DB) ([]Category, error) {
func GetAllCategories(db *sql.DB) ([]string, error) {
	//query := "SELECT id, name FROM categories"
	query := "SELECT name FROM categories"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//var categories []Category
	var categories []string
	for rows.Next() {
		//var category Category
		var category string
		//if err := rows.Scan(&category.ID, &category.Name); err != nil {
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	logger.Println("Successfully got all categories")
	return categories, nil
}
