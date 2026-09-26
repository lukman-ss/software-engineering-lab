# Docs vs Code Analysis

## Comparisons

1. **README vs Code**
   - README claims automated validation of structural invariants, monotonic numbering, and supersession lineage.
   - Code correctly implements these features in `internal/adr/linter.go`.
   - Result: MATCH.

2. **Engineering Notes vs Code**
   - Notes state regex is used instead of AST.
   - Code verifies this choice in `internal/adr/parser.go`.
   - Notes claim in-memory graph validation is concurrent.
   - Code verifies Goroutines and WaitGroups are correctly utilized.
   - Result: MATCH.

3. **Research Claims vs Implementation**
   - Research outlines the ADR cyclical lifecycle (`Proposed` -> `Accepted` -> `Superseded`).
   - Implementation uses these statuses in `internal/adr/models.go` and validates their invariants.
   - Demo code explicitly maps out the "Monolith First" to "Microservices" architectural progression recommended in the research.
   - Result: MATCH.

4. **Tests vs Claims**
   - Tests explicitly validate all constraints defined in the engineering design.
   - Result: MATCH.

## Conclusion

No significant discrepancies were found between the documentation, research claims, and the actual implementation. The code strictly adheres to the stated boundaries and features.