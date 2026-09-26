# Test Audit

Target Lab: labs/15-load-testing

## Test Coverage Overview

- `internal/loadtest/metrics_test.go`:
  - `TestCalculateMetrics`: Verifies known array distribution [1ms..100ms], checks exact min, max, avg, P50 (50ms), P95 (95ms), P99 (99ms).
  - `TestCalculateMetrics_Empty`: Verifies behavior on empty slice input (zeros returned, no panics).
  - `TestCalculateMetrics_Invariants`: Verifies monotonicity invariant `Min <= P50 <= P90 <= P95 <= P99 <= Max` on noisy unordered distributions.
- `tests/loadtest_test.go`:
  - `TestLoadTest_SmokeVsStress`: Spins up httptest server with 2 DB connections. Verifies that 10 VU stress scenario demonstrates significant P95 inflation over smoke scenario (1 VU) and tail latency exceeding average.
  - `TestLoadTest_ErrorCount`: Verifies HTTP 500 error recording without treating errors as successful latencies.
  - `TestServer_MethodNotAllowed`: Verifies HTTP 405 status code handling on non-POST methods.
  - `TestLoadTest_DialError`: Verifies connection failure handling (dial error recording).
  - `TestServer_ContextCanceled`: Verifies server immediately aborts when request context is canceled without returning 201 Created.

## Execution Verification

1. `go test -v ./...`
   Result: PASS (all tests pass across both packages).

2. `go test -race -v ./...`
   Result: PASS (clean race detector report, no data races).

3. `go run ./cmd/demo`
   Result: PASS (executes smoke and stress runs, displaying real, consistent metric shifts from ~21ms P95 to ~1.08s P95 under DB bottlenecking).

## Assessment

Coverage: Complete across happy path, edge cases, error tracking, context cancellations, and concurrency safety.
Assessment: PASS
