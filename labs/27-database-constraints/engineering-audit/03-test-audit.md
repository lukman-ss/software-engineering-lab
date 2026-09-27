# Test Audit

Target Lab: labs/27-database-constraints

## Coverage Summary

1. **NOT NULL Constraints (`TestNotNullConstraints`)**
   - Covered: Missing email on User, missing username on User, missing user_id on Order.
   - Assessment: PASS

2. **CHECK Constraints (`TestCheckConstraints`)**
   - Covered: Underage user (`age < 18`), invalid status (`banned`), order total zero/negative (`total_cents <= 0`).
   - Assessment: PASS

3. **UNIQUE Constraints (`TestUniqueConstraint`)**
   - Covered: First user insert succeeds, duplicate email insert fails with constraint violation.
   - Assessment: PASS

4. **FOREIGN KEY Constraints (`TestForeignKeyConstraint`)**
   - Covered: Non-existent `user_id` insert rejected, existing `user_id` insert succeeds.
   - Assessment: PASS

5. **PARTIAL UNIQUE INDEX (`TestPartialUniqueIndex`)**
   - Covered: Active insert succeeds -> Duplicate active insert rejected -> Soft delete active user -> Re-insert active user succeeds -> Duplicate active insert rejected again -> Soft-deleted direct insert allowed.
   - Assessment: PASS

6. **Concurrency Stress Test — Safe Store (`TestConcurrentRegistration_Safe_EnforcesUniqueness`)**
   - Covered: 20 concurrent goroutines attempting to register the exact same email.
   - Verified: Exactly 1 success, 19 rejected with unique violation, database contains exactly 1 user row.
   - Assessment: PASS

7. **Concurrency Stress Test — Unsafe Store (`TestConcurrentRegistration_Unsafe_SuffersRaceCondition`)**
   - Covered: 20 concurrent goroutines attempting to register the exact same email using application-level check.
   - Verified: Proves race condition occurrence by resulting in >1 duplicate user entries.
   - Assessment: PASS

8. **Error Classification (`TestErrorClassification`)**
   - Covered: Checks `dberr.IsConstraintViolation` against SQLSTATE constants (`23502`, `23503`, `23505`, `23514`).
   - Assessment: PASS

## Execution Results

Command: `go test -count=1 -v ./...`
Status: PASS
Output: 8 passed in internal/store package.

Command: `go test -count=1 -race ./...`
Status: PASS (No data races detected).

Command: `go run ./cmd/demo`
Status: PASS (Exhibits expected demo output for all 5 constraint categories).
