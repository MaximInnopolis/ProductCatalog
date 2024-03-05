package main

import (
	"log"
	//"net/http"

	//"github.com/MaximInnopolis/ProductCatalog/internal/api"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
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

	//http.HandleFunc("/categories", api.ListCategoriesHandler(database.GetDB()))
	//http.HandleFunc("/categories/", api.ListProductsInCategoryHandler(database.GetDB()))
	//
	//// Запуск HTTP сервера
	//logger.Println("Server started on port 8080")
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
