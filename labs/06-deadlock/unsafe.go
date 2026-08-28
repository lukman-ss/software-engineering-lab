package deadlock

import (
	"context"
	"database/sql"
)

type DeadlockUnsafe struct {
	db *sql.DB
}

func NewDeadlockUnsafe(db *sql.DB) *DeadlockUnsafe {
	return &DeadlockUnsafe{db: db}
}

// TransferUnsafe transfers funds by locking accounts in caller-specified order.
// If different callers lock A→B and B→A concurrently, deadlock can occur.
func (t *DeadlockUnsafe) Transfer(ctx context.Context, fromID, toID int, amount float64) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lock accounts in ambiguous order (caller-dependent)
	var fromBal, toBal float64
	if err := tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", fromID).Scan(&fromBal); err != nil {
		return err
	}

	if err := tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", toID).Scan(&toBal); err != nil {
		return err
	}

	fromBal -= amount
	toBal += amount

	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = $1 WHERE id = $2", fromBal, fromID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = $1 WHERE id = $2", toBal, toID)
	if err != nil {
		return err
	}

	return tx.Commit()
}