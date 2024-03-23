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
	categoriesRouter := router.PathPrefix("/categories").Subrouter()
	categoriesRouter.HandleFunc("/list", h.GetCategoriesHandler).Methods("GET") // READ

	createCategoryHandler := http.HandlerFunc(h.CreateCategoryHandler)
	categoriesRouter.Handle("/new", h.RequireValidTokenMiddleware(createCategoryHandler)).Methods("POST")

	updateCategoryHandler := http.HandlerFunc(h.UpdateCategoryHandler)
	categoriesRouter.Handle("/{name}", h.RequireValidTokenMiddleware(updateCategoryHandler)).Methods("PUT")

	deleteCategoryHandler := http.HandlerFunc(h.DeleteCategoryHandler)
	categoriesRouter.Handle("/{name}", h.RequireValidTokenMiddleware(deleteCategoryHandler)).Methods("DELETE")

	//CRUD product
	productsRouter := router.PathPrefix("/products").Subrouter()
	productsRouter.HandleFunc("/{name}", h.GetProductsByCategoryHandler).Methods("GET") // READ

	createProductHandler := http.HandlerFunc(h.CreateProductHandler)
	productsRouter.Handle("/new", h.RequireValidTokenMiddleware(createProductHandler)).Methods("POST")

	updateProductHandler := http.HandlerFunc(h.UpdateProductHandler)
	productsRouter.Handle("/", h.RequireValidTokenMiddleware(updateProductHandler)).Methods("PUT")

	deleteProductHandler := http.HandlerFunc(h.DeleteProductHandler)
	productsRouter.Handle("/{name}", h.RequireValidTokenMiddleware(deleteProductHandler)).Methods("DELETE")

	// Auth router
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", h.RegisterUserHandler).Methods("POST")
	authRouter.HandleFunc("/login", h.LoginUserHandler).Methods("POST")
}

func (h *Handler) StartServer() {
	router := mux.NewRouter()
	// Middleware for processing request ID
	router.Use(h.RequestIDMiddleware)
	h.InitRoutes(router)

	// HTTP server start
	logger.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
