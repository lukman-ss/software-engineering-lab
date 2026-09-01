// Package race explores application-level race conditions in concurrent systems.
package race

import (
	"context"
	"sync"
)

// AtomicInventory is a production-safe implementation using atomic conditional update.
//
// Pattern SQL (PostgreSQL):
//
//	UPDATE inventory_products
//	SET stock = stock - 1
//	WHERE id = $1 AND stock > 0
//	RETURNING stock;
//
// Interpretation:
//
//   - 1 row affected = decrement success
//   - 0 rows affected = out of stock / condition not met
//
// Benefits:
// - No SELECT before UPDATE (no stale read window)
// - Operation performed in a single atomic database statement
// - Concurrent execution handled by the database engine
//
// This is a recommended production pattern for inventory decrement.
type AtomicInventory struct {
	mu    sync.RWMutex
	stock int
}

// NewAtomicInventory membuat instance baru.
func NewAtomicInventory(initialStock int) *AtomicInventory {
	return &AtomicInventory{stock: initialStock}
}

// DecrementStock mengurangi stock secara atomik dengan kondisi stock > 0.
// Return: (newStock, true) jika sukses, (0, false) jika gagal.
//
// Ini analog dengan:
//
//	UPDATE products SET stock = stock - 1 WHERE id = $1 AND stock > 0 RETURNING stock
//	IF rows_affected == 1 THEN success ELSE rejected
func (a *AtomicInventory) DecrementStock(_ context.Context, _ string) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Cek condition: stock > 0
	if a.stock <= 0 {
		// 0 rows affected → rejected
		return 0, ErrOutOfStock
	}

	// Atomic decrement
	a.stock--

	// 1 row affected → success
	return a.stock, nil
}
