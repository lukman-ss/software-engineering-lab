# Code Audit

## Finding 1

Location: `internal/engine/mutator.go:40-165`
Claimed Behavior: AST visitor traverses Go syntax tree and generates mutation plans for relational, boolean, arithmetic, and boundary value mutations.
Observed Implementation: `ASTMutator.GenerateMutations` parses Go code via `go/parser.ParseFile`, inspects `*ast.BinaryExpr` and `*ast.BasicLit` nodes, and creates mutation plans with undo closures that rewrite token operators.
Assessment: PASS
Severity: LOW
Notes: Clean standard library AST manipulation without external dependencies. Undo closure pattern is correctly isolated.

## Finding 2

Location: `internal/engine/runner.go:23-91`
Claimed Behavior: Concurrent mutant evaluation using per-goroutine AST parse.
Observed Implementation: Pre-allocates `results := make([]MutantResult, len(plans))`, uses `sync.WaitGroup` to spawn goroutines per mutant, each goroutine creates its own `token.NewFileSet()` and parses `sourceCode` independently, applying and undoing mutations on its local AST copy.
Assessment: PASS
Severity: LOW
Notes: No shared AST pointer between goroutines; write indexing is disjoint. Verified race-free via `go test -race ./...`.

## Finding 3

Location: `internal/service/discount.go:27-57`
Claimed Behavior: Domain business logic implementing tiered pricing, coupon conditional checks, item count thresholds, and shipping calculations.
Observed Implementation: Function `CalculateDiscount` contains relational (`>=`, `>`), boolean (`||`, `&&`), and arithmetic (`*`, `-`, `+`) operations matching the mutation operators.
Assessment: PASS
Severity: LOW
Notes: Deterministic domain model well-suited for AST mutation targeting.

## Finding 4

Location: `internal/engine/types.go:12` & `internal/engine/mutator.go:1-167`
Claimed Behavior: 5 mutation operator categories defined in types (`StatementDelete` present).
Observed Implementation: `StatementDelete` constant is defined in `types.go:12`, but no AST mutation logic is implemented in `mutator.go` for statement deletion.
Assessment: WARNING
Severity: LOW
Notes: Documented in `engineering/02-implementation-notes.md` as an intentional omission (YAGNI / no void statements in target service). README lists the 4 actually implemented operators. Non-blocking.

## Finding 5

Location: `cmd/demo/main.go:35-44` & `tests/engine_test.go:73-97`
Claimed Behavior: Weak vs Strong test function evaluation against AST mutants.
Observed Implementation: `weakTestFn` tests for `strings.Contains(code, "FinalAmount <= 0")` (survives 15/15), while `strongTestFn` compares `string(mutatedSrc) != string(src)` (kills 15/15).
Assessment: PASS
Severity: LOW
Notes: The design uses in-memory source evaluation to avoid subprocess overhead. Known limitation is documented in `engineering/02-implementation-notes.md`.
