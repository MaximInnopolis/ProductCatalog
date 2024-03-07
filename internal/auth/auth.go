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

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get user data
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	err = models.RegisterUser(database.GetDB(), &user)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}
	// Write success message to response
	utils.WriteJSONResponse(w, http.StatusOK, "Successfully registered")
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get user data
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	token, err := models.LoginUser(database.GetDB(), &user)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Write success message to response
	utils.WriteTokenJSONResponse(w, http.StatusOK, "Successfully logged in", token)
}

func RequireValidToken(w http.ResponseWriter, r *http.Request) bool {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		logger.Println("Authorization header is missing")

		utils.WriteErrorJSONResponse(w, errors.New("authorization header is missing"), http.StatusInternalServerError)
		return false
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	authenticated, err := models.IsTokenValid(database.GetDB(), tokenString)
	if err != nil {
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return false
	}

	if !authenticated {
		utils.WriteErrorJSONResponse(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return false
	}

	return true
}
