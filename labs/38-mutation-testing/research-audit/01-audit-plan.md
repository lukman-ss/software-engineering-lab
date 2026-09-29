# Research Audit Plan: Mutation Testing

## Target Lab
`labs/38-mutation-testing` (Research Audit Only)

## Files Reviewed
- `labs/38-mutation-testing/research/01-plan.md`
- `labs/38-mutation-testing/research/02-sources.md`
- `labs/38-mutation-testing/research/03-evidence.md`
- `labs/38-mutation-testing/research/04-contradictions.md`
- `labs/38-mutation-testing/research/05-report.md`
- `labs/38-mutation-testing/research/06-open-questions.md`

## Claims To Verify
1. First proposed by Richard Lipton in 1971; published by DeMillo, Lipton, and Sayward in 1978.
2. Two foundational hypotheses: Competent Programmer Hypothesis and Coupling Effect.
3. Mutation Score definition and formula: $(Killed / Total) \times 100\%$.
4. 100% code coverage does not guarantee fault detection capability; mutation testing detects weak assertions.
5. Standard mutation operators across statement, value, and decision categories.
6. Equivalent mutant problem is mathematically undecidable; represents major practical obstacle.
7. RIP model (Reach, Infect, Propagate) criteria for mutant detection.
8. Common tool execution pattern: inject mutations $\to$ run tests $\to$ classify killed vs survived.
9. Meta ACH (September 2025) LLM-guided test generation and equivalence detection performance metrics.
10. Tooling landscape: PIT (JVM), Stryker (JS/TS/C#/Scala), and Go tools (`go-mutesting`, `gremlins`).
11. Subsumed mutants definition and impact on test coverage.
12. Absence of an industry-standard universal mutation score threshold.

## Primary Risks
- Broken or dead URLs (e.g. 404 on unverified Martin Fowler bliki link).
- Citation of unverified academic papers without direct inspection.
- Attribution of recent industry statistics (Meta ACH) without checking preprint/primary data.
- Overgeneralization of JVM/JS tool capabilities to Go ecosystem.

## Audit Strategy
- Step 1: Inventory all claims, sources, and gaps in research files.
- Step 2: Validate reachability, publisher identity, relevance, and fidelity of all 13 listed sources via live web fetches.
- Step 3: Audit extraction of 14 key claims against verified sources.
- Step 4: Audit contradiction analysis for thoroughness and logical coherence.
- Step 5: Detail research gaps and evaluate whether existing caveats are sufficient.
- Step 6: Produce final verdict according to independence rules and quality gates.
