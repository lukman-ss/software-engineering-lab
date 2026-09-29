# Research Topic
Property-Based Testing (PBT): Invariant Verification, Shrinking Mechanisms, and Edge Case Discovery in Go (`testing/quick` and `gopter`).

# Objective
Provide rigorous research backing, foundational theory, formal invariant patterns, framework comparison, shrinking algorithmic details, and edge case dynamics to guide senior engineering implementation without writing production code.

# Research Questions
1. What are the formal foundations of Property-Based Testing (Claessen & Hughes QuickCheck) vs Example-Based Testing and coverage-guided fuzzing?
2. How does test-case shrinking work algorithmically (integrated shrinking vs separate shrinkers)?
3. What are the canonical property patterns (Roundtrip, Idempotence, Metamorphic/Oracle, Invariant preservation) and their mathematical formulations?
4. How do Go implementations compare (`testing/quick` standard library vs `leanovate/gopter` vs Go native fuzzing `testing.F`) in terms of generator composability, shrinking capability, and type constraints?
5. What are the common failure modes and domain edge cases in financial currency formatting and interval merging under PBT?

# Search Strategy
- Primary: Academic papers (Claessen & Hughes 2000), Go official documentation (`testing/quick`, Go fuzzing design doc).
- Tier 1/2: Reputable industry papers/engineering specifications (Hypothesis documentation, Gopter repository and docs, Scott Wlaschin's Property-Based Testing design patterns).
- Cross-check: Compare standard library limitations against dedicated PBT frameworks.

# Expected Primary Sources
- Claessen, K., & Hughes, J. (2000). QuickCheck: a lightweight tool for random testing of Haskell programs.
- Go standard library documentation: `testing/quick`.
- Gopter documentation and source repository: `github.com/leanovate/gopter`.
- Go design document: Fuzzing is Beta Ready (Go 1.18).

# Risks / Unknowns
- Standard library `testing/quick` is legacy and lacks built-in custom shrinkers for composite types compared to `gopter`.
- Differences between PBT generators and coverage-guided fuzz engines (libFuzzer/Go fuzzing) in invariant exploration.
