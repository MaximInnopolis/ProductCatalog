package main

import (
	"github.com/MaximInnopolis/ProductCatalog/internal/handler"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/repository"
	"github.com/MaximInnopolis/ProductCatalog/internal/service"
	"github.com/MaximInnopolis/ProductCatalog/internal/utils"
	"github.com/MaximInnopolis/ProductCatalog/scripts"
	"log"
	"os"
)

func main() {
	// Logger initialization
	defer func() {
		if err := logger.Close(); err != nil {
			log.Fatalf("Failed to close log file: %v", err)
		}
	}()

	// Load env data
	utils.LoadEnv()

	// Initialize database connection
	db, err := repository.NewDB(os.Getenv("DB_PATH"))
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repository.Close(db)

	// Migrate
	if err := scripts.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	logger.Println("Migration completed successfully")

	//// Rollback
	//if err := scripts.Rollback(db); err != nil {
	//	log.Fatalf("Failed to rollback: %v", err)
	//}
	//
	//log.Println("Rollback completed successfully")

	//// Database records creation (Not necessary)
	//if err := scripts.CreateRecords(db); err != nil {
	//	log.Fatalf("Failed to create records: %v", err)
	//}
	//
	//logger.Println("Records created successfully")

	// Start data collection
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	scripts.StartDataCollection(services)

	// Register HTTP request handlers
	handlers.StartServer()
}
