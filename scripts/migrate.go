package scripts

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"

	_ "github.com/mattn/go-sqlite3"
)

// Migrate migrates database schema by creating necessary tables if they do not exist
// takes *sql.DB parameter representing database connection
// If any error occurs during table creation, logs error and returns it
// Otherwise logs successful creation of each table and returns nil
func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE 
        )
    `)
	if err != nil {
		logger.Println("Error creating 'categories' table:", err)
		return err
	}

	logger.Println("Table 'categories' created successfully")

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE 
        )
    `)
	if err != nil {
		logger.Println("Error creating 'products' table:", err)
		return err
	}

	logger.Println("Table 'products' created successfully")

	_, err = db.Exec(`
	   CREATE TABLE IF NOT EXISTS product_categories (
	       product_id INTEGER,
	       category_id INTEGER,
	       FOREIGN KEY(product_id) REFERENCES products(id),
	       FOREIGN KEY(category_id) REFERENCES categories(id),
	       PRIMARY KEY (product_id, category_id)
	   )
	`)
	if err != nil {
		logger.Println("Error creating 'product_categories' table:", err)
		return err
	}

	logger.Println("Table 'product_categories' created successfully")

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT,
            password TEXT
        )
    `)
	if err != nil {
		logger.Println("Error creating 'users' table:", err)
		return err
	}

	logger.Println("Table 'users' created successfully")

	return nil
}
