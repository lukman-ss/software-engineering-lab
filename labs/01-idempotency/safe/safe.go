// Package safe implements idempotent payment processing.
// It demonstrates: idempotency keys, payload hashing, unique constraints,
// database transactions, concurrent request handling, response caching,
// and conflict detection.
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
	// ErrRequestInProgress is returned when a concurrent request with the same
	// idempotency key is currently being processed.
	ErrRequestInProgress = errors.New("request is already in progress")
	// ErrLeaseExpired is returned when a processing record's lease has expired.
	ErrLeaseExpired = errors.New("processing lease expired")
	// ErrDuplicateOrder is returned when order has already been paid.
	ErrDuplicateOrder = errors.New("order already paid")
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
	StatusFailed     IdempotencyStatus = "failed"
)

// DefaultTTL is the default time-to-live for idempotency keys after completion.
const DefaultTTL = 72 * time.Hour // 3 days covers most client retry windows

// DefaultLeaseDuration is how long a processing lease is held before another worker can take over.
const DefaultLeaseDuration = 30 * time.Second

// TTLConfig holds configurable TTL values for different entity types.
type TTLConfig struct {
	// Payment is the TTL for completed payment idempotency keys.
	// For financial records, prioritize auditability over cleanup.
	// Default: 72h (covers client retry, allows post-mortem analysis)
	Payment time.Duration

	// InvoiceCreation is the TTL for invoice creation idempotency.
	// Invoices may be reviewed before finalization.
	// Default: 168h (7 days - allows multi-day review)
	InvoiceCreation time.Duration

	// FileUpload is the TTL for file upload idempotency.
	// Smaller payloads, faster cleanup OK.
	// Default: 24h (user typically waits for upload confirmation)
	FileUpload time.Duration

	// GenericCommand is the TTL for non-financial commands.
	// Ephemeral operations.
	// Default: 1h (quick operations)
	GenericCommand time.Duration
}

// DefaultTTLConfig returns sensible defaults for each entity type.
func DefaultTTLConfig() TTLConfig {
	return TTLConfig{
		Payment:         DefaultTTL,
		InvoiceCreation: 7 * 24 * time.Hour,
		FileUpload:      24 * time.Hour,
		GenericCommand:  1 * time.Hour,
	}
}

// Config holds the idempotency configuration.
type Config struct {
	TTL       TTLConfig
	LeaseDuration time.Duration // Override for processing leases
}

// PaymentRecord is the full database record including idempotency metadata.
type PaymentRecord struct {
	ID             string            `json:"id"`
	IdempotencyKey string            `json:"idempotency_key"`
	RequestHash    string            `json:"request_hash"`
	OrderID        string            `json:"order_id"` // For business constraint uniqueness
	Status         IdempotencyStatus `json:"status"`
	LockedAt       time.Time         `json:"locked_at"` // When 'processing' lease started
	ResponseStatus int               `json:"response_status"`
	ResponseJSON   string            `json:"response_json"`
	Amount         int64             `json:"amount"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	ExpiresAt      time.Time         `json:"expires_at"` // Idempotency key TTL
}

// Gateway simulates an external payment gateway.
type Gateway interface {
	Charge(ctx context.Context, req PaymentRequest) (PaymentResult, error)
}

// mockGateway simulates a payment gateway with call tracking.
type mockGateway struct {
	mu          sync.Mutex
	chargeCount int64
	latency     time.Duration
}

func (g *mockGateway) Charge(ctx context.Context, req PaymentRequest) (PaymentResult, error) {
	g.mu.Lock()
	g.chargeCount++
	count := g.chargeCount
	g.mu.Unlock()

	if g.latency > 0 {
		time.Sleep(g.latency)
	}

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

// MockGatewayWithLatency creates a mock gateway with simulated processing delay.
func MockGatewayWithLatency(d time.Duration) Gateway {
	return &mockGateway{latency: d}
}

// GatewayChargeCount returns the charge count for a mock gateway.
func GatewayChargeCount(g Gateway) int64 {
	if mg, ok := g.(*mockGateway); ok {
		return mg.ChargeCount()
	}
	return 0
}

// Store is the persistence layer interface for idempotency and payments.
type Store interface {
	TryAcquire(ctx context.Context, record PaymentRecord) (bool, error)
	TryAcquireLease(ctx context.Context, idempotencyKey, requestHash string, orderID string) (bool, error)
	UpdateCompleted(ctx context.Context, key string, result PaymentResult, responseStatus int) error
	UpdateCompletedResponse(ctx context.Context, key string, responseJSON string) error
	UpdateFailed(ctx context.Context, key string) error
	GetByIdempotencyKey(ctx context.Context, key string) (*PaymentRecord, error)
	IsOrderPaid(ctx context.Context, orderID string) (string, error)
	UpsertOrder(ctx context.Context, orderID, paymentID string) error

	// Cleanup removes expired idempotency keys. Returns count of removed records.
	Cleanup(ctx context.Context) (int, error)

	// ExpiredKeysCount returns count of keys past their expires_at time.
	ExpiredKeysCount(ctx context.Context) int
}

type dbStore struct {
	mu       sync.RWMutex
	payments map[string]PaymentRecord
	orders   map[string]string // order_id -> payment_id (business constraint)
}

func newDBStore() *dbStore {
	return &dbStore{
		payments: make(map[string]PaymentRecord),
		orders:   make(map[string]string),
	}
}

func (s *dbStore) TryAcquire(ctx context.Context, record PaymentRecord) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.payments[record.IdempotencyKey]
	if exists {
		if existing.RequestHash != record.RequestHash {
			return false, ErrIdempotencyKeyConflict
		}
		return false, nil
	}

	s.payments[record.IdempotencyKey] = record
	return true, nil
}

func (s *dbStore) TryAcquireLease(ctx context.Context, idempotencyKey, requestHash, orderID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.payments[idempotencyKey]
	if !exists {
		return false, nil
	}

	if existing.RequestHash != requestHash {
		return false, ErrIdempotencyKeyConflict
	}

	// Lease takeover: check if processing lease has expired
	if existing.Status == StatusProcessing {
		if time.Since(existing.LockedAt) > LeaseDuration {
			existing.LockedAt = time.Now()
			existing.Status = StatusProcessing
			existing.OrderID = orderID
			s.payments[idempotencyKey] = existing
			return true, nil
		}
	}

	return false, nil
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
	record.ID = result.PaymentID
	record.UpdatedAt = time.Now()
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
	record.UpdatedAt = time.Now()
	s.payments[key] = record
	return nil
}

func (s *dbStore) UpdateFailed(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, exists := s.payments[key]
	if !exists {
		return fmt.Errorf("record not found for key: %s", key)
	}

	record.Status = StatusFailed
	record.UpdatedAt = time.Now()
	s.payments[key] = record
	return nil
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

func (s *dbStore) IsOrderPaid(ctx context.Context, orderID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	paymentID, exists := s.orders[orderID]
	if !exists {
		return "", nil
	}
	return paymentID, nil
}

func (s *dbStore) UpsertOrder(ctx context.Context, orderID, paymentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existingPaymentID, exists := s.orders[orderID]; exists {
		if existingPaymentID != paymentID {
			return ErrDuplicateOrder
		}
		return nil
	}
	s.orders[orderID] = paymentID
	return nil
}

func (s *dbStore) ExpiredKeysCount(ctx context.Context) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	count := 0
	for _, record := range s.payments {
		if record.ExpiresAt.Before(now) {
			count++
		}
	}
	return count
}

func (s *dbStore) Cleanup(ctx context.Context) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	removed := 0
	for key, record := range s.payments {
		if record.ExpiresAt.Before(now) {
			delete(s.payments, key)
			removed++
		}
	}
	return removed, nil
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

// IsIdempotentKeyConflict checks if an error is an idempotency key conflict.
func IsIdempotentKeyConflict(err error) bool {
	return errors.Is(err, ErrIdempotencyKeyConflict)
}

// Service processes payments with full idempotency guarantees, failure recovery, and business constraints.
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
	hash := hashRequest(method, path, req)

	record := PaymentRecord{
		IdempotencyKey: idempotencyKey,
		RequestHash:    hash,
		OrderID:        req.OrderID,
		Status:         StatusProcessing,
		LockedAt:       time.Now(),
		Amount:         req.Amount,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	acquired, err := s.store.TryAcquire(ctx, record)
	if err != nil {
		if errors.Is(err, ErrIdempotencyKeyConflict) {
			return PaymentResult{}, 409, ErrIdempotencyKeyConflict
		}
		return PaymentResult{}, 500, fmt.Errorf("try acquire key: %w", err)
	}

	if !acquired {
		existing, err := s.store.GetByIdempotencyKey(ctx, idempotencyKey)
		if err != nil {
			return PaymentResult{}, 500, fmt.Errorf("get existing key: %w", err)
		}

		if existing.RequestHash != hash {
			return PaymentResult{}, 409, ErrIdempotencyKeyConflict
		}

		switch existing.Status {
		case StatusCompleted:
			var res PaymentResult
			if err := json.Unmarshal([]byte(existing.ResponseJSON), &res); err != nil {
				return PaymentResult{}, 500, fmt.Errorf("unmarshal cached response: %w", err)
			}
			return res, existing.ResponseStatus, nil

		case StatusProcessing:
			takenOver, takeOverErr := s.store.TryAcquireLease(ctx, idempotencyKey, hash, req.OrderID)
			if takeOverErr != nil {
				return PaymentResult{}, 500, fmt.Errorf("lease takeover check: %w", takeOverErr)
			}
			if takenOver {
				// Taken over expired lease, proceed to process
			} else {
				deadline := time.Now().Add(5 * time.Second)
				for time.Now().Before(deadline) {
					time.Sleep(50 * time.Millisecond)
					current, err := s.store.GetByIdempotencyKey(ctx, idempotencyKey)
					if err != nil {
						continue
					}
					if current.Status == StatusCompleted {
						var res PaymentResult
						if err := json.Unmarshal([]byte(current.ResponseJSON), &res); err != nil {
							return PaymentResult{}, 500, fmt.Errorf("unmarshal cached: %w", err)
						}
						return res, current.ResponseStatus, nil
					}
					if current.Status == StatusFailed {
						return PaymentResult{}, 500, fmt.Errorf("previous concurrent attempt failed")
					}
				}
				return PaymentResult{}, 409, fmt.Errorf("request timeout waiting for in-flight processing")
			}

		case StatusFailed:
			_, resetErr := s.store.TryAcquire(ctx, PaymentRecord{
				IdempotencyKey: idempotencyKey,
				RequestHash:    hash,
				Status:         StatusProcessing,
				LockedAt:       time.Now(),
			})
			if resetErr != nil {
				return PaymentResult{}, 500, fmt.Errorf("reset for retry: %w", resetErr)
			}
		}
	}

	// 2. Business constraint check (unique order_id)
	if req.OrderID != "" {
		paidPaymentID, err := s.store.IsOrderPaid(ctx, req.OrderID)
		if err != nil {
			return PaymentResult{}, 500, fmt.Errorf("order check failed: %w", err)
		}
		if paidPaymentID != "" {
			_ = s.store.UpdateFailed(ctx, idempotencyKey)
			return PaymentResult{}, 409, ErrDuplicateOrder
		}
	}

	// 3. Charge gateway
	gatewayResult, chargeErr := s.gateway.Charge(ctx, req)
	if chargeErr != nil {
		_ = s.store.UpdateFailed(ctx, idempotencyKey)
		return PaymentResult{}, 500, fmt.Errorf("gateway charge failed: %w", chargeErr)
	}

	responseJSON, err := json.Marshal(gatewayResult)
	if err != nil {
		_ = s.store.UpdateFailed(ctx, idempotencyKey)
		return PaymentResult{}, 500, fmt.Errorf("marshal response: %w", err)
	}

	// 4. Update completed response & transaction boundary
	if err := s.store.UpdateCompleted(ctx, idempotencyKey, gatewayResult, 200); err != nil {
		return PaymentResult{}, 500, fmt.Errorf("update completed: %w", err)
	}

	if err := s.store.UpdateCompletedResponse(ctx, idempotencyKey, string(responseJSON)); err != nil {
		return PaymentResult{}, 500, fmt.Errorf("store response: %w", err)
	}

	// 5. Enforce business uniqueness (unique order_id)
	if req.OrderID != "" {
		if err := s.store.UpsertOrder(ctx, req.OrderID, gatewayResult.PaymentID); err != nil {
			if errors.Is(err, ErrDuplicateOrder) {
				return PaymentResult{}, 409, ErrDuplicateOrder
			}
			return PaymentResult{}, 500, fmt.Errorf("enforce order uniqueness: %w", err)
		}
	}

	return gatewayResult, 200, nil
}

func (s *Service) CountPayments() int {
	s.store.(*dbStore).mu.RLock()
	defer s.store.(*dbStore).mu.RUnlock()
	return len(s.store.(*dbStore).payments)
}

func (s *Service) CountOrders() int {
	s.store.(*dbStore).mu.RLock()
	defer s.store.(*dbStore).mu.RUnlock()
	return len(s.store.(*dbStore).orders)
}

// IsOrderPaidError checks if error is due to duplicate order payment.
func IsOrderPaidError(err error) bool {
	return errors.Is(err, ErrDuplicateOrder)
}
