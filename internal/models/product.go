package models

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
)

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetProductsByCategory(db *sql.DB, categoryName string) ([]string, error) {
	query := `
        SELECT p.name
        FROM products p
        JOIN product_categories pc ON p.id = pc.product_id
        JOIN categories c ON pc.category_id = c.id
        WHERE c.name = ?
    `

	rows, err := db.Query(query, categoryName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []string
	for rows.Next() {
		var product string
		if err := rows.Scan(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	logger.Println("Successfully got products in category", categoryName)
	return products, nil
}
