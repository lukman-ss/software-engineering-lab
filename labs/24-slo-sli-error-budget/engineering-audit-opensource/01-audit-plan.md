# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Audit Date: 2026-09-27
Output Dir: labs/24-slo-sli-error-budget/engineering-audit-opensource/

## Implementation Files

- internal/metrics/tracker.go (WindowTracker: sliding-window bucketed event tracker)
- internal/slo/evaluator.go (Evaluator: SLI ratio, error budget, CanDeploy freeze)
- internal/alerting/engine.go (AlertEngine: burn-rate calc, multi-window AND-gate check)
- cmd/demo/main.go (4-phase demo: baseline, incident, alerts, criticality comparison)
- tests/slo_test.go (6 tests)
- go.mod (module labs/24-slo-sli-error-budget, go 1.22)

## Tests

- TestMetricsWindowTracker, TestSLOEvaluator, TestAlertEngineBurnRate,
  TestOutOfOrderTimestamps, TestEvaluatorZeroTraffic, TestConcurrencyMetrics

## Executable/Demo

- cmd/demo: baseline 1000 good reqs, incident 90 good + 10 bad, burn-rate check,
  Payment 99.9% vs Reports 95% comparison

## Approved Research Inputs

Per pipeline override: research/content NOT audited in this stage. Skipped.

## Main Claims To Verify

1. SLI = good/total ratio over sliding window
2. Error budget = (1 - SLO) * total, consumed per bad event
3. Release freeze: CanDeploy=false when budget exhausted
4. Multi-window burn-rate alerting (fast 14.4x / slow 6x) with transient-spike rejection
5. Stricter SLO for critical endpoints vs lax SLO
6. Thread-safe concurrent recording; zero-traffic and eviction edge cases

## Commands To Run

```bash
cd labs/24-slo-sli-error-budget
go vet ./...
go build ./...
go test -count=1 -v ./...
go test -count=1 -race ./...
go run ./cmd/demo
```

## Primary Risks

- Burn-rate rule window fields possibly decorative (engine uses injected trackers)
- Freeze boundary (remaining == 0) float-sensitive
- No minimum-sample guard -> tiny traffic can page
- Dead Config.LatencyThreshold field
