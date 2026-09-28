# Test Audit

## Test Suite Execution

Commands:
```bash
go test -v ./...
go test -count=1 -race -v ./...
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
ok  	github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store	1.088s
```

## Coverage Analysis

1. **Happy Path**: Tested (valid user registration, valid order creation, partial index re-registration after soft delete).
2. **Failure Path**: Tested (empty fields for NOT NULL, out-of-range age and invalid status for CHECK, missing parent row for FOREIGN KEY, duplicate email for UNIQUE).
3. **Edge Cases**: Tested (soft-deleted row permitting reuse of email, while active duplicate email is rejected).
4. **Concurrency**: Tested (50 concurrent goroutines against `SafeStore` with 1 winner and 49 rejections; same concurrency against `UnsafeStore` demonstrating multi-write race anomaly).
5. **Race Detector**: Ran with `-race` flag, 0 data races detected.

## Assessment

Assessment: PASS
Severity: LOW
Notes: Test suite is robust, reproducible, and verifies all claimed integrity guarantees.
