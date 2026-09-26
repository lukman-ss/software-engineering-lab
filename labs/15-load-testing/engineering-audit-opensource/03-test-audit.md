# Test Audit

## Test Suite Coverage

The lab contains two test packages:
- `internal/loadtest` unit tests: metrics, invariants, RPS, empty/single cases
- `tests` integration tests: actual HTTP load against the server

### Loadtest Unit Tests (metrics_test.go)
- `TestCalculateMetrics`: correctness for 100-sample ascending series (Min/Max/Avg/P50/P95/P99)
- `TestCalculateMetrics_Empty`: zero-sample edge case
- `TestCalculateMetrics_SingleSample`: one-sample invariants (min=max=percentiles)
- `TestCalculateMetrics_RPS`: requests-per-second arithmetic
- `TestCalculateMetrics_Invariants`: monotonic-min-to-max and percentile ordering (Min ≤ P50 ≤ P90 ≤ P95 ≤ P99 ≤ Max)
These five tests cover the deterministic statistical core. They cover happy path and edge cases (empty, single) but do NOT inject out-of-order samples to test sorting stability; they trust Go's `sort.Slice`. Acceptable for lab scale.

### Integration Tests (loadtest_test.go)
There are seven integration tests exercising actual concurrency and the server:
1. `TestLoadTest_SmokeVsStress`: 1 VU (smoke) vs 10 VU (stress) against a 2-slot DB server; asserts stress P95 > smoke P95 and stress P95 > stress Avg (tail spikes under load).
2. `TestLoadTest_ErrorCount`: mock 500 server; asserts all requests error (Success=0, Error=Total).
3. `TestServer_MethodNotAllowed`: GET on POST-only endpoint → 405.
4. `TestLoadTest_DialError`: connect to dead port (127.0.0.1:1); asserts all attempts fail (Success=0, Error=Total).
5. `TestServer_ContextCanceled`: pre-canceled context; asserts no 201 response.
6. `TestServer_MaxDBConnectionsBound`: 15 concurrent VU vs 3-slot server; uses polling + atomic to confirm active connections never exceed max (isolation of semaphore).
7. `TestLoadTest_SuccessAndErrorInvariant`: ensures Success + Error = Total for all runs.

### Execution Results
All tests pass (`go test -v ./...`); race detector clean (`go test -race ./...`).
```
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest	(cached)
?   	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/server	[no test files]
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/tests	(cached)
```

### Specific Coverage Checks
- Happy path: smoke test (low VUs, no queue) passes via `TestLoadTest_SmokeVsStress` (stress half) + demo.
- Failure path: error-count test, dial-error test, method-not-allowed test, context-canceled test.
- Edge cases: empty/single latency slices via unit tests.
- Transitions: smoke-vs-stress directly contrasts baseline (no queuing) vs saturation.
- Recovery/rollback: not applicable — no persistent state to recover.
- Concurrency: `TestServer_MaxDBConnectionsBound` confirms semaphore bound under 15 VUs; `-race` clean.
- Negative cases: covered by error-path tests.

### Quality and Gaps
- No flaky tests observed (multiple runs stable). Timing-dependent assertions (`stressP95 > smokeP95`) are robust due to order-of-magnitude gap (21ms vs 700ms+). The gap is large enough to survive CI jitter.
- Missing property-based tests for the percentile function (e.g., quicksort-style invariant over random inputs) — but the existing unit tests (empty/single/100-inv/arbitrary 12-sample invariants) are sufficient for the lab's scope.
- No explicit test of the `Context` cancellation inside the load generator loop (runner.go's `select { case <-ctx.Done(): ... }`), though `TestServer_ContextCanceled` tests server-side cancellation and the runner does check `ctx.Err()` on transport errors.
- Overall, the test suite validates the claimed behavior: smoke vs stress latency divergence, error counting, semaphore bounds, and metric correctness. No HIGH or CRITICAL gaps.