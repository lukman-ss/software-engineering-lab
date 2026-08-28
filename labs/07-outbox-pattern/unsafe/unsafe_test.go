package unsafe_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/stdlib"

	_ "github.com/lukman-ss/software-engineering-lab/labs/07-outbox-pattern/unsafe/order"
	"github.com/lukman-ss/software-engineering-lab/labs/07-outbox-pattern/unsafe"
)

// setupTestDB creates a test database connection.
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "se_lab"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Skipf("database not available: %v", err)
	}

	// Create orders table if not exists
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			id VARCHAR(36) PRIMARY KEY,
			customer_id VARCHAR(36) NOT NULL,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "TRUNCATE TABLE orders")
		db.Close()
	})

	return db
}

// getEnv reads an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestDualWriteProblem demonstrates the dual write problem:
// Order is created in DB but event publishing fails -> inconsistent state.
func TestDualWriteProblem(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a publisher that fails after first successful publish
	publisher := &unsafe.MockEventPublisher{
		FailAfter: 0, // Fail immediately on first publish
	}

	service := unsafe.NewUnsafeOrderService(nil, publisher, db)

	ctx := context.Background()
	order := unsafe.Order{
		ID:         "order-001",
		CustomerID: "customer-001",
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	// Attempt to create order
	_, err := service.CreateOrder(ctx, order)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Check database: order was inserted BEFORE event publish failed
	// This is the dual write problem - order exists but event was never sent
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM orders WHERE id = $1", order.ID).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count orders: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 orphaned order in DB, got %d", count)
	}

	// Verify no events were published
	events := publisher.GetPublishedEvents()
	if len(events) != 0 {
		t.Fatalf("expected 0 published events, got %d", len(events))
	}

	// The order exists in DB but consumers never received the event
	t.Log("DUAL WRITE PROBLEM DEMONSTRATED:")
	t.Logf("  - Order %s exists in database (count: %d)", order.ID, count)
	t.Logf("  - No events were published (count: %d)", len(events))
	t.Log("  - System is now inconsistent: order exists but consumers don't know about it")
}

// TestDualWriteProblemPartialFailure demonstrates partial failure.
func TestDualWriteProblemPartialFailure(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Publisher fails after 2 successful publishes
	publisher := &unsafe.MockEventPublisher{
		FailAfter: 2,
	}

	service := unsafe.NewUnsafeOrderService(nil, publisher, db)

	ctx := context.Background()

	// Create 3 orders
	orders := []unsafe.Order{
		{ID: "order-001", CustomerID: "customer-001", Status: "pending", CreatedAt: time.Now()},
		{ID: "order-002", CustomerID: "customer-002", Status: "pending", CreatedAt: time.Now()},
		{ID: "order-003", CustomerID: "customer-003", Status: "pending", CreatedAt: time.Now()},
	}

	successCount := 0
	for i, order := range orders {
		_, err := service.CreateOrder(ctx, order)
		if err != nil {
			t.Logf("Order %s failed: %v", order.ID, err)
			// Check if order was still created in DB
			var count int
			err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM orders WHERE id = $1", order.ID).Scan(&count)
			if err != nil {
				t.Fatalf("failed to count orders: %v", err)
			}
			if count > 0 {
				t.Logf("  BUT order %s exists in DB (orphaned!)", order.ID)
			}
		} else {
			successCount++
		}
	}

	// Verify events published
	events := publisher.GetPublishedEvents()
	t.Logf("Successfully created: %d orders, %d events published", successCount, len(events))

	// Count total orders in DB
	var totalOrders int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM orders").Scan(&totalOrders)
	if err != nil {
		t.Fatalf("failed to count orders: %v", err)
	}

	// The problem: total orders != successfully published events
	if totalOrders != len(events) {
		t.Logf("INCONSISTENCY DETECTED:")
		t.Logf("  - Orders in DB: %d", totalOrders)
		t.Logf("  - Events published: %d", len(events))
		t.Logf("  - Difference: %d orders have no corresponding events", totalOrders-len(events))
	}
}

// TestDualWriteProblemEventFirst shows the reverse problem.
func TestDualWriteProblemEventFirst(t *testing.T) {
	// If we publish event FIRST then insert to DB:
	// - Event succeeds, DB insert fails -> event sent for non-existent order
	// This is equally problematic
	t.Log("REVERSE PROBLEM:")
	t.Log("  If we publish event FIRST then insert to DB:")
	t.Log("  - Event succeeds, DB insert fails -> consumers see event for non-existent order")
	t.Log("  This is equally problematic")
}

// TestDualWriteProblemNoAtomicity demonstrates lack of atomicity.
func TestDualWriteProblemNoAtomicity(t *testing.T) {
	t.Log("ROOT CAUSE ANALYSIS:")
	t.Log("")
	t.Log("The dual write problem occurs because:")
	t.Log("  1. Database and message queue are separate systems")
	t.Log("  2. No distributed transaction (2PC/XA) is used")
	t.Log("  3. Operations execute sequentially, not atomically")
	t.Log("")
	t.Log("Possible failure scenarios:")
	t.Log("  Scenario 1: DB commit succeeds, event publish fails")
	t.Log("    -> Order exists, but no event sent (INCONSISTENT)")
	t.Log("")
	t.Log("  Scenario 2: Event publish succeeds, DB commit fails")
	t.Log("    -> Event sent, but order doesn't exist (INCONSISTENT)")
	t.Log("")
	t.Log("  Scenario 3: Both succeed (happy path)")
	t.Log("    -> Consistent state, but achieved by chance")
	t.Log("")
	t.Log("  Scenario 4: Both fail")
	t.Log("    -> Inconsistent state, but at least not misleading")
}

// TestDualWriteProblemWithRetry shows retry doesn't help.
func TestDualWriteProblemWithRetry(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Publisher that always fails
	publisher := &unsafe.MockEventPublisher{
		FailAfter: -1, // Always fail
	}

	service := unsafe.NewUnsafeOrderService(nil, publisher, db)

	ctx := context.Background()
	order := unsafe.Order{
		ID:         "order-retry-001",
		CustomerID: "customer-001",
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	// Even with retries, the order will be orphaned
	for attempt := 1; attempt <= 3; attempt++ {
		_, err := service.CreateOrder(ctx, order)
		if err != nil {
			t.Logf("Attempt %d failed: %v", attempt, err)
		}
	}

	// Check DB
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM orders WHERE id = $1", order.ID).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count orders: %v", err)
	}

	if count > 0 {
		t.Logf("RETRY DOESN'T HELP:")
		t.Logf("  Order %s was created %d times in DB", order.ID, count)
		t.Logf("  No events were ever published")
		t.Logf("  Retry only makes the problem worse (duplicate orders!)")
	}
}

// MockStdlibDriver is a no-op driver registration for tests.
func init() {
	// Register the pgx stdlib driver
}
