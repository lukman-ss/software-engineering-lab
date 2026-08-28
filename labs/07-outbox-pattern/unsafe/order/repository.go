// Package order provides the unsafe order repository implementation
// for the dual write problem demonstration.
package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Order represents an order entity in the database.
type Order struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// ErrOrderNotFound is returned when order is not found.
var ErrOrderNotFound = errors.New("order not found")

// UnsafeOrderRepository implements OrderRepository with plain SQL inserts.
// This is used to observe the dual write problem: when events fail to publish,
// orders are orphaned (exist in DB but no event was sent).
type UnsafeOrderRepository struct {
	db *sql.DB
}

// NewUnsafeOrderRepository creates a new unsafe repository.
func NewUnsafeOrderRepository(db *sql.DB) *UnsafeOrderRepository {
	return &UnsafeOrderRepository{db: db}
}

// Create inserts the order into the orders table.
func (r *UnsafeOrderRepository) Create(ctx context.Context, order Order) error {
	query := `
		INSERT INTO orders (id, customer_id, status, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query,
		order.ID,
		order.CustomerID,
		order.Status,
		order.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}
	return nil
}

// FindByID retrieves an order by ID.
func (r *UnsafeOrderRepository) FindByID(ctx context.Context, id string) (Order, error) {
	query := `
		SELECT id, customer_id, status, created_at
		FROM orders
		WHERE id = $1
	`
	var order Order
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.CustomerID,
		&order.Status,
		&order.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Order{}, ErrOrderNotFound
		}
		return Order{}, fmt.Errorf("failed to find order: %w", err)
	}
	return order, nil
}

// Count returns the total number of orders.
func (r *UnsafeOrderRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM orders`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count orders: %w", err)
	}
	return count, nil
}