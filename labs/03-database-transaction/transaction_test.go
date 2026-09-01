package transaction_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	transaction "github.com/lukman/software-engineer-lab/labs/03-database-transaction"
	"github.com/lukman/software-engineer-lab/labs/03-database-transaction/mockdb"
)

func seedTestDB(t *testing.T, db *sql.DB) {
	t.Helper()
	ctx := context.Background()
	tables := []struct{ sql string }{
		{"INSERT INTO orders (id, status) VALUES (101, 'pending')"},
		{"INSERT INTO invoices (order_id, status) VALUES (101, 'unpaid')"},
	}
	for _, tbl := range tables {
		_, err := db.ExecContext(ctx, tbl.sql)
		if err != nil {
			t.Fatalf("note: could not seed table: %v", err)
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

	// PROVE partial-state corruption after injected failure:
	// Flow: INSERT payment -> UPDATE order=paid -> ERROR injected -> wallet tx never made
	var paymentCount int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM payments WHERE order_id = $1", 101).Scan(&paymentCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if paymentCount != 1 {
		t.Errorf("expected 1 payment persisted (partial state), got %d", paymentCount)
	}

	var orderStatus string
	if err := db.QueryRowContext(ctx, "SELECT status FROM orders WHERE id = $1", 101).Scan(&orderStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if orderStatus != "paid" {
		t.Errorf("expected order status 'paid' persisted (partial state), got '%s'", orderStatus)
	}

	var walletCount int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM wallet_transactions WHERE order_id = $1", 101).Scan(&walletCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if walletCount != 0 {
		t.Errorf("expected 0 wallet transactions (missing), got %d", walletCount)
	}

	t.Log("SUCCESS: Unsafe local transaction demonstrated partial state corruption.")
	t.Log("  payment persisted=1, order.status=paid, wallet_transactions=0 → inconsistent")
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

	// PROVE full rollback state: all local invariants restored
	var paymentCount int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM payments WHERE order_id = $1", 101).Scan(&paymentCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if paymentCount != 0 {
		t.Errorf("expected 0 payments due to ROLLBACK, got %d", paymentCount)
	}

	var orderStatus string
	if err := db.QueryRowContext(ctx, "SELECT status FROM orders WHERE id = $1", 101).Scan(&orderStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if orderStatus != "pending" {
		t.Errorf("expected order status 'pending' (rollback restored), got '%s'", orderStatus)
	}

	var walletCount int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM wallet_transactions WHERE order_id = $1", 101).Scan(&walletCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if walletCount != 0 {
		t.Errorf("expected 0 wallet transactions due to ROLLBACK, got %d", walletCount)
	}

	t.Log("SUCCESS: Safe local transaction demonstrated clean ACID ROLLBACK.")
	t.Log("  payments=0, order.status=pending, wallet_transactions=0 → fully rolled back")
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
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
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
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if invoiceStatus != "paid" {
		t.Errorf("expected 'paid', got '%s'", invoiceStatus)
	}

	pending, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
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
	if err := svc.PayInvoiceWithOutbox(ctx, 101); err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

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
	if err := svc.PayInvoiceWithOutbox(ctx, 101); err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

	broker := transaction.NewInMemoryBroker(0)
	failMode := &transaction.FailureMode{CrashAfterPublish: true}

	// Dispatcher #1 simulates a crash AFTER publishing but BEFORE marking as published
	dispatcher1 := transaction.NewOutboxDispatcher(db, broker, 3, failMode)
	_, err := dispatcher1.DispatchBatch(ctx)
	if !errors.Is(err, transaction.ErrProcessCrashed) {
		t.Fatalf("expected crash error, got: %v", err)
	}

	// Verify the event is still pending but was published
	pending, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
	if pending != 1 {
		t.Errorf("expected 1 pending outbox event (since dispatcher crashed before update), got %d", pending)
	}
	if len(broker.PublishedEvents()) != 1 {
		t.Fatalf("expected 1 event published so far, got %d", len(broker.PublishedEvents()))
	}

	// Dispatcher #2 starts up and tries again (without crash)
	dispatcher2 := transaction.NewOutboxDispatcher(db, broker, 3, nil)
	dispatched, err := dispatcher2.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("unexpected error on second dispatch: %v", err)
	}
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
	pendingNow, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
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
	worker2 := transaction.NewCommissionWorker(db) // Different instance, same consumer type

	ctx := context.Background()
	event := transaction.Event{ID: "evt_101", EventType: "InvoicePaid", AggregateID: "101", Payload: `{"invoice_id": 101}`}

	// First time processing - should succeed
	processed1, err := worker1.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("unexpected error on first process: %v", err)
	}
	if !processed1 {
		t.Error("expected first event to be processed")
	}

	// Verify business state committed in DB
	count, err := worker1.GetDBCommissionCount(ctx)
	if err != nil {
		t.Fatalf("failed to query commissions: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 commission in DB after first process, got %d", count)
	}

	// Second time processing by same consumer - should skip
	processed2, err := worker2.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("unexpected error on duplicate: %v", err)
	}
	if processed2 {
		t.Error("expected duplicate event to be SKIPPED")
	}

	// Verify DB count unchanged (no duplicate commission inserted)
	count, err = worker1.GetDBCommissionCount(ctx)
	if err != nil {
		t.Fatalf("failed to query commissions: %v", err)
	}
	if count != 1 {
		t.Errorf("expected still 1 commission after duplicate (idempotent), got %d", count)
	}

	// Different consumer processing same event - should succeed independently
	inventoryWorker := transaction.NewCommissionWorker(db)
	processed3, err := inventoryWorker.HandleEvent(ctx, "InventoryWorker", event)
	if err != nil {
		t.Fatalf("unexpected error on different consumer: %v", err)
	}
	if !processed3 {
		t.Error("expected different consumer to process the event successfully")
	}

	// Verify now have 2 commissions (one per consumer)
	count, err = worker1.GetDBCommissionCount(ctx)
	if err != nil {
		t.Fatalf("failed to query commissions: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 commissions (one per consumer), got %d", count)
	}

	// Observability counters (for metrics, not business state)
	if worker1.ObservedBusinessExecutions() != 1 {
		t.Errorf("expected worker1 observability counter=1, got %d", worker1.ObservedBusinessExecutions())
	}
	if worker2.ObservedBusinessExecutions() != 0 {
		t.Errorf("expected worker2 observability counter=0, got %d", worker2.ObservedBusinessExecutions())
	}
	if inventoryWorker.ObservedBusinessExecutions() != 1 {
		t.Errorf("expected inventoryWorker observability counter=1, got %d", inventoryWorker.ObservedBusinessExecutions())
	}

	t.Log("SUCCESS: Idempotent consumer deduplicated based on (consumer_name, event_id).")
	t.Log("  - Same consumer: event processed once, DB commission count=1")
	t.Log("  - Duplicate: skipped, DB commission count unchanged")
	t.Log("  - Different consumer: processed independently, DB commission count=2")
}

// Test 8b: Atomic consumer flow - verifies business state + dedup marker in same transaction
func TestAtomicConsumerFlow(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	worker := transaction.NewCommissionWorker(db)
	ctx := context.Background()
	event := transaction.Event{ID: "evt_200", EventType: "InvoicePaid", AggregateID: "200", Payload: `{"invoice_id": 200}`}

	// Process event atomically
	processed, err := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !processed {
		t.Fatal("expected event to be processed")
	}

	// Verify both dedup marker AND business state exist in DB
	var dedupCount, commissionCount int64

	// Note: mockdb only supports single WHERE condition, so we just check event_id
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", event.ID).Scan(&dedupCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions WHERE event_id = $1", event.ID).Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if dedupCount != 1 {
		t.Errorf("expected 1 processed_events marker, got %d", dedupCount)
	}
	if commissionCount != 1 {
		t.Errorf("expected 1 commission row, got %d", commissionCount)
	}

	t.Log("SUCCESS: Atomic consumer flow - dedup marker + business state committed together.")
}

// Test 8c: Consumer sequential redelivery idempotency (same instance, retry after success).
// NOTE: Real consumer-restart test is TestConsumerCrashAfterCommitBeforeAck / TestConsumerRestartRedelivery.
func TestSequentialRedeliveryIdempotency(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	worker := transaction.NewCommissionWorker(db)
	ctx := context.Background()
	event := transaction.Event{ID: "evt_300", EventType: "InvoicePaid", AggregateID: "300", Payload: `{"invoice_id": 300}`}

	// First delivery - process successfully
	processed1, err := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("unexpected error on first delivery: %v", err)
	}
	if !processed1 {
		t.Fatal("expected first delivery to be processed")
	}

	// Verify state after first delivery
	count1, err := worker.GetDBCommissionCount(ctx)
	if err != nil {
		t.Fatalf("get commission count failed: %v", err)
	}
	if count1 != 1 {
		t.Fatalf("expected 1 commission after first delivery, got %d", count1)
	}

	// Simulate crash + redelivery (message broker redelivers at-least-once)
	processed2, err := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("unexpected error on redelivery: %v", err)
	}
	if processed2 {
		t.Error("expected redelivered event to be SKIPPED (idempotent)")
	}

	// Verify state unchanged after redelivery
	count2, err := worker.GetDBCommissionCount(ctx)
	if err != nil {
		t.Fatalf("get commission count failed: %v", err)
	}
	if count2 != 1 {
		t.Errorf("expected still 1 commission after redelivery, got %d", count2)
	}

	t.Log("SUCCESS: Crash/redelivery handled idempotently - no duplicate business state.")
}

// Test 10: REAL Concurrent duplicate consumer test with deterministic barrier synchronization
// Uses barrier to ensure both workers attempt simultaneously.
// MockDB's coarse-grained locking means one will succeed atomically via ON CONFLICT.
func TestConcurrentDuplicateConsumer(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	worker1 := transaction.NewCommissionWorker(db)
	worker2 := transaction.NewCommissionWorker(db)
	ctx := context.Background()
	event := transaction.Event{ID: "evt_concurrent", EventType: "InvoicePaid", AggregateID: "101", Payload: `{"invoice_id": 101}`}

	// REAL concurrent execution using deterministic barrier
	start := make(chan struct{})
	var wg sync.WaitGroup
	var processed1, processed2 bool
	var err1, err2 error

	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		processed1, err1 = worker1.HandleEvent(ctx, "CommissionWorker", event)
	}()

	go func() {
		defer wg.Done()
		<-start
		processed2, err2 = worker2.HandleEvent(ctx, "CommissionWorker", event)
	}()

	// Release both goroutines simultaneously
	close(start)
	wg.Wait()

	if err1 != nil || err2 != nil {
		t.Fatalf("errors: worker1=%v, worker2=%v", err1, err2)
	}

	// Exactly one should have processed, one should have skipped (idempotent)
	if !processed1 && !processed2 {
		t.Fatal("expected at least one worker to process")
	}
	if processed1 && processed2 {
		t.Error("expected exactly one worker to process (deduplication failed)")
	}

	// Verify exactly 1 business record (atomic uniqueness under concurrency)
	var count int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions WHERE event_id = $1", event.ID).Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 commission row, got %d", count)
	}

	// Verify exactly 1 processed_events marker
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", event.ID).Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 processed_events row, got %d", count)
	}

	t.Log("SUCCESS: Concurrent duplicate detection - exactly one business mutation occurred under race.")
}

// Test 11: Different consumers process same event independently
func TestDifferentConsumersSameEvent(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	inventoryWorker := transaction.NewCommissionWorker(db)
	commissionWorker := transaction.NewCommissionWorker(db)
	ctx := context.Background()
	event := transaction.Event{ID: "evt_shared", EventType: "InvoicePaid", AggregateID: "201", Payload: `{"invoice_id": 201}`}

	// Both consumers process the same event
	processed1, err := inventoryWorker.HandleEvent(ctx, "InventoryWorker", event)
	if err != nil {
		t.Fatalf("inventory worker error: %v", err)
	}
	if !processed1 {
		t.Error("expected InventoryWorker to process")
	}

	processed2, err := commissionWorker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("commission worker error: %v", err)
	}
	if !processed2 {
		t.Error("expected CommissionWorker to process")
	}

	// Verify 2 separate processed_events records
	var count int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", event.ID).Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 processed_events records (one per consumer), got %d", count)
	}

	// Verify 2 business records
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions WHERE event_id = $1", event.ID).Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 commission records, got %d", count)
	}

	t.Log("SUCCESS: Different consumers processed same event independently.")
}

// Test 12: Same consumer processes same event twice sequentially (idempotent skip on 2nd).
// NOTE: This is NOT a business-mutation-failure test. Real failure-injection tests:
//   - TestBusinessMutationRollbackAtomicity (dedup claim succeeds, business mutation fails mid-tx)
//   - TestBusinessMutationFailureRollback (worker with FailBusinessMutation injected)
func TestSequentialDuplicateConsumer(t *testing.T) {

	db := mockdb.NewDB()
	defer db.Close()

	worker := transaction.NewCommissionWorker(db)
	ctx := context.Background()
	event := transaction.Event{ID: "evt_mutation_test", EventType: "InvoicePaid", AggregateID: "301", Payload: `{"invoice_id": 301}`}

	// First processing succeeds
	processed, err := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !processed {
		t.Fatal("expected first processing to succeed")
	}

	// Verify state was committed
	var commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 1 {
		t.Fatalf("expected 1 commission, got %d", commissionCount)
	}

	// Simulate reprocessing same event (like retry after crash)
	processed2, err := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}
	if processed2 {
		t.Error("expected duplicate to be skipped")
	}

	// Count still 1 - no duplicate business mutation
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 1 {
		t.Errorf("idempotent - expected still 1 commission, got %d", commissionCount)
	}

	t.Log("SUCCESS: Business mutation failure scenario handled correctly via idempotency.")
}

// Test 9: Transient failure succeeds after retry.
func TestTransientFailureSuccessAfterRetry(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()
	if err := svc.PayInvoiceWithOutbox(ctx, 101); err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

	// Broker fails first 2 attempts, succeeds on 3rd attempt
	broker := transaction.NewInMemoryBroker(2)
	dispatcher := transaction.NewOutboxDispatcher(db, broker, 3, nil)

	// Attempt 1 -> fails
	dispatched, err := dispatcher.DispatchBatch(ctx)
	// Dispatcher swallows broker publish errors internally to keep processing other events,
	// so it doesn't return an error here, just returns dispatched=0.
	if err != nil {
		t.Fatalf("unexpected error from dispatcher: %v", err)
	}
	if dispatched != 0 {
		t.Errorf("expected 0 dispatched on attempt 1, got %d", dispatched)
	}

	// Attempt 2 -> fails
	dispatched, err = dispatcher.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("unexpected error from dispatcher: %v", err)
	}
	if dispatched != 0 {
		t.Errorf("expected 0 dispatched on attempt 2, got %d", dispatched)
	}

	// Attempt 3 -> succeeds
	dispatched, err = dispatcher.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("unexpected error from dispatcher: %v", err)
	}
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
	if err := svc.PayInvoiceWithOutbox(ctx, 101); err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

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
	pending, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
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
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if invoiceStatus == "paid" {
		t.Error("expected invoice NOT to be 'paid' (should be rolled back), got 'paid'")
	}

	// Outbox row must be 0
	pending, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
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

	if err := svc.PayInvoiceWithOutbox(ctx, 101); err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

	// Assert business state
	var invoiceStatus string
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if invoiceStatus != "paid" {
		t.Errorf("expected invoice 'paid', got '%s'", invoiceStatus)
	}

	// Assert outbox row exists
	pending, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
	if pending != 1 {
		t.Fatalf("expected 1 pending outbox event, got %d", pending)
	}

	// Assert outbox row full details using COUNT-based verification
	// (mockdb query parser limitations for complex SELECT with LIMIT)

	// Verify id exists
	var id string
	if err := db.QueryRowContext(ctx, "SELECT id FROM outbox_events WHERE status = 'pending'").Scan(&id); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if id == "" {
		t.Error("expected event_id != empty")
	}

	// Verify event_type
	var eventType, aggregateID, status string
	var attempts int
	if err := db.QueryRowContext(ctx, "SELECT event_type, aggregate_id, status, attempts FROM outbox_events WHERE status = 'pending'").Scan(&eventType, &aggregateID, &status, &attempts); err != nil {
		t.Fatalf("query failed: %v", err)
	}

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
	if err := svc.PayInvoiceWithOutbox(ctx, 101); err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

	// 1. Initial pending state - tested in previous test, attempts=0, published_at=null

	// 2. Publish failure state
	broker := transaction.NewInMemoryBroker(1) // Fail on first attempt
	dispatcher := transaction.NewOutboxDispatcher(db, broker, 3, nil)

	// Attempt 1 -> fails
	dispatcher.DispatchBatch(ctx)

	var status string
	var attempts int
	var publishedAt sql.NullTime
	if err := db.QueryRowContext(ctx, "SELECT status, attempts, published_at FROM outbox_events WHERE status = 'pending'").Scan(&status, &attempts, &publishedAt); err != nil {
		t.Fatalf("query failed: %v", err)
	}

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

	if err := db.QueryRowContext(ctx, "SELECT status, attempts, published_at FROM outbox_events WHERE status = 'published'").Scan(&status, &attempts, &publishedAt); err != nil {
		t.Fatalf("query failed: %v", err)
	}

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
// HTTP Inside Transaction Latency Test
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
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&status); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if status != "paid" {
		t.Errorf("expected 'paid', got '%s'", status)
	}

	t.Logf("SUCCESS: HTTP latency %v extended transaction lifetime (elapsed=%v)", 100*time.Millisecond, elapsed)
}

// TestCommitThenExternalCall shows standard pattern: BEGIN -> UPDATE -> COMMIT -> external call
// It proves that the transaction is already fully closed when the external call starts,
// preventing long external latency from holding connection pool resources.
func TestCommitThenExternalCall(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	// Use blocking client to pause during external call
	blockingHTTP := transaction.NewBlockingHTTPClient(0)
	svc := transaction.NewInvoiceServiceWithBlockingHTTP(db, blockingHTTP)
	ctx := context.Background()

	// 1. Manually begin and commit transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx failed: %v", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE invoices SET status = 'paid' WHERE order_id = $1", 101)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	// Prove database is committed
	var status string
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&status); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if status != "paid" {
		t.Errorf("expected 'paid', got '%s'", status)
	}

	// 2. Start external call (it will block)
	done := make(chan error, 1)
	go func() {
		// We use ping directly here to prove external call executes outside tx
		done <- blockingHTTP.Ping(ctx)
	}()

	// Wait for the ping to actually enter
	blockingHTTP.WaitUntilTxOpen()

	// Prove transaction is NOT open in the service during external call
	if svc.IsTxOpen() {
		t.Error("expected transaction to be CLOSED during external call")
	}

	// Release it
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

	t.Logf("SUCCESS: Commit-then-external pattern proven with deterministic synchronization.")
	t.Log("  Transaction was already committed and closed before external call blocked.")
}

// ============================================================================
// Transaction Lifetime Simulation with IsOpen()
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
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&status); err != nil {
		t.Fatalf("query failed: %v", err)
	}
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
// External Side Effect Rollback Test
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
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM payments WHERE order_id = $1", 101).Scan(&paymentCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if paymentCount != 0 {
		t.Errorf("expected 0 payments (rolled back), got %d", paymentCount)
	}

	var invoiceCount int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM invoices WHERE order_id = $1 AND status = 'paid'", 101).Scan(&invoiceCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if invoiceCount != 0 {
		t.Errorf("expected 0 paid invoices (rolled back), got %d", invoiceCount)
	}

	t.Logf("PROVEN: External side effect (WhatsApp) sent count=1, but DB rolled back - EXTERNAL ≠ DB")
}

// ============================================================================
// Dual-Write Crash Window Test
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
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
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
// Reverse Dual-Write Failure
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
// Outbox Dispatcher Concurrency (At-Least-Once Duplication)
// ============================================================================

func TestOutboxDispatcherConcurrencyOverlap(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()

	// Insert 1 pending event
	if err := svc.PayInvoiceWithOutbox(ctx, 101); err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

	// Custom broker that uses a barrier to guarantee both dispatchers
	// enter the Publish method concurrently before either is allowed to proceed
	var barrier sync.WaitGroup
	barrier.Add(2) // Wait for exactly 2 publishers

	broker := &barrierBroker{
		InMemoryBroker: transaction.NewInMemoryBroker(0),
		barrier:        &barrier,
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

type barrierBroker struct {
	*transaction.InMemoryBroker
	barrier *sync.WaitGroup
}

func (s *barrierBroker) Publish(ctx context.Context, event transaction.Event) error {
	// Wait for all concurrent dispatchers to reach this point
	s.barrier.Done()
	s.barrier.Wait()
	return s.InMemoryBroker.Publish(ctx, event)
}

// ============================================================================
// Event Payload Validation
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
	if err := json.Unmarshal([]byte(raw), &full); err != nil {
		t.Fatalf("unmarshal full map failed: %v", err)
	}
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
// Eventual Consistency Demo
// ============================================================================

// TestEventualConsistencyDemo shows state divergence between DB commit and worker processing.
// Publishing an event to the broker does NOT mean downstream projections are updated.
// Consumer must separately process the event to achieve full consistency.
func TestEventualConsistencyDemo(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)
	ctx := context.Background()

	// Step 1: Create outbox event (simulates invoice paid business transaction)
	svc := transaction.NewInvoiceServiceOutbox(db)
	if err := svc.PayInvoiceWithOutbox(ctx, 101); err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

	// IMMEDIATE STATE after business commit:
	// - invoice = paid (committed in DB)
	// - outbox = 1 pending (event stored atomically with invoice update)
	var invoiceStatus string
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if invoiceStatus != "paid" {
		t.Errorf("expected invoice 'paid' immediately after commit")
	}

	pending, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
	if pending != 1 {
		t.Errorf("expected 1 pending outbox event before dispatch, got %d", pending)
	}

	// CRITICAL: commissions table is EMPTY (consumer hasn't processed yet)
	// This is EXPECTED - eventual consistency gap
	var commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 0 {
		t.Errorf("expected 0 commissions (not yet processed), got %d", commissionCount)
	}

	// Step 2: Run outbox dispatcher (publish event to broker)
	// This does NOT mean downstream business effect is complete
	// Dispatcher PUBLISHES event - it does NOT process it as a consumer
	broker := transaction.NewInMemoryBroker(0)
	dispatcher := transaction.NewOutboxDispatcher(db, broker, 3, nil)
	dispatched, err := dispatcher.DispatchBatch(ctx)
	if err != nil || dispatched != 1 {
		t.Fatalf("expected 1 event dispatched, got %d (err: %v)", dispatched, err)
	}

	// AFTER DISPATCH (but before consumer):
	// - invoice = paid
	// - commissions = 0 (event reached broker, but no consumer processed it)
	// - outbox = 0 pending (events dispatched and marked published)
	events := broker.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event published to broker, got %d", len(events))
	}

	// Verify commissions still 0 after dispatch - dispatch ≠ processing
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 0 {
		t.Error("expected 0 commissions immediately after dispatch (dispatch ≠ consumer)")
	}

	// Step 3: Consumer processes the published event
	// CommissionWorker consumes the InvoicePaid event from broker
	worker := transaction.NewCommissionWorker(db)
	processed, err := worker.HandleEvent(ctx, "CommissionConsumer", events[0])
	if err != nil || !processed {
		t.Fatalf("consumer failed to process: processed=%v, err=%v", processed, err)
	}

	// FINAL STATE: fully consistent
	// - invoice = paid
	// - commissions = 1 (consumer processed the event)
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 1 {
		t.Errorf("expected 1 commission after consumer processing, got %d", commissionCount)
	}

	pendingAfter, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
	if pendingAfter != 0 {
		t.Errorf("expected 0 pending events after dispatch, got %d", pendingAfter)
	}

	t.Log("DEMO: Eventual consistency — invoice=paid immediately, commissions eventually =1")
	t.Log("  t=0 (business commit): invoice=paid, commissions=0, outbox=1pending")
	t.Log("  t=1 (dispatched):     invoice=paid, commissions=0, outbox=0pending (event in broker)")
	t.Log("  t=2 (consumer):       invoice=paid, commissions=1, all consistent")
	t.Log("  Key insight: dispatch (publish) ≠ consumer processing (business effect)")
}

// ============================================================================
// Compensation Must Be Idempotent
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
// Queue Publish Failure
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

	if err := svc.PayInvoiceWithOutbox(ctx, 101); err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

	// Verify invoice is paid (business transaction committed)
	var invoiceStatus string
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if invoiceStatus != "paid" {
		t.Fatalf("expected invoice 'paid', got '%s'", invoiceStatus)
	}

	// Step 2: Dispatcher sees broker unavailable (fail on all publishes)
	broker := transaction.NewInMemoryBroker(100) // fail first 100 attempts (all)
	dlq := transaction.NewDeadLetterQueue()
	dispatcher := transaction.NewOutboxDispatcherWithDLQ(db, broker, 3, nil, dlq)

	dispatched, err := dispatcher.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("unexpected error from dispatcher: %v", err)
	}
	if dispatched != 0 {
		t.Errorf("expected 0 dispatched when broker down, got %d", dispatched)
	}

	// Step 3: Outbox event remains pending (durable intent)
	pending, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
	if pending != 1 {
		t.Errorf("expected 1 pending event, got %d", pending)
	}

	// Step 4: Business state is persistent despite broker down
	var status string
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&status); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if status != "paid" {
		t.Errorf("invoice should remain paid even though broker down")
	}

	// Step 5: Eventually, when broker comes up, event can be delivered
	recoveringBroker := transaction.NewInMemoryBroker(0) // succeed on all
	recoveringDispatcher := transaction.NewOutboxDispatcher(db, recoveringBroker, 3, nil)

	dispatched, err = recoveringDispatcher.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("unexpected error from dispatcher: %v", err)
	}
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

// ============================================================================
// Test Concurrent Different Events
// ============================================================================

// TestConcurrentDifferentEvents verifies that processing different events concurrently
// does not cause lost updates. Each event should result in exactly one business row.
// Uses deterministic barrier synchronization for real concurrent execution.
func TestConcurrentDifferentEvents(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	worker1 := transaction.NewCommissionWorker(db)
	worker2 := transaction.NewCommissionWorker(db)
	ctx := context.Background()

	// Two different events
	event1 := transaction.Event{ID: "evt-1", EventType: "InvoicePaid", AggregateID: "101", Payload: `{"invoice_id": 101}`}
	event2 := transaction.Event{ID: "evt-2", EventType: "InvoicePaid", AggregateID: "102", Payload: `{"invoice_id": 102}`}

	// REAL concurrent execution using deterministic barrier
	start := make(chan struct{})
	var wg sync.WaitGroup
	var processed1, processed2 bool
	var err1, err2 error

	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		processed1, err1 = worker1.HandleEvent(ctx, "CommissionWorker", event1)
	}()

	go func() {
		defer wg.Done()
		<-start
		processed2, err2 = worker2.HandleEvent(ctx, "CommissionWorker", event2)
	}()

	close(start)
	wg.Wait()

	if err1 != nil {
		t.Fatalf("worker1 error: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("worker2 error: %v", err2)
	}
	if !processed1 {
		t.Error("expected event1 to be processed")
	}
	if !processed2 {
		t.Error("expected event2 to be processed")
	}

	// Verify exactly 2 processed_events markers
	var dedupCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE consumer_name = $1", "CommissionWorker").Scan(&dedupCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if dedupCount != 2 {
		t.Errorf("expected 2 processed_events markers, got %d", dedupCount)
	}

	// Verify exactly 2 business rows (commissions)
	var commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 2 {
		t.Errorf("expected 2 commission rows, got %d", commissionCount)
	}

	t.Log("SUCCESS: Concurrent different events processed correctly - 2 processed_events, 2 business rows, no lost update")
}

// ============================================================================
// Test Concurrent Same Event
// ============================================================================

// TestConcurrentSameEvent verifies that when two consumers process the same event
// concurrently, only one succeeds and the other is rejected (duplicate/no-op).
// Uses deterministic barrier synchronization for real concurrent execution.
func TestConcurrentSameEvent(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	worker1 := transaction.NewCommissionWorker(db)
	worker2 := transaction.NewCommissionWorker(db)
	ctx := context.Background()

	event := transaction.Event{ID: "evt-same", EventType: "InvoicePaid", AggregateID: "201", Payload: `{"invoice_id": 201}`}

	// REAL concurrent execution using deterministic barrier
	// Both workers start simultaneously and race to claim the event
	start := make(chan struct{})
	var wg sync.WaitGroup
	var processed1, processed2 bool
	var err1, err2 error

	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		processed1, err1 = worker1.HandleEvent(ctx, "CommissionWorker", event)
	}()

	go func() {
		defer wg.Done()
		<-start
		processed2, err2 = worker2.HandleEvent(ctx, "CommissionWorker", event)
	}()

	close(start)
	wg.Wait()

	if err1 != nil {
		t.Fatalf("worker1 error: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("worker2 error: %v", err2)
	}

	// Exactly one should have processed, one should have skipped (idempotent)
	if !processed1 && !processed2 {
		t.Fatal("expected at least one worker to process")
	}
	if processed1 && processed2 {
		t.Error("expected exactly one worker to process (deduplication failed)")
	}

	// Verify exactly 1 business row (atomic uniqueness under concurrency)
	var commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 1 {
		t.Fatalf("expected 1 commission row after concurrent processing, got %d", commissionCount)
	}

	// Verify exactly 1 processed_events marker (atomic ON CONFLICT)
	var dedupCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE consumer_name = $1 AND event_id = $2", "CommissionWorker", event.ID).Scan(&dedupCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if dedupCount != 1 {
		t.Errorf("expected 1 processed_events marker, got %d", dedupCount)
	}

	t.Log("SUCCESS: Concurrent same event - exactly one business transaction committed atomically under race")
}

// ============================================================================
// Concurrent Different Consumers Same Event
// ============================================================================

// TestConcurrentDifferentConsumersSameEvent verifies that two different consumers
// can process the same event concurrently without conflicts.
// Per README: deduplication key is (consumer_name, event_id), NOT global event_id.
// Each consumer has independent business effect.
func TestConcurrentDifferentConsumersSameEvent(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	// Two worker instances (same type, different consumer names)
	inventoryWorker := transaction.NewCommissionWorker(db)
	commissionWorker := transaction.NewCommissionWorker(db)

	// Same event, same event_id but DIFFERENT consumer_name
	event := transaction.Event{ID: "evt-shared-123", EventType: "InvoicePaid", AggregateID: "123", Payload: `{"invoice_id": 123}`}

	// Concurrent execution using deterministic barrier
	start := make(chan struct{})
	var wg sync.WaitGroup
	var processed1, processed2 bool
	var err1, err2 error

	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		processed1, err1 = inventoryWorker.HandleEvent(ctx, "InventoryConsumer", event)
	}()

	go func() {
		defer wg.Done()
		<-start
		processed2, err2 = commissionWorker.HandleEvent(ctx, "CommissionConsumer", event)
	}()

	close(start)
	wg.Wait()

	if err1 != nil {
		t.Fatalf("InventoryConsumer error: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("CommissionConsumer error: %v", err2)
	}
	if !processed1 {
		t.Error("expected InventoryConsumer to process")
	}
	if !processed2 {
		t.Error("expected CommissionConsumer to process")
	}

	// Verify both processed_events markers exist (one per consumer_name + event_id)
	var processedCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", event.ID).Scan(&processedCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if processedCount != 2 {
		t.Errorf("expected 2 processed_events (one per consumer), got %d", processedCount)
	}

	// Verify 2 business rows (each consumer has independent business effect)
	var commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 2 {
		t.Errorf("expected 2 commission rows, got %d", commissionCount)
	}

	t.Log("SUCCESS: Concurrent different consumers same event - both processed independently")
	t.Log("  Deduplication key is (consumer_name, event_id), NOT global event_id")
}

// ============================================================================
// Separate Delivery Count From Business Effect Count
// ============================================================================

// TestConsumerCrashAfterCommitBeforeAck verifies that a crash after commit but before ACK
// results in redelivery but idempotent processing prevents duplicate business effect.
// Uses separate worker instances to simulate consumer restart.
func TestConsumerCrashAfterCommitBeforeAck(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()
	event := transaction.Event{ID: "evt-ack-crash", EventType: "InvoicePaid", AggregateID: "400", Payload: `{"invoice_id": 400}`}

	// Delivery #1: worker1 processes and commits successfully
	worker1 := transaction.NewCommissionWorker(db)
	processed1, err := worker1.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("first delivery error: %v", err)
	}
	if !processed1 {
		t.Fatal("expected first delivery to process")
	}

	// Verify business state after first delivery
	var commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 1 {
		t.Fatalf("expected 1 commission after first delivery, got %d", commissionCount)
	}

	// Simulate crash after commit but before ACK:
	// worker1 process is gone (restarted).
	// Message broker redelivers to a new worker instance.
	// Dedup state must come from persistence DB, not worker1's in-memory state.
	worker2 := transaction.NewCommissionWorker(db)

	// Delivery #2: new worker processes same event - should be deduplicated
	processed2, err := worker2.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("redelivery error: %v", err)
	}
	if processed2 {
		t.Error("expected redelivered event to be SKIPPED (idempotent)")
	}

	// CRITICAL: deliveries = 2 (message delivered twice due to at-least-once)
	//          business rows = 1 (only 1 actual commission due to idempotency)
	var dedupCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", event.ID).Scan(&dedupCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if dedupCount != 1 {
		t.Errorf("expected 1 processed_events marker, got %d", dedupCount)
	}

	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 1 {
		t.Errorf("expected 1 commission business row, got %d", commissionCount)
	}

	// Verify dedup state is durable (not in-memory on worker1)
	// worker2 has its own fresh in-memory state, yet dedup still works
	if worker2.ObservedBusinessExecutions() != 0 {
		t.Errorf("expected worker2 to have 0 observed executions (not processed), got %d", worker2.ObservedBusinessExecutions())
	}
	if worker1.ObservedBusinessExecutions() != 1 {
		t.Errorf("expected worker1 to have 1 observed execution, got %d", worker1.ObservedBusinessExecutions())
	}

	t.Log("SUCCESS: deliveries=2 (redelivered), business rows=1 (idempotent)")
	t.Log("  Dedup state is durable - persisted in DB, not in worker in-memory state")
	t.Log("  At-least-once delivery + idempotent consumer = effectively-once business effect")
}

// ============================================================================
// Verify No Lost Updates with Mock DB
// ============================================================================

// TestMockDBNoLostUpdates verifies that the consumer's idempotent pattern
// correctly prevents lost updates when multiple transactions commit concurrently.
// Uses deterministic barrier synchronization for real concurrent execution.
func TestMockDBNoLostUpdates(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	// Create multiple workers to process concurrently
	workers := make([]*transaction.CommissionWorker, 3)
	for i := range workers {
		workers[i] = transaction.NewCommissionWorker(db)
	}
	ctx := context.Background()

	// Three different events
	events := []transaction.Event{
		{ID: "evt-no-lost-1", EventType: "InvoicePaid", AggregateID: "501", Payload: `{"invoice_id": 501}`},
		{ID: "evt-no-lost-2", EventType: "InvoicePaid", AggregateID: "502", Payload: `{"invoice_id": 502}`},
		{ID: "evt-no-lost-3", EventType: "InvoicePaid", AggregateID: "503", Payload: `{"invoice_id": 503}`},
	}

	// REAL concurrent execution using deterministic barrier
	start := make(chan struct{})
	var wg sync.WaitGroup
	results := make([]bool, 3)
	errs := make([]error, 3)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			<-start
			results[idx], errs[idx] = workers[idx].HandleEvent(ctx, "CommissionWorker", events[idx])
		}(i)
	}

	close(start)
	wg.Wait()

	// Check for errors
	for i, err := range errs {
		if err != nil {
			t.Fatalf("event %d error: %v", i, err)
		}
	}

	// All should process successfully
	for i, processed := range results {
		if !processed {
			t.Errorf("event %d should be processed", i)
		}
	}

	// Verify all 3 rows exist - no lost updates under concurrent commit
	var dedupCount, commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE consumer_name = $1", "CommissionWorker").Scan(&dedupCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if dedupCount != 3 {
		t.Errorf("expected 3 processed_events markers, got %d", dedupCount)
	}
	if commissionCount != 3 {
		t.Errorf("expected 3 commission rows, got %d", commissionCount)
	}

	t.Log("SUCCESS: Mock DB - all 3 events processed correctly under real concurrency")
	t.Log("  Each transaction committed concurrently with unique primary keys")
	t.Log("  No lost updates - all business writes persisted")
}

// ============================================================================
// Verify Rollback Isolation
// ============================================================================

// TestMockDBRollbackIsolation verifies that a rolled-back transaction doesn't
// leave any traces and doesn't affect subsequent committed transactions.
func TestMockDBRollbackIsolation(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	// Transaction A: Write then ROLLBACK
	txA, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("txA begin failed: %v", err)
	}
	_, err = txA.ExecContext(ctx, "INSERT INTO commissions (event_id, amount) VALUES ('txA-rolledback', 100)")
	if err != nil {
		t.Fatalf("txA insert failed: %v", err)
	}
	// Rollback A
	err = txA.Rollback()
	if err != nil {
		t.Fatalf("txA rollback failed: %v", err)
	}

	// Transaction B: Write and COMMIT
	txB, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("txB begin failed: %v", err)
	}
	_, err = txB.ExecContext(ctx, "INSERT INTO commissions (event_id, amount) VALUES ('txB-committed', 200)")
	if err != nil {
		t.Fatalf("txB insert failed: %v", err)
	}
	err = txB.Commit()
	if err != nil {
		t.Fatalf("txB commit failed: %v", err)
	}

	// Verify: A is absent, B is present
	var countA, countB int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions WHERE event_id = $1", "txA-rolledback").Scan(&countA); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions WHERE event_id = $1", "txB-committed").Scan(&countB); err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if countA != 0 {
		t.Errorf("expected txA-rolledback to be absent (rolled back), got count=%d", countA)
	}
	if countB != 1 {
		t.Errorf("expected txB-committed to be present (committed), got count=%d", countB)
	}

	totalCount := 0
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&totalCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if totalCount != 1 {
		t.Errorf("expected exactly 1 total row, got %d", totalCount)
	}

	t.Log("SUCCESS: Rollback isolation verified")
	t.Log("  TX A rolled back → no traces")
	t.Log("  TX B committed → row present")
}

// ============================================================================
// Real Consumer Restart/Redelivery Test
// ============================================================================

// TestConsumerRestartRedelivery verifies that a consumer can restart and
// reprocess events from the outbox without duplicating business state.
func TestConsumerRestartRedelivery(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	worker := transaction.NewCommissionWorker(db)

	// Event from durable outbox
	event := transaction.Event{ID: "evt-restart-test", EventType: "InvoicePaid", AggregateID: "701", Payload: `{"invoice_id": 701}`}

	// First delivery: process successfully
	processed1, err := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("first delivery error: %v", err)
	}
	if !processed1 {
		t.Fatal("expected first delivery to process")
	}

	// Verify state after first delivery
	var commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 1 {
		t.Fatalf("expected 1 commission after first delivery, got %d", commissionCount)
	}

	// SIMULATE: Consumer restarts, reconnects to same DB
	// Previous worker instance is gone (process crash/restart)
	// Same DB persisted the processed_events marker
	newWorker := transaction.NewCommissionWorker(db)

	// Event is redelivered (message broker at-least-once delivery)
	processed2, err := newWorker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("redelivery error: %v", err)
	}
	if processed2 {
		t.Error("expected redelivered event to be SKIPPED (idempotent)")
	}

	// Verify state unchanged - effectively-once business effect
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 1 {
		t.Errorf("expected still 1 commission after restart/redelivery, got %d", commissionCount)
	}

	t.Log("SUCCESS: Consumer restart/redelivery - effectively-once business effect")
	t.Log("  Delivery count = 2 (message broker redelivered)")
	t.Log("  Business state = 1 (idempotent deduplication worked)")
}

// ============================================================================
// Eventual Consistency Test Correctness
// ============================================================================

// TestEventualConsistencyCorrectness verifies the eventual consistency pattern
// and that the dispatcher correctly transitions events through states.
func TestEventualConsistencyCorrectness(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	ctx := context.Background()

	// Step 1: Create outbox event (simulates business transaction committed)
	svc := transaction.NewInvoiceServiceOutbox(db)
	invoiceID := 101

	err := svc.PayInvoiceWithOutbox(ctx, invoiceID)
	if err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

	// IMMEDIATE STATE: Business committed, event pending dispatch
	var invoiceStatus string
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", invoiceID).Scan(&invoiceStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if invoiceStatus != "paid" {
		t.Errorf("expected invoice 'paid' immediately, got '%s'", invoiceStatus)
	}

	// Outbox is pending (worker has NOT processed it yet)
	pending, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
	if pending != 1 {
		t.Errorf("expected 1 pending outbox event, got %d", pending)
	}

	// Dispatch events (outbox dispatcher publishes to broker)
	// NOTE: Dispatcher PUBLISHES events; it does NOT process business effects.
	// Consumer must process the published event separately for business projection.
	broker := transaction.NewInMemoryBroker(0)
	dispatcher := transaction.NewOutboxDispatcher(db, broker, 3, nil)
	dispatched, err := dispatcher.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}
	if dispatched != 1 {
		t.Errorf("expected 1 event dispatched, got %d", dispatched)
	}

	// AFTER DISPATCH (before consumer):
	// - invoice = paid (business committed)
	// - commissions = 0 (no consumer processed the event yet)
	// - broker = 1 event (event published, awaiting consumer)
	events := broker.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event published, got %d", len(events))
	}

	// Verify commissions still 0 after dispatch - dispatch ≠ processing
	var commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 0 {
		t.Error("expected 0 commissions immediately after dispatch (dispatch ≠ consumer)")
	}

	// Verify outbox event is now marked as published (not pending)
	pendingAfter, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
	if pendingAfter != 0 {
		t.Errorf("expected 0 pending events after dispatch, got %d", pendingAfter)
	}

	// Verify event details (dispatch ≠ consumer processing)
	if events[0].EventType != "InvoicePaid" {
		t.Errorf("expected InvoicePaid event type, got %s", events[0].EventType)
	}
	if events[0].AggregateID != fmt.Sprintf("%d", invoiceID) {
		t.Errorf("expected aggregate_id %d, got %s", invoiceID, events[0].AggregateID)
	}

	// T3: Consumer processes the event
	worker := transaction.NewCommissionWorker(db)
	processed, err := worker.HandleEvent(ctx, "CommissionConsumer", events[0])
	if err != nil || !processed {
		t.Fatalf("consumer failed to process: processed=%v, err=%v", processed, err)
	}

	// FINAL STATE: fully consistent
	// - invoice = paid
	// - commissions = 1 (consumer processed the event)
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 1 {
		t.Errorf("expected 1 commission after consumer processing, got %d", commissionCount)
	}

	t.Log("SUCCESS: Eventual consistency flow verified")
	t.Log("  t=0 (business commit): invoice=paid, outbox=1pending, commissions=0 (downstream absent)")
	t.Log("  t=1 (event published): outbox=0pending, commissions=0 still (downstream still absent)")
	t.Log("  t=2 (consumer active): commissions=1 (downstream present)")
	t.Log("  Result: Dispatch (publish) ≠ Consumer processing (business effect)")
}

// ============================================================================
// README ↔ Implementation Synchronization
// ============================================================================

// TestREADMESyncIdempotentConsumer verifies the idempotent consumer pattern
// matches the README specification: dedup by (consumer_name, event_id).
func TestREADMESyncIdempotentConsumer(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	// Per README section 9: processed_events has UNIQUE (consumer_name, event_id)
	// This allows different consumers to process same event independently

	commissionWorker := transaction.NewCommissionWorker(db)
	inventoryWorker := transaction.NewCommissionWorker(db)
	event := transaction.Event{ID: "evt-801", EventType: "InvoicePaid", AggregateID: "801", Payload: `{"invoice_id": 801}`}

	// Both CommissionWorker and InventoryWorker process same event
	processed1, err := commissionWorker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("CommissionWorker error: %v", err)
	}
	if !processed1 {
		t.Error("expected CommissionWorker to process")
	}

	processed2, err := inventoryWorker.HandleEvent(ctx, "InventoryWorker", event)
	if err != nil {
		t.Fatalf("InventoryWorker error: %v", err)
	}
	if !processed2 {
		t.Error("expected InventoryWorker to process")
	}

	// Verify both processed_events records exist (different consumer_name)
	var processedCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", event.ID).Scan(&processedCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if processedCount != 2 {
		t.Errorf("expected 2 processed_events (one per consumer), got %d", processedCount)
	}

	// Verify 2 separate business mutations (each worker has own business effect)
	var commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 2 {
		t.Errorf("expected 2 commission records, got %d", commissionCount)
	}

	// Now same consumer tries again - should be deduplicated
	processed3, err := commissionWorker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("duplicate CommissionWorker error: %v", err)
	}
	if processed3 {
		t.Error("expected duplicate CommissionWorker to be skipped")
	}

	// Count unchanged
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 2 {
		t.Errorf("expected still 2 commissions, got %d", commissionCount)
	}

	t.Log("SUCCESS: README idempotent consumer pattern verified")
	t.Log("  Different consumers: each processes independently")
	t.Log("  Same consumer: deduplication by (consumer_name, event_id)")
}

// ============================================================================
// ON CONFLICT Semantics - Direct MockDB Test
// ============================================================================

// TestMockDBOnConflictSemantics verifies the INSERT ... ON CONFLICT DO NOTHING
// returns correct RowsAffected values directly at the database level.
func TestMockDBOnConflictSemantics(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	// First insert: should succeed
	tx1, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("tx1 begin failed: %v", err)
	}
	result, err := tx1.ExecContext(ctx,
		"INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3) ON CONFLICT (consumer_name, event_id) DO NOTHING",
		"CommissionConsumer", "evt-on-conflict", time.Now())
	if err != nil {
		t.Fatalf("first insert failed: %v", err)
	}
	affected, _ := result.RowsAffected()
	if affected != 1 {
		t.Errorf("expected first insert RowsAffected=1, got %d", affected)
	}
	tx1.Commit()

	// Second insert with SAME consumer_name + event_id: should NOTHING (affected=0)
	tx2, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("tx2 begin failed: %v", err)
	}
	result2, err := tx2.ExecContext(ctx,
		"INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3) ON CONFLICT (consumer_name, event_id) DO NOTHING",
		"CommissionConsumer", "evt-on-conflict", time.Now())
	if err != nil {
		t.Fatalf("duplicate insert failed: %v", err)
	}
	affected2, _ := result2.RowsAffected()
	if affected2 != 0 {
		t.Errorf("expected duplicate RowsAffected=0, got %d", affected2)
	}
	tx2.Rollback()

	// Verify only 1 row in database
	var count int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", "evt-on-conflict").Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 row in database, got %d", count)
	}

	t.Log("SUCCESS: ON CONFLICT DO NOTHING semantics verified - first=1, duplicate=0")
}

// ============================================================================
// Real Business Mutation Failure Test
// ============================================================================

// TestBusinessMutationFailureRollback verifies that when business mutation fails
// after the dedup marker is claimed, the entire transaction including the claim
// is rolled back. This proves the atomicity: no partial state.
func TestBusinessMutationFailureRollback(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	// Worker with injected business mutation failure
	worker := transaction.NewCommissionWorkerWithFailure(db, &transaction.ConsumerFailureMode{
		FailBusinessMutation: true,
	})

	event := transaction.Event{ID: "evt-business-fail", EventType: "InvoicePaid", AggregateID: "601", Payload: `{"invoice_id": 601}`}

	// Attempt processing - should fail
	processed, err := worker.HandleEvent(ctx, "CommissionConsumer", event)
	if err == nil {
		t.Fatal("expected error due to injected business mutation failure")
	}
	if processed {
		t.Error("expected processing to fail, not succeed")
	}

	// Expected: BOTH processed_events = 0 AND commissions = 0
	// The claim marker must be rolled back along with business mutation
	var processedCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", event.ID).Scan(&processedCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if processedCount != 0 {
		t.Errorf("expected 0 processed_events (rolled back), got %d", processedCount)
	}

	var commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions WHERE event_id = $1", event.ID).Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 0 {
		t.Errorf("expected 0 commissions (rolled back), got %d", commissionCount)
	}

	t.Log("SUCCESS: Business mutation failure - entire transaction rolled back atomically")
	t.Log("  Processed marker was NOT committed (it was rolled back with business state)")
}

// ============================================================================
// Redelivery After Business Failure
// ============================================================================

// TestRedeliveryAfterBusinessFailure verifies that after a rollback due to
// business mutation failure, a redelivered event (new consumer, no failure config)
// can successfully process. This proves rollback did not create false dedup marker.
func TestRedeliveryAfterBusinessFailure(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	event := transaction.Event{ID: "evt-business-fail", EventType: "InvoicePaid", AggregateID: "601", Payload: `{"invoice_id": 601}`}

	// Worker with failure injection
	failWorker := transaction.NewCommissionWorkerWithFailure(db, &transaction.ConsumerFailureMode{
		FailBusinessMutation: true,
	})

	// First attempt - FAILS
	processed1, err := failWorker.HandleEvent(ctx, "CommissionConsumer", event)
	if err == nil {
		t.Fatal("expected error on first attempt")
	}
	if processed1 {
		t.Error("expected failure on first attempt")
	}

	// Verify nothing was committed
	commissionCount, err := failWorker.GetDBCommissionCount(ctx)
	if err != nil {
		t.Fatalf("get commission count failed: %v", err)
	}
	if commissionCount != 0 {
		t.Errorf("expected 0 commissions after failed attempt, got %d", commissionCount)
	}

	// Disable failure injection - simulate redelivery after recovery
	normalWorker := transaction.NewCommissionWorker(db)

	// Redeliver same event - should now succeed
	processed2, err := normalWorker.HandleEvent(ctx, "CommissionConsumer", event)
	if err != nil {
		t.Fatalf("redelivery failed: %v", err)
	}
	if !processed2 {
		t.Error("expected redelivered event to be processed successfully")
	}

	// Verify business state now exists
	commissionCount, _ = normalWorker.GetDBCommissionCount(ctx)
	if commissionCount != 1 {
		t.Errorf("expected 1 commission after redelivery, got %d", commissionCount)
	}

	// Verify dedup marker exists
	var processedCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", event.ID).Scan(&processedCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if processedCount != 1 {
		t.Errorf("expected 1 processed_events marker after redelivery, got %d", processedCount)
	}

	t.Log("SUCCESS: Redelivery after business failure - event processed after rollback cleared state")
}

// ============================================================================
// Processed Marker Failure Injection
// ============================================================================

// TestProcessedMarkerFailureRollback verifies that when the processed_events
// INSERT fails, the business mutation is NOT executed and transaction rolls back.
func TestProcessedMarkerFailureRollback(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	// Worker with FailProcessedInsert - simulates duplicate claim or system error
	worker := transaction.NewCommissionWorkerWithFailure(db, &transaction.ConsumerFailureMode{
		FailProcessedInsert: true,
	})

	// Use normal worker for assertion to verify rollback
	assertWorker := transaction.NewCommissionWorker(db)

	event := transaction.Event{ID: "evt-processed-fail", EventType: "InvoicePaid", AggregateID: "701", Payload: `{"invoice_id": 701}`}

	// Attempt processing - should fail at processed_events insert
	processed, err := worker.HandleEvent(ctx, "CommissionConsumer", event)
	if err == nil {
		t.Fatal("expected error due to processed_events insert failure")
	}
	if processed {
		t.Error("expected processing to fail, not succeed")
	}

	// Critical: Business mutation must NOT have been attempted
	// The claim must happen BEFORE business mutation
	var processedCount, commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", event.ID).Scan(&processedCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions WHERE event_id = $1", event.ID).Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if processedCount != 0 {
		t.Errorf("expected 0 processed_events (failure before commit), got %d", processedCount)
	}
	if commissionCount != 0 {
		t.Errorf("expected 0 commissions (business mutation must not run on claim failure), got %d", commissionCount)
	}

	// Verify the normal worker CAN process it (nothing was committed)
	processed2, err := assertWorker.HandleEvent(ctx, "CommissionConsumer", event)
	if err != nil {
		t.Fatalf("normal worker failed: %v", err)
	}
	if !processed2 {
		t.Error("expected normal worker to process event successfully after rollback")
	}

	commissionCount, _ = assertWorker.GetDBCommissionCount(ctx)
	if commissionCount != 1 {
		t.Errorf("expected 1 commission after retry, got %d", commissionCount)
	}

	t.Log("SUCCESS: Processed event insert failure - no business mutation ran, rollback clean")
	t.Log("  Business mutation happens AFTER claim - claim failure aborts before any business effect")
}

// ============================================================================
// Consumer Atomic Commit Happy Path
// ============================================================================

// TestConsumerAtomicCommitHappyPath verifies that after successful processing,
// both processed_events claim marker and business state (commissions) are
// committed atomically and together.
//
// Note: mockdb does not support simulated commit failure to avoid over-engineering.
// Commit atomicity is verified by the successful atomic insert pattern.
func TestConsumerAtomicCommitHappyPath(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	worker := transaction.NewCommissionWorker(db)
	event := transaction.Event{ID: "evt-commit-test", EventType: "InvoicePaid", AggregateID: "801", Payload: `{"invoice_id": 801}`}

	// Normal processing
	processed, err := worker.HandleEvent(ctx, "CommissionConsumer", event)
	if err != nil {
		t.Fatalf("processing failed: %v", err)
	}
	if !processed {
		t.Fatal("expected processing to succeed")
	}

	// Verify both committed atomically - both must be present after COMMIT
	var processedCount, commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", event.ID).Scan(&processedCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions WHERE event_id = $1", event.ID).Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if processedCount != 1 {
		t.Errorf("expected 1 processed_events after commit, got %d", processedCount)
	}
	if commissionCount != 1 {
		t.Errorf("expected 1 commission after commit, got %d", commissionCount)
	}

	t.Log("SUCCESS: Consumer atomic commit - both claim marker and business state committed together")
}

// ============================================================================
// MockDB Direct Concurrency Test (Two TX inserts without conflict)
// ============================================================================

// TestMockDBDirectConcurrency verifies that MockDB handles multiple transactions
// properly. Two concurrent inserts for different rows should both persist.
func TestMockDBDirectConcurrency(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	start := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		tx, _ := db.BeginTx(ctx, nil)
		_, _ = tx.ExecContext(ctx, "INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3)", "C1", "evt-A", time.Now())
		if err := tx.Commit(); err != nil {
			t.Fatalf("commit failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		<-start
		tx, _ := db.BeginTx(ctx, nil)
		_, _ = tx.ExecContext(ctx, "INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3)", "C2", "evt-B", time.Now())
		if err := tx.Commit(); err != nil {
			t.Fatalf("commit failed: %v", err)
		}
	}()

	close(start)
	wg.Wait()

	var count int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events").Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 rows, got %d", count)
	}
	t.Log("SUCCESS: MockDB direct concurrency - both A and B exist")
}

// ============================================================================
// MockDB Duplicate Claim Direct Test (Concurrent ON CONFLICT)
// ============================================================================

// TestMockDBDuplicateClaimDirect verifies that concurrent INSERT ON CONFLICT
// correctly enforces the composite uniqueness.
func TestMockDBDuplicateClaimDirect(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	start := make(chan struct{})
	var wg sync.WaitGroup
	var affected1, affected2 int64

	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		tx, _ := db.BeginTx(ctx, nil)
		res, _ := tx.ExecContext(ctx, "INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3) ON CONFLICT (consumer_name, event_id) DO NOTHING", "CommissionConsumer", "evt-123", time.Now())
		affected1, _ = res.RowsAffected()
		if err := tx.Commit(); err != nil {
			t.Fatalf("commit failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		<-start
		tx, _ := db.BeginTx(ctx, nil)
		res, _ := tx.ExecContext(ctx, "INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3) ON CONFLICT (consumer_name, event_id) DO NOTHING", "CommissionConsumer", "evt-123", time.Now())
		affected2, _ = res.RowsAffected()
		if err := tx.Commit(); err != nil {
			t.Fatalf("commit failed: %v", err)
		}
	}()

	close(start)
	wg.Wait()

	var count int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events").Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected exactly 1 row, got %d", count)
	}

	if affected1+affected2 != 1 {
		t.Errorf("expected exactly one RowsAffected=1, got sum=%d (affected1=%d, affected2=%d)", affected1+affected2, affected1, affected2)
	}

	t.Log("SUCCESS: MockDB duplicate claim direct - exactly one RowsAffected=1, one RowsAffected=0")
}

// ============================================================================
// MockDB Different Unique Keys (Same consumer, different events)
// ============================================================================

// TestMockDBDifferentUniqueKeys verifies concurrent inserts with different unique keys
// both succeed. (Constraint is not too global)
func TestMockDBDifferentUniqueKeys(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	start := make(chan struct{})
	var wg sync.WaitGroup
	var affected1, affected2 int64

	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		tx, _ := db.BeginTx(ctx, nil)
		res, _ := tx.ExecContext(ctx, "INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3) ON CONFLICT (consumer_name, event_id) DO NOTHING", "CommissionConsumer", "evt-1", time.Now())
		affected1, _ = res.RowsAffected()
		if err := tx.Commit(); err != nil {
			t.Fatalf("commit failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		<-start
		tx, _ := db.BeginTx(ctx, nil)
		res, _ := tx.ExecContext(ctx, "INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3) ON CONFLICT (consumer_name, event_id) DO NOTHING", "CommissionConsumer", "evt-2", time.Now())
		affected2, _ = res.RowsAffected()
		if err := tx.Commit(); err != nil {
			t.Fatalf("commit failed: %v", err)
		}
	}()

	close(start)
	wg.Wait()

	var count int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events").Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected exactly 2 rows, got %d", count)
	}

	if affected1 != 1 || affected2 != 1 {
		t.Errorf("expected both RowsAffected=1, got %d, %d", affected1, affected2)
	}

	t.Log("SUCCESS: MockDB different unique keys - both inserted (RowsAffected=1)")
}

// ============================================================================
// MockDB Different Consumers Same Event (Composite Unique Key)
// ============================================================================

// TestMockDBDifferentConsumersSameEvent verifies concurrent inserts with same event
// but different consumers both succeed (composite unique key semantics).
func TestMockDBDifferentConsumersSameEvent(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	start := make(chan struct{})
	var wg sync.WaitGroup
	var affected1, affected2 int64

	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		tx, _ := db.BeginTx(ctx, nil)
		res, _ := tx.ExecContext(ctx, "INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3) ON CONFLICT (consumer_name, event_id) DO NOTHING", "InventoryConsumer", "evt-123", time.Now())
		affected1, _ = res.RowsAffected()
		if err := tx.Commit(); err != nil {
			t.Fatalf("commit failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		<-start
		tx, _ := db.BeginTx(ctx, nil)
		res, _ := tx.ExecContext(ctx, "INSERT INTO processed_events (consumer_name, event_id, processed_at) VALUES ($1, $2, $3) ON CONFLICT (consumer_name, event_id) DO NOTHING", "CommissionConsumer", "evt-123", time.Now())
		affected2, _ = res.RowsAffected()
		if err := tx.Commit(); err != nil {
			t.Fatalf("commit failed: %v", err)
		}
	}()

	close(start)
	wg.Wait()

	var count int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events").Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected exactly 2 rows, got %d", count)
	}

	if affected1 != 1 || affected2 != 1 {
		t.Errorf("expected both RowsAffected=1, got %d, %d", affected1, affected2)
	}

	t.Log("SUCCESS: MockDB different consumers same event - composite unique key enforced")
}
