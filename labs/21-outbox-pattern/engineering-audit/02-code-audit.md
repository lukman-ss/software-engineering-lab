# Code Audit

## Finding 1

Location: `internal/outbox/db.go:71-127`
Claimed Behavior: Atomic transaction commit and rollback across entity and outbox tables.
Observed Implementation: `SaveOrder` and `SaveOutbox` buffer mutations inside `stagedOrders` and `stagedOutbox`. `Commit` applies both maps atomically under `tx.db.mu.Lock()`. `Rollback` discards staged mutations.
Assessment: PASS
Severity: LOW
Notes: Properly enforces all-or-nothing semantics without partial write risk.

## Finding 2

Location: `internal/outbox/db.go:12-69`
Claimed Behavior: Thread-safe in-memory database operations.
Observed Implementation: All reads and writes to `orders` and `outbox` maps are protected by `sync.RWMutex`.
Assessment: PASS
Severity: LOW
Notes: Passed `-race` validation with concurrent workers.

## Finding 3

Location: `internal/outbox/service.go:16-83`
Claimed Behavior: Correctly contrasts atomic outbox persistence against vulnerable dual-write.
Observed Implementation: `CreateOrderWithOutbox` stages both entities within one `Tx`. `CreateOrderDualWriteNaive` commits DB transaction first, then attempts external broker publish, isolating the vulnerability.
Assessment: PASS
Severity: LOW
Notes: Clean, decoupled comparison.

## Finding 4

Location: `internal/outbox/relay.go:47-66`
Claimed Behavior: Decoupled polling relay worker dispatching messages to broker and marking them processed.
Observed Implementation: `PollAndDispatch` queries `GetPendingOutbox()`, publishes to broker, and sets status to `PROCESSED` only upon success. Failed dispatches leave messages pending for future cycles.
Assessment: PASS
Severity: LOW
Notes: Conforms to polling publisher pattern.

## Finding 5

Location: `internal/outbox/consumer.go:19-35`
Claimed Behavior: Idempotent message consumption deduplicating incoming events by ID.
Observed Implementation: `Handle` checks `processedIDs` map guarded by `sync.Mutex`. Duplicates return `false` without appending to `received`.
Assessment: PASS
Severity: LOW
Notes: Correctly models consumer-side deduplication.

## Finding 6

Location: `internal/outbox/relay.go:47-66`
Claimed Behavior: Relay error handling and retries.
Observed Implementation: Failed broker publish logs an error and retries on next poll interval with no exponential backoff or dead-letter queue.
Assessment: WARNING
Severity: LOW
Notes: Accurately scoped and noted under limitations in `engineering/02-implementation-notes.md`.
