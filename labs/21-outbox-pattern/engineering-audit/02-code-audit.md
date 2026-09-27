# Code Audit

Target Lab: labs/21-outbox-pattern

## Finding 1

Location: internal/outbox/db.go:84-140
Claimed Behavior: Atomic staging, commit, and rollback across multiple tables (`orders`, `outbox`) within a transaction.
Observed Implementation: `Tx` maintains `stagedOrders` and `stagedOutbox` maps under mutex protection. On `Commit()`, staged entries are flushed into DB under `db.mu.Lock()`. On `Rollback()`, staged entries are dropped. Subsequent mutations after close return `ErrTxClosed`.
Assessment: PASS
Severity: LOW
Notes: Correctly models atomic single-transaction boundaries in memory.

## Finding 2

Location: internal/outbox/service.go:18-53
Claimed Behavior: Order entity and Outbox message are generated and written atomically within a single transaction.
Observed Implementation: `CreateOrderWithOutbox` begins `tx`, validates payload serialization, writes both `SaveOrder` and `SaveOutbox`, rolling back on any failure before committing.
Assessment: PASS
Severity: LOW
Notes: Atomicity is guaranteed at transaction boundary.

## Finding 3

Location: internal/outbox/service.go:57-90
Claimed Behavior: Demonstrates dual-write inconsistency when publishing directly outside the transaction.
Observed Implementation: `CreateOrderDualWriteNaive` commits the DB transaction first, then attempts broker publication. If broker fails, DB contains uncommunicated state, demonstrating the dual-write flaw.
Assessment: PASS
Severity: LOW
Notes: Faithfully models the dual-write anti-pattern and resulting state divergence.

## Finding 4

Location: internal/outbox/relay.go:27-67
Claimed Behavior: Asynchronous polling relay periodically polls pending outbox entries, publishes to broker, and transitions status to processed. Supports safe start/stop and retries on broker error.
Observed Implementation: `Relay.Start()` and `Stop()` use `sync.Once` and channel closure. `PollAndDispatch()` queries `GetPendingOutbox()`, publishes to broker, and only marks processed upon successful broker acknowledgment.
Assessment: PASS
Severity: LOW
Notes: At-least-once delivery semantics preserved; failed broker dispatches remain pending for subsequent poll cycles.

## Finding 5

Location: internal/outbox/consumer.go:19-43
Claimed Behavior: Downstream consumer idempotency using message ID deduplication registry.
Observed Implementation: `Consumer` guards `processedIDs` map with `sync.Mutex`. `Handle()` checks whether `msg.ID` has already been recorded; if seen, returns `false` and skips; if new, records ID and stores message.
Assessment: PASS
Severity: LOW
Notes: Thread-safe idempotent consumer implementation correctly handles duplicate deliveries.

## Finding 6

Location: internal/outbox/broker.go:10-53
Claimed Behavior: Thread-safe mock broker with controllable failure simulation (`SetFailNext`).
Observed Implementation: `MockBroker` utilizes `sync.RWMutex`, returns copy of slice in `GetPublished()` to prevent slice race conditions, and consumes `failNext` flag on publish.
Assessment: PASS
Severity: LOW
Notes: Thread-safe and properly isolated for concurrent test execution.

## Finding 7

Location: internal/outbox/db.go:71-82
Claimed Behavior: Outbox cleanup routine purges processed messages.
Observed Implementation: `PurgeProcessedOutbox()` safely acquires write lock on DB and deletes entries marked `MessageStatusProcessed`, returning count of removed items.
Assessment: PASS
Severity: LOW
Notes: Simple and effective retention cleanup mechanism.
