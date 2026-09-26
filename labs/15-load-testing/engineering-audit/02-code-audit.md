# Engineering Code Audit

Target Lab: labs/15-load-testing

## Finding 1

Location: `internal/loadtest/runner.go:68-80`
Claimed Behavior: Load test runner executes HTTP requests concurrently per VU without shared lock contention.
Observed Implementation: Each VU goroutine appends results to a thread-local slice (`results[vuID]`). Slices are aggregated once on completion.
Assessment: PASS
Severity: LOW
Notes: Clean concurrent architecture that prevents mutex contention from distorting latency measurements.

## Finding 2

Location: `internal/loadtest/metrics.go:65-71`
Claimed Behavior: Accurate percentile computation.
Observed Implementation: Uses standard 0-indexed integer rounding: `idx := int(float64(len(sorted)-1) * (pct / 100.0))`.
Assessment: PASS
Severity: LOW
Notes: Integer truncating on zero-based indices satisfies testing requirements for sample sizes < 10,000 without requiring heavy external dependencies.

## Finding 3

Location: `internal/server/server.go:50-68`
Claimed Behavior: Server simulates database connection pool saturation causing tail latency spikes that are masked by average latency.
Observed Implementation: Uses a buffered channel semaphore and fixed query time (`20ms`). In a closed-loop VU load model, once steady state is reached, the wait time is uniform across all requests (`~200ms`), causing average latency and P95 latency to be nearly identical (~200ms vs ~211ms).
Assessment: WARNING
Severity: HIGH
Notes: Fails to demonstrate the claimed "masking effect" where average conceals tail latency spikes, because there is no tail distribution (all requests queue equally due to fixed durations in a closed workload model).
