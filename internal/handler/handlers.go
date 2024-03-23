package handler

import (
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}

// InitRoutes inits HTTP request handlers
func (h *Handler) InitRoutes(router *mux.Router) {
	//CRUD category
	//categoriesRouter := router.PathPrefix("/categories").Subrouter()
	//categoriesRouter.HandleFunc("/new", CreateCategoryHandler).Methods("POST")      // CREATE
	//categoriesRouter.HandleFunc("/list", GetCategoriesHandler).Methods("GET")       // READ
	//categoriesRouter.HandleFunc("/{name}", UpdateCategoryHandler).Methods("PUT")    // UPDATE
	//categoriesRouter.HandleFunc("/{name}", DeleteCategoryHandler).Methods("DELETE") // DELETE

	//CRUD product
	//productsRouter := router.PathPrefix("/products").Subrouter()
	//productsRouter.HandleFunc("/new", CreateProductHandler).Methods("POST")           // CREATE
	//productsRouter.HandleFunc("/{name}", GetProductsByCategoryHandler).Methods("GET") // READ
	//productsRouter.HandleFunc("", UpdateProductHandler).Methods("PUT")                // UPDATE
	//productsRouter.HandleFunc("/{name}", DeleteProductHandler).Methods("DELETE")      // DELETE

	// Auth router
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", h.RegisterUserHandler).Methods("POST")
	authRouter.HandleFunc("/login", h.LoginUserHandler).Methods("POST")
}

func (h *Handler) StartServer() {
	router := mux.NewRouter()
	// Middleware for processing request ID
	router.Use(RequestIDMiddleware)
	router.Use(RequireValidTokenMiddleware)

	h.InitRoutes(router)

	// HTTP server start
	logger.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
