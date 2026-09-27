# Code Audit

## Finding 1

Location: internal/outbox/db.go, lines 12-16 (DB struct)
Claimed Behavior: The DB struct provides thread-safe storage for orders and outbox messages using RWMutex.
Observed Implementation: DB uses a single RWMutex for both orders and outbox maps. GetOrder and GetOutbox use RLock, while Mutating methods (MarkOutboxProcessed, PurgeProcessedOutbox) use Lock. Tx uses a separate mutex for staging and delegates to DB mutex for commit.
Assessment: PASS
Severity: LOW
Notes: The locking strategy is correct and avoids deadlock by not holding DB mutex while acquiring Tx mutex (Tx methods lock Tx.mu then DB.mu only in Commit). However, note that GetPendingOutbox uses RLock and iterates while holding the lock, which is fine for short iterations.

## Finding 2

Location: internal/outbox/db.go, lines 84-140 (Tx struct and methods)
Claimed Behavior: Tx provides staged staging for orders and outbox, with atomic commit or rollback.
Observed Implementation: Tx.mu protects stagedOrders and stagedOutbox. Commit locks Tx.mu then DB.mu, copies staged data to DB, and sets closed=true. Rollback locks Tx.mu, sets closed=true, discards staged data. SaveOrder and SaveOutbox lock Tx.mu and check closed flag.
Assessment: PASS
Severity: LOW
Notes: The Tx implementation correctly isolates staged changes until commit. The double locking (Tx.mu then DB.mu) is safe because Tx.mu is always acquired first, preventing deadlock. The closed flag prevents reuse after commit/rollback.

## Finding 3

Location: internal/outbox/broker.go, lines 10-26 (MockBroker struct)
Claimed Behavior: MockBroker is thread-safe and simulates publish failures via failNext flag.
Observed Implementation: Uses RWMutex to protect published slice and failNext flag. Publish locks mutex, checks failNext, returns error if set, otherwise appends message. SetFailNext locks mutex to update flag.
Assessment: PASS
Severity: LOW
Notes: Correct use of mutex for thread safety. The failNext flag allows simulating a single failure for retry tests.

## Finding 4

Location: internal/outbox/service.go, lines 17-53 (CreateOrderWithOutbox)
Claimed Behavior: Writes order and outbox event in same transaction, rolling back on any error.
Observed Implementation: Begins Tx, creates order and outbox message, saves both via Tx, commits on success. On any error (including marshaling, SaveOrder, SaveOutbox), rolls back Tx and returns error.
Assessment: PASS
Severity: LOW
Notes: Proper error handling and rollback on all error paths. The rollback is invoked via _ = tx.Rollback() to ignore rollback errors (which are only possible if Tx already closed, but we check errors before calling Rollback in the same flow).

## Finding 5

Location: internal/outbox/service.go, lines 55-90 (CreateOrderDualWriteNaive)
Claimed Behavior: Demonstrates dual-write vulnerability by writing order in DB transaction then attempting broker write outside transaction.
Observed Implementation: Begins Tx, saves order, commits Tx, then creates message and calls broker.Publish. If broker fails, returns error but order remains in DB.
Assessment: PASS
Severity: LOW
Notes: Correctly models the dual-write problem. The function returns the broker error, allowing the demo to show inconsistency.

## Finding 6

Location: internal/outbox/relay.go, lines 9-24 (Relay struct)
Claimed Behavior: Relay polls pending outbox events at intervals, dispatches to broker, and marks as processed.
Observed Implementation: Relay holds DB, Broker, pollInterval, stopChan, and once guards for Start/Stop. Start launches a goroutine with ticker that calls PollAndDispatch. Stop closes stopChan.
Assessment: PASS
Severity: LOW
Notes: Proper use of sync.Once ensures Start and Stop are idempotent. The goroutine is properly stopped via channel close.

## Finding 7

Location: internal/outbox/relay.go, lines 50-66 (PollAndDispatch)
Claimed Behavior: Fetches pending outbox messages, attempts to publish each to broker, on success marks as processed and increments dispatched count.
Observed Implementation: Gets pending messages (db.GetPendingOutbox uses RLock). For each msg, calls broker.Publish. If no error, calls db.MarkOutboxProcessed and logs any error from marking but still counts as dispatched? Note: dispatched is incremented only if broker.Publish succeeds AND mark processed succeeds? Actually, code increments dispatched only if broker.Publish succeeds and then if marking processed succeeds? Wait: lines 55-61: if err == nil (broker success) then err = r.db.MarkOutboxProcessed(msg.ID); if err != nil { log ... } else { dispatched++ }. So dispatched counts only those successfully published AND marked processed.
Assessment: WARNING
Severity: MEDIUM
Notes: The relay counts a message as dispatched only if both broker publish and mark processed succeed. However, if broker publish succeeds but marking processed fails, the message will be retried in next poll (since status remains PENDING). This is correct behavior (at-least-once delivery) but the dispatched count may undercount. However, the function returns dispatched count for metrics; the test does not rely on this count. The behavior is acceptable for at-least-once semantics.

## Finding 7b (continuation)

Location: internal/outbox/relay.go, lines 62-64
Claimed Behavior: On broker publish failure, log error and do not mark processed.
Observed Implementation: Logs failure and does not increment dispatched or attempt to mark processed.
Assessment: PASS
Severity: LOW
Notes: Correctly leaves message as PENDING for retry.

## Finding 8

Location: internal/outbox/consumer.go, lines 5-9 (Consumer struct)
Claimed Behavior: Consumer tracks processed message IDs to enforce idempotency.
Observed Implementation: Uses mutex to protect processedIDs map and received slice. Handle checks if ID exists in map, if so returns false (duplicate), else sets true, appends to received, returns true.
Assessment: PASS
Severity: LOW
Notes: Standard idempotency implementation. Note that received slice is appended only for new messages, so it tracks unique messages processed.

## Finding 9

Location: internal/outbox/consumer.go, lines 18-31 (Handle)
Claimed Behavior: Returns true if processed, false if duplicate.
Observed Implementation: As described.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 10

Location: tests/outbox_test.go
Claimed Behavior: Test suite covers happy path, rollback, idempotency, dual-write failure, concurrent writes, purge, broker retry, concurrent consumers.
Observed Implementation: All tests pass, including race detector.
Assessment: PASS
Severity: LOW
Notes: Tests are comprehensive and cover the claimed behaviors. No races detected.

## Finding 11

Location: cmd/demo/main.go
Claimed Behavior: Demo shows dual-write inconsistency, transactional outbox success, and idempotent consumer handling duplicates.
Observed Implementation: Runs scenarios as described, prints expected outputs.
Assessment: PASS
Severity: LOW
Notes: Demo matches code behavior and is deterministic.

## Finding 12

Location: internal/outbox/db.go, lines 47-57 (GetPendingOutbox)
Claimed Behavior: Returns a slice of pending outbox messages.
Observed Implementation: Uses RLock, iterates over outbox map, appends those with Status == MessageStatusPending.
Assessment: PASS
Severity: LOW
Notes: Correctly returns pending messages. Note that it returns a copy of the slice, so external modifications do not affect internal state.

## Finding 13

Location: internal/outbox/db.go, lines 59-69 (MarkOutboxProcessed)
Claimed Behavior: Marks a specific outbox message as processed.
Observed Implementation: Locks mutex, checks existence, updates status, unlocks.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 14

Location: internal/outbox/db.go, lines 71-82 (PurgeProcessedOutbox)
Claimed Behavior: Removes processed outbox messages and returns count purged.
Observed Implementation: Locks mutex, iterates over outbox map, deletes entries with Status == MessageStatusProcessed, increments count.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 15

Location: internal/outbox/broker.go, lines 28-37 (Publish)
Claimed Behavior: Simulates broker failure when failNext is true, otherwise appends message.
Observed Implementation: Locks mutex, if failNext returns error and resets failNext to false, else appends message.
Assessment: PASS
Severity: LOW
Notes: The failNext flag is reset after failure, allowing exactly one failure simulation. This is suitable for tests.

## Finding 16

Location: internal/outbox/relay.go, lines 27-41 (Start)
Claimed Behavior: Starts the polling relay goroutine idempotently.
Observed Implementation: Uses sync.Once to launch goroutine that runs ticker loop calling PollAndDispatch.
Assessment: PASS
Severity: LOW
Notes: Correct idempotent start.

## Finding 17

Location: internal/outbox/relay.go, lines 44-48 (Stop)
Claimed Behavior: Stops the relay goroutine idempotently.
Observed Implementation: Uses sync.Once to close stopChan.
Assessment: PASS
Severity: LOW
Notes: Correct idempotent stop.

## Finding 18

Location: internal/outbox/service.go, lines 17-53 (CreateOrderWithOutbox) - error handling
Claimed Behavior: On error during order or outbox save, rolls back transaction.
Observed Implementation: After each SaveOrder and SaveOutbox, if error, calls _ = tx.Rollback() and returns error.
Assessment: PASS
Severity: LOW
Notes: The rollback error is ignored because the only possible error is ErrTxClosed, which would indicate a double rollback or commit, but in this flow we know tx is not closed. However, it is safe to ignore.

## Finding 19

Location: internal/outbox/service.go, lines 55-90 (CreateOrderDualWriteNaive) - error handling
Claimed Behavior: On broker publish error, returns error but order remains in DB.
Observed Implementation: After Tx commit, calls broker.Publish, if error returns error wrapped.
Assessment: PASS
Severity: LOW
Notes: Correctly models the dual-write failure.

## Finding 20

Location: internal/outbox/relay.go, lines 50-66 (PollAndDispatch) - error handling for MarkOutboxProcessed
Claimed Behavior: If marking processed fails, log error and do not count as dispatched.
Observed Implementation: Logs error and does not increment dispatched.
Assessment: PASS
Severity: LOW
Notes: The message remains PENDING and will be retried. This is correct.

## Finding 21

Location: internal/outbox/db.go, lines 84-140 (Tx) - potential issue
Claimed Behavior: Tx methods lock Tx.mu to protect staged data.
Observed Implementation: All Tx methods (SaveOrder, SaveOutbox, Commit, Rollback) lock Tx.mu.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 22

Location: internal/outbox/db.go, lines 112-129 (Commit)
Claimed Behavior: Commits staged data to DB under lock.
Observed Implementation: Locks Tx.mu, checks closed, sets closed=true, locks DB.mu, copies staged orders and outbox to DB maps, unlocks.
Assessment: PASS
Severity: LOW
Notes: The order of locking (Tx.mu then DB.mu) is consistent and avoids deadlock. However, note that the Tx.mu is held while acquiring DB.mu, which is fine as long as no other goroutine tries to acquire DB.mu then Tx.mu. The DB methods that need DB.mu (like GetPendingOutbox) only take DB.mu (RLock) and do not take Tx.mu, so no deadlock.

## Finding 23

Location: internal/outbox/db.go, lines 33-38 (GetOrder) and 40-45 (GetOutbox)
Claimed Behavior: Thread-safe reads.
Observed Implementation: Use RLock.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 24

Location: internal/outbox/db.go, lines 47-56 (GetPendingOutbox)
Claimed Behavior: Thread-safe iteration to collect pending messages.
Observed Implementation: Uses RLock, iterates over outbox map, checks status.
Assessment: PASS
Severity: LOW
Notes: Correct. Note that holding RLock for the duration of iteration is acceptable for small to medium maps; for very large maps it could block writers, but this is a demo.

## Finding 25

Location: internal/outbox/broker.go, lines 39-45 (GetPublished)
Claimed Behavior: Thread-safe copy of published slice.
Observed Implementation: Uses RLock, creates new slice, copies.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 26

Location: internal/outbox/consumer.go, lines 33-37 (GetReceivedCount) and 39-43 (IsProcessed)
Claimed Behavior: Thread-safe reads.
Observed Implementation: Use mutex lock.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 27

Location: internal/outbox/consumer.go, lines 18-31 (Handle)
Claimed Behavior: Thread-safe check-and-set for idempotency.
Observed Implementation: Locks mutex, checks map, updates if not present.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 28

Location: tests/outbox_test.go, line 19 (relay.Start) and line 20 (defer relay.Stop)
Claimed Behavior: Relay started before test operations and stopped after.
Observed Implementation: Relay started, deferred stop.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 29

Location: tests/outbox_test.go, line 35 (time.Sleep(50 * time.Millisecond))
Claimed Behavior: Wait for relay to process.
Observed Implementation: Fixed sleep.
Assessment: WARNING
Severity: LOW
Notes: Using fixed sleep can lead to flaky tests under slow CI. However, the test passes consistently in local runs. Could improve with polling for condition, but not required for correctness.

## Finding 30

Location: tests/outbox_test.go, line 87 (time.Sleep(30 * time.Millisecond))
Claimed Behavior: Wait for relay to process (or not).
Observed Implementation: Fixed sleep.
Assessment: WARNING
Severity: LOW
Notes: Same as above.

## Finding 31

Location: tests/outbox_test.go, line 176-183 (ConcurrentWrites test waiting for relay)
Claimed Behavior: Wait for relay to process all orders with deadline.
Observed Implementation: Uses deadline loop with sleep.
Assessment: PASS
Severity: LOW
Notes: Better than fixed sleep; uses condition checking.

## Finding 32

Location: tests/outbox_test.go, line 241-249 (Retry test waiting for relay)
Claimed Behavior: Wait for relay to process after failure.
Observed Implementation: Uses fixed sleeps.
Assessment: WARNING
Severity: LOW
Notes: Same as above.

## Finding 33

Location: internal/outbox/relay.go, line 30 (ticker := time.NewTicker(r.pollInterval))
Claimed Behavior: Polling interval ticker.
Observed Implementation: Creates new ticker, defers ticker.Stop().
Assessment: PASS
Severity: LOW
Notes: Correctly stops ticker when goroutine exits.

## Finding 34

Location: internal/outbox/relay.go, line 36 (select { <-ticker.C, <-r.stopChan })
Claimed Behavior: Wait for ticker tick or stop signal.
Observed Implementation: Uses select.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 35

Location: internal/outbox/service.go, line 28 (payload, err := json.Marshal(order))
Claimed Behavior: Marshal order to JSON for outbox payload.
Observed Implementation: Standard marshal.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 36

Location: internal/outbox/service.go, line 34-40 (msg construction)
Claimed Behavior: Construct outbox message with UUID-like ID.
Observed Implementation: Uses fmt.Sprintf("evt-%s", orderID) for ID.
Assessment: PASS
Severity: LOW
Notes: Deterministic for testing; in real world would use UUID.

## Finding 37

Location: internal/outbox/service.go, line 79 (payload for dual-write)
Claimed Behavior: Construct payload for dual-write message.
Observed Implementation: Uses fmt.Sprintf(`{"ID":"%s","Amount":%f}`, orderID, amount).
Assessment: PASS
Severity: LOW
Notes: Simplified payload for demo.

## Finding 38

Location: internal/outbox/model.go (entire file)
Claimed Behavior: Define domain types.
Observed Implementation: Defines OrderStatus, Order, MessageStatus, OutboxMessage.
Assessment: PASS
Severity: LOW
Notes: Simple and correct.

## Finding 39

Location: go.mod
Claimed Behavior: Go module definition.
Observed Implementation: module github.com/software-engineering-lab/labs/21-outbox-pattern, go 1.22.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 40

Location: README.md
Claimed Behavior: Describes architecture and how to run tests and demo.
Observed Implementation: Matches the code structure and commands.
Assessment: PASS
Severity: LOW
Notes: README is accurate.