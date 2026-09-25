# Code Audit

## Finding 1
Location: `internal/loadtest/runner.go:47-97`
Claimed Behavior: Concurrency orchestrator generating HTTP traffic without race conditions.
Observed Implementation: Uses per-VU slices `results[vuID]` to record latencies and errors. Aggregates post-execution after `wg.Wait()`. Uses custom transport to disable internal connection limits.
Assessment: PASS
Severity: LOW
Notes: Safe and clean concurrent execution. `ctx.Err() == nil` avoids false positive errors on timeout.

## Finding 2
Location: `internal/server/server.go:59-62`
Claimed Behavior: Server limits concurrency to simulate DB connection pool exhaustion.
Observed Implementation: Uses a buffered channel (`semaphore`) to limit active processing. Sleep simulates query duration.
Assessment: PASS
Severity: LOW
Notes: Effectively models resource saturation.

## Finding 3
Location: `internal/loadtest/metrics.go:44-62`
Claimed Behavior: Calculates min, max, avg, P50, P90, P95, P99 accurately.
Observed Implementation: Copies and sorts slice, computes sum and averages, and uses rank-based index mapping `idx := int(float64(len(sorted)-1) * (pct / 100.0))` for percentiles.
Assessment: PASS
Severity: LOW
Notes: Sorting slice is memory-heavy for massive scale but perfectly acceptable for lab constraints as documented.
