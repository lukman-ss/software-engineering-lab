package httputil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"labs/25-rate-limiting-and-backpressure/internal/ratelimit"
)

func TestRateLimitMiddleware_RFC6585(t *testing.T) {
	registry := ratelimit.NewRegistry(1, 10) // 1 token cap

	handler := RateLimitMiddleware(registry, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))

	// 1st request -> 200 OK
	req1 := httptest.NewRequest("GET", "/api/v1/resource", nil)
	req1.Header.Set("X-API-Key", "tenant-123")
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec1.Code)
	}

	// 2nd request -> 429 Too Many Requests
	req2 := httptest.NewRequest("GET", "/api/v1/resource", nil)
	req2.Header.Set("X-API-Key", "tenant-123")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests, got %d", rec2.Code)
	}

	retryAfter := rec2.Header().Get("Retry-After")
	if retryAfter == "" {
		t.Fatalf("expected Retry-After header in 429 response")
	}
}

func TestRateLimitMiddleware_AnonymousFallbackAndBody(t *testing.T) {
	registry := ratelimit.NewRegistry(1, 1)

	handler := RateLimitMiddleware(registry, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Request without X-API-Key should fallback to anonymous
	req1 := httptest.NewRequest("GET", "/api/v1/resource", nil)
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for 1st anonymous, got %d", rec1.Code)
	}

	// 2nd anonymous request -> 429
	req2 := httptest.NewRequest("GET", "/api/v1/resource", nil)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 for 2nd anonymous, got %d", rec2.Code)
	}

	if ct := rec2.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}
}
