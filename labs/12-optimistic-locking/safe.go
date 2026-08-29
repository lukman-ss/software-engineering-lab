package optimistic

import (
	"context"
	"database/sql"
	"time"
)

type WalletSafe struct {
	db *sql.DB
}

func NewWalletSafe(db *sql.DB) *WalletSafe {
	return &WalletSafe{db: db}
}

// DepositWithOptimisticLocking updates balance ensuring version matches.
// Returns ErrOptimisticLockConflict if rows affected == 0.
func (w *WalletSafe) Deposit(ctx context.Context, id int, amount float64) error {
	for attempt := 0; attempt < 3; attempt++ {
		var balance float64
		var version int
		err := w.db.QueryRowContext(ctx, "SELECT balance, version FROM wallets WHERE id = $1", id).Scan(&balance, &version)
		if err != nil {
			return err
		}

		newBalance := balance + amount

		result, err := w.db.ExecContext(ctx, `
			UPDATE wallets
			SET balance = $1, version = version + 1
			WHERE id = $2 AND version = $3
		`, newBalance, id, version)
		if err != nil {
			return err
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rows > 0 {
			return nil // success
		}

		// Conflict detected, backoff and retry
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempt+1) * 10 * time.Millisecond):
		}
	}

	return ErrOptimisticLockConflict
}
