# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Implementation Files:
- `internal/metrics/tracker.go`
- `internal/slo/evaluator.go`
- `internal/alerting/engine.go`
Tests:
- `tests/slo_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`
Main Claims To Verify:
1. SLI calculation implements ratio of good events over total valid events (`good / total`).
2. Error Budget dynamically computed as `(1 - TargetSLO) * total - bad`.
3. Release freeze / policy enforcement toggles `CanDeploy` when error budget <= 0.
4. Multi-window multi-burn-rate alerting triggers when both short and long rolling window burn rates exceed rule thresholds.
5. In-memory time-bucketed sliding window correctly handles concurrent writes, event ordering, and stale bucket eviction.
6. Endpoint criticality differentiation allows distinct SLO targets (e.g. 99.9% vs 95.0%).
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Concurrency race conditions or deadlocks during concurrent metric ingestion and eviction.
- Sliding window eviction inaccuracies with non-monotonic/out-of-order timestamps.
- Floating point precision truncation during SLI/Error budget arithmetic.
