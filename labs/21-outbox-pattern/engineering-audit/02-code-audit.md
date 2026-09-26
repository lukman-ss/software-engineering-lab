# Code Audit

## Finding 1

Location: `internal/outbox/db.go:25-31`, `84-140`
Claimed Behavior: Atomic persistence across orders and outbox records with transactional staging and commit/rollback.
Observed Implementation: `Tx` stages mutations in local maps under `Tx.mu`. `Commit()` acquires `DB.mu` lock and merges staged maps into main storage atomically. `Rollback()` marks tx closed and discards staged mutations.
Assessment: PASS
Severity: LOW
Notes: Correctly implements simulated in-memory ACID transaction staging.

## Finding 2

Location: `internal/outbox/service.go:18-53`, `57-90`
Claimed Behavior: Atomic outbox creation vs dual-write failure simulation.
Observed Implementation: `CreateOrderWithOutbox` binds `SaveOrder` and `SaveOutbox` in one transaction. `CreateOrderDualWriteNaive` commits order to DB first, then attempts broker publish; if broker publish fails, DB state persists while message broker receives nothing.
Assessment: PASS
Severity: LOW
Notes: Accurately models the exact vulnerability of dual-writes and the atomic solution.

## Finding 3

Location: `internal/outbox/relay.go:24-60`
Claimed Behavior: Decoupled asynchronous polling worker querying pending outbox records and publishing to broker.
Observed Implementation: Background goroutine polls DB pending records on ticker interval, calls `broker.Publish`, and marks records processed upon success. Handles stop signals via channel.
Assessment: PASS
Severity: LOW
Notes: Polling loop and state transitions are thread-safe and properly synchronized.

## Finding 4

Location: `internal/outbox/consumer.go:19-31`
Claimed Behavior: Downstream consumer idempotency enforcement using message ID deduplication.
Observed Implementation: `Consumer` maintains `processedIDs` map guarded by `sync.Mutex`. Duplicate delivery check returns `false` and ignores payload.
Assessment: PASS
Severity: LOW
Notes: Correctly verifies consumer-side deduplication guarantee.
