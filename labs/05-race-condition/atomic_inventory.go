// Package race explores application-level race conditions in concurrent systems.
package race

import (
	"context"
	"sync"
)

// AtomicInventory adalah implementasi safe menggunakan PostgreSQL-style atomic conditional update.
//
// Pola SQL:
//
//	UPDATE products SET stock = stock - 1 WHERE id = $1 AND stock > 0 RETURNING stock
//
// Interpretasi:
//
//	1 row affected  = decrement sukses
//	0 rows affected = stock habis / condition tidak terpenuhi
//
// Keuntungan:
// - Tidak ada SELECT terlebih dahulu (tidak ada stale read window)
// - Operasi dilakukan dalam satu statement atomik di database
// - Concurrent execution di-handle oleh database engine
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
