# Gap Analysis

## Allowed Gap Types
- MISSING_TEST
- BROKEN_IMPLEMENTATION
- DOC_CODE_MISMATCH
- RACE_CONDITION
- UNHANDLED_ERROR
- MISSING_EDGE_CASE
- IMPLEMENTATION_OVERCLAIM
- RESEARCH_MISMATCH
- FAKE_DEMO
- FAKE_BENCHMARK
- UNVERIFIED_RESULT

## Gaps Identified

### Gap 1

Type: MISSING_TEST
Location: `internal/outbox/relay.go` (PollAndDispatch) and `tests/outbox_test.go`
Priority: HIGH
Severity: HIGH
Description: The relay's at-least-once delivery / retry behavior is core to the outbox pattern but untested. There is no test that sets broker publish failure while the relay is running and verifies that:
- the outbox message stays PENDING,
- the message is retried on the next poll when the broker recovers,
- the message is eventually published and marked PROCESSED.
The code appears correct, but the behavior is unproven.
Rationale: "Core behavior unproven" per severity scale.

### Gap 2

Type: MISSING_TEST
Location: `tests/outbox_test.go` (TestTransactionalOutbox_ConcurrentWrites)
Priority: MEDIUM
Severity: MEDIUM
Description: The concurrent-writes test uses a single shared orderID ("o-concurrent-1") across all 10 goroutines, causing write overwrites and providing no assertion that distinct concurrent orders are all persisted/dispatched. The test is effectively a race-detector smoke test, not a correctness test. It proves absence of data races (corroborated by `-race` passing) but not concurrent write correctness for distinct keys.
Rationale: Incomplete coverage.

### Gap 3

Type: MISSING_TEST
Location: `internal/outbox/relay.go` (Stop) and `cmd/demo/main.go`
Priority: LOW
Severity: LOW
Description: Relay.Stop is not idempotent — double-close on stopChan would panic. Not exercised by tests/demo (Stop called once via defer). Also no goroutine-leak assertion for the relay worker goroutine.
Rationale: Latent defect, low severity.

### Gap 4

Type: MISSING_TEST
Location: `internal/outbox/service.go` (CreateOrderWithOutbox rollback path)
Priority: LOW
Severity: LOW
Description: The rollback-on-error path in CreateOrderWithOutbox (triggered by json.Marshal or Save failures) is unreachable in practice because json.Marshal of the Order struct cannot fail and the in-memory DB Save methods never return errors. The path exists but is dead code for this mock implementation; not a correctness risk.
Rationale: Dead/defensive path; low severity.

### Gap 5

Type: MISSING_EDGE_CASE
Location: `internal/outbox/relay.go` (PollAndDispatch)
Priority: LOW
Severity: LOW
Description: When the broker fails to publish a message, the relay logs and skips it but leaves it PENDING. If all broker publishes keep failing, the outbox table grows unbounded (no backoff, retry limit, or dead-letter handling). Acceptable for a demonstration/lab but omitted from documented guarantees.
Rationale: Minor; documented limits not claimed.

## Gaps Not Present

| Gap Type | Present? | Notes |
|---|---|---|
| BROKEN_IMPLEMENTATION | No | All implementation logic is correct (atomic commit, retry-leaves-pending, idempotency). |
| DOC_CODE_MISMATCH | No | README matches code and demo results. |
| IMPLEMENTATION_OVERCLAIM | No | README claims are bounded to the demo's in-memory simulation ("simulating BeginTx/Commit/Rollback"); no real DB claimed. |
| FAKE_DEMO | No | Demo runs, exit code 0, output matches description. |
| FAKE_BENCHMARK | No | No benchmarks referenced; none claimed. |
| UNVERIFIED_RESULT | No | All runtime results (tests, race, demo) were actually executed and recorded. |
| RACE_CONDITION | No | Confirmed absent via `go test -race` (PASS, 1.515s). |
| UNHANDLED_ERROR | No | Errors are propagated (Publish err returned, MarkOutboxProcessed err logged; not silently dropped). |
