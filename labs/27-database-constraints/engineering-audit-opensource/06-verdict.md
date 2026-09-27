# Engineering Audit Verdict

Target Lab: labs/27-database-constraints
Audit Date: 2026-09-27

## Summary

Code Files Reviewed: 6 (`internal/engine/engine.go`, `internal/store/store.go`, `internal/store/store_test.go`, `internal/model/model.go`, `internal/dberr/errors.go`, `cmd/demo/main.go`)
Tests Reviewed: 8 (`TestNotNullConstraints`, `TestCheckConstraints`, `TestUniqueConstraint`, `TestForeignKeyConstraint`, `TestPartialUniqueIndex`, `TestConcurrentRegistration_Safe_EnforcesUniqueness`, `TestConcurrentRegistration_Unsafe_SuffersRaceCondition`, `TestErrorClassification`)
Commands Executed: 5 (`go build ./...`, `go vet ./...`, `go test -v -count=1 ./...`, `go test -race -count=1 ./...`, `go run ./cmd/demo`)
Failures: 0
Warnings: 3 (2 MEDIUM + 1 LOW collection)

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS
Demo: PASS
Research Alignment: NOT_APPLICABLE (per pipeline override — implementation/tests only, research not audited)
Documentation Accuracy: WARNING

## Blocking Issues
None.

## Non-Blocking Issues
1. [MEDIUM] MISSING_TEST — no end-to-end SQLSTATE assertion in constraint tests; only `TestErrorClassification` checks codes on constructed errors. See `05-gaps.md#1`.
2. [MEDIUM] DOC_CODE_MISMATCH — `engineering/01-design.md` Architecture lists `internal/db`, `internal/errors`, `internal/domain`, `internal/service`; actual packages are `internal/engine`, `internal/dberr`, `internal/model`, `internal/store`. README itself accurate. See `04-docs-vs-code.md`.
3. [LOW] MISSING_EDGE_CASE + unsafe-test scheduler sensitivity — boundary values (age 18, total_cents ±1, soft-deleted FK target, SoftDeleteUser failure paths) untested; unsafe concurrency test relies on scheduling + 1ms sleep (passes consistently, race-clean, but theoretically flaky by design). See `03-test-audit.md` and `05-gaps.md`.

## Required Revisions
None required for approval. Recommended (non-blocking):
1. Add `dberr.IsConstraintViolation(err, expectedCode)` assertions to each constraint test to prove claimed SQLSTATE end-to-end.
2. Fix `engineering/01-design.md` package names to match actual layout (`internal/engine`, `internal/dberr`, `internal/model`, `internal/store`).
3. Optionally add boundary tests (age 18/17, total_cents 1/0/-1, FK to soft-deleted user, SoftDeleteUser nil/non-existent).

## Final Status

APPROVED_WITH_WARNINGS
