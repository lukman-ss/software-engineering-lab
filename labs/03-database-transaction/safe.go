package transaction

import (
	"context"
	"database/sql"
	"fmt"
)

// WithTx runs fn inside a database transaction, handling commit, rollback on error,
// and recovery on panic. Context is passed through.
func WithTx(ctx context.Context, db *sql.DB, fn func(ctx context.Context, tx *sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: %v (original error: %w)", rbErr, err)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("commit tx: %w", commitErr)
			}
		}
	}()

	err = fn(ctx, tx)
	return err
}

type PaymentServiceSafe struct {
	db *sql.DB
}

func NewPaymentServiceSafe(db *sql.DB) *PaymentServiceSafe {
	return &PaymentServiceSafe{db: db}
}

// ProcessPaymentSafe uses WithTx to ensure atomicity across payment, order, and wallet tx.
func (s *PaymentServiceSafe) ProcessPayment(ctx context.Context, orderID int, amount float64, injectError bool) error {
	return WithTx(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		// 1. Create payment
		_, err := tx.ExecContext(ctx, "INSERT INTO payments (order_id, amount, status) VALUES ($1, $2, 'completed')", orderID, amount)
		if err != nil {
			return err
		}

		// 2. Update order status
		_, err = tx.ExecContext(ctx, "UPDATE orders SET status = 'paid' WHERE id = $1", orderID)
		if err != nil {
			return err
		}

		if injectError {
			return ErrInjectedFailure
		}

		// 3. Create wallet transaction
		_, err = tx.ExecContext(ctx, "INSERT INTO wallet_transactions (order_id, amount, type) VALUES ($1, $2, 'credit')", orderID, amount)
		if err != nil {
			return err
		}

		return nil
	})
}
