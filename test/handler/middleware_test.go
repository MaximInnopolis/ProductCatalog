package handler_test

import (
	mock_service "github.com/MaximInnopolis/ProductCatalog/internal/service/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaximInnopolis/ProductCatalog/internal/handler"
	"github.com/MaximInnopolis/ProductCatalog/internal/service"
	"github.com/golang/mock/gomock"
)

// TestRequestIDMiddleware tests the RequestIDMiddleware function
func TestRequestIDMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h := &handler.Handler{}
	handlerFunc := h.RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		endpoint := ctx.Value("endpoint").(string)
		if endpoint != "/test" {
			t.Errorf("Expected endpoint /test, got %s", endpoint)
		}
	}))

	handlerFunc.ServeHTTP(rr, req)
}

// TestRequireValidTokenMiddleware tests the RequireValidTokenMiddleware function
func TestRequireValidTokenMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)
	mockService := &service.Service{
		Authorization: mockAuthService,
	}
	h := handler.NewHandler(&service.Service{Authorization: mockService})

	// Test case: Missing Authorization header
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handlerFunc := h.RequireValidTokenMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when Authorization header is missing")
	}))

	handlerFunc.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	// Test case: Invalid token
	req, err = http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer invalid_token")

	rr = httptest.NewRecorder()

	mockAuthService.EXPECT().IsTokenValid(gomock.Any(), "invalid_token").Return(false, nil)

	handlerFunc.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, status)
	}
}
