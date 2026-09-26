# Engineering Changes Made

## Revision 1

Audit Issue: Lack of test exercising duplicate ADR ID validation branch in `internal/adr/linter.go:21-23`
Severity: MEDIUM
Files Changed: `tests/linter_test.go`
Action: Added `duplicate ADR ID` test case to `TestLinter_BrokenReferences`.
Verification: `go test -v -run TestLinter_BrokenReferences/duplicate_ADR_ID ./tests`
Status: RESOLVED

## Revision 2

Audit Issue: Lack of concurrency stress coverage for larger record batches under race detector
Severity: LOW
Files Changed: `tests/linter_test.go`
Action: Added `TestLinter_ConcurrencyStress` testing 100 pairwise superseded/superseding ADRs in parallel.
Verification: `go test -race -run TestLinter_ConcurrencyStress ./tests`
Status: RESOLVED
