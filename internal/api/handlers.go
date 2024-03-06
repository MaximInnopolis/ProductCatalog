package api

import (
	"github.com/MaximInnopolis/ProductCatalog/internal/auth"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHandlers() {
	router := mux.NewRouter()
	router.HandleFunc("/products/{categoryName}", GetProductsByCategoryHandler).Methods("GET")

	// Auth routes
	router.HandleFunc("/register", auth.RegisterUserHandler).Methods("POST")

	//CRUD category
	router.HandleFunc("/categories", CreateCategoryHandler).Methods("POST")                  // CREATE
	router.HandleFunc("/categories", GetCategoriesHandler).Methods("GET")                    // READ
	router.HandleFunc("/categories/{categoryName}", UpdateCategoryHandler).Methods("PUT")    // UPDATE
	router.HandleFunc("/categories/{categoryName}", DeleteCategoryHandler).Methods("DELETE") // DELETE

	http.Handle("/", router)
}
