# Engineering Audit Verdict

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Audit Date: 2026-09-27

## Summary

Code Files Reviewed:
- internal/inventory/model.go (17 lines)
- internal/inventory/service.go (49 lines)
- internal/inventory/store.go (173 lines)
- cmd/demo/main.go (121 lines)
Tests Reviewed:
- tests/locking_test.go (6 tests, 171 lines)
Commands Executed:
- go build ./... → PASS (exit 0)
- go test -v ./... → PASS (6/6 tests)
- go test -race ./... → PASS (zero race warnings)
- go run ./cmd/demo → PASS (all 5 scenarios executed)
Failures: 0
Warnings: 3

## Quality Gates

Compilation: PASS
Tests: PASS (6/6 pass, but with weak assertions — see gaps)
Race Detector: PASS
Demo: PASS
Research Alignment: PASS
Documentation Accuracy: PASS (with minor non-determinism note)

## Blocking Issues
(None — no CRITICAL or HIGH severity issues identified)

## Non-Blocking Issues
1. [MEDIUM] TestNaiveLostUpdate assertion too weak — only checks `stock != 50`.
2. [MEDIUM] TestOptimisticLockingWithRetry uses implementation counter as source of truth
   instead of independent invariant.
3. [MEDIUM] Missing test coverage for `ErrInvalidQuantity` (all 4 deduct methods).
4. [MEDIUM] Missing test coverage for `ErrNotFound` (all 4 deduct methods).
5. [LOW] Missing test for `AtomicDeduct` insufficient-stock path.
6. [LOW] Engineering execution-result.md recorded non-deterministic values (conflict count,
   elapsed time) without caveat.

## Required Revisions
None required for APPROVED_WITH_WARNINGS status.
Recommended (non-blocking):
- Strengthen `TestNaiveLostUpdate` assertion (assert `stock > 50`).
- Refactor `TestOptimisticLockingWithRetry` to use independent invariant.
- Add edge-case tests for `ErrInvalidQuantity` and `ErrNotFound`.
- Annotate engineering execution-result with non-determinism caveat.

## Final Status

APPROVED_WITH_WARNINGS

The implementation is correct, compiles, passes all tests including the race detector,
and the demo produces real output matching the described behavior. The three locking
strategies (pessimistic, optimistic+retry, atomic) are correctly implemented and
demonstrated. Gaps are in test coverage strength, not in implementation correctness.
