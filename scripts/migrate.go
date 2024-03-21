package scripts

import (
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

type Migration struct {
	Version     int
	Description string
	Up          func(db *sql.DB) error
	Down        func(db *sql.DB) error
}

var migrations = []Migration{
	{
		Version:     1,
		Description: "Create categories table",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
                CREATE TABLE IF NOT EXISTS categories (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    name TEXT UNIQUE 
                )
            `)
			return err
		},
		Down: func(db *sql.DB) error {
			_, err := db.Exec("DROP TABLE IF EXISTS categories")
			return err
		},
	},
	{
		Version:     2,
		Description: "Create products table",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
        		CREATE TABLE IF NOT EXISTS products (
            		id INTEGER PRIMARY KEY AUTOINCREMENT,
            		name TEXT UNIQUE 
        		)
    		`)
			return err
		},
		Down: func(db *sql.DB) error {
			_, err := db.Exec("DROP TABLE IF EXISTS products")
			return err
		},
	},
	{
		Version:     3,
		Description: "Create product_categories table",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
			   CREATE TABLE IF NOT EXISTS product_categories (
				   product_id INTEGER,
				   category_id INTEGER,
				   FOREIGN KEY(product_id) REFERENCES products(id),
				   FOREIGN KEY(category_id) REFERENCES categories(id),
				   PRIMARY KEY (product_id, category_id)
			   )
			`)
			return err
		},
		Down: func(db *sql.DB) error {
			_, err := db.Exec("DROP TABLE IF EXISTS product_categories")
			return err
		},
	},
	{
		Version:     4,
		Description: "Create users table",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					username TEXT,
					password TEXT
				)
			`)
			return err
		},
		Down: func(db *sql.DB) error {
			_, err := db.Exec("DROP TABLE IF EXISTS users")
			return err
		},
	},
}

// Migrate migrates database schema by creating necessary tables if they do not exist
// takes *sql.DB parameter representing database connection
// If any error occurs during table creation, logs error and returns it
// Otherwise logs successful creation of each table and returns nil
func Migrate(db *sql.DB) error {
	// Create migration table if not exists
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS migrations (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            version INTEGER,
            description TEXT
        )
    `)
	if err != nil {
		logger.Println("Error creating 'migrations' table:", err)
		return err
	}
	logger.Println("Table 'migrations' created successfully")

	// Get current version
	var currentVersion int
	err = db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM migrations").Scan(&currentVersion)
	if err != nil {
		logger.Println("Error getting current migration version:", err)
		return err
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			logger.Println("Applying migration:", migration.Description)
			if err = migration.Up(db); err != nil {
				logger.Println("Error applying migration:", err)
				return err
			}

			// Record migration
			_, err = db.Exec("INSERT INTO migrations (version, description) VALUES (?, ?)",
				migration.Version, migration.Description)
			if err != nil {
				logger.Println("Error recording migration:", err)
				return err
			}

			logger.Println("Migration applied successfully")
		}
	}

	return nil
}

func Rollback(db *sql.DB) error {
	var version int
	err := db.QueryRow("SELECT MAX(version) FROM migrations").Scan(&version)
	if err != nil {
		logger.Println("Error getting current migration version:", err)
		return err
	}

	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]
		if migration.Version <= version {
			logger.Println("Rolling back migration:", migration.Description)
			if err := migration.Down(db); err != nil {
				logger.Println("Error rolling back migration:", err)
				return err
			}

			// Remove migration record
			_, err := db.Exec("DELETE FROM migrations WHERE version = ?", migration.Version)
			if err != nil {
				logger.Println("Error removing migration record:", err)
				return err
			}

			logger.Println("Migration rolled back successfully")

			return nil
		}
	}

	logger.Println("No migrations to rollback")

	return nil
}
