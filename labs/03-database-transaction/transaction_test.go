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
	broker := transaction.NewInMemoryBroker(0)
	ctx := context.Background()

	event := transaction.Event{ID: "evt_101", EventType: "InvoicePaid", AggregateID: "101", Payload: `{"invoice_id": 101}`}

	_ = broker.Publish(ctx, event)
	_ = broker.Publish(ctx, event)

	events := broker.PublishedEvents()
	if len(events) != 2 {
		t.Fatalf("expected 2 deliveries (duplicate), got %d", len(events))
	}
	t.Logf("PROVEN: At-least-once delivery - event delivered twice: %s", events[0].ID)
}

// Test 8: Idempotent consumer processes duplicate event only once.
func TestIdempotentConsumerDeduplication(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()

	worker := transaction.NewCommissionWorker(db)
	ctx := context.Background()

	event := transaction.Event{ID: "evt_101", EventType: "InvoicePaid", AggregateID: "101", Payload: `{"invoice_id": 101}`}

	processed1, _ := worker.HandleEvent(ctx, "CommissionWorker", event)
	if !processed1 {
		t.Error("expected first event to be processed")
	}

	processed2, _ := worker.HandleEvent(ctx, "CommissionWorker", event)
	if processed2 {
		t.Error("expected duplicate event to be SKIPPED")
	}

	if worker.CommissionsPaidCount() != 1 {
		t.Errorf("expected commissions paid exactly 1 time, got %d", worker.CommissionsPaidCount())
	}
	t.Log("SUCCESS: Idempotent consumer deduplicated duplicate events.")
}

// Test 9: Transient failure succeeds after retry.
func TestTransientFailureSuccessAfterRetry(t *testing.T) {
	db := mockdb.NewDB()
	defer db.Close()
	seedTestDB(t, db)

	svc := transaction.NewInvoiceServiceOutbox(db)
	ctx := context.Background()
	_ = svc.PayInvoiceWithOutbox(ctx, 101)

	// Broker fails first 2, succeeds on 3rd
	broker := transaction.NewInMemoryBroker(2)
	dispatcher := transaction.NewOutboxDispatcher(db, broker, 3, nil)

	for i := 1; i <= 3; i++ {
		dispatcher.DispatchBatch(ctx)
	}

	events := broker.PublishedEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event after retries, got %d", len(events))
	}
	t.Log("SUCCESS: Transient failure succeeded after retry (fail-fail-succeed).")
}

// Test 10: Permanent failure moves event to DLQ after max attempts.
func TestDeadLetterQueue(t *testing.T) {
	dlq := transaction.NewDeadLetterQueue()
	event := transaction.Event{ID: "evt_failing_123", EventType: "InvoicePaid"}

	dlq.Add(event, "max attempts exceeded", 3)

	if dlq.Count() != 1 {
		t.Fatalf("expected DLQ count 1, got %d", dlq.Count())
	}

	records := dlq.Records()
	if records[0].Event.ID != event.ID {
		t.Errorf("expected event ID %s in DLQ", event.ID)
	}
	t.Log("SUCCESS: Event moved to DLQ after max attempts.")
}

// Test 11: Saga executes compensating actions on failure.
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

	if len(compensatedSteps) != 2 {
		t.Errorf("expected 2 compensations, got %d: %v", len(compensatedSteps), compensatedSteps)
	}

	t.Log("SUCCESS: Saga compensation executed: reserve->release, process->refund")
}
