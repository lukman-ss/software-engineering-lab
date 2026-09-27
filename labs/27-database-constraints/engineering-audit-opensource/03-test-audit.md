# Test Audit

Test file: internal/store/store_test.go (9 tests)

Coverage Matrix:
- TestNotNullConstraints: tests NULL email, NULL username, NULL user_id in Order → asserts err != nil, no specific SQLSTATE asserted on store path.
- TestCheckConstraints: tests age<18, invalid status, total_cents<=0 → same.
- TestUniqueConstraint: single-dup email → asserts err != nil.
- TestForeignKeyConstraint: order with user_id=9999 (nonexistent) + valid case → asserts err != nil on first, asserts ID nonzero on second.
- TestPartialUniqueIndex: multi-step: active insert → dup active fails → soft-delete → insert succeeds → second active fails → soft-deleted insert (with DeletedAt set) succeeds.
- TestConcurrentRegistration_Safe_EnforcesUniqueness: 20 goroutines, same email, SafeStore → asserts successCount == 1 && errorCount == 19 && total users == 1.
- TestConcurrentRegistration_Unsafe_SuffersRaceCondition: 20 goroutines, same email, UnsafeStore → asserts users.Count > 1 (no upper bound).
- TestErrorClassification: builds four raw *ConstraintError and asserts dberr.IsConstraintViolation matches.

No tests assert positive success path for happy-path insert (ID > 0, no error). All "happy path" tests are indirect (TestForeignKeyConstraint success case, TestPartialUniqueIndex success step, concurrency successCount==1).

Test Quality:
- No flaky timing except UnsafeRace (sleep 1ms) → still passed consistently.
- Tests are deterministic w.r.t engine mutex.
- No table-driven subtests; duplication acceptable given small count.
- No explicit cleanup needed (each test constructs fresh Engine).

Concurrency safety:
- SafeStore test: 20 goroutines, exactly one winner, validated.
- UnsafeStore test: expects >1 duplicates, validated.
- Race detector on test suite: clean (0 races).

Edge Cases missing:
- Zero-value User.ID / Order.ID (engine assigns on zero) not explicitly tested.
- Negative age (int, not unsigned) already caught by age<18.
- Empty-string Username caught by NOT NULL.
- Very long strings (no length limit in model).
- DeletedAt pointer nil vs non-nil in partial index (covered).
- Soft-delete on already-soft-deleted user (no test; engine currently allows overwrite; not claimed).

Overclaim vs test:
- None: each claimed behavior has at least one test asserting the rejecting or accepting case.

False negative risk:
- LOW — test suite enforces each rejection branch.

Verdict:
- Core behavior (constraints + concurrency) proven; happy-path implicit.
- No test leaks state (t.Cleanup not needed).
- No negative test for context cancellation.