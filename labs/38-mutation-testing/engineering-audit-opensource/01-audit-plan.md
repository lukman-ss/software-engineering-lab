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
- labs/38-mutation-testing/research/05-report.md
- labs/38-mutation-testing/research-audit/07-verdict.md
- labs/38-mutation-testing/research-revision/03-revision-result.md

Main Claims To Verify:
1. AST-based mutation testing framework generates mutants across 4 active mutation operator types (Relational, Boolean, Arithmetic, Boundary Shift).
2. Weak test suite achieves 100% statement/line coverage on `internal/service/discount.go` while failing to kill mutants (low/0% mutation score).
3. Strong test suite tests all edge cases and boundary behaviors, achieving 100% mutation score.
4. Concurrent runner executes mutant evaluations safely without data races.
5. Demo output matches claimed behavior and runtime execution.
6. Documentation in README.md accurately reflects implementation and commands.

Commands To Run:
- `go test -v ./...`
- `go test -count=1 -race ./...`
- `go test -run=TestCalculateDiscount_Weak -coverprofile=weak_cov.out ./internal/service && go tool cover -func=weak_cov.out`
- `go run ./cmd/demo`

Primary Risks:
1. Disconnect between textual test simulation in engine runner demo and actual compiled test execution.
2. Incomplete operator implementations (e.g., StatementDelete declared in types but not emitted).
3. Data races in concurrent runner goroutines or mutant slice allocations.
4. Accuracy of 100% statement coverage claim on the weak test suite.
