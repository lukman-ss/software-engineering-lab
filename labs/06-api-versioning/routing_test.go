package api_versioning

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRouting_ErrorBehaviors memverifikasi konsistensi handling error API
// dan memastikan response JSON sesuai dengan standar.
func TestRouting_ErrorBehaviors(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)
	mux.HandleFunc("/api/v2/invoices/", V2Handler)

	server := httptest.NewServer(mux)
	defer server.Close()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedKey    string // string JSON key error, etc
	}{
		{
			name:           "Valid ID returns 200",
			method:         http.MethodGet,
			path:           "/api/v1/invoices/1001",
			expectedStatus: http.StatusOK,
			expectedKey:    "id",
		},
		{
			name:           "Missing ID returns 400",
			method:         http.MethodGet,
			path:           "/api/v1/invoices/",
			expectedStatus: http.StatusBadRequest,
			expectedKey:    "error",
		},
		{
			name:           "Invalid numeric ID returns 400",
			method:         http.MethodGet,
			path:           "/api/v1/invoices/abc",
			expectedStatus: http.StatusBadRequest,
			expectedKey:    "error",
		},
		{
			name:           "Extra path segments returns 400",
			method:         http.MethodGet,
			path:           "/api/v1/invoices/1001/extra",
			expectedStatus: http.StatusBadRequest,
			expectedKey:    "error",
		},
		{
			name:           "Non-existent ID returns 404",
			method:         http.MethodGet,
			path:           "/api/v1/invoices/9999", // ID lain selain 1001
			expectedStatus: http.StatusNotFound,
			expectedKey:    "error",
		},
		{
			name:           "Unsupported method returns 405",
			method:         http.MethodPost,
			path:           "/api/v1/invoices/1001",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedKey:    "error",
		},
		{
			name:           "Unknown route returns 404 (handled by mux)",
			method:         http.MethodGet,
			path:           "/api/v3/invoices/1001",
			expectedStatus: http.StatusNotFound,
			expectedKey:    "", // mux doesn't return JSON by default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, server.URL+tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Validate Content-Type for all responses mapped to handlers
			if resp.StatusCode != http.StatusNotFound || tt.path == "/api/v1/invoices/9999" {
				// Mux 404 is plain text, but our 404 is JSON
				if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
					t.Errorf("expected Content-Type application/json, got %s", ct)
				}
			}

			// Verify if error key exists in response body when expected
			if tt.expectedKey != "" {
				var raw map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
					t.Errorf("failed to decode JSON response: %v", err)
				} else {
					if _, ok := raw[tt.expectedKey]; !ok {
						t.Errorf("expected JSON to contain key %q, got: %v", tt.expectedKey, raw)
					}
				}
			}
		})
	}
}