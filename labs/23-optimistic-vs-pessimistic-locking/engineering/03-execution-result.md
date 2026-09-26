# Execution Result

## Build
Command:
```bash
go build ./...
```
Result:
```text
PASS (exit code 0)
```

## Tests
Command:
```bash
go test -v ./...
```
Result:
```text
?   	github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/internal/inventory	[no test files]
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
--- PASS: TestOptimisticLockingWithRetry (0.05s)
=== RUN   TestAtomicConditionalUpdate
--- PASS: TestAtomicConditionalUpdate (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/tests	0.386s
```

## Race Detector
Command:
```bash
go test -race ./...
```
Result:
```text
?   	github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/internal/inventory	[no test files]
ok  	github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/tests	2.510s
```

## Demo
Command:
```bash
go run ./cmd/demo
```
Result:
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
    Total Attempted Conflicts Retried: 61
    Actual Final Stock:   80 (SUCCESS - All retries eventually converged)
    Elapsed Time:         49.859125ms

[5] Atomic Single-Statement Operation (UPDATE ... WHERE stock >= qty):
    Initial Stock: 100
    Actual Final Stock:   50 (SUCCESS - Lockless Single Statement)

==========================================================
  Lab Execution Completed Successfully                    
==========================================================
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
