package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"labs/36-cors-and-csrf/internal/bank"
	"labs/36-cors-and-csrf/internal/cors"
	"labs/36-cors-and-csrf/internal/csrf"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("  LAB 36: CORS & CSRF DEMONSTRATION & PROOF      ")
	fmt.Println("==================================================")

	secret := []byte("demo-secret-key-32-bytes-long!")
	bankServer := bank.NewBankServer(secret)

	corsMW := cors.New(cors.Config{
		AllowedOrigins:   []string{"https://legit-app.com"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})
	csrfMW := csrf.NewMiddleware(bankServer.TokenManager())

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/csrf-token", bankServer.HandleGetCSRFToken)
	mux.HandleFunc("POST /api/transfer/vulnerable", bankServer.HandleTransferVulnerable)
	mux.Handle("POST /api/transfer/protected", csrfMW.Handler(http.HandlerFunc(bankServer.HandleTransferProtected)))
	handler := corsMW.Handler(mux)

	// 1. Initial State
	victim, _ := bankServer.GetAccount("acc-victim")
	attacker, _ := bankServer.GetAccount("acc-attacker")
	fmt.Printf("[1] Initial State:\n    Victim Balance  : $%d\n    Attacker Balance: $%d\n\n", victim.Balance, attacker.Balance)

	// 2. Myth Busting: CORS blocking read does NOT stop server state execution!
	fmt.Println("[2] DEMO: Cross-Origin Attack on Vulnerable Endpoint (Origin: https://evil.com)")
	formVulnerable := url.Values{}
	formVulnerable.Set("to", "acc-attacker")
	formVulnerable.Set("amount", "400")

	reqVulnerable := httptest.NewRequest(http.MethodPost, "/api/transfer/vulnerable", strings.NewReader(formVulnerable.Encode()))
	reqVulnerable.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqVulnerable.Header.Set("Origin", "https://evil.com") // Unauthorized origin!
	reqVulnerable.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	recVulnerable := httptest.NewRecorder()
	handler.ServeHTTP(recVulnerable, reqVulnerable)

	victim, _ = bankServer.GetAccount("acc-victim")
	attacker, _ = bankServer.GetAccount("acc-attacker")
	fmt.Printf("    HTTP Response Code : %d\n", recVulnerable.Code)
	fmt.Printf("    CORS Header Present: %q\n", recVulnerable.Header().Get("Access-Control-Allow-Origin"))
	fmtPrintfState("RESULT", "Transfer EXECUTED despite CORS header absence!", victim.Balance, attacker.Balance)

	// 3. Protected Endpoint Attack Attempt
	fmt.Println("\n[3] DEMO: Cross-Origin Attack on Protected Endpoint (No CSRF Token)")
	formProtected := url.Values{}
	formProtected.Set("to", "acc-attacker")
	formProtected.Set("amount", "500")

	reqProtected := httptest.NewRequest(http.MethodPost, "/api/transfer/protected", strings.NewReader(formProtected.Encode()))
	reqProtected.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqProtected.Header.Set("Origin", "https://evil.com")
	reqProtected.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	recProtected := httptest.NewRecorder()
	handler.ServeHTTP(recProtected, reqProtected)

	victim, _ = bankServer.GetAccount("acc-victim")
	attacker, _ = bankServer.GetAccount("acc-attacker")
	fmt.Printf("    HTTP Response Code : %d\n", recProtected.Code)
	fmtPrintfState("RESULT", "Transfer BLOCKED by Anti-CSRF Token Middleware!", victim.Balance, attacker.Balance)

	// 4. Legitimate Client Flow
	fmt.Println("\n[4] DEMO: Legitimate Client Flow with Valid CSRF Token")
	reqToken := httptest.NewRequest(http.MethodGet, "/api/csrf-token", nil)
	reqToken.Header.Set("Origin", "https://legit-app.com")
	reqToken.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	recToken := httptest.NewRecorder()
	handler.ServeHTTP(recToken, reqToken)

	var tokenResp map[string]string
	json.NewDecoder(recToken.Body).Decode(&tokenResp)
	validToken := tokenResp["csrf_token"]
	fmt.Printf("    Obtained Signed CSRF Token: %s...\n", validToken[:25])

	formLegit := url.Values{}
	formLegit.Set("to", "acc-attacker")
	formLegit.Set("amount", "100")
	formLegit.Set("csrf_token", validToken)

	reqLegit := httptest.NewRequest(http.MethodPost, "/api/transfer/protected", strings.NewReader(formLegit.Encode()))
	reqLegit.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqLegit.Header.Set("Origin", "https://legit-app.com")
	reqLegit.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

	recLegit := httptest.NewRecorder()
	handler.ServeHTTP(recLegit, reqLegit)

	victim, _ = bankServer.GetAccount("acc-victim")
	attacker, _ = bankServer.GetAccount("acc-attacker")
	fmtPrintfState("RESULT", "Legitimate Transfer Succeeded!", victim.Balance, attacker.Balance)
	fmt.Println("\n==================================================")
}

func fmtPrintfState(tag, msg string, vBal, aBal int64) {
	fmt.Printf("    [%s] %s\n", tag, msg)
	fmt.Printf("    Victim Balance  : $%d\n", vBal)
	fmt.Printf("    Attacker Balance: $%d\n", aBal)
}
