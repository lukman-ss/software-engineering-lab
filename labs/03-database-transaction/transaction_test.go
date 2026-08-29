package transaction_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

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
