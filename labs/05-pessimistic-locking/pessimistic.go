package pessimistic

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type WalletPessimistic struct {
	db *sql.DB
}

func NewWalletPessimistic(db *sql.DB) *WalletPessimistic {
	return &WalletPessimistic{db: db}
}

// DepositPessimistic uses SELECT ... FOR UPDATE to lock the row during transaction.
func (w *WalletPessimistic) Deposit(ctx context.Context, id int, amount float64, holdDuration time.Duration) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var balance float64
	// SELECT ... FOR UPDATE blocks concurrent readers/writers from locking or modifying this row
	err = tx.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", id).Scan(&balance)
	if err != nil {
		return err
	}

	if holdDuration > 0 {
		time.Sleep(holdDuration) // Simulate processing / waiting / lock duration impact experiment
	}

	newBalance := balance + amount

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = $1 WHERE id = $2", newBalance, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
