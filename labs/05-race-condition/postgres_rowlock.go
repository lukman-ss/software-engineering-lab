//go:build integration
// +build integration

package race

import (
	"context"
	"database/sql"
	"fmt"
)

// PostgresRowLockRepository is a production-safe implementation using pessimistic row locking.
//
// Pattern: SELECT ... FOR UPDATE inside a transaction takes a row-level lock.
// Other transactions attempting a conflicting lock on the same row will block
// until the holding transaction commits or rolls back.
//
// This is a recommended production pattern when business logic requires reading
// state before deciding the new value (complex read-modify-write) where an
// atomic update is not sufficient.
type PostgresRowLockRepository struct {
	db *sql.DB
}

// NewPostgresRowLockRepository creates a new repository with row locking.
func NewPostgresRowLockRepository(db *sql.DB) *PostgresRowLockRepository {
	return &PostgresRowLockRepository{db: db}
}

// TrySell acquires a row lock, checks stock, then decrements.
//
// Transaction A acquires the row lock via SELECT ... FOR UPDATE.
// Transaction B blocks on SELECT ... FOR UPDATE (same row) until A commits.
// After A commits (stock = 0), B reads stock = 0 → CHECK fails → Reject.
//
// Note: Non-locking SELECTs (without FOR UPDATE) using MVCC snapshots are
// still allowed to read the row while A holds the lock; only conflicting
// writers/lockers are blocked.
func (r *PostgresRowLockRepository) TrySell(ctx context.Context, productID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Lock the row — blocks other FOR UPDATE transactions
	var stock int
	err = tx.QueryRowContext(ctx,
		"SELECT stock FROM inventory_products WHERE id = $1 FOR UPDATE",
		productID).Scan(&stock)
	if err != nil {
		return fmt.Errorf("lock and read: %w", err)
	}

	if stock <= 0 {
		return ErrOutOfStock
	}

	// Decrement inside the same transaction — lock prevents concurrent modification
	_, err = tx.ExecContext(ctx,
		"UPDATE inventory_products SET stock = $1 WHERE id = $2",
		stock-1, productID)
	if err != nil {
		return fmt.Errorf("update stock: %w", err)
	}

	return tx.Commit()
}

// GetStock reads the current stock value for verification.
func (r *PostgresRowLockRepository) GetStock(ctx context.Context, productID string) (int, error) {
	var stock int
	err := r.db.QueryRowContext(ctx,
		"SELECT stock FROM inventory_products WHERE id = $1",
		productID).Scan(&stock)
	if err != nil {
		return 0, fmt.Errorf("get stock: %w", err)
	}
	return stock, nil
}
