# Gap Analysis

Target Lab: labs/38-mutation-testing

## Identified Gaps

### 1. Statement Delete Operator Unimplemented
- Type: `MISSING_IMPLEMENTATION`
- Location: `internal/engine/types.go:13`, `engineering/01-design.md:17`
- Description: `StatementDelete` is declared in `types.go` and listed in design doc, but not handled by AST walker in `mutator.go`.
- Severity: LOW
- Mitigation: Documented as intentional limitation in `02-implementation-notes.md`. README.md does not claim statement deletion.

### 2. Test Oracle In-Memory Source Comparison Trade-Off
- Type: `IMPLEMENTATION_OVERCLAIM`
- Location: `cmd/demo/main.go:35-43`, `tests/engine_test.go:73-98`
- Description: The runner passes mutated source code to in-memory `TestFunc` closures instead of compiling binaries and running `go test` via subprocess.
- Severity: MEDIUM
- Mitigation: Explicitly noted in `02-implementation-notes.md:49-56` under Trade-offs. Acceptable for lab scale.

### 3. Missing Degenerate/Negative Input Test Cases
- Type: `MISSING_TEST`
- Location: `internal/service/discount_strong_test.go`
- Description: No tests for zero values (`TotalAmount: 0`, `ItemCount: 0`), negative amounts, or unassigned tier string.
- Severity: LOW
- Mitigation: Core happy-path and boundary test cases are complete and achieve 100% statement coverage.

### 4. Unused Mutex in Runner
- Type: `CODE_QUALITY`
- Location: `internal/engine/runner.go:16`
- Description: Struct `Runner` defines `mu sync.Mutex` which is never acquired or locked.
- Severity: LOW
- Mitigation: Harmless dead code field, does not affect functionality or safety.

## Summary

No HIGH or CRITICAL severity gaps identified.
No race conditions, broken code, fake benchmarks, or unverified output found.
