## Revision 1

Audit Issue: Gap 1 (Empty body text under sections passing validation)
Severity: LOW
Files Changed: 
- `internal/adr/parser.go`
- `tests/parser_test.go`
Action: Modified the parser to track active sections and enforce that at least one non-empty line of text follows each mandatory header (`## Context`, `## Decision`, `## Consequences`). Added test cases `empty context`, `empty decision`, `empty consequences`.
Verification: `go test -v ./...`
Status: RESOLVED

## Revision 2

Audit Issue: Gap 2 (Full cycle detection missing)
Severity: LOW
Files Changed: 
- `internal/adr/linter.go`
- `tests/linter_test.go`
Action: Added a cycle detection algorithm using 3-color graph traversal to detect multi-hop cycles in the `SupersededBy` chain (e.g. A -> B -> C -> A). Added test case `cyclical supersession` testing a 3-node cycle.
Verification: `go test -v ./...` and `go test -race ./...`
Status: RESOLVED
