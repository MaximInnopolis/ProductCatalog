package handler

import (
	"bytes"
	"encoding/json"
	"github.com/MaximInnopolis/ProductCatalog/internal/handler"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaximInnopolis/ProductCatalog/internal/model"
	"github.com/MaximInnopolis/ProductCatalog/internal/service"
	mock_service "github.com/MaximInnopolis/ProductCatalog/internal/service/mock"
	"github.com/golang/mock/gomock"
)

func TestRegisterUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)
	mockService := &service.Service{
		Authorization: mockAuthService,
	}

	h := handler.NewHandler(mockService)

	// Mock user data
	user := &model.User{
		Username: "testuser",
		Password: "testpassword",
	}

	// Mock request body
	reqBody, _ := json.Marshal(user)

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP response recorder
	w := httptest.NewRecorder()

	// Set expectation on the mock Authorization service
	mockAuthService.EXPECT().CreateUser(gomock.Any(), user).Return(nil)

	// Call the handler function
	h.RegisterUserHandler(w, req)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestLoginUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)
	mockService := &service.Service{
		Authorization: mockAuthService,
	}

	h := handler.NewHandler(mockService)

	// Mock user data
	user := &model.User{
		Username: "testuser",
		Password: "testpassword",
	}

	// Mock request body
	reqBody, _ := json.Marshal(user)

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP response recorder
	w := httptest.NewRecorder()

	// Mock token
	token := "mocktoken"

	// Set expectation on the mock Authorization service
	mockAuthService.EXPECT().GenerateToken(gomock.Any(), user).Return(token, nil)

	// Call the handler function
	h.LoginUserHandler(w, req)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
