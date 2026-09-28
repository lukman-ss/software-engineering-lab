# Implementation Notes

## Files Added
- `go.mod`: Go module specification (Go 1.22+ standard stdlib implementation).
- `internal/currency/currency.go`: Implementation of `NaiveCurrency` (float64 dollar formatter/parser) and `RobustAmount` (integer cents precision storage).
- `internal/currency/currency_test.go`: Example-based tests passing naive currency, and PBT verifying roundtrip invariant `ParseRobust(Format(amount)) == amount` with custom generator.
- `internal/interval/interval.go`: Implementation of `NaiveMerge` (fails on unsorted intervals) and `RobustMerge` (sort-based merge algorithm).
- `internal/interval/interval_test.go`: Idempotence invariant `Merge(Merge(x)) == Merge(x)` and non-overlapping invariants verified with `testing/quick`.
- `internal/shrinker/shrinker.go`: Binary sectioning and element-magnitude reduction shrinking algorithm.
- `internal/shrinker/shrinker_test.go`: Unit test verifying convergence to minimal single-element counterexample `[-1]`.
- `cmd/demo/main.go`: Execution binary running side-by-side comparison of Example-Based vs PBT and outputting detailed shrinking traces.

## Core Design Decisions
1. **Standard Library Only (`testing/quick`)**: Zero third-party dependencies used to showcase pure Go standard library PBT capabilities.
2. **Custom Type Generators**: Implemented `Generate(r *rand.Rand, size int) reflect.Value` for `RobustAmount`, `NaiveFloat`, and `IntervalSlice` to ensure boundary cases (0, negative, unsorted, wide ranges) are exercised.
3. **Multi-Strategy Counterexample Shrinking**: Created explicit shrinking logic combining chunk removal (`remove-right-half`, `remove-left-half`), element removal (`remove-elem`), and binary integer reduction toward zero (`reduce-val-idx`).

## Implementation-Specific Choices
- Implementation Decision: Interval inputs cover negative coordinates (`-20` to `100`) to test range boundaries.
- Implementation Decision: Integer cents stored as `int64` to handle multi-billion dollar amounts without overflow.

## Known Limitations
- Go stdlib `testing/quick` is officially frozen; advanced features like regex generation or stateful state-machine testing require third-party tools like `gopter` in production environments.
- Shrinking algorithm in `internal/shrinker` is specific to integer slice invariants; general reflection-based shrinking across arbitrary nested structs requires reflection trees.

## Trade-offs
- Random generation adds ~300-500ms to test suite execution compared to instant hand-picked unit tests, but uncovers non-obvious bugs.

## What Is Demonstrated
- Roundtrip Invariant (`Decode(Encode(x)) == x`).
- Idempotence Invariant (`f(f(x)) == f(x)`).
- Oracle / Equivalence Invariant (`f_naive(x) == f_robust(x)`).
- Automatic counterexample shrinking from 10 elements down to a minimal single element reproducer `[-1]`.

## What Is Not Demonstrated
- Stateful model-based testing (requires state-machine generators).
