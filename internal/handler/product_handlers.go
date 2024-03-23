package handler

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
	"github.com/MaximInnopolis/ProductCatalog/internal/utils"
	"net/http"
)

// requestData stores data received in request body in CreateProductHandler and UpdateProductHandler
var requestData struct {
	Name       string           `json:"Name"`
	Categories []model.Category `json:"Categories"`
}

// CreateProductHandler handles requests to create new product
// Parses request body to extract product data and creates new product in database
// If successful, writes success message to response; otherwise writes error response
func (h *Handler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Printf(ctx, "Creating new product")

	// Parse request body to get product data
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	// Create product struct and insert product name from request data
	product := model.Product{Name: requestData.Name}
	// Add product to database
	err = h.service.Product.CreateProduct(ctx, &product, requestData.Categories)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}
	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusCreated, "Product created")
}

// UpdateProductHandler handles requests to update existing product
// Parses request body to extract product data and updates corresponding product in database
// If successful, writes success message to response; otherwise writes error response
func (h *Handler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Printf(ctx, "Updating product")

	// Parse request body to get product data
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	// Create product struct and insert product name from request data
	product := model.Product{Name: requestData.Name}
	// Update product in database
	err = h.service.Product.UpdateProduct(ctx, &product, requestData.Categories)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}
	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusOK, "Product updated")
}

// DeleteProductHandler handles requests to delete existing product
// First extracts product name from request
// Then deletes product from database
// If successful, writes success message to response; otherwise writes error response
func (h *Handler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Printf(ctx, "Deleting product")

	// Extract product name from request
	productName := GetNameFromRequest(r)

	// Delete product from database
	err := h.service.Product.DeleteProduct(ctx, productName)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusOK, "Product deleted")
}

// GetProductsByCategoryHandler handles requests to retrieve list of products by category
// First extracts category name from request
// Then retrieves list of products associated with specified category name from database
// If successful, writes list of products in JSON format to response; otherwise writes error response
func (h *Handler) GetProductsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Printf(ctx, "Getting products by category")

	// Extract category name from request
	categoryName := GetNameFromRequest(r)

	// Get list of products by category name from database
	products, err := h.service.Product.GetProductsByCategory(ctx, categoryName)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Output in JSON format list of categories
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
