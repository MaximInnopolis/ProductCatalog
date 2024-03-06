package models

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

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
