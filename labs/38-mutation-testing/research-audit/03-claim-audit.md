# Claim Audit

## Claim 1

Claim:
Mutation testing was originally proposed by Richard Lipton in 1971 and first developed and published by DeMillo, Lipton, and Sayward in 1978.

Location:
`research/03-evidence.md:5-11`, `research/05-report.md:9`

Evidence Provided:
Wikipedia citation of DeMillo, Lipton, Sayward (1978) and Mutation 2000.

Source:
Source 2 (Wikipedia), referencing Source 1 (DeMillo et al. 1978).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately documented and explicitly flagged as secondary attribution from Wikipedia.

---

## Claim 2

Claim:
The two core theoretical hypotheses of mutation testing are the "competent programmer hypothesis" and the "coupling effect hypothesis".

Location:
`research/03-evidence.md:15-21`, `research/05-report.md:51-57`

Evidence Provided:
Wikipedia section citing DeMillo et al. (1978), Offutt (1992), and Acree et al. (1979).

Source:
Source 2 (Wikipedia).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Theoretical definitions match standard software engineering literature.

---

## Claim 3

Claim:
The mutation score formula is: Mutation Score = (Killed Mutants / Total Mutants) × 100%, and measures test suite quality.

Location:
`research/03-evidence.md:25-32`, `research/05-report.md:27-33`

Evidence Provided:
PIT documentation, Wikipedia, Martin Fowler bliki.

Source:
Source 2 (Wikipedia), Source 3 (Martin Fowler), Source 4 (PIT).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Standard definition across industry and academia.

---

## Claim 4

Claim:
Traditional code coverage can reach 100% while test assertions are weak or missing; mutation testing exposes this false confidence.

Location:
`research/03-evidence.md:35-40`, `research/05-report.md:15-21`

Evidence Provided:
PIT homepage statement and Stryker sandwich paste analogy.

Source:
Source 4 (PIT), Source 6 & 7 (Stryker).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Core thesis supported directly by verified tool documentation.

---

## Claim 5

Claim:
Mutation operators systematically modify code (statement deletion, arithmetic replacement, relational replacement, boolean replacement, statement insertion).

Location:
`research/03-evidence.md:44-50`, `research/05-report.md:39-45`

Evidence Provided:
Wikipedia citing Hamimoune & Falah (2016), PIT mutator groups, Stryker output.

Source:
Source 2 (Wikipedia), Source 5 (PIT FAQ), Source 6 (Stryker).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Traditional and class/method level operator taxonomies are verified.

---

## Claim 6

Claim:
Equivalent mutants are semantically identical to original code, making them impossible to kill; determining equivalence is mathematically undecidable.

Location:
`research/03-evidence.md:54-59`, `research/05-report.md:75-81`

Evidence Provided:
Wikipedia citing Frankl et al. (1997) and Meta ACH paper citing theoretical limits.

Source:
Source 2 (Wikipedia), Source 8 (Meta Engineering Blog).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Undecidability stems from Rice's Theorem/Halting Problem in program semantics.

---

## Claim 7

Claim:
Strong mutation testing requires all three RIP conditions: Reach, Infect, and Propagate; weak mutation requires only Reach and Infect.

Location:
`research/03-evidence.md:63-69`, `research/05-report.md:63-69`

Evidence Provided:
Wikipedia citing Offutt & Untch (2000).

Source:
Source 2 (Wikipedia).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
RIP model accurately documented.

---

## Claim 8

Claim:
Standard mutation testing workflow across tools: generate mutants -> run test suite against mutants -> classify as killed or survived.

Location:
`research/03-evidence.md:73-79`, `research/05-report.md:9`

Evidence Provided:
PIT and Stryker documentation.

Source:
Source 4 (PIT), Source 6 (Stryker).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately describes mechanics across bytecode and AST tools.

---

## Claim 9

Claim:
Mutation testing historically faced five industrial adoption barriers: scalability, unrealistic mutants, equivalent mutants, computational cost, and overstretched testing efforts.

Location:
`research/03-evidence.md:83-89`, `research/05-report.md:87`

Evidence Provided:
Mark Harman's keynote and blog post on Meta ACH.

Source:
Source 8 (Meta Engineering Blog).

Source Actually Supports Claim:
YES

Classification:
INTERPRETATION

Severity:
MEDIUM

Notes:
Categorization is Harman/Meta's formulation, though widely acknowledged as general pain points.

---

## Claim 10

Claim:
Mutation testing is a white-box testing technique whose purpose is to help develop effective regression tests.

Location:
`research/03-evidence.md:93-99`, `research/05-report.md:9`

Evidence Provided:
Wikipedia citing Ostrand (2002) and Misra (2003).

Source:
Source 2 (Wikipedia).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Standard categorization.

---

## Claim 11

Claim:
Meta ACH uses LLMs to generate targeted mutants and tests, achieving 73% engineer acceptance, 36% privacy relevance, and 0.95/0.96 precision/recall on equivalence detection with preprocessing.

Location:
`research/03-evidence.md:113-119`, `research/05-report.md:87-93`

Evidence Provided:
Meta Engineering blog post by Mark Harman (Sep 30, 2025).

Source:
Source 8 (Meta Engineering Blog).

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
Numbers are specific to Meta's internal trial (Oct-Dec 2024, Kotlin, privacy domain) and should not be generalized as universal LLM mutation performance. Research report appropriately flags confidence as MEDIUM.

---

## Claim 12

Claim:
Practical tooling is strongest for JVM (PIT) and JS/TS (Stryker), while Go has limited mature mutation testing tooling.

Location:
`research/05-report.md:99-105`, `research/06-open-questions.md:5`

Evidence Provided:
PIT FAQ, Stryker documentation, ecosystem survey.

Source:
Source 5 (PIT FAQ), Source 6 (Stryker docs).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurate reflection of the state of tooling in 2026.
