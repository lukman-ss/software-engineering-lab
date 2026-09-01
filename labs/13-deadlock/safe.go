package deadlock

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

var ErrDeadlock = errors.New("deadlock detected")

// IsDeadlockError checks if the error is a Postgres deadlock (40P01) or serialization failure (40001).
func IsDeadlockError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == pgerrcode.DeadlockDetected || pqErr.Code == pgerrcode.SerializationFailure
	}
	return false
}

type DeadlockSafe struct {
	db *sql.DB
}

func NewDeadlockSafe(db *sql.DB) *DeadlockSafe {
	return &DeadlockSafe{db: db}
}

// TransferSafe uses deterministic lock ordering: always lock the smaller ID first.
func (t *DeadlockSafe) Transfer(ctx context.Context, fromID, toID int, amount float64) error {
	if fromID == toID {
		return errors.New("cannot transfer to self")
	}

	firstID, secondID := fromID, toID
	if fromID > toID {
		firstID, secondID = toID, fromID
	}

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var firstBal, secondBal float64
	if err := tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", firstID).Scan(&firstBal); err != nil {
		return err
	}
	if err := tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", secondID).Scan(&secondBal); err != nil {
		return err
	}

	// Apply delta based on actual transfer direction
	if firstID == fromID {
		firstBal -= amount
		secondBal += amount
	} else {
		firstBal += amount
		secondBal -= amount
	}

	if firstBal < 0 || secondBal < 0 {
		return errors.New("insufficient funds")
	}

	if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = $1 WHERE id = $2", firstBal, firstID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = $1 WHERE id = $2", secondBal, secondID); err != nil {
		return err
	}

	return tx.Commit()
}

// TransferWithRetry adds bounded retries only for transient deadlock/serialization errors.
func (t *DeadlockSafe) TransferWithRetry(ctx context.Context, fromID, toID int, amount float64) error {
	const maxRetries = 3

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := t.Transfer(ctx, fromID, toID, amount)
		if err == nil {
			return nil
		}
		// Retry ONLY on transient deadlock/serialization errors
		if !IsDeadlockError(err) {
			return err
		}
		lastErr = err

		// Bounded backoff with jitter, but respect context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempt+1) * 20 * time.Millisecond):
		}
	}

	if lastErr != nil {
		return fmt.Errorf("%w: %v", ErrDeadlock, lastErr)
	}
	return nil
}
