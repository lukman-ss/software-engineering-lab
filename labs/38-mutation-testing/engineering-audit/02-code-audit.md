# Code Audit

Target Lab: labs/38-mutation-testing

## Finding 1: AST Mutator Implementation & Concurrency Safety

Location: `internal/engine/mutator.go`, `internal/engine/runner.go`
Claimed Behavior: Safe, AST-based mutation generation and concurrent mutant evaluation.
Observed Implementation:
- `mutator.go` uses Go stdlib (`go/ast`, `go/parser`, `go/token`, `go/format`) to construct mutation plans with `Apply` closures.
- `runner.go` spawns goroutines per mutant. Crucially, inside each goroutine, `parser.ParseFile(fset, "source.go", sourceCode, 0)` is called with a local `token.FileSet`, ensuring AST nodes and file positions are distinct per goroutine.
- Results array is pre-allocated by size (`results := make([]MutantResult, len(plans))`) and indexed by worker thread, preventing race conditions on slice mutation.
Assessment: PASS
Severity: LOW
Notes: Clean concurrent implementation verified by `go test -race ./...`.

## Finding 2: Standard Statement Deletion Mutation Type Unimplemented

Location: `internal/engine/types.go`, `internal/engine/mutator.go`
Claimed Behavior: Documented support for 5 mutation categories including `StatementDelete`.
Observed Implementation:
- `StatementDelete` is defined as a `MutationType` constant in `types.go:12`.
- `mutator.go` does not implement an `ast.Inspect` branch for statement deletion node types (`*ast.ExprStmt`, `*ast.AssignStmt`, etc.).
- Engineering documentation (`02-implementation-notes.md:44`) explicitly disclaims statement deletion as registered but not exercised due to `discount.go` structure. README (`README.md:48`) lists only 4 implemented mutation operators.
Assessment: WARNING
Severity: LOW
Notes: Minor mismatch between constant enum and mutator inspection logic, though correctly noted in engineering limitations.

## Finding 3: Simplified Test Oracle (AST String Diff vs Runtime Execution)

Location: `cmd/demo/main.go:41`, `tests/engine_test.go:92`
Claimed Behavior: Demo and engine tests evaluate strong vs weak test suites.
Observed Implementation:
- In `cmd/demo/main.go` and `tests/engine_test.go`, the "strong test function" is defined as `string(mutatedSrc) != string(src)`.
- This means any mutant generated is considered "killed" by the strong suite because the AST output differs textually from the original source.
- While the actual unit tests in `internal/service/discount_strong_test.go` independently verify comprehensive boundary assertions on `CalculateDiscount`, the demo runner's strong suite oracle relies on source text inequality rather than executing the test suite against a dynamic binary or interpreted engine.
Assessment: WARNING
Severity: MEDIUM
Notes: The engineering notes explicitly document this trade-off (`02-implementation-notes.md:49-57`). The core concept (demonstrating weak vs strong coverage assertion gap) is still proven, but the demo oracle is synthetic.
