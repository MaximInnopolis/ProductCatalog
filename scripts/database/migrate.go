package main

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open SQLite database connection
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Migrate
	if err := migrate(db); err != nil {
		log.Fatal(err)
	}

	log.Println("Migration completed successfully")
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT
        )
    `)
	if err != nil {
		return err
	}

	logger.Println("Table 'categories' created successfully")

	// TODO: change relation muliple-multiple
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            category_id INTEGER, // TODO: remove
            FOREIGN KEY(category_id) REFERENCES categories(id) // TODO: remove
        )
    `)
	if err != nil {
		return err
	}

	logger.Println("Table 'products' created successfully")

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	// TODO: look closely
	//_, err = db.Exec(`
	//    CREATE TABLE IF NOT EXISTS product_categories (
	//        product_id INTEGER,
	//        category_id INTEGER,
	//        FOREIGN KEY(product_id) REFERENCES products(id),
	//        FOREIGN KEY(category_id) REFERENCES categories(id),
	//        PRIMARY KEY (product_id, category_id)
	//    )
	//`)
	//if err != nil {
	//	return err
	//}
	//
	//logger.Println("Table 'product_categories' created successfully")

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT,
            password TEXT
        )
    `)
	if err != nil {
		return err
	}

	logger.Println("Table 'users' created successfully")

	return nil
}
