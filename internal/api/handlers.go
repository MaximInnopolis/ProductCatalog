package api

import (
	"net/http"
)

func RegisterHandlers() {
	http.HandleFunc("/categories", GetCategoriesHandler)
	http.HandleFunc("/categories/{categoryName}", GetProductsInCategoryHandler)
}
