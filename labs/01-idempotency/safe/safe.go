// Package safe implements idempotent payment processing.
// Demonstrates core idempotency concepts: idempotency keys, payload hashing,
// unique constraints, response caching, and conflict detection.
package safe

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrIdempotencyKeyConflict is returned when the same idempotency key
	// is used with a different payload.
	ErrIdempotencyKeyConflict = errors.New("idempotency key used with different payload")
	// ErrMissingIdempotencyKey is returned when no key is provided.
	ErrMissingIdempotencyKey = errors.New("idempotency key is required")
	// ErrRequestInProgress is returned when a concurrent request with the same
	// idempotency key is currently being processed.
	ErrRequestInProgress = errors.New("request already in progress")
)

// PaymentRequest represents a payment request from a client.
type PaymentRequest struct {
	OrderID string `json:"order_id"`
	Amount  int64  `json:"amount"`
}

// PaymentResult represents the stored payment and its gateway response.
type PaymentResult struct {
	PaymentID  string `json:"payment_id"`
	Status     string `json:"status"`
	Amount     int64  `json:"amount"`
	ExternalID string `json:"external_id"`
}

// IdempotencyStatus represents the state of an idempotency key.
type IdempotencyStatus string

const (
	StatusProcessing IdempotencyStatus = "processing"
	StatusCompleted  IdempotencyStatus = "completed"
)

// DefaultTTL is the default time-to-live for idempotency keys after completion.
const DefaultTTL = 72 * time.Hour

// PaymentRecord is the full database record including idempotency metadata.
type PaymentRecord struct {
	IdempotencyKey string            `json:"idempotency_key"`
	RequestHash    string            `json:"request_hash"`
	OrderID        string            `json:"order_id"`
	Status         IdempotencyStatus `json:"status"`
	ResponseStatus int               `json:"response_status"`
	ResponseJSON   string            `json:"response_json"`
	Amount         int64             `json:"amount"`
	CreatedAt      time.Time         `json:"created_at"`
	ExpiresAt      time.Time         `json:"expires_at"`
}

// Gateway simulates an external payment gateway.
type Gateway interface {
	Charge(ctx context.Context, req PaymentRequest) (PaymentResult, error)
}

// mockGateway simulates a payment gateway with call tracking.
type mockGateway struct {
	mu          sync.Mutex
	chargeCount int64
}

func (g *mockGateway) Charge(ctx context.Context, req PaymentRequest) (PaymentResult, error) {
	g.mu.Lock()
	g.chargeCount++
	count := g.chargeCount
	g.mu.Unlock()

	return PaymentResult{
		PaymentID:  fmt.Sprintf("pay_%d", count),
		Status:     "succeeded",
		Amount:     req.Amount,
		ExternalID: fmt.Sprintf("ext_%d", count),
	}, nil
}

func (g *mockGateway) ChargeCount() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.chargeCount
}

// MockGateway creates a test gateway with call tracking.
func MockGateway() Gateway {
	return &mockGateway{}
}

// GatewayChargeCount returns the charge count for a mock gateway.
func GatewayChargeCount(g Gateway) int64 {
	if mg, ok := g.(*mockGateway); ok {
		return mg.ChargeCount()
	}
	return 0
}

// Store is the persistence layer interface for idempotency and payments.
// Implementation uses in-memory map to simulate database unique constraint.
type Store interface {
	TryInsert(ctx context.Context, record PaymentRecord) (bool, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*PaymentRecord, error)
	UpdateCompleted(ctx context.Context, key string, result PaymentResult, responseStatus int) error
	UpdateCompletedResponse(ctx context.Context, key string, responseJSON string) error
}

type dbStore struct {
	mu       sync.RWMutex
	payments map[string]PaymentRecord
}

func newDBStore() *dbStore {
	return &dbStore{
		payments: make(map[string]PaymentRecord),
	}
}

func (s *dbStore) TryInsert(ctx context.Context, record PaymentRecord) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.payments[record.IdempotencyKey]
	if exists {
		if existing.RequestHash != record.RequestHash {
			return false, ErrIdempotencyKeyConflict
		}
		return false, nil // already exists with same hash, not a conflict
	}

	s.payments[record.IdempotencyKey] = record
	return true, nil
}

func (s *dbStore) GetByIdempotencyKey(ctx context.Context, key string) (*PaymentRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, exists := s.payments[key]
	if !exists {
		return nil, nil
	}
	return &record, nil
}

func (s *dbStore) UpdateCompleted(ctx context.Context, key string, result PaymentResult, responseStatus int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, exists := s.payments[key]
	if !exists {
		return fmt.Errorf("record not found for key: %s", key)
	}

	record.Status = StatusCompleted
	record.ResponseStatus = responseStatus
	s.payments[key] = record
	return nil
}

func (s *dbStore) UpdateCompletedResponse(ctx context.Context, key string, responseJSON string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, exists := s.payments[key]
	if !exists {
		return fmt.Errorf("record not found for key: %s", key)
	}

	record.ResponseJSON = responseJSON
	s.payments[key] = record
	return nil
}

// NewStore creates a new payment store.
func NewStore() Store {
	return newDBStore()
}

func canonicalPayload(req PaymentRequest) (string, error) {
	m := map[string]interface{}{
		"order_id": req.OrderID,
		"amount":   req.Amount,
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("marshal canonical payload: %w", err)
	}
	return string(b), nil
}

func hashRequest(method, path string, req PaymentRequest) string {
	canonical, err := canonicalPayload(req)
	if err != nil {
		return ""
	}
	fingerprint := fmt.Sprintf("%s %s %s", method, path, canonical)
	h := sha256.Sum256([]byte(fingerprint))
	return hex.EncodeToString(h[:])
}

// Service demonstrates application-level idempotent payment processing.
type Service struct {
	gateway Gateway
	store   Store
}

// NewService creates a new idempotent payment service.
func NewService(g Gateway) *Service {
	return &Service{
		gateway: g,
		store:   NewStore(),
	}
}

// ProcessPayment charges the gateway idempotently.
func (s *Service) ProcessPayment(ctx context.Context, method, path, idempotencyKey string, req PaymentRequest) (PaymentResult, int, error) {
	// Validate idempotency key is provided
	if idempotencyKey == "" {
		return PaymentResult{}, 400, ErrMissingIdempotencyKey
	}

	hash := hashRequest(method, path, req)

	// Try to insert new record atomically
	record := PaymentRecord{
		IdempotencyKey: idempotencyKey,
		RequestHash:    hash,
		OrderID:        req.OrderID,
		Status:         StatusProcessing,
		Amount:         req.Amount,
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(DefaultTTL),
	}

	inserted, err := s.store.TryInsert(ctx, record)
	if err != nil {
		return PaymentResult{}, 409, err
	}

	// If not inserted, check existing record
	if !inserted {
		existing, err := s.store.GetByIdempotencyKey(ctx, idempotencyKey)
		if err != nil {
			return PaymentResult{}, 500, fmt.Errorf("get existing key: %w", err)
		}

		// Payload hash mismatch - reject with 409 Conflict
		if existing.RequestHash != hash {
			return PaymentResult{}, 409, ErrIdempotencyKeyConflict
		}

		// Already completed - return cached response (response replay)
		if existing.Status == StatusCompleted {
			var res PaymentResult
			if err := json.Unmarshal([]byte(existing.ResponseJSON), &res); err != nil {
				return PaymentResult{}, 500, fmt.Errorf("unmarshal cached response: %w", err)
			}
			return res, existing.ResponseStatus, nil
		}

		// Still processing - do not execute again
		if existing.Status == StatusProcessing {
			return PaymentResult{}, 409, ErrRequestInProgress
		}
	}

	if !inserted {
		return PaymentResult{}, 409, ErrRequestInProgress
	}

	// Charge gateway
	gatewayResult, chargeErr := s.gateway.Charge(ctx, req)
	if chargeErr != nil {
		return PaymentResult{}, 500, fmt.Errorf("gateway charge failed: %w", chargeErr)
	}

	responseJSON, err := json.Marshal(gatewayResult)
	if err != nil {
		return PaymentResult{}, 500, fmt.Errorf("marshal response: %w", err)
	}

	// Update completed response
	if err := s.store.UpdateCompleted(ctx, idempotencyKey, gatewayResult, 200); err != nil {
		return PaymentResult{}, 500, fmt.Errorf("update completed: %w", err)
	}

	if err := s.store.UpdateCompletedResponse(ctx, idempotencyKey, string(responseJSON)); err != nil {
		return PaymentResult{}, 500, fmt.Errorf("store response: %w", err)
	}

	return gatewayResult, 200, nil
}

func (s *Service) CountPayments() int {
	s.store.(*dbStore).mu.Lock()
	defer s.store.(*dbStore).mu.Unlock()
	return len(s.store.(*dbStore).payments)
}
