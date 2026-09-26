# Engineering Revision Plan

Target Lab: labs/17-architecture-decision-record
Previous Verdict: APPROVED_WITH_WARNINGS

## Blocking Issues
None.

## Non-Blocking Issues
1. `tests/linter_test.go` lacks test coverage for duplicate ADR ID rejection (`internal/adr/linter.go:21-23`).
2. `tests/linter_test.go` lacks concurrency stress coverage for larger record batches under race detector.

## Files To Change
- `tests/linter_test.go`

## Tests To Add/Modify
- Add test case `duplicate ADR ID` in `TestLinter_BrokenReferences`.
- Add test `TestLinter_ConcurrencyStress` to test validation on larger batches of records with mutual references.

## Validation Commands
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
