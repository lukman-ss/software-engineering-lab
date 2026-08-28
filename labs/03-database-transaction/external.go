package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// === Error Definitions ===

var (
	ErrExternalServiceFailed = errors.New("external WhatsApp / ERP service failed after payment created")
	ErrHTTPCallFailed        = errors.New("external HTTP call failed inside transaction")
	ErrProcessCrashed        = errors.New("process crashed after DB commit but before event publish")
	ErrInvoicePaid           = errors.New("invoice already paid")
	ErrBrokerFailed          = errors.New("message broker publish failed (temporary)")
)

// === 3. WhatsApp Client (External Side Effect Demo) ===

type WhatsAppClient struct {
	mu           sync.Mutex
	sentMessages []string
	failureCount int64
	failAfter    int
}

func NewWhatsAppClient(failAfter int) *WhatsAppClient {
	return &WhatsAppClient{failAfter: failAfter}
}

func (w *WhatsAppClient) SendPaymentNotification(ctx context.Context, orderID int, amount float64) error {
	count := atomic.AddInt64(&w.failureCount, 1)
	if w.failAfter > 0 && count >= int64(w.failAfter) {
		return ErrExternalServiceFailed
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	w.sentMessages = append(w.sentMessages, fmt.Sprintf("Order %d: Paid %.2f", orderID, amount))
	return nil
}

func (w *WhatsAppClient) SentCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.sentMessages)
}

func (w *WhatsAppClient) GetMessages() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]string, len(w.sentMessages))
	copy(out, w.sentMessages)
	return out
}

// ============================================================================
// Distributed Service - External Side Effect Limitation Demo
// ============================================================================

type DistributedOrderService struct {
	db       *sql.DB
	whatsapp *WhatsAppClient
}

func NewDistributedOrderService(db *sql.DB, whatsapp *WhatsAppClient) *DistributedOrderService {
	return &DistributedOrderService{db: db, whatsapp: whatsapp}
}

func (s *DistributedOrderService) ProcessPaymentWithExternalSideEffect(ctx context.Context, orderID int, amount float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, "INSERT INTO payments (order_id, amount, status) VALUES ($1, $2, 'completed')", orderID, amount)
	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE invoices SET status = 'paid' WHERE order_id = $1", orderID)
	if err != nil {
		return fmt.Errorf("update invoice: %w", err)
	}

	err = s.whatsapp.SendPaymentNotification(ctx, orderID, amount)
	if err != nil {
		return fmt.Errorf("external whatsapp call failed: %w", err)
	}

	return errors.New("simulated ERP integration error — triggering database ROLLBACK")
}

// ============================================================================
// 4. HTTP Client - Anti-Pattern Demonstration
// ============================================================================

type HTTPClient struct {
	mu            sync.Mutex
	callCount     int64
	latencyMillis int
	failAfter     int
}

func NewHTTPClient(latencyMillis int, failAfter int) *HTTPClient {
	return &HTTPClient{latencyMillis: latencyMillis, failAfter: failAfter}
}

func (c *HTTPClient) Ping(ctx context.Context) error {
	count := atomic.AddInt64(&c.callCount, 1)
	if c.failAfter > 0 && count >= int64(c.failAfter) {
		return ErrHTTPCallFailed
	}

	if c.latencyMillis > 0 {
		select {
		case <-time.After(time.Duration(c.latencyMillis) * time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

type InvoiceServiceWithHTTPCall struct {
	db   *sql.DB
	http *HTTPClient
}

func NewInvoiceServiceWithHTTPCall(db *sql.DB, http *HTTPClient) *InvoiceServiceWithHTTPCall {
	return &InvoiceServiceWithHTTPCall{db: db, http: http}
}

func (s *InvoiceServiceWithHTTPCall) PayInvoiceWithHTTPInsideTx(ctx context.Context, invoiceID int, failAfterDBCommit bool) (time.Duration, error) {
	start := time.Now()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE invoices SET status = 'paid' WHERE order_id = $1", invoiceID)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("update invoice: %w", err)
	}

	err = s.http.Ping(ctx)
	elapsed := time.Since(start)

	if err != nil {
		_ = tx.Rollback()
		return elapsed, fmt.Errorf("http call failed: %w", err)
	}

	if failAfterDBCommit {
		_ = tx.Commit()
		return elapsed, ErrProcessCrashed
	}

	if err := tx.Commit(); err != nil {
		return elapsed, fmt.Errorf("commit: %w", err)
	}

	return elapsed, nil
}

// ============================================================================
// 5. Dual-Write Problem Simulation
// ============================================================================

type Event struct {
	ID          string
	EventType   string
	AggregateID string
	Payload     string
}

type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}

type InMemoryBroker struct {
	mu        sync.Mutex
	published []Event
	failCalls int64
	failUpTo  int
}

func NewInMemoryBroker(failUpTo int) *InMemoryBroker {
	return &InMemoryBroker{failUpTo: failUpTo}
}

func (b *InMemoryBroker) Publish(ctx context.Context, event Event) error {
	count := atomic.AddInt64(&b.failCalls, 1)
	if int(count) <= b.failUpTo {
		return ErrBrokerFailed
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.published = append(b.published, event)
	return nil
}

func (b *InMemoryBroker) PublishedEvents() []Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]Event, len(b.published))
	copy(out, b.published)
	return out
}

type InvoiceServiceDualWrite struct {
	db        *sql.DB
	publisher EventPublisher
}

func NewInvoiceServiceDualWrite(db *sql.DB, publisher EventPublisher) *InvoiceServiceDualWrite {
	return &InvoiceServiceDualWrite{db: db, publisher: publisher}
}

func (s *InvoiceServiceDualWrite) PayInvoiceDualWrite(ctx context.Context, invoiceID int, crashAfterCommit bool) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE invoices SET status = 'paid' WHERE order_id = $1", invoiceID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("update invoice: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	if crashAfterCommit {
		return ErrProcessCrashed
	}

	event := Event{
		ID:          fmt.Sprintf("evt_%d", invoiceID),
		EventType:   "InvoicePaid",
		AggregateID: fmt.Sprintf("%d", invoiceID),
		Payload:     fmt.Sprintf(`{"invoice_id": %d}`, invoiceID),
	}
	return s.publisher.Publish(ctx, event)
}

// ============================================================================
// 6 & 7. Transactional Outbox Pattern & Dispatcher with Retry & Duplicate Delivery
// ============================================================================

type InvoiceServiceOutbox struct {
	db *sql.DB
}

func NewInvoiceServiceOutbox(db *sql.DB) *InvoiceServiceOutbox {
	return &InvoiceServiceOutbox{db: db}
}

func (s *InvoiceServiceOutbox) PayInvoiceWithOutbox(ctx context.Context, invoiceID int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE invoices SET status = 'paid' WHERE order_id = $1", invoiceID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("update invoice: %w", err)
	}

	eventID := fmt.Sprintf("evt_%d", invoiceID)
	eventPayload := fmt.Sprintf(`{"invoice_id": %d, "status": "paid"}`, invoiceID)

	_, err = tx.ExecContext(ctx, `
		INSERT INTO outbox_events
		(id, event_type, aggregate_id, payload, status, attempts, created_at, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULL)
	`, eventID, "InvoicePaid", fmt.Sprintf("%d", invoiceID), eventPayload, "pending", 0, time.Now())

	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("insert outbox: %w", err)
	}

	return tx.Commit()
}

func (s *InvoiceServiceOutbox) CountOutboxEvents(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM outbox_events WHERE status = 'pending'").Scan(&count)
	return count, err
}

// OutboxDispatcher reads pending events, publishes them via EventPublisher, and marks them published.
// Supports retry mechanism with max attempts and duplicate delivery simulation.
type OutboxDispatcher struct {
	db          *sql.DB
	broker      EventPublisher
	maxAttempts int
	simulateCrashBeforeMark bool
}

func NewOutboxDispatcher(db *sql.DB, broker EventPublisher, maxAttempts int, crashBeforeMark bool) *OutboxDispatcher {
	return &OutboxDispatcher{
		db:                     db,
		broker:                 broker,
		maxAttempts:            maxAttempts,
		simulateCrashBeforeMark: crashBeforeMark,
	}
}

func (d *OutboxDispatcher) DispatchBatch(ctx context.Context) (int, error) {
	rows, err := d.db.QueryContext(ctx, `
		SELECT id, event_type, aggregate_id, payload, attempts
		FROM outbox_events
		WHERE status = 'pending'
	`)
	if err != nil {
		return 0, fmt.Errorf("query outbox: %w", err)
	}
	defer rows.Close()

	dispatched := 0
	for rows.Next() {
		var id, eventType, aggregateID, payload string
		var attempts int
		if err := rows.Scan(&id, &eventType, &aggregateID, &payload, &attempts); err != nil {
			return dispatched, fmt.Errorf("scan: %w", err)
		}

		if attempts >= d.maxAttempts {
			continue // Max attempts reached, skip or move to dead letter queue
		}

		event := Event{
			ID:          id,
			EventType:   eventType,
			AggregateID: aggregateID,
			Payload:     payload,
		}

		// Try publishing with retry
		pubErr := d.broker.Publish(ctx, event)
		attempts++

		if pubErr != nil {
			// Increment attempts in DB
			_, _ = d.db.ExecContext(ctx, "UPDATE outbox_events SET attempts = $1 WHERE id = $2", attempts, id)
			continue
		}

		// Simulate crash before marking as published -> causes duplicate delivery on retry!
		if d.simulateCrashBeforeMark {
			return dispatched, ErrProcessCrashed
		}

		// Mark as published
		_, err = d.db.ExecContext(ctx, "UPDATE outbox_events SET status = 'published', attempts = $1 WHERE id = $2", attempts, id)
		if err != nil {
			return dispatched, fmt.Errorf("mark published: %w", err)
		}

		dispatched++
	}

	return dispatched, nil
}

// ============================================================================
// 9. Idempotent Consumer Implementation
// ============================================================================

type CommissionWorker struct {
	db             *sql.DB
	commissionsPaid int64
	mu             sync.Mutex
}

func NewCommissionWorker(db *sql.DB) *CommissionWorker {
	return &CommissionWorker{db: db}
}

// HandleEvent processes the event idempotently using processed_events deduplication table.
// Connects to Lab #01 Idempotency concepts (unique constraint on consumer_name, event_id).
func (c *CommissionWorker) HandleEvent(ctx context.Context, consumerName string, event Event) (bool, error) {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("begin tx: %w", err)
	}

	// 1. Check & insert into processed_events (deduplication / idempotency check)
	// In PostgreSQL, UNIQUE(consumer_name, event_id) prevents race conditions.
	// In our mockdb, we check and insert.
	var exists bool
	row := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE consumer_name = $1 AND event_id = $2", consumerName, event.ID)
	var count int64
	if err := row.Scan(&count); err == nil && count > 0 {
		exists = true
	}

	if exists {
		_ = tx.Rollback()
		return false, nil // Already processed, safely skipped (Idempotent!)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3)", consumerName, event.ID, time.Now())
	if err != nil {
		_ = tx.Rollback()
		return false, fmt.Errorf("insert processed_events: %w", err)
	}

	// 2. Perform business operation (e.g., pay commission once)
	c.mu.Lock()
	c.commissionsPaid++
	c.mu.Unlock()

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit: %w", err)
	}

	return true, nil // Processed successfully for the first time
}

func (c *CommissionWorker) CommissionsPaidCount() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.commissionsPaid
}
