# Research Audit Plan

Target Lab: labs/38-mutation-testing
Audit Scope: Research Files Only (PIPELINE OVERRIDE)
Audit Directory: labs/38-mutation-testing/research-audit/

## Target Lab Overview

The target lab explores **Mutation Testing**, focusing on evaluating test suite quality beyond line/branch code coverage metrics, identifying unverified code paths (false confidence from 100% code coverage), understanding mutation operators, RIP model theoretical foundations, equivalent mutants, practical tooling (PIT, Stryker, go-mutesting, gremlins), and LLM-guided mutation testing (Meta ACH).

## Files Reviewed

- `research/01-plan.md` — Initial research plan & questions
- `research/02-sources.md` — Catalog of 13 sources (academic, vendors, blogs, preprints, community tools)
- `research/03-evidence.md` — 14 evidence statements compiled from sources
- `research/04-contradictions.md` — Analysis of 3 potential tension points and areas of agreement
- `research/05-report.md` — Comprehensive research synthesis with 11 findings
- `research/06-open-questions.md` — Log of open questions, weak evidence, and research gaps
- `research-revision/01-revision-plan.md` — Prior revision plan addressing 8 audit issues
- `research-revision/02-changes-made.md` — Detailed log of 6 revision actions
- `research-revision/03-revision-result.md` — Summary of revision results and readiness state

## Claims To Verify

1. Code coverage measures execution, not fault detection capability (100% coverage false confidence).
2. Mutation Score formula = (Killed Mutants / Total Mutants) × 100%.
3. Mutation operators systematically simulate typical human errors (statement deletion, operator replacement, etc.).
4. Theoretical foundation relies on the Competent Programmer Hypothesis and Coupling Effect (DeMillo et al. 1978; Offutt 1992).
5. The RIP model (Reach, Infect, Propagate) defines the conditions required to kill a mutant.
6. Equivalent mutant detection is mathematically undecidable and poses a major practical barrier.
7. Meta ACH uses LLMs to generate targeted mutants/tests, achieving 73% engineer acceptance, 36% privacy relevance, and 0.95/0.96 precision/recall with static analysis preprocessing.
8. Meta/Harman's formulation of the "Five Barriers" to mutation testing adoption.
9. Go mutation testing tools (`go-mutesting`, `gremlins`) exist as community projects but lack feature parity with JVM (PIT) or JS/TS (Stryker) tools.
10. Subsumed mutants exist at the same source location as another mutant and do not contribute to coverage metrics.
11. No universal industry-standard threshold exists for mutation scores.

## Code To Execute

Per PIPELINE OVERRIDE, source code and test execution are excluded from this research-only audit phase. Verification is restricted to URL validation, source corroboration, claim-to-evidence mapping, and documentation consistency.

## Primary Risks

- Unverified academic citations (DeMillo et al. 1978, Jia & Harman 2009) cited via secondary references.
- Single-source industrial claims regarding Meta ACH performance (arXiv preprint abstract vs. Meta blog post).
- Citation of pre-publication draft material (Martin Fowler bliki entry).
- Potential overgeneralization of implementation-specific trial data to universal technical rules.

## Audit Strategy

1. Audit source metadata, reachable URLs, and source tiers (Tier 1 primary vs. Tier 2 community/secondary).
2. Audit each claim in `03-evidence.md` and `05-report.md` for accuracy, attribution, classification, and severity.
3. Verify consistency between research findings, evidence logs, open questions, and revision documents.
4. Assess research completeness and explicitly record remaining gaps and non-blocking caveats.
5. Formulate an evidence-based final audit verdict (`07-verdict.md`).
