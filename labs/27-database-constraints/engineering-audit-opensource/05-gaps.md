# Gap Analysis — labs/27-database-constraints

## Gap 1
Type: UNVERIFIED_RESULT
Location: `internal/store/store_test.go` — missing boundary happy-path asserts
Severity: LOW
Description: Failure boundaries tested (age 16, total_cents 0) but no success asserts for age==18, total_cents==1, statuses suspended/pending.
Recommendation: Add one happy-path boundary test.

## Gap 2
Type: MISSING_TEST
Location: `internal/engine/engine.go:102-120` (`SoftDeleteUser` error paths)
Severity: LOW
Description: Missing-ID and nil-DeletedAt errors untested; cross-path case (full-unique user soft-deleted, then reuse blocked by `emailIndex`) untested.
Recommendation: Add error-path tests; document that soft-delete reuse applies to the partial-index path only.

## Gap 3
Type: UNHANDLED_ERROR (ergonomic, not silent)
Location: `internal/dberr/errors.go:83-99` (`MapToDomainError`)
Severity: MEDIUM
Description: Mapped error drops SQLSTATE; `errors.As` classification impossible after mapping. Tests never cover mapped output.
Recommendation: Wrap with `%w` or document display-only contract; add mapping test.

## Non-gaps (explicitly checked, none found)
- RACE_CONDITION: `go test -race` clean; engine mutex covers all maps.
- FAKE_DEMO / FAKE_BENCHMARK: demo output reproduced live; no benchmarks claimed.
- BROKEN_IMPLEMENTATION / MISSING_TEST (core): all five constraint types plus concurrency pair proven.
- DOC_CODE_MISMATCH: none — README/engineering notes match code and live output.
