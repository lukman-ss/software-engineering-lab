# Engineering Revision Plan

Target Lab: labs/17-architecture-decision-record
Previous Verdict: APPROVED (with non-blocking gaps)

## Blocking Issues
None.

## Non-Blocking Issues
1. Section content completeness check omitted in parser (headers validated, but empty section text accepted).
2. Deep cycle detection in supersession lineage graph omitted (only direct 1:1 bidirectional links validated).

## Files To Change
- `internal/adr/parser.go`
- `internal/adr/linter.go`

## Tests To Add/Modify
- `tests/parser_test.go` (Add test cases for empty Context, Decision, Consequences sections)
- `tests/linter_test.go` (Add test case for multi-hop cyclical supersession)

## Validation Commands
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
