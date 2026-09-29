# Test Audit

Target Lab: labs/38-mutation-testing

## Test Files

- `internal/service/discount_weak_test.go` — 4 test cases
- `internal/service/discount_strong_test.go` — 7 table-driven test cases
- `tests/engine_test.go` — 3 test functions

## Actual Execution Results

```
go test -v ./...
ok  labs/38-mutation-testing/internal/service
ok  labs/38-mutation-testing/tests

go test -race ./...
ok  labs/38-mutation-testing/internal/service
ok  labs/38-mutation-testing/tests
```

All 15 test sub-cases passed. Zero race conditions detected.

## Statement Coverage

```
go test -coverprofile=coverage.out ./internal/service
coverage: 100.0% of statements
CalculateDiscount: 100.0%
```

100% statement coverage confirmed on `discount.go`.

## Coverage by Category

### Happy Path
- PASS: All three customer tier paths (VIP, Premium, Standard) exercised.
- PASS: Coupon addition with >2 items covered.
- PASS: Free shipping threshold at 200.0 (exact boundary).

### Failure Path
- PASS: `weakTestFn` fails to detect any mutation (0% score).
- PASS: Weak test passes against original — deliberate design to show false confidence.

### Edge Cases
- PASS: `Coupon with <=2 items does NOT add extra rate` covers boundary.
- PASS: `Free shipping threshold at exact 200 final amount` — `200.0 * 0.05 = 10.0 → final = 190.0`, correctly expects no free shipping.
- PASS: Exact equality boundary `>=` tested with amount at exactly 500, 1000, 200.

### Boundary / Relational Transition
- PASS: Strong test covers `>=1000`, `>=500`, `>=100` thresholds directly.
- WARNING: No test for amount just below boundary (999.99, 499.99, 99.99) or exactly 0.

### Recovery / Rollback (engine)
- PASS: Runner re-parses fresh AST per goroutine. No shared mutable state.
- PASS: `undo()` functions exist and are called even on render errors (runner.go:56).

### Concurrency
- PASS: `go test -race ./...` passes clean.

### Negative / Degenerate Cases
- MISSING: No test for `ItemCount <= 0`, `TotalAmount = 0`, or empty `Tier` string.
- MISSING: No test for the edge where coupon is `true` but `ItemCount == 0`.

### Engine-Level Tests
- PASS: `TestEngine_GeneratesMutants` validates all 4 operator categories appear.
- PASS: `TestEngine_MutationScoreDifference` confirms `strongScore > weakScore` and `strongScore == 100.0`.
- WARNING: `TestDomainService_DirectCalculation` (tests/engine_test.go:112-126) expects `DiscountRate == 0.25`. VIP tier starts at 0.20, coupon adds 0.05 → 0.25. Correct. But this test is in the `tests` package, not `service`, and exercises service directly — minor scope concern (low severity).

## Summary

| Category             | Status       |
|---------------------|--------------|
| Happy Path           | PASS         |
| Exact Boundary       | PASS         |
| Weak Suite Exposed   | PASS         |
| Strong Suite Kills   | PASS         |
| Race Detection       | PASS         |
| Engine Mutation Ops  | PASS         |
| Degenerate Inputs    | NOT COVERED  |
| Below-Boundary Cases | NOT COVERED  |
