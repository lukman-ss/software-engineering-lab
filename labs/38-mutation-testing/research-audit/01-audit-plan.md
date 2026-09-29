# Audit Plan

Target Lab: labs/38-mutation-testing
Files Reviewed:
- `labs/38-mutation-testing/research/01-plan.md`
- `labs/38-mutation-testing/research/02-sources.md`
- `labs/38-mutation-testing/research/03-evidence.md`
- `labs/38-mutation-testing/research/04-contradictions.md`
- `labs/38-mutation-testing/research/05-report.md`
- `labs/38-mutation-testing/research/06-open-questions.md`
- `labs/38-mutation-testing/research-revision/01-revision-plan.md`
- `labs/38-mutation-testing/research-revision/02-changes-made.md`
- `labs/38-mutation-testing/research-revision/03-revision-result.md`

Claims To Verify:
1. Mutation testing originally proposed by DeMillo, Lipton, Sayward (1978) & Lipton (1971).
2. Code coverage measures execution; mutation score measures verification capability.
3. Core hypotheses: Competent Programmer Hypothesis and Coupling Effect.
4. RIP model requires Reach, Infect, Propagate for strong mutation testing.
5. Equivalent mutants detection is mathematically undecidable and a major industrial barrier.
6. Subsumed mutants do not contribute to coverage metrics.
7. Meta ACH uses LLMs for targeted mutation testing (73% acceptance rate, 36% privacy-relevant, 0.95/0.96 precision/recall with preprocessing).
8. Tooling maturity varies (PIT for JVM, Stryker for JS/TS, limited maturity for Go like go-mutesting and gremlins).

Code To Execute:
N/A (PIPELINE OVERRIDE: Audit research only. Do not audit implementation/code in this stage).

Primary Risks:
- Secondary-attributed foundational academic citations (DeMillo 1978, Jia & Harman 2009 unread directly).
- Draft status of Martin Fowler bliki page.
- Single-vendor industrial report (Meta ACH engineering blog & arXiv preprint).
- Overgeneralization of implementation-specific data (e.g. Meta ACH Android Kotlin privacy testing).

Audit Strategy:
- Verify reachability and relevance of all cited sources via direct web fetches and bibliographic check.
- Evaluate each major claim against cited evidence and note attribution limitations.
- Analyze internal and source consistency across research artifacts.
- Record all research gaps and render evidence-based verdict.
