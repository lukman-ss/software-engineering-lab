# Test Audit

## Test File
tests/locking_test.go (6 test functions)

## Executed Results

### go test ./...
Result: PASS
```
ok  github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/tests  (time recorded)
```
All 6 tests pass. No failures.

### go test -race ./...
Result: PASS
```
ok  github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/tests  (cache hit)
```
Race detector reported zero warnings.

### go run ./cmd/demo
Result: PASS (runs and terminates successfully)
Output recorded in 02-code-audit.md Finding 8.

## Test Coverage Analysis

### Test 1: TestNaiveLostUpdate — happy/positive path for anomaly
- 50 goroutines deduct 1 from stock 100.
- Asserts final stock != 50 (anomaly demonstrated).
- PASS. Demonstrates lost update.
- Coverage: happy path for naive anomaly. PASS.

### Test 2: TestPessimisticLocking — invariant preservation
- 50 goroutines deduct 1 from stock 100.
- Asserts final stock == 50.
- PASS. Invariant maintained.
- Coverage: happy path + invariant. PASS.

### Test 3: TestPessimisticLockingInsufficientStock — failure path
- Stock seeded to 2. First deduct of 2 succeeds. Second deduct of 1 returns ErrInsufficientStock.
- PASS. Correct error handling.
- Coverage: negative/edge case (stock boundary). PASS.

### Test 4: TestOptimisticLockingConflict — conflict detection
- 20 goroutines deduct 1 from stock 100.
- Asserts no data corruption: successCount + finalStock == 100.
- Asserts conflictCount > 0 (at least one conflict detected).
- PASS. Conflict detection and invariant hold.
- Coverage: concurrency + conflict detection + invariant. PASS.

### Test 5: TestOptimisticLockingWithRetry — retry convergence
- 20 goroutines deduct 1 from stock 100 with maxRetries=10.
- Asserts final stock == 100 - Optimistically counter.
- PASS. Retry converges to consistent state.
- Coverage: retry logic + invariant. PASS.

### Test 6: TestAtomicConditionalUpdate — atomicity
- 50 goroutines deduct 1 from stock 100 via AtomicDeduct.
- Asserts final stock == 50.
- PASS. Atomic operations serialize safely.
- Coverage: happy path + invariant. PASS.

## Coverage Summary
- Happy path: all 5 strategies (naive, pessimistic, optimistic-direct, optimistic-retry, atomic). PASS.
- Failure path: insufficient stock for pessimistic. PASS (single failure path; optimistic and atomic error paths not tested).
- Edge cases: stock boundary tested once (Test 3). Quantity <= 0 (ErrInvalidQuantity) not tested. Product not found (ErrNotFound) not tested.
- Transitions: optimistic conflict-to-retry not explicitly tested as state transition beyond aggregate invariant.
- Recovery/rollback: not applicable (no rollback needed on errors).
- Concurrency: 2 of 6 tests use high concurrency (50 goroutines). All pass under -race. PASS.
- Negative cases: only one (insufficient stock).

## Gaps
- Missing test: ErrInvalidQuantity for each strategy (input validation).
- Missing test: ErrNotFound for missing product.
- Missing test: optimistic retry exhaustion (maxRetries reached returns ErrOptimisticLock).
- Missing test: atomic insufficient stock failure path under concurrent overdraft.
- No table-driven test structure; tests are ad-hoc per scenario.

Assessment: Test suite is adequate to prove core claims for the three locking strategies. Race detector passes. Edge cases for invalid inputs are uncovered. No fabricated or hardcoded results.
