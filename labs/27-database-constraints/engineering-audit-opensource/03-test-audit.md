# Test Audit — labs/27-database-constraints

Commands executed (fresh, `-count=1`):
- `go test ./...` → PASS (store ok, others no test files)
- `go test -race ./...` → PASS, no data race
- `go vet ./...` → PASS
- `go test -v -count=1 ./...` → 8/8 PASS:
  TestNotNullConstraints, TestCheckConstraints, TestUniqueConstraint,
  TestForeignKeyConstraint, TestPartialUniqueIndex,
  TestConcurrentRegistration_Safe_EnforcesUniqueness,
  TestConcurrentRegistration_Unsafe_SuffersRaceCondition,
  TestErrorClassification
- `go run ./cmd/demo` → PASS, output matches engineering/03 (plus taxonomy line)

Coverage:
- Happy path: PASS — first unique insert, valid FK order, post-soft-delete reuse verified.
- Failure path: PASS — NOT NULL (email, username, user_id), CHECK (age, status, total_cents), UNIQUE, FK orphan, partial duplicate all assert error.
- Edge cases: WARNING — no boundary happy-path asserts (age==18, total_cents==1, status suspended/pending). Failure boundaries tested, success boundaries not.
- Transitions: PASS — partial lifecycle 6-step (insert → reject → soft-delete → reuse → reject → deleted-row bypass) is thorough.
- Recovery/rollback: WARNING — failed inserts never assert state unchanged (count/index intact). Concurrency test asserts final count==1, single-path tests do not.
- Concurrency: PASS — safe (exactly 1/20 succeeds) and unsafe (duplicates >1) both proven; `-race` clean.
- Negative cases: PASS — duplicate, orphan, invalid enum all covered.
- Error mapping: WARNING — `TestErrorClassification` covers constructors + `IsConstraintViolation` only; `MapToDomainError` output never tested and drops SQLSTATE (see code-audit Finding 5).
- Untested units: `SoftDeleteUser` error paths (missing ID, nil DeletedAt), `InsertUserUnsafe` direct, `GetUser` — LOW.

Assessment: suite proves core claims. Passing verdict is earned, not vacuous. Gaps are boundary/mapping coverage, not core behavior.
