# Open Questions & Next Research Directions

## Unanswered Questions
1. **Performance Cost in CI/CD**: What is the optimal iteration count (e.g., 100 vs 1,000 vs 10,000 runs) balancing test duration against rare bug discovery probability in PR pipelines?
2. **Stateful / Model-Based Testing**: How does Gopter handle state machine verification (Commands pattern) when testing stateful backend systems like transactional databases or distributed caches?

## Weak Evidence / Claims Needing Deeper Research
- Exact benchmark overhead between reflection-based `testing/quick` vs combinator-based `gopter` under heavy allocations.

## Next Research Directions for Implementation Phase
- Structuring an Interval Merger exercise with invariants:
  1. Non-overlapping result: $\forall i \neq j, \text{interval}_i \cap \text{interval}_j = \emptyset$.
  2. Complete coverage: $\bigcup \text{inputs} = \bigcup \text{outputs}$.
  3. Minimality: Output length is minimal.
- Structuring a Financial Currency Formatter exercise testing string roundtrip and decimal precision invariants without floating-point drift.
