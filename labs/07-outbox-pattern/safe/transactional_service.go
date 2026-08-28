// Package safe implements the transactional outbox pattern for safe event publishing.
package safe

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Order represents an order entity.
type Order struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// OrderCreatedEvent is the event published when order is created.
type OrderCreatedEvent struct {
	OrderID    string    `json:"order_id"`
	CustomerID string    `json:"customer_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// OutboxEvent represents an entry in the outbox table.
type OutboxEvent struct {
	ID            string         `json:"id"`
	AggregateType string         `json:"aggregate_type"`
	AggregateID   string         `json:"aggregate_id"`
	EventType     string         `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
	CreatedAt     time.Time      `json:"created_at"`
	PublishedAt   *time.Time     `json:"published_at"`
	Attempts      int            `json:"attempts"`
	NextAttemptAt time.Time      `json:"next_attempt_at"`
}

// SafeOrderService implements the transactional outbox pattern.
// The key insight: both order insertion AND outbox event insertion
// happen in the SAME database transaction, ensuring atomicity.
type SafeOrderService struct {
	db *sql.DB
}

// NewSafeOrderService creates a new service with transactional outbox.
func NewSafeOrderService(db *sql.DB) *SafeOrderService {
	return &SafeOrderService{db: db}
}

// CreateOrder creates an order with a corresponding outbox event
// in a single database transaction.
func (s *SafeOrderService) CreateOrder(ctx context.Context, order Order) (Order, error) {
	// Generate unique IDs
	orderID := uuid.New().String()
	eventID := uuid.New().String()

	// Convert order to JSON for event payload
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return Order{}, fmt.Errorf("failed to marshal order: %w", err)
	}

	// Execute both operations in the same transaction
	// This is the key: atomic operation that cannot fail partially
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Order{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Insert order
	_, err = tx.ExecContext(ctx,
		`INSERT INTO orders (id, customer_id, status, created_at)
		 VALUES ($1, $2, $3, $4)`,
		orderID, order.CustomerID, order.Status, time.Now(),
	)
	if err != nil {
		_ = tx.Rollback()
		return Order{}, fmt.Errorf("failed to insert order: %w", err)
	}

	// Insert outbox event (same transaction!)
	_, err = tx.ExecContext(ctx,
		`INSERT INTO outbox_events
		 (id, aggregate_type, aggregate_id, event_type, payload, created_at, published_at, attempts, next_attempt_at)
		 VALUES ($1, $2, $3, $4, $5::jsonb, $6, NULL, 0, $7)`,
		eventID,
		"Order",
		orderID,
		"OrderCreated",
		orderJSON,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		_ = tx.Rollback()
		return Order{}, fmt.Errorf("failed to insert outbox event: %w", err)
	}

	// Commit both changes atomically
	if err := tx.Commit(); err != nil {
		return Order{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	order.ID = orderID
	return order, nil
}

// FindOrder retrieves an order by ID.
func (s *SafeOrderService) FindOrder(ctx context.Context, id string) (Order, error) {
	query := `
		SELECT id, customer_id, status, created_at
		FROM orders
		WHERE id = $1
	`
	var order Order
	err := s.db.QueryRowContext(ctx, query, id).Scan(
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

// ErrOrderNotFound is returned when order is not found.
var ErrOrderNotFound = errors.New("order not found")

// GetOutboxEventCount returns the count of unpublished events.
func (s *SafeOrderService) GetOutboxEventCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM outbox_events WHERE published_at IS NULL`,
	).Scan(&count)
	return count, err
}

// GetOrderCount returns the total number of orders.
func (s *SafeOrderService) GetOrderCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders`).Scan(&count)
	return count, err
}