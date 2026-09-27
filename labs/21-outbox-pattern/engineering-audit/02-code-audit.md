# Engineering Code Audit

Target Lab: labs/21-outbox-pattern

## Finding 1

Location: `internal/outbox/db.go:25-31`, `internal/outbox/db.go:112-140`
Claimed Behavior: Atomic persistence across orders and outbox messages within single transaction boundary (`BeginTx`, `Commit`, `Rollback`).
Observed Implementation: `Tx` maintains local staging maps (`stagedOrders`, `stagedOutbox`) protected by a mutex. On `Commit()`, holding the DB write lock, all staged mutations are flushed into DB state atomically. On `Rollback()`, staged mutations are discarded without touching DB state.
Assessment: PASS
Severity: LOW
Notes: Clean in-memory transaction simulation matching isolated state transition semantics.

## Finding 2

Location: `internal/outbox/service.go:18-53`
Claimed Behavior: Atomic persistence of order and outbox records in `CreateOrderWithOutbox`.
Observed Implementation: Creates an order, marshals payload, stages both order and `OutboxMessage` within transaction, and commits atomically. Any failure rolls back the transaction.
Assessment: PASS
Severity: LOW
Notes: Properly sets `Status: MessageStatusPending` and formats outbox ID consistently (`evt-%s`).

## Finding 3

Location: `internal/outbox/service.go:57-90`
Claimed Behavior: Naive dual write commits DB first and then publishes to broker, exposing inconsistency if broker fails.
Observed Implementation: Commits order to DB first, then invokes `broker.Publish(msg)`. If broker fails, DB contains order while broker received nothing.
Assessment: PASS
Severity: LOW
Notes: Accurately replicates dual-write hazard scenario.

## Finding 4

Location: `internal/outbox/relay.go:27-48`, `internal/outbox/relay.go:50-67`
Claimed Behavior: Asynchronous polling relay reads pending outbox messages and dispatches them to broker, marking them processed.
Observed Implementation: Polling loop runs in separate goroutine using `time.Ticker` and clean shutdown via `stopChan` and `sync.Once`. `PollAndDispatch()` queries `GetPendingOutbox()`, publishes each message, and marks processed in DB upon success.
Assessment: PASS
Severity: LOW
Notes: Graceful start/stop protected against double calls with `sync.Once`.

## Finding 5

Location: `internal/outbox/consumer.go:18-31`
Claimed Behavior: Idempotent message consumption by deduplicating event IDs.
Observed Implementation: `Consumer.Handle()` checks `processedIDs` map under mutex. If already present, skips execution and returns `false`. If new, records ID and returns `true`.
Assessment: PASS
Severity: LOW
Notes: Concurrency-safe deduplication implementation.

## Finding 6

Location: `internal/outbox/db.go:71-82`
Claimed Behavior: Ability to purge processed outbox records to prevent unbounded table growth.
Observed Implementation: `PurgeProcessedOutbox()` safely acquires DB lock and deletes entries where status is `MessageStatusProcessed`, returning the purged count.
Assessment: PASS
Severity: LOW
Notes: Good implementation of outbox table retention/cleanup requirement.
