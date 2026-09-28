package csrf

import (
	"crypto/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestTokenManager_GenerateAndValidate(t *testing.T) {
	secret := make([]byte, 32)
	rand.Read(secret)
	tm := NewTokenManager(secret, 10*time.Minute)

	sessionID := "user-session-12345"
	token := tm.GenerateToken(sessionID)

	if err := tm.ValidateToken(token, sessionID); err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}

	// Session mismatch
	if err := tm.ValidateToken(token, "other-session"); err == nil {
		t.Fatalf("expected error on session mismatch")
	}

	// Tampered token
	tampered := token[:len(token)-4] + "AAAA"
	if err := tm.ValidateToken(tampered, sessionID); err == nil {
		t.Fatalf("expected error on tampered token")
	}
}

func TestTokenManager_ExpiredToken(t *testing.T) {
	secret := make([]byte, 32)
	rand.Read(secret)
	tm := NewTokenManager(secret, 10*time.Millisecond)

	sessionID := "user-session-12345"
	token := tm.GenerateToken(sessionID)

	time.Sleep(20 * time.Millisecond)
	if err := tm.ValidateToken(token, sessionID); err != ErrExpiredToken {
		t.Fatalf("expected ErrExpiredToken, got %v", err)
	}
}

func TestCSRFMiddleware_FormPost(t *testing.T) {
	secret := make([]byte, 32)
	rand.Read(secret)
	tm := NewTokenManager(secret, 10*time.Minute)
	mw := NewMiddleware(tm)

	handler := mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("action executed"))
	}))

	sessionID := "session-abc"
	token := tm.GenerateToken(sessionID)

	// Valid Form Post
	form := url.Values{}
	form.Set("csrf_token", token)
	req := httptest.NewRequest(http.MethodPost, "/transfer", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}

	// Missing CSRF Token (Simulated Attack)
	attackForm := url.Values{}
	attackForm.Set("amount", "100")
	reqAttack := httptest.NewRequest(http.MethodPost, "/transfer", strings.NewReader(attackForm.Encode()))
	reqAttack.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqAttack.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})

	recAttack := httptest.NewRecorder()
	handler.ServeHTTP(recAttack, reqAttack)

	if recAttack.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for missing token, got %d", recAttack.Code)
	}
}

func TestFetchMetadataMiddleware(t *testing.T) {
	handler := FetchMetadataMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Cross-site POST rejected
	req := httptest.NewRequest(http.MethodPost, "/action", nil)
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for cross-site POST, got %d", rec.Code)
	}

	// Same-origin POST allowed
	req2 := httptest.NewRequest(http.MethodPost, "/action", nil)
	req2.Header.Set("Sec-Fetch-Site", "same-origin")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 for same-origin POST, got %d", rec2.Code)
	}
}
