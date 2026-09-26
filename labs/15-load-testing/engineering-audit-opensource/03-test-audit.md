# Test Audit

## Test Coverage Summary

| Test File | Tests |
|-----------|-------|
| internal/loadtest/metrics_test.go | TestCalculateMetrics, TestCalculateMetrics_Empty, TestCalculateMetrics_Invariants (3) |
| tests/loadtest_test.go | TestLoadTest_SmokeVsStress, TestLoadTest_ErrorCount, TestServer_MethodNotAllowed, TestLoadTest_DialError, TestServer_ContextCanceled (5) |

Total: 8 tests, all passing.

## Finding 1
Location: tests/loadtest_test.go:14-58
Claimed Behavior: Smoke test should have lower latency than stress test; stress test should show tail latency exceeding average
Observed Implementation: The test compares smoke P95 vs stress P95, and stress P95 vs stress Avg. Both assertions pass.
Assessment: PASS
Severity: LOW
Notes: Validates core research claims: (1) stress causes higher tail latency, (2) averages mask tail spikes. Demonstrates smoke vs stress differentiation.

## Finding 2
Location: internal/loadtest/metrics_test.go:8-41
Claimed Behavior: Percentile calculations (P50, P95, P99, Avg) compute correctly
Observed Implementation: Uses 100 latency samples (1ms to 100ms) with known expected values. All assertions pass.
Assessment: PASS
Severity: LOW
Notes: Deterministic test with mathematical verification of percentile correctness.

## Finding 3
Location: tests/loadtest_test.go:60-85
Claimed Behavior: HTTP 500 errors are counted as errors, not successes
Observed Implementation: Mock server returns 500 for all requests. Test verifies ErrorCount == TotalRequests and SuccessCount == 0.
Assessment: PASS
Severity: LOW
Notes: Covers failure path. All errors are correctly counted.

## Finding 4
Location: tests/loadtest_test.go:87-101
Claimed Behavior: GET requests to /booking return 405 Method Not Allowed
Observed Implementation: Sends GET request, verifies response code is 405.
Assessment: PASS
Severity: LOW
Notes: Covers edge case for HTTP method validation.

## Finding 5
Location: tests/loadtest_test.go:103-120
Claimed Behavior: Dial errors are counted as errors
Observed Implementation: Requests sent to localhost:1 (unreachable port). Verifies ErrorCount == TotalRequests.
Assessment: PASS
Severity: LOW
Notes: Covers network failure path (connection refused).

## Finding 6
Location: tests/loadtest_test.go:122-137
Claimed Behavior: Server aborts processing when request context is already cancelled
Observed Implementation: Creates a context, cancels it immediately, sends POST request. Verifies response is not 201 Created.
Assessment: PASS
Severity: LOW
Notes: Covers context cancellation handling in server handler.

## Finding 7
Location: internal/loadtest/metrics_test.go:53-79
Claimed Behavior: Percentile ordering invariants hold (Min <= P50 <= P90 <= P95 <= P99 <= Max)
Observed Implementation: Tests with 12 unordered latency values. Verifies all ordering invariants.
Assessment: PASS
Severity: LOW
Notes: Validates that sorted percentile output maintains proper ordering.

## Finding 8
Location: tests/loadtest_test.go
Missing Coverage: No unit tests directly cover server.go's random slow-query behavior (10% chance at 25x duration)
Assessment: WARNING
Severity: MEDIUM
Notes: The server's random query degradation is never tested in isolation. The SmokeVsStress test implicitly covers it but does not specifically validate this behavior. A dedicated test for the slow-query simulation would strengthen coverage.

## Finding 9
Location: tests/loadtest_test.go
Missing Coverage: No test validates P99 specifically shows degradation under stress
Assessment: WARNING
Severity: MEDIUM
Notes: The SmokeVsStress test checks P95 degradation but not P99. The design doc claims P99 should spike significantly under stress. The demo output shows P99=1486ms (stress) vs P50=21.2ms (smoke), but this is not asserted in any test.

## Finding 10
Location: tests/loadtest_test.go
Missing Coverage: No negative test for missing/empty URL or zero VUs in loadtest.Config
Assessment: WARNING
Severity: LOW
Notes: NewRunner handles VUs <= 0 by defaulting to 1, but this behavior is not tested. An empty URL would cause all dial errors but this path is not explicitly tested as a configuration edge case.

## Finding 11
Location: tests/loadtest_test.go
Missing Coverage: No test validates that TotalRequests = SuccessCount + ErrorCount invariant holds
Assessment: WARNING
Severity: LOW
Notes: The metrics calculation ensures this invariant (total = len(latencies) + errors), but no test explicitly asserts `TotalRequests == SuccessCount + ErrorCount`. This is a basic safety invariant that should be verified.
