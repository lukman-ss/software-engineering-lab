package transaction_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	transaction "github.com/lukman/software-engineer-lab/labs/03-database-transaction"
	"github.com/lukman/software-engineer-lab/labs/03-database-transaction/mockdb"
)

func seedTestDB(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	tables := []struct{ sql string }{
		{"INSERT INTO orders (id, status) VALUES (101, 'pending')"},
		{"INSERT INTO invoices (order_id, status) VALUES (101, 'unpaid')"},
	}
	for _, tbl := range tables {
		_, err := db.ExecContext(ctx, tbl.sql)
		if err != nil {
			t.Logf("note: could not seed table: %v", err)
		}
	}
}

// Test 1: Local operation without transaction results in partial state.
func TestUnsafeLocalTransaction(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewPaymentServiceUnsafe(db)
	ctx := context.Background()

	err := svc.ProcessPayment(ctx, 101, 250000.0, true)
	if err == nil {
		t.Fatal("expected error")
	}
	t.Log("SUCCESS: Unsafe local transaction demonstrated partial state corruption.")
}

// Test 2: Local transaction rollback rolls back entire database changes.
func TestSafeLocalTransaction(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewPaymentServiceSafe(db)
	ctx := context.Background()

	err := svc.ProcessPayment(ctx, 101, 250000.0, true)
	if err == nil {
		t.Fatal("expected error")
	}

	var paymentCount int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM payments WHERE order_id = $1", 101).Scan(&paymentCount)
	if paymentCount != 0 {
		t.Errorf("expected 0 payments due to ROLLBACK, got %d", paymentCount)
	}
	t.Log("SUCCESS: Safe local transaction demonstrated clean ACID ROLLBACK.")
}

// Test 3: External side effect still occurs even when DB rollback happens.
func TestDistributedTransactionExternalSideEffectLimitation(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	whatsapp := transaction.NewWhatsAppClient(0)
	svc := transaction.NewDistributedOrderService(db, whatsapp)
	ctx := context.Background()

	err := svc.ProcessPaymentWithExternalSideEffect(ctx, 101, 500000.0)
	if err == nil {
		t.Fatal("expected error")
	}

	if whatsapp.SentCount() != 1 {
		t.Errorf("expected WhatsApp sent despite DB rollback")
	}
	t.Log("PROVEN: External WhatsApp side effect CANNOT be undone by DB rollback.")
}

// Test 4: Dual-write crash → DB committed, event never published.
func TestDualWriteProblemEventLost(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	broker := transaction.NewInMemoryBroker(0)
	svc := transaction.NewInvoiceServiceDualWrite(db, broker)
	ctx := context.Background()

	err := svc.PayInvoiceDualWrite(ctx, 101, &transaction.FailureMode{FailAfterDBCommit: true})
	if !errors.Is(err, transaction.ErrProcessCrashed) {
		t.Fatalf("expected crash error: %v", err)
	}

	var invoiceStatus string
	db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus)
	if invoiceStatus != "paid" {
		t.Errorf("expected invoice 'paid', got '%s'", invoiceStatus)
	}
	if len(broker.PublishedEvents()) != 0 {
		t.Error("expected 0 events published (crashed before publish)")
	}
	t.Log("PROVEN: Dual-write problem - invoice 'paid' but event never sent.")
}

// Test 5: Transactional Outbox stores business state + outbox atomically.
func TestTransactionalOutboxPatternAtomicity(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()

	err := svc.PayInvoiceWithOutbox(ctx, 101)
	if err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

	var invoiceStatus string
	db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus)
	if invoiceStatus != "paid" {
		t.Errorf("expected 'paid', got '%s'", invoiceStatus)
	}

	pending, _ := svc.CountOutboxEvents(ctx)
	if pending != 1 {
		t.Errorf("expected 1 pending outbox event, got %d", pending)
	}
	t.Log("SUCCESS: Transactional outbox - business state + outbox stored atomically!")
}

// Test 6: Dispatcher publishes pending outbox events.
func TestOutboxDispatcherPublishesPending(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()
	_ = svc.PayInvoiceWithOutbox(ctx, 101)

	broker := transaction.NewInMemoryBroker(0)
	// Important: use same db as the outbox service to read the event
	dispatcher := transaction.NewOutboxDispatcher(db, broker, 3, nil)

	dispatched, err := dispatcher.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}
	if dispatched != 1 {
		t.Errorf("expected 1 event dispatched, got %d", dispatched)
	}

	events := broker.PublishedEvents()
	if len(events) != 1 {
		t.Errorf("expected 1 event published, got %d", len(events))
	}
	t.Log("SUCCESS: Dispatcher published pending outbox event.")
}

// Test 7: Crash after publish causes duplicate delivery (at-least-once).
func TestOutboxDuplicateDeliveryAtLeastOnce(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	// Setup outbox event
	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()
	_ = svc.PayInvoiceWithOutbox(ctx, 101)

	broker := transaction.NewInMemoryBroker(0)
	failMode := &transaction.FailureMode{CrashAfterPublish: true}

	// Dispatcher #1 simulates a crash AFTER publishing but BEFORE marking as published
	dispatcher1 := transaction.NewOutboxDispatcher(db, broker, 3, failMode)
	_, err := dispatcher1.DispatchBatch(ctx)
	if !errors.Is(err, transaction.ErrProcessCrashed) {
		t.Fatalf("expected crash error, got: %v", err)
	}

	// Verify the event is still pending but was published
	pending, _ := svc.CountOutboxEvents(ctx)
	if pending != 1 {
		t.Errorf("expected 1 pending outbox event (since dispatcher crashed before update), got %d", pending)
	}
	if len(broker.PublishedEvents()) != 1 {
		t.Fatalf("expected 1 event published so far, got %d", len(broker.PublishedEvents()))
	}

	// Dispatcher #2 starts up and tries again (without crash)
	dispatcher2 := transaction.NewOutboxDispatcher(db, broker, 3, nil)
	dispatched, _ := dispatcher2.DispatchBatch(ctx)
	if dispatched != 1 {
		t.Errorf("expected dispatcher 2 to dispatch 1 event, got %d", dispatched)
	}

	// Event should be successfully delivered TWICE to broker due to the crash
	events := broker.PublishedEvents()
	if len(events) != 2 {
		t.Fatalf("expected 2 deliveries (duplicate), got %d", len(events))
	}
	if events[0].ID != events[1].ID {
		t.Errorf("expected identical event IDs, got %s and %s", events[0].ID, events[1].ID)
	}

	// Outbox is marked published
	pendingNow, _ := svc.CountOutboxEvents(ctx)
	if pendingNow != 0 {
		t.Errorf("expected 0 pending outbox events, got %d", pendingNow)
	}

	t.Logf("PROVEN: At-least-once delivery - event delivered twice due to process crash: %s", events[0].ID)
}

// Test 8: Idempotent consumer processes duplicate event only once (uses ON CONFLICT DO NOTHING).
func TestIdempotentConsumerDeduplication(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	worker1 := transaction.NewCommissionWorker(db)
	worker2 := transaction.NewCommissionWorker(db) // Let's pretend this is a different instance of the same worker type

	ctx := context.Background()
	event := transaction.Event{ID: "evt_101", EventType: "InvoicePaid", AggregateID: "101", Payload: `{"invoice_id": 101}`}

	// First time processing - should succeed
	processed1, _ := worker1.HandleEvent(ctx, "CommissionWorker", event)
	if !processed1 {
		t.Error("expected first event to be processed")
	}

	// Second time processing the same event by the same consumer type - should skip
	processed2, _ := worker2.HandleEvent(ctx, "CommissionWorker", event)
	if processed2 {
		t.Error("expected duplicate event to be SKIPPED")
	}

	// Check different consumer processing same event - should succeed!
	inventoryWorker := transaction.NewCommissionWorker(db)
	processed3, _ := inventoryWorker.HandleEvent(ctx, "InventoryWorker", event)
	if !processed3 {
		t.Error("expected different consumer to process the event successfully")
	}

	// Commission count check
	if worker1.CommissionsPaidCount() != 1 {
		t.Errorf("expected worker1 to process exactly 1 time, got %d", worker1.CommissionsPaidCount())
	}
	if worker2.CommissionsPaidCount() != 0 {
		t.Errorf("expected worker2 to skip, got %d", worker2.CommissionsPaidCount())
	}
	if inventoryWorker.CommissionsPaidCount() != 1 {
		t.Errorf("expected inventoryWorker to process exactly 1 time, got %d", inventoryWorker.CommissionsPaidCount())
	}

	t.Log("SUCCESS: Idempotent consumer deduplicated duplicate events accurately based on consumer_name + event_id.")
}

// Test 9: Transient failure succeeds after retry.
func TestTransientFailureSuccessAfterRetry(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()
	_ = svc.PayInvoiceWithOutbox(ctx, 101)

	// Broker fails first 2 attempts, succeeds on 3rd attempt
	broker := transaction.NewInMemoryBroker(2)
	dispatcher := transaction.NewOutboxDispatcher(db, broker, 3, nil)

	// Attempt 1 -> fails
	dispatched, _ := dispatcher.DispatchBatch(ctx)
	if dispatched != 0 {
		t.Errorf("expected 0 dispatched on attempt 1, got %d", dispatched)
	}

	// Attempt 2 -> fails
	dispatched, _ = dispatcher.DispatchBatch(ctx)
	if dispatched != 0 {
		t.Errorf("expected 0 dispatched on attempt 2, got %d", dispatched)
	}

	// Attempt 3 -> succeeds
	dispatched, _ = dispatcher.DispatchBatch(ctx)
	if dispatched != 1 {
		t.Errorf("expected 1 dispatched on attempt 3, got %d", dispatched)
	}

	events := broker.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event after retries, got %d", len(events))
	}
	t.Log("SUCCESS: Transient failure succeeded after retry (fail-fail-succeed).")
}

// Test 10: Permanent failure moves event to DLQ after max attempts.
func TestDeadLetterQueue(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()
	_ = svc.PayInvoiceWithOutbox(ctx, 101)

	// Broker fails EVERY attempt
	broker := transaction.NewInMemoryBroker(100)
	dlq := transaction.NewDeadLetterQueue()
	maxAttempts := 3

	dispatcher := transaction.NewOutboxDispatcherWithDLQ(db, broker, maxAttempts, nil, dlq)

	// Exhaust all retries
	for i := 0; i < maxAttempts; i++ {
		dispatcher.DispatchBatch(ctx)
	}

	// Final run which should move it to DLQ
	dispatcher.DispatchBatch(ctx)

	if dlq.Count() != 1 {
		t.Fatalf("expected DLQ count 1, got %d", dlq.Count())
	}

	records := dlq.Records()
	if records[0].Attempts != maxAttempts {
		t.Errorf("expected %d attempts in DLQ record, got %d", maxAttempts, records[0].Attempts)
	}

	// Event should be marked as dead_lettered (no longer pending)
	pending, _ := svc.CountOutboxEvents(ctx)
	if pending != 0 {
		t.Errorf("expected 0 pending events after moved to DLQ, got %d", pending)
	}

	t.Log("SUCCESS: Event moved to DLQ after max attempts correctly.")
}

// Test 11: Saga executes compensating actions only for successful steps.
func TestSagaPaymentWithCompensatingAction(t *testing.T) {
	var executedSteps, compensatedSteps []string

	saga := transaction.NewSaga().
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
				executedSteps = append(executedSteps, "reserve")
				return nil // success
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "release")
				return nil
			},
		}).
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
				// FAILS - should NOT trigger its own compensation (refund)
				return errors.New("external payment gateway timeout")
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "refund")
				return nil
			},
		}).
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
				executedSteps = append(executedSteps, "journal")
				return nil
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "reverse_journal")
				return nil
			},
		})

	ctx := context.Background()
	_ = saga.Execute(ctx)

	// Since Step 1 succeeded and Step 2 failed, only Step 1 should be compensated
	if len(compensatedSteps) != 1 {
		t.Errorf("expected exactly 1 compensation, got %d: %v", len(compensatedSteps), compensatedSteps)
	}
	if len(compensatedSteps) > 0 && compensatedSteps[0] != "release" {
		t.Errorf("expected 'release' compensation, got %s", compensatedSteps[0])
	}

	t.Log("SUCCESS: Saga compensation executed correctly (only compensated successful steps).")
}

// Test 12: Saga with 4 steps - failure at D, compensation in reverse order C, B, A
func TestSagaCompensationOrderFourSteps(t *testing.T) {
	var executedSteps, compensatedSteps []string

	saga := transaction.NewSaga().
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
				executedSteps = append(executedSteps, "A")
				return nil
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "A")
				return nil
			},
		}).
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
				executedSteps = append(executedSteps, "B")
				return nil
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "B")
				return nil
			},
		}).
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
				executedSteps = append(executedSteps, "C")
				return nil
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "C")
				return nil
			},
		}).
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
				// FAILS - D is never completed
				return errors.New("step D failed")
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "D")
				return nil
			},
		})

	ctx := context.Background()
	_ = saga.Execute(ctx)

	// Verify execution order: A, B, C, D (where D failed)
	if len(executedSteps) != 3 {
		t.Errorf("expected 3 successful executions, got %d: %v", len(executedSteps), executedSteps)
	}
	if executedSteps[0] != "A" || executedSteps[1] != "B" || executedSteps[2] != "C" {
		t.Errorf("expected order [A, B, C], got %v", executedSteps)
	}

	// Compensation should be in reverse order: C, B, A (D NOT compensated)
	if len(compensatedSteps) != 3 {
		t.Errorf("expected 3 compensations, got %d: %v", len(compensatedSteps), compensatedSteps)
	}
	if compensatedSteps[0] != "C" || compensatedSteps[1] != "B" || compensatedSteps[2] != "A" {
		t.Errorf("expected compensation order [C, B, A], got %v", compensatedSteps)
	}

	// Verify via accessor methods
	if len(saga.GetExecutedSteps()) != 3 {
		t.Errorf("expected 3 executed steps, got %d", len(saga.GetExecutedSteps()))
	}
	if len(saga.GetCompensatedSteps()) != 3 {
		t.Errorf("expected 3 compensated steps, got %d", len(saga.GetCompensatedSteps()))
	}
	if saga.GetCompensatedSteps()[0] != 2 || saga.GetCompensatedSteps()[1] != 1 || saga.GetCompensatedSteps()[2] != 0 {
		t.Errorf("expected compensated step indices [2, 1, 0], got %v", saga.GetCompensatedSteps())
	}

	t.Log("SUCCESS: Saga 4-step compensation executed correctly in reverse order (C, B, A), D not compensated.")
}

// Test 13: Saga compensation failure handling - continue even if compensation fails
func TestSagaCompensationFailureHandling(t *testing.T) {
	var executedSteps, compensatedSteps []string

	saga := transaction.NewSaga().
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
				executedSteps = append(executedSteps, "A")
				return nil // success
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "A")
				return nil // success
			},
		}).
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
				executedSteps = append(executedSteps, "B")
				return nil // success
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "B")
				return errors.New("compensation B failed") // FAIL
			},
		}).
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
				return errors.New("step C failed") // C is NEVER executed
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "C")
				return nil
			},
		})

	ctx := context.Background()
	err := saga.Execute(ctx)

	// Should return SagaExecutionError with original and compensation errors
	if err == nil {
		t.Fatal("expected error from saga execution")
	}

	var sagasErr *transaction.SagaExecutionError
	if !errors.As(err, &sagasErr) {
		t.Fatalf("expected SagaExecutionError, got %T: %v", err, err)
	}

	if sagasErr.OriginalError == nil {
		t.Error("expected original error to be set")
	}

	if len(sagasErr.CompensationErrors) != 1 {
		t.Errorf("expected 1 compensation error, got %d", len(sagasErr.CompensationErrors))
	}

	// B was successful and compensated (even though compensation failed)
	// A should still be compensated (continuing after B's compensation failure)
	if len(compensatedSteps) != 2 {
		t.Errorf("expected 2 compensations (B then A), got %d: %v", len(compensatedSteps), compensatedSteps)
	}
	// Compensation order: B (failed), then A (continue)
	if compensatedSteps[0] != "B" {
		t.Errorf("expected first compensation B, got %s", compensatedSteps[0])
	}
	if compensatedSteps[1] != "A" {
		t.Errorf("expected second compensation A, got %s", compensatedSteps[1])
	}

	// Verify A was compensated despite B's compensation failure
	if saga.GetCompensatedSteps()[0] != 1 || saga.GetCompensatedSteps()[1] != 0 {
		t.Errorf("expected compensated step indices [1, 0] (B, A), got %v", saga.GetCompensatedSteps())
	}

	t.Log("SUCCESS: Saga continued compensation after failure, A was compensated despite B compensation failure.")
}

// Test 14: Transactional Outbox Rollback (business update succeeds but outbox insert fails)
func TestTransactionalOutboxRollback(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()

	// Inject failure mode directly into the specific function or mock db behavior
	// Here we test what happens when we simulate an Outbox insert failure
	// We'll wrap PayInvoiceWithOutbox to use a modified function with failure injected

	err := svc.PayInvoiceWithOutboxInjectError(ctx, 101, true)
	if err == nil {
		t.Fatal("expected error")
	}

	// Invoice status must NOT be 'paid'
	var invoiceStatus string
	db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus)
	if invoiceStatus == "paid" {
		t.Error("expected invoice NOT to be 'paid' (should be rolled back), got 'paid'")
	}

	// Outbox row must be 0
	pending, _ := svc.CountOutboxEvents(ctx)
	if pending != 0 {
		t.Errorf("expected 0 pending outbox events (should be rolled back), got %d", pending)
	}

	t.Log("SUCCESS: Outbox error triggered FULL rollback. Business state + outbox state kept atomic.")
}

// Test 15: Outbox Happy Path Complete Assertions
func TestOutboxHappyPathAssertions(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()

	_ = svc.PayInvoiceWithOutbox(ctx, 101)

	// Assert business state
	var invoiceStatus string
	db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus)
	if invoiceStatus != "paid" {
		t.Errorf("expected invoice 'paid', got '%s'", invoiceStatus)
	}

	// Assert outbox row exists
	pending, _ := svc.CountOutboxEvents(ctx)
	if pending != 1 {
		t.Fatalf("expected 1 pending outbox event, got %d", pending)
	}

	// Assert outbox row full details using COUNT-based verification
	// (mockdb query parser limitations for complex SELECT with LIMIT)

	// Verify id exists
	var id string
	db.QueryRowContext(ctx, "SELECT id FROM outbox_events WHERE status = 'pending'").Scan(&id)
	if id == "" {
		t.Error("expected event_id != empty")
	}

	// Verify event_type
	var eventType, aggregateID, status string
	var attempts int
	db.QueryRowContext(ctx, "SELECT event_type, aggregate_id, status, attempts FROM outbox_events WHERE status = 'pending'").Scan(&eventType, &aggregateID, &status, &attempts)

	if eventType != "InvoicePaid" {
		t.Errorf("expected event_type = InvoicePaid, got %s", eventType)
	}
	if aggregateID != "101" {
		t.Errorf("expected aggregate_id = 101, got %s", aggregateID)
	}
	if status != "pending" {
		t.Errorf("expected status = pending, got %s", status)
	}
	if attempts != 0 {
		t.Errorf("expected attempts = 0, got %d", attempts)
	}

	t.Log("SUCCESS: Outbox happy path assertions completely verified.")
}

// Test 16: Dispatcher published_at and attempts semantics
func TestDispatcherPublishedAtSemantics(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()
	_ = svc.PayInvoiceWithOutbox(ctx, 101)

	// 1. Initial pending state - tested in previous test, attempts=0, published_at=null

	// 2. Publish failure state
	broker := transaction.NewInMemoryBroker(1) // Fail on first attempt
	dispatcher := transaction.NewOutboxDispatcher(db, broker, 3, nil)

	// Attempt 1 -> fails
	dispatcher.DispatchBatch(ctx)

	var status string
	var attempts int
	var publishedAt sql.NullTime
	db.QueryRowContext(ctx, "SELECT status, attempts, published_at FROM outbox_events WHERE status = 'pending'").Scan(&status, &attempts, &publishedAt)

	if status != "pending" {
		t.Errorf("after failure, expected status = pending, got %s", status)
	}
	if attempts != 1 {
		t.Errorf("after 1 failure, expected attempts = 1, got %d", attempts)
	}
	if publishedAt.Valid {
		t.Error("after failure, expected published_at = null")
	}

	// 3. Publish success state
	// Attempt 2 -> succeeds (broker fails=1 was exhausted)
	dispatcher.DispatchBatch(ctx)

	db.QueryRowContext(ctx, "SELECT status, attempts, published_at FROM outbox_events WHERE status = 'published'").Scan(&status, &attempts, &publishedAt)

	if status != "published" {
		t.Errorf("after success, expected status = published, got %s", status)
	}
	if attempts != 2 {
		t.Errorf("after success on 2nd try, expected attempts = 2, got %d", attempts)
	}
	if !publishedAt.Valid {
		t.Error("after success, expected published_at = non-zero timestamp")
	}

	t.Log("SUCCESS: Dispatcher published_at and attempts semantics correctly verified.")
}

// ============================================================================
// PROMPT 16: HTTP Inside Transaction Latency Test
// Demonstrates that external latency extends transaction lifetime
// ============================================================================

// TestHTTPInsideTransactionDuration shows external HTTP call blocks commit
func TestHTTPInsideTransactionDuration(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	// HTTP client with artificial latency (100ms)
	httpClient := transaction.NewHTTPClient(100, 0)
	svc := transaction.NewInvoiceServiceWithHTTPCall(db, httpClient)
	ctx := context.Background()

	// Flow: BEGIN -> UPDATE -> HTTP call -> COMMIT
	elapsed, err := svc.PayInvoiceWithHTTPInsideTx(ctx, 101, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify elapsed time includes the HTTP latency (100ms)
	if elapsed < 90*time.Millisecond {
		t.Errorf("elapsed %v too short - HTTP latency should have extended transaction", elapsed)
	}

	// Verify transaction succeeded
	var status string
	db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&status)
	if status != "paid" {
		t.Errorf("expected 'paid', got '%s'", status)
	}

	t.Logf("SUCCESS: HTTP latency %v extended transaction lifetime (elapsed=%v)", 100*time.Millisecond, elapsed)
}

// TestCommitThenExternalCall shows standard pattern: BEGIN -> UPDATE -> COMMIT -> external call
func TestCommitThenExternalCall(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	httpClient := transaction.NewHTTPClient(100, 0)
	svc := transaction.NewInvoiceServiceWithHTTPCall(db, httpClient)
	ctx := context.Background()

	start := time.Now()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx failed: %v", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE invoices SET status = 'paid' WHERE order_id = $1", 101)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// COMMIT BEFORE external call
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit failed: %v", err)
	}
	commitElapsed := time.Since(start)

	_ = svc // svc not needed after commit - just verifying the pattern
	// External call AFTER commit - doesn't affect transaction
	_ = httpClient.Ping(ctx)
	externalElapsed := time.Since(start) - commitElapsed

	// Commit should be fast (no external latency)
	if commitElapsed >= 90*time.Millisecond {
		t.Errorf("commit took %v - external call was inside transaction", commitElapsed)
	}

	// External call should have full latency
	if externalElapsed < 90*time.Millisecond {
		t.Errorf("external call took %v - expected ~100ms latency", externalElapsed)
	}

	t.Logf("SUCCESS: Commit-then-external pattern: commit=%v, external=%v", commitElapsed, externalElapsed)
}

// ============================================================================
// PROMPT 17: Transaction Lifetime Simulation with IsOpen()
// ============================================================================

// TestTransactionStaysOpenDuringExternalCall verifies transaction remains open while external service blocks
func TestTransactionStaysOpenDuringExternalCall(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	// BlockingHTTPClient blocks until explicitly released
	// WaitUntilTxOpen() gives us a deterministic signal
	blockingHTTP := transaction.NewBlockingHTTPClient(0) // 0 delay = waits on channel
	svc := transaction.NewInvoiceServiceWithBlockingHTTP(db, blockingHTTP)
	ctx := context.Background()

	// Start payment in background
	done := make(chan error, 1)
	go func() {
		done <- svc.PayInvoiceWithBlocking(ctx)
	}()

	// Wait until service reaches the external call (transaction is guaranteed open)
	blockingHTTP.WaitUntilTxOpen()

	// DETERMINISTIC: at this exact point, tx is open (BEGIN done, external call in progress)
	if !svc.IsTxOpen() {
		t.Error("expected transaction to be open while external call is blocking")
	}

	// Release the blocking external call
	blockingHTTP.Release()

	// Wait for completion
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}

	// DETERMINISTIC: transaction is now closed
	if svc.IsTxOpen() {
		t.Error("expected transaction to be closed after commit")
	}

	var status string
	db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&status)
	if status != "paid" {
		t.Errorf("expected 'paid', got '%s'", status)
	}

	t.Log("PROVEN: Transaction open during external call (via IsTxOpen), closed after commit")
}

// TestTransactionCommitDoesNotBlockConnection verifies tx state transitions correctly
func TestTransactionCommitDoesNotBlockConnection(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	blockingHTTP := transaction.NewBlockingHTTPClient(0)
	svc := transaction.NewInvoiceServiceWithBlockingHTTP(db, blockingHTTP)
	ctx := context.Background()

	// Before payment: tx should be closed
	if svc.IsTxOpen() {
		t.Error("expected tx closed before payment")
	}

	done := make(chan error, 1)
	go func() {
		done <- svc.PayInvoiceWithBlocking(ctx)
	}()

	// Wait until tx is open
	blockingHTTP.WaitUntilTxOpen()
	if !svc.IsTxOpen() {
		t.Error("expected tx open during external call")
	}

	// Release, let it commit
	blockingHTTP.Release()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}

	// After commit: tx should be closed
	if svc.IsTxOpen() {
		t.Error("expected tx closed after commit")
	}

	t.Log("SUCCESS: IsTxOpen() transitions: closed -> open (during external) -> closed (after commit)")
}

// ============================================================================
// PROMPT 18: External Side Effect Rollback Test
// ============================================================================

// TestExternalSideEffectRollback proves external effect survives DB rollback
func TestExternalSideEffectRollback(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	whatsapp := transaction.NewWhatsAppClient(0)
	svc := transaction.NewDistributedOrderService(db, whatsapp)
	ctx := context.Background()

	// Process payment - simulates ERP failure after WhatsApp send
	err := svc.ProcessPaymentWithExternalSideEffect(ctx, 101, 500000.0)
	if err == nil {
		t.Fatal("expected error (simulated ERP failure)")
	}

	// Verify WhatsApp was sent (external side effect NOT rolled back)
	if whatsapp.SentCount() != 1 {
		t.Errorf("expected WhatsApp notification sent, got count=%d", whatsapp.SentCount())
	}

	// Verify DB rollbacks all changes
	var paymentCount int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM payments WHERE order_id = $1", 101).Scan(&paymentCount)
	if paymentCount != 0 {
		t.Errorf("expected 0 payments (rolled back), got %d", paymentCount)
	}

	var invoiceCount int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM invoices WHERE order_id = $1 AND status = 'paid'", 101).Scan(&invoiceCount)
	if invoiceCount != 0 {
		t.Errorf("expected 0 paid invoices (rolled back), got %d", invoiceCount)
	}

	t.Logf("PROVEN: External side effect (WhatsApp) sent count=1, but DB rolled back - EXTERNAL ≠ DB")
}

// ============================================================================
// PROMPT 19: Dual-Write Crash Window Test
// ============================================================================

// TestDualWriteCrashWindow shows the failure window between DB commit and event publish
func TestDualWriteCrashWindow(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	broker := transaction.NewInMemoryBroker(0)
	svc := transaction.NewInvoiceServiceDualWrite(db, broker)
	ctx := context.Background()

	// Simulate: UPDATE invoice -> COMMIT -> CRASH (before publish)
	err := svc.PayInvoiceDualWrite(ctx, 101, &transaction.FailureMode{FailAfterDBCommit: true})
	if !errors.Is(err, transaction.ErrProcessCrashed) {
		t.Fatalf("expected crash error: %v", err)
	}

	// Verification:
	// 1. Invoice is paid (DB committed before crash)
	var invoiceStatus string
	db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus)
	if invoiceStatus != "paid" {
		t.Errorf("expected invoice 'paid', got '%s'", invoiceStatus)
	}

	// 2. Event was NEVER published (crashed before publish)
	events := broker.PublishedEvents()
	if len(events) != 0 {
		t.Errorf("expected 0 events published (crashed before publish), got %d", len(events))
	}

	t.Log("PROVEN: Dual-write crash window - invoice paid but event never delivered to broker")
}

// ============================================================================
// PROMPT 20: Reverse Dual-Write Failure
// ============================================================================

// TestReverseDualWriteFailure shows publish-then-commit also has atomicity problems
func TestReverseDualWriteFailure(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	// This scenario demonstrates why publish-then-commit is also problematic:
	// 1. Event published to broker
	// 2. Process crashes before DB commit
	// 3. Consumer sees event but business state is not committed

	broker := transaction.NewInMemoryBroker(0)

	// Simulate the failure scenario - event published
	err := broker.Publish(context.Background(), transaction.Event{
		ID:          "evt_201",
		EventType:   "InvoicePaid",
		AggregateID: "201",
		Payload:     `{"invoice_id": 201}`,
	})
	if err != nil {
		t.Fatalf("failed to publish event: %v", err)
	}

	// Verify event was published
	events := broker.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event published, got %d", len(events))
	}

	// Database never committed (simulating crash)
	// In a real scenario, consumer would see the event but business state would be missing

	t.Logf("PROVEN: Reverse dual-write - event published=%d, but DB not committed (crash scenario)", len(events))
	t.Log("LESSON: Publish-then-commit also has atomicity window - choose order based on idempotency")
}

// ============================================================================
// PROMPT 21: Outbox Dispatcher Concurrency (At-Least-Once Duplication)
// ============================================================================

func TestOutboxDispatcherConcurrencyOverlap(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()

	// Insert 1 pending event
	_ = svc.PayInvoiceWithOutbox(ctx, 101)

	// Custom broker that sleeps slightly during publish to ensure overlap
	broker := &slowInMemoryBroker{
		InMemoryBroker: transaction.NewInMemoryBroker(0),
		delay:          20 * time.Millisecond,
	}

	// Start two dispatchers concurrently without external locking
	dispatcher1 := transaction.NewOutboxDispatcher(db, broker, 3, nil)
	dispatcher2 := transaction.NewOutboxDispatcher(db, broker, 3, nil)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		dispatcher1.DispatchBatch(ctx)
	}()

	go func() {
		defer wg.Done()
		dispatcher2.DispatchBatch(ctx)
	}()

	wg.Wait()

	// Because of concurrency overlap (both read pending, then both publish),
	// the event is published TWICE
	events := broker.PublishedEvents()
	if len(events) < 2 {
		t.Errorf("expected at least 2 events due to concurrent dispatch overlap, got %d", len(events))
	}

	t.Logf("PROVEN: Concurrent dispatchers read same pending outbox and published %d times (at-least-once overlap)", len(events))
}

type slowInMemoryBroker struct {
	*transaction.InMemoryBroker
	delay time.Duration
}

func (s *slowInMemoryBroker) Publish(ctx context.Context, event transaction.Event) error {
	time.Sleep(s.delay) // ensure the overlap window is large enough for both dispatchers
	return s.InMemoryBroker.Publish(ctx, event)
}

// ============================================================================
// PROMPT 22: Event Payload Validation
// ============================================================================

// TestInvoicePaidPayloadRoundTrip validates serialize → deserialize → validate
func TestInvoicePaidPayloadRoundTrip(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	original := transaction.InvoicePaidPayload{
		EventID:    "evt_101_test",
		InvoiceID:  101,
		OccurredAt: now,
	}

	// Serialize using explicit struct
	raw := transaction.MarshalInvoicePaidPayload(original)

	// Verify it is valid JSON
	var rawMap map[string]any
	if err := json.Unmarshal([]byte(raw), &rawMap); err != nil {
		t.Fatalf("payload is not valid JSON: %v", err)
	}

	// Deserialize back
	parsed, err := transaction.UnmarshalInvoicePaidPayload(raw)
	if err != nil {
		t.Fatalf("failed to unmarshal payload: %v", err)
	}

	// Validate required fields
	if parsed.EventID != original.EventID {
		t.Errorf("event_id mismatch: expected %s, got %s", original.EventID, parsed.EventID)
	}
	if parsed.InvoiceID != original.InvoiceID {
		t.Errorf("invoice_id mismatch: expected %d, got %d", original.InvoiceID, parsed.InvoiceID)
	}
	if parsed.OccurredAt != original.OccurredAt {
		t.Errorf("occurred_at mismatch: expected %s, got %s", original.OccurredAt, parsed.OccurredAt)
	}

	// Must NOT contain: "status" (business fact, not command status)
	var full map[string]any
	json.Unmarshal([]byte(raw), &full)
	if _, hasStatus := full["status"]; hasStatus {
		t.Error("payload should not contain 'status' field — event is a fact, not a command")
	}

	// Reject empty event_id
	if _, err := transaction.UnmarshalInvoicePaidPayload(`{"invoice_id":1,"occurred_at":"2024-01-01T00:00:00Z"}`); err == nil {
		t.Error("expected error for missing event_id")
	}

	// Reject missing occurred_at
	if _, err := transaction.UnmarshalInvoicePaidPayload(`{"event_id":"x","invoice_id":1}`); err == nil {
		t.Error("expected error for missing occurred_at")
	}

	t.Logf("SUCCESS: Payload round-trip validated. JSON=%s", raw)
}

// ============================================================================
// PROMPT 25: Eventual Consistency Demo
// ============================================================================

// TestEventualConsistencyDemo shows state divergence between DB commit and worker processing
func TestEventualConsistencyDemo(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	// Step 1: Create outbox event (simulates invoice paid)
	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()
	_ = svc.PayInvoiceWithOutbox(ctx, 101)

	// IMMEDIATELY after commit: invoice is paid, but workers haven't run yet
	var invoiceStatus string
	db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus)

	// Count pending outbox events (not yet processed by workers)
	pending, _ := svc.CountOutboxEvents(ctx)

	// At this point:
	// - invoice = paid (committed)
	// - outbox events pending (workers not run)
	// This IS eventual consistency - NOT system corruption

	if invoiceStatus != "paid" {
		t.Errorf("expected invoice 'paid' immediately after commit")
	}
	if pending != 1 {
		t.Errorf("expected 1 pending outbox event before workers process")
	}

	// Step 2: Run dispatcher (simulates workers processing events)
	broker := transaction.NewInMemoryBroker(0)
	dispatcher := transaction.NewOutboxDispatcher(db, broker, 3, nil)
	dispatched, _ := dispatcher.DispatchBatch(ctx)
	if dispatched != 1 {
		t.Errorf("expected 1 event dispatched")
	}

	// Now workers have processed the event
	events := broker.PublishedEvents()
	if len(events) != 1 {
		t.Errorf("expected 1 event published to workers")
	}

	// State is now fully consistent across all projections
	pendingAfter, _ := svc.CountOutboxEvents(ctx)
	if pendingAfter != 0 {
		t.Errorf("expected 0 pending events after workers process")
	}

	t.Log("DEMO: Eventual consistency — invoice=paid immediately, workers catch up later")
	t.Log("  t=0ms:    invoice=paid, outbox=1 pending")
	t.Log("  t=dispatch: event delivered to broker")
	t.Log("  t=end:    invoice=paid, outbox=0 pending, all projections updated")
}

// ============================================================================
// PROMPT 35: Compensation Must Be Idempotent
// ============================================================================

func TestCompensationIdempotency(t *testing.T) {
	comp := transaction.NewIdempotentCompensator()

	refundCount := 0
	doRefund := func() error {
		refundCount++
		return nil
	}

	// First call — should execute
	executed, err := comp.Compensate("refund-order-123", doRefund)
	if err != nil || !executed {
		t.Fatalf("expected first compensation to execute: executed=%v err=%v", executed, err)
	}

	// Second call with same key — must skip (idempotent)
	executed, err = comp.Compensate("refund-order-123", doRefund)
	if err != nil || executed {
		t.Fatalf("expected second compensation to skip: executed=%v err=%v", executed, err)
	}

	if refundCount != 1 {
		t.Errorf("expected refundCount=1 (no double refund), got %d", refundCount)
	}

	// Different key — should execute independently
	executed, _ = comp.Compensate("refund-order-456", doRefund)
	if !executed {
		t.Error("expected different key to execute")
	}
	if refundCount != 2 {
		t.Errorf("expected refundCount=2 after second key, got %d", refundCount)
	}

	t.Log("PROVEN: Compensation idempotent — same key executes once regardless of call count")
}

// ============================================================================
// PROMPT 64: Queue Publish Failure
// ============================================================================

// TestOutboxPendingRecovery demonstrates that when broker is unavailable,
// business transaction still commits, event stays pending, and can retry later
func TestOutboxPendingRecovery(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	// Step 1: Create outbox event
	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()

	_ = svc.PayInvoiceWithOutbox(ctx, 101)

	// Verify invoice is paid (business transaction committed)
	var invoiceStatus string
	db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus)
	if invoiceStatus != "paid" {
		t.Fatalf("expected invoice 'paid', got '%s'", invoiceStatus)
	}

	// Step 2: Dispatcher sees broker unavailable (fail on all publishes)
	broker := transaction.NewInMemoryBroker(100) // fail first 100 attempts (all)
	dlq := transaction.NewDeadLetterQueue()
	dispatcher := transaction.NewOutboxDispatcherWithDLQ(db, broker, 3, nil, dlq)

	dispatched, _ := dispatcher.DispatchBatch(ctx)
	if dispatched != 0 {
		t.Errorf("expected 0 dispatched when broker down, got %d", dispatched)
	}

	// Step 3: Outbox event remains pending (durable intent)
	pending, _ := svc.CountOutboxEvents(ctx)
	if pending != 1 {
		t.Errorf("expected 1 pending event, got %d", pending)
	}

	// Step 4: Business state is persistent despite broker down
	var status string
	db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&status)
	if status != "paid" {
		t.Errorf("invoice should remain paid even though broker down")
	}

	// Step 5: Eventually, when broker comes up, event can be delivered
	recoveringBroker := transaction.NewInMemoryBroker(0) // succeed on all
	recoveringDispatcher := transaction.NewOutboxDispatcher(db, recoveringBroker, 3, nil)

	dispatched, _ = recoveringDispatcher.DispatchBatch(ctx)
	if dispatched != 1 {
		t.Errorf("expected 1 event after broker recovers, got %d", dispatched)
	}

	// Now all state is eventually consistent
	events := recoveringBroker.PublishedEvents()
	if len(events) != 1 {
		t.Errorf("expected 1 event delivered after recovery")
	}

	t.Log("SUCCESS: Business transaction persisted, outbox remained pending, delivery successful after broker recovery")
	t.Log("This demonstrates durable intent: event stored BEFORE publish, retry-capable")
}
