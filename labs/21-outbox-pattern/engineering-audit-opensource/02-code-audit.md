# Code Audit

Target Lab: labs/21-outbox-pattern

---

## Finding 1

Location: internal/outbox/db.go:Commit() — lines 112-129
Claimed Behavior: Order save and outbox message save are written atomically within a single transaction. Either both persist or neither persists.
Observed Implementation: `Tx` stages mutations into `stagedOrders` and `stagedOutbox` maps. On `Commit()`, it sets `tx.closed = true` then acquires `tx.db.mu.Lock()` and writes both staged maps into `db.orders` and `db.outbox` under the same lock. Rollback discards staged data without touching `db`.
Assessment: PASS
Severity: LOW
Notes: Lock ordering is tx.mu -> db.mu. No reverse ordering exists (Tx methods only lock tx.mu; DB methods only lock db.mu independently). No deadlock risk. The atomicity guarantee is sound for an in-memory store — both maps are updated under the same db.mu lock before the lock is released.

## Finding 2

Location: internal/outbox/db.go:Rollback() — lines 131-139
Claimed Behavior: Rollback discards staged mutations, no writes reach the persistent store.
Observed Implementation: Sets `tx.closed = true` and returns. Staged maps are never written to db.
Assessment: PASS
Severity: LOW
Notes: Correct implementation — Rollback simply abandons staged data. Test verifies this.

## Finding 3

Location: internal/outbox/service.go:CreateOrderWithOutbox() — lines 18-53
Claimed Behavior: Business data and outbox event written in same transaction atomically; JSON marshal failure → rollback.
Observed Implementation: Uses single BeginTx. Marshal failure rolls back. SaveOrder/SaveOutbox errors roll back. Otherwise commits. Order and event both persisted or neither.
Assessment: PASS
Severity: LOW
Notes: Correctly demonstrates atomicity. Message ID uses pattern `evt-<orderID>`.

## Finding 4

Location: internal/outbox/relay.go:PollAndDispatch() — lines 50-67
Claimed Behavior: Polling relay reads PENDING messages, publishes to broker, marks as PROCESSED. On publish failure, message remains PENDING and is retried next cycle.
Observed Implementation: GetPendingOutbox() returns snapshot. For each msg, attempts broker.Publish. On success, MarkOutboxProcessed. On failure, logs and leaves message PENDING. Next poll retries.
Assessment: PASS
Severity: MEDIUM
Notes: Correct at-least-once semantics. No exponential backoff — retries every pollInterval. No retry limit or dead-letter queue. Message remains PENDING on publish failure and is retried. Idempotency relies on consumer being idempotent, not on preventing duplicate publishes.

## Finding 5

Location: internal/outbox/relay.go:Start() — lines 27-42; Stop() — lines 44-48
Claimed Behavior: Relay starts/stops idempotently. Single goroutine with ticker.
Observed Implementation: Start() uses sync.Once. Stop() uses sync.Once + close(stopChan). Goroutine listens on ticker.C and stopChan.
Assessment: PASS
Severity: LOW
Notes: Start/Stop idempotency verified by test (TestTransactionalOutbox_RelayRetryAfterBrokerFailure calls Start() and Stop() twice). Ticker-driven polling.

## Finding 6

Location: internal/outbox/consumer.go:Handle() — lines 19-31
Claimed Behavior: Consumer deduplicates messages by ID. Returns true if processed, false if duplicate.
Observed Implementation: Checks processedIDs map. If already present, returns false. Otherwise marks processed, appends to received, returns true.
Assessment: PASS
Severity: LOW
Notes: Thread-safe with sync.Mutex. Correctly implements idempotency for at-least-once delivery.

## Finding 7

Location: internal/outbox/service.go:CreateOrderDualWriteNaive() — lines 57-90
Claimed Behavior: Demonstrates dual-write flaw — DB commit succeeds, broker publish fails → order in DB but no event sent.
Observed Implementation: Commits DB tx, then attempts broker.Publish. On failure, returns error but DB change persists. Order exists, event never sent.
Assessment: PASS
Severity: LOW
Notes: Correctly demonstrates the dual-write problem. The inconsistency (order persisted, event lost) is the core motivation for the outbox pattern.

## Finding 8

Location: internal/outbox/db.go:GetPendingOutbox() — lines 47-57
Claimed Behavior: Polling publisher implementation, not CDC/log-tailing.
Observed Implementation: Scans all outbox records on each poll, filters PENDING. Returns slice copy.
Assessment: PASS (implementation matches documented trade-off)
Severity: LOW
Notes: Design doc (engineering/01-design.md) mentions SQLite + CDC, but implementation-notes and master-draft acknowledge in-memory mock with polling. This is a deliberate simplification (see DOC_CODE_MISMATCH in 04-docs-vs-code.md).

## Finding 9

Location: internal/outbox/model.go — all types
Claimed Behavior: Outbox message has unique event ID. Consumer tracks by ID.
Observed Implementation: OutboxMessage.ID is caller-set (`evt-<orderID>`). No UUID enforcement. Consumer uses ID as dedup key.
Assessment: WARNING
Severity: MEDIUM
Notes: Research claims "Event payload marshaled into standard JSON format with unique UUID event identifiers" but implementation uses deterministic string IDs `evt-<orderID>`. In concurrent test, all workers use same orderID causing key collision — works by overwrite, not by uniqueness guarantee. Acceptable for demo but not production-grade.

## Finding 10

Location: internal/outbox/relay.go:PollAndDispatch() — dispatch loop
Claimed Behavior: At-least-once delivery, idempotent consumer handles duplicates.
Observed Implementation: If Publish succeeds but MarkOutboxProcessed fails (e.g., msg not found), message stays PENDING and will be re-published. If Publish succeeds then MarkOutboxProcessed succeeds, next poll won't re-fetch.
Assessment: PASS
Severity: LOW
Notes: Correct behavior. The window between Publish success and MarkOutboxProcessed is where duplicates can occur — consumer idempotency covers this gap.

## Finding 11

Location: internal/outbox/db.go:MarkOutboxProcessed() — lines 59-69
Claimed Behavior: Idempotent marking of outbox message as PROCESSED.
Observed Implementation: Checks existence, sets Status to PROCESSED, writes back to map.
Assessment: PASS
Severity: LOW
Notes: Returns error if message not found. Called by relay after successful publish.

## Finding 12

Location: internal/outbox/db.go:PurgeProcessedOutbox() — lines 71-82
Claimed Behavior: Cleanup of processed outbox messages.
Observed Implementation: Deletes all outbox records with PROCESSED status. Returns count purged.
Assessment: PASS
Severity: LOW
Notes: Function exists and is tested. Not called in demo or by relay (manual operation as documented).

## Finding 13

Location: internal/outbox/broker.go:MockBroker.Publish() — lines 28-37
Claimed Behavior: Simulates publish failures via failNext flag.
Observed Implementation: If failNext is true, returns BrokerError, clears flag. Otherwise appends to published.
Assessment: PASS
Severity: LOW
Notes: Simple but sufficient for fault-injection testing. Thread-safe.

## Finding 14

Location: internal/outbox/relay.go (PollAndDispatch + GetPendingOutbox)
Claimed Behavior: Relay and DB access are thread-safe.
Observed Implementation: GetPendingOutbox takes RLock, returns copy. MarkOutboxProcessed takes Lock. Relay goroutine runs concurrently with writers.
Assessment: PASS
Severity: LOW
Notes: Race detector passes. Lock discipline is consistent.

## Finding 15

Location: internal/outbox/consumer.go
Claimed Behavior: Consumer is thread-safe.
Observed Implementation: All methods use sync.Mutex.
Assessment: PASS
Severity: LOW
Notes: Not tested under concurrency in current test suite, but lock discipline is sound.

## Finding 16

Location: internal/outbox/db.go:Commit() — lines 112-129
Claimed Behavior: Error propagation from staged save to commit.
Observed Implementation: None — Commit returns nil unconditionally on success path. The only error is ErrTxClosed if called twice.
Assessment: WARNING
Severity: MEDIUM
Notes: In a real DB, Commit could fail (disk full, constraint violation). The in-memory implementation cannot fail to write to maps. This is an inherent limitation of the in-memory simulation and does not represent a defect in the pattern implementation.

## Finding 17

Location: internal/outbox/relay.go:Start() goroutine
Claimed Behavior: No goroutine leak on Stop.
Observed Implementation: Stop closes stopChan; goroutine's select picks it up and returns. Ticker deferred Stop() called.
Assessment: PASS
Severity: LOW
Notes: Clean shutdown. No goroutine leak.

## Finding 18

Location: cmd/demo/main.go:Scenario 3 — lines 56-59
Claimed Behavior: Duplicate delivery simulation.
Observed Implementation: Reuses published[0] (same message ID) sent to consumer.Handle again.
Assessment: PASS
Severity: LOW
Notes: Correctly demonstrates idempotency — same ID rejected on second delivery.