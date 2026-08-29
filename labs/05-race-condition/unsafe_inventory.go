// Package race explores application-level race conditions in concurrent systems.
package race

import (
	"context"
	"sync"
)

// UnsafeInventoryRepository menunjukkan check-then-act pattern yang TIDAK thread-safe
// untuk business state, meskipun masing-masing operasi read/write-nya aman secara Go memory model.
//
// Problem: READ → CHECK → CALCULATE → WRITE adalah sebuah transaksi yang tidak atoms.
// Antara READ dan WRITE, goroutine lain dapat mengubah state.
type UnsafeInventoryRepository struct {
	mu          sync.RWMutex // Hanya melindungi map access, bukan business logic
	stocks      map[string]int
	readyToRead chan string // Channel untuk sinkronisasi testing
}

// NewUnsafeInventoryRepository membuat repository dengan kontrol sinkronisasi.
func NewUnsafeInventoryRepository(initialStock map[string]int, readyToRead chan string) *UnsafeInventoryRepository {
	return &UnsafeInventoryRepository{
		stocks:      initialStock,
		readyToRead: readyToRead,
	}
}

// GetStock membaca stock.
// Mengirimkan productID ke channel untuk barrier sinkronisasi.
func (m *UnsafeInventoryRepository) GetStock(_ context.Context, productID string) (int, error) {
	// Signal bahwa request siap membaca
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Non-blocking send untuk barrier - hanya yang terblokir dalam test
	select {
	case m.readyToRead <- productID:
	default:
	}

	return m.stocks[productID], nil
}

// SetStock menulis stock.
func (m *UnsafeInventoryRepository) SetStock(_ context.Context, productID string, stock int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stocks[productID] = stock
	return nil
}
