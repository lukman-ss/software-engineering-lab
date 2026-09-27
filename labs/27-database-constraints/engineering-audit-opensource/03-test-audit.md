# Test Audit

## Commands Executed
- `go build ./...` → exit 0 (no output)
- `go vet ./...` → exit 0 (no output)
- `go test -v -count=1 ./...` → 8 tests PASS, total ok.
- `go test -race -count=1 ./...` → ok (no races), timing ~1.1s.
- `go run ./cmd/demo` → exit 0; output recorded (NOT NULL, CHECK, FOREIGN KEY, PARTIAL UNIQUE, 50-goroutine concurrency stress with 1 success + 49 rejections, SQLSTATE taxonomy verification).

Actual test output (from execution):

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
```

Actual race output:

```text
ok  	github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store	1.101s
```

Actual demo output (verified live, exit 0):

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

## Coverage Assessment

| Requirement | Test | Status |
|---|---|---|
| NOT NULL email empty | `TestNotNullConstraints` (RegisterUser missing email) | PASS — `err != nil` asserted |
| NOT NULL username empty | `TestNotNullConstraints` (RegisterUser missing username) | PASS — `err != nil` asserted |
| NOT NULL orders.user_id == 0 | `TestNotNullConstraints` (CreateOrder missing user_id) | PASS — `err != nil` asserted |
| CHECK age < 18 | `TestCheckConstraints` (age 16) | PASS — `err != nil` asserted |
| CHECK invalid status | `TestCheckConstraints` (status "banned") | PASS — `err != nil` asserted |
| CHECK total_cents == 0 | `TestCheckConstraints` (order total 0) | PASS — `err != nil` asserted |
| UNIQUE duplicate email | `TestUniqueConstraint` | PASS — second insert asserts `err != nil` |
| FOREIGN KEY missing parent | `TestForeignKeyConstraint` (user_id 99999) | PASS — `err != nil` asserted |
| FOREIGN KEY valid parent | `TestForeignKeyConstraint` (create user, then order) | PASS — asserts `o.ID != 0` |
| PARTIAL UNIQUE active dup rejected | `TestPartialUniqueIndex` step 2 | PASS |
| PARTIAL UNIQUE soft delete → re-insert allowed | `TestPartialUniqueIndex` steps 3-4 | PASS |
| PARTIAL UNIQUE second active after re-insert rejected | `TestPartialUniqueIndex` step 5 | PASS |
| PARTIAL UNIQUE soft-deleted direct insert bypasses index | `TestPartialUniqueIndex` step 6 | PASS |
| Concurrency safe exactly-1 | `TestConcurrentRegistration_Safe_EnforcesUniqueness` (20 goroutines) | PASS |
| Concurrency unsafe duplicates | `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` (20 goroutines) | PASS |
| Error classification 23502/05/14/03 | `TestErrorClassification` | PASS |

## Strengths
- Each claimed constraint type has at least one happy-path and one failure-path test.
- Concurrency coverage exists for both safe and unsafe paths (20 goroutines each; demo scales to 50).
- Partial index lifecycle fully exercised (insert → duplicate reject → soft delete → re-insert allow → second active reject → deleted insert bypass).

## Weaknesses / Missing Tests
1. **SQLSTATE not asserted in constraint tests**: all functional tests use `err == nil` / `err != nil` only. Only `TestErrorClassification` constructs errors directly and calls `IsConstraintViolation`. No test asserts that e.g. a duplicate email from `SafeStore.RegisterUser` carries SQLSTATE `23505` (or at least that `MapToDomainError` output contains the right domain prefix). This weakens proof that the engine returns the *claimed* codes end-to-end. Severity: MEDIUM — behavior proven but code mapping unproven.
2. **Boundary values untested**: `Age == 18` (should pass), `Age == 17` (should fail), `TotalCents == 1` vs `0` vs `-1`, empty status `""` (fails via CHECK, not NOT NULL). None explicitly asserted. Severity: LOW.
3. **No test for `SoftDeleteUser` failure paths**: non-existent ID, nil `DeletedAt`. Untested. Severity: LOW.
4. **No test for `InsertUserUnsafe` bypass semantics**: unsafe insert does not touch indexes; after unsafe duplicates, a subsequent `SafeStore.RegisterUser` with same email would behave based only on `emailIndex` (empty), allowing insert — untested, could confuse readers. Untested interaction. Severity: LOW.
5. **Unsafe concurrency test is scheduler-sensitive by design**: asserts `count > 1`. The 1ms sleep makes duplication near-certain, and race-clean run observed, but a pathological scheduler could theoretically serialize all 20 goroutines (check sees prior insert, rejects) yielding count == 1 and flaky failure. Not observed; acceptable for demo. Documented as WARNING not FAIL.
6. **`dberr` package has no dedicated test file**: classification covered only via `store_test`. Fine but package-level edge (unknown code, nil error, `MapToDomainError` default branch) untested. Severity: LOW.
