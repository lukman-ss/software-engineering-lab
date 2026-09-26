# Engineering Audit Verdict

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Audit Date: 2026-09-26

## Summary

Code Files Reviewed:
- internal/inventory/model.go
- internal/inventory/store.go
- internal/inventory/service.go
- cmd/demo/main.go

Tests Reviewed:
- tests/locking_test.go

Commands Executed:
- `go test ./...` → PASS
- `go test -race ./...` → PASS (1.125s)
- `go build ./...` → PASS
- `go run ./cmd/demo` → PASS

Failures:
- None

Warnings:
- Version not incremented by non-optimistic write paths (simulation fidelity)
- Missing edge case tests (zero/negative quantity, non-existent product)
- Demo records timing-dependent values as static

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS
Demo: PASS
Research Alignment: PASS (core behavior matches claims)
Documentation Accuracy: PASS (minor wording notes on "lockless" atomic)

## Blocking Issues
None

## Non-Blocking Issues
1. **Simulation Fidelity (MEDIUM)**: Optimistic locking does not detect changes made by non-optimistic write paths because only `OptimisticDeduct` increments the product version. This does not affect isolated strategy tests but limits simulation realism.
2. **Missing Edge Case Tests (MEDIUM)**: No tests for zero/negative quantity or non-existent product across any deduction strategy.
3. **Documentation Precision (LOW)**: README and demo output describe atomic operations as "lockless" or "single-statement," while implementation uses a mutex. Conflict counts and elapsed times in recorded demo output are timing-dependent but presented as static values.

## Required Revisions
1. Add tests for zero/negative quantity input (expecting `ErrInvalidQuantity`).
2. Add tests for non-existent product ID (expecting `ErrNotFound`).
3. Clarify in documentation that atomic operations are mutex-guarded simulations, not literally lockless SQL statements.
4. Note in engineering/03-execution-result.md that conflict counts and elapsed times vary per run due to concurrency timing.

## Final Status
APPROVED

### Rationale
- Code compiles successfully
- All required tests pass (including race detector)
- Core behavior is proven: 
  - Lost update anomaly demonstrated (naive strategy)
  - Pessimistic locking prevents conflicts (exact invariant)
  - Optimistic locking detects conflicts and converges with retry
  - Atomic operations maintain consistency under concurrency
- README accurately describes usage and expected behavior
- No unresolved HIGH/CRITICAL issues
- Warnings are non-blocking and relate to simulation fidelity or test coverage, not correctness of claimed behavior