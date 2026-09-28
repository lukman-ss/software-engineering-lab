# Open Questions

## Unanswered Questions

1. **What is the quantitative bug-detection improvement of PBT over example-based testing?** No controlled benchmark or A/B study was found that measures bug discovery rate per hour of test-writing time. The Hypothesis Corpus could be mined for failure rates, but no published analysis exists yet.

2. **How do engineers discover good properties for a new function?** The lab spec says "3 invariant utama" — but how to find them? Property discovery methods (reverse-engineering from docstrings, metamorphic relations, algebraic specifications) are under-documented.

3. **What is the maintenance cost of property suites vs. example-based test suites?** Do PBT properties age well (fewer false positives over time), or do they require constant adjustment as the codebase evolves? No empirical study was found.

4. **How effective is PBT for stateful systems vs. pure functions?** The Hypothesis Corpus shows stateful tests exist but are rare. What property patterns work best for stateful code?

5. **What is the optimal number of test iterations for CI?** 100 is the default everywhere; is 1,000 (as in the lab spec) better? Does CI time justify higher iteration counts?

## Weak Evidence

1. **PBT effectiveness in Go specifically**: Most empirical evidence comes from Python (Hypothesis), JavaScript (fast-check), and Haskell (QuickCheck). Go has fewer documented case studies (gopter README lists stars but no bug-findings catalog comparable to fast-check's track record). The lab targets Go (`testing/quick` / `gopter`), so effectiveness claims for Go are weaker.

2. **"Designed for bugs" claims**: fast-check claims its generators are "designed for bugs" with boundary upweighting. However, no independent verification was found of how much more effective this bias is vs. uniform random generation.

3. **Hypothesis Corpus findings**: The dataset was collected in October 2025 and published on HuggingFace in 2026. No peer-reviewed analysis of the corpus was found at time of research. Claims about test patterns and failure rates from this dataset are therefore based on the dataset description, not published analysis.

## Claims Needing Deeper Research

1. **Byte-stream vs. value-tree shrinking**: Hypothesis's byte-stream approach is claimed to avoid local minima better than QuickCheck's type-based shrinking. No formal comparison study was found.

2. **Distribution control**: Hypothesis explicitly refuses to let users control distribution. Is this a sound principle, or does it prevent useful targeting in some domains?

3. **Shrinking guarantees**: Is there a formal guarantee that shrinking always terminates? Hypothesis and proptest rely on well-founded orderings; QuickCheck's `shrink :: a -> [a]` must return finite lists for termination. What happens when shrinkers are poorly written?

4. **Integration with formal verification**: The boundary between PBT and tools like Liquid Haskell, Dafny, or TLA+ is unclear. When does PBT become "formal enough"?

## Possible Next Research Directions

1. **Mine the Hypothesis Corpus** for: failure rate by property pattern type, most common strategies, stateful test proportion, median test case count to failure.
2. **Compare shrinking quality** across Hypothesis (byte-stream), proptest (value-tree), and QuickCheck (type-based) on the same buggy function.
3. **Build an empirical benchmark** of Go PBT (testing/quick vs. gopter) on real Go codebases to measure bug-finding effectiveness in the Go ecosystem.
4. **Study property discovery heuristics**: What techniques do experienced PBT users use to identify useful invariants? Can they be taught or automated?
5. **Investigate flakiness mitigation**: How do frameworks handle non-deterministic properties (e.g., map iteration order, floating-point rounding)?