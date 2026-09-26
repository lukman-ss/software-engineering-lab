# Code Audit

## Finding 1

Location: `internal/loadtest/runner.go:48-53, 68`
Claimed Behavior: Concurrently collect VU metrics without lock contention or race conditions.
Observed Implementation: Each VU goroutine appends locally to its own slice (`lats`) and writes to an indexed position in `results[vuID]` upon exit. WaitGroup ensures completion before combining.
Assessment: PASS
Severity: LOW
Notes: Clean concurrent architecture; avoids mutex contention under load.

## Finding 2

Location: `internal/loadtest/metrics.go:45-64`
Claimed Behavior: Accurately compute percentiles (P50, P90, P95, P99) and summary statistics.
Observed Implementation: Makes a defensive copy of latencies slice, sorts it via `sort.Slice`, and calculates percentiles by mapping `(len-1) * pct / 100`. Returns zero metrics gracefully on empty input.
Assessment: PASS
Severity: LOW
Notes: Correct standard percentile calculation and defensive against mutation of input slices.

## Finding 3

Location: `internal/server/server.go:61-76`
Claimed Behavior: Simulate DB connection saturation with context cancellation support.
Observed Implementation: Uses buffered channel semaphore to bound concurrent DB slots. Accurately aborts on `r.Context().Done()`. Defers slot release.
Assessment: PASS
Severity: LOW
Notes: Properly handles context cancellation and limits concurrency to `MaxDBConnections`.
