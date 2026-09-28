package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS_NoOrigin(t *testing.T) {
	mw := New(Config{AllowedOrigins: []string{"https://bank.com"}})
	handler := mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("expected no CORS header for non-CORS request")
	}
}

func TestCORS_DisallowedOrigin_Preflight(t *testing.T) {
	mw := New(Config{AllowedOrigins: []string{"https://trusted.com"}})
	handler := mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/data", nil)
	req.Header.Set("Origin", "https://evil.com")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for disallowed preflight origin, got %d", rec.Code)
	}
}

func TestCORS_Preflight_Success(t *testing.T) {
	mw := New(Config{
		AllowedOrigins: []string{"https://app.bank.com"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"X-CSRF-Token", "Content-Type"},
		MaxAge:         3600,
	})
	handler := mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/data", nil)
	req.Header.Set("Origin", "https://app.bank.com")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content for preflight, got %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "https://app.bank.com" {
		t.Fatalf("unexpected origin header: %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Header().Get("Access-Control-Max-Age") != "3600" {
		t.Fatalf("unexpected max age: %s", rec.Header().Get("Access-Control-Max-Age"))
	}
}

func TestCORS_Credentials_With_Wildcard_DisallowedInSpec(t *testing.T) {
	// When credentials are true, origin must be reflected explicitly, never wildcard "*"
	mw := New(Config{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	handler := mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	req.Header.Set("Origin", "https://anywhere.com")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") == "*" {
		t.Fatalf("spec violation: Access-Control-Allow-Origin must not be '*' when credentials are true")
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "https://anywhere.com" {
		t.Fatalf("expected reflected origin, got %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatalf("expected credentials true")
	}
}
