# Code Audit

## Finding 1

Location: `internal/outbox/db.go:84-140`
Claimed Behavior: In-memory transactional database simulating atomic `BeginTx`, `Commit`, and `Rollback` across `orders` and `outbox` records.
Observed Implementation: `Tx` stages mutations in `stagedOrders` and `stagedOutbox` maps under mutex protection. On `Commit()`, mutations are atomically applied to `db.orders` and `db.outbox` under `db.mu` write lock. On `Rollback()`, staged maps are discarded.
Assessment: PASS
Severity: LOW
Notes: Correctly models basic transactional staging and isolation in-memory.

## Finding 2

Location: `internal/outbox/service.go:18-53`
Claimed Behavior: Atomic creation of order and outbox record within single transaction.
Observed Implementation: Service begins transaction, stages `Order` and `OutboxMessage`, rolls back on marshal or save error, and commits transaction atomically.
Assessment: PASS
Severity: LOW
Notes: Properly handles error rollback paths and commits atomically.

## Finding 3

Location: `internal/outbox/service.go:57-90`
Claimed Behavior: Direct dual-write implementation demonstrates dual-write vulnerability when broker publish fails after DB commit.
Observed Implementation: `CreateOrderDualWriteNaive` commits DB transaction first, then calls `broker.Publish()`. When `broker.Publish()` fails, DB contains order while broker has zero messages.
Assessment: PASS
Severity: LOW
Notes: Accurately reproduces dual-write inconsistency scenario.

## Finding 4

Location: `internal/outbox/relay.go:43-59`
Claimed Behavior: Polling worker querying pending outbox records, dispatching to broker, and updating message status to `PROCESSED`.
Observed Implementation: `PollAndDispatch()` retrieves pending messages via `GetPendingOutbox()`, publishes each message, and calls `MarkOutboxProcessed(msg.ID)` on successful publish. On publish failure, error is logged and message remains `PENDING`.
Assessment: PASS
Severity: LOW
Notes: Implements polling publisher at-least-once delivery semantics cleanly.

## Finding 5

Location: `internal/outbox/consumer.go:19-31`
Claimed Behavior: Downstream consumer enforcing idempotency via event ID tracking.
Observed Implementation: `Consumer.Handle()` checks `processedIDs` under mutex lock. If ID exists, returns `false` (skipped duplicate). If new, records ID and stores message.
Assessment: PASS
Severity: LOW
Notes: Safe concurrent idempotent handling.

## Finding 6

Location: `internal/outbox/db.go:71-82`
Claimed Behavior: Outbox purge mechanism removes processed outbox records.
Observed Implementation: `PurgeProcessedOutbox()` acquires write lock and deletes all records where `Status == MessageStatusProcessed`.
Assessment: PASS
Severity: LOW
Notes: Performs correct in-memory deletion of processed entries.
