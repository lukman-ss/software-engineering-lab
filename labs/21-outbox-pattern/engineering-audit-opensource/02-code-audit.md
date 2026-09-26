# Code Audit

## Finding 1: Atomic Transaction Commit

Location: internal/outbox/db.go:112-129, internal/outbox/service.go:17-53
Claimed Behavior: Order and outbox message are persisted atomically in a single DB transaction
Observed Implementation: Tx struct stages WritesOrder and SaveOutbox mutations in separate maps, Commit() locks both Tx mutex and DB mutex, writes all staged orders then all staged outbox messages, sets closed=true before applying. Rollback() sets closed=true and discards staged mutations.
Assessment: PASS
Severity: LOW
Notes: Staged mutation maps correctly isolate in-flight writes. No partial commit possible since Commit applies all staged writes atomically under DB lock.

## Finding 2: Rollback Discards No Order Saved

Location: internal/outbox/db.go:59-68, internal/outbox/service.go:17-53
Claimed Behavior: On rollback, neither order nor outbox is persisted
Observed Implementation: Rollback does not touch DB maps; staged mutations discarded when Tx is garbage collected.
Assessment: PASS
Severity: LOW
Notes: Correct behavior. Test TestTransactionalOutbox_Rollback confirms no order/outbox persisted after rollback.

## Finding 3: Relay Polling and Dispatch Marking

Location: internal/outbox/relay.go:43-60
Claimed Behavior: Relay polls pending outbox records, publishes to broker, marks as processed
Observed Implementation: PollAndDispatch fetches PENDING messages, publishes each to broker. On success, calls MarkOutboxProcessed. On broker publish failure, logs and leaves status unchanged (re-poll). On mark-processed failure, logs error but still increments dispatched (BUG - message may not be persisted as PROCESSED).
Assessment: WARNING
Severity: MEDIUM
Notes: Race condition: dispatched++ is incremented even when MarkOutboxProcessed fails (relay.go:50-54). This could cause double-publish on retry because the outbox record stays PENDING, contradicting the at-least-once guarantee of idempotent consumers. Fix: only increment dispatched after successful MarkOutboxProcessed.

## Finding 4: Relay Stop Race

Location: internal/outbox/relay.go:24-41
Claimed Behavior: Relay stops cleanly on Stop()
Observed Implementation: Start() launches goroutine that selects on ticker.C and stopChan. Stop() closes stopChan. Goroutine returns when stopChan closed.
Assessment: PASS
Severity: LOW
Notes: No deadlock. Goroutine exits cleanly.

## Finding 5: Consumer Idempotency

Location: internal/outbox/consumer.go:19-31
Claimed Behavior: Consumer tracks processed event IDs for idempotency
Observed Implementation: processedIDs map checked under mutex; duplicates return false, new messages accepted and stored.
Assessment: PASS
Severity: LOW
Notes: Correct idempotent consumer implementation.

## Finding 6: Broker Fault Injection

Location: internal/outbox/broker.go:22-37
Claimed Behavior: MockBroker simulates publish failures via SetFailNext
Observed Implementation: failNext bool toggled; Publish returns BrokerError and resets failNext to false on next call.
Assessment: PASS
Severity: LOW
Notes: Supports dual-write failure test and relay retry scenarios.

## Finding 7: No Retry Logic in Relay

Location: internal/outbox/relay.go:43-60
Claimed Behavior: (Design) Relay retries on transient failures
Observed Implementation: Relay does not retry failed publishes. On broker failure, message remains PENDING and is re-attempted on next poll interval.
Assessment: PASS (with note)
Severity: LOW
Notes: Re-poll provides implicit retry. No explicit retry/backoff. Acceptable for demo purposes.

## Finding 8: Dual-Write Demonstration Correctness

Location: internal/outbox/service.go:55-90, cmd/demo/main.go:19-27
Claimed Behavior: Dual-write naive path commits order to DB then attempts broker publish; if broker fails, DB inconsistency results
Observed Implementation: Service saves order in tx, commits tx, then attempts broker.Publish. On broker failure returns error but DB already committed.
Assessment: PASS
Severity: LOW
Notes: Correctly demonstrates dual-write anti-pattern.

## Finding 9: Concurrent Write Consistency

Location: internal/outbox/db.go:59-68, internal/outbox/service.go:17-53
Claimed Behavior: Concurrent order creation is thread-safe
Observed Implementation: Each BeginTx() creates independent Tx with its own staged maps; Commit locks DB mutex. GetPendingOutbox uses RLock. MarkOutboxProcessed uses Lock.
Assessment: WARNING
Severity: MEDIUM
Notes: Concurrent orders with same ID will have last-writer-wins on commit. Test uses same orderID across workers (TestTransactionalOutbox_ConcurrentWrites), so final DB state has only one unique order. Not a data race but may not reflect realistic concurrency. No explicit conflict/duplicate detection.

## Finding 10: PurgeProcessed Outbox

Location: internal/outbox/db.go:71-82
Claimed Behavior: PurgeProcessedOutbox removes PROCESSED messages
Observed Implementation: Iterates outbox map, deletes PROCESSED entries under Lock.
Assessment: PASS
Severity: LOW
Notes: Correct cleanup behavior.