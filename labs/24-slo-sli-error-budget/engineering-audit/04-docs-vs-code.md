# Documentation vs Code Audit

Target Lab: `labs/24-slo-sli-error-budget`

## Comparison Matrix

| Aspect | README / Design Claim | Implementation / Demo Output | Match Status |
|---|---|---|---|
| Structure | `internal/metrics`, `internal/slo`, `internal/alerting`, `cmd/demo`, `tests/` | Exact directory and package structure exists | MATCH |
| Metrics Engine | Sliding-window time-bucketed event tracker measuring good vs total events | `metrics.WindowTracker` with configurable window/bucket size and predicate `isGoodEvent` | MATCH |
| SLO Evaluator | Evaluates SLI ratio, remaining error budget, and release freeze policy | `slo.Evaluator` returns `Status` with SLI, Error Budget, and `CanDeploy` boolean | MATCH |
| Burn Rate Alerting | Multi-window burn-rate alert calculator evaluating fast and slow budget burn rates | `alerting.AlertEngine` evaluating short and long window burn rates against factor thresholds | MATCH |
| Demo Execution | Run `go run ./cmd/demo` to demonstrate baseline, severe incident, and alert triggering | Real executable simulation generating verified output across 3 phases | MATCH |
| Test Commands | `go test ./...` and `go test -race ./...` | Tests pass cleanly with `-race` | MATCH |

## Discrepancy Findings

### 1. Endpoint Criticality Claim in Design Doc
- Location: `engineering/01-design.md:10`
- Claim: "Endpoint Criticality Bucketing: Critical endpoints (e.g. Payment) configured with stricter SLOs (99.9%) compared to non-critical endpoints (e.g. Reports: 95.0%)."
- Reality: Demo and tests define a payment service SLO (99.9%), but there is no separate secondary non-critical endpoint evaluator instantiated in demo.
- Assessment: WARNING (LOW) - Minor omission in demo demonstration, but design pattern is fully supported by configuring multiple `slo.Evaluator` instances.

### 2. Demo Alert Triggering Output in Phase 3
- Output from running `go run ./cmd/demo`:
  ```text
  >>> ALERT TRIGGERED: [TICKET] Slow Burn Alert (6.0x - 5% in 6h) | ShortBurn: 9.09x | LongBurn: 9.09x (Threshold: 6.00x)
  ```
  Both short and long trackers receive the same incident events in demo simulation, leading to identical short and long burn rates (9.09x), triggering the 6.0x rule but not the 14.4x rule. The calculation is mathematically exact and consistent with the code.
- Assessment: PASS
