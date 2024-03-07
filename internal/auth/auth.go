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

	err = models.RegisterUser(database.GetDB(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write the success message to the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Successfully registered"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	// Parse request body to get user data
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := models.LoginUser(database.GetDB(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Define token format
	type TokenResponse struct {
		Token string `json:"token"`
	}

	// Create response struct with the token
	response := TokenResponse{
		Token: token,
	}

	_, err = w.Write([]byte("Successfully logged in\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode response to JSON and send it in response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		http.Error(w, "You're Unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}
