package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lukman/software-engineer-lab/pkg/middleware"
)

func TestRequestIDMiddleware(t *testing.T) {
	// Handler that echoes the request ID from context
	handler := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetRequestID(r.Context())
		if reqID == "" {
			t.Fatal("expected request ID in context")
		}
		w.Header().Set("X-Echo-ID", reqID)
		w.WriteHeader(http.StatusOK)
	}))

	// Test 1: Generate new ID when none provided
	req1 := httptest.NewRequest("GET", "/test", nil)
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	respID1 := rec1.Header().Get(middleware.HeaderRequestID)
	if respID1 == "" {
		t.Fatal("expected X-Request-ID response header")
	}

	// Test 2: Preserve valid incoming ID
	req2 := httptest.NewRequest("GET", "/test", nil)
	customID := "custom-trace-id-123"
	req2.Header.Set(middleware.HeaderRequestID, customID)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	respID2 := rec2.Header().Get(middleware.HeaderRequestID)
	if respID2 != customID {
		t.Fatalf("expected preserved ID %s, got %s", customID, respID2)
	}

	// Test 3: Reject invalid incoming ID and generate new
	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.Header.Set(middleware.HeaderRequestID, "invalid id with spaces & symbols!")
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)

	respID3 := rec3.Header().Get(middleware.HeaderRequestID)
	if respID3 == "" || respID3 == "invalid id with spaces & symbols!" {
		t.Fatalf("expected valid generated ID, got %s", respID3)
	}

	t.Log("Request ID middleware tests passed successfully")
}
