package bank

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"labs/36-cors-and-csrf/internal/csrf"
)

type Account struct {
	ID      string `json:"id"`
	Owner   string `json:"owner"`
	Balance int64  `json:"balance"`
}

type BankServer struct {
	mu           sync.RWMutex
	accounts     map[string]*Account
	sessions     map[string]string // cookie session_id -> account_id
	tokenManager *csrf.TokenManager
}

func NewBankServer(secret []byte) *BankServer {
	return &BankServer{
		accounts: map[string]*Account{
			"acc-victim":   {ID: "acc-victim", Owner: "Victim User", Balance: 1000},
			"acc-attacker": {ID: "acc-attacker", Owner: "Attacker User", Balance: 50},
		},
		sessions: map[string]string{
			"session-victim-secret": "acc-victim",
		},
		tokenManager: csrf.NewTokenManager(secret, 1*time.Hour),
	}
}

func (b *BankServer) TokenManager() *csrf.TokenManager {
	return b.tokenManager
}

func (b *BankServer) GetAccount(id string) (*Account, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	acc, ok := b.accounts[id]
	if !ok {
		return nil, false
	}
	// return copy
	return &Account{ID: acc.ID, Owner: acc.Owner, Balance: acc.Balance}, true
}

func (b *BankServer) authenticate(r *http.Request) (*Account, string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		return nil, "", fmt.Errorf("missing session cookie")
	}

	b.mu.RLock()
	accountID, ok := b.sessions[cookie.Value]
	b.mu.RUnlock()

	if !ok {
		return nil, "", fmt.Errorf("invalid session")
	}

	acc, ok := b.GetAccount(accountID)
	if !ok {
		return nil, "", fmt.Errorf("account not found")
	}
	return acc, cookie.Value, nil
}

// HandleGetCSRFToken returns a fresh signed CSRF token for the authenticated user.
func (b *BankServer) HandleGetCSRFToken(w http.ResponseWriter, r *http.Request) {
	_, sessionID, err := b.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token := b.tokenManager.GenerateToken(sessionID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"csrf_token": token})
}

// HandleBalance returns the user's account details.
func (b *BankServer) HandleBalance(w http.ResponseWriter, r *http.Request) {
	acc, _, err := b.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

// HandleTransferVulnerable processes transfer WITHOUT CSRF verification (VULNERABLE).
func (b *BankServer) HandleTransferVulnerable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	acc, _, err := b.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	toID := r.FormValue("to")
	amountStr := r.FormValue("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	sender := b.accounts[acc.ID]
	recipient, ok := b.accounts[toID]
	if !ok {
		http.Error(w, "Recipient not found", http.StatusBadRequest)
		return
	}

	if sender.Balance < amount {
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
		return
	}

	sender.Balance -= amount
	recipient.Balance += amount

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("Transferred %d from %s to %s", amount, sender.ID, recipient.ID),
		"balance": sender.Balance,
	})
}

// HandleTransferProtected processes transfer WITH CSRF verification (PROTECTED).
func (b *BankServer) HandleTransferProtected(w http.ResponseWriter, r *http.Request) {
	// Business logic is identical to HandleTransferVulnerable; CSRF middleware handles protection wrapper.
	b.HandleTransferVulnerable(w, r)
}
