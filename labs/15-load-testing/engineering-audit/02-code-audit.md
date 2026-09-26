# Code Audit

Target Lab: labs/15-load-testing

## Finding 1

Location: internal/server/server.go:66-71
Claimed Behavior: Bounded database connection capacity simulating semaphore pool limit.
Observed Implementation: Channel semaphore `s.semaphore = make(chan struct{}, cfg.MaxDBConnections)` buffers up to `MaxDBConnections`. Acquisition uses non-blocking select against `r.Context().Done()`:
```go
select {
case s.semaphore <- struct{}{}:
case <-r.Context().Done():
    return
}
defer func() { <-s.semaphore }()
```
Assessment: PASS
Severity: LOW
Notes: Correctly releases token via deferred receive and respects request context cancellation.

## Finding 2

Location: internal/loadtest/runner.go:48-53, 107-112
Claimed Behavior: Lock-free per-VU data collection during load test execution.
Observed Implementation: Each VU worker routine appends to local slice `lats` and records its own error count, then places its struct into pre-allocated slice `results[vuID]` upon `ctx.Done()`. Aggregation happens sequentially after `wg.Wait()`.
Assessment: PASS
Severity: LOW
Notes: Eliminates mutex contention overhead during benchmarking and avoids race conditions.

## Finding 3

Location: internal/loadtest/metrics.go:67-73
Claimed Behavior: Accurate percentile computation on sorted latency slices.
Observed Implementation:
```go
func percentile(sorted []time.Duration, pct float64) time.Duration {
    if len(sorted) == 0 {
        return 0
    }
    idx := int(float64(len(sorted)-1) * (pct / 100.0))
    return sorted[idx]
}
```
Assessment: PASS
Severity: LOW
Notes: `len(sorted)-1` scaling prevents index out of range errors for boundary values `pct=0` and `pct=100`, returning valid durations on single-element slices.

## Finding 4

Location: internal/loadtest/runner.go:72-76, 85-89
Claimed Behavior: Only legitimate runtime errors tracked; context cancellation at test end not counted as failed requests.
Observed Implementation: Checks `if ctx.Err() == nil` before incrementing `errs++` on request creation or execution error.
Assessment: PASS
Severity: LOW
Notes: Invariant `TotalRequests == SuccessCount + ErrorCount` holds across lifecycle.
