package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaximInnopolis/ProductCatalog/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandlers(t *testing.T) {
	var testCases []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}

	// Start test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.RegisterHandlers()
	}))
	defer server.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, server.URL+tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			http.DefaultServeMux.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}
