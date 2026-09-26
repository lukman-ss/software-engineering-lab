# Code Audit Findings

Target Lab: labs/21-outbox-pattern

## Finding 1

Location: `internal/outbox/db.go:99-127`
Claimed Behavior: Atomic commit and rollback across domain and outbox state.
Observed Implementation: `Tx` stages mutations in `stagedOrders` and `stagedOutbox`. During `Commit`, it acquires `db.mu.Lock()` and transfers staged entries into `db.orders` and `db.outbox` within the lock. `Rollback` marks the transaction closed and discards staged entries without touching `db`.
Assessment: PASS
Severity: LOW
Notes: Correctly models transactional staging and atomicity for in-memory simulation.

## Finding 2

Location: `internal/outbox/relay.go:43-60`
Claimed Behavior: Decoupled relay polls pending events, dispatches to broker, and marks status processed.
Observed Implementation: `PollAndDispatch` fetches pending outbox messages from `db.GetPendingOutbox()`, publishes each message via `broker.Publish(msg)`, and on success calls `db.MarkOutboxProcessed(msg.ID)`. If publish fails, the message remains `PENDING` for next polling iteration.
Assessment: PASS
Severity: LOW
Notes: Properly adheres to at-least-once delivery semantics.

## Finding 3

Location: `internal/outbox/consumer.go:19-31`
Claimed Behavior: Idempotent processing of duplicate message deliveries.
Observed Implementation: `Consumer.Handle` locks mutex, verifies whether `msg.ID` exists in `processedIDs`, returns `false` if seen, and records it only once if unseen.
Assessment: PASS
Severity: LOW
Notes: Idempotency is thread-safe and verified.

## Finding 4

Location: `internal/outbox/broker.go:28-37`
Claimed Behavior: Injectable broker failures to prove dual-write state inconsistency.
Observed Implementation: `MockBroker.Publish` uses mutex-protected `failNext` flag to simulate downstream broker unavailability, returning `BrokerError`.
Assessment: PASS
Severity: LOW
Notes: Reliable mock behavior without hidden side effects.

## Finding 5

Location: `internal/outbox/service.go:57-89`
Claimed Behavior: Dual-write naive implementation demonstrates failure when broker write fails post DB commit.
Observed Implementation: Commits order to DB first, then invokes `broker.Publish(msg)`. If `failNext` is active, DB contains order while broker has zero records, replicating real-world partial failure.
Assessment: PASS
Severity: LOW
Notes: Clear pedagogical contrast to transactional outbox pattern.
