# Test Audit

Target Lab: labs/23-optimistic-vs-pessimistic-locking

## Test Coverage Overview

File: `tests/locking_test.go`
Scenarios covered:
1. `TestNaiveLostUpdate`: Happy/anomaly path. 50 parallel goroutines execute uncoordinated read-modify-write. Asserts final stock != 50 (lost update reproduced).
2. `TestPessimisticLocking`: Happy path concurrency. 50 goroutines deduct 1 stock each. Asserts exact final stock == 50.
3. `TestPessimisticLockingInsufficientStock`: Negative path. Deducts initial stock, validates subsequent attempt returns `ErrInsufficientStock`.
4. `TestOptimisticLockingConflict`: Direct conflict path. 20 concurrent goroutines without retry. Verifies at least 1 conflict rejected and exact stock invariant `successCount + stock == 100`.
5. `TestOptimisticLockingWithRetry`: Convergence path. 20 goroutines with backoff retry. Verifies all updates eventually succeed and state remains consistent.
6. `TestAtomicConditionalUpdate`: Atomic path. 50 goroutines deduct atomically. Asserts exact final stock == 50.

## Execution Verification

### `go test -v ./...` Output
```text
=== RUN   TestNaiveLostUpdate
    locking_test.go:37: Lost Update Demonstrated: 50 deduct calls occurred, but final stock is 97 (expected 50 under proper locking)
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
ok  	github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/tests	0.061s
```

### `go test -race ./...` Output
```text
PASS
ok  	github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/tests	1.170s
```

## Assessment
- All 6 tests pass without race warnings.
- Both positive, negative, and concurrent conflict paths are asserted.
- Tests prove claimed behavior.
