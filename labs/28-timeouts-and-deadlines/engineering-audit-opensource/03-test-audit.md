# Test Audit — labs/28-timeouts-and-deadlines

## Coverage Matrix

| Package | Unit Tests | Concurrent | Failure/Edge | Transitions |
|---|---|---|---|---|
| deadline | success, timeout, parent-inherit | no | timeout, parent ctx | deadline vs fn-done race |
| retry | first-try, retry-until-success, exceed-max, ctx-canceled, jitter-bounds | no | zero-config defaults, ctx cancel | full retry loop |
| circuit | closed→open→half→closed, half-open-failure-trip, zero-config default | no | none | all state transitions |
| idempotency | get/set, expiration, lazy-eviction, concurrent access | yes (50x2) | lazy eviction correctness | expiration |
| integration | retry+circuit, idempotent-retry | no | backend-down persistence | end-to-end retry + CB trip |

## Finding 1 — Race detector clean

Location: go test -race ./... (executed)
Claimed Behavior: engineering/01-design.md Success Criteria 2 requires race detector clean under concurrent ops.
Observed Result: all packages ok (cached, then uncached run PASS).
Assessment: PASS
Severity: LOW
Notes: `internal/idempotency` concurrent test exercises real contention without races. No race-sensitive goroutine leak check in deadline package.

## Finding 2 — Happy path

Location: retry/retry_test.go:10, circuit/circuit_test.go:9, idempotency/idempotency_test.go:9
Claimed Behavior: Success transitions recorded correctly.
Observed Implementation: All positive paths verified.
Assessment: PASS
Severity: NONE
Notes: n/a.

## Finding 3 — Failure / negative path

Location: retry/retry_test.go:43, circuit/circuit_test.go:30, deadline/deadline_test.go:20
Claimed Behavior: MaxRetries exceeded returns joined error; OPEN breaker rejects; budget exceeded returns DeadlineExceeded.
Observed Implementation: Verified.
Assessment: PASS
Severity: NONE
Notes: `errors.Is` checks robust.

## Finding 4 — State-transition coverage

Location: circuit/circuit_test.go:9-49, 51-81
Claimed Behavior: All three transitions documented.
Observed Implementation: CLOSED→OPEN via failure threshold; OPEN→HALF_OPEN via cooldown; HALF_OPEN→CLOSED via success; HALF_OPEN→OPEN on failure.
Assessment: PASS (with caveat)
Severity: LOW
Notes: No test asserts HALF_OPEN→OPEN via *success-count-incomplete-stays-HALF_OPEN* (implicitly: one success then OPEN still OPEN). Boundary `failures == threshold` exact boundary tested (2=2). HALF_OPEN probe concurrency not asserted.

## Finding 5 — Edge cases missing

Location: whole suite
Claimed Behavior: Robustness under load.
Observed Implementation: No deadline test for fn racing ctx-done simultaneously with fn returning nil (leak path). No circuit test with zero-cooldown concurrency. No retry test for ctx-canceled *during backoff with fn returning nil first*.
Assessment: WARNING
Severity: LOW
Notes: Edge-case gaps are gaps, not failures of implemented claims.

## Summary

Suite proves all four claimed behaviors in single-threaded execution. Concurrency proven only for idempotency store. No test fabricates outcome — all run via `go test`.
