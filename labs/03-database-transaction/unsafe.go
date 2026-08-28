package transaction

import (
	"context"
	"database/sql"
	"errors"
)

var ErrInjectedFailure = errors.New("injected failure after payment creation")

type PaymentServiceUnsafe struct {
	db *sql.DB
}

func NewPaymentServiceUnsafe(db *sql.DB) *PaymentServiceUnsafe {
	return &PaymentServiceUnsafe{db: db}
}

// ProcessPaymentUnsafe performs payment, updates order, and creates wallet tx
// WITHOUT a database transaction. If failure happens after payment, partial state is left.
func (s *PaymentServiceUnsafe) ProcessPayment(ctx context.Context, orderID int, amount float64, injectError bool) error {
	// 1. Create payment
	_, err := s.db.ExecContext(ctx, "INSERT INTO payments (order_id, amount, status) VALUES ($1, $2, 'completed')", orderID, amount)
	if err != nil {
		return err
	}

	// 2. Update order status
	_, err = s.db.ExecContext(ctx, "UPDATE orders SET status = 'paid' WHERE id = $1", orderID)
	if err != nil {
		return err
	}

	if injectError {
		return ErrInjectedFailure
	}

	// 3. Create wallet transaction
	_, err = s.db.ExecContext(ctx, "INSERT INTO wallet_transactions (order_id, amount, type) VALUES ($1, $2, 'credit')", orderID, amount)
	if err != nil {
		return err
	}

	return nil
}
