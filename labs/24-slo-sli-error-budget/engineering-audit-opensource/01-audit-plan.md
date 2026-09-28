# Engineering Audit Plan

Target Lab: `labs/24-slo-sli-error-budget`
Audit Date: 2026-09-28
Audit Output Dir: `labs/24-slo-sli-error-budget/engineering-audit-opensource/`
Scope (pipeline override): implementation and tests only. No research/content audit. No code modification.

## Implementation Files

- `internal/metrics/tracker.go` (128 lines) — sliding-window time-bucketed event tracker
- `internal/slo/evaluator.go` (71 lines) — SLI ratio, error budget, release-freeze policy
- `internal/alerting/engine.go` (89 lines) — burn-rate calculation, multi-window alert check
- `cmd/demo/main.go` (151 lines) — 4-phase executable demo
- `go.mod` — module `labs/24-slo-sli-error-budget`, go 1.22, stdlib only

## Tests

- `tests/slo_test.go` (236 lines, 6 tests): TestMetricsWindowTracker, TestSLOEvaluator,
  TestAlertEngineBurnRate, TestOutOfOrderTimestamps, TestEvaluatorZeroTraffic, TestConcurrencyMetrics

## Executable/Demo

- `cmd/demo`: baseline traffic → incident burn → burn-rate alerts → endpoint criticality comparison

## Approved Research Inputs

Not in scope per pipeline override. Research verdict (`research-audit/07-verdict.md`) is APPROVED;
used only as background, not audited here.

## Main Claims To Verify

1. SLI computed as good/total ratio over a sliding window (`internal/slo`, `internal/metrics`).
2. Error budget = `(1 - target) * total - bad`; freeze releases (`CanDeploy=false`) when exhausted.
3. Multi-window burn-rate alerting (fast 14.4x / slow 6x) requiring both windows over threshold.
4. Endpoint criticality: stricter SLO (99.9%) vs lenient SLO (95.0%).
5. Thread-safe concurrent recording (race detector clean).
6. Demo output is real; README commands work; execution-result record is accurate.

## Commands To Run (from the lab dir — module-scoped, `./...` matches nothing outside it)

```bash
cd /Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/24-slo-sli-error-budget
go build ./...
go test -count=1 -v ./...
go test -count=1 -race ./...
go run ./cmd/demo
```

## Primary Risks

- `BurnRateRule.LongWindow/ShortWindow/BudgetConsumedPct` may be dead fields (engine uses two fixed trackers).
- Demo short/long windows may contain identical events (no real window separation in demo).
- Design doc claims (100% coverage, histogram buckets, recovery phase) may overclaim implementation.
- Concurrency test assertion may be weaker than exact counts.
