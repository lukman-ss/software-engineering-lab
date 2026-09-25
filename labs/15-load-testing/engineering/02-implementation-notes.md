# Implementation Notes

## Files Added
- `internal/server/server.go`: HTTP handler simulating an endpoint with a constrained connection pool.
- `internal/loadtest/metrics.go`: Thread-safe latency collector and percentile calculator.
- `internal/loadtest/runner.go`: Concurrency orchestrator generating HTTP traffic and aggregating results.
- `tests/loadtest_test.go`: Integration tests confirming metrics tracking and resource exhaustion behaviors.
- `cmd/demo/main.go`: Executable runner proving baseline vs. saturated load scenarios.

## Core Design Decisions
- Hand-rolled percentile calculation over standard arrays rather than adding `HdrHistogram` dependency. Exact sorting is feasible for test durations.
- Semaphore pattern (buffered channel) used in the server handler to simulate database connection pool bottlenecks, accurately mimicking tail latency growth when saturated.
- Per-goroutine slices in load generator to avoid mutex contention, aggregating once on completion.

## Implementation-Specific Choices
- Wait durations in server mock are fixed (`20ms`), making the tail latency strictly a function of queuing time when VUs exceed max DB connections.
- Load generator executes loops with `time.Since` for durations rather than pre-generating requests to reduce memory bloat.

## Known Limitations
- Percentile algorithm is exact-sort, memory footprint scales linearly with request count; unsuitable for hour-long multi-million RPS benchmarks, but perfect for a 2-second lab test.
- Simulated external network latency is omitted to isolate the connection pool bottleneck.

## Trade-offs
- Used standard `net/http` client which carries its own connection pooling limits; overridden using a custom `http.Transport` to ensure the load generator does not become the bottleneck.

## What Is Demonstrated
- Measuring P50, P95, and P99 percentiles manually.
- Smoke load executing linearly (fast average, fast P95).
- Stress load queuing behind a saturation point (connection pool limit), forcing P95 to severely degrade.

## What Is Not Demonstrated
- Distributed load generation (multi-node).
- Real database lock contention or CPU saturation (mocked via semaphore and time.Sleep).
