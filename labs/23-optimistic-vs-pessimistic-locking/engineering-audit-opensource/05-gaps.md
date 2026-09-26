# Gap Analysis

## Identified Gaps

| ID | Type | Location | Description | Severity |
|----|------|----------|-------------|----------|
| G1 | MISSING_TEST | tests/locking_test.go | No test for `ErrInvalidQuantity` (qty <= 0) on any deduction method | MEDIUM |
| G2 | MISSING_TEST | tests/locking_test.go | No test for `ErrNotFound` on any deduction method for non-existent product ID | LOW |
| G3 | MISSING_TEST | tests/locking_test.go | No test for retry exhaustion in `DeductOptimisticWithRetry` (returning `ErrOptimisticLock` when maxRetries exceeded) | HIGH |
| G4 | MISSING_TEST | tests/locking_test.go | No assertion of convergence in `TestOptimisticLockingWithRetry` (does not assert all goroutines succeeded) | HIGH |
| G5 | MISSING_TEST | tests/locking_test.go | No test verifying `NaivelyDrawn` counter in `TestNaiveLostUpdate` (to prove 50 deductions were attempted) | MEDIUM |
| G6 | WEAK_TEST | tests/locking_test.go | `TestOptimisticLockingConflict` relies on probabilistic conflict assertion (`conflictCount > 0`) | MEDIUM |
| G7 | DOC_CODE_MISMATCH | engineering/01-design.md:33 | Claims "4 scenarios" but demo has 5 sections ([1] through [5]) | LOW |
| G8 | WEAK_TEST | tests/locking_test.go | `TestNaiveLostUpdate` only asserts `stock != 50`, not the specific lost-update value (should be 98-99) | MEDIUM |
| G9 | IMPLEMENTATION_OVERCLAIM | N/A | None found. All implementation claims are verified by tests or demo. | N/A |

### G1: Missing `ErrInvalidQuantity` Test
**Location**: Missing in tests/locking_test.go  
**Description**: None of the test cases verify that `DeductNaive`, `DeductPessimistic`, `DeductOptimisticDirect`, `DeductOptimisticWithRetry`, or `DeductAtomic` return `ErrInvalidQuantity` when `qty <= 0`.  
**Severity**: MEDIUM  
**Impact**: Error-handling paths are untested. A regression in input validation would not be caught.

### G2: Missing `ErrNotFound` Test
**Location**: Missing in tests/locking_test.go  
**Description**: None of the test cases verify deduction methods return `ErrNotFound` when the product ID has never been seeded.  
**Severity**: LOW  
**Impact**: Low because seeding is required before deduction in the API, but it's still an untested error path.

### G3: Missing Retry Exhaustion Test
**Location**: Missing in tests/locking_test.go  
**Description**: No test verifies that `DeductOptimisticWithRetry` returns `ErrOptimisticLock` when `maxRetries`maxRetries` is exhausted (e.g., with `maxRetries=0` or under extreme contention).  
**Severity**: HIGH  
**Impact**: The retry loop's termination condition is untested. A bug causing infinite retry or premature success would not be caught.

### G4: Weak Convergence Assertion in Retry Test
**Location**: `tests/locking_test.go:140-144` (`TestOptimisticLockingWithRetry`)  
**Description**: The test only asserts the accounting invariant `p.Stock == 100 - store.Optimistically`. It does not assert that all 20 goroutines succeeded (i.e., `store.Optimistically == 20`). The invariant holds even if some goroutines gave up due to retry exhaustion, making the test unable to detect a broken retry implementation.  
**Severity**: HIGH  
**Impact**: False confidence in retry logic. If `DeductOptimisticWithRetry` incorrectly returned `ErrOptimisticLock` on the first attempt, the test would still pass because `Optimistically` would be <20 and stock would adjust accordingly.

### G5: Missing `NaivelyDrawn` Assertion in Lost Update Test
**Location**: `tests/locking_test.go:10-38` (`TestNaiveLostUpdate`)  
**Description**: The test proves lost update by showing `stock != 50`, but does not verify that 50 deductions were recorded via the `NaivelyDrawn` counter.  
**Severity**: MEDIUM  
**Impact**: The test cannot distinguish between "fewer deductions attempted" vs "deductions attempted but lost." Verifying `NaivelyDrawn == 50` would strengthen the proof.

### G6: Probabilistic Conflict Assertion
**Location**: `tests/locking_test.go:117` (`TestOptimisticLockingConflict`)  
**Description**: The test asserts `conflictCount > 0`, which relies on timing-dependent behavior (20 goroutines with 50µs sleep). While virtually guaranteed, it is not deterministic.  
**Severity**: MEDIUM  
**Impact**: In extremely unlikely scheduling scenarios, the test could flake. Strengthening with additional invariants would improve reliability.

### G7: Documentation-Code Mismatch
**Location**: `engineering/01-design.md:33`  
**Description**: States "Visual CLI runner showing all 4 scenarios" but the demo (`cmd/demo/main.go`) prints 5 distinct sections: [1] Naive, [2] Pessimistic, [3] Optimistic Direct, [4] Optimistic With Retry, [5] Atomic.  
**Severity**: LOW  
**Impact**: Minor documentation inaccuracy that could cause confusion.

### G8: Weak Lost Update Assertion
**Location**: `tests/locking_test.go:34` (`TestNaiveLostUpdate`)  
**Description**: The test asserts `p.Stock == 50` → FAIL (expected lost update). While correct, it does not specify the expected lost-update value (which should be 98-99 given the 100µs sleep).  
**Severity**: MEDIUM  
**Impact**: Less precise than asserting the exact expected outcome of the demonstrated anomaly.

### G9: No Implementation Overclaim
**Description**: All implementation claims (lost update, pessimistic locking, optimistic conflict detection, retry convergence, atomic updates) are verified by either tests or demo output. No evidence of the implementation claiming behavior it does not possess.  
**Severity**: N/A  

## Gap Severity Summary
- **CRITICAL**: 0
- **HIGH**: 2 (G3, G4)
- **MEDIUM**: 5 (G1, G5, G6, G8, plus one implicit from research alignment)
- **LOW**: 2 (G2, G7)