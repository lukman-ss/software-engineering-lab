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

// === Failure Injection (Poin 22) ===

type FailureMode struct {
	FailAfterDBCommit   bool
	FailExternalService bool
	FailPublishAttempts int
	CrashAfterPublish   bool
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
		_ = tx.Commit()
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
	db          *sql.DB
	http        *BlockingHTTPClient
	txOpen      int32 // 1 = transaction is open, 0 = closed; accessed via atomic
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
// 6 & 7. Transactional Outbox Pattern & Dispatcher with Retry Mechanism
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

// OutboxDispatcher publishes pending events with retry support and optional DLQ integration
type OutboxDispatcher struct {
	db          *sql.DB
	broker      EventPublisher
	maxAttempts int
	failMode    *FailureMode
	dlq         *DeadLetterQueue
}

func NewOutboxDispatcher(db *sql.DB, broker EventPublisher, maxAttempts int, failMode *FailureMode) *OutboxDispatcher {
	return &OutboxDispatcher{db: db, broker: broker, maxAttempts: maxAttempts, failMode: failMode}
}

func NewOutboxDispatcherWithDLQ(db *sql.DB, broker EventPublisher, maxAttempts int, failMode *FailureMode, dlq *DeadLetterQueue) *OutboxDispatcher {
	return &OutboxDispatcher{db: db, broker: broker, maxAttempts: maxAttempts, failMode: failMode, dlq: dlq}
}

func (d *OutboxDispatcher) DispatchBatch(ctx context.Context) (int, error) {
	var dispatched int

	// Get ALL pending event IDs first to avoid infinite loops on failing events
	rows, err := d.db.QueryContext(ctx, "SELECT id, event_type, aggregate_id, payload, attempts FROM outbox_events WHERE status = 'pending'")
	if err != nil {
		return 0, fmt.Errorf("query pending: %w", err)
	}

	type pendingEvent struct {
		id, eventType, aggregateID, payload string
		attempts int
	}
	var pending []pendingEvent

	for rows.Next() {
		var e pendingEvent
		if err := rows.Scan(&e.id, &e.eventType, &e.aggregateID, &e.payload, &e.attempts); err == nil {
			pending = append(pending, e)
		}
	}
	rows.Close()

	for _, e := range pending {
		id := e.id
		eventType := e.eventType
		aggregateID := e.aggregateID
		payload := e.payload
		attempts := e.attempts

		if attempts >= d.maxAttempts {
			// Move to DLQ if there's a broker failure after max attempts
			if d.dlq != nil {
				event := Event{ID: id, EventType: eventType, AggregateID: aggregateID, Payload: payload}
				d.dlq.Add(event, "max attempts exceeded", attempts)
				_, _ = d.db.ExecContext(ctx, "UPDATE outbox_events SET status = 'dead_lettered' WHERE id = $1", id)
				log.Printf("[DISPATCHER] event=%s moved to DLQ after %d attempts", id, attempts)
			}
			continue
		}

		event := Event{
			ID:          id,
			EventType:   eventType,
			AggregateID: aggregateID,
			Payload:     payload,
		}

		log.Printf("[DISPATCHER] publishing event=%s attempt=%d", event.ID, attempts+1)
		pubErr := d.broker.Publish(ctx, event)
		attempts++

		if pubErr != nil {
			log.Printf("[DISPATCHER] publish failed for event=%s: %v", event.ID, pubErr)
			_, _ = d.db.ExecContext(ctx, "UPDATE outbox_events SET attempts = $1 WHERE id = $2", attempts, id)
			continue
		}

		if d.failMode != nil && d.failMode.CrashAfterPublish {
			log.Printf("[DISPATCHER] simulated crash after publish for event=%s (not marking as published)", event.ID)
			// Crash happens BEFORE marking as published
			// This simulates: publish succeeded but process crashed before UPDATE
			// On restart, dispatcher sees status='pending' and will re-publish
			return dispatched, ErrProcessCrashed
		}

		now := time.Now()
			_, err = d.db.ExecContext(ctx, "UPDATE outbox_events SET status = 'published', published_at = $1, attempts = $2 WHERE id = $3", now, attempts, id)
		if err != nil {
			return dispatched, fmt.Errorf("mark published: %w", err)
		}

		log.Printf("[DISPATCHER] event=%s marked as published", event.ID)
		dispatched++
	}

	return dispatched, nil
}

// ============================================================================
// 9. Idempotent Consumer Implementation (FIXED: consumer_name + event_id)
// ============================================================================

// processedEventsMutex is a package-level mutex to ensure atomic check-and-insert
// for idempotency across multiple worker instances. This simulates database-level
// unique constraint enforcement that would normally be handled by the DB.
var processedEventsMutex sync.Mutex

type CommissionWorker struct {
	db              *sql.DB
	commissionsPaid int64
	mu              sync.Mutex
}

func NewCommissionWorker(db *sql.DB) *CommissionWorker {
	return &CommissionWorker{db: db}
}

// HandleEvent processes an event idempotently using consumer_name + event_id
// Uses a package-level mutex to ensure atomicity of the check-and-insert pattern.
// In a real database with UNIQUE(consumer_name, event_id), this would be atomic.
func (c *CommissionWorker) HandleEvent(ctx context.Context, consumerName string, event Event) (bool, error) {
	log.Printf("[CONSUMER:%s] processing %s", consumerName, event.ID)

	// Use package-level mutex to ensure atomicity of check-and-insert
	processedEventsMutex.Lock()
	defer processedEventsMutex.Unlock()

	// Check if already processed by this consumer
	// We query by event_id and check consumer_name in code since mockdb's WHERE parser is limited
	rows, err := c.db.QueryContext(ctx, "SELECT consumer_name FROM processed_events WHERE event_id = $1", event.ID)
	if err != nil {
		return false, fmt.Errorf("query processed_events: %w", err)
	}

	alreadyProcessed := false
	for rows.Next() {
		var existingConsumer string
		if err := rows.Scan(&existingConsumer); err == nil {
			if existingConsumer == consumerName {
				alreadyProcessed = true
				break
			}
		}
	}
	rows.Close()

	if alreadyProcessed {
		log.Printf("[CONSUMER:%s] duplicate %s skipped", consumerName, event.ID)
		return false, nil // Idempotent: skip duplicate
	}

	_, err = c.db.ExecContext(ctx, "INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3)", consumerName, event.ID, time.Now())
	if err != nil {
		return false, fmt.Errorf("insert processed_events: %w", err)
	}

	c.mu.Lock()
	c.commissionsPaid++
	c.mu.Unlock()

	log.Printf("[CONSUMER:%s] completed %s", consumerName, event.ID)
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
	Event    Event
	Reason   string
	Attempts int
	FailedAt time.Time
}

type DeadLetterQueue struct {
	mu      sync.Mutex
	records []DLQRecord
}

func NewDeadLetterQueue() *DeadLetterQueue {
	return &DeadLetterQueue{}
}

func (dlq *DeadLetterQueue) Add(event Event, reason string, attempts int) {
	dlq.mu.Lock()
	defer dlq.mu.Unlock()
	dlq.records = append(dlq.records, DLQRecord{
		Event:    event,
		Reason:   reason,
		Attempts: attempts,
		FailedAt: time.Now(),
	})
	log.Printf("[DLQ] %s moved after %d attempts", event.ID, attempts)
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

// ============================================================================
// 10. Saga Pattern & Compensating Transactions (FIXED: correct compensation order)
// ============================================================================

type SagaStep struct {
	Action     func(ctx context.Context) error
	Compensate func(ctx context.Context) error
}

type SagaExecutionError struct {
	OriginalError       error
	CompensationErrors  []error
}

func (e *SagaExecutionError) Error() string {
	return fmt.Sprintf("saga failed: original=%v, compensation_errors=%v", e.OriginalError, e.CompensationErrors)
}

type Saga struct {
	steps          []SagaStep
	executed       []int
	compensated    []int
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
			// Continue compensating even if a compensation fails (PROMPT 11)
			var compErrors []error
			// Find the index of the last successfully executed step
			for j := len(s.executed) - 1; j >= 0; j-- {
				if compErr := s.steps[s.executed[j]].Compensate(ctx); compErr != nil {
					compErrors = append(compErrors, compErr)
				}
				s.compensated = append(s.compensated, s.executed[j])
			}
			// Note: s.executed remains as-is to show what executed before failure
			_ = s.executed // reference to avoid unused variable warning
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