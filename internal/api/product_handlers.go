package api

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/auth"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"github.com/MaximInnopolis/ProductCatalog/internal/utils"
	"net/http"
)

var requestData struct {
	Name       string            `json:"Name"`
	Categories []models.Category `json:"Categories"`
}

// CreateProductHandler creates new product
func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	// Check if token is valid
	if !auth.RequireValidToken(w, r) {
		return
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	product := models.Product{Name: requestData.Name}
	err = models.AddProduct(database.GetDB(), &product, requestData.Categories)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}
	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusCreated, "Product created")
}

// UpdateProductHandler updates existing product
func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {

	// Check if token is valid
	if !auth.RequireValidToken(w, r) {
		return
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	product := models.Product{Name: requestData.Name}
	err = models.UpdateProduct(database.GetDB(), &product, requestData.Categories)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}
	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusOK, "Product updated")
}

// DeleteProductHandler deletes an existing product
func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {

	// Extract product name from request
	productName := GetNameFromRequest(r)

	// Check if token is valid
	if !auth.RequireValidToken(w, r) {
		return
	}

	err := models.DeleteProduct(database.GetDB(), productName)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusOK, "Product deleted")
}

// GetProductsByCategoryHandler returns product list of concrete category
func GetProductsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extract category name from request
	categoryName := GetNameFromRequest(r)

	products, err := models.GetProductsByCategory(database.GetDB(), categoryName)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
