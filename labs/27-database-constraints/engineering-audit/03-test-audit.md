# Test Audit Report

## Test Coverage Assessment

### 1. Happy Path & Constraint Enforcements
- `TestNotNullConstraints`:
  - Missing email -> Rejection with NOT NULL constraint error (PASS).
  - Missing username -> Rejection with NOT NULL constraint error (PASS).
  - Missing order user_id -> Rejection with NOT NULL constraint error (PASS).
- `TestCheckConstraints`:
  - Underage (`age < 18`) -> Rejection with CHECK violation (PASS).
  - Invalid enum status (`status = "banned"`) -> Rejection with CHECK violation (PASS).
  - Non-positive order amount (`total_cents <= 0`) -> Rejection with CHECK violation (PASS).
- `TestUniqueConstraint`:
  - First registration succeeds (PASS).
  - Subsequent registration with same email fails with UNIQUE constraint violation (PASS).
- `TestForeignKeyConstraint`:
  - Non-existent parent UserID (`99999`) -> Rejection with FK violation (PASS).
  - Valid existing UserID -> Success with generated order ID (PASS).
- `TestPartialUniqueIndex`:
  - Insert first active user -> Success (PASS).
  - Insert duplicate active user -> Rejected (PASS).
  - Soft delete first active user -> Success (PASS).
  - Insert new active user with same email -> Success (PASS).
  - Insert duplicate active user again -> Rejected (PASS).
  - Insert explicitly soft-deleted user with same email -> Success (PASS).

### 2. Concurrency & Race Condition Verification
- `TestConcurrentRegistration_Safe_EnforcesUniqueness`:
  - 20 concurrent goroutines attempt to register identical email.
  - Result: Exactly 1 success, 19 rejected with unique constraint violation, store contains 1 record.
  - Race Detector: Clean pass without race warnings.
- `TestConcurrentRegistration_Unsafe_SuffersRaceCondition`:
  - 20 concurrent goroutines attempt to register identical email using UnsafeStore.
  - Result: Multiple records (>1) inserted due to check-then-act race condition, demonstrating failure of app-level checks.

### 3. Error Taxonomy & Mapping
- `TestErrorClassification`:
  - Tests `IsConstraintViolation` helper across all supported SQLSTATE codes (`23502`, `23505`, `23514`, `23503`).

## Execution Log

```text
=== RUN   TestNotNullConstraints
--- PASS: TestNotNullConstraints (0.00s)
=== RUN   TestCheckConstraints
--- PASS: TestCheckConstraints (0.00s)
=== RUN   TestUniqueConstraint
--- PASS: TestUniqueConstraint (0.00s)
=== RUN   TestForeignKeyConstraint
--- PASS: TestForeignKeyConstraint (0.00s)
=== RUN   TestPartialUniqueIndex
--- PASS: TestPartialUniqueIndex (0.00s)
=== RUN   TestConcurrentRegistration_Safe_EnforcesUniqueness
--- PASS: TestConcurrentRegistration_Safe_EnforcesUniqueness (0.00s)
=== RUN   TestConcurrentRegistration_Unsafe_SuffersRaceCondition
--- PASS: TestConcurrentRegistration_Unsafe_SuffersRaceCondition (0.00s)
=== RUN   TestErrorClassification
--- PASS: TestErrorClassification (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store	0.306s
```

Race detector test:
```text
ok  	github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store	1.318s
```
Result: 0 race detected.
