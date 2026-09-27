# Gap Analysis

| ID | Type | Location | Severity | Summary |
|----|------|----------|----------|---------|
| G1 | DOC_CODE_MISMATCH | engineering/03-execution-result.md "Tests" section | LOW | Printed `go test -v` list omits TestConcurrentRegistration_Unsafe_SuffersRaceCondition. Test exists and passes; only the log omits it. |
| G2 | DOC_CODE_MISMATCH | engineering/01-design.md Test Strategy | LOW | Design doc uses singular names (TestNotNullConstraint) while code/tests use plural (TestNotNullConstants). |
| G3 | MISSING_EDGE_CASE | engine.go | LOW | No UPDATE path or constraint re-validation on update. README scope is INSERT/soft-delete only; acceptable for lab. |
| G4 | MISSING_EDGE_CASE | engine.go | LOW | No UNIQUE NULL-distinct handling (Postgres allows multiple NULLs in UNIQUE); in-memory uses empty-string. Limitation, not defect. |
| G5 | MISSING_EDGE_CASE | store.go | LOW | FK insert order has no ON DELETE action. Out of scope per README. |
| G6 | UNHANDLED_ERROR | store.go | MEDIUM | ctx.Context accepted by all SafeStore methods but never respected (no ctx.Done check). Demo cancels none, but real callers could leak. |
| G7 | MISSING_TEST | store.go / engine.go | MEDIUM | No test asserting FK with valid user creates order with non-zero ID and persisted — actually TestForeignKeyConstraint checks o.ID==0; covered. No gap. (refuted) |
| G8 | MISSING_TEST | engine.go | LOW | No test for partial-index insert of soft-deleted email reusing email while an active one already exists (engine rejects if any active exists). TestPartialUniqueIndex step 6 inserts soft-deleted only. Could add negative case. |
| G9 | MISSING_TEST | dberr | LOW | No test asserting MapToDomainError for each SQLSTATE returns distinct message (only that error != nil at store layer). Mapping correctness asserted only via TestErrorClassification (IsConstraintViolation), not message content. |

No BROKEN_IMPLEMENTATION, RACE_CONDITION, UNHANDLED_ERROR (at engine), MISSING_EDGE_CASE blocking, IMPLEMENTATION_OVERCLAIM, RESEARCH_MISMATCH, FAKE_DEMO, FAKE_BENCHMARK, or UNVERIFIED_RESULT.

Demo output is deterministic and reproducible (all operations synchronous except the controlled concurrency section). Race detector clean.
