package scripts

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"net/http"
	"time"
)

// StartDataCollection starts collecting data from source and saving in database
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

// collectAndSaveData collects data from source and saves it to database
func CollectAndSaveProducts(db *sql.DB) error {
	logger.Println("Started collecting....")
	resp, err := http.Get("https://emojihub.yurace.pro/api/all")
	if err != nil {
		logger.Println(err)
		return err
	}
	defer resp.Body.Close()

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
		err = models.AddProduct(db, &product, category)
		if err != nil {
			continue
		}
	}

	fmt.Println("Products collected and saved successfully")
	return nil
}
