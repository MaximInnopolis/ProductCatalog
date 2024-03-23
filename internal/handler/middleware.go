package handler

import (
	"context"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/utils"
	"net/http"
	"strings"
)

// RequestIDMiddleware assigns endpoint to each request and adds it to the request context
func (h *Handler) RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract endpoint information from request
		endpoint := r.URL.Path
		ctx := context.WithValue(r.Context(), "endpoint", endpoint)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireValidTokenMiddleware checks if request contains valid authentication token
// Extracts token from Authorization header and verifies its validity
// If token is missing or invalid, writes appropriate error response; otherwise pass further
func (h *Handler) RequireValidTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract endpoint information from request
		endpoint := r.URL.Path

		// Define list of endpoints that require authorization
		protectedEndpoints := []string{
			"/categories/new",    // Creating category
			"/categories/{name}", // Updating and deleting category
			"/products/new",      // Creating product
			"/products/{name}",   // Updating and deleting product
		}

		// Check if authorization is required for current endpoint
		requiresAuth := false
		for _, protectedEndpoint := range protectedEndpoints {
			if endpoint == protectedEndpoint {
				requiresAuth = true
				break
			}
		}

		// If authorization is required for the current endpoint, check
		if requiresAuth {

			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			// Check if the Authorization header is missing
			if authHeader == "" {
				logger.Printf(ctx, "Authorization header is missing")
				utils.WriteErrorJSONResponse(w, errors.New("authorization header is missing"), http.StatusInternalServerError)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Check if token is valid
			authenticated, err := h.service.IsTokenValid(ctx, tokenString)
			if err != nil {
				utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
				return
			}

			// If token is not valid, return unauthorized status
			if !authenticated {
				utils.WriteErrorJSONResponse(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}
		}

		// If token is valid pass request further
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
