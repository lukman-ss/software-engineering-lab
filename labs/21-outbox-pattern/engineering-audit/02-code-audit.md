# Code Audit

## Finding 1

Location: `internal/outbox/db.go:84-140`
Claimed Behavior: Atomic transaction management with isolated staging, commit, and rollback.
Observed Implementation: `Tx` struct buffers modifications in `stagedOrders` and `stagedOutbox` protected by `sync.Mutex`. Calling `Commit()` locks the underlying `DB.mu` and applies staged mutations. Calling `Rollback()` discards staged maps and marks transaction closed.
Assessment: PASS
Severity: LOW
Notes: Clean in-memory transaction abstraction with explicit closed state guard (`ErrTxClosed`).

---

## Finding 2

Location: `internal/outbox/service.go:18-53`
Claimed Behavior: Atomically persists order entity and outbox event in single transaction.
Observed Implementation: `CreateOrderWithOutbox` creates `Tx`, serializes `Order` into JSON payload, stages both order and outbox record with status `PENDING`, and calls `tx.Commit()`. Rollback is invoked on any staging or serialization error.
Assessment: PASS
Severity: LOW
Notes: Payload format correctly matches `Order` struct.

---

## Finding 3

Location: `internal/outbox/service.go:57-89`
Claimed Behavior: Demonstrates dual-write inconsistency when broker write fails after DB commit.
Observed Implementation: Naive write commits order to DB first, then attempts `broker.Publish()`. If broker returns error, error is returned but DB order persists without corresponding event emitted.
Assessment: PASS
Severity: LOW
Notes: Accurately models dual-write failure mode.

---

## Finding 4

Location: `internal/outbox/relay.go:43-60`
Claimed Behavior: Polling relay retrieves pending outbox records, dispatches to broker, and transitions status to `PROCESSED`.
Observed Implementation: `PollAndDispatch` calls `db.GetPendingOutbox()`, publishes each message via `broker.Publish()`, and calls `db.MarkOutboxProcessed(msg.ID)` on successful publication. Errors during publish prevent status transition, leaving records pending for retry.
Assessment: PASS
Severity: LOW
Notes: At-least-once delivery loop functions correctly.

---

## Finding 5

Location: `internal/outbox/consumer.go:19-31`
Claimed Behavior: Idempotent message consumption by deduplicating on message ID.
Observed Implementation: `Consumer.Handle` locks internal mutex, checks `processedIDs[msg.ID]`, skips duplicate delivery returning `false`, or records message and returns `true`.
Assessment: PASS
Severity: LOW
Notes: Idempotent consumer design verified and thread-safe.

---

## Finding 6

Location: `internal/outbox/db.go:71-82`
Claimed Behavior: Purging processed outbox entries to maintain table size without dropping pending records.
Observed Implementation: `PurgeProcessedOutbox` safely locks `db.mu` and deletes entries with `MessageStatusProcessed`.
Assessment: PASS
Severity: LOW
Notes: Correct deletion logic.

---

## Finding 7

Location: `internal/outbox/relay.go:24-37`
Claimed Behavior: Clean start/stop lifecycle management for asynchronous polling goroutine.
Observed Implementation: Uses `time.NewTicker` with `stopChan` select and proper `ticker.Stop()` on exit.
Assessment: PASS
Severity: LOW
Notes: Goroutine leak prevented on `Stop()`.
