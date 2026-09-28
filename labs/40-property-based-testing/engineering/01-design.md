# Engineering Design

Target Lab: labs/40-property-based-testing
Research Status: APPROVED

## Concept To Prove
Property-Based Testing (PBT) vs. Example-Based Testing:
1. Example-based tests rely on hand-picked input/output cases that miss unconsidered boundary conditions, edge cases, and arithmetic pitfalls.
2. Property-Based Testing defines universal invariants across input domains and verifies them across hundreds or thousands of randomized inputs.
3. Shrinking mechanics systematically reduce complex failing counterexamples to minimal reproducers.
4. Three canonical invariant categories:
   - Roundtrip: `Decode(Encode(x)) == x`
   - Idempotence: `f(f(x)) == f(x)`
   - Equivalence / Test Oracle: `f_optimized(x) == f_reference(x)`

## Expected Behavior
- Example-based test suite passes cleanly on flawed implementations by only testing "happy" and known examples.
- Property-based test suite fails on flawed implementations by exploring randomized spaces and boundary values (e.g., negative numbers, sub-cent fractions, large ints, empty inputs).
- A counterexample shrinker isolates the minimal input causing invariant failure.
- Corrected implementations pass all properties across 1,000 randomized iterations without regression.

## Failure Scenario
- Naive Currency Formatter / Parser:
  - Fails roundtrip on floating-point precision loss (e.g. `$0.07 * 100` rounding issues or negative currency formatting).
- Naive Interval Merging:
  - Fails idempotence or equivalence on unsorted, overlapping, or single-point intervals when edge conditions are unhandled.
- Run-Length Encoder (RLE):
  - Fails roundtrip on digits inside raw text or empty strings without proper escaping.

## Success Criteria
1. Domain modules implemented with both flawed (naive) and verified (robust) variants or clear switchable bugs.
2. Example-based unit tests demonstrate false confidence by passing naive implementations.
3. Property-based test suite (using Go stdlib `testing/quick` and custom generator + shrinker) catches bugs in naive implementations and passes robust implementations across 1,000 iterations.
4. Minimal counterexample shrinking demonstrated with trace logs showing input reduction.
5. All tests pass (`go test -race ./...`) and demo (`go run ./cmd/demo`) executes deterministically with detailed stdout report.

## Architecture
```text
labs/40-property-based-testing/
├── go.mod
├── README.md
├── cmd/
│   └── demo/
│       └── main.go
├── internal/
│   ├── currency/
│   │   ├── currency.go
│   │   └── currency_test.go
│   ├── interval/
│   │   ├── interval.go
│   │   └── interval_test.go
│   └── shrinker/
│       ├── shrinker.go
│       └── shrinker_test.go
└── engineering/
    ├── 01-design.md
    ├── 02-implementation-notes.md
    └── 03-execution-result.md
```

## Components
- `internal/currency`:
  - Models monetary amounts as integer cents to avoid IEEE 754 precision errors.
  - Naive implementation uses `float64` leading to precision loss.
  - Roundtrip invariant: `ParseCents(FormatCents(c)) == c`.
- `internal/interval`:
  - Merges intervals `[start, end]`.
  - Naive implementation misses unordered slices or boundary merges.
  - Idempotence invariant: `Merge(Merge(intervals)) == Merge(intervals)`.
  - Oracle invariant: `MergeOptimized(intervals) == MergeReference(intervals)`.
- `internal/shrinker`:
  - Demonstrates counterexample shrinking algorithm: binary search and list element elimination to reduce a failing array/integer to minimal reproducible input.
- `cmd/demo`:
  - Runs side-by-side comparison of Example-Based Testing vs. Property-Based Testing.
  - Runs shrinking demonstration showing step-by-step reduction of failure.

## Test Strategy
- Unit tests (`*_test.go`) in each package:
  - Example-based tests showing what typical unit tests cover.
  - Property-based tests using `testing/quick` and custom property checks.
  - Shrinker unit tests verifying convergence to minimal counterexample.
- Race detector check with `go test -race ./...`.

## Execution Plan
1. Initialize Go module: `go.mod` (Go 1.22+).
2. Implement `internal/currency` (naive float vs robust int cents, parser, formatter).
3. Implement `internal/interval` (naive vs robust interval merge, idempotence & oracle).
4. Implement `internal/shrinker` (counterexample search & binary/element shrinking).
5. Add comprehensive unit & property tests.
6. Implement `cmd/demo/main.go`.
7. Execute `go test -race ./...` and `go run ./cmd/demo`.
8. Document implementation notes and execution results.

## Implementation Decisions
- Implementation-Specific Decision: Use Go standard library (`testing/quick`, `math/rand/v2` or `math/rand`, `reflect`) with zero third-party dependencies to maximize standard library leverage, maintainability, and reproducibility per workspace guidelines.
