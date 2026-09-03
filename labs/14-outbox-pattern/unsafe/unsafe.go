// Package unsafe demonstrates the Dual Write Problem:
// - Order is created in PostgreSQL
// - Event is published to queue
// If DB commit succeeds but event publish fails, we have inconsistency.
package unsafe

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
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

// EventPublisher simulates publishing events to a message queue.
type EventPublisher interface {
	Publish(ctx context.Context, event OrderCreatedEvent) error
}

// MockEventPublisher simulates an event publisher with failure injection.
type MockEventPublisher struct {
	publishedEvents []OrderCreatedEvent
	mu              sync.Mutex
	failAfter       int // Fail after N successful publishes
	failCount       int64
}

func (p *MockEventPublisher) Publish(ctx context.Context, event OrderCreatedEvent) error {
	count := atomic.AddInt64(&p.failCount, 1)
	if p.failAfter > 0 && count > int64(p.failAfter) {
		return errors.New("publish failed: connection lost")
	}
	p.mu.Lock()
	p.publishedEvents = append(p.publishedEvents, event)
	p.mu.Unlock()
	return nil
}

func (p *MockEventPublisher) GetPublishedEvents() []OrderCreatedEvent {
	p.mu.Lock()
	defer p.mu.Unlock()
	result := make([]OrderCreatedEvent, len(p.publishedEvents))
	copy(result, p.publishedEvents)
	return result
}

func (p *MockEventPublisher) GetPublishCount() int64 {
	return atomic.LoadInt64(&p.failCount)
}

// OrderRepository handles order persistence.
type OrderRepository interface {
	Create(ctx context.Context, order Order) error
	FindByID(ctx context.Context, id string) (Order, error)
	Count(ctx context.Context) (int64, error)
}

// UnsafeOrderService demonstrates the dual write problem.
// BUG: Order creation and event publishing happen sequentially,
// NOT in a transaction. This means:
// - DB commit succeeds, event publish fails -> INCONSISTENT STATE
// - Event publish succeeds, DB commit fails -> ORDER LOST
type UnsafeOrderService struct {
	repo      OrderRepository
	publisher EventPublisher
	db        *sql.DB
}

// NewUnsafeOrderService creates a service with dual write bug.
func NewUnsafeOrderService(repo OrderRepository, publisher EventPublisher, db *sql.DB) *UnsafeOrderService {
	return &UnsafeOrderService{
		repo:      repo,
		publisher: publisher,
		db:        db,
	}
}

// CreateOrder demonstrates the dual write bug.
// The order is inserted into DB first, THEN event is published.
// If publishing fails after DB insert, the order exists in DB
// but no event was published -> system is now inconsistent.
func (s *UnsafeOrderService) CreateOrder(ctx context.Context, order Order) (Order, error) {
	// Step 1: Insert order into database (autonomous transaction)
	if err := s.repo.Create(ctx, order); err != nil {
		return Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	// Step 2: Publish event to message queue
	// BUG: If this fails, order exists in DB but event was never sent
	// This creates the dual write problem
	event := OrderCreatedEvent{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
	}

	if err := s.publisher.Publish(ctx, event); err != nil {
		// LOG: We cannot roll back the DB transaction at this point
		// The order exists, but consumers won't see it
		return Order{}, fmt.Errorf("failed to publish order created event: %w", err)
	}

	return order, nil
}

// GetOrder retrieves an order by ID.
func (s *UnsafeOrderService) GetOrder(ctx context.Context, id string) (Order, error) {
	return s.repo.FindByID(ctx, id)
}

// ListEvents returns all published events (for testing).
func (p *MockEventPublisher) ListEvents() []OrderCreatedEvent {
	p.mu.Lock()
	defer p.mu.Unlock()
	result := make([]OrderCreatedEvent, len(p.publishedEvents))
	copy(result, p.publishedEvents)
	return result
}
