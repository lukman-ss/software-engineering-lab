# Engineering Audit Verdict

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Audit Date: 2026-09-27

## Summary

Code Files Reviewed: internal/inventory/model.go, internal/inventory/store.go, internal/inventory/service.go, cmd/demo/main.go
Tests Reviewed: tests/locking_test.go (6 tests)
Commands Executed: go vet ./... (exit 0), go test -v ./... (6/6 PASS), go test -race -count=1 ./... (PASS), go run ./cmd/demo (exit 0)
Failures: 0
Warnings: 7 gaps (4 MEDIUM missing-test, 3 LOW edge-case), none blocking

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS
Demo: PASS
Research Alignment: NOT_APPLICABLE (out of scope per pipeline override)
Documentation Accuracy: PASS

## Blocking Issues

None. No HIGH or CRITICAL gaps. All four claimed behaviors proven live:
- Naive: 50 deducts on stock 100 -> final 99 (lost update reproduced)
- Pessimistic: 50 deducts -> final 50 (exact invariant)
- Optimistic direct: 1 success / 19 conflicts, successCount + stock == 100 (guard holds)
- Optimistic retry: 20/20 converged (stock 80)
- Atomic: 50 deducts -> final 50 (exact invariant)

## Non-Blocking Issues

1. [MEDIUM, MISSING_TEST] No ErrInvalidQuantity test on any strategy (05-gaps.md GAP-1)
2. [MEDIUM, MISSING_TEST] No ErrNotFound test on any strategy (GAP-2)
3. [MEDIUM, MISSING_TEST] Insufficient-stock failure tested only for pessimistic locking (GAP-3)
4. [MEDIUM, MISSING_TEST] Retry exhaustion path untested; retry test discards errors (GAP-4, GAP-7)
5. [LOW, MISSING_EDGE_CASE] No post-failure stock-unchanged assertion (GAP-5)
6. [LOW, MISSING_EDGE_CASE] `conflictCount == 0` assertion timing-dependent; invariant assertion already covers safety deterministically (GAP-6)

## Required Revisions

None blocking. Recommended test hardening before publication:
1. Add invalid-quantity (0, negative) tests for all four strategies.
2. Add not-found tests for all four strategies.
3. Add insufficient-stock tests for naive, optimistic direct, retry, and atomic paths.
4. Add retry-exhaustion test (maxRetries=0 under contention -> ErrOptimisticLock) and assert retry errors instead of discarding.
5. Assert stock unchanged after each rejected deduct.

## Final Status

APPROVED_WITH_WARNINGS
