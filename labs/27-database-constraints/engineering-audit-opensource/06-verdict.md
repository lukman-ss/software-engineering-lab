# Engineering Audit Verdict

Target Lab: labs/27-database-constraints
Audit Date: 2026-09-26
Scope: Implementation + tests only (research/content not audited per pipeline override)

## Summary

Code Files Reviewed:
- internal/engine/engine.go (161 lines)
- internal/store/store.go (84 lines)
- internal/model/model.go (20 lines)
- internal/dberr/errors.go (100 lines)
- cmd/demo/main.go (99 lines)

Tests Reviewed:
- internal/store/store_test.go (8 tests)

Commands Executed:
- go build ./... -> PASS (exit 0)
- go test ./... -> PASS (cached)
- go test -v ./... -> PASS (8/8 tests, 0.312s)
- go test -race -count=1 ./... -> PASS (1.290s, no data race)
- go run ./cmd/demo -> ran; output reproduced (see below)
- go test -run TestConcurrentRegistration_Unsafe -count=10 ./internal/store/ -> PASS (0.148s)

Demo Output (reproduced verbatim):
```
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

Failures: None
Warnings: 2 LOW + 1 MEDIUM (see Non-Blocking Issues)

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS
Demo: PASS
Research Alignment: NOT_APPLICABLE (out of scope per pipeline override)
Documentation Accuracy: PASS (README matches code exactly; minor staleness in engineering/03-execution-result.md, see Non-Blocking)

## Blocking Issues
None (no HIGH/CRITICAL findings).

## Non-Blocking Issues
1. [LOW] engine.go:77-90 — InsertUser allows caller-supplied ID; no uniqueness guard on explicit non-zero IDs. A caller passing an existing non-zero ID will silently overwrite (map assignment). Not exercised by tests or README.
2. [MEDIUM] store.go:28-35 — UnsafeStore.RegisterUser iterates integer Primary Key range [1..count] assuming dense packing; holes from failed/tx-style inserts would cause false-negative duplicate checks. Acceptable as a conceptual demo, but the "vulnerable pattern" is not a true memory data race under the engine mutex.
3. [LOW] engineering/03-execution-result.md omits TestConcurrentRegistration_Unsafe_SuffersRaceCondition from its PASS list though the suite has 8 tests; stale execution log (actual output verified by this audit).

## Required Revisions
None for approval (out of scope to modify per pipeline override).

## Final Status

APPROVED_WITH_WARNINGS

The implementation compiles, all 8 tests pass including the -race detector, the demo output is reproduced verbatim, and the README matches the code. Two LOW and one MEDIUM informational findings exist around edge-case handling and a stale engineering log, but no core behavior is unproven or incorrect. Approved for Technical Writer handoff.
