## Gap Analysis

### MISSING_TEST: Relay retry after broker failure
- **Location**: tests/outbox_test.go
- **Description**: No test verifies that when `broker.Publish` returns an error, the outbox message remains in `PENDING` status and is successfully published on a subsequent poll cycle.
- **Severity**: MEDIUM
- **Notes**: The relay logs the failure and continues; the message is retried on the next tick. This behavior is implied by code but not explicitly tested.

### MISSING_TEST: Relay handling of MarkOutboxProcessed failure
- **Location**: tests/outbox_test.go
- **Description**: No test verifies behavior when `db.MarkOutboxProcessed` fails (e.g., message not found) after a successful `broker.Publish`.
- **Severity**: LOW
- **Notes**: Edge case where publish succeeds but status update fails; message remains PENDING and will be re-published, relying on consumer idempotency.

### MISSING_TEST: Stop() double-close behavior
- **Location**: internal/outbox/relay.go:39-41
- **Description**: No test calls `Stop()` twice to verify behavior; currently, a second call to `Stop()` will panic with "close of closed channel".
- **Severity**: LOW
- **Notes**: Latent bug not exercised by current tests/demo (single Stop per lifecycle).

### MISSING_TEST: Concurrent writes correctness assertions
- **Location**: tests/outbox_test.go:141-170
- **Description**: `TestTransactionalOutbox_ConcurrentWrites` runs 100 concurrent `CreateOrderWithOutbox` calls and relies solely on the race detector. It does not assert correctness outcomes such as: exactly 100 orders in DB, 100 outbox messages created, 100 broker publishes, or consumer receiving 100 unique messages.
- **Severity**: MEDIUM
- **Notes**: While no data races exist, the test does not prove functional correctness under concurrency.

### DOC_CODE_MISMATCH: README omits broker failure mechanism details
- **Location**: README.md
- **Description**: README mentions "Thread-safe mock message broker simulating publish failures" but does not detail the `SetFailNext`/`failNext` mechanism or the `BrokerError` type used to simulate failures.
- **Severity**: LOW
- **Notes**: Minor documentation gap; the implementation correctly simulates failures, but the README lacks specificity.

### Summary of Gap Types
- MISSING_TEST: 4 instances
- DOC_CODE_MISMATCH: 1 instance
- BROKEN_IMPLEMENTATION: 0
- RACE_CONDITION: 0
- UNHANDLED_ERROR: 0 (Stop() double-close categorized as MISSING_TEST for test omission)
- MISSING_EDGE_CASE: 0 (covered by above MISSING_TEST items)
- IMPLEMENTATION_OVERCLAIM: 0
- RESEARCH_MISMATCH: N/A (not audited per pipeline override)
- FAKE_DEMO: 0
- FAKE_BENCHMARK: 0
- UNVERIFIED_RESULT: 0