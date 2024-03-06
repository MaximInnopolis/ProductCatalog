package models

import (
	"database/sql"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
)

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// AddProduct creates a new product in the database and associates it with the specified categories
func AddProduct(db *sql.DB, product *Product, categories []Category) error {
	// Insert product into products table
	query := "INSERT INTO products (name) VALUES (?)"
	result, err := db.Exec(query, product.Name)
	if err != nil {
		logger.Println("Error inserting product into database:", err)
		return err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		logger.Println("Error getting last insert ID:", err)
		return err
	}

	// Insert associations into product_categories table
	for _, category := range categories {
		categoryID, err := getCategoryID(db, category.Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {

				// If category doesn't exist, add it to table 'categories'
				categoryID, err = AddCategory(db, &category)
				if err != nil {
					logger.Println("Error adding category:", err)
					return err
				}
			} else {
				logger.Println("Error getting category ID:", err)
				return err
			}
		}

		query = "INSERT INTO product_categories (product_id, category_id) VALUES (?, ?)"
		_, err = db.Exec(query, productID, categoryID)
		if err != nil {
			logger.Println("Error inserting association into product_categories:", err)
			return err
		}
	}

	logger.Println("Product added successfully")
	return nil
}

// getCategoryID retrieves ID of category from database
func getCategoryID(db *sql.DB, categoryName string) (int64, error) {
	query := "SELECT id FROM categories WHERE name = ?"
	row := db.QueryRow(query, categoryName)

	var categoryID int64
	err := row.Scan(&categoryID)
	if err != nil {
		logger.Printf("Category %v does not exist in database yet. Creating...:", categoryName)
		return 0, err
	}

	return categoryID, nil
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
		logger.Println("Error executing database query:", err)
		return nil, err
	}
	defer rows.Close()

	var products []string
	for rows.Next() {
		var product string
		if err := rows.Scan(&product); err != nil {
			logger.Println("Error scanning row from query result:", err)
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		logger.Println("Error processing query result:", err)
		return nil, err
	}

	logger.Println("Successfully got products in category", categoryName)
	return products, nil
}
