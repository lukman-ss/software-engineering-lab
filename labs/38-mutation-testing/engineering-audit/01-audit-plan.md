# Engineering Audit Plan

Target Lab: labs/38-mutation-testing
Implementation Files:
- internal/engine/types.go
- internal/engine/mutator.go
- internal/engine/runner.go
- internal/service/discount.go
- cmd/demo/main.go

Tests:
- internal/service/discount_weak_test.go
- internal/service/discount_strong_test.go
- tests/engine_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/01-plan.md
- research/05-report.md

Main Claims To Verify:
1. AST mutation engine correctly parses and mutates Go code across 4+ operator types (Relational, Boolean, Arithmetic, Boundary/Value).
2. Weak test suite achieves 100% line coverage on `discount.go` while missing injected mutations (yielding 0.00% mutation score).
3. Strong test suite detects all mutants (yielding 100.00% mutation score).
4. Go test execution (`go test ./...`) and race detector (`go test -race ./...`) pass cleanly.
5. Demo runner (`go run ./cmd/demo`) runs without error and outputs accurate score comparison.
6. Documentation (README.md, engineering notes) accurately describes implemented architecture and behavior without overclaiming.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
- `go test -coverprofile=coverage.out ./internal/service && go tool cover -func=coverage.out`

Primary Risks:
- Simulated test evaluation functions (`TestFunc`) in demo/engine test using AST string diff instead of dynamic runtime execution/assertion testing.
- Concurrency race conditions in runner AST modifications if FileSet or AST node modification shared across goroutines.
- Unhandled AST cases causing panic during AST inspection or formatting.
