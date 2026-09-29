# Documentation vs Code Audit

Target Lab: labs/38-mutation-testing

## Comparison Matrix

| Claim / Topic | Documentation Claim (README / Engineering) | Code Implementation | Status |
| :--- | :--- | :--- | :--- |
| **Formula** | $\text{Score} = (\text{Killed} / \text{Total}) \times 100\%$ | Implemented in `internal/engine/runner.go:103` | MATCH |
| **Operators** | README lists 4 operators: Relational, Boolean, Arithmetic, Boundary Shift | `internal/engine/mutator.go` implements exactly those 4 operators | MATCH |
| **Statement Deletion** | Mentioned in `01-design.md:17` and `types.go:12` | Omitted from AST mutator; noted in `02-implementation-notes.md:44` | MATCH (documented limitation) |
| **Weak Test Coverage** | Claims 100% statement coverage with weak assertions | Verified at 100.0% coverage via `go test -coverprofile` | MATCH |
| **Demo Output** | `03-execution-result.md` claims 15 mutants generated, Weak: 0/15 (0.00%), Strong: 15/15 (100.00%) | `go run ./cmd/demo` produces identical 15 mutant breakdown and scores | MATCH |
| **Race Detector** | Claims clean execution with zero data races | Verified via `go test -race ./...` (0 data races detected) | MATCH |

## Identified Discrepancies

- None between `README.md` and code. `README.md` cleanly documents the 4 active operators without claiming unimplemented features.
- Minor discrepancy in `01-design.md` mentioning statement deletion, but appropriately reconciled in `02-implementation-notes.md`.
