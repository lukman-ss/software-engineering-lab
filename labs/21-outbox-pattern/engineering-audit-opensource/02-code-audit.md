# Code Audit

Target: labs/21-outbox-pattern

Files reviewed: internal/outbox/db.go, internal/outbox/model.go, internal/outbox/service.go, internal/outbox/broker.go, internal/outbox/relay.go, internal/outbox/consumer.go, cmd/demo/main.go

## Finding 1

Location: internal/outbox/db.go:71-127 (Tx.Commit/Rollback/SaveOrder/SaveOutbox)
Claimed Behavior: Order and outbox records are written in a single transaction that commits atomically or rolls back atomically.
Observed Implementation: `BeginTx` returns a `Tx` holding staged `stagedOrders` and `stagedOutbox` maps. `SaveOrder`/`SaveOutbox` stage mutations locally under `tx.mu` (returning `ErrTxClosed` if already closed). `Commit` sets `closed=true`, then under `tx.db.mu` merges both staged maps into the live `db.orders`/`db.outbox`. `Rollback` sets `closed=true` and discards staged maps.
Assessment: PASS
Severity: HIGH
Notes: Lock ordering is always `tx.mu` -> `db.mu`; no path locks `db.mu` then `tx.mu`, so no deadlock. A `closed` flag correctly prevents double commit/rollback. The `tx.mu` guard and `db.mu` guard ensure no data race between staging and committing. This correctly models transactional atomicity for the two record types.

## Finding 2

Location: internal/outbox/service.go:17-53 (CreateOrderWithOutbox), :57-89 (CreateOrderDualWriteNaive)
Claimed Behavior: Transactional outbox path saves order + outbox event in one atomic transaction; naive dual-write path publishes to broker outside the DB transaction and is inconsistent on broker failure.
Observed Implementation: `CreateOrderWithOutbox` calls `BeginTx`, stages `SaveOrder` then `SaveOutbox`, marshals payload first (rolls back on marshal failure), and `Commit()`s both atomically. `CreateOrderDualWriteNaive` commits the order first, then calls `broker.Publish` outside the transaction; on publish failure the order is already committed.
Assessment: PASS
Severity: HIGH
Notes: The atomic outbox path is correct and demonstrably distinguishes itself from the naive dual-write path. Both code paths handle error propagation correctly (rollback on internal failure, error returned on external publish failure). `CreateOrderDualWriteNaive` does not stage an outbox message, which is the intended demonstration of the flaw.

## Finding 3

Location: internal/outbox/broker.go:10-53 (MockBroker)
Claimed Behavior: Thread-safe mock broker that simulates publish failures.
Observed Implementation: `MockBroker` guards `published` and `failNext` with `sync.RWMutex`. `SetFailNext` writes under `Lock`; `Publish` acquires `Lock`, checks/resets `failNext`, appends under `Lock`; `GetPublished` copies under `RLock`.
Assessment: PASS
Severity: MEDIUM
Notes: `SetFailNext` is safe to call concurrently with `Publish`/`GetPublished`. `failNext` is a single-fail flag (reset after one failed publish), which matches the documented "simulate broker down" intent. The copy in `GetPublished` prevents external mutation of internal state.

## Finding 4

Location: internal/outbox/relay.go:8-60 (Relay.Start/Stop/PollAndDispatch)
Claimed Behavior: Asynchronous polling worker that dispatches pending outbox records to the broker, implementing at-least-once delivery.
Observed Implementation: `Start` spawns a goroutine with a `time.Ticker` and a `stopChan`, dispatching on each tick or returning on stop. `PollAndDispatch` snapshots pending messages via `db.GetPendingOutbox`, publishes each, and on success marks the message processed; on publish failure the message is left `PENDING` (logged) so it can be retried on the next poll.
Assessment: PASS
Severity: MEDIUM
Notes: This correctly implements at-leased-once delivery: if publish fails or the relay crashes before `MarkOutboxProcessed`, the message remains pending and is re-dispatched later. The duplicate delivery produced by at-least-once is handled by the consumer's idempotency (Finding 5). Logs record publish/mark failures but the method does not propagate them to the caller (acceptable for a background worker, but errors are not surfaced).

## Finding 5

Location: internal/outbox/consumer.go:5-43 (Consumer.Handle/GetReceivedCount/IsProcessed)
Claimed Behavior: Downstream consumer enforces idempotency through event ID tracking.
Observed Implementation: `Consumer` guards `processedIDs` and `received` with `sync.Mutex`. `Handle` checks `processedIDs[msg.ID]`; if present returns `false` (duplicate safely skipped), otherwise records the ID and appends the message, returning `true`. `GetReceivedCount` and `IsProcessed` read under lock.
Assessment: PASS
Severity: HIGH
Notes: The duplicate-suppression logic is correct and thread-safe. The boolean return value cleanly distinguishes first delivery from duplicate, matching the demo's at-least-once scenario. `IsProcessed` is defined but unused by the demo/tests (minor dead surface, not a defect).

## Finding 6

Location: internal/outbox/relay.go:39-41 (Relay.Stop)
Claimed Behavior: Relay supports graceful shutdown.
Observed Implementation: `Stop` calls `close(r.stopChan)`. The `stopChan` is created once (unbuffered) in `NewRelay`.
Assessment: WARNING
Severity: LOW
Notes: `close(stopChan)` is not guarded against a second call. Calling `Stop()` twice would panic with "close of closed channel". In the current demo and tests `Stop` is called exactly once (via a single `defer`), so this is latent rather than triggered. `Stop` also does not block until the goroutine exits, so callers cannot be certain the worker has terminated before reusing shared state. Not a test failure today but a real concurrency hazard if the relay lifecycle is reused.

## Finding 7

Location: internal/outbox/relay.go:43-59 (PollAndDispatch)
Claimed Behavior: No message is silently lost on broker failure.
Observed Implementation: On `broker.Publish` failure or `MarkOutboxProcessed` failure the message is left `PENDING` and logged.
Assessment: PASS
Severity: MEDIUM
Notes: Retry semantics are correct for transient broker failures. There is no poison-message/backoff handling; a persistently failing message would be retried every poll interval indefinitely. This is acceptable for a demonstration lab but is a production gap (noted, not a defect).
