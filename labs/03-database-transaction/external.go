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

// === Event Types ===

type Event struct {
	ID          string
	EventType   string
	AggregateID string
	Payload     string
}

type OutboxEvent struct {
	ID          string
	EventType   string
	AggregateID string
	Payload     string
	Status      string
	Attempts    int
	CreatedAt   time.Time
	PublishedAt *time.Time
}

// ============================================================================
// 3. WhatsApp Client (External Side Effect Demo)
// ============================================================================

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
// Distributed Order Service - External Side Effect Limitation Demo
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
// 6 & 7. Transactional Outbox Pattern & Dispatcher with Retry Mechanism (Poin 10)
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

type OutboxDispatcher struct {
	db                    *sql.DB
	broker                EventPublisher
	maxAttempts           int
	simulateCrashBeforeMark bool
}

func NewOutboxDispatcher(db *sql.DB, broker EventPublisher, maxAttempts int, crashBeforeMark bool) *OutboxDispatcher {
	return &OutboxDispatcher{
		db:                    db,
		broker:                broker,
		maxAttempts:           maxAttempts,
		simulateCrashBeforeMark: crashBeforeMark,
	}
}

func (d *OutboxDispatcher) DispatchBatch(ctx context.Context) (int, error) {
	var dispatched int
	var pendingCount int64

	err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM outbox_events WHERE status = 'pending'").Scan(&pendingCount)
	if err != nil {
		return 0, fmt.Errorf("count pending: %w", err)
	}

	for i := 0; i < int(pendingCount); i++ {
		var id, eventType, aggregateID, payload string
		var attempts int

		err := d.db.QueryRowContext(ctx, "SELECT id, event_type, aggregate_id, payload, attempts FROM outbox_events WHERE status = 'pending' LIMIT 1").Scan(&id, &eventType, &aggregateID, &payload, &attempts)
		if err != nil {
			break
		}

		if attempts >= d.maxAttempts {
			continue
		}

		event := Event{
			ID:          id,
			EventType:   eventType,
			AggregateID: aggregateID,
			Payload:     payload,
		}

		pubErr := d.broker.Publish(ctx, event)
		attempts++

		if pubErr != nil {
			_, _ = d.db.ExecContext(ctx, "UPDATE outbox_events SET attempts = $1 WHERE id = $2", attempts, id)
			continue
		}

		if d.simulateCrashBeforeMark {
			return dispatched, ErrProcessCrashed
		}

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
	db              *sql.DB
	commissionsPaid int64
	mu              sync.Mutex
}

func NewCommissionWorker(db *sql.DB) *CommissionWorker {
	return &CommissionWorker{db: db}
}

func (c *CommissionWorker) HandleEvent(ctx context.Context, consumerName string, event Event) (bool, error) {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("begin tx: %w", err)
	}

	var count int64
	row := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", event.ID)
	if err := row.Scan(&count); err == nil && count > 0 {
		_ = tx.Rollback()
		return false, nil
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO processed_events (event_id) VALUES ($1)", event.ID)
	if err != nil {
		_ = tx.Rollback()
		return false, fmt.Errorf("insert processed_events: %w", err)
	}

	c.mu.Lock()
	c.commissionsPaid++
	c.mu.Unlock()

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit: %w", err)
	}

	return true, nil
}

func (c *CommissionWorker) CommissionsPaidCount() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.commissionsPaid
}

// ============================================================================
// 11. Dead Letter Queue (DLQ) Simulation
// ============================================================================

type DLQRecord struct {
	Event     Event
	Reason    string
	Attempts  int
	FailedAt  time.Time
}

type DeadLetterQueue struct {
	mu   sync.Mutex
	records []DLQRecord
}

func NewDeadLetterQueue() *DeadLetterQueue {
	return &DeadLetterQueue{}
}

func (dlq *DeadLetterQueue) Add(event Event, reason string, attempts int) {
	dlq.mu.Lock()
	defer dlq.mu.Unlock()
	dlq.records = append(dlq.records, DLQRecord{
		Event:     event,
		Reason:    reason,
		Attempts:  attempts,
		FailedAt:  time.Now(),
	})
}

func (dlq *DeadLetterQueue) Count() int {
	dlq.mu.Lock()
	defer dlq.mu.Unlock()
	return len(dlq.records)
}

func (dlq *DeadLetterQueue) Records() []DLQRecord {
	dlq.mu.Lock()
	defer dlq.mu.Unlock()
	out := make([]DLQRecord, len(dlq.records))
	copy(out, dlq.records)
	return out
}

// OutboxDispatcherWithDLQ extends OutboxDispatcher to move failed events to DLQ
type OutboxDispatcherWithDLQ struct {
	db          *sql.DB
	broker      EventPublisher
	maxAttempts int
	dlq         *DeadLetterQueue
}

func NewOutboxDispatcherWithDLQ(db *sql.DB, broker EventPublisher, maxAttempts int, dlq *DeadLetterQueue) *OutboxDispatcherWithDLQ {
	return &OutboxDispatcherWithDLQ{
		db:          db,
		broker:      broker,
		maxAttempts: maxAttempts,
		dlq:         dlq,
	}
}

// DispatchUntilDLQ tries to publish, and if max attempts exceeded, moves to DLQ.
// Returns the number of events successfully published.
func (d *OutboxDispatcherWithDLQ) DispatchUntilDLQ(ctx context.Context, event Event, maxRetries int) (bool, error) {
	for i := 0; i < maxRetries; i++ {
		err := d.broker.Publish(ctx, event)
		if err == nil {
			return true, nil
		}
	}
	// All retries exhausted - move to DLQ
	d.dlq.Add(event, "max attempts exceeded", maxRetries)
	return false, fmt.Errorf("event moved to DLQ after %d failed attempts", maxRetries)
}

// ============================================================================
// 12. Saga Pattern & Compensating Transactions
// ============================================================================

type ReserveResourceError struct{}

func (ReserveResourceError) Error() string { return "failed to reserve resource" }

// SagaStep represents a single step in a saga with its compensating action.
type SagaStep struct {
	Action         func(ctx context.Context) error
	Compensate     func(ctx context.Context) error
}

// Saga orchestrates a series of steps, running compensating actions on failure.
type Saga struct {
	steps    []SagaStep
	executed []int // indices of completed steps
}

func NewSaga() *Saga {
	return &Saga{}
}

func (s *Saga) Then(step SagaStep) *Saga {
	s.steps = append(s.steps, step)
	return s
}

// Execute runs all steps, invoking compensations on failure.
func (s *Saga) Execute(ctx context.Context) error {
	for i, step := range s.steps {
		if err := step.Action(ctx); err != nil {
			// Compensation: run in reverse order
			for j := i; j >= 0; j-- {
				_ = s.steps[j].Compensate(ctx)
			}
			return err
		}
		s.executed = append(s.executed, i)
	}
	return nil
}

// Simple example: Vendor Payment Saga
type PaymentReservation struct {
	Reserved bool
}

func (pr *PaymentReservation) Reserve(ctx context.Context) error {
	// Simulate deducting cash from available balance
	pr.Reserved = true
	return nil // Success!
}

func (pr *PaymentReservation) Release(ctx context.Context) error {
	// Compensating action: refund the reserved amount
	pr.Reserved = false
	return nil
}