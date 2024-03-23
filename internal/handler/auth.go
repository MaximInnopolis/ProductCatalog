package handler

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/utils"

	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"

	"net/http"
)

// RegisterUserHandler handles user registration requests
// Parses request body to extract user data, attempts to register user in database
// If successful, writes success message to response; otherwise writes error response
func (h *Handler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Printf(ctx, "Registering user")

	// Parse request body to get user data
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	// Add user to database
	//err = models.RegisterUser(ctx, &user)
	err = h.service.Authorization.CreateUser(ctx, &user)
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
func (h *Handler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request body to get user data
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Write error response with bad request status code
		utils.WriteErrorJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	// Login user and generate token
	//token, err := models.LoginUser(ctx, &user)
	token, err := h.service.Authorization.GenerateToken(ctx, &user)
	if err != nil {
		// Write error response with internal server error status code
		utils.WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Write success message to response
	utils.WriteTokenJSONResponse(w, http.StatusOK, "Successfully logged in", token)
}
