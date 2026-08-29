package race

import (
	"context"
	"sync"
)

// MockInventoryRepository adalah implementasi mock untuk testing.
// Menggunakan map dengan mutex untuk single-threaded safety,
// tetapi TIDAK melindungi terhadap application-level race condition.
type MockInventoryRepository struct {
	mu     sync.RWMutex
	stocks map[string]int
}

// NewMockInventoryRepository membuat mock repository baru.
func NewMockInventoryRepository(initialStock map[string]int) *MockInventoryRepository {
	return &MockInventoryRepository{
		stocks: initialStock,
	}
}

// GetStock membaca stock secara concurrent-safe (Go memory model).
// Namun nilai yang dibaca mungkin sudah stale karena race condition lainnya.
func (m *MockInventoryRepository) GetStock(_ context.Context, productID string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stocks[productID], nil
}

// SetStock menulis stock secara concurrent-safe (Go memory model).
// Namun tidak ada proteksi terhadap lost update dari banyak pembaca sebelumnya.
func (m *MockInventoryRepository) SetStock(_ context.Context, productID string, stock int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stocks[productID] = stock
	return nil
}

// StockSnapshot merekam snapshot stock untuk verifikasi invariant.
func (m *MockInventoryRepository) StockSnapshot(productID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stocks[productID]
}
