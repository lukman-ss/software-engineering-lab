# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
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
- `research-revision/03-revision-result.md`
- `research-audit/07-verdict.md`

Main Claims To Verify:
1. SLI Calculation as quantitative ratio of good/total events.
2. Error Budget calculation and dynamic deployment freeze enforcement (`CanDeploy=false` on budget exhaustion).
3. Multi-Window Multi-Burn-Rate alerting logic evaluating short & long windows concurrently.
4. Concurrency safety under multi-threaded request recording.
5. Out-of-order event ingestion handling and sliding window eviction correctness.

Commands To Run:
- `go test -count=1 ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent bucket modification.
- Mathematical precision/rounding edge cases (zero traffic division by zero, floating point rounding).
- Alerting engine false positives/negatives during window boundary transitions.
- Docs vs code mismatch on API definitions and demo behavior.
