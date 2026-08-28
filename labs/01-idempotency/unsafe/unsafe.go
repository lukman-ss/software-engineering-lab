// Package unsafe demonstrates the problem: without idempotency,
// retried payment requests cause duplicate charges.
package unsafe

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// PaymentRequest represents a payment request.
type PaymentRequest struct {
	OrderID string
	Amount  int64
}

// PaymentResult represents the result of a payment.
type PaymentResult struct {
	PaymentID  string
	Status     string
	Amount     int64
	ExternalID string
}

// Gateway simulates an external payment gateway.
type Gateway interface {
	Charge(ctx context.Context, req PaymentRequest) (PaymentResult, error)
}

// MockGateway simulates a payment gateway with call tracking.
type MockGateway struct {
	chargeCount   int64
	callHistory   []PaymentRequest
	mu            sync.Mutex
	processingSem chan struct{} // For synchronization in tests
}

func (g *MockGateway) Charge(ctx context.Context, req PaymentRequest) (PaymentResult, error) {
	atomic.AddInt64(&g.chargeCount, 1)
	g.mu.Lock()
	g.callHistory = append(g.callHistory, req)
	g.mu.Unlock()

	// Simulate external gateway response
	return PaymentResult{
		PaymentID:  fmt.Sprintf("pay_%d", atomic.LoadInt64(&g.chargeCount)),
		Status:     "succeeded",
		Amount:     req.Amount,
		ExternalID: fmt.Sprintf("ext_%d", g.chargeCount),
	}, nil
}

// ChargeCount returns the number of gateway calls.
func (g *MockGateway) ChargeCount() int64 {
	return atomic.LoadInt64(&g.chargeCount)
}

// GetCallHistory returns the history of all charge requests.
func (g *MockGateway) GetCallHistory() []PaymentRequest {
	g.mu.Lock()
	defer g.mu.Unlock()
	result := make([]PaymentRequest, len(g.callHistory))
	copy(result, g.callHistory)
	return result
}

// Service processes payments. This implementation is NOT idempotent.
type Service struct {
	gateway Gateway
	store   *store
}

// store holds payment results (simulating database rows).
type store struct {
	mu       sync.RWMutex
	payments map[string]PaymentResult // key: payment ID
}

func newStore() *store {
	return &store{
		payments: make(map[string]PaymentResult),
	}
}

// Service processes payments without idempotency protection.
func NewService(g Gateway) *Service {
	return &Service{
		gateway: g,
		store:   newStore(),
	}
}

// ProcessPayment charges the gateway and stores the result.
// BUG: This implementation has no idempotency. If called twice with
// the same payment request (due to client retry), it creates two charges.
func (s *Service) ProcessPayment(ctx context.Context, req PaymentRequest) (PaymentResult, error) {
	// Charge the external gateway
	result, err := s.gateway.Charge(ctx, req)
	if err != nil {
		return PaymentResult{}, fmt.Errorf("gateway charge failed: %w", err)
	}

	// Store the result (simulating database insert)
	s.store.mu.Lock()
	s.store.payments[result.PaymentID] = result
	s.store.mu.Unlock()

	return result, nil
}

// GetPayments returns all stored payments.
func (s *Service) GetPayments() []PaymentResult {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()

	results := make([]PaymentResult, 0, len(s.store.payments))
	for _, p := range s.store.payments {
		results = append(results, p)
	}
	return results
}

// CountPayments returns the number of stored payments.
func (s *Service) CountPayments() int {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	return len(s.store.payments)
}
