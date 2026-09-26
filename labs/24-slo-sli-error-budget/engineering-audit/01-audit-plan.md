# Engineering Audit Plan

Target Lab: `labs/24-slo-sli-error-budget`
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
- `research/03-evidence.md`
- `research-audit/07-verdict.md` (APPROVED)
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. SLI calculation implements Google SRE ratio `good_events / total_events`.
2. Error budget dynamically calculated as `(1 - SLO) * total_events` and consumed by failures/latency breaches.
3. Release freeze policy (`CanDeploy = false`) enforced when remaining error budget <= 0.
4. Multi-window multi-burn-rate alerting triggers only when both short and long rolling windows breach burn rate factor threshold.
5. In-memory sliding window bucket tracker handles out-of-order events, time eviction, and concurrent access safely without data races.
6. Endpoint criticality tiers allow distinct SLO/error budget tolerances (Payment 99.9% vs Reports 95.0%).
7. Go test suite compiles and runs cleanly with zero failures and passes `-race`.
8. Interactive demo runs and produces genuine output without simulated mocks or hardcoded fake results.

Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race -v -count=1 ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent event ingestion and bucket summary eviction.
- Out-of-order timestamp handling in `metrics.WindowTracker`.
- Zero-traffic edge case (division by zero in SLI or burn rate calculations).
- Discrepancy between Google SRE multi-window alerting formulation and implementation logic.
