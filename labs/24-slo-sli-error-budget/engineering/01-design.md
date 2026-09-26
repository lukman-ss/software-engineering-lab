# Engineering Design

Target Lab: `labs/24-slo-sli-error-budget`
Research Status: APPROVED

## Concept To Prove
1. **SLI Calculation as Good/Total Ratio**: Quantitative user-experience measurement using a ratio (`good_events / total_events`).
2. **Error Budget Management**: Error Budget calculated as `1 - SLO`. Consumed when requests fail or breach latency thresholds.
3. **Burn Rate Alerting**: Multi-window multi-burn-rate alerting logic detecting rapid budget consumption (e.g. 14.4x for fast burn, 6x for slow burn).
4. **Endpoint Criticality Bucketing**: Critical endpoints (e.g. Payment) configured with stricter SLOs (99.9%) compared to non-critical endpoints (e.g. Reports: 95.0%).

## Expected Behavior
- SLI accurately tracks good vs total events across rolling time windows.
- Error budget decreases when bad requests occur and triggers burn rate alerts when thresholds are breached.
- Dynamic error budget policy allows releases when budget > 0, and halts non-essential releases when budget is exhausted.

## Failure Scenario
- High failure rate or latency spike on critical endpoints rapidly consumes error budget, triggering multi-window burn rate alerts and enforcing a release freeze policy.

## Success Criteria
- 100% test coverage on core math and sliding window calculations.
- Go race detector passes cleanly without concurrency issues.
- Demonstration output shows real-time metric updates, error budget depletion, burn rate alert triggering, and release freeze enforcement.

## Architecture
- `internal/metrics`: Sliding-window event recorder (histogram latency buckets & success counts).
- `internal/slo`: SLO definition, SLI evaluation, Error Budget tracking, and Release Policy evaluator.
- `internal/alerting`: Multi-window burn-rate alert calculator.
- `cmd/demo`: Executable simulation demonstrating normal traffic, incident budget burn, alerting, and recovery.

## Components
- `Event`: Represents an HTTP request with status code, duration, endpoint, and timestamp.
- `WindowTracker`: Sliding time window tracker storing bucketed events.
- `SLOEvaluator`: Calculates SLI, remaining Error Budget, and current Burn Rate.
- `BurnRateAlertEngine`: Monitors short and long windows for fast/slow burn threshold breaches.

## Test Strategy
- Unit tests for sliding window bucket aggregation, edge cases (zero traffic, 100% errors, rolling window expiry).
- Concurrency test with parallel requests updating metrics concurrently (`go test -race`).
- Integration tests verifying multi-window burn rate alert triggering.

## Execution Plan
1. Create `go.mod` for `labs/24-slo-sli-error-budget`.
2. Implement core packages in `internal/metrics`, `internal/slo`, and `internal/alerting`.
3. Implement unit and concurrency tests in `tests/`.
4. Implement runnable interactive simulation in `cmd/demo/main.go`.
5. Run tests, race detector, and demo.
6. Record results in `engineering/02-implementation-notes.md` and `engineering/03-execution-result.md`.

## Implementation Decisions
- Standard library only (`time`, `sync`, `testing`, `fmt`).
- In-memory ring/time-bucketed window tracking to simulate 30-day/rolling windows in compressed real-time (e.g. 1-second = 1-hour scale for demo).
