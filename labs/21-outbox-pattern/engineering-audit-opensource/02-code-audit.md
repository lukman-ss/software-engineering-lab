## Finding 1
Location: internal/outbox/db.go:99-116, internal/outbox/service.go:18-53
Claimed Behavior: Order and outbox event are written atomically within a single local transaction; if commit succeeds both persist, if rollback neither persists.
Observed Implementation: `CreateOrderWithOutbox` creates a `Tx`, stages both `Order` and `OutboxMessage` into separate staged maps, then calls `Commit()`. `Commit()` acquires `db.mu.Lock()` and applies both staged maps to the persistent `orders`/`outbox` maps under a single lock hold, then marks `closed = true`. `Rollback()` discards staged mutations. The atomicity guarantee is correctly simulated.
Assessment: PASS
Severity: N/A
Notes: Since `BeginTx` and `Commit` are synchronous and all DB mutation occurs under `db.mu`, there is no window where the order is visible without the outbox message (or vice versa).

## Finding 2
Location: internal/outbox/relay.go:43-60
Claimed Behavior: Decoupled polling relay queries pending outbox records, dispatches them to the broker, and marks them as processed upon successful delivery. At-least-once delivery is implied.
Observed Implementation: `Relay.Start()` spawns a goroutine with a ticker. Each tick calls `PollAndDispatch()`, which calls `GetPendingOutbox()`, iterates pending messages, calls `broker.Publish(msg)`, and on success calls `db.MarkOutboxProcessed(msg.ID)`. If Publish fails, the message stays in PENDING status and will be re-picked on the next tick (at-least-once). If Publish succeeds but MarkOutboxProcessed fails or the process crashes before marking, the message remains PENDING and will be re-published; the idempotent consumer deduplicates.
Assessment: PASS
Severity: N/A
Notes: This correctly models at-least-once delivery. No transactional "in-flight" status exists between publish and mark, which is acceptable for the single-worker polling pattern but is a known limitation documented in engineering notes.

## Finding 3
Location: internal/outbox/consumer.go:18-31
Claimed Behavior: Consumer tracks event IDs to handle duplicate deliveries gracefully (idempotency).
Observed Implementation: `Consumer.Handle(msg)` acquires a mutex, checks `processedIDs[msg.ID]`, and returns `false` if already processed, otherwise records the ID and appends the message, returning `true`. `GetReceivedCount()` and `IsProcessed()` are also mutex-protected.
Assessment: PASS
Severity: N/A

## Finding 4
Location: internal/outbox/service.go:55-90, tests/outbox_test.go:117-139, cmd/demo/main.go:20-27
Claimed Behavior: Naive dual-write (DB commit then direct broker publish) leads to inconsistency when the broker is unavailable.
Observed Implementation: `CreateOrderDualWriteNaive` commits the order to the DB via Tx.Commit(), then calls `broker.Publish(msg)` outside the transaction. If Publish returns an error (broker unavailable), the function returns the error but the order remains persisted in the DB with no event sent. The test `TestDualWriteProblem_Failure` and the demo `Scenario 1` both verify that `dbFound = true` and `broker message count = 0`.
Assessment: PASS
Severity: N/A

## Finding 5
Location: internal/outbox/db.go, internal/outbox/broker.go, internal/outbox/consumer.go
Claimed Behavior: Thread-safe implementations with zero race conditions.
Observed Implementation: `DB` uses `sync.RWMutex`; `Tx` uses `sync.Mutex`; `MockBroker` uses `sync.RWMutex`; `Consumer` uses `sync.Mutex`. All map/slice access is guarded. `GetPendingOutbox` and `GetPublished` return copies.
Assessment: PASS
Severity: N/A
Notes: Verified with `go test -race -count=1 ./...` — no data races detected.

## Finding 6
Location: internal/outbox/relay.go:39-41
Claimed Behavior: `Stop()` gracefully terminates the relay goroutine.
Observed Implementation: `Stop()` calls `close(r.stopChan)`. If called twice, Go panics with "close of closed channel".
Assessment: WARNING
Severity: LOW
Notes: Not currently triggered by tests or demo (single Stop per lifecycle), but is a latent bug for any caller that invokes Stop() more than once. Should be guarded with sync.Once or a boolean flag.

## Finding 7
Location: internal/outbox/relay.go:43-60, internal/outbox/db.go:47-57, 59-69
Claimed Behavior: Relay correctly handles the case where a pending message is picked up but the broker fails.
Observed Implementation: When `broker.Publish` returns an error, the relay logs and continues. The message remains PENDING and will be retried on the next poll cycle. No retry-with-backoff is implemented.
Assessment: PASS
Severity: N/A
Notes: The lack of exponential backoff is documented as a known limitation in engineering/02-implementation-notes.md:27-29.