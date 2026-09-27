# Test Audit

Target Lab: labs/23-optimistic-vs-pessimistic-locking

## Test Coverage Overview

File: `tests/locking_test.go`

1. `TestNaiveLostUpdate`:
   - Spawns 50 concurrent goroutines deducting stock without synchronization.
   - Verifies that final stock is not 50 (demonstrates lost update anomaly).
   - Result: PASS

2. `TestPessimisticLocking`:
   - Spawns 50 concurrent goroutines deducting stock via pessimistic row lock.
   - Verifies that stock decrements strictly from 100 to 50 without error or lost update.
   - Result: PASS

3. `TestPessimisticLockingInsufficientStock`:
   - Negative test verifying that attempting to deduct more than available returns `ErrInsufficientStock`.
   - Result: PASS

4. `TestOptimisticLockingConflict`:
   - Spawns 20 concurrent goroutines attempting direct optimistic update without retry.
   - Verifies conflict generation (`conflictCount > 0`) and stock integrity invariant (`successCount + stock == 100`).
   - Result: PASS

5. `TestOptimisticLockingWithRetry`:
   - Spawns 20 concurrent goroutines attempting optimistic update with exponential backoff retry.
   - Verifies that retries converge and stock strictly equals initial stock minus successful deductions.
   - Result: PASS

6. `TestAtomicConditionalUpdate`:
   - Spawns 50 concurrent goroutines calling atomic single-statement updates.
   - Verifies accurate stock decrement from 100 to 50.
   - Result: PASS

## Concurrency and Race Detection

- Command: `go test -race -count=1 ./...`
- Status: PASS (0 race conditions detected)

## Execution Log

```text
=== RUN   TestNaiveLostUpdate
    locking_test.go:37: Lost Update Demonstrated: 50 deduct calls occurred, but final stock is 99 (expected 50 under proper locking)
--- PASS: TestNaiveLostUpdate (0.00s)
=== RUN   TestPessimisticLocking
--- PASS: TestPessimisticLocking (0.00s)
=== RUN   TestPessimisticLockingInsufficientStock
--- PASS: TestPessimisticLockingInsufficientStock (0.00s)
=== RUN   TestOptimisticLockingConflict
    locking_test.go:120: Optimistic Locking Conflict Demonstrated: 1 succeeded, 19 conflicts rejected
--- PASS: TestOptimisticLockingConflict (0.00s)
=== RUN   TestOptimisticLockingWithRetry
    locking_test.go:144: Optimistic Retry Successful: 20 total updates applied, retried conflicts resolved
--- PASS: TestOptimisticLockingWithRetry (0.03s)
=== RUN   TestAtomicConditionalUpdate
--- PASS: TestAtomicConditionalUpdate (0.00s)
PASS
```
