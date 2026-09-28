# Gaps

## GAP-1
Type: MISSING_EDGE_CASE
Location: `internal/deadline/deadline.go:18-20`
Severity: LOW
Description: Worker goroutine leaks if `fn` ignores `childCtx.Done()`. `select` returns on timeout but nothing stops `fn`; `done` buffer size 1 avoids sender block only if `fn` eventually returns.
Evidence: Code review; tests/demo all honor ctx so no leak observed; `go test -race` clean.
Recommendation: Document requirement that `fn` must honor ctx; no code change required for audit.

## GAP-2
Type: DOC_CODE_MISMATCH
Location: `engineering/01-design.md` Success Criteria 1 vs `internal/deadline/deadline.go`
Severity: LOW
Description: Design states "context cancellation terminates long-running downstream work without leakage" — implementation only guarantees this when downstream honors ctx.
Evidence: Finding 1 code audit; demo/test fns honor ctx.
Recommendation: Soften wording to "provided downstream honors context"; non-blocking.

No HIGH or CRITICAL gaps. No MISSING_TEST, BROKEN_IMPLEMENTATION, RACE_CONDITION, UNHANDLED_ERROR, FAKE_DEMO, FAKE_BENCHMARK, or UNVERIFIED_RESULT found.
