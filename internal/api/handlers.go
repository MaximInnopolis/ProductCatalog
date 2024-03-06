package api

import (
	"github.com/MaximInnopolis/ProductCatalog/internal/auth"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHandlers() {
	router := mux.NewRouter()
	router.HandleFunc("/categories", GetCategoriesHandler).Methods("GET")
	router.HandleFunc("/products/{categoryName}", GetProductsByCategoryHandler).Methods("GET")

	// Auth routes
	router.HandleFunc("/register", auth.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/login", auth.LoginUserHandler).Methods("POST")

	http.Handle("/", router)
}
