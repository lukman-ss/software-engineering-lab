package transaction

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

// === Failure Injection ===
// Deterministic failure controls for testing crash scenarios.

type FailureMode struct {
	// Outbox/Dispatcher related
	FailAfterDBCommit   bool // Simulate crash after DB commit but before publish
	FailExternalService bool
}

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

// InvoicePaidPayload is the canonical payload for InvoicePaid events.
// Use explicit structs, not fmt.Sprintf — guarantees valid JSON, correct types.
// ponytail: no encryption/PII scrubbing; add when regulatory required
type InvoicePaidPayload struct {
	EventID    string `json:"event_id"`
	InvoiceID  int    `json:"invoice_id"`
	OccurredAt string `json:"occurred_at"` // RFC3339
}

// MarshalInvoicePaidPayload serializes the payload and panics on failure (impossible for this struct)
func MarshalInvoicePaidPayload(p InvoicePaidPayload) string {
	b, err := json.Marshal(p)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal InvoicePaidPayload: %v", err))
	}
	return string(b)
}

// UnmarshalInvoicePaidPayload deserializes a raw JSON string into InvoicePaidPayload
func UnmarshalInvoicePaidPayload(raw string) (InvoicePaidPayload, error) {
	var p InvoicePaidPayload
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		return p, fmt.Errorf("invalid InvoicePaid payload: %w", err)
	}
	if p.EventID == "" {
		return p, fmt.Errorf("payload missing event_id")
	}
	if p.InvoiceID == 0 {
		return p, fmt.Errorf("payload missing invoice_id")
	}
	if p.OccurredAt == "" {
		return p, fmt.Errorf("payload missing occurred_at")
	}
	return p, nil
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
		if err := tx.Commit(); err != nil {
			return elapsed, fmt.Errorf("simulated crash pre-requisite commit failed: %w", err)
		}
		return elapsed, ErrProcessCrashed
	}

	if err := tx.Commit(); err != nil {
		return elapsed, fmt.Errorf("commit: %w", err)
	}

	return elapsed, nil
}

// ============================================================================
// 4b. Blocking HTTP Client for Transaction Lifetime Testing
// ============================================================================

// BlockingHTTPClient blocks until channel is closed - for testing transaction lifetime
type BlockingHTTPClient struct {
	blockChan    chan struct{}
	delay        time.Duration
	txOpenSignal chan struct{} // closed when Ping is called (i.e. tx is open)
}

func NewBlockingHTTPClient(delay time.Duration) *BlockingHTTPClient {
	return &BlockingHTTPClient{
		blockChan:    make(chan struct{}),
		delay:        delay,
		txOpenSignal: make(chan struct{}),
	}
}

// WaitUntilTxOpen blocks until the service reaches the external call (i.e. transaction is open)
func (c *BlockingHTTPClient) WaitUntilTxOpen() {
	<-c.txOpenSignal
}

func (c *BlockingHTTPClient) Ping(ctx context.Context) error {
	// Signal that we're inside the external call (transaction is open)
	select {
	case <-c.txOpenSignal:
		// already signalled
	default:
		close(c.txOpenSignal)
	}

	if c.delay > 0 {
		time.Sleep(c.delay)
		return nil
	}
	select {
	case <-c.blockChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *BlockingHTTPClient) Release() {
	select {
	case c.blockChan <- struct{}{}:
	default:
	}
}

// InvoiceServiceWithBlockingHTTP uses BlockingHTTPClient to simulate long-running external call
type InvoiceServiceWithBlockingHTTP struct {
	db     *sql.DB
	http   *BlockingHTTPClient
	txOpen int32 // 1 = transaction is open, 0 = closed; accessed via atomic
}

func NewInvoiceServiceWithBlockingHTTP(db *sql.DB, http *BlockingHTTPClient) *InvoiceServiceWithBlockingHTTP {
	return &InvoiceServiceWithBlockingHTTP{db: db, http: http}
}

// IsTxOpen returns true if a transaction is currently open in this service
func (s *InvoiceServiceWithBlockingHTTP) IsTxOpen() bool {
	return atomic.LoadInt32(&s.txOpen) == 1
}

func (s *InvoiceServiceWithBlockingHTTP) PayInvoiceWithBlocking(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	atomic.StoreInt32(&s.txOpen, 1)

	defer func() {
		atomic.StoreInt32(&s.txOpen, 0)
	}()

	_, err = tx.ExecContext(ctx, "UPDATE invoices SET status = 'paid' WHERE order_id = $1", 101)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("update invoice: %w", err)
	}

	// External call that blocks - transaction remains open
	err = s.http.Ping(ctx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("http ping failed: %w", err)
	}

	return tx.Commit()
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

func (s *InvoiceServiceDualWrite) PayInvoiceDualWrite(ctx context.Context, invoiceID int, failMode *FailureMode) error {
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

	if failMode != nil && failMode.FailAfterDBCommit {
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
// 6 & 7. Transactional Outbox Pattern (Concept) & Simplified Dispatcher
//
// NOTE: This is a minimal demonstration of the outbox concept for educational
// purposes. Production-grade dispatcher features — FOR UPDATE SKIP LOCKED
// concurrency control, exponential backoff retry, crash-after-publish recovery,
// Dead Letter Queue (DLQ), poison message handling, replay mechanisms —
// are covered in Lab 14 — Outbox Pattern.
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

	result, err := tx.ExecContext(ctx, "UPDATE invoices SET status = 'paid' WHERE order_id = $1", invoiceID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("update invoice: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		_ = tx.Rollback()
		return ErrInvoicePaid
	}

	eventID := fmt.Sprintf("evt_%d_%d", invoiceID, time.Now().UnixNano())
	eventPayload := MarshalInvoicePaidPayload(InvoicePaidPayload{
		EventID:    eventID,
		InvoiceID:  invoiceID,
		OccurredAt: time.Now().Format(time.RFC3339),
	})

	_, err = tx.ExecContext(ctx, `
		INSERT INTO outbox_events
		(id, event_type, aggregate_id, payload, status, attempts, created_at, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULL)
	`, eventID, "InvoicePaid", fmt.Sprintf("%d", invoiceID), eventPayload, "pending", 0, time.Now())

	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("insert outbox: %w", err)
	}

	log.Printf("[OUTBOX] InvoicePaid event stored for invoice %d", invoiceID)
	return tx.Commit()
}

func (s *InvoiceServiceOutbox) CountOutboxEvents(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM outbox_events WHERE status = 'pending'").Scan(&count)
	return count, err
}

// PayInvoiceWithOutboxInjectError demonstrates rollback when outbox insert fails
func (s *InvoiceServiceOutbox) PayInvoiceWithOutboxInjectError(ctx context.Context, invoiceID int, injectError bool) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	result, err := tx.ExecContext(ctx, "UPDATE invoices SET status = 'paid' WHERE order_id = $1", invoiceID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("update invoice: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		_ = tx.Rollback()
		return ErrInvoicePaid
	}

	// Try to update first (this should succeed)
	_, err = tx.ExecContext(ctx, "INSERT INTO outbox_events (id, event_type, aggregate_id, payload, status, attempts, created_at, published_at) VALUES ('bad_id', 'InvoicePaid', '101', '{}', 'pending', 0, NULL)")
	if injectError && err == nil {
		// Force an error without letting it succeed - duplicate key simulation
		err = errors.New("simulated outbox constraint violation")
	}

	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("insert outbox: %w", err)
	}

	log.Printf("[OUTBOX] InvoicePaid event stored for invoice %d", invoiceID)
	return tx.Commit()
}

// ============================================================================
// Simplified Outbox Dispatcher
//
// This dispatcher demonstrates the basic concept: read pending events from the
// outbox table, publish them to a broker, and mark them as published.
//
// Production-grade concerns (retry with backoff, concurrency control with
// FOR UPDATE SKIP LOCKED, crash recovery, DLQ handling) are NOT implemented here.
// See Lab 14 — Outbox Pattern for the full implementation.
// ============================================================================

type OutboxDispatcher struct {
	db     *sql.DB
	broker EventPublisher
}

func NewOutboxDispatcher(db *sql.DB, broker EventPublisher) *OutboxDispatcher {
	return &OutboxDispatcher{db: db, broker: broker}
}

func (d *OutboxDispatcher) DispatchBatch(ctx context.Context) (int, error) {
	var dispatched int

	rows, err := d.db.QueryContext(ctx, "SELECT id, event_type, aggregate_id, payload FROM outbox_events WHERE status = 'pending'")
	if err != nil {
		return 0, fmt.Errorf("query pending: %w", err)
	}

	type pendingEvent struct {
		id, eventType, aggregateID, payload string
	}
	var pending []pendingEvent

	for rows.Next() {
		var e pendingEvent
		if err := rows.Scan(&e.id, &e.eventType, &e.aggregateID, &e.payload); err == nil {
			pending = append(pending, e)
		}
	}
	rows.Close()

	for _, e := range pending {
		event := Event{
			ID:          e.id,
			EventType:   e.eventType,
			AggregateID: e.aggregateID,
			Payload:     e.payload,
		}

		log.Printf("[DISPATCHER] publishing event=%s", event.ID)
		if pubErr := d.broker.Publish(ctx, event); pubErr != nil {
			log.Printf("[DISPATCHER] publish failed for event=%s: %v — Lab 14 covers retry/backoff", event.ID, pubErr)
			continue
		}

		now := time.Now()
		_, err = d.db.ExecContext(ctx, "UPDATE outbox_events SET status = 'published', published_at = $1 WHERE id = $2", now, event.ID)
		if err != nil {
			log.Printf("[DISPATCHER] failed to mark event=%s as published: %v", event.ID, err)
			continue
		}

		log.Printf("[DISPATCHER] event=%s marked as published", event.ID)
		dispatched++
	}

	return dispatched, nil
}

// ============================================================================
// 9. Idempotent Consumer (Conceptual Demo)
//
// The consumer pattern: dedup marker + business state committed in a single
// transaction using INSERT ... ON CONFLICT DO NOTHING.
//
// Production concerns (per-consumer dedup keys, concurrency control, replay
// handling) are covered in Lab 14 — Outbox Pattern.
// ============================================================================

// CommissionWorker processes events idempotently.
// KEY DESIGN: business state (commissions) and processed_events are
// committed in a single transaction.
type CommissionWorker struct {
	db                         *sql.DB
	observedBusinessExecutions int64 // observability only, NOT source of truth
	mu                         sync.Mutex
}

func NewCommissionWorker(db *sql.DB) *CommissionWorker {
	return &CommissionWorker{db: db}
}

// HandleEvent processes an event idempotently using consumer_name + event_id.
// Uses atomic INSERT ON CONFLICT DO NOTHING pattern for race-free deduplication.
// Returns (true, nil) = processed, (false, nil) = duplicate skipped, (false, error) = failure.
func (c *CommissionWorker) HandleEvent(ctx context.Context, consumerName string, event Event) (bool, error) {
	log.Printf("[CONSUMER:%s] processing %s", consumerName, event.ID)

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("begin tx: %w", err)
	}

	// 1. Atomically try to claim the event (ON CONFLICT DO NOTHING)
	result, err := tx.ExecContext(ctx,
		"INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3) ON CONFLICT (consumer_name, event_id) DO NOTHING",
		consumerName, event.ID, time.Now())
	if err != nil {
		_ = tx.Rollback()
		return false, fmt.Errorf("insert processed_events: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		_ = tx.Rollback()
		log.Printf("[CONSUMER:%s] duplicate %s skipped", consumerName, event.ID)
		return false, nil // Idempotent skip
	}

	// 2. Business operation - INSERT into commissions table (within same transaction)
	_, err = tx.ExecContext(ctx,
		"INSERT INTO commissions (event_id, invoice_id, amount, created_at) VALUES ($1, $2, $3, $4)",
		event.ID, event.AggregateID, 10000.0, time.Now())
	if err != nil {
		_ = tx.Rollback()
		return false, fmt.Errorf("insert commission: %w", err)
	}

	// 3. Commit both dedup marker + business state atomically
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit consumer tx: %w", err)
	}

	// Observability counter (out of tx for metrics, NOT business state)
	c.mu.Lock()
	c.observedBusinessExecutions++
	c.mu.Unlock()

	log.Printf("[CONSUMER:%s] completed %s", consumerName, event.ID)
	return true, nil
}

func (c *CommissionWorker) ObservedBusinessExecutions() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.observedBusinessExecutions
}

// GetDBCommissionCount returns the actual business state count from the database
func (c *CommissionWorker) GetDBCommissionCount(ctx context.Context) (int64, error) {
	var count int64
	err := c.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&count)
	return count, err
}

// ============================================================================
// Saga Pattern & Compensating Transactions (Conceptual Demo)
//
// Saga compensation for semantic undo, not technical rollback.
// ============================================================================

type SagaStep struct {
	Action     func(ctx context.Context) error
	Compensate func(ctx context.Context) error
}

type SagaExecutionError struct {
	OriginalError      error
	CompensationErrors []error
}

func (e *SagaExecutionError) Error() string {
	return fmt.Sprintf("saga failed: original=%v, compensation_errors=%v", e.OriginalError, e.CompensationErrors)
}

type Saga struct {
	steps       []SagaStep
	executed    []int
	compensated []int
}

func NewSaga() *Saga {
	return &Saga{}
}

func (s *Saga) Then(step SagaStep) *Saga {
	s.steps = append(s.steps, step)
	return s
}

func (s *Saga) Execute(ctx context.Context) error {
	for i, step := range s.steps {
		if err := step.Action(ctx); err != nil {
			// Compensate only the steps that successfully completed
			// Steps are compensated in reverse order
			// Continue compensating even if a compensation fails
			var compErrors []error
			for j := len(s.executed) - 1; j >= 0; j-- {
				if compErr := s.steps[s.executed[j]].Compensate(ctx); compErr != nil {
					compErrors = append(compErrors, compErr)
				}
				s.compensated = append(s.compensated, s.executed[j])
			}
			_ = s.executed
			if len(compErrors) > 0 {
				return &SagaExecutionError{
					OriginalError:      err,
					CompensationErrors: compErrors,
				}
			}
			return err
		}
		s.executed = append(s.executed, i)
	}
	return nil
}

// GetExecutedSteps returns the list of successfully executed step indices
func (s *Saga) GetExecutedSteps() []int {
	return s.executed
}

// GetCompensatedSteps returns the list of compensated step indices in order of compensation
func (s *Saga) GetCompensatedSteps() []int {
	return s.compensated
}

// ============================================================================
// IdempotentCompensator: Compensation must be idempotent
// Running RefundPayment order-123 twice must not double-refund.
// ============================================================================

type IdempotentCompensator struct {
	mu       sync.Mutex
	executed map[string]bool
}

func NewIdempotentCompensator() *IdempotentCompensator {
	return &IdempotentCompensator{executed: make(map[string]bool)}
}

// Compensate runs action only if key has not been used before.
// Returns (executed, nil) if action was run, (false, nil) if skipped (idempotent).
func (c *IdempotentCompensator) Compensate(key string, action func() error) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.executed[key] {
		return false, nil // idempotent skip
	}
	if err := action(); err != nil {
		return false, err
	}
	c.executed[key] = true
	return true, nil
}
