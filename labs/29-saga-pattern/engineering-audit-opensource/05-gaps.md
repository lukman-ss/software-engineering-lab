# Gap Analysis

# Audit Execution Results

## Compilation
Command: go vet ./...
Result: PASS

## Tests
Command: go test ./...
Result: PASS

## Concurrency / Race Detection
Command: go test -race ./...
Result: PASS (no race conditions detected)

## Demo
Command: go run ./cmd/demo
Result: PASS (both scenarios execute as described)

---

## Identified Gaps

### Gap 1: OrderService ApproveOrder Lock Leak
Type: MISSING_EDGE_CASE / UNHANDLED_ERROR
Severity: MEDIUM
Location: internal/services/services.go lines 42-51

- ApproveOrder returns an error early without releasing the semantic lock if the order is already cancelled.
- The lock remains stuck for the orderID until CancelOrder is called.
- If the order was cancelled by a path that does not call CancelOrder (or a failure between CreateOrder and CancelOrder), the lock leaks permanently.
- The current test suite only tests CreateOrder/CreateOrder failure (lock set, then second CreateOrder fails). It does not test ApproveOrder failure after CreateOrder and verify lock release. The test TestChoreography_Flow and TestChoreography_FailureCompensates test approve via event bus but not the lock-leak scenario.

### Gap 2: No Test for Concurrent EventBus Access
Type: MISSING_TEST
Severity: LOW
Location: internal/saga/choreography.go / tests/saga_test.go

- EventBus Publish and Subscribe mutate shared handler map under mutex.
- While the implementation uses copy-on-publish (RWMutex + copy), there is no test exercising concurrent publish/subscribe to verify safety.
- The race detector does not catch data races in the test suite as written because no concurrent access is exercised.

### Gap 3: No Test for Mid-Sequence Compensation Failure with Subsequent Steps Running
Type: MISSING_TEST
Severity: LOW
Location: internal/saga/orchestrator.go / tests/saga_test.go

- TestOrchestrator_CompensationErrorPropagated tests a single compensation step failing.
- However, it does not verify that LIFO compensations continue for steps preceding the failing compensation (i.e., if compensation of step 2 fails, does step 1 still get compensated?).
- The orchestrator code in compensate() does NOT return early on compensation error (it collects errors and continues), which is correct. But this behavior is not explicitly verified by a test.
- This is a recovery/rollback edge case that should be tested.

### Gap 4: Compensation Context Discards Original Context
Type: UNHANDLED_ERROR (minor)
Severity: LOW
Location: internal/saga/orchestrator.go lines 63, 78

- Compensation uses context.Background() even when the original context was cancelled.
- If compensation steps should respect the original timeout/cancellation, this is incorrect.
- This is a design decision (compensations run to completion) but is not documented in code comments beyond "ponytail:" note in code.

### Gap 5: Orchestrator Logs Do Not Track Pending State
Type: MISSING_EDGE_CASE
Severity: LOW
Location: internal/saga/orchestrator.go Execute method

- Logs record only final status (EXECUTED, FAILED, COMPENSATED, COMPENSATE_FAILED).
- There is no PENDING/IN PROGRESS log entry, meaning logs don't show step start.
- This is acceptable for a demo but limits observability for auditing purposes.
- StatusPending is defined but never written.

### Gap 6: Missing Test for Choreography Compensation Error Propagation
Type: MISSING_TEST
Severity: LOW
Location: tests/saga_test.go

- TestChoreography_FailureCompensates tests the happy compensation path (refund + cancel).
- There is no test where a choreography compensation handler itself fails, and there's no mechanism for the event bus to aggregate/propagate compensation failures. This is a limitation of the choreography model as implemented.

# Summary of Gaps

| # | Gap | Type | Severity | Status |
|---|-----|------|----------|--------|
| 1 | OrderService ApproveOrder lock leak | MISSING_EDGE_CASE | MEDIUM | Open |
| 2 | No concurrent EventBus test | MISSING_TEST | LOW | Open |
| 3 | No mid-sequence compensation failure test | MISSING_TEST | LOW | Open |
| 4 | Compensation discards original context | UNHANDLED_ERROR | LOW | Open |
| 5 | Logs lack Pending state | MISSING_EDGE_CASE | LOW | Open |
| 6 | No choreography compensation error test | MISSING_TEST | LOW | Open |

# Overall Quality

The implementation matches its claimed behavior (orchestration + choreography, LIFO rollback, idempotency, semantic locking, demo happy/failure paths). All provided tests pass and no race conditions are detected.

The identified gaps are primarily around edge-case coverage and a potential lock-leak in ApproveOrder. None rise to the level of fabricated results or fake benchmarks. The MEDIUM severity item (lock leak) should be addressed before considering the lab fully trustworthy for technical writing.
