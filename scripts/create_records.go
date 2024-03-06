package scripts

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func CreateRecords(db *sql.DB) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
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
		return err
	}
	productID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Insert categories
	categories := []string{"Flour", "Food"}
	for _, categoryName := range categories {
		result, err = tx.Exec("INSERT INTO categories (name) VALUES (?) ON CONFLICT(name) DO NOTHING", categoryName)
		if err != nil {
			return err
		}
		categoryID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Insert into product_categories
		if productID == 0 || categoryID == 0 {
			return nil
		}

		_, err = tx.Exec("INSERT INTO product_categories (product_id, category_id) VALUES (?, ?)", productID, categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}
