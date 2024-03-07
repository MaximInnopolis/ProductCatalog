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

// AddProduct creates new product in database and associates it with specified categories
func AddProduct(db *sql.DB, product *Product, categories []Category) error {

	if len(categories) == 0 {
		logger.Println("Error adding product: product must have at least one category")
		return errors.New("product must have at least one category")
	}

	// Insert product into 'products' table
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

	// Insert associations into 'product_categories' table
	for _, category := range categories {
		categoryID, err := getCategoryID(db, category.Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {

				// If category doesn't exist, add it to 'categories' table
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

// UpdateProduct edits existing product in database
func UpdateProduct(db *sql.DB, product *Product, categories []Category) error {

	if len(categories) == 0 {
		logger.Println("Error adding product: product must have at least one category")
		return errors.New("product must have at least one category")
	}

	// Check if product exists in database
	productID, err := getProductID(db, product.Name)
	if err != nil {
		return err
	}

	// Get current categories associated with product
	currentCategories, err := getCategoriesByProductID(db, productID)
	if err != nil {
		return err
	}

	// Map to store category names to their IDs
	categoryIDMap := make(map[string]int64)

	// Insert new categories into 'categories' table if they don't exist yet
	for _, category := range categories {
		categoryID, err := getCategoryID(db, category.Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// If category doesn't exist, add it to 'categories' table
				categoryID, err = AddCategory(db, &category)
				if err != nil {
					return err
				}
			} else {
				logger.Println("Error getting category ID:", err)
				return err
			}
		}
		categoryIDMap[category.Name] = categoryID
	}

	// Remove categories that are no longer associated with product
	for _, currentCategory := range currentCategories {
		if _, exists := categoryIDMap[currentCategory]; !exists {
			if err := deleteProductCategory(db, productID, currentCategory); err != nil {
				return err
			}
		}
	}

	// Insert new associations into 'product_categories' table
	for _, category := range categories {
		categoryID := categoryIDMap[category.Name]
		if !contains(currentCategories, category.Name) {
			query := "INSERT INTO product_categories (product_id, category_id) VALUES (?, ?)"
			_, err = db.Exec(query, productID, categoryID)
			if err != nil {
				logger.Println("Error inserting association into product_categories:", err)
				return err
			}
		}
	}

	logger.Println("Product updated successfully")
	return nil
}

// DeleteProduct deletes specified product from database
func DeleteProduct(db *sql.DB, productName string) error {

	// Check if product exists in database
	productID, err := getProductID(db, productName)
	if err != nil {
		return err
	}

	query := "DELETE FROM products WHERE id = ?"
	result, err := db.Exec(query, productID)
	if err != nil {
		logger.Println("Error deleting product from database:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Println("Error getting rows affected:", err)
		return err
	}

	if rowsAffected == 0 {
		logger.Printf("Product %s not found in database", productName)
		return errors.New("no rows affected, product not found")
	}

	// Delete associated records from 'product_categories' table
	query = "DELETE FROM product_categories WHERE product_id = ?"
	_, err = db.Exec(query, productID)
	if err != nil {
		logger.Println("Error deleting product associations from database:", err)
		return err
	}

	logger.Println("Product deleted successfully")
	return nil
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
