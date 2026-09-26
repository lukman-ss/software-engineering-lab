# Code Audit: Correctness, Concurrency, and Failure Handling

## Finding 1: Atomicity and Transaction Simulation

Location: `labs/21-outbox-pattern/internal/outbox/db.go` (Tx.Commit, Tx.Rollback)

Claimed Behavior: The DB simulates an atomic transaction where both order and outbox writes are persisted together or neither is persisted.

Observed Implementation: 
- Tx stages order/outbox writes in local maps (`stagedOrders`, `stagedOutbox`).
- Commit acquires DB write lock and copies staged data to the DB's maps, marking Tx closed.
- Rollback sets Tx closed and discards staged maps without touching DB.
- Service.CreateOrderWithOutbox creates a Tx, saves order and outbox into Tx, then calls Commit (both succeed) or Rollback on any error.

Assessment: PASS. Atomicity is correctly simulated by staging then bulk-commit under a single write lock. Both writes are either committed together (under db.mu.Lock()) or discarded together (Rollback). No partial visibility of order without outbox or vice-versa.

Severity: -
Notes: None.

## Finding 2: Race Conditions in Concurrent Access

Location: `labs/21-outbox-pattern/internal/outbox/*.go` (all map accesses)

Claimed Behavior: Concurrent reads/writes to orders, outbox, broker, and consumer state are safe.

Observed Implementation:
- DB: orders/outbox maps guarded by sync.RWMutex (RLock for reads, Lock for writes/mutations).
- Tx: stagedOrders/stagedOutbox guarded by sync.Mutex.
- MockBroker: published/failNext guarded by sync.Mutex.
- Consumer: processedIDs/received guarded by sync.Mutex.
- Relay: Accesses DB via synchronized methods; single goroutine poll loop.

Assessment: PASS. All shared state accesses are properly synchronized with mutexes or RWMutexes. The `go test -race` passes, confirming no data races. Relay is single-threaded and only reads from DB (with RLock) and writes via synchronized MarkOutboxProcessed/Publish.

Severity: -
Notes: The concurrent test (`TestTransactionalOutbox_ConcurrentWrites`) uses the same orderID across goroutines, which stresses locking but overwrites the same key; nonetheless, race detector passes, indicating correct locking.

## Finding 3: Broker Failure Handling in Relay

Location: `labs/21-outbox-pattern/internal/outbox/relay.go` (PollAndDispatch)

Claimed Behavior: If broker publish fails, the outbox message remains PENDING for retry on next poll.

Observed Implementation:
- In `PollAndDispatch`, if `r.broker.Publish(msg)` returns non-nil error:
  - Logs the error
  - Does NOT call `r.db.MarkOutboxProcessed(msg.ID)`
  - Message remains with status `MessageStatusPending` in DB
  - Loop continues to next message (does not break or return early)
- Thus on next poll cycle, the same message will be seen as pending again and retried.

Assessment: PASS. The relay correctly leaves outbox messages pending on broker failure, enabling retry. This implements the at-least-once delivery guarantee.

Severity: -
Notes: This behavior is exercised in the code but not asserted by any test (see Test Audit).

## Finding 4: Consumer Idempotency

Location: `labs/21-outbox-pattern/internal/outbox/consumer.go` (Handle)

Claimed Behavior: Consumer ignores duplicate message deliveries based on event ID.

Observed Implementation:
- Consumer maintains a `processedIDs` map (protected by mu).
- On `Handle(msg)`: if `msg.ID` exists in map, return false (duplicate ignored); else mark as processed, store in received slice, return true.
- Test `TestTransactionalOutbox_Idempotency_DuplicateDelivery` verifies first Handle=true, second Handle=false.

Assessment: PASS. Idempotency is correctly implemented via event-ID deduplication.

Severity: -
Notes: None.

## Finding 5: Rollback Behavior

Location: `labs/21-outbox-pattern/internal/outbox/db.go` (Tx.Rollback)

Claimed Behavior: Tx.Rollback discards all staged order and outbox writes.

Observed Implementation:
- Rollback sets `tx.closed = true` and returns nil.
- Staged maps are left untouched but ignored on any further Tx operations (due to closed flag).
- New Tx begins with fresh empty staged maps.

Assessment: PASS. Rollback prevents staged data from ever being committed to DB. Service.CreateOrderWithOutbox calls Rollback on any error during save/marshal/commit sequence.

Severity: -
Notes: The marshal-error path in service (json.Marshal failure) is unreachable in practice with the given Order struct but exists as a defensive measure.

## Finding 6: PurgeProcessed Behavior

Location: `labs/21-outbox-pattern/internal/outbox/db.go` (PurgeProcessedOutbox)

Claimed Behavior: Removes all outbox messages with status PROCESSED from the DB.

Observed Implementation:
- Acquires DB write lock.
- Iterates over outbox map, deletes entries where msg.Status == MessageStatusProcessed.
- Returns count of purged messages.

Assessment: PASS. Correctly removes only processed messages, leaving pending ones intact.

Severity: -
Notes: Test `TestTransactionalOutbox_PurgeProcessed` verifies this behavior.

## Finding 7: Demo Matches Implementation

Location: `labs/21-outbox-pattern/cmd/demo/main.go`

Claimed Behavior: Demo illustrates dual-write inconsistency, then atomic outbox success with relay dispatch and consumer idempotency.

Observed Implementation:
- Scenario 1: Sets broker to fail next publish, calls CreateOrderDualWriteNaive (order write then broker write outside tx). Shows order persisted, broker empty.
- Scenario 2: Starts relay, calls CreateOrderWithOutbox (atomic Tx), sleeps for relay to poll, shows broker receives one message, consumer accepts it.
- Scenario 3: Re-delivers the same published message to consumer, shows consumer rejects duplicate.

Assessment: PASS. Demo output matches code behavior and proves the claimed properties.

Severity: -
Notes: None.

## Finding 8: Latent Defect - Relay Stop Double-Close Panic

Location: `labs/21-outbox-pattern/internal/outbox/relay.go` (Stop)

Claimed Behavior: Relay can be started and stopped cleanly.

Observed Implementation:
- Stop() closes the stopChan channel.
- If Stop() is called twice, the second call will panic: `close of closed channel`.

Assessment: WARNING. This is a latent defect not triggered by current usage (Stop called once via defer in tests/demo). However, if Stop were invoked multiple times (e.g., in error-handling paths), it would panic.

Severity: LOW
Notes: Could be fixed by checking if stopChan is already closed or using a closed channel sentinel. Not exercised in current code path.