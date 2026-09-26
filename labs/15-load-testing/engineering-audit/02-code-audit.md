# Code Audit

Target Lab: labs/15-load-testing
Audit Scope: `internal/server/server.go`, `internal/loadtest/runner.go`, `internal/loadtest/metrics.go`, `cmd/demo/main.go`

---

## Finding 1: Concurrency Control in Mock Server Connection Pool

Location: `internal/server/server.go:18,46,67,71`
Claimed Behavior: Bounded database connection pool capacity using a semaphore to simulate resource saturation under concurrent load.
Observed Implementation:
- Server initializes `semaphore: make(chan struct{}, cfg.MaxDBConnections)`.
- Handlers push to `s.semaphore` before simulated query and pop on completion (`defer func() { <-s.semaphore }()`).
- Context cancellation during queue waiting (`select { case s.semaphore <- struct{}{}: ... case <-r.Context().Done(): return }`) safely releases slots.
Assessment: PASS
Severity: LOW
Notes: Correctly models resource contention without leaking channel tokens on request abortion.

---

## Finding 2: Per-VU Metric Buffers and Race Prevention

Location: `internal/loadtest/runner.go:48-53,60-101,107-114`
Claimed Behavior: Virtual users run concurrently without mutex lock contention or data races during latency recording.
Observed Implementation:
- `results := make([]vuResult, r.cfg.VUs)` preallocates distinct per-VU result slots.
- Each goroutine (`vuID`) writes exclusively to its own slice/counters in `results[vuID]`.
- Aggregation is performed sequentially after `wg.Wait()`.
Assessment: PASS
Severity: LOW
Notes: Clean lock-free pattern for load generation. Confirmed race-free via `-race`.

---

## Finding 3: Percentile Calculation Correctness and Monotonicity

Location: `internal/loadtest/metrics.go:45-73`
Claimed Behavior: Accurately compute Min, Max, Average, P50, P90, P95, and P99 percentiles from collected latency durations.
Observed Implementation:
- Copies slice to avoid mutating input, sorts with standard `sort.Slice`.
- Percentile index computed via `int(float64(len(sorted)-1) * (pct / 100.0))`.
- Empty slice, single sample, and monotonic invariant checks are handled robustly.
Assessment: PASS
Severity: LOW
Notes: Index arithmetic maps correctly across 0 to len-1 bounds.

---

## Finding 4: HTTP Client Connection Pooling and Transport Tuning

Location: `internal/loadtest/runner.go:30-41`
Claimed Behavior: Avoid load runner client-side connection pooling bottlenecks during high-VU load tests.
Observed Implementation:
- Sets `MaxIdleConns: 1000` and `MaxIdleConnsPerHost: 1000` in `http.Transport`.
- Sets client timeout to 5 seconds to prevent unbounded hanging goroutines.
- Responses read and closed reliably with `io.Copy(io.Discard, resp.Body)` and `resp.Body.Close()`.
Assessment: PASS
Severity: LOW
Notes: Correctly isolates server-side bottleneck from client transport constraints.

---

## Finding 5: Context Termination and Timeout Handling in VU Loops

Location: `internal/loadtest/runner.go:55,66-70,85-88`
Claimed Behavior: Load test runs for the specified duration and terminates all workers gracefully without falsely reporting context expiration as server errors.
Observed Implementation:
- `context.WithTimeout(ctx, r.cfg.Duration)` controls loop lifecycle.
- Workers inspect `ctx.Err()` to distinguish intentional test termination from genuine HTTP/network failures.
Assessment: PASS
Severity: LOW
Notes: Prevents false positive error rate inflation when load run finishes.

---

## Finding 6: Demo Scenario Structure and Realism

Location: `cmd/demo/main.go:16-59`
Claimed Behavior: Live execution running Smoke (2 VUs) vs. Stress (50 VUs) scenarios against a 5-connection constrained server.
Observed Implementation:
- Starts actual `httptest.NewServer`, runs standard HTTP POST requests with realistic JSON payloads.
- Outputs clean tabular comparison with actual measured RPS, Average, P50, P95, and P99 metrics.
Assessment: PASS
Severity: LOW
Notes: No hardcoded output or synthetic results; output is computed directly from live test execution.
