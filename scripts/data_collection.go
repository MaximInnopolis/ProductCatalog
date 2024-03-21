package scripts

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"net/http"
	"time"
)

// StartDataCollection starts periodic data collection process to collect data from specified source and save to database
// sets up ticker to trigger data collection at regular intervals and runs collection process in goroutine
func StartDataCollection() {

	// Set up ticker for collecting data at intervals
	ticker := time.NewTicker(time.Hour)

	go func() {
		for {
			select {
			case <-ticker.C:
				// Collect and save data to database
				CollectAndSaveProducts(database.GetDB())
			}
		}
	}()
}

// CollectAndSaveProducts collects data from specified source, processes it, and saves to database
// sends HTTP request to source API to retrieve raw product data, processes data, and saves to database
// If any errors occur during process, logs error and continues processing other products
func CollectAndSaveProducts(db *sql.DB) error {
	ctx := context.WithValue(context.Background(), "endpoint", "product_collection")
	logger.Printf(ctx, "Started collecting....")
	// Send HTTP request to retrieve raw product data from source
	resp, err := http.Get("https://emojihub.yurace.pro/api/all")
	if err != nil {
		logger.Println(err)
		return err
	}
	defer resp.Body.Close()

	// Decode raw product data from response body
	var rawProducts []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&rawProducts)
	if err != nil {
		logger.Println(err)
		return err
	}

	// Save products to database
	for _, rawProduct := range rawProducts {
		productName, ok := rawProduct["name"].(string)
		if !ok {
			return errors.New("product name not found in source")
		}

		categoryName, ok := rawProduct["category"].(string)
		if !ok {
			return errors.New("category name not found in source")
		}

		product := models.Product{Name: productName}
		var category []models.Category

		category = append(category, models.Category{Name: categoryName})
		err = models.AddProduct(ctx, db, &product, category)
		if err != nil {
			continue
		}
	}

	logger.Printf(ctx, "Products collected and saved successfully")
	return nil
}
