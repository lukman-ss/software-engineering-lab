# Documentation vs Code Audit

Target Lab: `labs/38-mutation-testing`

## Comparison Matrix

| Claim / Item | Documentation Location | Code Implementation | Status | Notes |
| :--- | :--- | :--- | :--- | :--- |
| Project Structure | `README.md:15-31` | Repository tree | PASS | Exact file tree match |
| Implemented Operators | `README.md:48-53` | `internal/engine/mutator.go` | PASS | Relational, Boolean, Arithmetic, Boundary Shift documented and implemented |
| 100% Statement Coverage Weak Suite | `README.md:3`, `01-design.md:11` | `internal/service/discount_weak_test.go` | PASS | Verified 100.0% coverage via `go tool cover` |
| Mutation Score Calculation Formula | `README.md:9-11` | `internal/engine/runner.go:101-105` | PASS | `(killed / total) * 100.0` accurately calculated |
| Race Detection Cleanliness | `README.md:39`, `01-design.md:30` | `internal/engine/runner.go` | PASS | Verified with `go test -count=1 -race ./...` |
| CLI Demo Execution | `README.md:45` | `cmd/demo/main.go` | PASS | Matches output in `engineering/03-execution-result.md` |
| Statement Deletion Operator | `engineering/01-design.md:17`, `types.go:12` | `internal/engine/mutator.go` | WARNING | Omitted from AST walker; documented in `02-implementation-notes.md` |

## Summary of Discrepancies
- No blocking mismatches (`DOC_CODE_MISMATCH` = 0).
- `README.md` correctly lists only the 4 actively implemented operators, avoiding over-claiming.
