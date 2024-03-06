package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open database connection
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createRecords(db); err != nil {
		log.Fatal(err)
	}

	log.Println("Records created successfully")
}

func createRecords(db *sql.DB) error {
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
	result, err := tx.Exec("INSERT INTO products (name) VALUES (?)", "BRead")
	if err != nil {
		return err
	}
	productID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Insert categories
	categories := []string{"Flour"}
	for _, categoryName := range categories {
		result, err = tx.Exec("INSERT INTO categories (name) VALUES (?)", categoryName)
		if err != nil {
			return err
		}
		categoryID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Insert into product_categories
		_, err = tx.Exec("INSERT INTO product_categories (product_id, category_id) VALUES (?, ?)", productID, categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}
