# Test Audit

Target Lab: labs/23-optimistic-vs-pessimistic-locking

## Coverage Breakdown

| Test Name | File & Lines | Type / Scenario | Assertions & Invariants Checked | Result |
|---|---|---|---|---|
| `TestNaiveLostUpdate` | `tests/locking_test.go:10-38` | Concurrency / Negative Anomaly | 50 parallel goroutines; asserts `p.Stock != 50` proving lost update anomaly | PASS |
| `TestPessimisticLocking` | `tests/locking_test.go:40-67` | Concurrency / Happy Path | 50 parallel goroutines; asserts `p.Stock == 50` proving full serialization | PASS |
| `TestPessimisticLockingInsufficientStock` | `tests/locking_test.go:69-83` | Boundary / Failure Path | Initial stock 2, deduct 2 then 1; asserts `ErrInsufficientStock` returned | PASS |
| `TestOptimisticLockingConflict` | `tests/locking_test.go:85-121` | Concurrency / Conflict Path | 20 parallel goroutines without retry; asserts `conflictCount > 0` and stock invariant `successCount + stock == 100` | PASS |
| `TestOptimisticLockingWithRetry` | `tests/locking_test.go:123-145` | Concurrency / Convergence | 20 parallel goroutines with retry; asserts `p.Stock == 100 - Optimistically` | PASS |
| `TestAtomicConditionalUpdate` | `tests/locking_test.go:147-171` | Concurrency / Happy Path | 50 parallel goroutines; asserts `p.Stock == 50` | PASS |

## Execution Proof

### 1. `go test -v -count=1 ./...`
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
--- PASS: TestOptimisticLockingWithRetry (0.04s)
=== RUN   TestAtomicConditionalUpdate
--- PASS: TestAtomicConditionalUpdate (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/tests	0.369s
```

### 2. `go test -race -count=1 ./...`
```text
ok  	github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/tests	1.355s
```
Status: Zero data races detected.

### 3. `go run ./cmd/demo`
```text
==========================================================
  Optimistic vs Pessimistic Locking & Atomic Operations   
==========================================================

[1] Naive Read-Modify-Write (50 concurrent requests):
    Initial Stock: 100
    Expected Final Stock: 50
    Actual Final Stock:   99 (LOST UPDATE DETECTED!)

[2] Pessimistic Locking (SELECT ... FOR UPDATE):
    Initial Stock: 100
    Actual Final Stock:   50 (SUCCESS - Fully Synchronized)

[3] Optimistic Locking Direct (20 concurrent requests, no retry):
    Initial Stock: 100
    Successful Deductions: 1
    Rejected Conflicts:   19
    Actual Final Stock:   99 (SUCCESS - State Guarded, Zero Corruption)

[4] Optimistic Locking With Exponential Backoff Retry (20 requests):
    Initial Stock: 100
    Successful Deductions: 20
    Total Attempted Conflicts Retried: 40
    Actual Final Stock:   80 (SUCCESS - All retries eventually converged)
    Elapsed Time:         21.490208ms

[5] Atomic Single-Statement Operation (UPDATE ... WHERE stock >= qty):
    Initial Stock: 100
    Actual Final Stock:   50 (SUCCESS - Lockless Single Statement)

==========================================================
  Lab Execution Completed Successfully                    
==========================================================
```
Status: Executed successfully with real runtime metrics.
