# Research Plan

## Research Topic
Property-Based Testing (PBT) - Testing Invariants and Discovering Unforeseen Edge Cases

## Objective
To investigate the theory, practice, tools, and effectiveness of Property-Based Testing compared to traditional Example-Based Testing, with focus on invariant definitions, test case generation, shrinking mechanisms, and practical applications in software engineering.

## Research Questions

### Core Concepts
1. What are the theoretical foundations of Property-Based Testing (QuickCheck origin, algebraic specification)?
2. How does PBT differ fundamentally from Example-Based Testing in terms of test philosophy and coverage?
3. What constitutes a valid "property" or "invariant" suitable for PBT?

### Technical Mechanisms
4. How do test case generation algorithms work (random generation, coverage-guided, exhaustive for small domains)?
5. What is the shrinking process and how does it find minimal counterexamples?
6. How do different PBT frameworks implement generation and shrinking (QuickCheck, testing/quick, gopter, Hypothesis, fast-check, etc.)?

### Invariant Categories
7. What are the canonical categories of properties/invariants (roundtrip, idempotence, commutativity, associativity, monotonicity, equivalence/oracle)?
8. How to discover properties for a given function/system (property discovery techniques)?
9. What are common pitfalls in property formulation (over-specification, under-specification, flaky properties)?

### Practical Application
10. What are the best practices for integrating PBT into CI/CD pipelines?
11. How to handle stateful systems and side effects in PBT (state machine testing, model-based testing)?
12. What is the effectiveness of PBT in finding real bugs vs. example-based testing (empirical studies)?
13. How to combine PBT with other techniques (fuzzing, formal verification, metamorphic testing)?

### Tooling Ecosystem
14. Comparison of major PBT frameworks across languages (Go: testing/quick, gopter, gopter; Rust: proptest, quickcheck; Python: Hypothesis; JS/TS: fast-check, jqwik; Java: jqwik; Haskell: QuickCheck; Erlang: PropEr)
15. Integration with existing test infrastructure (test runners, coverage tools, property reporting)

## Search Strategy

### Primary Sources (Tier 1)
- Original QuickCheck papers (Claessen & Hughes, 2000, 2002, 2011)
- Official documentation of major frameworks (Hypothesis, fast-check, proptest, gopter, testing/quick)
- Academic papers on PBT effectiveness, shrinking algorithms, property discovery
- Language-specific standard library documentation (Go testing/quick, Rust proptest book)

### Secondary Sources (Tier 2)
- Reputable technical publications (ACM Queue, IEEE Software, CACM articles on PBT)
- Conference talks from Strange Loop, ICFP, Curry On, GopherCon, RustConf
- Framework-specific guides and tutorials by maintainers

### Community Sources (Tier 3) - for discovery only
- Blog posts by practitioners (e.g., Hillel Wayne, David MacIver, John Hughes)
- Reddit discussions (r/programming, r/haskell, r/rust, r/golang)
- Conference Q&A sessions

### Search Keywords
- "Property-Based Testing" "QuickCheck" "invariant" "shrinking" "test case generation"
- "metamorphic testing" "model-based testing" "stateful property-based testing"
- "Hypothesis" "fast-check" "proptest" "gopter" "testing/quick" comparative
- "property-based testing effectiveness" "bug finding" "empirical study"

## Expected Primary Sources

1. **QuickCheck: A Lightweight Tool for Random Testing of Haskell Programs** (Claessen & Hughes, 2000) - Foundational paper
2. **Testing Monadic Code with QuickCheck** (Claessen & Hughes, 2002)
3. **QuickCheck Specification-Based Testing** (Hughes, 2019) - Modern overview
4. **Hypothesis documentation** (David MacIver) - Advanced generation/shrinking
5. **proptest book** (Rust) - Comprehensive guide
6. **fast-check documentation** - Modern JS/TS implementation
7. **gopter GitHub + docs** - Go implementation
8. **Go testing/quick package docs** - Standard library approach
9. **PBT empirical studies** (e.g., "An Empirical Study on the Effectiveness of Property-Based Testing" - if exists)
10. **Metamorphic Testing papers** (Chen et al.) - Related technique

## Risks / Unknowns

1. **Limited empirical data**: Few large-scale controlled studies comparing PBT vs example-based bug detection rates
2. **Framework-specific behavior**: Shrinking algorithms differ significantly; results may not generalize
3. **Property discovery gap**: How engineers actually find good properties is under-documented
4. **Stateful testing complexity**: Model-based/stateful PBT has less standardization
5. **Language bias**: Most literature focuses on Haskell; Go/Rust/JS ecosystems have different ergonomics
6. **Maintenance burden**: Property maintenance cost vs. example test maintenance - limited data
7. **Flaky properties**: Non-determinism in generation/shrinking can cause CI flakiness - mitigation strategies needed
8. **Integration with formal methods**: Boundary between PBT and formal verification (e.g., Liquid Haskell, Dafny) unclear