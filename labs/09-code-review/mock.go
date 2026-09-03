package codereview

import (
	"context"
	"sync"
)

type IdempotencyRecord struct {
	RequestHash string
	Processed   bool
	Response    *CheckoutResponse
}

type MockProductRepository struct {
	mu    sync.RWMutex
	stock map[string]int
}

func NewMockProductRepository(initialStock map[string]int) *MockProductRepository {
	return &MockProductRepository{stock: initialStock}
}

func (r *MockProductRepository) GetProduct(ctx context.Context, productID string) (*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stock, exists := r.stock[productID]
	if !exists {
		return nil, ErrProductNotFound
	}
	return &Product{ID: productID, Stock: stock, UnitPrice: 1000}, nil
}

func (r *MockProductRepository) UpdateStock(ctx context.Context, productID string, delta int) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	stock, exists := r.stock[productID]
	if !exists {
		return 0, ErrProductNotFound
	}
	newStock := stock + delta
	r.stock[productID] = newStock
	return newStock, nil
}

func (r *MockProductRepository) GetStock(ctx context.Context, productID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stock, exists := r.stock[productID]
	if !exists {
		return 0, ErrProductNotFound
	}
	return stock, nil
}

type MockOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		orders: make(map[string]*Order),
	}
}

func (r *MockOrderRepository) CreateOrder(ctx context.Context, order *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *MockOrderRepository) GetOrders(ctx context.Context) []*Order {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Order, 0, len(r.orders))
	for _, o := range r.orders {
		result = append(result, o)
	}
	return result
}

type MockIdempotencyRepository struct {
	mu      sync.RWMutex
	records map[string]*IdempotencyRecord
}

func NewMockIdempotencyRepository() *MockIdempotencyRepository {
	return &MockIdempotencyRepository{records: make(map[string]*IdempotencyRecord)}
}

func (r *MockIdempotencyRepository) TryInsert(ctx context.Context, key string, hash string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, exists := r.records[key]; exists {
		if existing.RequestHash != hash {
			return false, ErrIdempotencyConflict
		}
		return false, nil
	}
	r.records[key] = &IdempotencyRecord{RequestHash: hash, Processed: false}
	return true, nil
}

func (r *MockIdempotencyRepository) MarkProcessed(ctx context.Context, key string, resp *CheckoutResponse) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if rec, exists := r.records[key]; exists {
		rec.Processed = true
		rec.Response = resp
	}
	return nil
}

func (r *MockIdempotencyRepository) Get(ctx context.Context, key string) (*IdempotencyRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if rec, exists := r.records[key]; exists {
		return rec, nil
	}
	return nil, nil
}

type MockLogger struct {
	mu   sync.Mutex
	logs []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{logs: make([]string, 0)}
}

func (l *MockLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = append(l.logs, msg)
}

func (l *MockLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = append(l.logs, "ERROR: "+msg)
}

func (l *MockLogger) GetLogs() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	result := make([]string, len(l.logs))
	copy(result, l.logs)
	return result
}

type MockNotificationSender struct {
	mu   sync.Mutex
	sent []string
}

func NewMockNotificationSender() *MockNotificationSender {
	return &MockNotificationSender{sent: make([]string, 0)}
}

func (s *MockNotificationSender) SendOrderConfirmation(ctx context.Context, userID string, orderID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sent = append(s.sent, orderID)
	return nil
}

func (s *MockNotificationSender) GetSent() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]string, len(s.sent))
	copy(result, s.sent)
	return result
}