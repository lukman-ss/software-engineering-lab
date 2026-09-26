# Engineering Audit Verdict

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Audit Date: 2026-09-26

## Summary

Code Files Reviewed:
- internal/inventory/model.go
- internal/inventory/store.go
- internal/inventory/service.go
- cmd/demo/main.go
- go.mod

Tests Reviewed:
- tests/locking_test.go

Commands Executed:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo

Failures:
- None (all tests pass, build succeeds, race detector clean, demo runs)

Warnings:
- TestNaiveLostUpdate uses inverted assertion (only fails if lost update does NOT occur); timing-dependent pass/fail.
- Missing negative/error-path tests for all strategies (invalid ID, invalid quantity).
- Missing optimistic oversell boundary tests (concurrent drain to zero).
- Demo/Test counters vary due to scheduling; recorded values in engineering/03-execution-result.md are one possible outcome.

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS
Demo: PASS
Research Alignment: NOT_APPLICABLE (PIPELINE OVERRIDE: implementation and tests only)
Documentation Accuracy: PASS (README matches implementation/demo; engineering notes consistent)

## Blocking Issues
1. None

## Non-Blocking Issues
1. TestNaiveLostUpdate assertion invert (LOW): test could hide regression if timing accidentally yields stock==50; better to assert stock < 50 (any lost update) or compute expected lost range.
2. Missing negative tests (LOW): no coverage for ErrInvalidQuantity, ErrNotFound per strategy.
3. Missing optimistic oversell edge test (LOW): no test verifies optimistic paths never produce negative stock under concurrent drain-to-zero.

## Required Revisions
1. Strengthen TestNaiveLostUpdate: assert final stock < initial stock - (goroutines - 1) or simply stock < 50 for 50 goroutines.
2. Add Test<Strategy>InvalidQuantity for each deduction method.
3. Add Test<Strategy>NotFound for each deduction method.
4. Add TestOptimisticNeverOversells: many goroutines decrement from low stock; assert stock >= 0 and no lost increments.

## Final Status

APPROVED

All quality gates passed:
- Code compiles without errors.
- Test suite passes (including race detector).
- Demo executes and demonstrates claimed behaviors (lost update, pessimistic correctness, optimistic conflict detection/retry, atomic correctness).
- No unresolved HIGH/CRITICAL issues; only LOW-risk gaps in test coverage.
- Documentation (README, engineering notes) matches implementation and observed demo output.