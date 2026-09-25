# Engineering Design

Target Lab: labs/15-load-testing
Research Status: APPROVED

## Concept To Prove
- Load testing reveals boundaries, saturation, and degradation patterns.
- Average response time conceals tail latency spikes; percentiles (P95, P99) are necessary to uncover degradation.
- Incremental test stages (smoke vs stress) differentiate baseline performance from resource exhaustion.
- Downstream resource saturation (e.g., database connection pool exhaustion) causes non-linear latency degradation for tail requests.

## Expected Behavior
- **Smoke Load (low VUs)**: All requests process within normal latency limits. Average and P95 are close. Error rate is 0%.
- **Stress Load (high VUs)**: Concurrent requests exceed server resource capacity (connection pool limit). Tail requests queue up. P95 and P99 latency spikes significantly, while average latency degrades less severely, proving the masking effect of averages.

## Failure Scenario
- Under excessive concurrent load, queuing behind a constrained resource (connection pool) causes high latency and timeouts for the 95th percentile.

## Success Criteria
- Automated benchmarks and tests execute without external dependencies.
- Load generator computes Min, Max, Average, P50, P90, P95, and P99 latencies accurately.
- Demonstration clearly contrasts Smoke test metrics against Stress test metrics.
- All tests pass with zero race conditions (`go test -race ./...`).

## Architecture
- `internal/server`: HTTP server implementing `POST /booking` with simulated resource constraints (bounded worker/connection pool).
- `internal/loadtest`: Self-contained load test harness that spawns VUs, executes requests, tracks latencies, and computes statistical percentiles.
- `cmd/demo`: Executable running a smoke test followed by a stress test, printing comparative metric tables.

## Components
1. `BookingServer`: HTTP handler with a configurable concurrency limit (semaphore) simulating a database connection pool.
2. `LoadTester`: Concurrency orchestrator generating HTTP traffic with specified virtual users (VUs) and iterations.
3. `MetricsAggregator`: Thread-safe latency collector sorting durations to derive accurate percentiles.

## Test Strategy
- Unit tests: Verify statistical calculations (P50, P95, P99, Avg).
- Integration tests: Run actual HTTP load tests against the mock server to confirm metric collection and error handling.
- Concurrency test: Run race detector to verify safe concurrent updates in load generation.

## Execution Plan
1. Initialize Go module with standard library.
2. Implement metrics aggregator and percentile logic.
3. Implement constrained server.
4. Implement load test runner.
5. Create comprehensive test suite.
6. Verify via `go test -race ./...`.
7. Implement runnable demo.

## Implementation Decisions
- Standard library `net/http` and `sync` used exclusively; no third-party dependencies required.
- Ponytail: Exact sorting used for percentiles instead of HdrHistogram; suitable for the test scale (under 10,000 samples).
