package auth

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
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

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	// Parse request body to get user data
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	user.Token = tokenString

	authenticated, err := models.AuthenticateUser(database.GetDB(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !authenticated {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
