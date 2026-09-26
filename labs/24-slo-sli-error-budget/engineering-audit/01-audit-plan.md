# Engineering Audit Plan

Target Lab: `labs/24-slo-sli-error-budget`
Implementation Files:
- `internal/metrics/tracker.go`
- `internal/slo/evaluator.go`
- `internal/alerting/engine.go`
- `cmd/demo/main.go`

Tests:
- `tests/slo_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research-revision/03-revision-result.md`
- `research-audit/07-verdict.md`

Main Claims To Verify:
1. SLI calculation implements good/total event ratio tracking over sliding time windows (`internal/metrics` + `internal/slo`).
2. Error budget dynamically calculated as `(1 - SLO) * total_events - bad_events` and enforces freeze policy when exhausted.
3. Multi-window multi-burn-rate alerting triggers only when both short and long window burn rates exceed rule thresholds.
4. Concurrency safety: metrics aggregation and sliding window tracker are thread-safe under concurrent recording.
5. Out-of-order event insertion and timestamp eviction function correctly.
6. Zero-traffic edge cases do not trigger divide-by-zero panics and default to valid state.
7. Endpoint criticality comparison (e.g. 99.9% vs 95.0%) behaves according to configured thresholds.

Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent slice mutations in `WindowTracker`.
- Divide-by-zero panics in SLI or Burn Rate calculations on empty windows.
- Out-of-order event eviction bugs causing stale bucket retention or slice index corruption.
- Mismatch between demo output recorded in docs and live demo execution.
