# Docs vs Code Audit

Target Lab: `labs/24-slo-sli-error-budget`

## Comparison Matrix

| Aspect | Documentation / Research Claim | Code / Demo Implementation | Status |
|---|---|---|---|
| Project Structure | `internal/metrics`, `internal/slo`, `internal/alerting`, `cmd/demo`, `tests/` | Identical packages and layout exist in repo | MATCH |
| SLI Metric Model | Good events / Total events ratio using sliding window | Implemented in `tracker.go` and `evaluator.go` | MATCH |
| Error Budget Tracking | Budget = `(1 - SLO) * total - bad`; freeze when budget <= 0 | Implemented in `evaluator.go:49-57` | MATCH |
| Multi-Window Multi-Burn-Rate Alerting | Both short and long window burn rates must breach threshold factor | Implemented in `engine.go:73` | MATCH |
| Demo Execution | 3 Phases (Baseline -> Severe Incident -> Alerting) | Implemented in `cmd/demo/main.go`, output is 100% reproducible | MATCH |
| Execution Output in Engineering Notes | Claims passing test suite and specific demo logs | Verified identical output when running `go test` and `go run ./cmd/demo` | MATCH |

## Observed Mismatches

None. Documentation accurately describes implementation and demo output.
