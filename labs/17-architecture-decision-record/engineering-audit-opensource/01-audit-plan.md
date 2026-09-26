# Engineering Audit Plan

Target Lab: labs/17-architecture-decision-record

Implementation Files:
- internal/adr/models.go (29 lines)
- internal/adr/parser.go (83 lines)
- internal/adr/linter.go (82 lines)
- cmd/demo/main.go (104 lines)

Tests:
- tests/parser_test.go (TestParse_Valid, TestParse_Superseded, TestParse_Supersedes, TestParse_Invalid/3 subcases)
- tests/linter_test.go (TestLinter_ValidSequence, TestLinter_BrokenReferences/4 subcases)

Executable/Demo: cmd/demo/main.go (Monolith -> Microservices, 3 ADRs)

Approved Research Inputs: SKIPPED per pipeline override (implementation + tests only)

Main Claims To Verify:
1. Parser extracts ID, Title, Status, SupersededBy, Supersedes from markdown
2. Linter accepts valid sequence (monotonic 1..n, bidirectional supersession links)
3. Linter rejects broken refs, mismatched links, non-monotonic numbering, invalid status
4. Concurrent validation is race-safe
5. Demo output is real
6. README matches code

Commands To Run:
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo

Primary Risks:
- Test suite passes but misses duplicate-ID, self-loop, empty/nil edge cases
- Parser strictness vs claimed robustness
- strings.Title deprecation (Go 1.22 still compiles, future risk)
