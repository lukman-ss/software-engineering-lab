# Test Audit

Lab: 27-database-constraints  
Test file: `internal/store/store_test.go`  
Tested functions: constraint checks, concurrency, error mapping.

## Summary
- 8 unit tests covering all claimed behaviors.
- All pass (`go test -v ./...`).
- Race detector clean (`go test -race ./...`).
- Test `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` is timing-dependent but reliably demonstrates the intended race.

## Happypass Coverage

| Test | Constraint / Path | Happy Path | Notes |
|------|-------------------|------------|-------|
| `TestNotNullConstraints` | NOT NULL (email, username, user_id) | ❌ intentional failure (missing fields) | Confirms rejection works |
| `TestCheckConstraints` | CHECK (age >= 18, status enum, total_cents > 0) | ❌ intentional failure | Boundary/ enum checks correct |
| `TestUniqueConstraint` | UNIQUE (full email index) | First insert succeeds, second fails | Atomic check+insert inside engine lock |
| `TestForeignKeyConstraint` | FOREIGN KEY (user_id → users) | Valid FK succeeds, invalid fails | Lookup+insert atomic |
| `TestPartialUniqueIndex` | PARTIAL UNIQUE (WHERE deleted_at IS NULL) | Active insert succeeds, re-active fails; soft-deleted insert succeeds; post-soft-delete active succeeds; second re-active fails | Validates conditional uniqueness flow |
| `TestConcurrentRegistration_Safe_EnforcesUniqueness` | UNIQUE under 50 concurrent goroutines | Exactly 1 success, 49 unique violations | Uses engine's write-lock + unique check |
| `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` | Unsafe app-check (read-then-write) | >1 rows inserted (race win) | Demonstrates TOCTOU |
| `TestErrorClassification` | `IsConstraintViolation` for each SQLSTATE | Returns true for wrapped `ConstraintError` | Validates helpers |

All tests follow the pattern: arrange (engine + store), act (one or more calls), assert (error nil/not-nil, count, specific messages where relevant).

## Failure Path Coverage

- Every constraint type has at least one test exercising the rejection (`err != nil`).
- Concurrency unsafe test asserts that race condition leads to duplicate insert (count > 1).
- Concurrency safe test asserts exactly one winner and N-1 errors.
- Error classification test ensures unwrapping and state detection works.

## Edge Cases / Boundary Conditions

- **NOT NULL**: empty string caught (as in demo). Null pointer not relevant (Go strings). Zero-valued `int64` (user_id) caught.
- **CHECK**:  
  - `age == 18` accepted (`u.Age < 18` is strict; good).  
  - `age == 17` rejected.  
  - Status enum: exact strings `"active"`, `"suspended"`, `"pending"` accepted; any other string rejected.  
  - `total_cents == 0` rejected (`<= 0` check).  
- **UNIQUE**: duplicate detection is case-sensitive (maps use Go string equality).  
- **PARTIAL UNIQUE**:  
  - Soft-deleted row (`DeletedAt != nil`) does NOT block a new active insert with same email.  
  - Soft-deleted row does NOT block another soft-deleted insert with same email (index ignores `DeletedAt != nil`).  
  - Active row blocks another active insert.  
  - Active row blocks a soft-delete+re-insert? No: soft delete removes from index first, then re-insert active allowed (see test step 4-5).  
- **FOREIGN KEY**: zero `user_id` caught by NOT NULL first; non-existent positive ID caught by lookup.  
- **Concurrency**:  
  - Safe: 20-50 goroutine runs consistently produce 1 success (verified x5 under `-race`).  
  - Unsafe: 20-50 goroutine runs consistently produce count > 1 under the artificial 1ms sleep. Without sleep, flaky — but the test includes the sleep deliberately.

## Negative Cases

All negative cases above covered. No constraint type is missing a rejection test.

## Transitions

Only relevant transition is soft delete → re-register active. Covered by `TestPartialUniqueIndex` steps 3→4→5.

## Recovery / Rollback

No multi-statement transactions modeled; each engine call is atomic. No rollback logic to test.

## Concurrency Safety (Test Perspective)

- The two concurrency tests together validate that:  
  1. The unsafe path *can* lose (and does, reliably with the sleep).  
  2. The safe path *never* loses (deterministic win + N-1 errors).  
- No test for mixed safe/unsafe concurrency; out of scope.

## Test Quality Notes

- No table-driven tests; each constraint has a dedicated function. Acceptable given small number.
- Helper functions not extracted (e.g., makeUser) — minor duplication but clear.
- `TestErrorClassification` only tests the `New*` constructors and `IsConstraintViolation`. It does **not** test `MapToDomainError`. This is a minor gap (the mapper is used by `SafeStore` but not asserted in any test).
- `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` includes the artificial `time.Sleep` which makes the test less pure but more reliably demonstrates the intended lesson. Acceptable as a teaching device.
- No test for the **cross-index path** (Finding 11 in code audit). This is a missing edge case: mixing `RegisterUser` and `RegisterUserPartial` for the same email could violate uniqueness silently.

## Verdict on Tests

- **Strengths**: covers every claimed constraint behavior, both concurrency paths, error mapping skeleton. Deterministic where it matters (safe concurrency). No flaky assertions (except the unsafe test which is *supposed* to be non-deterministic but biased toward failure by the sleep).
- **Weaknesses**: no direct assertion on `MapToDomainError` output; no test for mixed registration paths; unsafe test leans on a timing crutch (but that crutch is in the production code too, so it tests exactly what ships).
- Overall: tests prove the claimed behavior exists and is reproducible under the harness. No evidence of inflated or fake results.