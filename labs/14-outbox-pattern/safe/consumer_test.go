package safe_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/lukman-ss/software-engineering-lab/labs/14-outbox-pattern/safe"
)

// setupConsumerTestDB initializes the database for consumer tests.
func setupConsumerTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := "postgres://postgres:postgres@localhost:5432/se_lab?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Skipf("database not available: %v", err)
	}

	// Create consumer_processed_events table for idempotency deduplication
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS consumer_processed_events (
			event_id VARCHAR(36) PRIMARY KEY,
			aggregate_id VARCHAR(36) NOT NULL,
			event_type VARCHAR(100) NOT NULL,
			processed_at TIMESTAMP WITH TIME ZONE NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("failed to create consumer_processed_events table: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "TRUNCATE TABLE orders, outbox_events, consumer_processed_events")
		db.Close()
	})

	return db
}

// TestAtLeastOnceDelivery documents that outbox still produces duplicates.
func TestAtLeastOnceDelivery(t *testing.T) {
	t.Log("AT-LEAST-ONCE DELIVERY REALITY:")
	t.Log("")
	t.Log("The outbox pattern provides AT-LEAST-ONCE delivery semantics.")
	t.Log("This means:")
	t.Log("  - Every event WILL be delivered at least once")
	t.Log("  - Events MAY be delivered MULTIPLE times")
	t.Log("")
	t.Log("Reasons for duplicates:")
	t.Log("  1. Worker publishes event, crashes BEFORE marking published_at")
	t.Log("     -> Event remains in outbox -> retry publishes again")
	t.Log("")
	t.Log("  2. Broker returns success, but ACK was lost")
	t.Log("     -> Consumer receives duplicate from broker")
	t.Log("")
	t.Log("  3. Network partition causes duplicate message delivery")
	t.Log("     -> Consumer sees duplicate")
	t.Log("")
	t.Log("  4. Consumer processes event, crashes before committing position")
	t.Log("     -> Rebalancing causes re-delivery of same event")
	t.Log("")
	t.Log("CONCLUSION: Consumers MUST be idempotent.")
}

// TestIdempotentConsumerSuccess verifies idempotent processing.
func TestIdempotentConsumerSuccess(t *testing.T) {
	db := setupConsumerTestDB(t)
	consumer := safe.NewIdempotentConsumer(db)
	ctx := context.Background()

	eventID := "test-event-001"
	aggregateID := "order-001"
	eventType := "OrderCreated"

	// Process same event multiple times
	for i := 0; i < 3; i++ {
		err := consumer.ProcessEvent(ctx, eventID, eventType, aggregateID)
		if err != nil {
			t.Fatalf("unexpected error processing event: %v", err)
		}
	}

	// Verify idempotency log shows only 1 unique entry
	var count int
	err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM consumer_processed_events WHERE event_id = $1`,
		eventID,
	).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 unique entry, got %d", count)
	}

	t.Logf("IDEMPOTENCY VERIFIED: Event %s processed %d times but only 1 deduplication entry created",
		eventID, consumer.GetProcessCount(eventID))
}

// TestIdempotentConsumerWithRealOutboxIntegration tests full flow.
func TestIdempotentConsumerWithRealOutboxIntegration(t *testing.T) {
	db := setupConsumerTestDB(t)
	service := safe.NewSafeOrderService(db)
	consumer := safe.NewIdempotentConsumer(db)
	ctx := context.Background()

	// Create order (inserts into orders and outbox_events)
	order := safe.Order{CustomerID: "customer-001", Status: "pending"}
	created, err := service.CreateOrder(ctx, order)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	// Simulate consumer processing the same event multiple times
	// (as would happen with at-least-once delivery)
	eventID := created.ID // In real scenario, this would be the outbox event ID
	evtType := "OrderCreated"

	for i := 0; i < 5; i++ {
		err := consumer.ProcessEvent(ctx, eventID+"-evt", evtType, created.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// Verify only 1 unique event processed
	var dupCount int
	err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM consumer_processed_events WHERE aggregate_id = $1`,
		created.ID,
	).Scan(&dupCount)
	if err != nil {
		t.Fatalf("failed to count: %v", err)
	}

	if dupCount != 1 {
		t.Errorf("expected 1 deduplicated entry, got %d", dupCount)
	}

	t.Log("INTEGRATION TEST PASSED: Consumer handled duplicate processing gracefully")
}

// TestDeduplicationStrategies documents different approaches.
func TestDeduplicationStrategies(t *testing.T) {
	t.Log("DEDUPLICATION STRATEGIES:")
	t.Log("")
	t.Log("Strategy 1: Database unique constraint")
	t.Log("  INSERT INTO consumer_events ... ON CONFLICT (event_id) DO NOTHING")
	t.Log("  Pros: Simple, reliable, works with any event source")
	t.Log("  Cons: Requires write access, storage overhead")
	t.Log("")
	t.Log("  PostgreSQL:")
	t.Log("    CREATE TABLE processed_events (event_id VARCHAR(36) PRIMARY KEY, processed_at TIMESTAMP);")
	t.Log("    INSERT INTO processed_events VALUES ($1, NOW()) ON CONFLICT DO NOTHING;")
	t.Log("")
	t.Log("Strategy 2: Pre-check with existence query")
	t.Log("  SELECT 1 FROM processed_events WHERE event_id = $1;")
	t.Log("  INSERT INTO processed_events ... (if not exists)")
	t.Log("  Pros: Can skip processing before expensive operations")
	t.Log("  Cons: Still needs transaction, race condition if not serialized")
	t.Log("")
	t.Log("Strategy 3: Application-level deduplication cache")
	t.Log("  Use Redis SET with TTL: SETEX processed:{event_id} 86400 1")
	t.Log("  Pros: Fast, works with Kafka offsets, stateless workers")
	t.Log("  Cons: Cache miss can still cause duplicate write")
	t.Log("")
	t.Log("Strategy 4: Idempotent operations only (best practice)")
	t.Log("  Design operations so: process(same_event) == process(different_events)")
	t.Log("  Example: INSERT INTO orders ... ON CONFLICT DO NOTHING")
	t.Log("  Pros: No deduplication overhead, naturally concurrent")
	t.Log("  Cons: Requires careful domain modeling, not always possible")
	t.Log("")
	t.Log("RECOMMENDATION: Use Strategy 1 (DB unique constraint) as primary defense,")
	t.Log("            combined with idempotent operations where possible.")
}

// TestWhatStillCanFail documents remaining failure modes.
func TestWhatStillCanFail(t *testing.T) {
	t.Log("AT-LIKE-ONCE REALITY - WHAT STILL CAN FAIL:")
	t.Log("")
	t.Log("Even with outbox and idempotent consumers, these failures remain:")
	t.Log("")
	t.Log("1. COMPLETE MESSAGE BROKER OUTAGE")
	t.Log("   - Outbox workers queue up events in DB")
	t.Log("   - Consumer cannot receive events")
	t.Log("   - Mitigation: Dead letter queue after max attempts")
	t.Log("")
	t.Log("2. ORCHESTRATOR/LEADER ELECTION FAILURES")
	t.Log("   - Multiple workers process same batch before SKIP LOCKED")
	t.Log("   - Race condition edge cases in high-load scenarios")
	t.Log("   - Mitigation: Higher isolation levels, careful testing")
	t.Log("")
	t.Log("3. CONSUMER LOGIC BUGS")
	t.Log("   - Deduplication works, but business logic has bugs")
	t.Log("   - Example: Double-decrementing a counter due to race")
	t.Log("   - Mitigation: Idempotent operations, integration tests")
	t.Log("")
	t.Log("4. EVENT SCHEMA EVOLUTION")
	t.Log("   - Consumer expects schema v2, receives v1")
	t.Log("   - Breaking changes cause silent failures")
	t.Log("   - Mitigation: Schema versioning, backward compatibility")
	t.Log("")
	t.Log("5. DEAD LETTER QUEUE PROCESSING")
	t.Log("   - Events that exceed max_attempts need manual intervention")
	t.Log("   - Forgotten DLQ events cause data gaps")
	t.Log("   - Mitigation: Monitoring, alerts, periodic review")
	t.Log("")
	t.Log("6. CLOCK DRIFT IN RETRY SCHEDULING")
	t.Log("   - next_attempt_at computed from wall clock")
	t.Log("   - If clock jumps, events may retry too soon/slowly")
	t.Log("   - Mitigation: Use monotonic time, NTP sync")
	t.Log("")
	t.Log("7. OUTBOX BACKUP/RELATIONSHIP ACCURACY")
	t.Log("   - Outbox is a backup log, not primary data source")
	t.Log("   - If primary data is corrupted, outbox alone cannot restore it")
	t.Log("   - Mitigation: Regular backups, audits")
}
