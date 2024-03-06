package api

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"github.com/gorilla/mux"
	"net/http"
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

// GetProductsByCategoryHandler returns product list of concrete category
func GetProductsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryName := getCategoryNameFromRequest(r)
	products, err := models.GetProductsByCategory(database.GetDB(), categoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// getCategoryNameFromRequest retrieves category name from URL
func getCategoryNameFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	categoryName := vars["categoryName"]
	return categoryName
}
