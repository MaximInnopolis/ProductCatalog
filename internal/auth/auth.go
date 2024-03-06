package auth

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"net/http"
	"strings"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	// Parse request body to get user data
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := models.RegisterUser(database.GetDB(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(token)
}

func RequireValidToken(w http.ResponseWriter, r *http.Request) bool {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		logger.Println("Authorization header is missing")
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return false
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	authenticated, err := models.IsTokenValid(database.GetDB(), tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	if !authenticated {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}
