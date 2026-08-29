package optimistic

import (
	"context"
	"database/sql"
	"errors"
)

var ErrOptimisticLockConflict = errors.New("optimistic lock conflict: version mismatch")

type WalletUnsafe struct {
	db *sql.DB
}

func NewWalletUnsafe(db *sql.DB) *WalletUnsafe {
	return &WalletUnsafe{db: db}
}

// DepositUnsafe reads balance and version, computes new balance, and updates
// without checking version. Causes lost updates under concurrency.
func (w *WalletUnsafe) Deposit(ctx context.Context, id int, amount float64) error {
	var balance float64
	var version int
	err := w.db.QueryRowContext(ctx, "SELECT balance, version FROM wallets WHERE id = $1", id).Scan(&balance, &version)
	if err != nil {
		return err
	}

	newBalance := balance + amount

	_, err = w.db.ExecContext(ctx, "UPDATE wallets SET balance = $1, version = version + 1 WHERE id = $2", newBalance, id)
	return err
}
