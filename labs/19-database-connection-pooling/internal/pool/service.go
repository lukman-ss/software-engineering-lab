package pool

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrAcquireTimeout = errors.New("timeout acquiring connection from pool")
)

type OrderService struct {
	db *sql.DB
}

func NewOrderService(db *sql.DB) *OrderService {
	return &OrderService{db: db}
}

// ProcessOrderSafe performs external I/O outside of the DB transaction/connection.
func (s *OrderService) ProcessOrderSafe(ctx context.Context, orderID int, externalCall func() error) error {
	// 1. External Call first (or after), never while holding the connection
	if externalCall != nil {
		if err := externalCall(); err != nil {
			return err
		}
	}

	// 2. Short, bounded DB operation
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx, "UPDATE orders SET status = 'completed' WHERE id = ?", orderID)
	return err
}

// ProcessOrderUnsafeLeak holds the DB connection open during an external network call.
func (s *OrderService) ProcessOrderUnsafeLeak(ctx context.Context, orderID int, externalCall func() error) error {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx, "UPDATE orders SET status = 'processing' WHERE id = ?", orderID)
	if err != nil {
		return err
	}

	// External call while holding connection open!
	if externalCall != nil {
		if err := externalCall(); err != nil {
			return err
		}
	}

	return nil
}
