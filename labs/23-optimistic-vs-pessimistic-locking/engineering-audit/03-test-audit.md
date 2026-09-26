# Test Audit

Target Lab: `labs/23-optimistic-vs-pessimistic-locking`

## Execution Verification

### 1. Standard Unit & Concurrency Tests
Command:
```bash
go test -v -count=1 ./...
```
Result: PASS
```text
=== RUN   TestNaiveLostUpdate
    locking_test.go:37: Lost Update Demonstrated: 50 deduct calls occurred, but final stock is 98 (expected 50 under proper locking)
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
--- PASS: TestOptimisticLockingWithRetry (0.11s)
=== RUN   TestAtomicConditionalUpdate
--- PASS: TestAtomicConditionalUpdate (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/tests	0.114s
```

### 2. Race Detector Execution
Command:
```bash
go test -race -count=1 ./...
```
Result: PASS (Zero race conditions detected)
```text
ok  	github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/tests	1.127s
```

### 3. CLI Demo Execution
Command:
```bash
go run ./cmd/demo
```
Result: PASS
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
    Total Attempted Conflicts Retried: 42
    Actual Final Stock:   80 (SUCCESS - All retries eventually converged)
    Elapsed Time:         28.396375ms

[5] Atomic Single-Statement Operation (UPDATE ... WHERE stock >= qty):
    Initial Stock: 100
    Actual Final Stock:   50 (SUCCESS - Lockless Single Statement)

==========================================================
  Lab Execution Completed Successfully                    
==========================================================
```

## Coverage & Invariant Matrix

| Test Name | Focus Path | Invariant Checked | Result |
|---|---|---|---|
| `TestNaiveLostUpdate` | Failure / Anomaly | Final stock > 50 despite 50 deducts | PASS |
| `TestPessimisticLocking` | Happy Path / Concurrency | Exact stock (50) under 50 goroutines | PASS |
| `TestPessimisticLockingInsufficientStock` | Negative / Edge Case | Stock boundary error when stock < qty | PASS |
| `TestOptimisticLockingConflict` | Conflict Detection | Successes + Stock == 100; Conflicts > 0 | PASS |
| `TestOptimisticLockingWithRetry` | Recovery & Convergence | Stock == 100 - Optimistically applied | PASS |
| `TestAtomicConditionalUpdate` | Happy Path / Concurrency | Stock exactly matches 100 - 50 = 50 | PASS |
