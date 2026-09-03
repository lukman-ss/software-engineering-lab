package validation_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lukman-ss/software-engineering-lab/pkg/validation"
)

type testPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Count int    `json:"count"`
}

func TestReadJSONValid(t *testing.T) {
	body := `{"name":"test user","email":"test@example.com","count":42}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	var payload testPayload
	err := validation.ReadJSON(req, &payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if payload.Name != "test user" {
		t.Fatalf("expected name 'test user', got %s", payload.Name)
	}
	if payload.Email != "test@example.com" {
		t.Fatalf("expected email 'test@example.com', got %s", payload.Email)
	}
	if payload.Count != 42 {
		t.Fatalf("expected count 42, got %d", payload.Count)
	}
}

func TestReadJSONMalformed(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{"empty body", "", true},
		{"truncated JSON", `{"name": "test"`, true},
		{"malformed object", `{invalid}`, true},
		{"trailing data", `{"name":"test"} garbage`, true},
		{"unquoted key", `{name: "test"}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			var payload testPayload
			err := validation.ReadJSON(req, &payload)
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestReadJSONUnknownField(t *testing.T) {
	body := `{"name":"test","email":"test@example.com","count":1,"unknown_field":"should fail"}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	var payload testPayload
	err := validation.ReadJSON(req, &payload)
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestReadJSONBodyTooLarge(t *testing.T) {
	// Create a body larger than MaxBodySize
	largeBody := bytes.Repeat([]byte("x"), validation.MaxBodySize+1)
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(largeBody))
	req.Header.Set("Content-Type", "application/json")

	var payload testPayload
	err := validation.ReadJSON(req, &payload)
	if err == nil {
		t.Fatal("expected error for body too large")
	}
}

func TestReadJSONTypeMismatch(t *testing.T) {
	body := `{"name":"test","email":"test@example.com","count":"not a number"}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	var payload testPayload
	err := validation.ReadJSON(req, &payload)
	if err == nil {
		t.Fatal("expected error for type mismatch")
	}
}
