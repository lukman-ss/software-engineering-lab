# Research Audit Plan — Lab 26: Contract Testing

## Target Lab
`labs/26-contract-testing`

## Pipeline Override Notice
- Research audit only.
- Implementation/code files are excluded from this audit stage per pipeline override.
- Research files are evaluated as-is without modification.
- Output directory: `labs/26-contract-testing/research-audit/`

## Files Reviewed
- `labs/26-contract-testing/research/01-plan.md`
- `labs/26-contract-testing/research/02-sources.md`
- `labs/26-contract-testing/research/03-evidence.md`
- `labs/26-contract-testing/research/04-contradictions.md`
- `labs/26-contract-testing/research/05-report.md`
- `labs/26-contract-testing/research/06-open-questions.md`
- Cross-reference file: `labs/06-api-versioning/README.md` (Source 7)

## Claims To Verify
1. Contract testing validates inter-application communication in isolation, not just domestic component correctness.
2. Consumer-Driven Contracts (CDC) express expectations driven by the consumer, verifying only used subset of interactions.
3. Code-based contract testing differs fundamentally from schema/OpenAPI testing (specification vs example/conversation, ambiguity vs concrete guarantees).
4. CI/CD integration models: build failure / deployment gating via contract verification.
5. Contract testing principles apply equally to asynchronous message/event architectures.
6. Taxonomies of breaking vs additive changes and classification of the lab's three changes (enum casing change, field rename, primitive type change).
7. Contract test scope: contracts need not be complete snapshots of payloads, only what consumers require.
8. Trade-offs around provider states and test data complexity.

## Code To Execute
- None (research-only audit phase per pipeline override).

## Primary Risks
- Verification of external citations against live sources for quote accuracy and context preservation.
- Discrepancy between Martin Fowler's classic CDC view (contract tests might not strictly break regular deployment builds immediately) vs modern tooling (Pact `can-i-deploy` blocking build).
- Verification of whether local references (Lab 06) are appropriate evidence for technical claims.
- Detection of unverified claims, overgeneralizations, or temporal anomalies.

## Audit Strategy
1. Live network verification of all cited web URLs via WebFetch tool.
2. Direct comparison of quotations and claims in `02-sources.md`, `03-evidence.md`, and `05-report.md` against authoritative source text.
3. Detailed analysis of source tiers, scope boundaries, and potential cherry-picking.
4. Evaluation of contradictions and open questions recorded by research agent.
5. Formal grading across quality gates and issuance of final verdict.
