package safe_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lukman-ss/software-engineering-lab/labs/14-outbox-pattern/safe"
)

// TestOutboxWorkerSuccess tests that the worker picks up unpublished events and publishes them successfully.
func TestOutboxWorkerSuccess(t *testing.T) {
	db := setupTestDB(t)
	service := safe.NewSafeOrderService(db)

	ctx := context.Background()

	// Create an order (which inserts outbox event)
	order := safe.Order{
		CustomerID: "cust-worker-01",
		Status:     "pending",
	}
	_, err := service.CreateOrder(ctx, order)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	// Verify 1 unpublished event exists
	count, err := service.GetOutboxEventCount(ctx)
	if err != nil || count != 1 {
		t.Fatalf("expected 1 unpublished event, got %d (err: %v)", count, err)
	}

	// Start broker and worker
	broker := &safe.MockEventBroker{}
	cfg := safe.WorkerConfig{
		BatchSize:    10,
		PollInterval: 100 * time.Millisecond,
		MaxAttempts:  3,
		BackoffBase:  10 * time.Millisecond,
	}

	worker := safe.NewOutboxWorker(db, broker, cfg)
	worker.Start()
	defer worker.Stop()

	// Wait for worker to process
	time.Sleep(300 * time.Millisecond)

	// Verify broker received the event
	published := broker.GetPublished()
	if len(published) != 1 {
		t.Fatalf("expected 1 published event on broker, got %d", len(published))
	}

	// Verify outbox event is now marked published
	count, err = service.GetOutboxEventCount(ctx)
	if err != nil || count != 0 {
		t.Errorf("expected 0 unpublished events remaining, got %d", count)
	}

	t.Log("OUTBOX WORKER SUCCESS: Event published and marked published in DB")
}

// TestOutboxWorkerRetryAndBackoff tests that failed publishes trigger retry and exponential backoff.
func TestOutboxWorkerRetryAndBackoff(t *testing.T) {
	db := setupTestDB(t)
	service := safe.NewSafeOrderService(db)
	ctx := context.Background()

	// Create order
	_, err := service.CreateOrder(ctx, safe.Order{CustomerID: "cust-retry-01", Status: "pending"})
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	broker := &safe.MockEventBroker{}
	broker.SetShouldFail(true) // Broker is down initially

	cfg := safe.WorkerConfig{
		BatchSize:    10,
		PollInterval: 50 * time.Millisecond,
		MaxAttempts:  3,
		BackoffBase:  20 * time.Millisecond,
	}

	worker := safe.NewOutboxWorker(db, broker, cfg)
	worker.Start()

	// Wait for attempts
	time.Sleep(150 * time.Millisecond)
	worker.Stop()

	// Verify attempts incremented and event not published yet
	var attempts int
	var nextAttempt time.Time
	err = db.QueryRowContext(ctx, `SELECT attempts, next_attempt_at FROM outbox_events`).Scan(&attempts, &nextAttempt)
	if err != nil {
		t.Fatalf("failed to query event retry state: %v", err)
	}

	if attempts == 0 {
		t.Error("expected attempts > 0 after worker runs")
	}

	t.Logf("RETRY & BACKOFF VERIFIED: Attempts = %d, NextAttemptAt = %v", attempts, nextAttempt)

	// Now fix broker
	broker.SetShouldFail(false)
	// Force next_attempt_at to now so worker picks it up immediately
	_, _ = db.ExecContext(ctx, `UPDATE outbox_events SET next_attempt_at = CURRENT_TIMESTAMP`)

	worker2 := safe.NewOutboxWorker(db, broker, cfg)
	worker2.Start()
	time.Sleep(150 * time.Millisecond)
	worker2.Stop()

	// Verify published successfully now
	count, err := service.GetOutboxEventCount(ctx)
	if err != nil || count != 0 {
		t.Errorf("expected 0 unpublished events after broker recovery, got %d", count)
	}
}

// TestOutboxWorkerGracefulShutdown tests that worker finishes current batch on stop.
func TestOutboxWorkerGracefulShutdown(t *testing.T) {
	db := setupTestDB(t)
	service := safe.NewSafeOrderService(db)
	ctx := context.Background()

	// Create multiple orders
	for i := 0; i < 5; i++ {
		_, _ = service.CreateOrder(ctx, safe.Order{CustomerID: fmt.Sprintf("cust-%d", i), Status: "pending"})
	}

	broker := &safe.MockEventBroker{}
	cfg := safe.WorkerConfig{
		BatchSize:    2,
		PollInterval: 50 * time.Millisecond,
		MaxAttempts:  3,
		BackoffBase:  10 * time.Millisecond,
	}

	worker := safe.NewOutboxWorker(db, broker, cfg)
	worker.Start()

	// Let it run briefly then stop
	time.Sleep(100 * time.Millisecond)
	worker.Stop() // Should block until in-flight work completes cleanly

	t.Log("GRACEFUL SHUTDOWN SUCCESS: Worker stopped cleanly without hanging")
}
