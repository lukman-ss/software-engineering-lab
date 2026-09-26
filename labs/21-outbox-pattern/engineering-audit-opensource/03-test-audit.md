## Test Coverage Analysis

### Tests Reviewed
File: `tests/outbox_test.go` (170 lines)

Tests present:
1. `TestTransactionalOutbox_HappyPath` (lines 11-11-59)
   - Happy path: CreateOrderWithOutbox, verify order persisted, relay dispatches, broker receives, consumer processes.
2. `TestTransactionalOutbox_Rollback` (lines-61-90)
   - Rollback path: BeginTx, SaveOrder, SaveOutbox, Rollback, verify nothing persisted, no broker publish.
3. `TestTransactionalOutbox_Idempotency_DuplicateDelivery` (lines-92-115)
   - Consumer idempotency: Handle same message twice, second call returns false, received count = 1.
4. `TestDualWriteProblem_Failure` (lines-116-139)
   - Dual-write failure: Broker failNext=true, CreateOrderDualWriteNaive, error returned, order in DB, no broker publish.
5. `TestTransactionalOutbox_ConcurrentWrites` (lines-141-170)
   - Concurrency: 10 goroutines each creating 10 orders via CreateOrderWithOutbox, wait, sleep for relay, no assertions on final state.

### Findings

#### Test Cases Verified
- **Happy path**: ✓ Covered by `TestTransactionalOutbox_HappyPath`
- **Rollback safety**: ✓ Covered by `TestTransactionalOutbox_Rollback` (verifies atomic rollback discards both order and outbox)
- **Duplicate delivery / Idempotency**: ✓ Covered by `TestTransactionalOutbox_Idempotency_DuplicateDelivery` (consumer-level) and partially by relay not marking as processed on broker failure (indirect)
- **Dual-write failure demonstration**: ✓ Covered by `TestDualWriteProblem_Failure` and demo Scenario 1
- **Concurrent execution safety**: Partially covered by `TestTransactionalOutbox_ConcurrentWrites` (runs with race detector but asserts nothing about correctness)

#### Missing or Weak Test Coverage
1. **MISSING_TEST**: Relay retry after broker failure
   - No test verifies that when `broker.Publish` fails, the message remains PENDING and is successfully published on a subsequent poll cycle.
   - The `failNext` mechanism is only used in dual-write test; not exercised in relay context.
   - **Severity**: MEDIUM (core behavior of at-least-once delivery not explicitly tested)

2. **MISSING_TEST**: Relay handling of MarkOutboxProcessed failure
   - No test verifies behavior when `db.MarkOutboxProcessed` fails (e.g., message not found) after successful broker publish.
   - **Severity**: LOW (edge case, but part of at-least-once delivery flow)

3. **WEAK_TEST**: ConcurrentWrites test lacks correctness assertions
   - `TestTransactionalOutbox_ConcurrentWrites` runs 100 concurrent CreateOrderWithOutbox calls and only relies on the race detector.
   - It does not assert that exactly 100 orders exist in DB, 100 outbox messages are created, 100 broker publishes occur, or that the consumer receives 100 unique messages.
   - **Severity**: MEDIUM (concurrent correctness not proven)

4. **WEAK_TEST**: No test for Stop() double-close
   - No test calls `Stop()` twice to verify behavior (currently would panic).
   - **Severity**: LOW (latent bug)

### Test Quality Assessment
- The test suite covers the primary failure and success scenarios for the outbox pattern.
- It correctly verifies atomicity (happy path), rollback, dual-write inconsistency, and consumer idempotency.
- The concurrent test is weak—it ensures no data races but not that concurrent transactions produce the expected outcome.
- The relay's retry logic (at-least-once delivery) is implied by code but not explicitly tested.
- All tests pass and the race detector passes.