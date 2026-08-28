package safe

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
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
	ID            string          `json:"id"`
	AggregateType string          `json:"aggregate_type"`
	AggregateID   string          `json:"aggregate_id"`
	EventType     string          `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
	CreatedAt     time.Time       `json:"created_at"`
	PublishedAt   *time.Time      `json:"published_at"`
	Attempts      int             `json:"attempts"`
	NextAttemptAt time.Time       `json:"next_attempt_at"`
}

// EventPublisher defines the interface for publishing events to an external message broker.
type EventPublisher interface {
	Publish(ctx context.Context, event OutboxEvent) error
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
	orderID := generateID()
	eventID := generateID()

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
		`INSERT INTO orders (id, customer_id, status, created_at) VALUES ($1, $2, $3, $4)`,
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

// OutboxWorker processes unpublished outbox events with batching, retry, backoff,
// context cancellation, and graceful shutdown.
type OutboxWorker struct {
	db           *sql.DB
	publisher    EventPublisher
	batchSize    int
	pollInterval time.Duration
	maxAttempts  int
	backoffBase  time.Duration

	// Coordination for graceful shutdown
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// WorkerConfig holds configuration for the OutboxWorker.
type WorkerConfig struct {
	BatchSize    int
	PollInterval time.Duration
	MaxAttempts  int
	BackoffBase  time.Duration
}

// DefaultWorkerConfig returns sensible defaults.
func DefaultWorkerConfig() WorkerConfig {
	return WorkerConfig{
		BatchSize:    50,
		PollInterval: 1 * time.Second,
		MaxAttempts:  5,
		BackoffBase:  2 * time.Second,
	}
}

// NewOutboxWorker creates a new outbox worker instance.
func NewOutboxWorker(db *sql.DB, publisher EventPublisher, cfg WorkerConfig) *OutboxWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &OutboxWorker{
		db:           db,
		publisher:    publisher,
		batchSize:    cfg.BatchSize,
		pollInterval: cfg.PollInterval,
		maxAttempts:  cfg.MaxAttempts,
		backoffBase:  cfg.BackoffBase,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start runs the outbox worker in the background.
func (w *OutboxWorker) Start() {
	w.wg.Add(1)
	go w.run()
}

// Stop initiates graceful shutdown, waiting for in-flight batches to complete.
func (w *OutboxWorker) Stop() {
	w.cancel()
	w.wg.Wait()
}

func (w *OutboxWorker) run() {
	defer w.wg.Done()

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	log.Printf("[OutboxWorker] Started polling every %v", w.pollInterval)

	for {
		select {
		case <-w.ctx.Done():
			log.Println("[OutboxWorker] Shutting down gracefully...")
			return
		case <-ticker.C:
			if err := w.processBatch(w.ctx); err != nil {
				if !errors.Is(err, context.Canceled) {
					log.Printf("[OutboxWorker] Error processing batch: %v", err)
				}
			}
		}
	}
}

// processBatch fetches and processes a batch of unpublished events.
func (w *OutboxWorker) processBatch(ctx context.Context) error {
	query := `
		SELECT id, aggregate_type, aggregate_id, event_type, payload, created_at, published_at, attempts, next_attempt_at
		FROM outbox_events
		WHERE published_at IS NULL AND next_attempt_at <= CURRENT_TIMESTAMP
		ORDER BY created_at ASC
		LIMIT $1
	`

	rows, err := w.db.QueryContext(ctx, query, w.batchSize)
	if err != nil {
		return fmt.Errorf("failed to query outbox events: %w", err)
	}
	defer rows.Close()

	var events []OutboxEvent
	for rows.Next() {
		var e OutboxEvent
		var pubAt sql.NullTime
		if err := rows.Scan(
			&e.ID, &e.AggregateType, &e.AggregateID, &e.EventType,
			&e.Payload, &e.CreatedAt, &pubAt, &e.Attempts, &e.NextAttemptAt,
		); err != nil {
			return fmt.Errorf("failed to scan outbox event: %w", err)
		}
		if pubAt.Valid {
			e.PublishedAt = &pubAt.Time
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("row iteration error: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	log.Printf("[OutboxWorker] Processing batch of %d events", len(events))

	for _, event := range events {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := w.processEvent(ctx, event); err != nil {
			log.Printf("[OutboxWorker] Failed to publish event %s: %v", event.ID, err)
		}
	}

	return nil
}

// processEvent publishes an individual event and updates its status or applies backoff.
func (w *OutboxWorker) processEvent(ctx context.Context, event OutboxEvent) error {
	err := w.publisher.Publish(ctx, event)

	now := time.Now()
	if err != nil {
		newAttempts := event.Attempts + 1
		backoffDuration := w.backoffBase * time.Duration(1<<(newAttempts-1))
		nextAttempt := now.Add(backoffDuration)

		_, updateErr := w.db.ExecContext(ctx,
			`UPDATE outbox_events SET attempts = $1, next_attempt_at = $2 WHERE id = $3`,
			newAttempts, nextAttempt, event.ID,
		)
		if updateErr != nil {
			return fmt.Errorf("failed to update event retry state: %w", updateErr)
		}

		return fmt.Errorf("publish failed (attempt %d/%d, retry at %v): %w",
			newAttempts, w.maxAttempts, nextAttempt, err)
	}

	_, updateErr := w.db.ExecContext(ctx,
		`UPDATE outbox_events SET published_at = $1, attempts = attempts + 1 WHERE id = $2`,
		now, event.ID,
	)
	if updateErr != nil {
		return fmt.Errorf("failed to mark event as published: %w", updateErr)
	}

	log.Printf("[OutboxWorker] Event %s (%s) published successfully", event.ID, event.EventType)
	return nil
}

// MockEventBroker simulates a message queue broker.
type MockEventBroker struct {
	published []OutboxEvent
	mu        sync.Mutex
	shouldFail bool
}

func (b *MockEventBroker) Publish(ctx context.Context, event OutboxEvent) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.shouldFail {
		return errors.New("simulated broker connection failure")
	}

	b.published = append(b.published, event)
	return nil
}

func (b *MockEventBroker) GetPublished() []OutboxEvent {
	b.mu.Lock()
	defer b.mu.Unlock()
	res := make([]OutboxEvent, len(b.published))
	copy(res, b.published)
	return res
}

func (b *MockEventBroker) SetShouldFail(fail bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.shouldFail = fail
}

// generateID generates a simple unique ID using UUID.
func generateID() string {
	return uuid.New().String()
}