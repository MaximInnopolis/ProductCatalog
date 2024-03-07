package utils

import (
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"net/http"
)

func WriteErrorJSONResponse(w http.ResponseWriter, err error, status int) {
	response := models.ErrorResponse{Error: err.Error()}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
	}
}

func WriteJSONResponse(w http.ResponseWriter, status int, message string) {
	response := models.ResponseMessage{Message: message}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
	}
}

func WriteTokenJSONResponse(w http.ResponseWriter, status int, message, token string) {
	response := models.TokenResponseMessage{Message: message, Token: token}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		WriteErrorJSONResponse(w, err, http.StatusInternalServerError)
	}
}
