package scripts

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

func CreateRecords(db *sql.DB) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		logger.Println("Error beginning transaction:", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Insert product
	result, err := tx.Exec("INSERT INTO products (name) VALUES (?) ON CONFLICT(name) DO NOTHING", "Bread")
	if err != nil {
		logger.Println("Error inserting product:", err)
		return err
	}
	productID, err := result.LastInsertId()
	if err != nil {
		logger.Println("Error getting last inserted product ID:", err)
		return err
	}

	// Insert categories
	categories := []string{"Flour", "Food"}
	for _, categoryName := range categories {
		result, err = tx.Exec("INSERT INTO categories (name) VALUES (?) ON CONFLICT(name) DO NOTHING", categoryName)
		if err != nil {
			logger.Println("Error inserting category:", err)
			return err
		}
		categoryID, err := result.LastInsertId()
		if err != nil {
			logger.Println("Error getting last inserted category ID:", err)
			return err
		}

		// Insert into product_categories
		if productID == 0 || categoryID == 0 {
			return nil
		}

		_, err = tx.Exec("INSERT INTO product_categories (product_id, category_id) VALUES (?, ?)", productID, categoryID)
		if err != nil {
			logger.Println("Error inserting into product_categories:", err)
			return err
		}
	}

	return nil
}
