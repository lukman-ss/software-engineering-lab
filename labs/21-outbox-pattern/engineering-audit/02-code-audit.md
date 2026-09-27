# Code Audit

Target Lab: labs/21-outbox-pattern

## Finding 1

Location: `internal/outbox/db.go:25-31`, `84-140`
Claimed Behavior: Transactional atomicity across orders and outbox messages (staging mutations until Commit, discarding on Rollback).
Observed Implementation: `Tx` struct buffers modifications in `stagedOrders` and `stagedOutbox` maps under transaction mutex. On `Commit()`, lock on DB is acquired and maps are applied atomically to DB state. On `Rollback()`, staged maps are discarded without touching DB.
Assessment: PASS
Severity: LOW
Notes: Clean thread-safe in-memory transaction boundary implementation.

## Finding 2

Location: `internal/outbox/service.go:18-53`
Claimed Behavior: Atomic insertion of Order and OutboxMessage within the same database transaction.
Observed Implementation: Creates `Order` and `OutboxMessage`, stages both in `tx`, and commits. If JSON serialization or staging fails, rollback is triggered.
Assessment: PASS
Severity: LOW
Notes: Correctly demonstrates transactional outbox write pattern.

## Finding 3

Location: `internal/outbox/relay.go:27-66`
Claimed Behavior: Asynchronous polling relay reads pending outbox messages, publishes to broker, and marks status as PROCESSED. Handles start/stop idempotency.
Observed Implementation: `Start()` and `Stop()` protected by `sync.Once`. Ticker runs background goroutine. `PollAndDispatch()` queries `GetPendingOutbox()`, publishes each, and calls `MarkOutboxProcessed()`. If publish fails, message remains pending for subsequent retry.
Assessment: PASS
Severity: LOW
Notes: Concurrency safe and properly handles transient dispatch failures.

## Finding 4

Location: `internal/outbox/consumer.go:19-31`
Claimed Behavior: Idempotent message consumption using processed event ID tracking.
Observed Implementation: `Consumer.Handle` checks `processedIDs[msg.ID]` under mutex lock. If already processed, returns false (skipped). If new, sets `processedIDs[msg.ID] = true` and appends message.
Assessment: PASS
Severity: LOW
Notes: Accurately implements at-least-once deduplication logic.

## Finding 5

Location: `internal/outbox/service.go:57-90`
Claimed Behavior: Naive dual write demonstrates vulnerability when broker fails after DB commit.
Observed Implementation: `CreateOrderDualWriteNaive` commits `Order` to DB first, then attempts `broker.Publish()`. If broker fails, DB has order record but broker has no message.
Assessment: PASS
Severity: LOW
Notes: Effectively proves dual-write inconsistency flaw.
