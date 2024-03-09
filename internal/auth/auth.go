package auth

import (
	"encoding/json"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"github.com/MaximInnopolis/ProductCatalog/internal/utils"
	"net/http"
	"strings"
)

// RegisterUserHandler handles user registration requests
// Parses request body to extract user data, attempts to register user in database
// If successful, writes success message to response; otherwise writes error response
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get user data
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	// Add user to database
	err = models.RegisterUser(database.GetDB(), &user)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}
	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusOK, "Successfully registered")
}

// LoginUserHandler handles user login requests
// Parses request body to extract user data, attempts to log in the user, and generate token
// If successful, writes success message along with the token to the response; otherwise writes error response
func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get user data
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	// Login user and generate token
	token, err := models.LoginUser(database.GetDB(), &user)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Write success message to response
	utils.WriteTokenJSONResponse(w, http.StatusOK, "Successfully logged in", token)
}

// RequireValidToken checks if request contains valid authentication token
// Extracts token from Authorization header and verifies its validity against the database
// If token is missing or invalid, writes appropriate error response and returns false; otherwise returns true
func RequireValidToken(w http.ResponseWriter, r *http.Request) bool {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	// Check if the Authorization header is missing
	if authHeader == "" {
		logger.Println("Authorization header is missing")
		utils.WriteErrorJSONResponse(w, errors.New("authorization header is missing"), http.StatusInternalServerError)
		return false
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Check if token is valid
	authenticated, err := models.IsTokenValid(database.GetDB(), tokenString)
	if err != nil {
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return false
	}

	// If token is not valid, return unauthorized status
	if !authenticated {
		utils.WriteErrorJSONResponse(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return false
	}

	return true
}
