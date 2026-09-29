# Source Audit Report

## Source 1

Claimed Title: Hints on test data selection: Help for the practicing programmer
Claimed Publisher: IEEE Computer (DeMillo, Lipton, Sayward)
URL: NOT VERIFIED (original paper not opened; DOI not confirmed)

Reachable: NO (not fetched directly; DOI behind paywall; paper text not accessed)

Source Type: PRIMARY

Relevant: YES

Supports Claimed Topic: PARTIAL

Problems:
- Research Agent correctly flagged this as NOT VERIFIED for direct quotations.
- Evidence 1 states: "Paper text NOT accessed directly. All claims attributed to this paper in this research come from Wikipedia's reference list and summaries."
- Wikipedia's reference list correctly identifies the paper: R. A. DeMillo, R. J. Lipton, F. G. Sayward. IEEE Computer, 11(4):34-41, April 1978. This bibliographic record is real.
- Claims attributed to this paper (CPH, coupling effect, mutation score) are intermediated exclusively through Wikipedia; cannot verify original wording.

Assessment: WARNING — Bibliographic citation is plausible and correctly attributed. Content cannot be verified from the primary source. Research Agent correctly disclosed this limitation. No fabrication detected.

---

## Source 2

Claimed Title: Mutation testing
Claimed Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Mutation_testing

Reachable: YES (verified via direct fetch)

Source Type: SECONDARY

Relevant: YES

Supports Claimed Topic: YES

Problems:
- Wikipedia page confirmed accessible. Last edited 6 September 2026 (confirmed by page footer).
- Content verified: CPH, coupling effect, RIP model, mutation operators, equivalent mutants, subsumed mutants, three types of mutation testing (statement/value/decision) — all confirmed present in article.
- Source is crowdsourced encyclopedia. Research Agent correctly classified as Tier 2, cross-referenced with primary literature.

Assessment: PASS

---

## Source 3

Claimed Title: Mutation Testing (bliki entry)
Claimed Publisher: Martin Fowler
URL: https://martinfowler.com/bliki/MutationTesting.html

Reachable: NO — 404 HTTP error returned during this audit.

Source Type: PRIMARY (claimed Tier 1 authoritative expert)

Relevant: PARTIAL (intended to support mutation testing cost history, tool speedups, LLM relevance)

Supports Claimed Topic: NOT VERIFIED

Problems:
- URL returns 404 at time of audit. Research Agent noted the page is marked "DRAFT" at the time it was accessed.
- Research report uses this source to corroborate mutation score definition and historical cost. If page is not reachable, it cannot be verified.
- Research Agent explicitly labeled it as pre-publication draft and qualified all citations from it.
- The evidence citations from this source are supplementary; core claims are supported by Wikipedia and PIT.
- Possible the draft page was later removed or unpublished by the author.

Assessment: FAIL — URL not reachable at audit time. All claims that relied solely on this source must be reviewed. Research Agent had appropriately caveated this source as DRAFT and used it only for corroboration, not sole support.

---

## Source 4

Claimed Title: PIT Mutation Testing (homepage)
Claimed Publisher: PIT Project
URL: https://pitest.org/

Reachable: YES (verified via direct fetch)

Source Type: PRIMARY (official tool documentation)

Relevant: YES

Supports Claimed Topic: YES

Problems:
- Content confirmed: "Traditional test coverage... measures only which code is executed... does not check that your tests are actually able to detect faults" — verbatim match to research claim.
- Confirmed: killed/survived mutant classification workflow present.
- Confirmed: "mutation testing is the gold standard against which all other types of coverage are measured".
- No discrepancies detected between research claims and actual page content.

Assessment: PASS

---

## Source 5

Claimed Title: PIT FAQ
Claimed Publisher: PIT Project
URL: https://pitest.org/faq/

Reachable: YES (verified via direct fetch)

Source Type: PRIMARY (official tool documentation)

Relevant: YES

Supports Claimed Topic: YES

Problems:
- Confirmed: "Mutation testing is a computationally expensive process" — verbatim match.
- Confirmed: mutator groups DEFAULTS, STRONGER, ALL — present.
- Confirmed: Java and Kotlin languages listed as supported. PIT does NOT list Go support.
- Confirmed: "The quality of your tests can be gauged from the percentage of mutations killed." — verbatim match.
- Confirmed: bytecode mutation in memory, never written to disk.
- Research claim about "not writing mutated code to disk" confirmed from source.

Assessment: PASS

---

## Source 6

Claimed Title: What is mutation testing? (Introduction)
Claimed Publisher: Stryker Mutator
URL: https://stryker-mutator.io/docs/

Reachable: YES (verified via direct fetch)

Source Type: PRIMARY (official tool documentation)

Relevant: YES

Supports Claimed Topic: YES

Problems:
- Confirmed: sandwich/paste analogy present verbatim: "Code coverage would tell you the bread is 80% covered with paste. Mutation testing, on the other hand, would tell you it is actually chocolate paste."
- Confirmed: `user.age >= 18` example with four mutants: `> 18`, `< 18`, `false`, `true`.
- Confirmed: mutator names `BinaryOperator` and `RemoveConditionals` in output example.
- Confirmed: killed/survived classification.

Assessment: PASS

---

## Source 7

Claimed Title: Welcome to the RoboCoasters — An introduction to mutation testing
Claimed Publisher: Stryker Mutator
URL: https://stryker-mutator.io/docs/General/example/

Reachable: YES (verified via direct fetch)

Source Type: PRIMARY (official tool documentation)

Relevant: YES

Supports Claimed Topic: YES

Problems:
- Confirmed: "How code coverage of 100% could mean only 60% is tested" — verbatim subtitle present.
- Confirmed: reproducible example, 100% code coverage + ~60% mutation score demonstrated.
- Research Agent correctly cited this as concrete validation of core lab claim.

Assessment: PASS

---

## Source 8

Claimed Title: LLMs Are the Key to Mutation Testing and Better Compliance
Claimed Publisher: Engineering at Meta (by Mark Harman)
URL: https://engineering.fb.com/2025/09/30/security/llms-are-the-key-to-mutation-testing-and-better-compliance/

Reachable: YES (verified via direct fetch)

Source Type: PRIMARY (industry engineering blog; primary industry source)

Relevant: YES

Supports Claimed Topic: YES

Problems:
- Confirmed: five barriers to mutation testing at scale are explicitly named in the article.
- Confirmed: ACH system uses LLM-based equivalence detector.
- Confirmed: precision 0.79 / recall 0.47, rising to 0.95 / 0.96 with static analysis preprocessing — verbatim match from article body.
- Confirmed: "Determining whether a mutant is equivalent or not is known to be mathematically undecidable" — present.
- Confirmed: 73% acceptance rate, 36% privacy-relevant — present.
- Confirmed: trial scope October–December 2024 on Facebook, Instagram, WhatsApp, wearables platforms.
- No statistics discrepancies detected.

Assessment: PASS

---

## Source 9

Claimed Title: An Analysis and Survey of the Development of Mutation Testing (bibliographic record)
Claimed Publisher: IEEE Transactions on Software Engineering (Jia & Harman, 2009)
URL: doi:10.1109/TSE.2010.62 (not opened)

Reachable: NOT VERIFIED (DOI not fetched; paper text not accessed)

Source Type: PRIMARY (peer-reviewed academic survey)

Relevant: PARTIAL

Supports Claimed Topic: PARTIAL

Problems:
- Research Agent correctly noted "Paper NOT opened directly; cite as secondary-attributed only."
- Not used as sole evidence for any critical claim; referenced as corroborating bibliography from Wikipedia citation list.
- Wikipedia article references this paper with a valid archived PDF link (Semantic Scholar verified by Wikipedia). The Archived PDF link present in Wikipedia footnote is accessible indirectly.

Assessment: WARNING — Not directly verified. Used only for bibliographic attribution support, not as primary claim evidence. Acceptable given correct disclosure by Research Agent.

---

## Source 10

Claimed Title: go-mutesting — Mutation testing for Go source code
Claimed Publisher: zimmski (GitHub repository)
URL: https://github.com/zimmski/go-mutesting

Reachable: YES (verified via direct fetch)

Source Type: SECONDARY (community open-source project)

Relevant: YES

Supports Claimed Topic: YES

Problems:
- Confirmed: "go-mutesting is a framework for performing mutation testing on Go source code."
- Confirmed: branch, expression, statement mutators present (branch/if, branch/else, branch/case; expression/comparison, expression/remove; statement/remove).
- Confirmed: AST-level mutation (source file replacement workflow).
- Note: Research Agent describes this as "exec-based mutation workflow" which is accurate per README.
- Stars: 677, Forks: 59. Community-tier tool with active maintenance.

Assessment: PASS

---

## Source 11

Claimed Title: gremlins — A mutation testing tool for Go
Claimed Publisher: go-gremlins (GitHub organization)
URL: https://github.com/go-gremlins/gremlins

Reachable: YES (verified via direct fetch)

Source Type: SECONDARY (community open-source project)

Relevant: YES

Supports Claimed Topic: YES

Problems:
- Confirmed: "A mutation testing tool for Go" with CLI and documentation at gremlins.dev.
- Confirmed: tool categorizes mutations as: RUNNABLE, NOT COVERED, KILLED, LIVED, TIMED OUT, NOT VIABLE.
- Important note: Gremlins README explicitly states "As of now, Gremlins doesn't work very well on very big Go modules, mainly because a run can take hours to complete." Research Agent's assessment that it lacks production maturity is corroborated.
- Stars: 433, Forks: 48. Still in 0.x.x pre-stability releases.

Assessment: PASS

---

## Source 12

Claimed Title: Mutation-Guided LLM-based Test Generation at Meta (arXiv preprint)
Claimed Publisher: arXiv (cs.SE, cs.AI, cs.LG)
URL: https://arxiv.org/abs/2501.12862

Reachable: YES (verified via direct fetch)

Source Type: PRIMARY (academic preprint; accepted to FSE 2025 Industry Track)

Relevant: YES

Supports Claimed Topic: YES

Problems:
- Confirmed: abstract text verified. Exact match of statistics: 9,095 mutants, 571 privacy-hardening test cases, 10,795 Android Kotlin classes, 7 software platforms.
- Confirmed: precision 0.79 / recall 0.47 (rising to 0.95 / 0.96 with preprocessing).
- Confirmed: 73% engineer acceptance, 36% judged privacy-relevant.
- Confirmed: submitted to FSE 2025 Industry Track (indicated in metadata "Submitted to FSE 2025 Industry Track").
- Authors confirmed: Christopher Foster, Abhishek Gulati, Mark Harman, Inna Harper, Ke Mao, Jillian Ritchey, Hervé Robert, Shubho Sengupta.
- All statistics reported by Research Agent match the arXiv abstract precisely.

Assessment: PASS

---

## Source 13

Claimed Title: Mutation testing (Wikipedia — subsumed mutants section)
Claimed Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Mutation_testing

Reachable: YES (same page as Source 2; verified)

Source Type: SECONDARY

Relevant: YES

Supports Claimed Topic: YES

Problems:
- Confirmed: subsumed mutants section present in the live Wikipedia article.
- Content matches research claim verbatim: "In addition to equivalent mutants, there are subsumed mutants which are mutants that exist in the same source code location as another mutant, and are said to be 'subsumed' by the other mutant."

Assessment: PASS
