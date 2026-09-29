# Audit Plan

## Target Lab
`labs/38-mutation-testing`

## Pipeline Mode
RESEARCH ONLY (per PIPELINE OVERRIDE instructions: "Audit research only. Do not audit implementation/code in this stage. Do not modify research files. Write all audit output to: labs/38-mutation-testing/research-audit/").

## Files Reviewed
1. `labs/38-mutation-testing/research/01-plan.md`
2. `labs/38-mutation-testing/research/02-sources.md`
3. `labs/38-mutation-testing/research/03-evidence.md`
4. `labs/38-mutation-testing/research/04-contradictions.md`
5. `labs/38-mutation-testing/research/05-report.md`
6. `labs/38-mutation-testing/research/06-open-questions.md`
7. `labs/38-mutation-testing/research-revision/01-revision-plan.md`
8. `labs/38-mutation-testing/research-revision/02-changes-made.md`
9. `labs/38-mutation-testing/research-revision/03-revision-result.md`

## Claims To Verify
- Code coverage measures execution only, not fault detection capability.
- Mutation Score formula: (Killed / Total) × 100%.
- Core theoretical hypotheses: Competent Programmer Hypothesis and Coupling Effect.
- Three conditions of the RIP model (Reach, Infect, Propagate) for mutant killing.
- Equivalent mutants: definition, theoretical undecidability, and heuristic practical approaches.
- Mutation operators taxonomy (statement, value, decision, operator replacement).
- Meta ACH system findings: LLM equivalence detection metrics and empirical trial metrics.
- Go mutation testing tooling status: `go-mutesting` and `gremlins` availability vs JVM/JS tool maturity.
- Subsumed mutants definition and impact on metrics.
- Industry standard thresholds: absence of universal threshold standard.

## Code To Execute
None. PIPELINE OVERRIDE: Audit research only. Implementation code not present or out of scope for research-stage audit.

## Primary Risks
- Inaccessible or dead URLs (e.g. Fowler bliki draft 404 or path shifting).
- Foundational papers unread directly (DeMillo et al. 1978, Jia & Harman 2009 cited via secondary bibliography).
- Single-vendor industrial findings (Meta ACH) presented as universal facts if not properly scoped.
- Premature publication status for preprints (arXiv:2501.12862 submitted to FSE 2025).

## Audit Strategy
1. Live network verification of all referenced URLs.
2. Fact vs interpretation claim classification across research files.
3. Analysis of citations for secondary attribution disclosure.
4. Evaluation of prior revisions against identified research gaps.
5. Verdict determination based on evidence completeness.
