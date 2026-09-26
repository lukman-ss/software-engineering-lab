# Test Audit

## Coverage Matrix

| Scenario | Test | Status |
|---|---|---|
| Naive lost update | TestNaiveLostUpdate | present |
| Pessimistic invariant | TestPessimisticLocking | present |
| Pessimistic insufficient stock | TestPessimisticLockingInsufficientStock | present |
| Optimistic conflict detection | TestOptimisticLockingConflict | present |
| Optimistic retry convergence | TestOptimisticLockingWithRetry | present |
| Atomic invariant | TestAtomicConditionalUpdate | present |

## Test-by-test

### TestNaiveLostUpdate
- Happy path: 50 goroutines deduct 1. PASSES when stock != 50.
- Flaw: inverted assertion. Test cannot FAIL if lost update does not manifest; it only fails if naive path happens to be fully correct (stock==50). Nondeterministic timing means test may mask regression.
- Edge cases: single-threaded lost-update not tested.
- Assessment: WEAK.

### TestPessimisticLocking
- Asserts exact invariant stock==50. Correct and deterministic under locking.
- Concurrent 50 goroutines, errors flagged. Good.

### TestPessimisticLockingInsufficientStock
- Single-threaded oversell guard (stock 2, second 1 fails). Positive coverage of failure path.
- Missing: concurrent oversell boundary for pessimistic.

### TestOptimisticLockingConflict
- Asserts invariant successCount+stock==100 and conflictCount>0.
- Flaw: `TestOptimisticLockingWithRetry` (same store counter reuse) — note `store.OptimisticFails` is separate. Invariant is correct for no-retry path.
- Nondeterminism: relies on at least one conflict occurring — reliable with 50µs delay + 20 goroutines, but not guaranteed in principle.

### TestOptimisticLockingWithRetry
- Asserts consistency via counter: `p.Stock == 100 - store.Optimistically`.
- Flaw: `store.Optimistically` is sum of applied quantities across retries; if a retry's committed update increments counter, equality holds. Observed: 20 applied. Assertion correct.
- No assertion that stock >= 0 under concurrent drain-to-zero (oversell protection for retry path not directly tested).

### TestAtomicConditionalUpdate
- Asserts exact stock==50 under 50 goroutines. Deterministic and correct.

## Gaps
- No negative test: invalid quantity (ErrInvalidQuantity) not asserted for any strategy.
- No negative test: missing product ID not asserted for any strategy.
- No oversell boundary test for optimistic path (concurrent drain to 0).
- No concurrency test for `Get` / `Seed` read-while-write.

## Race Detector
Command: `go test -race -v ./...`
Result: `ok ... tests ...` — PASS, zero race warnings.
