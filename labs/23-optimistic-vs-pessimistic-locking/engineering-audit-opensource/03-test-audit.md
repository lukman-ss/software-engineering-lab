## Finding 1

Location: tests/locking_test.go:10-38 (TestNaiveLostUpdate)
Covered Behavior: Concurrency happy path for naive strategy (demonstrating lost update).
Test Coverage:
- ✅ Launches 50 goroutines each calling DeductNaive(1,1).
- ✅ Waits for completion.
- ✅ Checks final stock is NOT 50 (lost update manifested).
- ❌ Does not verify that exactly 50 deduct calls succeeded (they ignore return error, which is always nil for naive with sufficient stock).
- ❌ Does not test edge cases: insufficient stock, invalid quantity, zero quantity.
- ❌ Does not verify stock unchanged on failure (not applicable as naive never fails with sufficient stock).
Assessment: PARTIAL
Missing: Negative test cases (insufficient stock, invalid quantity). Edge case boundary values.

## Finding 2

Location: tests/locking_test.go:40-67 (TestPessimisticLocking)
Covered Behavior: Concurrency happy path for pessimistic locking.
Test Coverage:
- ✅ 50 goroutines each DeductPessimistic(2,1).
- ✅ All return nil (checked).
- ✅ Final stock = 50.
- ❌ Does not test insufficient stock edge case (covered in separate test).
- ❌ Does not test invalid quantity or zero quantity.
Assessment: PARTIAL
Missing: Invalid quantity and zero quantity tests.

## Finding 3

Location: tests/locking_test.go:69-83 (TestPessimisticLockingInsufficientStock)
Covered Behavior: Failure path for pessimistic locking (insufficient stock).
Test Coverage:
- ✅ First deduct of 2 succeeds (stock 2->0).
- ✅ Second deduct of 1 returns ErrInsufficientStock.
- ❌ Does NOT verify that stock remains 0 after the failed deduct (should be unchanged).
- ❌ Does not test that multiple failures leave stock unchanged.
- ❌ Does not test negative/zero quantity.
Assessment: PARTIAL
Missing: Post-failure stock validation. Negative/zero quantity tests.

## Finding 4

Location: tests/locking_test.go:85-121 (TestOptimisticLockingConflict)
Covered Behavior: Concurrency for optimistic locking direct (no retry).
Test Coverage:
- ✅ 20 goroutines each DeductOptimisticDirect(3,1).
- ✅ Tracks success and conflict counts via mutex.
- ✅ Verifies invariant: successCount + finalStock == 100 (no lost updates).
- ✅ Checks at least one conflict occurred (flaky but usually passes).
- ❌ Does not test insufficient stock scenario.
- ❌ Does not test invalid quantity or zero quantity.
- ❌ Does not test product not found.
Assessment: PARTIAL
Missing: Error condition tests (insufficient stock, invalid quantity). Edge case boundary values.

## Finding 5

Location: tests/locking_test.go:123-145 (TestOptimisticLockingWithRetry)
Covered Behavior: Concurrency for optimistic locking with retry.
Test Coverage:
- ✅ 20 goroutines each DeductOptimisticWithRetry(4,1,10).
- ✅ All return nil (retry hides conflicts).
- ✅ Verifies final stock = 100 - Optimistically counter (no lost updates).
- ❌ Does not test retry exhaustion (should return ErrOptimisticLock after maxRetries).
- ❌ Does not test insufficient stock scenario.
- ❌ Does not test invalid quantity or zero quantity.
- ❌ Does not test product not found.
Assessment: PARTIAL
Missing: Retry exhaustion test. Error condition tests (insufficient stock, invalid quantity). Edge case boundary values.

## Finding 6

Location: tests/locking_test.go:147-171 (TestAtomicConditionalUpdate)
Covered Behavior: Concurrency for atomic single-statement update.
Test Coverage:
- ✅ 50 goroutines each DeductAtomic(5,1).
- ✅ All return nil (checked via t.Errorf inside goroutine).
- ✅ Final stock = 50.
- ❌ Does not test insufficient stock edge case.
- ❌ Does not test invalid quantity or zero quantity.
- ❌ Does not test product not found.
Assessment: PARTIAL
Missing: Insufficient stock, invalid quantity, zero quantity tests.

## Finding 7

Location: tests/locking_test.go (general)
Covered Behavior: Missing test cases across all strategies.
Missing Coverage:
- ❌ No tests for ErrInsufficientStock on naive, optimistic direct, optimistic with retry, atomic strategies.
- ❌ No tests for ErrInvalidQuantity (qty <= 0) on any strategy.
- ❌ No tests for ErrNotFound (product not found) on any strategy.
- ❌ No tests for boundary condition: qty == stock (should succeed, stock -> 0).
- ❌ No tests for retry exhaustion: DeductOptimisticWithRetry with maxRetries=0 and contention should return ErrOptimisticLock.
- ❌ No tests verifying stock unchanged on failure (except indirectly via invariants in some tests).
Assessment: PARTIAL
Overall: Tests demonstrate core concurrency behavior but lack comprehensive error and edge case coverage.

## Finding 8

Location: tests/locking_test.go:94-106 (TestOptimisticLockingConflict conflict counting)
Covered Behavior: Conflict counting logic.
Test Coverage:
- ✅ Uses mutex to protect conflictCount and successCount increments.
- ✅ Increments conflictCount on ErrOptimisticLock.
- ❌ The assertion `if conflictCount == 0 { t.Fatalf(...) }` is flaky; depends on timing to generate at least one conflict. With 20 goroutines and 50μs sleep, highly likely but not guaranteed. Could spuriously fail in slow environments.
Assessment: WARNING
Severity: MEDIUM
Notes: Test may be flaky under heavy load or slow CI. Better to run enough goroutines to make conflict virtually certain, or remove the assertion and rely on the invariant which already proves correctness.