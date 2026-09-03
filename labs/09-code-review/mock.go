package codereview

import (
	"context"
	"sync"
	"sync/atomic"
)

const (
	IdempotencyStatusProcessing = "PROCESSING"
	IdempotencyStatusCompleted  = "COMPLETED"
)

type IdempotencyRecord struct {
	RequestHash string
	Status      string
	Response    *CheckoutResponse
}

type MockProductRepository struct {
	mu            sync.RWMutex
	stock         map[string]int
	getCalls      atomic.Int64
	batchGetCalls atomic.Int64
	lastBatchIDs  []string
	readHook      func(productID string)
	getError      error
	batchGetError error
	reserveError  error
	updateError   error
}

func NewMockProductRepository(initialStock map[string]int) *MockProductRepository {
	stockCopy := make(map[string]int, len(initialStock))
	for k, v := range initialStock {
		stockCopy[k] = v
	}
	return &MockProductRepository{stock: stockCopy}
}

func (r *MockProductRepository) GetLastBatchIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]string, len(r.lastBatchIDs))
	copy(res, r.lastBatchIDs)
	return res
}

func (r *MockProductRepository) SetReadHook(hook func(productID string)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.readHook = hook
}

func (r *MockProductRepository) SetGetError(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.getError = err
}

func (r *MockProductRepository) SetBatchGetError(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.batchGetError = err
}

func (r *MockProductRepository) GetCalls() int64 {
	return r.getCalls.Load()
}

func (r *MockProductRepository) BatchGetCalls() int64 {
	return r.batchGetCalls.Load()
}

func (r *MockProductRepository) GetProduct(ctx context.Context, productID string) (*Product, error) {
	r.getCalls.Add(1)

	r.mu.RLock()
	hook := r.readHook
	err := r.getError
	stock, exists := r.stock[productID]
	r.mu.RUnlock()

	if hook != nil {
		hook(productID)
	}

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrProductNotFound
	}
	return &Product{ID: productID, Stock: stock, UnitPrice: 1000}, nil
}

func (r *MockProductRepository) GetProducts(ctx context.Context, productIDs []string) (map[string]*Product, error) {
	r.batchGetCalls.Add(1)

	r.mu.Lock()
	r.lastBatchIDs = make([]string, len(productIDs))
	copy(r.lastBatchIDs, productIDs)
	err := r.batchGetError
	r.mu.Unlock()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if err != nil {
		return nil, err
	}

	result := make(map[string]*Product, len(productIDs))
	for _, id := range productIDs {
		stock, exists := r.stock[id]
		if !exists {
			return nil, ErrProductNotFound
		}
		result[id] = &Product{ID: id, Stock: stock, UnitPrice: 1000}
	}
	return result, nil
}

func (r *MockProductRepository) UpdateStock(ctx context.Context, productID string, delta int) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.updateError != nil {
		return 0, r.updateError
	}

	stock, exists := r.stock[productID]
	if !exists {
		return 0, ErrProductNotFound
	}
	newStock := stock + delta
	r.stock[productID] = newStock
	return newStock, nil
}

func (r *MockProductRepository) ReserveStock(ctx context.Context, productID string, quantity int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.reserveError != nil {
		return r.reserveError
	}

	stock, exists := r.stock[productID]
	if !exists {
		return ErrProductNotFound
	}
	if stock < quantity {
		return ErrInsufficientStock
	}
	r.stock[productID] = stock - quantity
	return nil
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
	mu      sync.RWMutex
	orders  map[string]*Order
	failErr error
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		orders: make(map[string]*Order),
	}
}

func (r *MockOrderRepository) FailNextCreate(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.failErr = err
}

func (r *MockOrderRepository) consumeCreateError() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.failErr != nil {
		err := r.failErr
		r.failErr = nil
		return err
	}
	return nil
}

func (r *MockOrderRepository) CreateOrder(ctx context.Context, order *Order) error {
	if err := r.consumeCreateError(); err != nil {
		return err
	}
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

type MockTransactionManager struct {
	mu       sync.Mutex
	products *MockProductRepository
	orders   *MockOrderRepository
}

func NewMockTransactionManager(products *MockProductRepository, orders *MockOrderRepository) *MockTransactionManager {
	return &MockTransactionManager{
		products: products,
		orders:   orders,
	}
}

type mockTx struct {
	parent   *MockTransactionManager
	stockMut map[string]int
	orderMut *Order
}

func (tm *MockTransactionManager) WithinTransaction(ctx context.Context, fn func(tx CheckoutTx) error) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tx := &mockTx{
		parent:   tm,
		stockMut: make(map[string]int),
	}

	if err := fn(tx); err != nil {
		return err
	}

	return tx.commit(ctx)
}

func (tx *mockTx) ReserveStock(ctx context.Context, productID string, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	tx.parent.products.mu.RLock()
	currentStock, exists := tx.parent.products.stock[productID]
	tx.parent.products.mu.RUnlock()

	if !exists {
		return ErrProductNotFound
	}

	accumulated := tx.stockMut[productID]
	if currentStock-accumulated < quantity {
		return ErrInsufficientStock
	}

	tx.stockMut[productID] = accumulated + quantity
	return nil
}

func (tx *mockTx) CreateOrder(ctx context.Context, order *Order) error {
	if err := tx.parent.orders.consumeCreateError(); err != nil {
		return err
	}
	tx.orderMut = order
	return nil
}

func (tx *mockTx) commit(ctx context.Context) error {
	tx.parent.products.mu.Lock()
	defer tx.parent.products.mu.Unlock()

	for pID, deduct := range tx.stockMut {
		curr, exists := tx.parent.products.stock[pID]
		if !exists || curr < deduct {
			return ErrInsufficientStock
		}
	}

	for pID, deduct := range tx.stockMut {
		tx.parent.products.stock[pID] -= deduct
	}

	if tx.orderMut != nil {
		if err := tx.parent.orders.CreateOrder(ctx, tx.orderMut); err != nil {
			for pID, deduct := range tx.stockMut {
				tx.parent.products.stock[pID] += deduct
			}
			return err
		}
	}

	return nil
}

type MockIdempotencyRepository struct {
	mu           sync.RWMutex
	records      map[string]*IdempotencyRecord
	claimErr     error
	markErr      error
	releaseError error
}

func NewMockIdempotencyRepository() *MockIdempotencyRepository {
	return &MockIdempotencyRepository{records: make(map[string]*IdempotencyRecord)}
}

func (r *MockIdempotencyRepository) FailNextClaim(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.claimErr = err
}

func (r *MockIdempotencyRepository) FailNextMarkCompleted(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.markErr = err
}

func (r *MockIdempotencyRepository) FailNextRelease(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.releaseError = err
}

func (r *MockIdempotencyRepository) consumeClaimError() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.claimErr != nil {
		err := r.claimErr
		r.claimErr = nil
		return err
	}
	return nil
}

func (r *MockIdempotencyRepository) consumeMarkError() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.markErr != nil {
		err := r.markErr
		r.markErr = nil
		return err
	}
	return nil
}

func (r *MockIdempotencyRepository) consumeReleaseError() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.releaseError != nil {
		err := r.releaseError
		r.releaseError = nil
		return err
	}
	return nil
}

func (r *MockIdempotencyRepository) Claim(ctx context.Context, key string, hash string) (string, *CheckoutResponse, error) {
	if err := r.consumeClaimError(); err != nil {
		return "", nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.records[key]
	if exists {
		if existing.RequestHash != hash {
			return "", nil, ErrIdempotencyConflict
		}
		if existing.Status == IdempotencyStatusCompleted && existing.Response != nil {
			return IdempotencyStatusCompleted, existing.Response, nil
		}
		return IdempotencyStatusProcessing, nil, ErrDuplicateRequest
	}

	r.records[key] = &IdempotencyRecord{
		RequestHash: hash,
		Status:      IdempotencyStatusProcessing,
	}

	return IdempotencyStatusProcessing, nil, nil
}

func (r *MockIdempotencyRepository) MarkCompleted(ctx context.Context, key string, resp *CheckoutResponse) error {
	if err := r.consumeMarkError(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if rec, exists := r.records[key]; exists {
		rec.Status = IdempotencyStatusCompleted
		rec.Response = resp
	}
	return nil
}

func (r *MockIdempotencyRepository) Release(ctx context.Context, key string) error {
	if err := r.consumeReleaseError(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, key)
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
