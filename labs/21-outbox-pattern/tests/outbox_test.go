package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/software-engineering-lab/labs/21-outbox-pattern/internal/outbox"
)

func TestTransactionalOutbox_HappyPath(t *testing.T) {
	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	relay := outbox.NewRelay(db, broker, 10*time.Millisecond)
	service := outbox.NewOrderService(db)
	consumer := outbox.NewConsumer()

	relay.Start()
	defer relay.Stop()

	// 1. Service writes order and outbox atomically
	err := service.CreateOrderWithOutbox("o-1", "c-1", 100.50)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify DB state before relay processes it
	order, ok := db.GetOrder("o-1")
	if !ok || order.Amount != 100.50 {
		t.Fatalf("expected order to be persisted in db")
	}

	// 2. Wait for relay to poll and dispatch
	time.Sleep(50 * time.Millisecond)

	// Verify broker received
	published := broker.GetPublished()
	if len(published) != 1 {
		t.Fatalf("expected 1 message published to broker, got %d", len(published))
	}
	if published[0].EventType != "OrderCreated" {
		t.Fatalf("expected OrderCreated event")
	}

	// Verify outbox status updated to processed
	msg, ok := db.GetOutbox("evt-o-1")
	if !ok || msg.Status != outbox.MessageStatusProcessed {
		t.Fatalf("expected outbox message to be marked processed")
	}

	// 3. Consumer handles idempotently
	processed := consumer.Handle(published[0])
	if !processed {
		t.Fatalf("expected consumer to process new message")
	}
	if consumer.GetReceivedCount() != 1 {
		t.Fatalf("expected consumer to have 1 message")
	}
}

func TestTransactionalOutbox_Rollback(t *testing.T) {
	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	relay := outbox.NewRelay(db, broker, 10*time.Millisecond)

	relay.Start()
	defer relay.Stop()

	tx := db.BeginTx()
	_ = tx.SaveOrder(outbox.Order{ID: "o-2", Amount: 50})
	_ = tx.SaveOutbox(outbox.OutboxMessage{ID: "evt-o-2", Status: outbox.MessageStatusPending})

	// Simulate error and rollback
	_ = tx.Rollback()

	// Verify nothing was persisted
	_, ok := db.GetOrder("o-2")
	if ok {
		t.Fatalf("expected order to not be saved")
	}
	_, ok = db.GetOutbox("evt-o-2")
	if ok {
		t.Fatalf("expected outbox to not be saved")
	}

	time.Sleep(30 * time.Millisecond)
	if len(broker.GetPublished()) > 0 {
		t.Fatalf("expected no messages sent to broker")
	}
}

func TestTransactionalOutbox_Idempotency_DuplicateDelivery(t *testing.T) {
	consumer := outbox.NewConsumer()

	msg := outbox.OutboxMessage{
		ID:        "evt-duplicate",
		EventType: "OrderCreated",
		Payload:   "data",
	}

	p1 := consumer.Handle(msg)
	if !p1 {
		t.Fatalf("expected first delivery to be processed")
	}

	// Simulate relay duplicate delivery (e.g., relay crashed before marking processed)
	p2 := consumer.Handle(msg)
	if p2 {
		t.Fatalf("expected duplicate delivery to be rejected by consumer")
	}

	if consumer.GetReceivedCount() != 1 {
		t.Fatalf("expected exactly 1 message processed despite duplicate delivery")
	}
}

func TestDualWriteProblem_Failure(t *testing.T) {
	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	service := outbox.NewOrderService(db)

	broker.SetFailNext(true) // simulate broker down

	err := service.CreateOrderDualWriteNaive(broker, "o-bug", "c-2", 200)
	if err == nil {
		t.Fatalf("expected error from dual write naive approach")
	}

	// Assert dual-write inconsistency
	_, ok := db.GetOrder("o-bug")
	if !ok {
		t.Fatalf("expected order to be saved in DB")
	}

	if len(broker.GetPublished()) != 0 {
		t.Fatalf("expected no message published to broker")
	}
	// System is now inconsistent. Order exists, event never sent.
}

func TestTransactionalOutbox_ConcurrentWrites(t *testing.T) {
	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	relay := outbox.NewRelay(db, broker, 5*time.Millisecond)
	service := outbox.NewOrderService(db)

	relay.Start()
	defer relay.Stop()

	var wg sync.WaitGroup
	workers := 10
	ordersPerWorker := 10
	totalExpected := workers * ordersPerWorker

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(wID int) {
			defer wg.Done()
			for j := 0; j < ordersPerWorker; j++ {
				orderID := fmt.Sprintf("o-concurrent-%d-%d", wID, j)
				err := service.CreateOrderWithOutbox(
					orderID,
					"c-multi",
					float64(wID*100+j),
				)
				if err != nil {
					t.Errorf("failed to create order: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()

	// Wait for relay to process all orders
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if len(broker.GetPublished()) == totalExpected {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	published := broker.GetPublished()
	if len(published) != totalExpected {
		t.Fatalf("expected %d messages published, got %d", totalExpected, len(published))
	}

	pending := db.GetPendingOutbox()
	if len(pending) != 0 {
		t.Fatalf("expected 0 pending messages, got %d", len(pending))
	}
}

func TestTransactionalOutbox_PurgeProcessed(t *testing.T) {
	db := outbox.NewDB()
	tx := db.BeginTx()
	_ = tx.SaveOutbox(outbox.OutboxMessage{ID: "m1", Status: outbox.MessageStatusPending})
	_ = tx.SaveOutbox(outbox.OutboxMessage{ID: "m2", Status: outbox.MessageStatusProcessed})
	_ = tx.Commit()

	purged := db.PurgeProcessedOutbox()
	if purged != 1 {
		t.Fatalf("expected 1 record purged, got %d", purged)
	}

	_, ok := db.GetOutbox("m1")
	if !ok {
		t.Fatalf("expected pending message to remain")
	}

	_, ok = db.GetOutbox("m2")
	if ok {
		t.Fatalf("expected processed message to be purged")
	}
}

func TestTransactionalOutbox_RelayRetryAfterBrokerFailure(t *testing.T) {
	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	relay := outbox.NewRelay(db, broker, 10*time.Millisecond)
	service := outbox.NewOrderService(db)
	consumer := outbox.NewConsumer()

	broker.SetFailNext(true)

	err := service.CreateOrderWithOutbox("o-retry", "c-retry", 55.0)
	if err != nil {
		t.Fatalf("expected order creation to succeed: %v", err)
	}

	relay.Start()
	relay.Start() // verify idempotency of Start()
	defer func() {
		relay.Stop()
		relay.Stop() // verify idempotency of Stop()
	}()

	// First poll fails due to SetFailNext
	time.Sleep(25 * time.Millisecond)

	msg, ok := db.GetOutbox("evt-o-retry")
	if !ok {
		t.Fatalf("expected outbox message to exist")
	}

	// Wait for subsequent successful poll
	time.Sleep(50 * time.Millisecond)

	published := broker.GetPublished()
	if len(published) != 1 {
		t.Fatalf("expected 1 message published after retry, got %d", len(published))
	}

	msg, ok = db.GetOutbox("evt-o-retry")
	if !ok || msg.Status != outbox.MessageStatusProcessed {
		t.Fatalf("expected message to be marked processed after retry")
	}

	if !consumer.Handle(published[0]) {
		t.Fatalf("expected consumer to process retried message")
	}
}

func TestTransactionalOutbox_ConcurrentConsumers(t *testing.T) {
	consumer := outbox.NewConsumer()
	var wg sync.WaitGroup
	workers := 10
	messagesPerWorker := 20

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(wID int) {
			defer wg.Done()
			for j := 0; j < messagesPerWorker; j++ {
				msg := outbox.OutboxMessage{
					ID:        fmt.Sprintf("evt-shared-%d", j), // Overlapping IDs to test concurrent deduplication
					EventType: "OrderCreated",
					Payload:   "test-data",
				}
				consumer.Handle(msg)
			}
		}(i)
	}

	wg.Wait()

	// Only messagesPerWorker unique IDs should have been accepted
	if consumer.GetReceivedCount() != messagesPerWorker {
		t.Fatalf("expected exactly %d unique messages processed, got %d", messagesPerWorker, consumer.GetReceivedCount())
	}
}
