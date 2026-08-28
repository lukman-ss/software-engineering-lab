package safe

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

// Consumer simulates a downstream message consumer that processes outbox events.
// Because the outbox pattern guarantees AT-LEAST-ONCE delivery (e.g. worker crashes
// after publishing to broker but before marking `published_at`, or network partition),
// consumers WILL receive duplicate events.
//
// Therefore, the consumer MUST be IDEMPOTENT.
type IdempotentConsumer struct {
	db          *sql.DB
	processedMu sync.Mutex
	processed   map[string]int // Tracks event ID processing count in memory for testing
}

// NewIdempotentConsumer creates a new consumer.
func NewIdempotentConsumer(db *sql.DB) *IdempotentConsumer {
	return &IdempotentConsumer{
		db:        db,
		processed: make(map[string]int),
	}
}

// ProcessEvent processes an incoming event idempotently using a deduplication table
// or unique constraint on the consumer side.
func (c *IdempotentConsumer) ProcessEvent(ctx context.Context, eventID string, eventType string, aggregateID string) error {
	c.processedMu.Lock()
	c.processed[eventID]++
	c.processedMu.Unlock()

	// Idempotency check via consumer_processed_events table:
	// INSERT INTO consumer_processed_events (event_id, processed_at) VALUES ($1, NOW())
	// ON CONFLICT (event_id) DO NOTHING;
	// If no rows inserted, it's a duplicate! We acknowledge and ignore successfully.

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	// Try inserting into processed events log
	var inserted bool
	err = tx.QueryRowContext(ctx, `
		INSERT INTO consumer_processed_events (event_id, aggregate_id, event_type, processed_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (event_id) DO NOTHING
		RETURNING true
	`, eventID, aggregateID, eventType).Scan(&inserted)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to check idempotency log: %w", err)
	}

	if !inserted {
		// Duplicate delivery detected! Skip processing business logic.
		// fmt.Printf("[Consumer] Duplicate event %s detected, skipping.\n", eventID)
		_ = tx.Commit()
		return nil
	}

	// Perform actual business logic (e.g. fulfill order, provision inventory)
	// fmt.Printf("[Consumer] Processing event %s for order %s\n", eventID, aggregateID)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit consumer transaction: %w", err)
	}

	return nil
}

// GetProcessCount returns how many times an event was passed to the consumer.
func (c *IdempotentConsumer) GetProcessCount(eventID string) int {
	c.processedMu.Lock()
	defer c.processedMu.Unlock()
	return c.processed[eventID]
}
