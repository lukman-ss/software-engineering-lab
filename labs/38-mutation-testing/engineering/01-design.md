# Engineering Design

Target Lab: labs/38-mutation-testing
Research Status: APPROVED

## Concept To Prove
Prove that 100% line coverage does not guarantee fault detection ("false confidence"). Demonstrate how AST-based mutation operators inject faults, how tests catch (kill) or miss (survive) them, and how mutation score is computed:
$$\text{Mutation Score} = \left(\frac{\text{Killed Mutants}}{\text{Total Mutants}}\right) \times 100\%$$

## Expected Behavior
1. Given target business logic (e.g., banking/discount policy with boundary and boolean logic), a weak test suite achieves 100% line coverage but has weak assertions.
2. The mutation testing engine systematically generates 5 categories of mutants:
   - Relational operator replacement (`>`, `<`, `>=`, `<=`, `==`, `!=`)
   - Boolean flip (`&&` <-> `||`)
   - Arithmetic operator replacement (`+`, `-`, `*`, `/`)
   - Boundary/Value mutation (constants shifted $\pm 1$)
   - Statement deletion (void calls or assignments neutralized)
3. Running the weak test suite against all mutants results in surviving mutants (low mutation score despite 100% coverage).
4. Running an improved comprehensive test suite satisfies the RIP model (Reach, Infect, Propagate) and achieves a 100% mutation score by killing all non-equivalent mutants.

## Failure Scenario
- Weak tests pass against mutated code where logic is inverted, revealing undetected bugs.
- Equivalent mutants (if generated) are identifiable as undetectable by design.

## Success Criteria
- Compiles cleanly on Go 1.22+.
- AST mutator parses Go code, generates valid mutant ASTs, and executes test functions.
- Weak test suite passes on original code (100% coverage) but yields low mutation score (<50%).
- Strong test suite kills all non-equivalent mutants (100% mutation score).
- Concurrency-safe mutant test execution with zero race conditions (`go test -race ./...`).

## Architecture
```
labs/38-mutation-testing/
├── cmd/demo/main.go            # CLI demo runner showing weak vs strong suite mutation analysis
├── internal/
│   ├── engine/                 # In-memory Go AST mutation engine & runner
│   │   ├── mutator.go          # AST visitor & operator replacement logic
│   │   ├── runner.go           # Mutant test execution & score calculation
│   │   └── types.go            # Mutant, MutationType, Result definitions
│   └── service/                # Target business logic to test
│       ├── discount.go         # Tiered discount calculator logic
│       ├── discount_weak_test.go   # Weak assertions (100% coverage, low mutation score)
│       └── discount_strong_test.go # Comprehensive assertions (kills mutants)
├── engineering/
└── go.mod
```

## Components
1. `internal/engine`: AST parser and mutator inspecting Go AST expressions and binary operations; test runner executing test suites against mutated candidates.
2. `internal/service`: Concrete pricing/discount rules with relational, arithmetic, and boolean conditionals.
3. `cmd/demo`: Terminal output presenting mutation score comparison between weak and comprehensive suites.

## Test Strategy
- Unit tests for the mutation engine itself (verifying each operator transforms AST correctly).
- Unit tests for the target domain service (`discount.go`).
- Comparative test demonstrating weak suite vs strong suite kill rates.

## Execution Plan
1. Initialize Go module `labs/38-mutation-testing`.
2. Implement domain logic in `internal/service/discount.go`.
3. Implement `internal/engine` types, AST mutators, and runner.
4. Implement test suites (`weak` vs `strong`) and engine unit tests.
5. Implement `cmd/demo/main.go` demonstrating coverage vs mutation score.
6. Verify via `go test -v -race ./...` and `go run ./cmd/demo`.

## Implementation Decisions
- Implementation Decision: In Go, we use standard library `go/ast`, `go/parser`, `go/token` to mutate AST nodes in memory, allowing fast and safe mutation without requiring shell spawning or slow disk writes.
