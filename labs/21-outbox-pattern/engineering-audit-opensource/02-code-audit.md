# Code Audit

## Finding 1

Location: internal/outbox/db.go:112-129 (Tx.Commit)
Claimed Behavior: Commit atomically writes staged orders and outbox records to the DB under a single lock.
Observed Implementation: Commit sets `closed=true` on the tx, then acquires `tx.db.mu.Lock()` (exclusive) and writes all staged orders and outbox entries. Lock ordering is tx.mu -> db.mu, consistent with no nested-acquisition violations elsewhere.
Assessment: PASS
Severity: LOW
Notes: Commit never returns an error (always nil). This is acceptable for an in-memory mock but would mask failures in a real DB. Deliberate simplification.

## Finding 2

Location: internal/outbox/db.go:131-139 (Tx.Rollback)
Claimed Behavior: Rollback discards all staged mutations so nothing is persisted.
Observed Implementation: Rollback sets `closed=true`, which causes any subsequent SaveOrder/SaveOutbox/Commit/Rollback to return ErrTxClosed. The staged data is local to the Tx and will be garbage-collected; it is never written to db.
Assessment: PASS
Severity: LOW
Notes: None.

## Finding 3

Location: internal/outbox/service.go:18-53 (CreateOrderWithOutbox)
Claimed Behavior: Writes order and outbox message atomically in the same transaction.
Observed Implementation: Calls BeginTx, stages order and outbox, calls Commit. If any step fails, rolls back and returns error. All three claims (atomicity, rollback on error) are satisfied.
Assessment: PASS
Severity: LOW
Notes: Uses a single Tx for both domain entity and event log. This is correct outbox pattern implementation.

## Finding 4

Location: internal/outbox/service.go:57-89 (CreateOrderDualWriteNaive)
Claimed Behavior: Demonstrates dual-write flaw - DB commit succeeds but broker publish fails, leaving inconsistent state.
Observed Implementation: Commits order to DB first, then attempts broker.Publish outside the transaction. On publish failure, returns error but order remains committed. This correctly demonstrates the non-atomicity problem.
Assessment: PASS
Severity: LOW
Notes: The naive payload JSON is hand-constructed (`fmt.Sprintf`) rather than `json.Marshal`, but this is acceptable for a demonstration of the dual-write flaw.

## Finding 5

Location: internal/outbox/relay.go:43-59 (PollAndDispatch)
Claimed Behavior: Polls pending outbox messages and dispatches them to the broker. On broker failure, message stays pending for retry. On mark-processed failure, logs the error.
Observed Implementation: Fetches pending messages, iterates, calls broker.Publish. On success, marks processed. On failure, logs and leaves message pending (retry on next poll). This is correct at-least-once delivery with retry.
Assessment: PASS
Severity: LOW
Notes: If Publish succeeds but MarkOutboxProcessed fails, the message is already in the broker (duplicate risk). Consumer idempotency handles this correctly per Finding 7.

## Finding 6

Location: internal/outbox/relay.go:39-41 (Relay.Stop)
Claimed Behavior: Stops the polling goroutine.
Observed Implementation: Closes stopChan. If Stop() is called more than once, `close(stopChan)` will panic ("close of closed channel").
Assessment: WARNING
Severity: LOW
Notes: This is a latent defect. The tests only call Stop() once via `defer relay.Stop()`, so it is not triggered. Not a blocking issue but should be guarded with sync.Once or a check.

## Finding 7

Location: internal/outbox/consumer.go:19-31 (Consumer.Handle)
Claimed Behavior: Processes messages idempotently. Returns true for new, false for duplicate.
Observed Implementation: Checks processedIDs map under mutex. If ID already seen, returns false (duplicate). Otherwise records ID, appends to received, returns true.
Assessment: PASS
Severity: LOW
Notes: Thread-safe via sync.Mutex. Correct idempotency enforcement.

## Finding 8

Location: internal/outbox/broker.go:22-37 (MockBroker.SetFailNext/Publish)
Claimed Behavior: Thread-safe mock broker with simulated publish failures.
Observed Implementation: failNext flag protected by mutex. Publish sets failNext=false and returns BrokerError on failure, otherwise appends to published list.
Assessment: PASS
Severity: LOW
Notes: Single-flag failure simulation is sufficient for demo/test purposes.

## Finding 9

Location: internal/outbox/db.go:47-57 (GetPendingOutbox)
Claimed Behavior: Returns all pending outbox messages.
Observed Implementation: Reads under RLock, filters by status. Returns a copy of values (OutboxMessage is a value type).
Assessment: PASS
Severity: LOW
Notes: Values are copied from the map, so the caller gets independent copies. No data leak.

## Finding 10

Location: internal/outbox/db.go:59-68 (MarkOutboxProcessed)
Claimed Behavior: Marks a message as processed.
Observed Implementation: Acquires write lock, checks existence, sets status, writes back. The `msg` struct is a value copy from the map (since OutboxMessage is a value type in the map), so the updated value must be written back to the map. This is done via `db.outbox[id] = msg`.
Assessment: PASS
Severity: LOW
Notes: Correct write-back of value type. Good.