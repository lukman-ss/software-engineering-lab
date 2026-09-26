# Test Audit

Test file reviewed: `tests/locking_test.go`
Test framework: Go `testing`
Commands executed:
- `go test -v -race ./...`
- Result: all 6 tests PASS, zero race warnings

---

## Test Inventory

| # | Test Function | Scenario | Goroutines | Concurrency |
|---|---|---|---|---|
| 1 | `TestNaiveLostUpdate` | Naive read-modify-write lost update | 50 | Yes |
| 2 | `TestPessimisticLocking` | Pessimistic locking correctness | 50 | Yes |
| 3 | `TestPessimisticLockingInsufficientStock` | Insufficient stock edge case | 1 (sequential) | No |
| 4 | `TestOptimisticLockingConflict` | Optimistic conflict detection | 20 | Yes |
| 5 | `TestOptimisticLockingWithRetry` | Optimistic retry convergence | 20 | Yes |
| 6 | `TestAtomicConditionalUpdate` | Atomic single-statement update | 50 | Yes |

---

## Test 1: TestNaiveLostUpdate

Location: `tests/locking_test.go:10-38`

**Happy path coverage**: Yes — 50 concurrent `DeductNaive` calls.
**Failure path coverage**: The test asserts that the "happy path" expected value (stock == 50) is NOT reached, proving the anomaly.
**Edge cases**: No.
**Concurrency**: Yes, 50 goroutines.
**Assertions**:
- `p.Stock == 50` → FAIL (expected lost update)
- Log: "50 deduct calls occurred, but final stock is X (expected 50)"

Assessment: PASS (demonstrates the claim)
Severity: MEDIUM
Notes: The test is probabilistic — the 100µs sleep in `NaiveDeduct` makes manifestation reliable, but the assertion only checks `stock != 50`, not the exact anomaly. The `NaivelyDrawn` counter (which should be 50) is never asserted, meaning the test cannot prove that 50 deductions were actually attempted and recorded. A stronger assertion would be `store.NaivelyDrawn == 50 && p.Stock < 99` (or `p.Stock == 99`). Observed: stock = 98-99, which correctly demonstrates lost update but with non-deterministic exact value.

---

## Test 2: TestPessimisticLocking

Location: `tests/locking_test.go:40-67`

**Happy path coverage**: Yes — 50 concurrent `DeductPessimistic` calls, all expected to succeed.
**Failure path coverage**: No.
**Edge cases**: No.
**Concurrency**: Yes, 50 goroutines.
**Assertions**:
- Each `DeductPessimistic` returns no error.
- `p.Stock == 50` (exact final value).

Assessment: PASS
Severity: LOW
Notes: Strong deterministic test. The 50 concurrent deductions with per-row locking must yield exactly stock 50. Any deviation indicates a concurrency bug. The race detector confirms no data races.

---

## Test 3: TestPessimisticLockingInsufficientStock

Location: `tests/locking_test.go:69-83`

**Happy path coverage**: Partial — first deduction succeeds.
**Failure path coverage**: Yes — second deduction returns `ErrInsufficientStock`.
**Edge cases**: Yes — boundary condition where stock exactly equals first deduction.
**Concurrency**: No — sequential.
**Assertions**:
- First `DeductPessimistic(20, 2)` returns nil.
- Second `DeductPessimistic(20, 1)` returns `ErrInsufficientStock`.

Assessment: PASS
Severity: LOW
Notes: Correctly tests the insufficient stock guard. Could be strengthened with a concurrent test where stock runs out mid-flight (race between deduction and check), but the pessimistic lock should prevent that.

---

## Test 4: TestOptimisticLockingConflict

Location: `tests/locking_test.go:85-121`

**Happy path coverage**: Yes — successful optimistic deductions counted.
**Failure path coverage**: Yes — conflicts counted via `ErrOptimisticLock`.
**Edge cases**: No.
**Concurrency**: Yes, 20 goroutines.
**Assertions**:
- `successCount + p.Stock == 100` (stock invariant: successful deductions + remaining stock = initial stock).
- `conflictCount > 0` (at least one conflict detected).

Assessment: PASS with WARNING
Severity: MEDIUM
Notes: The test asserts the invariant `successCount + p.Stock == 100`, which proves no data corruption occurred (no over-decrement). It also asserts `conflictCount > 0`, proving conflicts were detected. However, the `conflictCount > 0` assertion is probabilistic — with 20 goroutines and a 50µs sleep, at least one conflict is virtually guaranteed, but not deterministic. A stronger test would assert `successCount == int(store.Optimistically)` and check that `store.Optimistically + store.OptimisticFails + p.Stock == 100`.

---

## Test 5: TestOptimisticLockingWithRetry

Location: `tests/locking_test.go:123-145`

**Happy path coverage**: Yes — retry convergence.
**Failure path coverage**: No — retry exhaustion is not tested.
**Edge cases**: No.
**Concurrency**: Yes, 20 goroutines.
**Assertions**:
- `p.Stock == 100 - int(store.Optimistically)` (accounting invariant).

Assessment: WARNING
Severity: HIGH
Notes: **Critical weakness**: The test does not assert that all 20 goroutines succeeded. The invariant `p.Stock == 100 - store.Optimistically` holds regardless of how many goroutines succeeded vs. gave up due to retry exhaustion. If `DeductOptimisticWithRetry` silently returned `ErrOptimisticLock` without retrying (a broken implementation), the invariant would still hold because `Optimistically` only counts successful deductions. The test cannot distinguish correct retry behavior from a broken implementation that gives up. A stronger test would assert `store.Optimistically == 20` (all goroutines succeeded) and/or `p.Stock == 80`. Observed: `Optimistically = 20`, `Stock = 80`, confirming convergence in practice, but the test does not enforce it.

---

## Test 6: TestAtomicConditionalUpdate

Location: `tests/locking_test.go:147-171`

**Happy path coverage**: Yes — 50 concurrent `DeductAtomic` calls.
**Failure path coverage**: No.
**Edge cases**: No.
**Concurrency**: Yes, 50 goroutines.
**Assertions**:
- Each `DeductAtomic` returns no error.
- `p.Stock == 50` (exact final value).

Assessment: PASS
Severity: LOW
Notes: Strong deterministic test. The atomic guard (`s.mu` held for read-check-write) ensures no lost updates. The exact stock assertion proves correctness.

---

## Coverage Gaps

### Missing Tests

| Gap | Description | Severity |
|-----|-------------|----------|
| No `ErrInvalidQuantity` test for any method | `qty <= 0` returns `ErrInvalidQuantity` but is never tested | MEDIUM |
| No `ErrNotFound` test for non-existent product | Deducting from a product ID that was never seeded is never tested | LOW |
| No retry exhaustion test | `DeductOptimisticWithRetry` returning `ErrOptimisticLock` when `maxRetries` is exhausted is not tested | HIGH |
| No concurrent insufficient-stock test | Pessimistic or atomic deduction racing against stock depletion is not tested | LOW |
| No assertion that retry test converges | `TestOptimisticLockingWithRetry` does not assert `Optimistically == 20` or `Stock == 80` | HIGH |
| No negative case for optimistic direct | `DeductOptimisticDirect` error paths beyond `ErrOptimisticLock` are not tested | LOW |
| No concurrent pessimistic insufficient-stock | Concurrent pessimistic deductions that race to exhaust stock is not tested | LOW |
| `TestOptimisticLockingConflict` probabilistic | `conflictCount > 0` is timing-dependent, not deterministic | MEDIUM |

### Race Detector Results

`go test -race ./...` passed with zero warnings. No data races detected on `Store`, `Product`, `Service`, or counter fields. The `math/rand.Intn` call in the retry loop is safe for concurrent use under Go 1.20+ semantics (go.mod specifies `go 1.22`).
