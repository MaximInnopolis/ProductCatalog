package api

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/MaximInnopolis/ProductCatalog/internal/database"
)

// GetCategoriesHandler returns list of all categories
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetAllCategories(database.GetDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GetProductsInCategoryHandler returns product list of concrete category
func GetProductsInCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := getCategoryNameFromRequest(r)
	products, err := database.GetProductsByCategory(database.GetDB(), categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// getCategoryNameFromRequest retrieves category name from URL
func getCategoryNameFromRequest(r *http.Request) string {
	// Use Gorilla Mux to get URL parameters
	vars := mux.Vars(r)
	categoryName := vars["categoryName"]
	return categoryName
}
