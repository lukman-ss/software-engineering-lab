# Docs vs Code Audit

Target Lab: labs/38-mutation-testing

## Claim Verification

| Item | Claim in README | Reality in Code | Status |
|------|-----------------|-----------------|--------|
| **Directory Structure** | Shows `cmd/demo/main.go`, `internal/engine/*`, `internal/service/*`, `tests/engine_test.go`, `engineering/` | Matches file system layout exactly | MATCH |
| **Formula** | $\text{Score} = (\text{Killed}/\text{Total}) \times 100\%$ | `engine/runner.go:105` computes `float64(killed) / float64(total) * 100.0` | MATCH |
| **Operators** | 4 operators: Relational, Boolean, Arithmetic, Boundary Shift | `internal/engine/mutator.go` implements all 4 | MATCH |
| **Statement Deletion** | Omitted from README.md, listed in `01-design.md` and `types.go` | Not implemented in `mutator.go`, noted in `02-implementation-notes.md` | MINOR_MISMATCH (Internal docs only) |
| **Weak Test Coverage** | 100% line coverage | `go test -cover` verifies 100.0% statement coverage on `discount.go` | MATCH |
| **Weak vs Strong Contrast** | Weak suite scores 0%, strong suite scores 100% | `cmd/demo` and `tests/engine_test.go` yield 0% weak, 100% strong | MATCH |
| **Execution Commands** | `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | All 3 commands execute without error | MATCH |

## Demo Output Verification

README does not paste full demo output, but `03-execution-result.md` contains recorded output.

Comparing recorded output vs actual `go run ./cmd/demo` run:

- Total mutants: 15 generated (identical)
- Mutant breakdowns: exact line matches on Lines 31, 33, 35, 40, 44, 45, 48 (identical)
- Weak suite: Total 15, Killed 0, Survived 15, Score 0.00% (identical)
- Strong suite: Total 15, Killed 15, Survived 0, Score 100.00% (identical)

Zero fake or fabricated demo output detected.

## Discrepancies Found

1. `01-design.md:17` lists "Statement deletion (void calls or assignments neutralized)" under Expected Behavior, but `mutator.go` does not implement this operator. This is acknowledged and explained in `02-implementation-notes.md` and correctly excluded from `README.md`.
2. `tests/engine_test.go:73-98` and `cmd/demo/main.go:35-43` use simulated test oracle functions (`TestFunc`) based on source strings rather than dynamic out-of-process compilation. This trade-off is documented in `02-implementation-notes.md:45-56`.
