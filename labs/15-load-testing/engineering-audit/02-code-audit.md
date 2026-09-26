# Code Audit

Target Lab: labs/15-load-testing

## Finding 1

Location: `internal/server/server.go:61-81`
Claimed Behavior: Server enforces max concurrent DB operations via channel semaphore and respects request context cancellations during queueing and query execution.
Observed Implementation: Buffered channel semaphore of size `cfg.MaxDBConnections` properly guards critical section; uses `select` on `r.Context().Done()` during acquire and query sleep; timers properly stopped via `defer t.Stop()`.
Assessment: PASS
Severity: LOW
Notes: Correct context lifecycle management preventing goroutine leaks when client cancels.

## Finding 2

Location: `internal/loadtest/runner.go:48-114`
Claimed Behavior: Concurrency runner allocates separate per-VU metric slices to avoid lock contention during high throughput runs, and aggregates latencies and errors post-run.
Observed Implementation: Each worker goroutine writes exclusively to its indexed `results[vuID]` struct without shared mutex locks, synchronized via `sync.WaitGroup`. Latency slices aggregated after `wg.Wait()`.
Assessment: PASS
Severity: LOW
Notes: Clean concurrent architecture with zero race conditions detected under `go test -race`.

## Finding 3

Location: `internal/loadtest/metrics.go:44-73`
Claimed Behavior: Accurate calculation of percentiles (P50, P90, P95, P99), min, max, average, and RPS.
Observed Implementation: Uses sorted latency copies, correct 0-indexed linear rank indexing, defensive checks for empty input slices, and accurate integer duration arithmetic.
Assessment: PASS
Severity: LOW
Notes: Formula invariant property verified by unit tests.

## Finding 4

Location: `internal/loadtest/runner.go:30-34`
Claimed Behavior: HTTP transport configured to avoid client-side connection pooling bottlenecks hiding server saturation.
Observed Implementation: `MaxIdleConns` and `MaxIdleConnsPerHost` set to 1000 with a 5s client timeout.
Assessment: PASS
Severity: LOW
Notes: Sound load test client design.
