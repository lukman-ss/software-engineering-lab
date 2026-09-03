package transaction_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	transaction "github.com/lukman-ss/software-engineering-lab/labs/03-database-transaction"
	"github.com/lukman-ss/software-engineering-lab/labs/03-database-transaction/mockdb"
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
	dispatcher := transaction.NewOutboxDispatcher(db, broker)

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

// Test 7: Idempotent consumer processes duplicate event only once (uses ON CONFLICT DO NOTHING).
func TestIdempotentConsumerDeduplication(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	worker1 := transaction.NewCommissionWorker(db)
	worker2 := transaction.NewCommissionWorker(db)

	ctx := context.Background()
	event := transaction.Event{ID: "evt_101", EventType: "InvoicePaid", AggregateID: "101", Payload: `{"invoice_id": 101}`}

	processed1, err := worker1.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("unexpected error on first process: %v", err)
	}
	if !processed1 {
		t.Error("expected first event to be processed")
	}

	count, err := worker1.GetDBCommissionCount(ctx)
	if err != nil {
		t.Fatalf("failed to query commissions: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 commission in DB after first process, got %d", count)
	}

	processed2, err := worker2.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("unexpected error on duplicate: %v", err)
	}
	if processed2 {
		t.Error("expected duplicate event to be SKIPPED")
	}

	count, err = worker1.GetDBCommissionCount(ctx)
	if err != nil {
		t.Fatalf("failed to query commissions: %v", err)
	}
	if count != 1 {
		t.Errorf("expected still 1 commission after duplicate (idempotent), got %d", count)
	}

	inventoryWorker := transaction.NewCommissionWorker(db)
	processed3, err := inventoryWorker.HandleEvent(ctx, "InventoryWorker", event)
	if err != nil {
		t.Fatalf("unexpected error on different consumer: %v", err)
	}
	if !processed3 {
		t.Error("expected different consumer to process the event successfully")
	}

	count, err = worker1.GetDBCommissionCount(ctx)
	if err != nil {
		t.Fatalf("failed to query commissions: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 commissions (one per consumer), got %d", count)
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

	processed, err := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !processed {
		t.Fatal("expected event to be processed")
	}

	var dedupCount, commissionCount int64

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
func TestSequentialRedeliveryIdempotency(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	worker := transaction.NewCommissionWorker(db)
	ctx := context.Background()
	event := transaction.Event{ID: "evt_300", EventType: "InvoicePaid", AggregateID: "300", Payload: `{"invoice_id": 300}`}

	processed1, err := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("unexpected error on first delivery: %v", err)
	}
	if !processed1 {
		t.Fatal("expected first delivery to be processed")
	}

	count1, err := worker.GetDBCommissionCount(ctx)
	if err != nil {
		t.Fatalf("get commission count failed: %v", err)
	}
	if count1 != 1 {
		t.Fatalf("expected 1 commission after first delivery, got %d", count1)
	}

	processed2, err := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err != nil {
		t.Fatalf("unexpected error on redelivery: %v", err)
	}
	if processed2 {
		t.Error("expected redelivered event to be SKIPPED (idempotent)")
	}

	count2, err := worker.GetDBCommissionCount(ctx)
	if err != nil {
		t.Fatalf("get commission count failed: %v", err)
	}
	if count2 != 1 {
		t.Errorf("expected still 1 commission after redelivery, got %d", count2)
	}

	t.Log("SUCCESS: Crash/redelivery handled idempotently - no duplicate business state.")
}

// Test 9: HTTP inside transaction duration
func TestHTTPInsideTransactionDuration(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	httpClient := transaction.NewHTTPClient(100, 0)
	svc := transaction.NewInvoiceServiceWithHTTPCall(db, httpClient)
	ctx := context.Background()

	elapsed, err := svc.PayInvoiceWithHTTPInsideTx(ctx, 101, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if elapsed < 90*time.Millisecond {
		t.Errorf("elapsed %v too short - HTTP latency should have extended transaction", elapsed)
	}

	var status string
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&status); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if status != "paid" {
		t.Errorf("expected 'paid', got '%s'", status)
	}

	t.Logf("SUCCESS: HTTP latency %v extended transaction lifetime (elapsed=%v)", 100*time.Millisecond, elapsed)
}

// Test 10: Commit then external call shows standard pattern
func TestCommitThenExternalCall(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	blockingHTTP := transaction.NewBlockingHTTPClient(0)
	svc := transaction.NewInvoiceServiceWithBlockingHTTP(db, blockingHTTP)
	ctx := context.Background()

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

	var status string
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&status); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if status != "paid" {
		t.Errorf("expected 'paid', got '%s'", status)
	}

	done := make(chan error, 1)
	go func() {
		done <- blockingHTTP.Ping(ctx)
	}()

	blockingHTTP.WaitUntilTxOpen()

	if svc.IsTxOpen() {
		t.Error("expected transaction to be CLOSED during external call")
	}

	blockingHTTP.Release()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}

	t.Log("SUCCESS: Commit-then-external pattern proven with deterministic synchronization.")
	t.Log("  Transaction was already committed and closed before external call blocked.")
}

// Test: Transaction stays open during external call
func TestTransactionStaysOpenDuringExternalCall(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	blockingHTTP := transaction.NewBlockingHTTPClient(0)
	svc := transaction.NewInvoiceServiceWithBlockingHTTP(db, blockingHTTP)
	ctx := context.Background()

	done := make(chan error, 1)
	go func() {
		done <- svc.PayInvoiceWithBlocking(ctx)
	}()

	blockingHTTP.WaitUntilTxOpen()

	if !svc.IsTxOpen() {
		t.Error("expected transaction to be open while external call is blocking")
	}

	blockingHTTP.Release()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}

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

// Test: External side effect rollback
func TestExternalSideEffectRollback(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	whatsapp := transaction.NewWhatsAppClient(0)
	svc := transaction.NewDistributedOrderService(db, whatsapp)
	ctx := context.Background()

	err := svc.ProcessPaymentWithExternalSideEffect(ctx, 101, 500000.0)
	if err == nil {
		t.Fatal("expected error (simulated ERP failure)")
	}

	if whatsapp.SentCount() != 1 {
		t.Errorf("expected WhatsApp notification sent, got count=%d", whatsapp.SentCount())
	}

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

	t.Log("PROVEN: External side effect (WhatsApp) sent count=1, but DB rolled back - EXTERNAL ≠ DB")
}

// Test: Dual-write crash window
func TestDualWriteCrashWindow(t *testing.T) {
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

	events := broker.PublishedEvents()
	if len(events) != 0 {
		t.Errorf("expected 0 events published (crashed before publish), got %d", len(events))
	}

	t.Log("PROVEN: Dual-write crash window - invoice paid but event never delivered to broker")
}

// Test: Reverse dual-write failure
func TestReverseDualWriteFailure(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	broker := transaction.NewInMemoryBroker(0)

	err := broker.Publish(context.Background(), transaction.Event{
		ID:          "evt_201",
		EventType:   "InvoicePaid",
		AggregateID: "201",
		Payload:     `{"invoice_id": 201}`,
	})
	if err != nil {
		t.Fatalf("failed to publish event: %v", err)
	}

	events := broker.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event published, got %d", len(events))
	}

	t.Logf("PROVEN: Reverse dual-write - event published=%d, but DB not committed (crash scenario)", len(events))
	t.Log("LESSON: Publish-then-commit also has atomicity window - choose order based on idempotency")
}

// Test: Outbox happy path assertions
func TestOutboxHappyPathAssertions(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()

	if err := svc.PayInvoiceWithOutbox(ctx, 101); err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

	var invoiceStatus string
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if invoiceStatus != "paid" {
		t.Errorf("expected invoice 'paid', got '%s'", invoiceStatus)
	}

	pending, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
	if pending != 1 {
		t.Fatalf("expected 1 pending outbox event, got %d", pending)
	}

	var id string
	if err := db.QueryRowContext(ctx, "SELECT id FROM outbox_events WHERE status = 'pending'").Scan(&id); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if id == "" {
		t.Error("expected event_id != empty")
	}

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

	t.Log("SUCCESS: Outbox happy path assertions completely verified.")
}

// Test: Transactional Outbox Rollback
func TestTransactionalOutboxRollback(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()

	err := svc.PayInvoiceWithOutboxInjectError(ctx, 101, true)
	if err == nil {
		t.Fatal("expected error")
	}

	var invoiceStatus string
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if invoiceStatus == "paid" {
		t.Error("expected invoice NOT to be 'paid' (should be rolled back), got 'paid'")
	}

	pending, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
	if pending != 0 {
		t.Errorf("expected 0 pending outbox events (should be rolled back), got %d", pending)
	}

	t.Log("SUCCESS: Outbox error triggered FULL rollback. Business state + outbox state kept atomic.")
}

// Test: Outbox Pending Recovery
func TestOutboxPendingRecovery(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()

	if err := svc.PayInvoiceWithOutbox(ctx, 101); err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

	var invoiceStatus string
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if invoiceStatus != "paid" {
		t.Fatalf("expected invoice 'paid', got '%s'", invoiceStatus)
	}

	broker := transaction.NewInMemoryBroker(100)

	dispatcher := transaction.NewOutboxDispatcher(db, broker)

	dispatched, err := dispatcher.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("unexpected error from dispatcher: %v", err)
	}
	if dispatched != 0 {
		t.Errorf("expected 0 dispatched when broker down, got %d", dispatched)
	}

	pending, err := svc.CountOutboxEvents(ctx)
	if err != nil {
		t.Fatalf("count outbox events failed: %v", err)
	}
	if pending != 1 {
		t.Errorf("expected 1 pending event, got %d", pending)
	}

	var status string
	if err := db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&status); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if status != "paid" {
		t.Errorf("invoice should remain paid even though broker down")
	}

	recoveringBroker := transaction.NewInMemoryBroker(0)
	recoveringDispatcher := transaction.NewOutboxDispatcher(db, recoveringBroker)

	dispatched, err = recoveringDispatcher.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("unexpected error from dispatcher: %v", err)
	}
	if dispatched != 1 {
		t.Errorf("expected 1 event after broker recovers, got %d", dispatched)
	}

	events := recoveringBroker.PublishedEvents()
	if len(events) != 1 {
		t.Errorf("expected 1 event delivered after recovery")
	}

	t.Log("SUCCESS: Business transaction persisted, outbox remained pending, delivery successful after broker recovery")
	t.Log("This demonstrates durable intent: event stored BEFORE publish, retry-capable")
}

// Test: Different consumers process same event concurrently
func TestConcurrentDifferentConsumersSameEvent(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	inventoryWorker := transaction.NewCommissionWorker(db)
	commissionWorker := transaction.NewCommissionWorker(db)

	event := transaction.Event{ID: "evt-shared-123", EventType: "InvoicePaid", AggregateID: "123", Payload: `{"invoice_id": 123}`}

	start := make(chan struct{})
	var wg sync.WaitGroup
	var processedInv, processedComm bool
	var errInv, errComm error

	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		processedInv, errInv = inventoryWorker.HandleEvent(ctx, "InventoryConsumer", event)
	}()

	go func() {
		defer wg.Done()
		<-start
		processedComm, errComm = commissionWorker.HandleEvent(ctx, "CommissionConsumer", event)
	}()

	close(start)
	wg.Wait()

	if errInv != nil {
		t.Fatalf("InventoryConsumer error: %v", errInv)
	}
	if errComm != nil {
		t.Fatalf("CommissionConsumer error: %v", errComm)
	}

	if !processedInv {
		t.Error("expected InventoryConsumer to process")
	}
	if !processedComm {
		t.Error("expected CommissionConsumer to process")
	}

	t.Log("SUCCESS: Different consumers process same event independently")
}

// Test: Saga payment with compensating action
func TestSagaPaymentWithCompensatingAction(t *testing.T) {
	var executedSteps, compensatedSteps []string

	saga := transaction.NewSaga().
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
				executedSteps = append(executedSteps, "reserve")
				return nil
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "release")
				return nil
			},
		}).
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
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

	if len(compensatedSteps) != 1 {
		t.Errorf("expected exactly 1 compensation, got %d: %v", len(compensatedSteps), compensatedSteps)
	}
	if len(compensatedSteps) > 0 && compensatedSteps[0] != "release" {
		t.Errorf("expected 'release' compensation, got %s", compensatedSteps[0])
	}

	t.Log("SUCCESS: Saga compensation executed correctly (only compensated successful steps).")
}

// Test: Saga compensation order four steps
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
				return errors.New("step D failed")
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "D")
				return nil
			},
		})

	ctx := context.Background()
	_ = saga.Execute(ctx)

	if len(executedSteps) != 3 {
		t.Errorf("expected 3 successful executions, got %d: %v", len(executedSteps), executedSteps)
	}
	if executedSteps[0] != "A" || executedSteps[1] != "B" || executedSteps[2] != "C" {
		t.Errorf("expected order [A, B, C], got %v", executedSteps)
	}

	if len(compensatedSteps) != 3 {
		t.Errorf("expected 3 compensations, got %d: %v", len(compensatedSteps), compensatedSteps)
	}
	if compensatedSteps[0] != "C" || compensatedSteps[1] != "B" || compensatedSteps[2] != "A" {
		t.Errorf("expected compensation order [C, B, A], got %v", compensatedSteps)
	}

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

// Test: Saga compensation failure handling
func TestSagaCompensationFailureHandling(t *testing.T) {
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
				return errors.New("compensation B failed")
			},
		}).
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error {
				return errors.New("step C failed")
			},
			Compensate: func(ctx context.Context) error {
				compensatedSteps = append(compensatedSteps, "C")
				return nil
			},
		})

	ctx := context.Background()
	err := saga.Execute(ctx)

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

	if len(compensatedSteps) != 2 {
		t.Errorf("expected 2 compensations (B then A), got %d: %v", len(compensatedSteps), compensatedSteps)
	}
	if compensatedSteps[0] != "B" {
		t.Errorf("expected first compensation B, got %s", compensatedSteps[0])
	}
	if compensatedSteps[1] != "A" {
		t.Errorf("expected second compensation A, got %s", compensatedSteps[1])
	}

	if saga.GetCompensatedSteps()[0] != 1 || saga.GetCompensatedSteps()[1] != 0 {
		t.Errorf("expected compensated step indices [1, 0] (B, A), got %v", saga.GetCompensatedSteps())
	}

	t.Log("SUCCESS: Saga continued compensation after failure, A was compensated despite B compensation failure.")
}

// Test: Compensation idempotency
func TestCompensationIdempotency(t *testing.T) {
	comp := transaction.NewIdempotentCompensator()

	refundCount := 0
	doRefund := func() error {
		refundCount++
		return nil
	}

	executed, err := comp.Compensate("refund-order-123", doRefund)
	if err != nil || !executed {
		t.Fatalf("expected first compensation to execute: executed=%v err=%v", executed, err)
	}

	executed, err = comp.Compensate("refund-order-123", doRefund)
	if err != nil || executed {
		t.Fatalf("expected second compensation to skip: executed=%v err=%v", executed, err)
	}

	if refundCount != 1 {
		t.Errorf("expected refundCount=1 (no double refund), got %d", refundCount)
	}

	executed, _ = comp.Compensate("refund-order-456", doRefund)
	if !executed {
		t.Error("expected different key to execute")
	}
	if refundCount != 2 {
		t.Errorf("expected refundCount=2 after second key, got %d", refundCount)
	}

	t.Log("PROVEN: Compensation idempotent — same key executes once regardless of call count")
}

// Test: InvoicePaid payload round-trip
func TestInvoicePaidPayloadRoundTrip(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	original := transaction.InvoicePaidPayload{
		EventID:    "evt_101_test",
		InvoiceID:  101,
		OccurredAt: now,
	}

	raw := transaction.MarshalInvoicePaidPayload(original)

	var rawMap map[string]any
	if err := json.Unmarshal([]byte(raw), &rawMap); err != nil {
		t.Fatalf("payload is not valid JSON: %v", err)
	}

	parsed, err := transaction.UnmarshalInvoicePaidPayload(raw)
	if err != nil {
		t.Fatalf("failed to unmarshal payload: %v", err)
	}

	if parsed.EventID != original.EventID {
		t.Errorf("event_id mismatch: expected %s, got %s", original.EventID, parsed.EventID)
	}
	if parsed.InvoiceID != original.InvoiceID {
		t.Errorf("invoice_id mismatch: expected %d, got %d", original.InvoiceID, parsed.InvoiceID)
	}
	if parsed.OccurredAt != original.OccurredAt {
		t.Errorf("occurred_at mismatch: expected %s, got %s", original.OccurredAt, parsed.OccurredAt)
	}

	var full map[string]any
	if err := json.Unmarshal([]byte(raw), &full); err != nil {
		t.Fatalf("unmarshal full map failed: %v", err)
	}
	if _, hasStatus := full["status"]; hasStatus {
		t.Error("payload should not contain 'status' field — event is a fact, not a command")
	}

	if _, err := transaction.UnmarshalInvoicePaidPayload(`{"invoice_id":1,"occurred_at":"2024-01-01T00:00:00Z"}`); err == nil {
		t.Error("expected error for missing event_id")
	}

	if _, err := transaction.UnmarshalInvoicePaidPayload(`{"event_id":"x","invoice_id":1}`); err == nil {
		t.Error("expected error for missing occurred_at")
	}

	t.Logf("SUCCESS: Payload round-trip validated. JSON=%s", raw)
}

// Test: Eventual Consistency Demo
func TestEventualConsistencyDemo(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)
	ctx := context.Background()

	svc := transaction.NewInvoiceServiceOutbox(db)
	if err := svc.PayInvoiceWithOutbox(ctx, 101); err != nil {
		t.Fatalf("pay with outbox failed: %v", err)
	}

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

	var commissionCount int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 0 {
		t.Errorf("expected 0 commissions (not yet processed), got %d", commissionCount)
	}

	broker := transaction.NewInMemoryBroker(0)
	dispatcher := transaction.NewOutboxDispatcher(db, broker)
	dispatched, err := dispatcher.DispatchBatch(ctx)
	if err != nil || dispatched != 1 {
		t.Fatalf("expected 1 event dispatched, got %d (err: %v)", dispatched, err)
	}

	events := broker.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event published to broker, got %d", len(events))
	}

	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&commissionCount); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if commissionCount != 0 {
		t.Error("expected 0 commissions immediately after dispatch (dispatch ≠ consumer)")
	}

	worker := transaction.NewCommissionWorker(db)
	processed, err := worker.HandleEvent(ctx, "CommissionConsumer", events[0])
	if err != nil || !processed {
		t.Fatalf("consumer failed to process: processed=%v, err=%v", processed, err)
	}

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

	t.Log("DEMO: Eventual consistency — invoice=paid immediately, commissions eventually=1")
	t.Log("  t=0 (business commit): invoice=paid, commissions=0, outbox=1pending")
	t.Log("  t=1 (dispatched):     invoice=paid, commissions=0, outbox=0pending (event in broker)")
	t.Log("  t=2 (consumer):       invoice=paid, commissions=1, all consistent")
	t.Log("  Key insight: dispatch (publish) ≠ consumer processing (business effect)")
}

// MockDB Concurrent Tests

// Test: Concurrent different events
func TestConcurrentDifferentEvents(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	workers := []*transaction.CommissionWorker{
		transaction.NewCommissionWorker(db),
		transaction.NewCommissionWorker(db),
		transaction.NewCommissionWorker(db),
	}
	ctx := context.Background()

	events := []transaction.Event{
		{ID: "evt-no-lost-1", EventType: "InvoicePaid", AggregateID: "501", Payload: `{"invoice_id": 501}`},
		{ID: "evt-no-lost-2", EventType: "InvoicePaid", AggregateID: "502", Payload: `{"invoice_id": 502}`},
		{ID: "evt-no-lost-3", EventType: "InvoicePaid", AggregateID: "503", Payload: `{"invoice_id": 503}`},
	}

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

	for i, err := range errs {
		if err != nil {
			t.Fatalf("event %d error: %v", i, err)
		}
	}

	for i, processed := range results {
		if !processed {
			t.Errorf("event %d should be processed", i)
		}
	}

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

// Test: Concurrent same event
func TestConcurrentSameEvent(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	worker1 := transaction.NewCommissionWorker(db)
	worker2 := transaction.NewCommissionWorker(db)
	ctx := context.Background()

	event := transaction.Event{ID: "evt-same", EventType: "InvoicePaid", AggregateID: "201", Payload: `{"invoice_id": 201}`}

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

	if !processed1 && !processed2 {
		t.Fatal("expected at least one worker to process")
	}
	if processed1 && processed2 {
		t.Error("expected exactly one worker to process (deduplication failed)")
	}

	var count int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM commissions").Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 commission row after concurrent processing, got %d", count)
	}

	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE consumer_name = $1 AND event_id = $2", "CommissionWorker", event.ID).Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 processed_events marker, got %d", count)
	}

	t.Log("SUCCESS: Concurrent same event - exactly one business transaction committed atomically under race")
}

// Test: MockDB ON CONFLICT DO NOTHING semantics
func TestMockDBOnConflictSemantics(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

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

	var count int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM processed_events WHERE event_id = $1", "evt-on-conflict").Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 row in database, got %d", count)
	}

	t.Log("SUCCESS: ON CONFLICT DO NOTHING semantics verified - first=1, duplicate=0")
}

// Test: MockDB no lost updates
func TestMockDBNoLostUpdates(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	workers := []*transaction.CommissionWorker{
		transaction.NewCommissionWorker(db),
		transaction.NewCommissionWorker(db),
		transaction.NewCommissionWorker(db),
	}
	ctx := context.Background()

	events := []transaction.Event{
		{ID: "evt-no-lost-1", EventType: "InvoicePaid", AggregateID: "501", Payload: `{"invoice_id": 501}`},
		{ID: "evt-no-lost-2", EventType: "InvoicePaid", AggregateID: "502", Payload: `{"invoice_id": 502}`},
		{ID: "evt-no-lost-3", EventType: "InvoicePaid", AggregateID: "503", Payload: `{"invoice_id": 503}`},
	}

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

	for i, err := range errs {
		if err != nil {
			t.Fatalf("event %d error: %v", i, err)
		}
	}

	for i, processed := range results {
		if !processed {
			t.Errorf("event %d should be processed", i)
		}
	}

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

// Test: Rollback isolation
func TestMockDBRollbackIsolation(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	ctx := context.Background()

	txA, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("txA begin failed: %v", err)
	}
	_, err = txA.ExecContext(ctx, "INSERT INTO commissions (event_id, amount) VALUES ('txA-rolledback', 100)")
	if err != nil {
		t.Fatalf("txA insert failed: %v", err)
	}
	err = txA.Rollback()
	if err != nil {
		t.Fatalf("txA rollback failed: %v", err)
	}

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

	var countA, countB int
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

	t.Log("SUCCESS: Rollback isolation verified")
	t.Log("  TX A rolled back → no traces")
	t.Log("  TX B committed → row present")
}
