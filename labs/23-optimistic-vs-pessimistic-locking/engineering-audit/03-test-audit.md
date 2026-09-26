# Test Audit

## Test Suite Overview

File: `tests/locking_test.go`
Package: `tests`
Execution Tool: `go test` and `go test -race`

## Test Cases Evaluated

### 1. `TestNaiveLostUpdate`
- Target: `NaiveDeduct`
- Scenario: 50 concurrent goroutines deducting 1 unit each from initial stock of 100.
- Assertion: Verifies that final stock is NOT 50 due to lost updates (typically 98-99).
- Verdict: PASS
- Notes: Accurately proves vulnerability of uncoordinated read-modify-write patterns.

### 2. `TestPessimisticLocking`
- Target: `DeductPessimistic`
- Scenario: 50 concurrent goroutines deducting 1 unit each from initial stock of 100.
- Assertion: `p.Stock == 50` exactly.
- Verdict: PASS
- Notes: Validates absolute consistency under exclusive row locking.

### 3. `TestPessimisticLockingInsufficientStock`
- Target: `DeductPessimistic`
- Scenario: Product with stock 2; deduct 2 succeeds, subsequent deduct 1 fails with `ErrInsufficientStock`.
- Assertion: Checks error return and boundary enforcement.
- Verdict: PASS
- Notes: Validates edge case and rejection of overdrafts.

### 4. `TestOptimisticLockingConflict`
- Target: `DeductOptimisticDirect`
- Scenario: 20 concurrent goroutines attempting direct optimistic update without retries.
- Assertion: Invariant `successCount + stock == 100` holds, and `conflictCount > 0`.
- Verdict: PASS
- Notes: Proves version guard rejects stale concurrent writes while protecting data integrity.

### 5. `TestOptimisticLockingWithRetry`
- Target: `DeductOptimisticWithRetry`
- Scenario: 20 concurrent goroutines attempting optimistic update with exponential backoff retries.
- Assertion: `p.Stock == 100 - Optimistically` and all retried conflicts converge.
- Verdict: PASS
- Notes: Proves self-healing convergence of retry mechanism under contention.

### 6. `TestAtomicConditionalUpdate`
- Target: `DeductAtomic`
- Scenario: 50 concurrent goroutines executing atomic single-statement decrements.
- Assertion: `p.Stock == 50` exactly, error-free.
- Verdict: PASS
- Notes: Proves lockless conditional decrement semantics.

## Concurrency and Race Safety Execution Results

Command: `go test -v ./...`
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
--- PASS: TestOptimisticLockingWithRetry (0.04s)
=== RUN   TestAtomicConditionalUpdate
--- PASS: TestAtomicConditionalUpdate (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/tests	0.168s
```

Command: `go test -race ./...`
```text
ok  	github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/tests	1.176s
```
Race Detector Result: Zero warnings, zero data races detected.

## Demo Execution Result

Command: `go run ./cmd/demo`
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
    Elapsed Time:         41.633041ms

[5] Atomic Single-Statement Operation (UPDATE ... WHERE stock >= qty):
    Initial Stock: 100
    Actual Final Stock:   50 (SUCCESS - Lockless Single Statement)

==========================================================
  Lab Execution Completed Successfully                    
==========================================================
```
Demo Result: Real, reproducible, matches claimed behavior and test suite.
