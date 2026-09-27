# Docs vs Code Audit

Target Lab: `labs/24-slo-sli-error-budget`

## Comparison Matrix

| Item | Documented Claim | Observed Implementation | Result |
| :--- | :--- | :--- | :--- |
| **SLI Definition** | Good requests divided by total valid requests | `sli = float64(good) / float64(total)` in `evaluator.go:46` | MATCH |
| **Error Budget Definition** | `(1 - SLO) * total` events | `totalErrorBudget := (1.0 - TargetUptime) * float64(total)` in `evaluator.go:50` | MATCH |
| **Deployment Freeze Policy** | Halts deployments when budget exhausted (`budgetRemaining <= 0`) | `canDeploy = false` when `budgetRemaining <= 0` in `evaluator.go:55` | MATCH |
| **Burn Rate Alerting** | Multi-window burn rate alert fires only when short & long windows breach factor | `shortBurn >= rule.BurnRateFactor && longBurn >= rule.BurnRateFactor` in `engine.go:73` | MATCH |
| **Endpoint Criticality** | Strict SLO on critical endpoints (e.g. 99.9% Payment) vs relaxed on non-critical (95.0% Reports) | Demonstrated in `cmd/demo/main.go:117-147` | MATCH |
| **Demo Output Accuracy** | Demo output recorded in `engineering/03-execution-result.md` | Exact verbatim match with real `go run ./cmd/demo` execution output | MATCH |

## Discrepancies Found

None. README, engineering design, implementation notes, and demo logs precisely mirror code behavior.
