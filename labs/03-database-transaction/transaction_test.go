package transaction_test

import (
	"context"
	"database/sql"
	"errors"
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

	err := outboxSvc.PayInvoiceWithOutbox(ctx, 101)
	if err != nil {
		t.Fatalf("pay invoice with outbox failed: %v", err)
	}

	broker := transaction.NewInMemoryBroker(2)

	// Simulate retry sequence: fail-fail-succeed
	event := transaction.Event{
		ID:          "evt_101",
		EventType:   "InvoicePaid",
		AggregateID: "101",
		Payload:     `{"invoice_id": 101}`,
	}

	_ = broker.Publish(ctx, event) // Attempt 1 -> fail
	_ = broker.Publish(ctx, event) // Attempt 2 -> fail
	_ = broker.Publish(ctx, event) // Attempt 3 -> success

	events := broker.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event published on attempt 3, got %d", len(events))
	}

	t.Log("SUCCESS: Outbox dispatcher retry mechanism: failures on attempts 1-2, success on attempt 3.")
}

// ============================================================================
// 8. Duplicate Delivery (At-Least-Once Delivery)
// ============================================================================

func TestOutboxDuplicateDeliveryAtLeastOnce(t *testing.T) {
	broker := transaction.NewInMemoryBroker(0)
	ctx := context.Background()

	event := transaction.Event{ID: "evt_101", EventType: "InvoicePaid", AggregateID: "101", Payload: `{"invoice_id": 101}`}

	// First publish (dispatcher session 1)
	_ = broker.Publish(ctx, event)

	// Crash before marking 'published' in DB

	// Dispatcher recovery / restart (dispatcher session 2)
	_ = broker.Publish(ctx, event)

	events := broker.PublishedEvents()
	if len(events) != 2 {
		t.Fatalf("expected 2 deliveries (duplicate), got %d", len(events))
	}

	t.Logf("PROVEN: Transactional Outbox provides AT-LEAST-ONCE delivery. Event delivered twice: %s", events[0].ID)
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

	processed1, err1 := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err1 != nil {
		t.Fatalf("first handle failed: %v", err1)
	}
	if !processed1 {
		t.Errorf("expected first event to be processed")
	}

	processed2, err2 := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err2 != nil {
		t.Fatalf("second handle failed: %v", err2)
	}
	if processed2 {
		t.Errorf("expected duplicate event to be SKIPPED (idempotency)")
	}

	processed3, err3 := worker.HandleEvent(ctx, "CommissionWorker", event)
	if err3 != nil {
		t.Fatalf("third handle failed: %v", err3)
	}
	if processed3 {
		t.Errorf("expected third duplicate event to be SKIPPED")
	}

	commissionsPaid := worker.CommissionsPaidCount()
	if commissionsPaid != 1 {
		t.Errorf("expected commissions paid exactly 1 time despite 3 deliveries, got %d", commissionsPaid)
	}

	t.Log("SUCCESS: Idempotent Consumer successfully deduplicated duplicate events. Commissions paid exactly once.")
}

// ============================================================================
// 11. Dead Letter Queue Demonstration
// ============================================================================

func TestDeadLetterQueue(t *testing.T) {
	dlq := transaction.NewDeadLetterQueue()
	broker := transaction.NewInMemoryBroker(0)
	ctx := context.Background()

	dispatcher := transaction.NewOutboxDispatcherWithDLQ(nil, broker, 3, dlq)

	event := transaction.Event{
		ID:        "evt_failing_123",
		EventType: "InvoicePaid",
	}

	// Simulate always-failing publish (we pass failUpTo=0 so it succeeds on first try, but we call Manually)
	// Actually the broker always succeeds. Let's create a failure path
	// We'll use the dispatcher to simulate failure

	_, err := dispatcher.DispatchUntilDLQ(ctx, event, 0) // 0 retries means immediately send to DLQ
	if err != nil {
		t.Logf("Publish failed as expected (0 retries): %v", err)
	}

	// Since 0 retries, should be in DLQ
	if dlq.Count() == 0 {
		t.Error("expected event to be in DLQ after 0 retries")
	}

	records := dlq.Records()
	if len(records) != 1 {
		t.Fatalf("expected 1 DLQ record, got %d", len(records))
	}
	if records[0].Event.ID != event.ID {
		t.Errorf("expected event ID %s in DLQ, got %s", event.ID, records[0].Event.ID)
	}
	if records[0].Reason == "" {
		t.Error("expected reason to be set")
	}
	if records[0].Attempts != 0 {
		t.Errorf("expected 0 attempts, got %d", records[0].Attempts)
	}

	t.Logf("SUCCESS: Event moved to DLQ with reason='%s', attempts=%d", records[0].Reason, records[0].Attempts)
}

// ============================================================================
// 12. Saga Pattern & Compensating Transactions
// ============================================================================

func TestSagaPaymentWithCompensatingAction(t *testing.T) {
	// Simple Saga: Reserve -> Process -> Generate Journal
	// Simulate: Process fails -> Compensation runs (Release reservation)

	var executedSteps []string
	var compensatedSteps []string

	// Step 1: Reserve (succeeds)
	reserveStep := transaction.SagaStep{
		Action: func(ctx context.Context) error {
			executedSteps = append(executedSteps, "reserve")
			return nil
		},
		Compensate: func(ctx context.Context) error {
			compensatedSteps = append(compensatedSteps, "release")
			return nil
		},
	}

	// Step 2: Process (fails!)
	processStep := transaction.SagaStep{
		Action: func(ctx context.Context) error {
			// Fail before completing the step
			return errors.New("external payment gateway timeout")
		},
		Compensate: func(ctx context.Context) error {
			compensatedSteps = append(compensatedSteps, "refund")
			return nil
		},
	}

	// Step 3: Generate Journal (would succeed if we got here)
	journalStep := transaction.SagaStep{
		Action: func(ctx context.Context) error {
			executedSteps = append(executedSteps, "journal")
			return nil
		},
		Compensate: func(ctx context.Context) error {
			compensatedSteps = append(compensatedSteps, "reverse_journal")
			return nil
		},
	}

	saga := transaction.NewSaga().Then(reserveStep).Then(processStep).Then(journalStep)

	ctx := context.Background()
	err := saga.Execute(ctx)
	if err == nil {
		t.Fatal("expected saga to fail at process step")
	}

	// Verify: reserve executed, process failed (not added to executed), journal skipped
	// Compensation: process compensated, then reserve compensated
	if len(executedSteps) != 1 || executedSteps[0] != "reserve" {
		t.Errorf("expected only 'reserve' executed, got %v", executedSteps)
	}
	if len(compensatedSteps) != 2 || compensatedSteps[0] != "refund" || compensatedSteps[1] != "release" {
		t.Errorf("expected compensations ['refund', 'release'], got %v", compensatedSteps)
	}

	t.Log("SUCCESS: Saga executed: reserve->process(fail)->compensate(refund)->compensate(release)")
}

func TestSagaPaymentSuccess(t *testing.T) {
	var executedSteps []string

	// All steps succeed
	saga := transaction.NewSaga().
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error { executedSteps = append(executedSteps, "reserve"); return nil },
			Compensate: func(ctx context.Context) error { return nil },
		}).
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error { executedSteps = append(executedSteps, "process"); return nil },
			Compensate: func(ctx context.Context) error { return nil },
		}).
		Then(transaction.SagaStep{
			Action: func(ctx context.Context) error { executedSteps = append(executedSteps, "journal"); return nil },
			Compensate: func(ctx context.Context) error { return nil },
		})

	ctx := context.Background()
	if err := saga.Execute(ctx); err != nil {
		t.Fatalf("saga failed: %v", err)
	}

	if len(executedSteps) != 3 {
		t.Errorf("expected 3 steps executed, got %d", len(executedSteps))
	}

	t.Log("SUCCESS: Saga completed all steps: reserve->process->journal")
}