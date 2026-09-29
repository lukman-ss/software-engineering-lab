# Research Audit Plan

Target Lab: `/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/38-mutation-testing`
Audit Date: 2026-09-29
Auditor: Independent Technical Research Auditor

## 1. Inventory of Files Reviewed

- `README.md` (Lab overview, mutation testing formula, directory structure, implemented mutation operators)
- `research/01-plan.md` (Research topic, objectives, 10 research questions, search strategy, primary sources, risks)
- `research/02-sources.md` (13 sources documented with tiering, access dates, and status)
- `research/03-evidence.md` (14 claims linked to evidence items and sources)
- `research/04-contradictions.md` (Analysis of 3 potential tension points and areas of agreement)
- `research/05-report.md` (Comprehensive research report with 11 findings, executive summary, areas of agreement/disagreement, limitations, conclusion)
- `research/06-open-questions.md` (Unanswered questions, weak evidence items, deeper research topics, next directions)
- `research/runs/2026-09-29-mutation-testing/01-plan.md` (Run plan document)
- `research-revision/01-revision-plan.md` (Post-audit revision plan)
- `research-revision/02-changes-made.md` (Detailed records of revisions applied)
- `research-revision/03-revision-result.md` (Validation and readiness status)

## 2. Inventory of Claims to Verify

1. **Historical Origins**: First proposed by Richard Lipton in 1971; published by DeMillo, Lipton, and Sayward in 1978.
2. **Theoretical Foundations**: Two core hypotheses: Competent Programmer Hypothesis and Coupling Effect Hypothesis.
3. **Mutation Score Definition**: Formula is $(\text{Killed Mutants} / \text{Total Mutants}) \times 100\%$.
4. **Line Coverage Discrepancy**: 100% line coverage can fail to catch bugs when assertions are weak; mutation testing identifies this gap.
5. **Mutation Operators**: Common operators include statement deletion, arithmetic replacement, relational replacement, boolean replacement, and value modification.
6. **Equivalent Mutant Problem**: Determining semantic equivalence is mathematically undecidable; represents a primary practical obstacle.
7. **RIP Model**: Strong mutation requires Reach, Infect, and Propagate conditions; weak mutation requires Reach and Infect.
8. **Evaluation vs Generative Tool Workflow**: Evaluation tools follow generate $\to$ run $\to$ classify; generative systems (Meta ACH) use unkilled mutants to prompt test synthesis.
9. **Industrial Scaling Barriers**: Five barriers identified by Harman/Meta (scalability, unrealistic mutants, equivalent mutants, computational cost, overstretching).
10. **Testing Classification**: Mutation testing is white-box testing designed to improve regression test suites.
11. **Taxonomy of Modifications**: Categorized into statement mutation, value mutation, and decision mutation.
12. **Meta ACH Empirical Data**: Applied across 10,795 Android Kotlin classes, 9,095 mutants, 571 tests; 73% engineer acceptance, 36% privacy relevance, 0.95/0.96 precision/recall with preprocessing.
13. **Go Tooling Ecosystem**: Community tools exist (`go-mutesting`, `gremlins`), but lack the production maturity/CI parity of JVM (PIT) and JS/TS (Stryker).
14. **Subsumed Mutants**: Mutants at the same source location producing identical failure dynamics that do not contribute to unique coverage metrics.
15. **Score Thresholds**: No universal industry-standard threshold exists; determined empirically per risk profile.

## 3. Inventory of Source URLs Found

1. `https://en.wikipedia.org/wiki/Mutation_testing`
2. `https://martinfowler.com/bliki/MutationTesting.html`
3. `https://pitest.org/`
4. `https://pitest.org/faq/`
5. `https://stryker-mutator.io/docs/`
6. `https://stryker-mutator.io/docs/General/example/`
7. `https://engineering.fb.com/2025/09/30/security/llms-are-the-key-to-mutation-testing-and-better-compliance/`
8. `https://arxiv.org/abs/2501.12862`
9. `https://github.com/zimmski/go-mutesting`
10. `https://github.com/go-gremlins/gremlins`
11. `https://gremlins.dev`

## 4. Primary Risks Identified

- **Unreachable URL**: Martin Fowler's bliki URL returned HTTP 404 (draft status). Must verify if research properly disclaims or decouples from it.
- **Secondary Academic Citations**: Foundational papers (DeMillo 1978, Jia & Harman 2009) referenced via Wikipedia rather than direct DOI fetching.
- **Industry Report Bias**: Meta ACH empirical claims originate from a single corporate post and preprint abstract.
- **Overgeneralization of Thresholds**: Risk of attributing rigid mutation score target numbers without normative consensus.
