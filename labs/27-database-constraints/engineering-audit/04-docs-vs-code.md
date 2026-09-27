# Docs vs Code Analysis

## Comparisons

1. README.md vs Code
   - README specifies 5 implemented constraints: NOT NULL (`23502`), CHECK (`23514`), UNIQUE (`23505`), FOREIGN KEY (`23503`), PARTIAL UNIQUE INDEX (`WHERE deleted_at IS NULL`).
   - Implementation in `internal/engine/engine.go` implements all 5 constraints with matching SQLSTATE codes.
   - Result: MATCH.

2. Engineering Design vs Implementation
   - Design specified safe vs unsafe stores, in-memory ACID simulation, domain error mapping, and concurrency tests.
   - All modules (`internal/model`, `internal/dberr`, `internal/engine`, `internal/store`, `cmd/demo`) are present and implemented as designed.
   - Result: MATCH.

3. Research Claims vs Code Demonstration
   - Research Finding 1 & 2 (Application check-then-act fails; UNIQUE constraint succeeds atomically) -> verified by `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` and `TestConcurrentRegistration_Safe_EnforcesUniqueness`.
   - Research Finding 3 (NOT NULL) -> verified by `TestNotNullConstraints`.
   - Research Finding 4 (FOREIGN KEY) -> verified by `TestForeignKeyConstraint`.
   - Research Finding 5 (CHECK constraints row-scoped) -> verified by `TestCheckConstraints`.
   - Research Finding 6 (Partial Unique Index for soft-delete) -> verified by `TestPartialUniqueIndex`.
   - Research Finding 7 (SQLSTATE Class 23 error taxonomy) -> verified by `TestErrorClassification`.
   - Result: MATCH.

## Findings Summary
- DOC_CODE_MISMATCH: None.
- TEST_CLAIM_MISMATCH: None.
- RESEARCH_IMPLEMENTATION_MISMATCH: None.
