# Gap Analysis

Allowed gap types per instructions: MISSING_TEST, BROKEN_IMPLEMENTATION, DOC_CODE_MISMATCH, RACE_CONDITION, UNHANDLED_ERROR, MISSING_EDGE_CASE, IMPLEMENTATION_OVERCLAIM, RESEARCH_MISMATCH, FAKE_DEMO, FAKE_BENCHMARK, UNVERIFIED_RESULT.

Identified gaps:

1. Gap: MISSING_TEST
   Location: `internal/store/store_test.go`
   Description: No test asserts that constraint violations from `SafeStore.RegisterUser`/`SafeStore.CreateOrder` carry the expected SQLSTATE (23502, 23503, 23505, 23514) end-to-end. Functional tests only check `err != nil`; only `TestErrorClassification` uses constructed errors to test `IsConstraintViolation` and `MapToDomainError`. Thus proof that the engine returns the claimed SQLSTATE in practice is indirect (via demo error messages, which are after `MapToDomainError`). While demo shows user messages aligning with expectation, a test that verifies raw error codes strengthens the guarantee. Severity: MEDIUM. Note: `TestPartialUniqueIndex` checks error messages but not raw SQLSTATE. To close: add an assertion in each constraint test that `dberr.IsConstraintViolation(err, expectedSQLState)` is true, or that error unwraps to `*dberr.ConstraintError` with matching `.Code`.

2. Gap: MISSING_EDGE_CASE
   Location: Various constraint tests.
   Description: Boundary values not explicitly tested:
     - NOT NULL: empty string already tested (failure); whitespace-only string (e.g. `" "`) currently passes as non-empty because check is `== ""`; is that correct? Not asserted; but column is TEXT; whitespace is not NULL so constraint should pass — not a bug but worth a note.
     - CHECK age >= 18: age 17 (fail), 18 (pass) untested; only age 16 tested.
     - CHECK status enum: empty string `""` tested implicitly via default branch (fail); single character invalid untested.
     - CHECK total_cents > 0: total_cents = 1 (pass), 0 (fail), -1 (fail) — only zero tested.
     - FOREIGN KEY: referencing a soft-deleted user (DeletedAt set) should succeed because row still exists; untested.
   Severity: LOW (each individual missing edge low; collectively MEDIUM if many). No incorrect behavior asserted; only missing explicit verification.

3. Gap: DOC_CODE_MISMATCH
   Location: `engineering/01-design.md`
   Description: Documented package structure differs from actual: doc mentions `internal/db`, `internal/errors`, `internal/domain`, `internal/service`; actual are `internal/engine`, `internal/dberr`, `internal/model`, and logic in `internal/store`. This is documentation drift. Severity: MEDIUM (per instructions: DOC_CODE_MISMATCH allowed). Note: Conceptual architecture aligns; no code defect.

4. Gap: UNHANDLED_ERROR
   Location: None found. All error paths checked; engine returns errors, store maps or propagates; demo prints errors; tests expect errors. No swallowed errors.

5. Gap: RACE_CONDITION
   Location: The unsafe concurrency test intentionally exposes a race condition to demonstrate the point. This is not a bug; it is a deliberate design of the test to show vulnerability. However, the test itself has a flakiness risk: if the scheduler happens to serialize goroutines such that each sees the prior insert, count would be 1 and test would fail spuriously. The 1ms sleep mitigates but does not eliminate. In practice race-detector clean and test passes consistently; still noted. Severity: LOW (test flakiness, not production race). Not a BROKEN_IMPLEMENTATION.

6. Gap: IMPLEMENTATION_OVERCLAIM
   Location: None. Claims match implementation (engine enforces constraints atomically; unsafe store demonstrates race; safe store prevents it). No overstatement beyond what code shows.

7. Gap: RESEARCH_MISMATCH
   Location: N/A (per PIPELINE OVERRIDE: audit implementation and tests only; do not audit research/content). Skipped.

8. Gap: FAKE_DEMO / FAKE_BENCHMARK / UNVERIFIED_RESULT
   Location: None. Demo output recorded live; commands executed and output captured; no fabrication.

Summary of gap types with counts:
- MISSING_TEST: 1
- MISSING_EDGE_CASE: 1
- DOC_CODE_MISMATCH: 1