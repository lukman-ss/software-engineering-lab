# Code Audit

## Finding 1

Location: internal/outbox/db.go:48-57, 59-69
Claimed Behavior: Polling relay queries pending outbox records and updates record state to PROCESSED safely.
Observed Implementation: `DB.GetPendingOutbox` grabs `RLock()` and returns copies/slice of pending messages. `DB.MarkOutboxProcessed` acquires `Lock()` and updates `Status = MessageStatusProcessed`. Mutations on `Tx` are staged in local maps and atomically written to `DB.orders` and `DB.outbox` upon `Tx.Commit()` under `DB.mu.Lock()`.
Assessment: PASS
Severity: LOW
Notes: Thread-safe in-memory database simulation satisfies single-node transactional boundary requirements.

## Finding 2

Location: internal/outbox/service.go:18-53
Claimed Behavior: Atomic insert of Order entity and Outbox record within a single database transaction.
Observed Implementation: `OrderService.CreateOrderWithOutbox` creates `Tx`, stages both `Order` and `OutboxMessage` (`Status: MessageStatusPending`), and calls `tx.Commit()`. If marshaling or staging fails, `tx.Rollback()` is invoked.
Assessment: PASS
Severity: LOW
Notes: Atomicity guarantee is correctly structured.

## Finding 3

Location: internal/outbox/service.go:57-90
Claimed Behavior: Naive dual-write demonstrates state inconsistency on broker failure.
Observed Implementation: `OrderService.CreateOrderDualWriteNaive` commits `Order` to DB first, then attempts direct broker dispatch. If `broker.Publish` fails, error returns while `Order` remains persisted in DB.
Assessment: PASS
Severity: LOW
Notes: Clearly isolates and demonstrates the dual-write flaw.

## Finding 4

Location: internal/outbox/relay.go:43-60
Claimed Behavior: Asynchronous polling relay reads pending outbox messages, publishes to broker, and marks records processed.
Observed Implementation: `Relay.PollAndDispatch` fetches pending outbox events via `db.GetPendingOutbox()`, publishes each message via `broker.Publish(msg)`, and on success updates status via `db.MarkOutboxProcessed(msg.ID)`. On failure, message remains `PENDING` for subsequent poll cycles.
Assessment: PASS
Severity: LOW
Notes: Implements polling publisher pattern.

## Finding 5

Location: internal/outbox/consumer.go:19-31
Claimed Behavior: Idempotent consumer skips duplicate events based on message ID tracking.
Observed Implementation: `Consumer.Handle` locks mutex, verifies presence in `processedIDs[msg.ID]`, skips duplicate if present, otherwise marks ID as processed and appends to `received`.
Assessment: PASS
Severity: LOW
Notes: Standard message deduplication table pattern correctly implemented.

## Finding 6

Location: engineering/01-design.md:14,25 vs internal/outbox
Claimed Behavior: Design document mentions "Outbox cleanup job removes or archives processed outbox records after retention interval."
Observed Implementation: No cleanup worker or table purging function is implemented in `internal/outbox`. `db.go` only updates status to `PROCESSED`.
Assessment: WARNING
Severity: LOW
Notes: While documented as a potential expected behavior / success criteria in design doc, `engineering/02-implementation-notes.md` scopes the implementation to core outbox & idempotency. Does not block core verification.
