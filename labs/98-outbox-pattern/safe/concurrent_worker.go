package safe

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// ConcurrentOutboxWorker is an outbox worker designed to run in multiple instances
// concurrently without stepping on each other's toes, achieved via row-level locking.
type ConcurrentOutboxWorker struct {
	*OutboxWorker
	workerID string
}

// NewConcurrentOutboxWorker creates an outbox worker safe for multiple instances.
func NewConcurrentOutboxWorker(id string, db *sql.DB, publisher EventPublisher, cfg WorkerConfig) *ConcurrentOutboxWorker {
	baseWorker := NewOutboxWorker(db, publisher, cfg)
	return &ConcurrentOutboxWorker{
		OutboxWorker: baseWorker,
		workerID:     id,
	}
}

// Start runs the concurrent worker.
func (w *ConcurrentOutboxWorker) Start() {
	w.wg.Add(1)
	go w.run()
}

func (w *ConcurrentOutboxWorker) run() {
	defer w.wg.Done()

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	log.Printf("[ConcurrentOutboxWorker %s] Started polling every %v", w.workerID, w.pollInterval)

	for {
		select {
		case <-w.ctx.Done():
			log.Printf("[ConcurrentOutboxWorker %s] Shutting down...", w.workerID)
			return
		case <-ticker.C:
			if err := w.processBatch(w.ctx); err != nil {
				log.Printf("[ConcurrentOutboxWorker %s] Error: %v", w.workerID, err)
			}
		}
	}
}

// processBatch is overridden to use FOR UPDATE SKIP LOCKED
func (w *ConcurrentOutboxWorker) processBatch(ctx context.Context) error {
	// Begin transaction to hold the lock
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// THE MAGIC HAPPENS HERE: FOR UPDATE SKIP LOCKED
	// 1. FOR UPDATE: Lock the selected rows so other workers can't read/write them.
	// 2. SKIP LOCKED: If another worker already locked some rows, SKIP them instead of waiting!
	// This allows multiple workers to poll the same table concurrently without contention.
	query := `
		SELECT id, aggregate_type, aggregate_id, event_type, payload, created_at, published_at, attempts, next_attempt_at
		FROM outbox_events
		WHERE published_at IS NULL AND next_attempt_at <= CURRENT_TIMESTAMP
		ORDER BY created_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`

	rows, err := tx.QueryContext(ctx, query, w.batchSize)
	if err != nil {
		return fmt.Errorf("failed to query with SKIP LOCKED: %w", err)
	}
	defer rows.Close()

	var events []OutboxEvent
	for rows.Next() {
		var e OutboxEvent
		var pubAt sql.NullTime
		if err := rows.Scan(
			&e.ID, &e.AggregateType, &e.AggregateID, &e.EventType,
			&e.Payload, &e.CreatedAt, &pubAt, &e.Attempts, &e.NextAttemptAt,
		); err != nil {
			return fmt.Errorf("scan error: %w", err)
		}
		if pubAt.Valid {
			e.PublishedAt = &pubAt.Time
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows iteration error: %w", err)
	}

	if len(events) == 0 {
		return nil // Nothing to do
	}

	// log.Printf("[Worker %s] Locked %d events for processing", w.workerID, len(events))

	now := time.Now()
	for _, event := range events {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Publish event
		pubErr := w.publisher.Publish(ctx, event)

		if pubErr != nil {
			// Backoff logic
			newAttempts := event.Attempts + 1
			backoffDuration := w.backoffBase * time.Duration(1<<(newAttempts-1))
			nextAttempt := now.Add(backoffDuration)

			_, updateErr := tx.ExecContext(ctx,
				`UPDATE outbox_events SET attempts = $1, next_attempt_at = $2 WHERE id = $3`,
				newAttempts, nextAttempt, event.ID,
			)
			if updateErr != nil {
				log.Printf("[Worker %s] Failed to update retry status for %s: %v", w.workerID, event.ID, updateErr)
			}
		} else {
			// Success
			_, updateErr := tx.ExecContext(ctx,
				`UPDATE outbox_events SET published_at = $1, attempts = attempts + 1 WHERE id = $2`,
				now, event.ID,
			)
			if updateErr != nil {
				log.Printf("[Worker %s] Failed to mark event %s published: %v", w.workerID, event.ID, updateErr)
			}
		}
	}

	// Commit transaction to release locks and save updates
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch: %w", err)
	}

	return nil
}
