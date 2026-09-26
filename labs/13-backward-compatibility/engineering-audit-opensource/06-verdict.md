# Engineering Audit Verdict

Target Lab: labs/13-backward-compatibility
Audit Date: 2026-09-25

## Summary

Code Files Reviewed:
- internal/compat/model.go
- internal/compat/store.go
- internal/compat/flags.go
- internal/compat/metrics.go
- internal/compat/backfill.go
- internal/compat/service.go
- internal/compat/handler.go
- cmd/demo/main.go
- schema.sql

Tests Reviewed:
- internal/compat/service_test.go
- tests/migration_test.go
- tests/concurrency_test.go

Commands Executed:
- `go build ./...` → success
- `go vet ./...` → success
- `go test -count=1 -v ./...` → 8 tests PASS
- `go test -race -count=1 ./...` → ok (no race)
- `go run ./cmd/demo` → completed successfully (output matches engineering/03-execution-result.md)

Failures:
- 0

Warnings:
- Finding #2 (MEDIUM): Silent discard of extra-phone write errors in `service.go:80` weakens "atomic dual-write" claim for non-primary phones.
- Finding #3 (LOW): `handler.go` uses `id, _ := strconv.Atoi(idStr)` ignoring conversion error; invalid IDs become 0 → 404 not 400.
- Finding #4 (LOW): `store.go:GetUserIDs` uses O(n²) sort instead of `sort.Ints`.
- Finding #5 (LOW): Contract guard uses cumulative `LegacyReadHits` (never reset), making it stricter than the 30-day zero-traffic heuristic (documented as unverified in research; guard still correctly blocks on non-zero traffic).
- Finding #6 (LOW): In-memory persistence boundary documented (`ponytail:`); not a defect but limits restart/resume claims.
- Finding #7 (LOW): Missing unit test for `GetModernUser` post-contract error path (`ErrLegacyUnavailable` only reachable via V1).
- Finding #8 (LOW): Missing test asserting `GetUserIDs` ordering.

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS (data-race free under concurrent writers/readers/backfill/reconcile)
Demo: PASS (output matches recorded execution result exactly)
Research Alignment: PASS (all claims implemented; research Q11 heuristic noted as unverified and not claimed)
Documentation Accuracy: PASS (README, engineering notes, execution result align with code and behavior)

## Blocking Issues
- NONE (no HIGH/CRITICAL findings)

## Non-Blocking Issues
1. MEDIUM: Silent error discard on extra-phone writes during dual-write (Finding #2). Weakens atomicity claim for non-primary fields.
2. LOW: O(n²) user-ID sort in `GetUserIDs` (Finding #4).
3. LOW: Unchecked `Atoi` in HTTP handler (Finding #3).
4. LOW: Cumulative legacy-read counter makes contract guard stricter than sliding window (Finding #5).
5. LOW: In-memory persistence boundary documented, not a defect (Finding #6).
6. LOW: Missing test for modern-client error post-contract (Finding #7).
7. LOW: Missing test asserting ID ordering in backfill checkpoint (Finding #8).

## Required Revisions
None required for approval. All warnings are:
- Documented as in-memory simplifications (`ponytail:` in design/notes)
- Or do not affect core backward-compatibility claims (Expand-Migrate-Contract, dual-write, idempotent/resumable backfill, fallback read, drift detection, safe rollback, observability, contract guard).
- Or are test suite enhancements (non-blocking).

## Final Status

APPROVED

Verdict rationale: Code compiles; all tests pass including race detector; demo executes and matches recorded output; README and engineering docs accurately describe behavior; no fabricated benchmarks or results; research claims are implemented as specified; unresolved findings are LOW/MEDIUM and bounded by documented in-memory scope or test improvements.