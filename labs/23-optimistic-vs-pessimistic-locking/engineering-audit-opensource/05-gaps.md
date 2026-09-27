# Gap Analysis

## Gaps Identified

### 1. MISSING_TEST: ErrInvalidQuantity not tested
- The service validates `qty <= 0` in all four strategies but no test exercises this error path.
- Severity: LOW
- Impact: Input validation correctness is assumed but unproven.

### 2. MISSING_TEST: ErrNotFound not tested
- Store returns ErrNotFound for unknown product ID, but no test verifies this.
- Severity: LOW
- Impact: Error propagation for missing entities unverified.

### 3. MISSING_TEST: Optimistic retry exhaustion
- DeductOptimisticWithRetry returns ErrOptimisticLock when maxRetries exhausted, but no test forces this condition.
- Severity: MEDIUM
- Impact: Retry boundary behavior unproven; could silently loop or panic under high contention.

### 4. MISSING_TEST: Atomic insufficient stock under concurrency
- AtomicDeduct returns ErrInsufficientStock when stock < qty, but no test verifies this under concurrent overdraft (e.g., 200 goroutines deducting 1 from stock 100).
- Severity: MEDIUM
- Impact: Atomic operation correctness only proven for non-overdraft case.

### 5. MISSING_TEST: Pessimistic overdraft under concurrency
- Test 3 only tests serial overdraft. Concurrent overdraft (50 goroutines, stock 2) could expose race in the stock check. (Likely safe due to row lock, but untested.)
- Severity: LOW
- Impact: Pessimistic overdraft invariant under concurrency unproven.

### 6. DOC_CODE_MISMATCH: SQL terminology vs in-memory simulation
- README references `SELECT ... FOR UPDATE` and `UPDATE ... SET ... WHERE stock >= N` SQL syntax.
- Implementation uses Go `sync.Mutex` and in-memory maps.
- Design doc acknowledges this as a deliberate scoping decision.
- Severity: LOW
- Impact: Reader may believe real database is used. Mitigated by engineering design doc transparency.

## No Critical Gaps
- No RACE_CONDITION (race detector passed).
- No UNHANDLED_ERROR (errors propagated consistently).
- No BROKEN_IMPLEMENTATION (all strategies execute correctly per tests and demo).
- No FAKE_DEMO (demo output verified by execution).
- No FAKE_BENCHMARK (no benchmarks present).
- No UNVERIFIED_RESULT (all results reproduced at runtime).
- No IMPLEMENTATION_OVERCLAIM (implementation scope matches design doc).
- No RESEARCH_MISMATCH (implementation aligns with research claims).

## Summary
Test suite proves core behavior (three locking strategies + naive baseline) under concurrency with race detector. Gap is in edge-case and error-path coverage (4 of 6 strategies lack error-path tests for invalid inputs and insufficient stock). These are test gaps, not implementation bugs.
