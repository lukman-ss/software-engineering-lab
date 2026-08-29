package transaction

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

var ErrInjectedFailure = errors.New("injected failure after payment creation")

// PaymentServiceUnsafe demonstrates partial state corruption without transaction.
type PaymentServiceUnsafe struct{ db *sql.DB }

func NewPaymentServiceUnsafe(db *sql.DB) *PaymentServiceUnsafe {
	return &PaymentServiceUnsafe{db: db}
}

func (s *PaymentServiceUnsafe) ProcessPayment(ctx context.Context, orderID int, amount float64, injectError bool) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO payments (order_id, amount, status) VALUES ($1, $2, 'completed')", orderID, amount)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, "UPDATE orders SET status = 'paid' WHERE id = $1", orderID)
	if err != nil {
		return err
	}

	if injectError {
		return ErrInjectedFailure
	}

	_, err = s.db.ExecContext(ctx, "INSERT INTO wallet_transactions (order_id, amount, type) VALUES ($1, $2, 'credit')", orderID, amount)
	return err
}

// PaymentServiceSafe demonstrates clean ACID rollback with local transaction.
type PaymentServiceSafe struct{ db *sql.DB }

func NewPaymentServiceSafe(db *sql.DB) *PaymentServiceSafe {
	return &PaymentServiceSafe{db: db}
}

func (s *PaymentServiceSafe) ProcessPayment(ctx context.Context, orderID int, amount float64, injectError bool) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, "INSERT INTO payments (order_id, amount, status) VALUES ($1, $2, 'completed')", orderID, amount)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE orders SET status = 'paid' WHERE id = $1", orderID)
	if err != nil {
		return err
	}

	if injectError {
		return ErrInjectedFailure
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO wallet_transactions (order_id, amount, type) VALUES ($1, $2, 'credit')", orderID, amount)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("[PAYMENT] transaction committed")
	return nil
}
