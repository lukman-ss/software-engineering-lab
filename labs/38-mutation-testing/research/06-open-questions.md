# Open Questions

## Unanswered Questions

1. **Go mutation testing tooling gap**: What is the current state of mutation testing tools for Go (e.g., GoMutator, mull)? How do they compare to PIT/Stryker in terms of operator coverage and performance? NOT VERIFIED.

2. **Mutation score vs defect correlation**: What is the empirical correlation between mutation score and real-world defect density? Does a 90%+ mutation score reliably predict fewer production bugs? NOT VERIFIED.

3. **Optimal mutation operator selection**: Which mutation operators provide the highest ROI (kill rate per computational cost)? Are some operators consistently ineffective across languages? NOT VERIFIED.

4. **LLM-based mutation testing reliability**: How reliable is ACH's LLM-based test generation compared to human-written tests? What are the false positive/negative rates for the generated tests? Limited evidence from Meta trial only; no independent verification.

5. **Mutation testing in microservices**: How does mutation testing apply to distributed systems where failures manifest across service boundaries rather than within a single service? NOT VERIFIED.

## Weak Evidence

1. **Meta ACH statistics**: The 73% acceptance rate and 36% privacy relevance claims come from a single industry blog post (Meta Engineering, September 2025). No independent verification or peer-reviewed publication accessed. Confidence is MEDIUM at best.

2. **Mutation score thresholds**: No consensus found on what mutation score threshold constitutes "good enough" — sources did not agree on a standard target (80%, 85%, 90% variously proposed elsewhere but NOT VERIFIED here).

3. **Fowler bliki draft status**: Martin Fowler's mutation testing page is explicitly marked as DRAFT at the time of access. It should not be cited as finalized evidence until the author removes that notice.

## Claims Needing Deeper Research

1. **Statement mutation effectiveness**: The paper "Empirical Evaluation of the Statement Deletion Mutation Operator" (Deng et al., 2013) was cited in Wikipedia but NOT opened. Specific findings about statement deletion's kill rate remain uncovered.

2. **Higher-order mutants**: Wikipedia references higher-order mutants (mutants with more than one mutation) as supporting the coupling effect, but empirical data on their effectiveness is NOT VERIFIED here.

3. **Weak vs strong mutation trade-offs**: Practical guidance on when to use weak mutation (Reach+Infect) vs strong mutation (RIP) is not well-documented in the sources that were successfully fetched.

4. **Subsumed mutants and trivial/new mutation score refinements**: Wikipedia mentions subsumed mutants and many secondary sources discuss variants such as New Mutation Score — these threads were not fully explored.

## Possible Next Research Directions

1. **Build a Go mutation testing harness**: Following the lab specification, implement a simple mutation testing framework for Go functions (pricing engine or credit scoring logic).

2. **Compare mutation scores across test styles**: Run mutation testing on codebases with TDD tests vs. non-TDD tests to measure the mutation score difference.

3. **Investigate equivalent mutant heuristics for Go**: Evaluate whether LLM equivalence detection (as in ACH) can be ported to Go mutation workflows.

4. **CI/CD integration study**: Evaluate how mutation testing integrates into real CI/CD pipelines — build time impact, flakiness, and developer adoption patterns.
