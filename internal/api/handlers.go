package api

import (
	"github.com/MaximInnopolis/ProductCatalog/internal/auth"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHandlers() {
	router := mux.NewRouter()
	router.HandleFunc("/products/{name}", GetProductsByCategoryHandler).Methods("GET")

	// Auth routes
	router.HandleFunc("/register", auth.RegisterUserHandler).Methods("POST")

	//CRUD category
	router.HandleFunc("/category", CreateCategoryHandler).Methods("POST")          // CREATE
	router.HandleFunc("/categories", GetCategoriesHandler).Methods("GET")          // READ
	router.HandleFunc("/category/{name}", UpdateCategoryHandler).Methods("PUT")    // UPDATE
	router.HandleFunc("/category/{name}", DeleteCategoryHandler).Methods("DELETE") // DELETE

	//CRUD product
	router.HandleFunc("/product", CreateProductHandler).Methods("POST")          // CREATE
	router.HandleFunc("/product", UpdateProductHandler).Methods("PUT")           // UPDATE
	router.HandleFunc("/product/{name}", DeleteProductHandler).Methods("DELETE") // DELETE

	http.Handle("/", router)
}
