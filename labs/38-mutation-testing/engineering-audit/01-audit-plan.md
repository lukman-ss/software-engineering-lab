# Engineering Audit Plan

Target Lab: labs/38-mutation-testing
Implementation Files:
- internal/engine/mutator.go
- internal/engine/runner.go
- internal/engine/types.go
- internal/service/discount.go
- cmd/demo/main.go

Tests:
- internal/service/discount_weak_test.go
- internal/service/discount_strong_test.go
- tests/engine_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/05-report.md
- engineering/01-design.md
- engineering/02-implementation-notes.md

Main Claims To Verify:
1. AST mutator correctly parses and mutates Go code across relational, boolean, arithmetic, and boundary categories.
2. Weak unit test suite achieves 100% statement coverage on discount.go but fails to kill AST mutations.
3. Strong unit test suite catches boundary conditions and kills all generated non-equivalent mutants (100% score).
4. Parallel mutant test execution in runner is concurrency-safe and passes `go test -race ./...`.
5. Demo CLI executes and displays real mutation contrast results without mock data.
6. Documentation in README.md matches code structure and execution commands.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
- `go test -coverprofile=coverage.out ./internal/service && go tool cover -func=coverage.out`

Primary Risks:
- Thread-safety issues or state corruption during parallel AST parsing and mutation execution.
- Statement deletion operator claimed in design/types but missing implementation in engine walker.
- Test function approximations in demo/engine tests vs actual unit test assertions.
