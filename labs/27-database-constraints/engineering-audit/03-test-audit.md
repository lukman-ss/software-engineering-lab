# Test Audit

Target Lab: labs/27-database-constraints

## Test Suite Overview

Test file: `internal/store/store_test.go`
Tests present:
1. `TestNotNullConstraints`: Verifies missing `Email`, missing `Username`, and missing `UserID` in Order return NOT NULL errors.
2. `TestCheckConstraints`: Verifies underage (`Age < 18`), invalid `Status` string, and invalid `TotalCents <= 0` return CHECK violation errors.
3. `TestUniqueConstraint`: Verifies sequential duplicate insert returns UNIQUE violation error.
4. `TestForeignKeyConstraint`: Verifies non-existent `UserID` insertion returns FOREIGN KEY error, and valid `UserID` succeeds.
5. `TestPartialUniqueIndex`: Verifies single active user constraint, rejection of duplicate active user, soft deletion behavior, re-registration of soft-deleted email, and coexistence of multiple inactive records.
6. `TestConcurrentRegistration_Safe_EnforcesUniqueness`: Spawns 20 goroutines attempting registration with identical email against `SafeStore`; verifies exactly 1 succeeds, 19 fail with constraint violation, and table row count is 1.
7. `TestErrorClassification`: Verifies `dberr.IsConstraintViolation` accurately matches `SQLStateNotNullViolation`, `SQLStateUniqueViolation`, `SQLStateCheckViolation`, and `SQLStateForeignKeyViolation`.

## Test Execution Results

Command:
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
=== RUN   TestErrorClassification
--- PASS: TestErrorClassification (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store	0.102s
```

Command:
```bash
go test -race ./...
```
Output:
```text
ok  	github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store	1.120s
```

## Coverage & Gap Assessment

- Happy path coverage: PASS
- Failure path coverage: PASS (NOT NULL, CHECK, UNIQUE, FOREIGN KEY, PARTIAL INDEX)
- Concurrency race detector: PASS (0 data races detected)
- Unsafe store test: In `engineering/01-design.md`, a test named `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` was planned but not implemented in `store_test.go`.
