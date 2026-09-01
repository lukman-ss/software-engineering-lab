// Package race explores application-level race conditions in concurrent systems.
package race

import (
	"context"
	"sync"
)

// UnsafeInventoryRepository demonstrates an educational unsafe check-then-act pattern.
//
// WARNING: This implementation is **EDUCATIONAL** — it intentionally shows how
// a race condition can corrupt business state. **DO NOT USE** for production
// inventory mutation. Use AtomicInventory (atomic conditional update) or
// PostgresRowLockRepository (SELECT ... FOR UPDATE) in production.
//
// Problem: READ → CHECK → CALCULATE → WRITE is not atomic. Between READ and WRITE,
// another goroutine may change the state.
type UnsafeInventoryRepository struct {
	mu     sync.RWMutex
	stocks map[string]int
}

// NewUnsafeInventoryRepository membuat repository dengan kontrol sinkronisasi.
func NewUnsafeInventoryRepository(initialStock map[string]int) *UnsafeInventoryRepository {
	return &UnsafeInventoryRepository{
		stocks: initialStock,
	}
}

// GetStock membaca stock.
func (m *UnsafeInventoryRepository) GetStock(_ context.Context, productID string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stocks[productID], nil
}

// SetStock menulis stock.
func (m *UnsafeInventoryRepository) SetStock(_ context.Context, productID string, stock int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stocks[productID] = stock
	return nil
}

// StockSnapshot merekam snapshot stock untuk verifikasi invariant.
func (m *UnsafeInventoryRepository) StockSnapshot(productID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stocks[productID]
}
