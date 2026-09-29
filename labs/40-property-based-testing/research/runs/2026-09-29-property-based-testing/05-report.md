# Research Report: Property-Based Testing (PBT)

## Research Question
How does Property-Based Testing (PBT) verify domain invariants, systematically discover unhandled edge cases via test-case shrinking, and how should it be structured in Go systems using `testing/quick` and `gopter`?

## Executive Summary
Example-Based Testing tests known points in input space, leaving blind spots at boundaries, edge cases, and combinations. Property-Based Testing specifies invariant properties that must hold across generated inputs. When a violation occurs, shrinking simplifies large counterexamples to their Minimal Failing Example (MFE). In Go, developers have three options: standard library `testing/quick` (minimal, reflection-based), `gopter` (extensible generators and shrinkers), and `testing.F` (coverage-guided fuzzing).

## Findings

### Finding 1: Core Mathematical Properties
Claim: Invariants fall into standard algebraic categories:
1. **Roundtrip / Invertibility**: `Decode(Encode(x)) == x`.
2. **Idempotence**: `f(f(x)) == f(x)`.
3. **Invariant Preservation**: `len(sort(xs)) == len(xs)` and elements are monotonically non-decreasing.
4. **Metamorphic / Test Oracle**: `fast_algorithm(x) == naive_algorithm(x)`.
5. **Hard to Prove, Easy to Verify**: Computing a graph coloring or optimal packing is NP-hard, but verifying feasibility is linear/polynomial.

Evidence: Claessen & Hughes (2000), Scott Wlaschin (2014).
Sources: QuickCheck paper, F# for Fun and Profit.
Confidence: HIGH

### Finding 2: Shrinking Mechanics
Claim: Shrinking isolates minimal failing inputs by applying deterministic reductions:
- Slices: bisection, element removal, recursive sub-slice reduction.
- Integers: binary search towards zero (`x / 2`, `x - 1`).
- Strings: character deletion, lowering character ranges.

Evidence: Gopter repo / QuickCheck.
Sources: `leanovate/gopter`.
Confidence: HIGH

### Finding 3: Framework Comparison in Go
Claim:
- `testing/quick`: Stdlib, no dependencies, but rigid generators and limited shrinking for custom structs.
- `gopter`: Full QuickCheck port with `gen.SliceOf`, `gen.IntRange`, custom shrinking trees (`shrink.Shrinker`), and condition filters.
- `testing.F`: Go 1.18+ native coverage-guided engine, ideal for byte/primitive boundaries, but lacks composable domain object generators.

Evidence: Go docs (`pkg.go.dev/testing/quick`, `go.dev/doc/security/fuzz`), Gopter docs.
Sources: Official Go Documentation, Gopter.
Confidence: HIGH

### Finding 4: Domain Case Studies (Currency & Interval Merger)
Claim:
1. **Currency Formatter**: Fails on negative zero (`-0.00`), floating point sub-cent representation errors (`0.1 + 0.2 != 0.3`), int64 overflow when parsing cents, localized grouping separators.
2. **Interval Merger**: Fails on inverted intervals `[end < start]`, open/closed boundary overlaps `[1, 3]` and `[3, 5]`, adjacent intervals, empty sets, and extreme ranges `[MinInt, MaxInt]`.

Evidence: Domain testing patterns in financial systems and interval scheduling.
Sources: Industry benchmarks and PBT case studies.
Confidence: HIGH

## Areas of Agreement
- Example-based tests are insufficient for complex invariants.
- Shrinking is essential to make counterexamples actionable for engineers.
- Invariants must be independent of implementation details.

## Areas of Disagreement
- Whether coverage-guided fuzzing makes classical generator-based PBT obsolete in Go. (Verdict: No, they serve complementary purposes).

## Limitations
- PBT tests can be slower than unit tests due to multi-thousand iteration cycles.
- Weakly constrained generators can waste cycles generating invalid inputs without reaching deep code branches.

## Conclusion
Property-Based Testing elevates testing from static verification of known cases to specification-level invariant validation. For Go engineering, `gopter` offers the most capable framework for complex domain invariants and custom shrinking, while `testing/quick` suffices for primitive algebraic functions.
