# Test Audit

Coverage of claims:
- NOT NULL: TestNotNullConstraints (3 subtests)
- CHECK: TestCheckConstraints (3 subtests: age, status, order total)
- UNIQUE: TestUniqueConstraint
- FOREIGN KEY: TestForeignKeyConstraint
- PARTIAL UNIQUE INDEX: TestPartialUniqueIndex
- CONCURRENCY (safe store enforces 1, 49 rejects): TestConcurrentRegistration_Safe_EnforcesUniqueness
- CONCURRENCY (unsafe suffers >1): TestConcurrentRegistration_Unsafe_SuffersRaceCondition
- ERROR CLASSIFICATION: TestErrorClassification

Assessment per test:
TestNotNullConstraints: missing email/username/user_id -> error != nil. PASS.
TestCheckConstraints: age<18 -> error != nil; status invalid -> error != nil; order total<=0 -> error != nil. PASS.
TestUniqueConstraint: insert duplicate email -> error != nil. PASS.
TestForeignKeyConstraint: order user_id=9999 -> error != nil; valid ref -> order ID>0. PASS.
TestPartialUniqueIndex:
1. Insert active user -> ok.
2. Duplicate active -> err != nil.
3. Soft-delete active -> no err.
4. Re-insert active same email -> ok.
5. Second active duplicate -> err != nil.
6. Insert soft-deleted same email (DeletedAt set) -> ok (bypasses partial index). PASS.
TestConcurrentRegistration_Safe_EnforcesUniqueness:
50 goroutines, same email, safe store -> exactly 1 success, 49 errors, engine.GetUsersCount()==1. PASS.
TestConcurrentRegistration_Unsafe_SuffersRaceCondition:
50 goroutines, same email, unsafe store -> GetUsersCount()>1 (i.e., duplicate created). PASS.
TestErrorClassification:
Construct each *Violation and test IsConstraintViolation(err, respective SQLState). PASS.

Weaknesses/coverage gaps:
- No test for NULL (Go nil) vs empty string: fields are string, not pointer; Go zero value "" triggers NOT NULL. Acceptable.
- No test for UPDATE path; only INSERT via RegisterUser/CreateOrder/RegisterUserPartial (no update methods in store). Engine supports only INSERT/soft-delete; no UPDATE constraint re-check. Demo never updates. Note: constraint enforcement for UPDATE is not required per README.
- No test for CHECK (status nullable?) — status is string, not pointer; empty string fails CHECK (not IN active/suspended/pending). Acceptable.
- No test for multi-column constraints; all constraints single-column. Acceptable per spec.
- No test for deferrable constraints; not implemented.
- No test for FK actions (ON DELETE CASCADE etc).
- No test for NULL in unique index (Postgres allows multiple NULLs; engine rejects empty string as duplicate — simulator limitation).
- Concurrency test uses 20 goroutines (Unsafe) and 50 (Safe) — reasonable; no variability (deterministic due to sleep). Still proves race window exists.
- No test for engine-level direct InsertUser vs store wrappers; but store is thin wrapper.
- Race detector passes -> no data races observed.

Test names: store_test.go uses plural suffixes (Constraints) while engineering/01-design.md mentions singular (TestNotNullConstraint). Minor mismatch, no functional impact.

Tests are fast (<0.1s) and race-clean (--race ok). They validate all advertised constraint behaviors and the unsafe vs safe dichotomy.

Result: tests are sufficient to prove claims; minor gaps acceptable for lab scope.