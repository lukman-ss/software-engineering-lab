package bank

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestBankTransfer_VulnerableEndpoint(t *testing.T) {
	server := NewBankServer([]byte("01234567890123456789012345678901"))

	form := url.Values{}
	form.Set("to", "acc-attacker")
	form.Set("amount", "200")

	// Attacker cross-origin form post carrying victim session cookie
	req := httptest.NewRequest(http.MethodPost, "/transfer/vulnerable", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://evil.com")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	rec := httptest.NewRecorder()
	server.HandleTransferVulnerable(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on vulnerable endpoint, got %d", rec.Code)
	}

	victim, _ := server.GetAccount("acc-victim")
	if victim.Balance != 800 {
		t.Fatalf("expected balance 800, got %d", victim.Balance)
	}

	attacker, _ := server.GetAccount("acc-attacker")
	if attacker.Balance != 250 {
		t.Fatalf("expected balance 250, got %d", attacker.Balance)
	}
}
