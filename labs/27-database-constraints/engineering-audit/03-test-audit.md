# Test Audit

## Executed Commands & Results

### Standard Tests
```bash
go test -v ./...
```
Output:
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
ok  	github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store	0.091s
```

### Race Detector
```bash
go test -race ./...
```
Output:
```text
ok  	github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store	(clean exit, no race conditions detected)
```

### Demo Execution
```bash
go run ./cmd/demo
```
Output:
```text
=================================================================
LAB 27: DATABASE CONSTRAINTS & DATA INTEGRITY DEMONSTRATION
=================================================================

[1] DEMONSTRATING NOT NULL CONSTRAINT (SQLSTATE 23502)
Attempt insert with missing email -> Error: invalid input: mandatory field is missing (rule: users_email_not_null)

[2] DEMONSTRATING CHECK CONSTRAINT (SQLSTATE 23514)
Attempt insert with age=15 (CHECK age >= 18) -> Error: validation failed: value outside permissible boundary (rule: users_age_check)
Attempt insert with invalid status -> Error: validation failed: value outside permissible boundary (rule: users_status_check)

[3] DEMONSTRATING FOREIGN KEY CONSTRAINT (SQLSTATE 23503)
Attempt insert order for non-existent UserID=9999 -> Error: reference error: referenced entity does not exist (rule: fk_orders_user)

[4] DEMONSTRATING PARTIAL UNIQUE INDEX (WHERE deleted_at IS NULL)
Created active user ID=1 (alice@company.com)
Duplicate active user insert rejected -> Error: conflict: resource with this unique attribute already exists (rule: users_active_email_idx)
Soft-deleted user ID=1 (deleted_at set)
New active user re-using email after soft delete -> Created user ID=2 (alice@company.com)

[5] CONCURRENCY STRESS TEST: 50 CONCURRENT REGISTRATIONS FOR SAME EMAIL
Results:
  - Total Goroutines: 50
  - Successful Registrations: 1
  - Rejected with UNIQUE VIOLATION (23505): 49
  - Database Integrity Intact: true

SQLSTATE Taxonomy Verification: code=23505 isUniqueViolation=true
```

## Test Coverage Evaluation
- Happy Path: Verified for all constraint types and successful lifecycle insertions.
- Failure / Negative Path: Verified for NOT NULL, CHECK boundary (age, status, total), UNIQUE, and FK.
- Edge Cases: Partial unique index covers active duplicate rejection, soft delete transition, reuse after soft delete, and multiple soft-deleted rows.
- Concurrency: Verified safe store under 20 concurrent goroutines (and demo with 50 goroutines), and proved race condition in unsafe store.
- Error Mapping: Verified SQLSTATE predicate functions across all constraint codes.
