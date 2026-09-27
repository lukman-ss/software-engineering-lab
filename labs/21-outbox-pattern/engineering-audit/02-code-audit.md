# Code Audit

## Finding 1

Location: internal/outbox/db.go:47-57
Claimed Behavior: Polling relay queries pending outbox records safely.
Observed Implementation: `GetPendingOutbox()` takes an RLock and returns a slice of pending messages. However, there is no row locking or status transition during query (no `FOR UPDATE` equivalent or claim phase), making concurrent relay workers prone to polling the exact same pending messages simultaneously.
Assessment: WARNING
Severity: MEDIUM
Notes: Single relay worker functions correctly, but if multiple relay instances run concurrently, redundant publish attempts occur (relying solely on consumer idempotency).

## Finding 2

Location: internal/outbox/relay.go:43-59
Claimed Behavior: Outbox relay handles publishing and marking outbox processed with proper error propagation/retries.
Observed Implementation: `PollAndDispatch()` logs errors when publishing fails or marking processed fails, but does not implement explicit retry exponential backoff or max retry limits.
Assessment: PASS
Severity: LOW
Notes: Sufficient for in-memory lab demonstration, but lacks production retry policies.

## Finding 3

Location: internal/outbox/db.go:120-128
Claimed Behavior: In-memory DB simulates transactional commit across orders and outbox maps atomically.
Observed Implementation: `Commit()` locks DB mutex and writes staged map entries to DB maps in memory.
Assessment: PASS
Severity: LOW
Notes: In-memory simulation correctly models atomic staged commits for unit test purposes, though original design doc mentioned modernc.org/sqlite.

## Finding 4

Location: internal/outbox/consumer.go:19-31
Claimed Behavior: Idempotent consumer tracks processed IDs thread-safely.
Observed Implementation: `Handle()` uses `sync.Mutex` lock to check `processedIDs[msg.ID]`, returning false on duplicate and recording new messages.
Assessment: PASS
Severity: LOW
Notes: Clean thread-safe implementation of consumer idempotency.
