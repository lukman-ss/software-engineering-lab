# Test Audit

## Coverage Analysis

1. **NOT NULL Constraints (`TestNotNullConstraints`)**
   - Asserts missing Email on user registration returns non-nil error.
   - Asserts missing Username on user registration returns non-nil error.
   - Asserts missing UserID on order creation returns non-nil error.
   - Status: PASS.

2. **CHECK Constraints (`TestCheckConstraints`)**
   - Asserts `age < 18` is rejected.
   - Asserts invalid status string is rejected.
   - Asserts `total_cents <= 0` order is rejected.
   - Status: PASS.

3. **UNIQUE Constraints (`TestUniqueConstraint`)**
   - Asserts duplicate email registration fails on second attempt.
   - Status: PASS.

4. **FOREIGN KEY Constraints (`TestForeignKeyConstraint`)**
   - Asserts order creation for non-existent `UserID=99999` fails.
   - Asserts valid `UserID` succeeds.
   - Status: PASS.

5. **Partial Unique Index (`TestPartialUniqueIndex`)**
   - Asserts duplicate active email fails.
   - Asserts soft-deleting initial user allows subsequent active user with same email.
   - Asserts duplicate second active user fails.
   - Asserts inserting soft-deleted user directly bypasses active partial index constraint.
   - Status: PASS.

6. **Concurrency Race Safety (`TestConcurrentRegistration_Safe_EnforcesUniqueness`)**
   - Launches 20 concurrent goroutines trying to insert identical email via `SafeStore`.
   - Asserts exactly 1 registration succeeds and 19 fail.
   - Asserts database count is exactly 1.
   - Status: PASS.

7. **Concurrency Race Vulnerability (`TestConcurrentRegistration_Unsafe_SuffersRaceCondition`)**
   - Launches 20 concurrent goroutines trying to insert identical email via `UnsafeStore`.
   - Asserts race condition causes `count > 1` duplicate insertions.
   - Status: PASS.

8. **Error Classification (`TestErrorClassification`)**
   - Validates SQLSTATE code mappings for `23502`, `23503`, `23505`, and `23514`.
   - Status: PASS.

## Execution Output

`go test -v ./...`: PASS (8 tests executed, 0 failures)
`go test -count=1 -race ./...`: PASS (0 data races detected)
