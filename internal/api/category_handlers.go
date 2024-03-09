package api

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/auth"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"github.com/MaximInnopolis/ProductCatalog/internal/utils"
	"github.com/gorilla/mux"
	"net/http"
)

// CreateCategoryHandler handles requests to create new category
// First checks if token is valid
// Then parses request body to extract category data
// If successful, adds category to database and writes success message to response
// If any errors occur, writes error response
func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Check if token is valid
	if !auth.RequireValidToken(w, r) {
		return
	}

	// Parse request body to get category data
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	// Add category to database
	_, err = models.AddCategory(database.GetDB(), &category)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusCreated, "Category created")
}

// GetCategoriesHandler handles requests to retrieve list of all categories
// Fetches list of categories from database
// If successful, outputs list of categories in JSON format
// If any errors occur, writes error response
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	// Get list of categories from database
	categories, err := models.GetAllCategories(database.GetDB())
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Output in JSON format list of categories
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// UpdateCategoryHandler handles requests to update existing category
// Extracts category name from request and checks if token is valid
// If successful, parses request body to get category data and updates category in database
// If any errors occur during process, writes error response
func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extract category name from request
	categoryName := GetNameFromRequest(r)

	// Check if token is valid
	if !auth.RequireValidToken(w, r) {
		return
	}

	// Parse request body to get category data
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	// Update category in database
	err = models.UpdateCategory(database.GetDB(), categoryName, &category)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusOK, "Category updated")
}

// DeleteCategoryHandler handles requests to delete specified category
// Extracts category name from request, checks if token is valid,
// and then deletes category from database.
// If any errors occur during process, writes error response
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
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusOK, "Category deleted")
}

// GetNameFromRequest retrieves category name from URL
func GetNameFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	categoryName := vars["name"]
	return categoryName
}
