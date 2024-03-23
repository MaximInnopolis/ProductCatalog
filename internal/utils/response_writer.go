package utils

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
	"net/http"
)

// WriteErrorJSONResponse writes error JSON response with given error message and status code to provided http.ResponseWriter
func WriteErrorJSONResponse(w http.ResponseWriter, err error, status int) {
	response := model.ErrorResponse{Error: err.Error()}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
	}
}

// WriteJSONResponse writes JSON response with given message and status code to provided http.ResponseWriter
func WriteJSONResponse(w http.ResponseWriter, status int, message string) {
	response := model.ResponseMessage{Message: message}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
	}
}

// WriteTokenJSONResponse writes JSON response with given message, status code, and JWT token to provided http.ResponseWriter

func WriteTokenJSONResponse(w http.ResponseWriter, status int, message, token string) {
	response := model.TokenResponseMessage{Message: message, Token: token}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
	}
}
