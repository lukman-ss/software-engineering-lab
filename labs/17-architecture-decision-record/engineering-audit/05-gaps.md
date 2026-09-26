# Gap Analysis

## MISSING_TEST
**Severity:** MEDIUM
**Component:** `tests/linter_test.go`
**Description:** `internal/adr/linter.go` contains error-handling logic for duplicate ADR IDs (lines 21-23). However, there is no corresponding test in `linter_test.go` to prove that the linter correctly rejects collections with duplicate IDs.
