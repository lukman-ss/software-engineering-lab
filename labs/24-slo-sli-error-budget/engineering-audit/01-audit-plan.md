# Engineering Audit Plan

Target Lab: `labs/24-slo-sli-error-budget`
Implementation Files:
- `internal/metrics/tracker.go`
- `internal/slo/evaluator.go`
- `internal/alerting/engine.go`
- `cmd/demo/main.go`
- `go.mod`

Tests:
- `tests/slo_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`

Main Claims To Verify:
1. SLI is calculated strictly as `good_events / total_events` across a time window.
2. Error budget is calculated as `(1 - SLO) * total_events`, budget consumed equals `bad_events`, and negative remaining budget enforces `CanDeploy = false`.
3. Out-of-order events are inserted and aggregated into correct time buckets rather than appended or dropped.
4. Window eviction drops buckets older than `windowSize` relative to query/record timestamps.
5. Multi-window multi-burn-rate alerting requires both short window and long window burn rates to exceed the burn rate threshold before triggering.
6. Endpoint criticality differentiation correctly configures different SLO targets (99.9% vs 95.0%) with distinct error tolerance thresholds.
7. Concurrency safety: metrics aggregation via `WindowTracker` is safe under concurrent reader/writer goroutines without data races.
8. Demo executes realistically without mocked or hardcoded static output.

Commands To Run:
- `go test -count=1 -v ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Slice re-allocation/reslicing bugs in `WindowTracker.Record` during out-of-order insertion.
- Eviction bounds calculation during sliding window evaluation.
- Division by zero in SLI or BurnRate calculations under zero traffic.
- Discrepancies between demo output and README / engineering claims.
