package auth_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MaximInnopolis/ProductCatalog/internal/auth"
	"github.com/MaximInnopolis/ProductCatalog/internal/database"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	"github.com/stretchr/testify/assert"
)

// tokenValidator interface represents function for token validity check
type tokenValidator interface {
	IsTokenValid(db *sql.DB, tokenString string) (bool, error)
}

// mockTokenValidator implements tokenValidator interface for mocking IsTokenValid function
type mockTokenValidator struct{}

// IsTokenValid always returns true without error
func (m *mockTokenValidator) IsTokenValid(db *sql.DB, tokenString string) (bool, error) {
	return true, nil
}

func TestRegisterUserHandler_Success(t *testing.T) {
	// Create test user
	user := models.User{Username: "testuser", Password: "testpassword"}
	userJSON, _ := json.Marshal(user)

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP response recorder
	rr := httptest.NewRecorder()

	// Mock database
	database.Init(":memory:")
	defer database.Close()

	// Execute HTTP handler
	http.HandlerFunc(auth.RegisterUserHandler).ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check response body
	expectedBody := `{"message":"Successfully registered"}`
	assert.Equal(t, expectedBody, rr.Body.String())
}

func TestLoginUserHandler_Success(t *testing.T) {
	// Create test user
	user := models.User{Username: "testuser", Password: "testpassword"}
	userJSON, _ := json.Marshal(user)

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP response recorder
	rr := httptest.NewRecorder()

	// Mock database
	database.Init(":memory:")
	defer database.Close()

	// Execute HTTP handler
	http.HandlerFunc(auth.LoginUserHandler).ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check response body
	assert.Contains(t, rr.Body.String(), "Successfully logged in")
}

func TestRequireValidToken_Success(t *testing.T) {
	// Create valid token
	validToken := "valid_token"

	// Mock HTTP request with valid token in Authorization header
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+validToken)

	// Mock HTTP response recorder
	rr := httptest.NewRecorder()

	// Mock database
	database.Init(":memory:")
	defer database.Close()

	mockValidator := &mockTokenValidator{}

	// Execute function
	authenticated := RequireValidTokenMocked(rr, req, mockValidator)

	// Check if authenticated is true
	assert.True(t, authenticated)
}

func TestRequireValidToken_Failure(t *testing.T) {
	// Mock HTTP request without Authorization header
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP response recorder
	rr := httptest.NewRecorder()

	// Execute function
	authenticated := auth.RequireValidToken(rr, req)

	// Check if authenticated is false
	assert.False(t, authenticated)

	// Check status code in response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func RequireValidTokenMocked(w http.ResponseWriter, r *http.Request, validator tokenValidator) bool {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		logger.Println("Authorization header is missing")

		utils.WriteErrorJSONResponse(w, errors.New("authorization header is missing"), http.StatusInternalServerError)
		return false
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	authenticated, err := validator.IsTokenValid(database.GetDB(), tokenString)
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
