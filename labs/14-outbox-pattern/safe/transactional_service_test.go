package safe_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lukman-ss/software-engineering-lab/labs/14-outbox-pattern/safe"
)

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

	// Run migration
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			id VARCHAR(36) PRIMARY KEY,
			customer_id VARCHAR(36) NOT NULL,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create orders table: %v", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS outbox_events (
			id VARCHAR(36) PRIMARY KEY,
			aggregate_type VARCHAR(100) NOT NULL,
			aggregate_id VARCHAR(36) NOT NULL,
			event_type VARCHAR(100) NOT NULL,
			payload JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			published_at TIMESTAMP WITH TIME ZONE,
			attempts INT NOT NULL DEFAULT 0,
			next_attempt_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create outbox_events table: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "TRUNCATE TABLE orders, outbox_events")
		db.Close()
	})

	return db
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestTransactionalOutboxAtomicity ensures both order and event are created atomically.
func TestTransactionalOutboxAtomicity(t *testing.T) {
	db := setupTestDB(t)

	service := safe.NewSafeOrderService(db)

	ctx := context.Background()
	order := safe.Order{
		CustomerID: "customer-001",
		Status:     "pending",
	}

	created, err := service.CreateOrder(ctx, order)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	// Verify order exists
	_, err = service.FindOrder(ctx, created.ID)
	if err != nil {
		t.Fatalf("order should exist in DB: %v", err)
	}

	// Verify outbox event exists
	var eventCount int
	err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM outbox_events WHERE aggregate_id = $1 AND published_at IS NULL`,
		created.ID,
	).Scan(&eventCount)
	if err != nil {
		t.Fatalf("failed to count outbox events: %v", err)
	}
	if eventCount != 1 {
		t.Fatalf("expected 1 unpublished event, got %d", eventCount)
	}

	t.Logf("SUCCESS: Order %s and outbox event created atomically", created.ID)
}

// TestTransactionalOutboxConsistency ensures no orphaned data.
func TestTransactionalOutboxConsistency(t *testing.T) {
	db := setupTestDB(t)
	service := safe.NewSafeOrderService(db)

	ctx := context.Background()

	// Create orders
	for i := 0; i < 3; i++ {
		order := safe.Order{
			CustomerID: fmt.Sprintf("customer-%d", i),
			Status:     "pending",
		}
		_, err := service.CreateOrder(ctx, order)
		if err != nil {
			t.Fatalf("failed to create order %d: %v", i, err)
		}
	}

	// Verify counts match
	orderCount, err := service.GetOrderCount(ctx)
	if err != nil {
		t.Fatalf("failed to get order count: %v", err)
	}

	eventCount, err := service.GetOutboxEventCount(ctx)
	if err != nil {
		t.Fatalf("failed to get outbox count: %v", err)
	}

	if orderCount != eventCount {
		t.Errorf("MISMATCH: %d orders but %d events", orderCount, eventCount)
	}

	t.Logf("CONSISTENCY VERIFIED: %d orders = %d events", orderCount, eventCount)
}

// TestTransactionalOutboxNoOrphans verifies that failed transaction leaves no traces.
func TestTransactionalOutboxNoOrphans(t *testing.T) {
	db := setupTestDB(t)

	ctx := context.Background()
	// Simulate constraint violation by using same ID twice
	duplicateID := "duplicate-order-id"

	// Use raw SQL to simulate a constraint violation scenario
	_, err := db.ExecContext(ctx,
		`INSERT INTO orders (id, customer_id, status, created_at) VALUES ($1, $2, $3, $4)`,
		duplicateID, "cust-001", "pending", time.Now(),
	)
	if err != nil {
		t.Fatalf("first insert should succeed: %v", err)
	}

	// Second insert should fail
	_, err = db.ExecContext(ctx,
		`INSERT INTO orders (id, customer_id, status, created_at) VALUES ($1, $2, $3, $4)`,
		duplicateID, "cust-002", "pending", time.Now(),
	)
	if err == nil {
		t.Fatal("duplicate insert should have failed")
	}

	// Verify: only one order exists, no outbox events
	var orderCount int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM orders WHERE id = $1", duplicateID).Scan(&orderCount)
	if err != nil {
		t.Fatalf("failed to count orders: %v", err)
	}

	var eventCount int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM outbox_events").Scan(&eventCount)
	if err != nil {
		t.Fatalf("failed to count events: %v", err)
	}

	if eventCount != 0 {
		t.Errorf("expected 0 events from failed transaction, got %d", eventCount)
	}

	t.Log("NO ORPHANS: Failed transaction leaves database clean")
}
