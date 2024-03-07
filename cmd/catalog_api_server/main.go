package main

import (
	"github.com/MaximInnopolis/ProductCatalog/internal/api"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/scripts"
	"log"
)

func main() {
	// Logger initialization
	defer func() {
		if err := logger.Close(); err != nil {
			log.Fatalf("Failed to close log file: %v", err)
		}
	}()

	// Database initialization
	dbPath := "./data/database.db"
	if err := database.Init(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Migrate
	if err := scripts.Migrate(database.GetDB()); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	log.Println("Migration completed successfully")

	// Database records creation (Not necessary)
	if err := scripts.CreateRecords(database.GetDB()); err != nil {
		log.Fatalf("Failed to create records: %v", err)
	}

	log.Println("Records created successfully")

	api.RegisterHandlers()
}
