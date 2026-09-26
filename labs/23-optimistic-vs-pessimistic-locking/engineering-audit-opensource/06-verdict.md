# Engineering Audit Verdict

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Audit Date: 2026-09-26

## Summary

Code Files Reviewed: 4 (`internal/inventory/model.go`, `internal/inventory/service.go`, `internal/inventory/store.go`, `cmd/demo/main.go`)
Tests Reviewed: 1 (`tests/locking_test.go`, 6 test functions)
Commands Executed:
- `go test ./...` → PASS
- `go test -v -race ./...` → PASS (6/6 tests, 0 race warnings)
- `go run ./cmd/demo` → PASS (all 5 scenarios produce expected output)
- `go vet ./...` → PASS
Failures: 0
Warnings: 2 (HIGH severity test coverage gaps, 1 LOW doc mismatch)

## Quality Gates

| Gate | Status |
|------|--------|
| Compilation | PASS |
| Tests | PASS |
| Race Detector | PASS |
| Demo | PASS |
| Research Alignment | PASS |
| Documentation Accuracy | WARNING |

## Blocking Issues

1. **G3 (HIGH)**: No test for retry exhaustion — `DeductOptimisticWithRetry` returning `ErrOptimisticLock` when `maxRetries` is exhausted. The retry loop's failure-termination condition is untested.
2. **G4 (HIGH)**: `TestOptimisticLockingWithRetry` does not assert all goroutines succeeded. The invariant `p.Stock == 100 - store.Optimistically` holds even if some goroutines gave up without retrying, making the test unable to detect a broken retry implementation.

## Non-Blocking Issues

1. **G7 (LOW)**: `engineering/01-design.md:33` claims "4 scenarios" but the demo has 5 printed sections.
2. **G6 (MEDIUM)**: `TestOptimisticLockingConflict` uses a probabilistic assertion (`conflictCount > 0`) rather than a deterministic one.
3. **G8 (MEDIUM)**: `TestNaiveLostUpdate` asserts only `stock != 50` rather than the specific expected lost-update value (98-99).
4. **G5 (MEDIUM)**: `TestNaiveLostUpdate` does not verify the `NaivelyDrawn` counter to prove all 50 deductions were attempted.
5. **G1 (MEDIUM)**: No test for `ErrInvalidQuantity` on any deduction method.
6. **G2 (LOW)**: No test for `ErrNotFound` on any deduction method for non-existent product.

## Required Revisions

1. Add a test that verifies `DeductOptimisticWithRetry` returns `ErrOptimisticLock` when `maxRetries` is exhausted (e.g., `maxRetries=0` with a guaranteed conflict, or inject a scenario where retries exhaust).
2. Strengthen `TestOptimisticLockingWithRetry` to assert `store.Optimistically == 20` (or `p.Stock == 80`), proving all goroutines converged successfully.
3. (Optional) Add tests for `ErrInvalidQuantity` and `ErrNotFound` error paths.
4. Fix `engineering/01-design.md:33` to say "5 scenarios" instead of "4 scenarios."

## Final Status

APPROVED_WITH_WARNINGS

**Rationale**: The implementation is correct and all code audits pass. The demo produces real, reproducible output proving all five concurrency scenarios. The race detector confirms zero data races. However, two HIGH-severity test coverage gaps remain unresolved: the retry-exhaustion failure path is untested, and the retry-convergence test cannot detect a broken retry implementation. The implementation itself is sound; the tests simply do not fully prove the claimed behavior of the retry logic. Per the audit standard, APPROVED is withheld until these high-severity test gaps are addressed.
