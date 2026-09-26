# Code Audit

## Finding 1

Location: `internal/server/server.go`
Claimed Behavior: Simulates a constrained database connection pool using a semaphore, causing queuing latency.
Observed Implementation: Uses a buffered channel `semaphore <- struct{}{}` sized to `MaxDBConnections`. Implements 10% penalty (`dur * 25`) if `activeReq` > `MaxDBConnections`.
Assessment: PASS
Severity: LOW
Notes: Implementation correctly mocks backpressure and resource exhaustion.

## Finding 2

Location: `internal/loadtest/runner.go`
Claimed Behavior: Generates load safely without lock contention.
Observed Implementation: Spawns exactly `VUs` goroutines. Each VU stores results in a pre-allocated slice `results[vuID]`. Uses `sync.WaitGroup` to wait for completion before aggregating.
Assessment: PASS
Severity: LOW
Notes: Clean concurrent design. Prevents client-side bottlenecks. Custom `http.Transport` correctly configured to allow up to 1000 connections.

## Finding 3

Location: `internal/loadtest/metrics.go`
Claimed Behavior: Accurately calculates Min, Max, Avg, P50, P90, P95, P99.
Observed Implementation: Sorts a copy of latencies slice. Index calculation `int(float64(len(sorted)-1) * (pct / 100.0))` is standard and correct for percentile sampling. Average handles `0` requests safely.
Assessment: PASS
Severity: LOW
Notes: Explicitly calls out memory scaling trade-off (linear exact sort). Perfect for the lab constraints.
