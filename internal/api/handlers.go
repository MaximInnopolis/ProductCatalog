package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHandlers() {
	router := mux.NewRouter()
	router.HandleFunc("/categories", GetCategoriesHandler).Methods("GET")
	router.HandleFunc("/products/{categoryName}", GetProductsByCategoryHandler).Methods("GET")
	http.Handle("/", router)
}
