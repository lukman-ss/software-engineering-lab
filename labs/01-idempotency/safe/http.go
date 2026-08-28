package safe

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// HTTPHandler adapts the payment service to HTTP.
type HTTPHandler struct {
	svc *Service
}

func NewPaymentHandler(svc *Service) *HTTPHandler {
	return &HTTPHandler{svc: svc}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/payments":
		h.handlePayment(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *HTTPHandler) handlePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		idempotencyKey = r.URL.Query().Get("idempotency_key")
	}
	if idempotencyKey == "" {
		http.Error(w, "missing Idempotency-Key header or query param", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req PaymentRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "parse body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.OrderID == "" {
		http.Error(w, "order_id required", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		http.Error(w, "amount must be positive", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Pass method and path for canonical request fingerprinting
	result, status, err := h.svc.ProcessPayment(ctx, r.Method, r.URL.Path, idempotencyKey, req)
	if err != nil {
		if IsIdempotentKeyConflict(err) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "payment failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}