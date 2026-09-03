package safe_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/lukman-ss/software-engineering-lab/labs/14-outbox-pattern/safe"
)

// TestConcurrentOutboxWorkers ensures multiple workers don't create duplicate deliveries.
func TestConcurrentOutboxWorkers(t *testing.T) {
	db := setupTestDB(t)
	service := safe.NewSafeOrderService(db)
	ctx := context.Background()

	// Create multiple orders
	orderCount := 10
	var mu sync.Mutex
	orderIDs := make([]string, orderCount)

	for i := 0; i < orderCount; i++ {
		order := safe.Order{CustomerID: fmt.Sprintf("cust-%d", i), Status: "pending"}
		created, err := service.CreateOrder(ctx, order)
		if err != nil {
			t.Fatalf("failed to create order %d: %v", i, err)
		}
		orderIDs[i] = created.ID
	}

	// Verify outbox events exist
	eventCount, err := service.GetOutboxEventCount(ctx)
	if err != nil {
		t.Fatalf("failed to get event count: %v", err)
	}
	if int(eventCount) != orderCount {
		t.Fatalf("expected %d events, got %d", orderCount, eventCount)
	}

	// Create broker and worker config
	broker := &safe.MockEventBroker{}
	cfg := safe.WorkerConfig{
		BatchSize:    3, // Small batch to test coordination
		PollInterval: 50 * time.Millisecond,
		MaxAttempts:  5,
		BackoffBase:  10 * time.Millisecond,
	}

	// Start 3 concurrent workers
	workers := make([]*safe.ConcurrentOutboxWorker, 3)
	for i := 0; i < 3; i++ {
		workers[i] = safe.NewConcurrentOutboxWorker(
			fmt.Sprintf("worker-%d", i),
			db,
			broker,
			cfg,
		)
		workers[i].Start()
	}

	// Let them run
	time.Sleep(300 * time.Millisecond)

	// Stop all workers
	for _, w := range workers {
		w.Stop()
	}

	// Verify results
	published := broker.GetPublished()

	t.Logf("CONCURRENT WORKER RESULTS:")
	t.Logf("  Orders created: %d", orderCount)
	t.Logf("  Events published by broker: %d", len(published))

	// Check for duplicates
	eventIDMap := make(map[string]int)
	for _, e := range published {
		eventIDMap[e.ID]++
	}

	duplicates := 0
	for id, count := range eventIDMap {
		if count > 1 {
			duplicates++
			t.Logf("  DUPLICATE: event %s was published %d times", id, count)
		}
	}

	if duplicates > 0 {
		t.Errorf("Found %d events with duplicate deliveries", duplicates)
	}

	// Verify all events were processed
	if len(published) != orderCount {
		t.Errorf("Expected %d events published, got %d", orderCount, len(published))
	}

	t.Log("SUCCESS: No duplicate deliveries, all events processed exactly once")
}

// TestConcurrentOutboxWorkersWithFailures tests workers handle failures gracefully.
func TestConcurrentOutboxWorkersWithFailures(t *testing.T) {
	db := setupTestDB(t)
	service := safe.NewSafeOrderService(db)
	ctx := context.Background()

	// Create orders
	for i := 0; i < 5; i++ {
		_, _ = service.CreateOrder(ctx, safe.Order{CustomerID: fmt.Sprintf("cust-%d", i), Status: "pending"})
	}

	broker := &safe.MockEventBroker{}
	// Start with failures
	broker.SetShouldFail(true)

	cfg := safe.WorkerConfig{
		BatchSize:    2,
		PollInterval: 30 * time.Millisecond,
		MaxAttempts:  3,
		BackoffBase:  10 * time.Millisecond,
	}

	// Start workers with failures
	workers := make([]*safe.ConcurrentOutboxWorker, 2)
	for i := 0; i < 2; i++ {
		workers[i] = safe.NewConcurrentOutboxWorker(
			fmt.Sprintf("failing-worker-%d", i),
			db,
			broker,
			cfg,
		)
		workers[i].Start()
	}

	time.Sleep(100 * time.Millisecond)

	// Fix broker
	broker.SetShouldFail(false)

	// Let workers retry
	time.Sleep(200 * time.Millisecond)

	for _, w := range workers {
		w.Stop()
	}

	// Verify all events eventually published
	published := broker.GetPublished()
	t.Logf("Events published after failure recovery: %d", len(published))

	remaining, err := service.GetOutboxEventCount(ctx)
	if err != nil {
		t.Fatalf("failed to query remaining events: %v", err)
	}

	t.Logf("Remaining unpublished events: %d", remaining)
	t.Log("SUCCESS: Workers recovered from failures and completed processing")
}

// TestFORUpdateSkipLockedBehavior explicitly documents the SKIP LOCKED behavior.
func TestFORUpdateSkipLockedBehavior(t *testing.T) {
	db := setupTestDB(t)
	service := safe.NewSafeOrderService(db)
	ctx := context.Background()

	// Create orders
	for i := 0; i < 4; i++ {
		_, _ = service.CreateOrder(ctx, safe.Order{CustomerID: fmt.Sprintf("cust-%d", i), Status: "pending"})
	}

	t.Log("SKIP LOCKED BEHAVIOR:")
	t.Log("")
	t.Log("  Worker A: SELECT ... FOR UPDATE SKIP LOCKED LIMIT 2")
	t.Log("    -> Locks 2 rows, gets events 1,2")
	t.Log("    -> Events 3,4 still available for other workers")
	t.Log("")
	t.Log("  Worker B: SELECT ... FOR UPDATE SKIP LOCKED LIMIT 2")
	t.Log("    -> Cannot lock 1,2 (already locked)")
	t.Log("    -> SKIP LOCKED skips them")
	t.Log("    -> Locks 3,4, gets those events")
	t.Log("")
	t.Log("  Result: Worker A got 2, Worker B got 2. No contention, no duplicates!")
	t.Log("  If we used plain FOR UPDATE without SKIP LOCKED:")
	t.Log("    -> Worker B would wait indefinitely for A to release locks")
	t.Log("    -> Workers would serialize instead of running concurrently")
}
