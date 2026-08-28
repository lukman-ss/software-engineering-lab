package transaction_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	transaction "github.com/lukman/software-engineer-lab/labs/03-database-transaction"
	"github.com/lukman/software-engineer-lab/labs/03-database-transaction/mockdb"
)

func seedTestDB(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	tables := []struct {
		sql string
	}{
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

// ============================================================================
// 1. Unsafe Local Transaction - Partial State Corruption
// ============================================================================

func TestUnsafeLocalTransaction(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewPaymentServiceUnsafe(db)
	ctx := context.Background()

	err := svc.ProcessPayment(ctx, 101, 250000.0, true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var paymentCount int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM payments WHERE order_id = $1", 101).Scan(&paymentCount)
	if err != nil {
		t.Fatalf("query payments failed: %v", err)
	}
	if paymentCount != 1 {
		t.Errorf("expected 1 payment in partial state, got %d", paymentCount)
	}

	t.Log("SUCCESS: Unsafe local transaction demonstrated partial state corruption.")
}

// ============================================================================
// 2. Safe Local Transaction - ACID Rollback
// ============================================================================

func TestSafeLocalTransaction(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewPaymentServiceSafe(db)
	ctx := context.Background()

	err := svc.ProcessPayment(ctx, 101, 250000.0, true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var paymentCount int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM payments WHERE order_id = $1", 101).Scan(&paymentCount)
	if err != nil {
		t.Fatalf("query payments failed: %v", err)
	}
	if paymentCount != 0 {
		t.Errorf("expected 0 payments due to ROLLBACK, got %d", paymentCount)
	}

	t.Log("SUCCESS: Safe local transaction demonstrated clean ACID ROLLBACK.")
}

// ============================================================================
// 3. Distributed Transaction - External Side Effect Cannot Be Rolled Back
// ============================================================================

func TestDistributedTransactionExternalSideEffectLimitation(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	whatsapp := transaction.NewWhatsAppClient(0)
	svc := transaction.NewDistributedOrderService(db, whatsapp)
	ctx := context.Background()

	err := svc.ProcessPaymentWithExternalSideEffect(ctx, 101, 500000.0)
	if err == nil {
		t.Fatal("expected error from simulated ERP failure, got nil")
	}

	var paymentCount int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM payments WHERE order_id = $1", 101).Scan(&paymentCount)
	if err != nil {
		t.Fatalf("query payments failed: %v", err)
	}
	if paymentCount != 0 {
		t.Errorf("expected database payment to be rolled back (count = 0), got %d", paymentCount)
	}

	sentCount := whatsapp.SentCount()
	if sentCount != 1 {
		t.Errorf("expected 1 WhatsApp notification sent despite DB rollback, got %d", sentCount)
	}

	t.Logf("PROVEN: DB rollback successfully occurred (0 payments), BUT external WhatsApp side effect CANNOT be undone! Message sent: %q", whatsapp.GetMessages()[0])
}

// ============================================================================
// 4. HTTP Call Inside Transaction - Anti-Pattern Demonstration
// ============================================================================

func TestHTTPCallInsideTransactionLatency(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	httpClient := transaction.NewHTTPClient(50, 0)
	svc := transaction.NewInvoiceServiceWithHTTPCall(db, httpClient)

	ctx := context.Background()
	elapsed, err := svc.PayInvoiceWithHTTPInsideTx(ctx, 101, false)
	if err != nil {
		t.Fatalf("pay invoice failed: %v", err)
	}

	if elapsed < 40*time.Millisecond {
		t.Errorf("expected transaction to take at least 40ms (due to HTTP latency), got %v", elapsed)
	}

	t.Logf("SUCCESS: HTTP call inside transaction extended transaction duration to %v", elapsed)
}

// ============================================================================
// 5. Dual-Write Problem - Event Lost After DB Commit
// ============================================================================

func TestDualWriteProblemEventLost(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	broker := transaction.NewInMemoryBroker(0)
	svc := transaction.NewInvoiceServiceDualWrite(db, broker)

	ctx := context.Background()
	err := svc.PayInvoiceDualWrite(ctx, 101, true)
	if err == nil {
		t.Fatal("expected ErrProcessCrashed, got nil")
	}

	var invoiceStatus string
	err = db.QueryRowContext(ctx, "SELECT status FROM invoices WHERE order_id = $1", 101).Scan(&invoiceStatus)
	if err != nil {
		t.Fatalf("query invoices failed: %v", err)
	}
	if invoiceStatus != "paid" {
		t.Errorf("expected invoice status 'paid' (DB commit happened), got '%s'", invoiceStatus)
	}

	publishedCount := len(broker.PublishedEvents())
	if publishedCount != 0 {
		t.Errorf("expected 0 events published (process crashed), got %d", publishedCount)
	}

	t.Log("PROVEN: Dual-write problem - invoice is 'paid' but NO Event was published!")
}

// ============================================================================
// 6 & 7. Transactional Outbox & Dispatcher with Retry Mechanism (Poin 10)
// ============================================================================

func TestOutboxDispatcherWithRetry(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	outboxSvc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()

	// 1. Store business state + outbox event atomically
	err := outboxSvc.PayInvoiceWithOutbox(ctx, 101)
	if err != nil {
		t.Fatalf("pay invoice with outbox failed: %v", err)
	}

	// 2. Broker fails first 2 attempts, succeeds on 3rd attempt (failUpTo = 2)
	broker := transaction.NewInMemoryBroker(2)
	dispatcher := transaction.NewOutboxDispatcher(db, broker, 3, false)

	// Attempt 1: Broker fails, attempts incremented to 1
	_, err = dispatcher.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("batch 1 failed: %v", err)
	}
	if len(broker.PublishedEvents()) != 0 {
		t.Errorf("expected 0 published events on attempt 1, got %d", len(broker.PublishedEvents()))
	}

	// Attempt 2: Broker fails, attempts incremented to 2
	_, err = dispatcher.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("batch 2 failed: %v", err)
	}
	if len(broker.PublishedEvents()) != 0 {
		t.Errorf("expected 0 published events on attempt 2, got %d", len(broker.PublishedEvents()))
	}

	// Attempt 3: Broker succeeds! Event published, status marked published
	_, err = dispatcher.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("batch 3 failed: %v", err)
	}

	events := broker.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event published on attempt 3, got %d", len(events))
	}
	if events[0].EventType != "InvoicePaid" {
		t.Errorf("expected event type 'InvoicePaid', got %s", events[0].EventType)
	}

	// Verify pending count is now 0
	pending, _ := outboxSvc.CountOutboxEvents(ctx)
	if pending != 0 {
		t.Errorf("expected 0 pending events after successful dispatch, got %d", pending)
	}

	t.Log("SUCCESS: Outbox dispatcher successfully retried (Attempt 1 fail → Attempt 2 fail → Attempt 3 success) with max attempts limit.")
}

// ============================================================================
// 8. Duplicate Delivery (At-Least-Once Delivery)
// ============================================================================

func TestOutboxDuplicateDeliveryAtLeastOnce(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	outboxSvc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()

	_ = outboxSvc.PayInvoiceWithOutbox(ctx, 101)

	broker := transaction.NewInMemoryBroker(0) // always succeeds
	// Simulate crash before marking as published (crashBeforeMark = true)
	dispatcherCrashing := transaction.NewOutboxDispatcher(db, broker, 3, true)

	// Dispatcher publishes event but crashes before marking published
	_, err := dispatcherCrashing.DispatchBatch(ctx)
	if err == nil {
		t.Fatal("expected ErrProcessCrashed, got nil")
	}

	// Event was published to broker...
	events1 := broker.PublishedEvents()
	if len(events1) != 1 {
		t.Fatalf("expected 1 event published, got %d", len(events1))
	}

	// ...BUT event is STILL PENDING in outbox because it crashed before mark published!
	pending, _ := outboxSvc.CountOutboxEvents(ctx)
	if pending != 1 {
		t.Errorf("expected 1 pending event remaining in outbox, got %d", pending)
	}

	// When dispatcher recovers / restarts (normal dispatcher):
	dispatcherHealthy := transaction.NewOutboxDispatcher(db, broker, 3, false)
	_, err = dispatcherHealthy.DispatchBatch(ctx)
	if err != nil {
		t.Fatalf("healthy dispatcher failed: %v", err)
	}

	// Event is published a SECOND time! (Duplicate delivery)
	events2 := broker.PublishedEvents()
	if len(events2) != 2 {
		t.Fatalf("expected 2 total deliveries (duplicate delivery), got %d", len(events2))
	}

	if events2[0].ID != events2[1].ID {
		t.Errorf("expected same event ID for duplicate delivery")
	}

	t.Logf("PROVEN: Transactional Outbox provides AT-LEAST-ONCE delivery. Event delivered twice: %s and %s", events2[0].ID, events2[1].ID)
}

// ============================================================================
// 9. Idempotent Consumer & Deduplication (Connected to Lab #01)
// ============================================================================

func TestIdempotentConsumerDeduplication(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	worker := transaction.NewCommissionWorker(db)
	ctx := context.Background()

	event := transaction.Event{
		ID:          "evt_invoice_101",
		EventType:   "InvoicePaid",
		AggregateID: "101",
		Payload:     `{"invoice_id": 101}`,
	}

	// Receive duplicate delivery of the exact same event ID
	processed1, err1 := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err1 != nil {
		t.Fatalf("first handle failed: %v", err1)
	}
	if !processed1 {
		t.Errorf("expected first event to be processed")
	}

	// Second delivery (duplicate from At-Least-Once delivery mechanism)
	processed2, err2 := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err2 != nil {
		t.Fatalf("second handle failed: %v", err2)
	}
	if processed2 {
		t.Errorf("expected duplicate event to be SKIPPED (idempotency)")
	}

	// Third delivery (another duplicate)
	processed3, err3 := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err3 != nil {
		t.Fatalf("third handle failed: %v", err3)
	}
	if processed3 {
		t.Errorf("expected third duplicate event to be SKIPPED")
	}

	// Verify business logic executed EXACTLY ONCE
	commissionsPaid := worker.CommissionsPaidCount()
	if commissionsPaid != 1 {
		t.Errorf("expected commissions paid exactly 1 time despite 3 deliveries, got %d", commissionsPaid)
	}

	t.Log("SUCCESS: Idempotent Consumer successfully deduplicated duplicate events via processed_events unique constraint. Commissions paid exactly once.")
}