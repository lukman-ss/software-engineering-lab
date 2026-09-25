# Test Audit

Target Lab: labs/18-deadlock

## Test Suite Execution

Command: `go test -count=1 -v ./...`
Status: PASS
Duration: 0.463s

Command: `go test -count=1 -race -v ./...`
Status: PASS
Duration: 1.456s

## Coverage Analysis

### 1. Happy Path
- `TestLockOrderingPreventsDeadlock`: PASS
- Validates successful concurrent transfers when lock ordering is maintained. Balance invariants hold.

### 2. Failure Path / Abort Mechanism
- `TestDeadlockOccurrence`: PASS
- Validates that naive circular transfers encounter `ErrDeadlock` due to deadlock timeout.

### 3. Recovery / Retry
- `TestRetryRecoversDeadlock`: PASS
- Validates that retry loops resolve temporary deadlocks under contention.

### 4. Empirical Correlation
- `TestTransactionDurationImpact`: PASS
- Compares deadlock rates across 10 iterations between 5ms delay vs 0ms delay. Confirms longer lock hold times increase deadlock frequency.

### 5. Edge Cases & Negative Paths
- Context pre-cancellation: NOT COVERED
- Insufficient balance: NOT COVERED
- Self-transfer (`accA == accB`): NOT COVERED
- Assessment: WARNING (Non-blocking: lab scope focuses on concurrency and lock dynamics rather than full banking business logic).
