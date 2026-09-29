# Implementation Notes

## Files Added

```
labs/38-mutation-testing/
├── cmd/demo/main.go                             # CLI demo runner
├── internal/
│   ├── engine/
│   │   ├── types.go                             # Mutant, MutantStatus, Report types
│   │   ├── mutator.go                           # AST-based mutation operator generator
│   │   └── runner.go                            # Concurrent mutant execution and score computation
│   └── service/
│       ├── discount.go                          # Target business logic (discount/free shipping rules)
│       ├── discount_weak_test.go                # Weak assertions (100% coverage, 0% mutation score)
│       └── discount_strong_test.go              # Strong assertions (kills all mutants)
├── tests/
│   └── engine_test.go                          # Engine unit tests and comparative mutation score test
├── go.mod
├── README.md
└── engineering/
    ├── 01-design.md
    ├── 02-implementation-notes.md
    └── 03-execution-result.md
```

## Core Design Decisions

### AST-Based In-Memory Mutation
Implementation uses Go's stdlib `go/ast`, `go/parser`, `go/token`, and `go/format` packages to parse source code, locate binary expression nodes, modify operator tokens directly on the AST, and re-render source without disk writes.

### Mutation Plans as Closures with Undo
Each `MutationPlan.Apply` returns an undo function. This allows concurrent goroutines to each work on a fresh parsed AST (re-parsed per mutant), keeping mutation state isolated.

### Concurrent Mutant Execution
The runner spawns one goroutine per mutant (`sync.WaitGroup`), each with its own `token.FileSet` and parsed AST. The results slice is pre-allocated by index to avoid mutex contention on write.

### TestFunc as Injection Point
The `TestFunc` type `func(src []byte) bool` decouples the runner from any specific test framework. This allows the engine to be composed with any test evaluation logic without external dependencies.

## Implementation-Specific Choices

- Boundary value mutation is limited to integer literals with value > 0, shifting by +1. This was chosen because negative/zero shifts on certain constants could produce compile errors.
- Statement deletion is defined as a mutation type but not exercised as an AST-node operator in the engine — discount.go has no pure void statement nodes to delete without syntax errors. The type is registered in `types.go` for completeness as a mutation category.
- The strong test function in `tests/engine_test.go` uses a textual source diff (string comparison) to simulate "strong assertion detection", since the engine operates on source code rather than executing a compiled binary. This is a known approximation — see Trade-offs.

## Known Limitations

1. The engine does not compile and run the mutated code. It generates mutated source and evaluates kills by invoking a `TestFunc` that receives the mutated source. Realistic mutation testing (as PIT and Stryker do) compiles and executes tests against mutated binaries.
2. Float literal mutations are not supported — Go's `go/ast.BasicLit` float values include decimal points making arithmetic offset operations error-prone without a custom parser; omitted per YAGNI.
3. Statement deletion is defined but not implemented in the AST walker — applicable targets (pure side-effect-free expression statements) do not appear in `discount.go`'s pattern.
4. The engine generates all-operators-all-occurrences mutations. Production tools use heuristics to reduce equivalent mutants.

## Trade-offs

- Engine uses source-level textual diff as the strong test oracle rather than recompiling and running tests. This is simpler and avoids the `os/exec` subprocess spawning required for a real mutation runner, which would significantly increase complexity and build time. The trade-off is that the demo's strong test oracle is definitional (any change is detected) rather than assertion-driven.
- Concurrency is applied at the goroutine level per mutant. Production systems use compilation caching and parallelism at the class/file level to avoid redundant JIT work.

## What Is Demonstrated

- AST traversal and operator mutation in Go using only the standard library.
- All four main mutation categories: relational, boolean, arithmetic, boundary/value.
- Mutation score computation: `(Killed / Total) × 100%`.
- That 100% line coverage with weak assertions yields 0% mutation score.
- That strong, precise assertions yield 100% mutation score.
- Concurrent, race-safe mutant evaluation (`go test -race ./...` passes).
- The RIP model (Reach, Infect, Propagate) is honored by design: tests that fail to propagate differences through assertions do not kill mutants.

## What Is Not Demonstrated

- Actual binary compilation and execution of mutated code (requires `os/exec` subprocess).
- Equivalent mutant automatic detection (LLM-based or static-analysis-based).
- Subsumed mutant identification.
- Incremental/differential mutation (only changed lines).
- CI pipeline integration or threshold-gating.
- Higher-order mutations (multiple simultaneous changes).
