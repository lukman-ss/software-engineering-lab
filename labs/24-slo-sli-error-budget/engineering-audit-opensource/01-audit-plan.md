# Engineering Audit Plan

Target Lab: `labs/24-slo-sli-error-budget`
Implementation Files:
- `internal/metrics/tracker.go` — `WindowTracker` sliding time-bucketed event recorder (Record/evictStaleLocked/Summary).
- `internal/slo/evaluator.go` — `Evaluator` computing SLI ratio, error budget (remaining/consumed), and `CanDeploy` release-freeze policy.
- `internal/alerting/engine.go` — `AlertEngine` multi-window (short/long) burn-rate checker with `BurnRateRule` thresholds (14.4x fast/PAGE, 6.0x slow/TICKET).
- `cmd/demo/main.go` — executable simulation: baseline traffic, incident error-budget depletion, burn-rate alerting, and endpoint-criticality comparison.
Modules: single binary package under module `labs/24-slo-sli-error-budget` (go 1.22).

Tests: `tests/slo_test.go` (single test package covering metrics, SLO math, alerting, out-of-order timestamps, zero-traffic edge case, and concurrency).
Executable/Demo: `cmd/demo/main.go` (`go run ./cmd/demo`).
Approved Research Inputs: `research/05-report.md` (Google SRE Book Ch. 3/4/6, App. A; Prometheus alerting practices; Datadog SLO docs classified as implementation-specific). Key claims: error budget = 1 - SLO; burn-rate factors 14.4x (page) and 6.0x (ticket) are canonical Google SRE values; Datadog's percentage-formula and 1-6/6+ icon thresholds are explicitly vendor-specific.

Main Claims To Verify:
1. SLI = good / total across rolling time windows.
2. Error budget = (1 - SLO) * total; consumed by bad events; release frozen when remaining <= 0 (and total > 0).
3. Burn rate = (bad/total) / (1 - SLO); multi-window alert fires on fast (14.4x) and slow (6.0x) thresholds.
4. Endpoint-criticality bucketing supports different SLO targets (99.9% vs 95.0%).
5. Metrics tracker is thread-safe under concurrent writers (race-clean).
6. Demo output is real and reproducible.

Commands To Run:
- `go build ./...`
- `go vet ./...`
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
- `go test -coverpkg="./internal/..." ./tests/ -coverprofile` (to validate the "100% coverage" success criterion)

Primary Risks:
- Stale/mismatched recorded execution output vs current code (demo Phase 4 added post-record).
- Overclaim of "100% test coverage" on success criteria.
- Unused config fields (`BurnRateRule.LongWindow/ShortWindow/BudgetConsumedPct`, `slo.Config.LatencyThreshold`) — structural simplification.
- Defensive branches in `CalculateBurnRate`/`NewWindowTracker` not exercised by tests.
