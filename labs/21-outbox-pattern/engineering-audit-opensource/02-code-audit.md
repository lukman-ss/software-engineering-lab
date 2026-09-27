# Code Audit

## Finding 1
Location: internal/outbox/db.go:112-129 (Tx.Commit) and db.go:131-139 (Tx.Rollback)
Claimed Behavior: Order and outbox record persist atomically; rollback discards staged mutations.
Observed Implementation: Commit sets closed=true then under a single db.mu.Lock applies all staged orders + outbox. Rollback sets closed=true and discards staged maps. Staged maps are tx-owned, never mutated after BeginTx; no DB writes occur before commit.
Assessment: PASS
Severity: LOW
Notes: Atomicity hinges on holding db.mu for the whole commit window; acceptable for in-memory DB. Commit is non-atomic with respect to partial failure (e.g. if a map write panicked), but no such failure mode exists here.

## Finding 2
Location: internal/outbox/service.go:55-89 (CreateOrderDualWriteNaive)
Claimed Behavior: Demonstrates dual-write vulnerability — DB committed, then broker write fails leaving order persisted but event lost.
Observed Implementation: tx.Commit() runs first (order persisted), then broker.Publish fails, error returned without outbox record ever being written. Order remains in DB, zero broker messages.
Assessment: PASS
Severity: LOW
Notes: Correctly demonstrates the dual-write inconsistency anti-pattern. Matches test TestDualWriteProblem_Failure.

## Finding 3
Location: internal/outbox/relay.go:43-60 (PollAndDispatch)
Claimed Behavior: Relay polls PENDING outbox records and dispatches to broker; on success marks PROCESSED and counts dispatch; on failure logs.
Observed Implementation: GetPendingOutbox returns all PENDING msgs; for each, broker.Publish called; if nil error, MarkOutboxProcessed called and dispatched incremented; failures logged.
Assessment: PASS
Severity: MEDIUM
Notes: Publish-then-mark is NOT atomic within a single DB transaction. If process crashes between Publish success and MarkOutbox, the message is left PROCESSED-but-not-flagged (actually remains PENDING) causing a DUPLICATE publish to broker on next poll. This is the documented at-least-once semantics, acceptable in design but should be noted. No transactional guard wrapping publish+mark.

## Finding 4
Location: internal/outbox/relay.go:24-41 (Start/Stop)
Claimed Behavior: Asynchronous polling worker started/stopped safely.
Observed Implementation: Start spawns goroutine with ticker + stopChan. Stop closes stopChan.
Assessment: WARNING
Severity: MEDIUM
Notes: Stop() does not guard for multiple calls — close on already-closed stopChan panics. Start() called twice starts multiple goroutines sharing same stopChan (only one will receive close signal, others leak). No sync.Once / goroutine lifecycle tracking. Acceptable for limited demo but a latent bug.

## Finding 5
Location: internal/outbox/relay.go:44-58 / broker.go:28-37
Claimed Behavior: Published messages received by consumers in order.
Observed Implementation: Broker.published is a slice appended under lock; relay processes pending msgs in map-iteration order (GETPENDINGOUTBOX iterates db.outbox map — non-deterministic order). Consumer.Handle appends to received in that order.
Assessment: WARNING
Severity: LOW
Notes: Message processing order is non-deterministic due to map iteration. Doesn't break correctness (idempotency holds) but demo output order may vary. Documented nowhere.

## Finding 6
Location: internal/outbox/consumer.go:19-31
Claimed Behavior: Idempotent consumption — duplicate message ID rejected, first accepted.
Observed Implementation: processedIDs map under mutex; duplicate ID returns false, no append; first ID sets true and appends.
Assessment: PASS
Severity: LOW
Notes: Correct idempotency via event-ID tracking. Single-consumer only (per-instance lock). Documented semantics match.

## Finding 7
Location: internal/outbox/db.go:47-57 (GetPendingOutbox), db.go:59-69 (MarkOutboxProcessed)
Claimed Behavior: Relay only dispatches PENDING records; marks processed atomically.
Observed Implementation: GetPendingOutbox returns PENDING msgs. MarkOutboxProcessed sets status=PROCESSED under db.mu and overwrites db.outbox[id].
Assessment: PASS
Severity: LOW
Notes: Each DB mutation individually locked. No race on status transitions.

## Finding 8
Location: internal/outbox/relay.go:44-46
Claimed Behavior: At-least-once delivery until processed.
Observed Implementation: On Publish failure, message remains PENDING (not marked) and will be retried on next poll. On Publish success but MarkOutbox failure (logged), dispatched not incremented; message stays PENDING -> retried.
Assessment: PASS
Severity: LOW
Notes: Retry semantics correct; at-least-once behavior preserved.

## Finding 9
Location: internal/outbox/consumer.go / relay.go
Claimed Behavior: Relay delivers to consumer.
Observed Implementation: Relay.PollAndDispatch only publishes to Broker interface; it does NOT call Consumer.Handle. The Consumer is invoked manually by demo/tests on broker.GetPublished().
Assessment: WARNING
Severity: MEDIUM
NOTES: README claims relay "dispatching them to a message broker" (OK) but the consumer is described as a Subscriber receiving dispatched events. In actual code there is no wire from Broker->Consumer; consumers read broker.GetPublished() themselves. Architecture description is partially accurate (relay -> broker) but lacks the broker->consumer delivery leg. Not a bug, but an architectural gap vs README framing.
