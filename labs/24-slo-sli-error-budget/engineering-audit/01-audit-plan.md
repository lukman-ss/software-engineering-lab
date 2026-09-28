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
- `research-audit/07-verdict.md` (APPROVED)
- `research/05-report.md`
- `research/03-evidence.md`

Main Claims To Verify:
1. SLI ratio calculation follows `good_events / total_events` format.
2. Error Budget calculation follows `(1 - SLO) * total_events` and correctly triggers deployment freeze when exhausted.
3. Multi-window multi-burn-rate alerting triggers only when both short and long window burn rates exceed configured thresholds.
4. Window tracker correctly manages sliding window time buckets, evicts stale buckets, handles out-of-order events, and maintains thread safety under concurrent access.
5. Criticality tiers allow differentiated SLO targets (e.g. 99.9% critical vs 95.0% non-critical).
6. Demo output is authentic and matches recorded execution results.
7. Documentation in `README.md` accurately describes implementation, tests, and execution commands.

Commands To Run:
- `go build ./...`
- `go test -v -count=1 ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions in sliding window time-bucket mutations or concurrent reads.
- Time-based bucket truncation bugs during boundary shifts or stale bucket eviction.
- Out-of-order event ingestion leading to unsorted buckets or slice corruption.
- Division-by-zero or NaN when total traffic is zero.
