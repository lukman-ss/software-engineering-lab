# Claim Audit Report

## Claim 1

Claim: Mutation testing was first proposed by Richard Lipton in 1971 and published by DeMillo, Lipton, and Sayward in 1978.

Location: research/03-evidence.md — Evidence 1; research/05-report.md — Executive Summary

Evidence Provided: Wikipedia (citing DeMillo et al. 1978); Martin Fowler bliki (DRAFT — 404 at audit time)

Source: Wikipedia (Source 2); Martin Fowler (Source 3 — UNREACHABLE)

Source Actually Supports Claim: YES (via Wikipedia; NOT VERIFIED via Fowler bliki)

Classification: FACT

Severity: LOW

Notes: Wikipedia directly confirmed: "Mutation testing was originally proposed by Richard Lipton as a student in 1971, and first developed and published by DeMillo, Lipton and Sayward." The primary paper (Source 1) is not directly accessible but the bibliographic reference is verified in Wikipedia. Martin Fowler's corroboration source is now unreachable (404). Claim stands on Wikipedia alone; single encyclopedia source is not ideal for a historical fact claim, but is acceptable with correct disclosure.

---

## Claim 2

Claim: The two core hypotheses of mutation testing are the "competent programmer hypothesis" and the "coupling effect hypothesis."

Location: research/03-evidence.md — Evidence 2; research/05-report.md — Finding 4

Evidence Provided: Wikipedia (citing Ammann & Offutt 2008; Offutt 1992; Acree et al. 1979)

Source: Wikipedia (Source 2)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Verified directly from Wikipedia fetch: "Mutation testing is based on two hypotheses. The first is the competent programmer hypothesis... The second hypothesis is called the coupling effect." Language matches research evidence exactly. The underlying academic citations (Offutt 1992, Acree et al. 1979) not directly verified, but Wikipedia is a reliable intermediary for this well-established theoretical foundation.

---

## Claim 3

Claim: Mutation Score = (Killed Mutants / Total Mutants) × 100%. A high mutation score indicates robust test quality.

Location: research/03-evidence.md — Evidence 3; research/05-report.md — Finding 2

Evidence Provided: Wikipedia; Martin Fowler (DRAFT — unreachable); PIT FAQ

Source: Wikipedia (Source 2); PIT FAQ (Source 5)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Wikipedia confirmed: "The value of a test suite is measured by the percentage of mutants that it kills." PIT confirmed: "The quality of your tests can be gauged from the percentage of mutations killed." Formula stated in research is consistent with both sources. Martin Fowler source unreachable but claim supported by two independent verified sources.

---

## Claim 4

Claim: Code coverage can be 100% while test quality is poor; mutation testing exposes this discrepancy.

Location: research/03-evidence.md — Evidence 4; research/05-report.md — Finding 1

Evidence Provided: PIT homepage; Stryker docs (sandwich/paste analogy); Stryker RoboCoasters example

Source: PIT (Source 4); Stryker (Source 6); Stryker Example (Source 7)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: PIT confirmed verbatim: "Traditional test coverage... measures only which code is executed... does not check that your tests are actually able to detect faults." Stryker confirmed verbatim: sandwich/paste analogy. RoboCoasters example confirmed: "How code coverage of 100% could mean only 60% is tested." Strongest and most thoroughly evidenced claim in the research.

---

## Claim 5

Claim: Mutation operators include statement deletion, arithmetic operator replacement, relational operator replacement, Boolean replacement, and statement insertion.

Location: research/03-evidence.md — Evidence 5; research/05-report.md — Finding 3

Evidence Provided: Wikipedia (citing Hamimoune & Falah, 2016); PIT FAQ mutator groups; Stryker output examples

Source: Wikipedia (Source 2); PIT (Source 5); Stryker (Source 6)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Wikipedia confirmed verbatim list: statement deletion, statement duplication/insertion, replacement of Boolean subexpressions with true/false, replacement of arithmetic operations, replacement of Boolean relations (>, >=, ==, <=). PIT confirmed DEFAULTS/STRONGER/ALL groups. Stryker confirmed BinaryOperator and RemoveConditionals mutators. Claim well-supported.

---

## Claim 6

Claim: Equivalent mutants are a major obstacle in mutation testing because determining equivalence is mathematically undecidable.

Location: research/03-evidence.md — Evidence 6; research/05-report.md — Finding 6

Evidence Provided: Wikipedia (citing [18], [19]); Meta Engineering Blog

Source: Wikipedia (Source 2); Meta Engineering Blog (Source 8)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Wikipedia confirmed: "Equivalent mutants detection is one of the biggest obstacles to practical usage of mutation testing." Meta Engineering Blog confirmed: "Determining whether a mutant is equivalent or not is known to be mathematically undecidable." Both sources corroborate the claim independently. Research correctly identifies this as a theoretical undecidability (not a practical one).

---

## Claim 7

Claim: Strong mutation testing requires all three RIP conditions: Reach, Infect, and Propagate.

Location: research/03-evidence.md — Evidence 7; research/05-report.md — Finding 5

Evidence Provided: Wikipedia (citing Offutt & Untch, 2000 — "Mutation 2000")

Source: Wikipedia (Source 2)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Wikipedia confirmed verbatim: three conditions (Reach, Infect, Propagate), collectively called the RIP model, citing [9] (Offutt & Untch, 2000). Confirmed that weak mutation requires only Reach + Infect; strong mutation requires all three. Research evidence matches Wikipedia content exactly. The underlying paper (cs.gmu.edu/~offutt/rsrch/papers/mut00.pdf) is not directly accessed but Wikipedia's reference is specific and attributable.

---

## Claim 8

Claim: All major mutation testing tools follow the same pattern: generate mutants → run tests → classify as killed/survived.

Location: research/03-evidence.md — Evidence 8; research/05-report.md — Executive Summary

Evidence Provided: Wikipedia; PIT; Stryker

Source: Wikipedia (Source 2); PIT (Source 4); Stryker (Source 6)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Three independent sources confirm the same workflow pattern. Use of "all major" is technically an overgeneralization, but is effectively supported by the three most significant tools (PIT, Stryker, go-mutesting). The claim in the report is attributed correctly to "all major" tools, not "all possible" tools. Minor quibble only — not a material issue.

---

## Claim 9

Claim: Mutation testing faces five major barriers to industrial adoption: scalability, unrealistic mutants, equivalent mutants, computational cost, and overstretching testing efforts.

Location: research/03-evidence.md — Evidence 9; research/05-report.md — Finding 7

Evidence Provided: Meta Engineering Blog (explicitly lists five barriers); PIT FAQ (corroborates computational cost)

Source: Meta Engineering Blog (Source 8)

Source Actually Supports Claim: YES

Classification: INTERPRETATION

Severity: LOW

Notes: Meta Engineering Blog verified as listing exactly five barriers, with corresponding headers in the article. This framing is explicitly Meta/Harman's taxonomy, not a universal classification from the broader literature. Research Agent correctly noted: "The 'five barriers' framing is Harman/Meta's specific articulation, not a universal taxonomy." This is an implementation-specific classification being labeled as INTERPRETATION, which is appropriate and correctly disclosed.

---

## Claim 10

Claim: Mutation testing is classified as white-box testing and its purpose is to help develop effective regression tests.

Location: research/03-evidence.md — Evidence 10; research/05-report.md — Executive Summary

Evidence Provided: Wikipedia (citing Ostrand 2002; Misra 2003)

Source: Wikipedia (Source 2)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Wikipedia confirmed verbatim: "Mutation testing is a form of white-box testing. Its purpose is to help the tester develop effective regression tests..." Matches research claim precisely.

---

## Claim 11

Claim: There are three types of mutation testing: statement mutation, value mutation, and decision mutation.

Location: research/03-evidence.md — Evidence 11; research/05-report.md — Finding 3

Evidence Provided: Wikipedia (section "Types of mutation testing")

Source: Wikipedia (Source 2)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Wikipedia confirmed: three sections — Statement mutation, Value mutation, Decision mutation — with code examples for each. Research rated this MEDIUM confidence and noted "pedagogically useful; industrial tools use different categorizations." This qualifier is appropriate. The three-type taxonomy directly from Wikipedia content matches the claim.

---

## Claim 12

Claim: Meta ACH achieves 73% engineer acceptance, 36% privacy relevance, and 0.95/0.96 precision/recall with preprocessing.

Location: research/03-evidence.md — Evidence 12; research/05-report.md — Finding 7

Evidence Provided: arXiv abstract (Source 12); Meta Engineering Blog (Source 8)

Source: arXiv:2501.12862 (abstract); Meta Engineering Blog

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: arXiv abstract confirmed verbatim: "9,095 mutants and 571 privacy-hardening test cases... precision of 0.79 and a recall of 0.47 (rising to 0.95 and 0.96 with simple pre-processing)... engineers accepted 73% of its tests, judging 36% to privacy relevant." All statistics verified from both sources; no discrepancy detected. Research Agent correctly rated this MEDIUM confidence due to industry source nature, while documenting corroboration from arXiv.

---

## Claim 13

Claim: Go has two community-driven mutation testing tools (go-mutesting, gremlins) neither matching maturity of PIT or Stryker.

Location: research/03-evidence.md — Evidence 13; research/05-report.md — Finding 8

Evidence Provided: go-mutesting GitHub (Source 10); gremlins GitHub (Source 11)

Source: Both GitHub repositories

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Both repositories confirmed reachable and active. go-mutesting: 677 stars, community tool. gremlins: 433 stars, still in 0.x.x pre-1.0 releases (explicitly stated in README: "Gremlins is still in its 0.x.x release"). Gremlins README acknowledges "doesn't work very well on very big Go modules." Research Agent's maturity assessment is corroborated. PIT supports Java/Kotlin only (verified from FAQ), Stryker supports JS/TS/C#/Scala (no Go mentioned).

---

## Claim 14

Claim: No universally accepted minimum mutation score threshold exists; proposed values (80%, 85%, 90%) vary by project.

Location: research/05-report.md — Finding 11

Evidence Provided: PIT FAQ (no threshold stated); Stryker docs (no threshold stated)

Source: PIT (Source 5); Stryker (Source 6)

Source Actually Supports Claim: PARTIAL

Classification: INTERPRETATION

Severity: MEDIUM

Notes: Research Agent correctly identifies absence of threshold in PIT FAQ and Stryker docs. However, the claim that thresholds of "80%, 85%, 90% variously proposed elsewhere" is stated but NOT VERIFIED — no source for these specific numbers is provided. The claim itself is framed as an absence claim (no standard exists), which is supported by the sources. The parenthetical specific numbers (80%, 85%, 90%) are flagged in research/06-open-questions.md as "NOT VERIFIED." This parenthetical should not appear in the report without citation, but the core finding (no standard exists) is valid. Classified MEDIUM because the specific percentage values are unsourced.
