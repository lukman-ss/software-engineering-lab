# Test Audit: Coverage of Happy Path, Failure Path, Edge/Transition Cases

Test suite: `labs/21-outbox-pattern/tests/outbox_test.go`
Run results:

```bash
go test ./...           # PASS (tests)
go test -race ./...     # PASS (tests) 1.515s
```

## Finding 1: Happy Path Coverage

Location: `TestTransactionalOutbox_HappyPath`

Claimed Behavior: Full happy path — atomic order+outbox write, relay dispatch, broker publish, consumer accept.

Observed Implementation:
- Creates DB, broker, relay, service, consumer.
- Starts relay, calls CreateOrderWithOutbox.
- Asserts order persisted in DB (GetOrder).
- Sleeps 50ms, asserts broker received exactly 1 message (GetPublished, len==1, EventType==OrderCreated).
- Asserts outbox message marked PROCESSED (GetOutbox, Status==MessageStatusProcessed).
- Asserts consumer processed the message (Handle returns true, GetReceivedCount==1).

Assessment: PASS. Covers happy path end-to-end including persistence, dispatch, publish, and idempotent consumption.

Severity: -
Notes: Uses sleep-based polling (50ms) rather than deterministic synchronization. Works but is timing-dependent; could flake under slow CI.

## Finding 2: Rollback Coverage

Location: `TestTransactionalOutbox_Rollback`

Claimed Behavior: Transaction rollback discards staged data.

Observed Implementation:
- Begins a Tx, stages an Order and an OutboxMessage.
- Calls Rollback.
- Asserts neither order nor outbox persisted (GetOrder/GetOutbox return not found).
- Sleeps 30ms, asserts no messages sent to broker.

Assessment: PASS. Directly exercises the Tx rollback path and confirms staged data is discarded.

Severity: -
Notes: Does not test the service-level rollback path (CreateOrderWithOutbox rollback on marshal error), which is unreachable anyway.

## Finding 3: Idempotency / Duplicate Delivery Coverage

Location: `TestTransactionalOutbox_Idempotency_DuplicateDelivery`

Claimed Behavior: Consumer rejects duplicate message delivery.

Observed Implementation:
- Creates a consumer and a single message.
- First Handle returns true (new message).
- Second Handle returns false (duplicate rejected).
- Asserts GetReceivedCount == 1.

Assessment: PASS. Demonstrates idempotent consumer.

Severity: -
Notes: None.

## Finding 4: Dual-Write Failure (Naive Approach) Coverage

Location: `TestDualWriteProblem_Failure`

Claimed Behavior: Dual-write naive approach leaves DB/broker inconsistent on broker failure.

Observed Implementation:
- Sets broker failNext; calls CreateOrderDualWriteNaive.
- Asserts error returned.
- Asserts order persisted in DB (GetOrder found).
- Asserts broker has zero published messages.

Assessment: PASS. Proves the dual-write vulnerability.

Severity: -
Notes: None.

## Finding 5: Concurrency Coverage

Location: `TestTransactionalOutbox_ConcurrentWrites`

Claimed Behavior: Concurrent writes to the transactional DB do not cause data races.

Observed Implementation:
- 10 goroutines each call CreateOrderWithOutbox 10 times (100 total calls), but all use the SAME orderID ("o-concurrent-1") and thus the SAME outbox message ID ("evt-o-concurrent-1").
- Test waits, sleeps 50ms for relay to catch up.
- No assertion on persisted order/message counts or correctness — the test only checks that the race detector reports no races.

Assessment: WARNING. This is a race-smoke test, not a correctness test:
- All goroutines write the same key, so they overwrite each other (no meaningful concurrency of distinct records).
- No assertion that all intended writes are persisted or that all distinct events are dispatched.
- The naming implies concurrency correctness, but it merely survives concurrent writes via locking; it does not prove N distinct concurrent writes are all durable.

Severity: MEDIUM (incomplete coverage)
Notes: A stronger test would write distinct IDs and assert that all N orders/messages are persisted and dispatched. Under load, this test's lack of assertions weakens the concurrency guarantee.

## Finding 6: Purge Coverage

Location: `TestTransactionalOutbox_PurgeProcessed`

Claimed Behavior: Purge removes processed but not pending outbox messages.

Observed Implementation:
- Stages one pending (m1) and one processed (m2) message in a Tx, commits.
- Calls PurgeProcessedOutbox, asserts purged count == 1.
- Asserts m1 (pending) still retrievable, m2 (processed) removed.

Assessment: PASS. Covers purge for processed vs pending.

Severity: -
Notes: None.

## Finding 7: Missing Coverage - Relay Retry on Broker Failure

Location: `labs/21-outbox-pattern/internal/outbox/relay.go` (PollAndDispatch)

Claimed Behavior: (Should) Re-deliver outbox messages when broker publish fails, leaving them pending.

Observed Implementation: No test sets broker failure while the relay is running. The only broker-failure test (TestDualWriteProblem_Failure) uses the naive service path, not the relay. Nothing asserts that:
- a pending outbox message stays pending when Publish fails inside the relay, or
- the same message is retried and eventually succeeds on a subsequent poll.

Assessment: FAIL coverage. A core durability property of the outbox pattern (at-least-once delivery with retry) is exercised by code but entirely unproven by tests.

Severity: HIGH (core behavior unproven)
Notes: Test would need a broker that fails first N publishes then succeeds, with the relay running, asserting the message is retried and eventually published and marked processed.

## Finding 8: Missing Coverage - Relay Stop / Goroutine Cleanup

Location: `labs/21-outbox-pattern/internal/outbox/relay.go`

Observed Implementation: No test verifies that Stop() terminates the relay goroutine (no goroutine leak check) or that Stop is idempotent. Relays are started/stopped via defer in each test but the goroutine's exit is only timing-guaranteed, not asserted.

Assessment: WARNING. Weak coverage of lifecycle/cleanup.

Severity: LOW
Notes: Minor; no observed leak, but not asserted.

## Test Summary

| Test Name | Covers | Assessment |
|---|---|---|
| TestTransactionalOutbox_HappyPath | happy path | PASS |
| TestTransactionalOutbox_Rollback | rollback | PASS |
| TestTransactionalOutbox_Idempotency_DuplicateDelivery | idempotency | PASS |
| TestDualWriteProblem_Failure | naive failure path | PASS |
| TestTransactionalOutbox_ConcurrentWrites | concurrency (race smoke) | WARNING (weak assertion) |
| TestTransactionalOutbox_PurgeProcessed | purge | PASS |

Not covered: relay retry-on-broker-failure; relay stop/idempotency; service-level rollback path; distinct-concurrent-writes correctness.
