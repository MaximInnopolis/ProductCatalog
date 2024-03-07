package api

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/auth"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"github.com/MaximInnopolis/ProductCatalog/internal/utils"
	"github.com/gorilla/mux"
	"net/http"
)

// CreateCategoryHandler creates new category
func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Check if token is valid
	if !auth.RequireValidToken(w, r) {
		return
	}

	logger.Println("Handling CreateCategoryHandler") // Добавленный отладочный вывод

	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	_, err = models.AddCategory(database.GetDB(), &category)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusCreated, "Category created")
}

// GetCategoriesHandler returns list of all categories
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetAllCategories(database.GetDB())
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// UpdateCategoryHandler edits specified existing category
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
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	err = models.UpdateCategory(database.GetDB(), categoryName, &category)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusOK, "Category updated")
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
