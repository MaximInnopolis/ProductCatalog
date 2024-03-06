package api

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/auth"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"github.com/gorilla/mux"
	"net/http"
)

// CreateCategoryHandler create new category
func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Check if token is valid
	if !auth.RequireValidToken(w, r) {
		return
	}

	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = models.AddCategory(database.GetDB(), &category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

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

// UpdateCategoryHandler edit specified existing category
func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extract category name from request
	categoryName := GetNameFromRequest(r)

	// Check if token is valid
	if !auth.RequireValidToken(w, r) {
		return
	}

	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = models.UpdateCategory(database.GetDB(), categoryName, &category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteCategoryHandler deletes specified category
func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extract category name from request
	categoryName := GetNameFromRequest(r)

	// Check if token is valid
	if !auth.RequireValidToken(w, r) {
		return
	}

	// Delete category from database
	err := models.DeleteCategory(database.GetDB(), categoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetNameFromRequest retrieves category name from URL
func GetNameFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	categoryName := vars["name"]
	return categoryName
}
