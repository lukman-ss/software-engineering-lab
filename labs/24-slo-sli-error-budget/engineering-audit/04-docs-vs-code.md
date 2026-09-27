# Documentation vs Code Verification

## Artifact Comparison

| Claim / Section | Documented Location | Code / Reality | Status |
| :--- | :--- | :--- | :--- |
| SLI Ratio Definition | `README.md`, `engineering/01-design.md`, `research/05-report.md` | `internal/slo/evaluator.go:44-47` (`good / total`) | MATCH |
| Error Budget Formula | `README.md`, `engineering/01-design.md`, `research/05-report.md` | `internal/slo/evaluator.go:49-52` (`(1 - target) * total`) | MATCH |
| Release Freeze Behavior | `README.md`, `engineering/01-design.md` | `internal/slo/evaluator.go:54-57` (`budgetRemaining <= 0 -> CanDeploy = false`) | MATCH |
| Multi-Window Alerting | `engineering/01-design.md`, `engineering/02-implementation-notes.md` | `internal/alerting/engine.go:73` (`shortBurn >= factor && longBurn >= factor`) | MATCH |
| Test Commands | `README.md:15-18` | `go test ./...` and `go test -race ./...` run and pass directly | MATCH |
| Demo Execution Output | `engineering/03-execution-result.md:48-73` | Output of `go run ./cmd/demo` matches verbatim | MATCH |

## Discrepancies Found
- None. No `DOC_CODE_MISMATCH`, `TEST_CLAIM_MISMATCH`, or `RESEARCH_IMPLEMENTATION_MISMATCH` identified.
