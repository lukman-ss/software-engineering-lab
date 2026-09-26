# Implementation Notes

## Files Added
- `go.mod`: Go module specification (Go 1.22+).
- `internal/metrics/tracker.go`: Sliding window time-bucketed metric aggregation tracker.
- `internal/slo/evaluator.go`: SLI ratio calculator and Error Budget policy manager.
- `internal/alerting/engine.go`: Multi-window burn rate alert evaluator.
- `cmd/demo/main.go`: Interactive executable demonstration.
- `tests/slo_test.go`: Unit tests and concurrent race detection tests.
- `engineering/01-design.md`: Architecture and design document.
- `engineering/02-implementation-notes.md`: Implementation documentation.
- `engineering/03-execution-result.md`: Recorded build, test, race, and demo execution outputs.
- `README.md`: High-level overview and instructions.

## Core Design Decisions
1. **Ratio-Based SLI**: Evaluated as `good_events / total_events` for consistent metric aggregation.
2. **Error Budget Calculation**: Budget defined as `(1.0 - target_slo) * total_events`. Consumed on bad events or latency threshold breaches.
3. **Multi-Window Burn-Rate Alerting**: Multi-window check evaluating short and long windows concurrently to prevent alert fatigue and false positives.

## Implementation-Specific Choices
- In-memory time-bucketed ring buffer (`WindowTracker`) used to aggregate metrics without external dependencies like Prometheus or Redis.
- standard library only (`time`, `sync`, `testing`, `math`, `fmt`).

## Known Limitations
- Metrics are held entirely in-memory and will reset if the process restarts.

## Trade-offs
- Bucket resolution vs memory footprint: Sub-second bucket sizes provide higher granularity but consume slightly more memory. Default bucket size of 1-10 seconds chosen for optimal real-time efficiency.

## What Is Demonstrated
- Accurate SLI evaluation under normal and incident traffic.
- Error budget depletion leading to release freeze policy enforcement (`CanDeploy = false`).
- Burn rate calculation and alert triggering upon rapid error budget consumption.

## What Is Not Demonstrated
- Integration with external TSDBs (Prometheus/Datadog) or external alerting notification platforms (PagerDuty/Slack).
