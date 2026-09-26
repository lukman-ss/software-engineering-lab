## DOC_CODE_MISMATCH 1

Location: README.md lines 7-11 (Implemented Constraints list) vs code behavior
Claim in README: 
- NOT NULL (23502): ensures required columns cannot accept null/empty values.
- CHECK (23514): evaluates boolean predicate logic (e.g. age >= 18, status IN (...), total_cents > 0).
- UNIQUE (23505): enforces single occurrence across rows.
- FOREIGN KEY (23503): enforces referential integrity preventing orphan rows.
- PARTIAL UNIQUE INDEX: conditional uniqueness (WHERE deleted_at IS NULL) enabling soft delete re-registration.
Code Observation: 
- NOT NULL: engine treats empty string as violation (email=="", username=="") -> matches.
- CHECK: age<18, status not in enum, total_cents<=0 -> matches.
- UNIQUE: emailIndex and activeEmails maps enforce uniqueness -> matches.
- FOREIGN KEY: user_id lookup in e.users map -> matches.
- PARTIAL UNIQUE: activeEmails only updated when DeletedAt==NULL -> matches.
Assessment: PASS (no mismatch)

## DOC_CODE_MISMATCH 2

Location: README.md lines 13-25 (Running Tests/Demo) vs actual execution
Claim: go test ./... works; go test -race ./... works; go run ./cmd/demo works.
Observed: All three commands succeed (see engineering/audit-opensource/01-audit-plan.md for evidence).
Assessment: PASS (no mismatch)

## DOC_CODE_MISMATCH 3

Location: README.md line 9: "CHECK Constraints (23514): Evaluates boolean predicate logic on row data (e.g. age >= 18, status IN (...), total_cents > 0)."
Code: engine.go check for status uses switch case "active", "suspended", "pending". Matches claim.
Assessment: PASS

## DOC_CODE_MISMATCH 4

Location: README.md line 10: "UNIQUE Constraints (23505): Enforces single occurrence across rows and prevents concurrent read-then-write race conditions."
Code: InsertUser locks mutex, checks index, only one writer can win; concurrent test shows exactly 1 success.
Assessment: PASS

## DOC_CODE_MISMATCH 5

Location: README.md line 11: "PARTIAL UNIQUE INDEX: Demonstrates conditional uniqueness (WHERE deleted_at IS NULL) enabling soft delete re-registration while maintaining active uniqueness."
Code: SoftDeleteUser removes from activeEmails; RegisterUserPartial allows re-insert if deleted_at!=NULL.
Assessment: PASS

## TEST_CLAIM_MISMATCH

Location: store_test.go:210-239 (TestConcurrentRegistration_Safe_EnforcesUniqueness) vs claim in engineering/02-implementation-notes.md line 26: "Concurrency test proving that unsafe store suffers race condition duplicates while safe store (with DB constraints) guarantees exactly 1 record created and all concurrent duplicates return 23505."
Test: safe store asserts successCount==1 and errorCount==goroutines-1; unsafe store asserts final count>1.
Code: matches claim exactly.
Assessment: PASS

## RESEARCH_IMPLEMENTATION_MISMATCH

Location: out of scope per pipeline override (audit implementation/tests only)
Assessment: NOT_APPLICABLE