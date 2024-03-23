package api

import (
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// registerHandlers registers HTTP request handlers
func registerHandlers(router *mux.Router) {
	//CRUD category
	categoriesRouter := router.PathPrefix("/categories").Subrouter()
	categoriesRouter.HandleFunc("/new", CreateCategoryHandler).Methods("POST")      // CREATE
	categoriesRouter.HandleFunc("/list", GetCategoriesHandler).Methods("GET")       // READ
	categoriesRouter.HandleFunc("/{name}", UpdateCategoryHandler).Methods("PUT")    // UPDATE
	categoriesRouter.HandleFunc("/{name}", DeleteCategoryHandler).Methods("DELETE") // DELETE

	//CRUD product
	productsRouter := router.PathPrefix("/products").Subrouter()
	productsRouter.HandleFunc("/new", CreateProductHandler).Methods("POST")           // CREATE
	productsRouter.HandleFunc("/{name}", GetProductsByCategoryHandler).Methods("GET") // READ
	productsRouter.HandleFunc("", UpdateProductHandler).Methods("PUT")                // UPDATE
	productsRouter.HandleFunc("/{name}", DeleteProductHandler).Methods("DELETE")      // DELETE

	// Auth router
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", RegisterUserHandler).Methods("POST")
	authRouter.HandleFunc("/login", LoginUserHandler).Methods("POST")
}

func StartServer() {
	router := mux.NewRouter()
	// Middleware for processing request ID
	router.Use(RequestIDMiddleware)
	router.Use(RequireValidTokenMiddleware)

	registerHandlers(router)

	// HTTP server start
	logger.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
