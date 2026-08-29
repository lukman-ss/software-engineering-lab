// Package race explores application-level race conditions in concurrent systems.
package race

import (
	"context"
	"sync"
)

// UnsafeInventoryRepository menunjukkan check-then-act pattern yang TIDAK safe
// untuk business state, meskipun masing-masing operasi read/write-nya aman
// secara Go memory model.
//
// Problem: READ → CHECK → CALCULATE → WRITE adalah sebuah transaksi yang
// tidak atomic. Antara READ dan WRITE, goroutine lain dapat mengubah state.
//
// Ponytail: ini adalah demo unsafe — jangan pakai di production.
// Gunakan AtomicInventoryRepository untuk production.
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
