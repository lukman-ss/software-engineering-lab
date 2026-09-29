# Code Audit

Target Lab: labs/38-mutation-testing

## Finding 1

Location: internal/engine/runner.go:25-94
Claimed Behavior: Safe concurrent mutant evaluation without AST cross-contamination.
Observed Implementation: Runner re-parses fresh AST per mutant goroutine (`parser.ParseFile(fset, "source.go", sourceCode, 0)`), isolates mutated source rendering to separate `token.FileSet`, pre-allocates result slice by index, and waits with `sync.WaitGroup`.
Assessment: PASS
Severity: LOW
Notes: `sync.Mutex` on `Runner` struct is unused, but runner operations are entirely concurrent-safe and free from data races.

## Finding 2

Location: internal/engine/mutator.go:40-164
Claimed Behavior: Generates mutations for relational, boolean, arithmetic, and boundary values.
Observed Implementation: AST visitor traverses binary expressions (`token.GEQ`, `token.GTR`, `token.EQL`, `token.LOR`, `token.LAND`, `token.SUB`, `token.MUL`) and basic literals (`token.INT` > 0). Correctly applies AST transformations and returns working undo functions.
Assessment: PASS
Severity: LOW
Notes: Covers 4 primary mutation categories effectively.

## Finding 3

Location: internal/engine/types.go:12 vs internal/engine/mutator.go
Claimed Behavior: Statement deletion operator defined in design and enum.
Observed Implementation: `StatementDelete` is declared in `types.go:12` and referenced in `01-design.md`, but not implemented in `mutator.go`.
Assessment: WARNING
Severity: LOW
Notes: Documented as an intentional omission in `02-implementation-notes.md` due to absence of standalone void statements in `discount.go`. README.md correctly lists only the 4 implemented operators.

## Finding 4

Location: cmd/demo/main.go:35-43, tests/engine_test.go:73-98
Claimed Behavior: Evaluates weak test suite vs strong test suite effectiveness.
Observed Implementation: `weakTestFn` checks `strings.Contains(code, "FinalAmount <= 0")` and `strongTestFn` compares source diff `string(mutatedSrc) != string(src)`.
Assessment: WARNING
Severity: MEDIUM
Notes: The engine simulates test suite detection capability via source evaluation functions rather than compiling and executing mutated binaries out-of-process. While documented in `02-implementation-notes.md`, it represents an oracle approximation.

## Finding 5

Location: internal/service/discount.go:27-56
Claimed Behavior: Core domain pricing logic with multi-branch logic for testing.
Observed Implementation: Implements realistic business logic with customer tiers (VIP, Premium, Standard), threshold comparisons, coupon modifiers, and shipping eligibility. Clean, deterministic, and idiomatic.
Assessment: PASS
Severity: LOW
Notes: Provides balanced targets for relational, boolean, arithmetic, and boundary mutation operators.
