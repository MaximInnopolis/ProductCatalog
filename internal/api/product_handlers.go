package api

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/auth"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product := models.Product{Name: requestData.Name}
	err = models.AddProduct(database.GetDB(), &product, requestData.Categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

//// UpdateProductHandler updates an existing product
//func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
//	// Check if token is valid
//	if !auth.RequireValidToken(w, r) {
//		return
//	}
//
//	err := json.NewDecoder(r.Body).Decode(&requestData)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	product := models.Product{Name: requestData.Name}
//	err = models.UpdateProduct(database.GetDB(), &product, requestData.Categories)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}
//
//// DeleteProductHandler deletes an existing product
//func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
//	// Check if token is valid
//	if !auth.RequireValidToken(w, r) {
//		return
//	}
//
//	err := json.NewDecoder(r.Body).Decode(&requestData)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	err = models.DeleteProduct(database.GetDB(), requestData.Name)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}

// GetProductsByCategoryHandler returns product list of concrete category
func GetProductsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extract category name from request
	categoryName := GetNameFromRequest(r)

	products, err := models.GetProductsByCategory(database.GetDB(), categoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
