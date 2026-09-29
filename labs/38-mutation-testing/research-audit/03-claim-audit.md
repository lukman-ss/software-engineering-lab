# Claim Audit

Target Lab: `/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/38-mutation-testing`
Audit Date: 2026-09-29

---

## Claim 1
Claim: Mutation testing was first proposed by Richard Lipton in 1971 and first published by DeMillo, Lipton, and Sayward in 1978.
Location: `research/03-evidence.md` Evidence 1; `research/05-report.md` Executive Summary
Evidence Provided: Wikipedia citation of DeMillo et al. 1978; not directly accessed.
Source: Wikipedia (https://en.wikipedia.org/wiki/Mutation_testing)
Source Actually Supports Claim: YES — Wikipedia text confirmed: "Mutation testing was originally proposed by Richard Lipton as a student in 1971, and first developed and published by DeMillo, Lipton and Sayward."
Classification: FACT
Severity: LOW
Notes: Secondary attribution properly disclosed. The direct paper (IEEE Computer 11(4):34-41, April 1978) is bibliographically well-established and confirmed referenced in Wikipedia's footnotes.

---

## Claim 2
Claim: Mutation testing rests on two core hypotheses: the Competent Programmer Hypothesis and the Coupling Effect Hypothesis.
Location: `research/03-evidence.md` Evidence 2; `research/05-report.md` Finding 4
Evidence Provided: Wikipedia quote confirmed; sourced to Ammann & Offutt (2008), Offutt (1992), Acree et al. (1979).
Source: Wikipedia
Source Actually Supports Claim: YES — Confirmed verbatim from Wikipedia: "The first is the competent programmer hypothesis... The second hypothesis is called the coupling effect."
Classification: FACT
Severity: LOW
Notes: Both hypotheses are universally recognized in the mutation testing literature. Not independently verified from original papers, but consistent with established academic consensus.

---

## Claim 3
Claim: Mutation Score = (Killed Mutants / Total Mutants) × 100%. A high score indicates robust tests; a low score reveals weak assertions regardless of code coverage.
Location: `research/03-evidence.md` Evidence 3; `research/05-report.md` Finding 2; `README.md`
Evidence Provided: Wikipedia, PIT homepage. Martin Fowler URL cited but marked as "pre-publication draft" and now 404.
Source: Wikipedia + PIT (https://pitest.org/)
Source Actually Supports Claim: YES — Wikipedia: "The value of a test suite is measured by the percentage of mutants that it kills." PIT: "The quality of your tests can be gauged from the percentage of mutations killed."
Classification: FACT
Severity: LOW
Notes: Formula is correct and universally accepted. The Martin Fowler URL is 404 but claim is fully supported by the other two verified sources. Residual citation in report is properly labeled as unreachable.

---

## Claim 4
Claim: 100% line coverage can coexist with poor test quality; mutation testing exposes this gap.
Location: `research/03-evidence.md` Evidence 4; `research/05-report.md` Finding 1; `README.md` Overview
Evidence Provided: PIT direct quote; Stryker sandwich/paste analogy.
Source: PIT (https://pitest.org/), Stryker (https://stryker-mutator.io/docs/)
Source Actually Supports Claim: YES — Both sources confirmed verbatim as cited. Stryker RoboCoasters example demonstrates exactly 100% coverage + ~60% mutation score.
Classification: FACT
Severity: LOW
Notes: Central thesis of the lab; strongly supported by two independent primary sources.

---

## Claim 5
Claim: Standard mutation operators include statement deletion, arithmetic replacement (+↔*, -↔/), relational replacement (>↔>=, ==↔<=), Boolean replacement (&&↔||), and statement insertion.
Location: `research/03-evidence.md` Evidence 5; `research/05-report.md` Finding 3
Evidence Provided: Wikipedia (citing Hamimoune & Falah, 2016), PIT FAQ (DEFAULTS/STRONGER/ALL groups), Stryker (BinaryOperator, RemoveConditionals).
Source: Wikipedia, PIT FAQ, Stryker docs
Source Actually Supports Claim: YES — Wikipedia operator list confirmed. PIT FAQ mutator groups confirmed. Stryker mutator names confirmed.
Classification: FACT
Severity: LOW
Notes: Confirmed from 3 independent reachable sources. The lab's implemented operators (relational, boolean, arithmetic, boundary shift) map directly to these established operator categories.

---

## Claim 6
Claim: Equivalent mutants are semantically identical to original code, impossible to kill regardless of test quality, and determining equivalence is mathematically undecidable.
Location: `research/03-evidence.md` Evidence 6; `research/05-report.md` Finding 6
Evidence Provided: Wikipedia, Meta Engineering Blog (Harman).
Source: Wikipedia, Meta Engineering Blog (https://engineering.fb.com/2025/09/30/security/llms-are-the-key-to-mutation-testing-and-better-compliance/)
Source Actually Supports Claim: YES — Wikipedia: "Equivalent mutants detection is one of the biggest obstacles to practical usage of mutation testing." Meta blog: "Determining whether a mutant is equivalent or not is known to be mathematically undecidable."
Classification: FACT
Severity: LOW
Notes: Well-established theoretical result. Both sources confirmed. PIT FAQ further corroborates via its avoidance of static initializer mutations.

---

## Claim 7
Claim: For a test to kill a mutant (strong mutation), the RIP model requires: Reach (test executes mutated statement), Infect (input causes different program state), Propagate (incorrect state reaches test assertion).
Location: `research/03-evidence.md` Evidence 7; `research/05-report.md` Finding 5
Evidence Provided: Wikipedia (citing Offutt & Untch, 2000 — Mutation 2000).
Source: Wikipedia
Source Actually Supports Claim: YES — Confirmed verbatim from Wikipedia with the 3 numbered conditions.
Classification: FACT
Severity: LOW
Notes: Wikipedia also confirmed: "Weak mutation testing (or weak mutation coverage) requires that only the first and second conditions are satisfied." This distinction is correctly noted in research (weak mutation = Reach+Infect only).

---

## Claim 8
Claim: All major evaluation-focused mutation testing tools follow: generate mutants → run tests → classify (killed/survived). Generative tools (Meta ACH) extend this by using unkilled mutants to prompt LLM test synthesis.
Location: `research/03-evidence.md` Evidence 8; `research/05-report.md` Executive Summary
Evidence Provided: Wikipedia, PIT, Stryker.
Source: Wikipedia, PIT (https://pitest.org/), Stryker (https://stryker-mutator.io/docs/)
Source Actually Supports Claim: YES — Standard workflow confirmed from all 3 sources. The qualification of "evaluation-focused" tools vs generative tools (ACH) was added during revision and is accurate.
Classification: FACT
Severity: LOW
Notes: Revision correctly distinguished evaluation-focused vs generative tools. No overgeneralization remains.

---

## Claim 9
Claim: Large-scale adoption is limited by five barriers (scalability, unrealistic mutants, equivalent mutants, computational cost, overstretching), which LLMs can help overcome via mutation-guided test generation.
Location: `research/03-evidence.md` Evidence 9; `research/05-report.md` Finding 7
Evidence Provided: Meta Engineering Blog, PIT FAQ.
Source: Meta Engineering Blog (primary), PIT FAQ (corroborating)
Source Actually Supports Claim: YES — Five barriers confirmed by name in Meta blog with full explanations. The framing is Harman/Meta's characterization; research correctly notes "The 'five barriers' framing is Harman/Meta's specific articulation, not a universal taxonomy."
Classification: INTERPRETATION (the five-barriers framing is industry-specific, not universally defined)
Severity: LOW
Notes: The qualification is present in the research. The claim is well-attributed and clearly sourced to one organization's perspective.

---

## Claim 10
Claim: Mutation testing is white-box testing; its purpose is to develop effective regression tests.
Location: `research/03-evidence.md` Evidence 10; `research/05-report.md` Finding introduction
Evidence Provided: Wikipedia (citing Ostrand 2002; Misra 2003).
Source: Wikipedia
Source Actually Supports Claim: YES — Wikipedia confirmed: "Mutation testing is a form of white-box testing. Its purpose is to help the tester develop effective regression tests."
Classification: FACT
Severity: LOW
Notes: Straightforward, well-established classification. Confirmed.

---

## Claim 11
Claim: There are three types of mutation testing: statement mutation, value mutation, and decision mutation.
Location: `research/03-evidence.md` Evidence 11; `research/05-report.md` Finding 3 Notes
Evidence Provided: Wikipedia.
Source: Wikipedia
Source Actually Supports Claim: YES — Wikipedia section "Types of mutation testing" confirmed with all three types (statement, value, decision) with code examples.
Classification: FACT
Severity: LOW
Notes: Research correctly notes this taxonomy is "pedagogically useful; industrial tools use different categorizations."

---

## Claim 12
Claim: Meta ACH was applied to 10,795 Android Kotlin classes, generated 9,095 mutants and 571 privacy-hardening tests; LLM equivalence detector achieves 0.79/0.47 precision/recall raw, rising to 0.95/0.96 with preprocessing; engineers accepted 73% of tests, 36% privacy-relevant.
Location: `research/03-evidence.md` Evidence 12; `research/05-report.md` Finding 7
Evidence Provided: arXiv:2501.12862 abstract, Meta Engineering Blog.
Source: arXiv (https://arxiv.org/abs/2501.12862), Meta Blog (https://engineering.fb.com/2025/09/30/security/llms-are-the-key-to-mutation-testing-and-better-compliance/)
Source Actually Supports Claim: YES — All statistics confirmed verbatim from arXiv abstract. All statistics also confirmed from Meta blog.
Classification: FACT (empirical claim from single-source study, industry-reported)
Severity: LOW
Notes: Confidence is correctly marked MEDIUM in research due to industry-origin and lack of independent peer-reviewed replication. Statistics are confirmed from two sources (blog + preprint) but both originate from Meta. The preprint is submitted to FSE 2025 Industry Track, not yet independently peer-reviewed at time of research. This is a real limitation that is correctly documented.

---

## Claim 13
Claim: Go mutation testing tools (`go-mutesting`, `gremlins`) exist but lack production maturity, CI parity, or feature parity with PIT or Stryker.
Location: `research/03-evidence.md` Evidence 13; `research/05-report.md` Finding 8
Evidence Provided: go-mutesting GitHub README, gremlins GitHub README, PIT FAQ, Stryker docs.
Source: GitHub (https://github.com/zimmski/go-mutesting, https://github.com/go-gremlins/gremlins), PIT FAQ, Stryker docs
Source Actually Supports Claim: YES — Both GitHub repositories confirmed as reachable and active but community-level tools. Gremlins README explicitly states it's in 0.x.x (no backward compatibility guarantee) and "doesn't work very well on very big Go modules." PIT FAQ confirms Java/Kotlin only. Stryker supports JS/TS, C#, Scala — no Go.
Classification: FACT
Severity: LOW
Notes: One minor omission: research does not call out gremlins' explicit statement about being designed for "smallish modules" only. This is an additional nuance that would strengthen the maturity-gap claim. However, the overall conclusion is correct and supported.

---

## Claim 14
Claim: Subsumed mutants exist at the same source location as another mutant, are "subsumed" by it, and do not contribute to coverage metrics.
Location: `research/03-evidence.md` Evidence 14; `research/05-report.md` Finding 10
Evidence Provided: Wikipedia subsumed mutants section.
Source: Wikipedia
Source Actually Supports Claim: YES — Wikipedia text confirmed verbatim as quoted in research.
Classification: FACT
Severity: LOW
Notes: Advanced refinement concept. Correctly labeled as secondary concern in research.

---

## Claim 15
Claim: No universally accepted minimum mutation score threshold exists; targets vary by project risk profile, codebase maturity, and team disincentive risks.
Location: `research/05-report.md` Finding 11
Evidence Provided: PIT FAQ (no threshold stated), Stryker docs (no threshold stated). Previous revision removed unsourced numbers (80%, 85%, 90%) from an earlier draft.
Source: PIT FAQ (https://pitest.org/faq/), Stryker docs (https://stryker-mutator.io/docs/)
Source Actually Supports Claim: YES — Confirmed by absence of threshold guidance in both verified sources. Negative evidence (sources not prescribing a threshold) is a valid form of support for this claim.
Classification: FACT
Severity: LOW
Notes: Revision correctly removed unsourced specific percentages. The current finding is an honest reporting of what sources do and do not say. No overstated conclusions.

---

## Universal/Absolute Language Check

Searched all research files for "always", "never", "must", "all tools", "every" patterns:
- `research/04-contradictions.md`: "All authoritative sources agree on..." — describes areas of agreement; not an overclaim.
- Evidence 8 (pre-revision): Was "All major mutation testing tools follow..." — corrected to "All major evaluation-focused mutation testing tools follow..." during revision. Current version is appropriately qualified.
- No remaining universal claims found that lack source support.

## Numeric Threshold Check

- Previous draft contained unsourced "(80%, 85%, 90%)" in Finding 11.
- Revision 1 removed these numbers. Current `research/05-report.md` Finding 11 contains no specific numeric thresholds.
- `README.md` contains no threshold claims.
- PASS.
