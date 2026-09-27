# Docs vs Code

## README.md vs Implementation
- Claims 5 constraint types: NOT NULL (23502), CHECK (23514), UNIQUE (23505), FK (23503), PARTIAL UNIQUE INDEX. Code implements all.
- Running tests: says go test -v ./... and go test -race ./...; matches execution.
- Demo: says go run ./cmd/demo; matches.
- Implemented Constraints section lines 7-11 map 1:1 to test functions:
  1. NOT NULL → TestNotNullConstraints
  2. CHECK → TestCheckConstraints
  3. UNIQUE → TestUniqueConstraint
  4. FOREIGN KEY → TestForeignKeyConstraint
  5. PARTIAL UNIQUE INDEX → TestPartialUniqueIndex
- No overclaim: README does not claim UPDATE, cross-row, or serializable.
- Assertion: "NOT NULL ... cannot accept null/empty values." Code treats empty string "" as violation (no nil pointer). Acceptable in-memory proxy.
- Assertion: "CHECK Evaluates boolean predicate logic on row data (e.g. age >= 18, status IN (...), total_cents > 0)." Code matches.
- Assertion: "UNIQUE Enforces single occurrence across rows and prevents concurrent read-then-write race conditions." Code: engine.unique via emailIndex/activeEmails map; TestUniqueConstraint + concurrency test prove.
- Assertion: "FOREIGN KEY Enforces referential integrity preventing orphan rows." Code: InsertOrder checks e.users[o.UserID] exists; TestForeignKeyConstraint proves.
- Assertion: "PARTIAL UNIQUE INDEX: Demonstrates conditional uniqueness (WHERE deleted_at IS NULL) enabling soft delete re-registration while maintaining active uniqueness." Code: TestPartialUniqueIndex + demo steps 4-6 prove.
- No mention of SQLSTATE mapping; but demo prints domain errors; error classification in store_test.go covers mapping.

## Engineering notes (01-design.md, 02-implementation-notes.md, 03-execution-result.md) vs code/tests/demo
- Design doc Success Criteria #2: "Concurrency test proving that unsafe store suffers race condition duplicates while safe store (with DB constraints) guarantees exactly 1 record created and all concurrent duplicates return 23505." Code: TestConcurrentRegistration_Unsafe_SuffersRaceCondition + TestConcurrentRegistration_Safe_EnforcesUniqueness prove; demo [5] mirrors.
- Design doc #3: "Error classification accurately categorizing 23502, 23503, 23505, and 23514." Code: TestErrorClassification proves.
- Design doc #4: "Tests pass with go test ./... and go test -race ./...". Verified.
- Design doc #5: "Standalone demo in cmd/demo/main.go runs with clear output demonstrating safe vs unsafe concurrency and constraint violations." Verified.
- Implementation Notes: Files added list matches actual files created.
- Core Design Decisions: storage engine integrity via mutex/index layer; SQLSTATE standardization; partial unique index semantics — all present.
- Implementation-Specific Choices: in-memory simulation using stdlib sync.RWMutex, atomic.Int64, map — accurate.
- Known Limitations: does not connect to live PG; EXCLUDE omitted — accurate.
- Trade-offs: coarse-grained table locks — accurate.
- What Is Demonstrated: enumerates 7 items; all present in code/tests/demo.
- What Is Not Demonstrated: distributed constraints, full SQL parser — accurate.

## Execution-Result.md vs actual runs
- Build: go build ./... => PASS (exit code 0) matches.
- Tests: Lists 7 tests (omits TestConcurrentRegistration_Unsafe_SuffersRaceCondition). Actual go test -v ./... runs and PASSes all 8 tests. Omission in document is a DOC_CODE_MISMATCH (missing test line).
- Race Detector: lists ok internal/store 1.330s; matches actual run (1.101s) — PASS.
- Demo: output block exactly matches actual go run ./cmd/demo output (verified line-for-line). No deviation.
- Final Engineering Status: READY_FOR_ENGINEERING_AUDIT — accurate.

## Mismatches Found
1. DOC_CODE_MISMATCH: execution-result.md "Tests" section printed test list omits TestConcurrentRegistration_Unsafe_SuffersRaceCondition (though test actually ran and passed). This is a documentation gap only; test exists and validates claim.
2. Minor: design doc lists singular test names (TestNotNullConstraint) while code uses plural (TestNotNullConstants). No functional impact.
3. No evidence of fake benchmark/result: demo output deterministic, races present, tests pass under -race.
4. Implementation-specific decisions (pure Go in-memory) scoped correctly; does not overclaim to be real Postgres.

No HIGH/CRITICAL mismatches; all core behavior proven; docs align except for missing test line in execution-result.md and test naming inconsistency.