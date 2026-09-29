package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"labs/36-cors-and-csrf/internal/bank"
	"labs/36-cors-and-csrf/internal/cors"
	"labs/36-cors-and-csrf/internal/csrf"
)

func setupBankApp() (http.Handler, *bank.BankServer) {
	secret := []byte("a-very-secret-key-32-bytes-long!")
	bankServer := bank.NewBankServer(secret)

	corsMW := cors.New(cors.Config{
		AllowedOrigins:   []string{"https://victim-client.com"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "X-CSRF-Token", "X-Requested-With"},
		AllowCredentials: true,
	})

	csrfMW := csrf.NewMiddleware(bankServer.TokenManager())

	mux := http.NewServeMux()

	// 1. Balance endpoint
	mux.HandleFunc("GET /api/balance", bankServer.HandleBalance)

	// 2. Token endpoint
	mux.HandleFunc("GET /api/csrf-token", bankServer.HandleGetCSRFToken)

	// 3. Vulnerable transfer (No CSRF check, CORS wrapped)
	mux.HandleFunc("POST /api/transfer/vulnerable", bankServer.HandleTransferVulnerable)

	// 4. Protected transfer (CSRF token required, CORS wrapped)
	mux.Handle("POST /api/transfer/protected", csrfMW.Handler(http.HandlerFunc(bankServer.HandleTransferProtected)))

	// 5. Fetch-Metadata protected transfer
	mux.Handle("POST /api/transfer/fetch-metadata", csrf.FetchMetadataMiddleware(http.HandlerFunc(bankServer.HandleTransferProtected)))

	// 6. JSON-only custom header protected transfer
	customHeaderMW := csrf.RequireCustomHeaderMiddleware("X-Requested-With", "XMLHttpRequest")
	mux.Handle("POST /api/transfer/custom-header", customHeaderMW(http.HandlerFunc(bankServer.HandleTransferProtected)))

	return corsMW.Handler(mux), bankServer
}

func TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution(t *testing.T) {
	app, bankServer := setupBankApp()

	// Scenario: Attacker website sends POST to vulnerable endpoint with Victim's cookie.
	// Origin is evil.com (NOT allowed in CORS middleware).
	form := url.Values{}
	form.Set("to", "acc-attacker")
	form.Set("amount", "300")

	req := httptest.NewRequest(http.MethodPost, "/api/transfer/vulnerable", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://evil-attacker.com")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	// The server processed the transfer!
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK because CORS does NOT stop server execution, got %d", rec.Code)
	}

	// Verify funds were stolen
	victim, _ := bankServer.GetAccount("acc-victim")
	if victim.Balance != 700 {
		t.Fatalf("expected victim balance 700, got %d", victim.Balance)
	}

	// Notice: Even though CORS headers are omitted or rejected for evil-attacker.com,
	// the state mutation already completed in the database.
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("disallowed origin should not receive Access-Control-Allow-Origin")
	}
}

func TestIntegration_CSRF_Token_Prevents_Attack(t *testing.T) {
	app, bankServer := setupBankApp()

	// Attacker tries to forge request to /api/transfer/protected
	form := url.Values{}
	form.Set("to", "acc-attacker")
	form.Set("amount", "200")

	req := httptest.NewRequest(http.MethodPost, "/api/transfer/protected", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://evil-attacker.com")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	// Blocked by CSRF Middleware
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", rec.Code)
	}

	// Verify funds safe
	victim, _ := bankServer.GetAccount("acc-victim")
	if victim.Balance != 1000 {
		t.Fatalf("expected balance untouched (1000), got %d", victim.Balance)
	}
}

func TestIntegration_Legitimate_Flow_With_CSRF_Token(t *testing.T) {
	app, bankServer := setupBankApp()

	// 1. Fetch CSRF token as legitimate client
	reqToken := httptest.NewRequest(http.MethodGet, "/api/csrf-token", nil)
	reqToken.Header.Set("Origin", "https://victim-client.com")
	reqToken.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	recToken := httptest.NewRecorder()
	app.ServeHTTP(recToken, reqToken)

	if recToken.Code != http.StatusOK {
		t.Fatalf("failed to get token: %d", recToken.Code)
	}

	var resp map[string]string
	json.NewDecoder(recToken.Body).Decode(&resp)
	token := resp["csrf_token"]
	if token == "" {
		t.Fatalf("received empty token")
	}

	// 2. Perform protected transfer with valid token
	form := url.Values{}
	form.Set("to", "acc-attacker")
	form.Set("amount", "150")
	form.Set("csrf_token", token)

	reqTransfer := httptest.NewRequest(http.MethodPost, "/api/transfer/protected", strings.NewReader(form.Encode()))
	reqTransfer.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqTransfer.Header.Set("Origin", "https://victim-client.com")
	reqTransfer.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	recTransfer := httptest.NewRecorder()
	app.ServeHTTP(recTransfer, reqTransfer)

	if recTransfer.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for legit transfer, got %d (body: %s)", recTransfer.Code, recTransfer.Body.String())
	}

	victim, _ := bankServer.GetAccount("acc-victim")
	if victim.Balance != 850 {
		t.Fatalf("expected balance 850, got %d", victim.Balance)
	}
}

func TestIntegration_SecFetchSite_Protection(t *testing.T) {
	app, bankServer := setupBankApp()

	form := url.Values{}
	form.Set("to", "acc-attacker")
	form.Set("amount", "100")

	// Cross-site request rejected by Sec-Fetch-Site
	req := httptest.NewRequest(http.MethodPost, "/api/transfer/fetch-metadata", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for cross-site Sec-Fetch-Site, got %d", rec.Code)
	}

	victim, _ := bankServer.GetAccount("acc-victim")
	if victim.Balance != 1000 {
		t.Fatalf("expected balance untouched (1000), got %d", victim.Balance)
	}
}

func TestIntegration_Concurrency_RaceCondition(t *testing.T) {
	app, _ := setupBankApp()

	var wg sync.WaitGroup
	workers := 20

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/api/csrf-token", nil)
			req.Header.Set("Origin", "https://victim-client.com")
			req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})
			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, req)
		}()
	}

	wg.Wait()
}

func TestIntegration_CustomHeader_Protection(t *testing.T) {
	app, bankServer := setupBankApp()

	form := url.Values{}
	form.Set("to", "acc-attacker")
	form.Set("amount", "100")

	// Missing header -> 403
	req := httptest.NewRequest(http.MethodPost, "/api/transfer/custom-header", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for missing custom header, got %d", rec.Code)
	}

	victim, _ := bankServer.GetAccount("acc-victim")
	if victim.Balance != 1000 {
		t.Fatalf("expected balance 1000, got %d", victim.Balance)
	}

	// Valid header -> 200
	reqValid := httptest.NewRequest(http.MethodPost, "/api/transfer/custom-header", strings.NewReader(form.Encode()))
	reqValid.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqValid.Header.Set("X-Requested-With", "XMLHttpRequest")
	reqValid.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	recValid := httptest.NewRecorder()
	app.ServeHTTP(recValid, reqValid)

	if recValid.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for custom header request, got %d", recValid.Code)
	}

	victim, _ = bankServer.GetAccount("acc-victim")
	if victim.Balance != 900 {
		t.Fatalf("expected balance 900, got %d", victim.Balance)
	}
}

func TestIntegration_CSRF_Token_In_Header(t *testing.T) {
	app, bankServer := setupBankApp()

	// 1. Fetch CSRF token
	reqToken := httptest.NewRequest(http.MethodGet, "/api/csrf-token", nil)
	reqToken.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})
	recToken := httptest.NewRecorder()
	app.ServeHTTP(recToken, reqToken)

	var resp map[string]string
	json.NewDecoder(recToken.Body).Decode(&resp)
	token := resp["csrf_token"]

	// 2. Submit token via X-CSRF-Token header
	form := url.Values{}
	form.Set("to", "acc-attacker")
	form.Set("amount", "100")

	req := httptest.NewRequest(http.MethodPost, "/api/transfer/protected", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-Token", token)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for header-based CSRF token, got %d", rec.Code)
	}

	victim, _ := bankServer.GetAccount("acc-victim")
	if victim.Balance != 900 {
		t.Fatalf("expected balance 900, got %d", victim.Balance)
	}
}

func TestIntegration_CrossSession_Token_Reuse_Rejected(t *testing.T) {
	app, bankServer := setupBankApp()

	// Token generated for session-victim-secret
	token := bankServer.TokenManager().GenerateToken("session-victim-secret")

	// Attacker tries to use victim's token with attacker's session
	form := url.Values{}
	form.Set("to", "acc-attacker")
	form.Set("amount", "100")
	form.Set("csrf_token", token)

	req := httptest.NewRequest(http.MethodPost, "/api/transfer/protected", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-attacker-secret"})

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for cross-session CSRF token reuse, got %d", rec.Code)
	}
}

func TestIntegration_SecFetchSite_SameOrigin_Allowed(t *testing.T) {
	app, bankServer := setupBankApp()

	form := url.Values{}
	form.Set("to", "acc-attacker")
	form.Set("amount", "100")

	req := httptest.NewRequest(http.MethodPost, "/api/transfer/fetch-metadata", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for same-origin Sec-Fetch-Site, got %d", rec.Code)
	}

	victim, _ := bankServer.GetAccount("acc-victim")
	if victim.Balance != 900 {
		t.Fatalf("expected balance 900, got %d", victim.Balance)
	}
}
