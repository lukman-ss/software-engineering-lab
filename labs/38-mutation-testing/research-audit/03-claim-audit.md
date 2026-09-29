# Claim Audit

## Claim 1
Claim: Mutation testing was first proposed by Richard Lipton in 1971 and published by DeMillo, Lipton, and Sayward in 1978.

Location: research/03-evidence.md — Evidence 1; research/05-report.md — Executive Summary

Evidence Provided: Wikipedia citation referencing DeMillo, Lipton, Sayward 1978. Martin Fowler bliki (draft) also mentions historical timeline.

Source: Wikipedia (secondary) attributing DeMillo et al. 1978; original paper NOT opened.

Source Actually Supports Claim: YES
Wikipedia article directly read confirms: "Mutation testing was originally proposed by Richard Lipton as a student in 1971, and first developed and published by DeMillo, Lipton and Sayward."

Classification: FACT

Severity: LOW

Notes: Limitation explicitly disclosed. Bibliographic record confirmed via Wikipedia reading.

---

## Claim 2
Claim: Code coverage measures only what code is executed, not whether tests detect faults. Mutation testing exposes this discrepancy.

Location: research/03-evidence.md — Evidence 4; research/05-report.md — Finding 1

Evidence Provided: PIT homepage ("Traditional test coverage...measures only which code is executed by your tests. It does not check that your tests are actually able to detect faults in the executed code.") and Stryker sandwich analogy.

Source: PIT (https://pitest.org/), Stryker (https://stryker-mutator.io/docs/)

Source Actually Supports Claim: YES
PIT homepage directly verified; quote confirmed.

Classification: FACT

Severity: LOW

Notes: Strong evidence from authoritative primary sources. Both sources directly verified.

---

## Claim 3
Claim: Mutation score formula is (Killed Mutants / Total Mutants) × 100%.

Location: research/03-evidence.md — Evidence 3; research/05-report.md — Finding 2

Evidence Provided: Wikipedia: "The value of a test suite is measured by the percentage of mutants that it kills." PIT: "The quality of your tests can be gauged from the percentage of mutations killed."

Source: Wikipedia, PIT, Martin Fowler (draft)

Source Actually Supports Claim: YES
Wikipedia read directly confirms the formula. PIT homepage verified.

Classification: FACT

Severity: LOW

Notes: Universal convention across all sources. Formula representation not explicitly spelled out as "(Killed/Total)×100%" in sources, but is a direct mathematical interpretation of "percentage of mutants killed." No issues.

---

## Claim 4
Claim: Core theoretical foundations are the Competent Programmer Hypothesis and the Coupling Effect.

Location: research/03-evidence.md — Evidence 2; research/05-report.md — Finding 4

Evidence Provided: Wikipedia directly states both hypotheses, citing DeMillo et al. 1978, Offutt 1992, and Acree et al. 1979.

Source: Wikipedia (secondary) referencing multiple primaries.

Source Actually Supports Claim: YES
Wikipedia article directly read confirms exact text about both hypotheses.

Classification: FACT

Severity: LOW

Notes: Original academic papers not opened; Wikipedia text accessed and verified.

---

## Claim 5
Claim: The RIP model defines when a mutant is killed: Reach, Infect, Propagate.

Location: research/03-evidence.md — Evidence 7; research/05-report.md — Finding 5

Evidence Provided: Wikipedia cites Offutt & Untch, 2000 (Mutation 2000 paper). Exact three conditions listed verbatim from Wikipedia article.

Source: Wikipedia referencing Offutt & Untch (2000)

Source Actually Supports Claim: YES
Wikipedia text directly verified; RIP model is spelled out word-for-word.

Classification: FACT

Severity: LOW

Notes: PIT FAQ supports strong mutation via bytecode mutation. Foundational well-established model.

---

## Claim 6
Claim: Equivalent mutants are mathematically undecidable to detect and are the primary obstacle to practical adoption.

Location: research/03-evidence.md — Evidence 6; research/05-report.md — Finding 6

Evidence Provided: Wikipedia ("effort needed to check if mutants are equivalent or not can be very high, even for small programs"), Meta Engineering Blog ("Determining whether a mutant is equivalent or not is known to be mathematically undecidable").

Source: Wikipedia + Meta Engineering Blog

Source Actually Supports Claim: YES
Wikipedia text verified directly. arXiv abstract confirms undecidability context in ACH description.

Classification: FACT

Severity: LOW

Notes: Undecidability is a standard theoretical result; the Wikipedia article and Meta blog both support this without contradiction.

---

## Claim 7
Claim: Meta ACH achieved 73% engineer acceptance, 36% privacy-relevant tests, and LLM equivalence detector precision 0.95/recall 0.96 with preprocessing.

Location: research/03-evidence.md — Evidence 12; research/05-report.md — Finding 7

Evidence Provided: Meta Engineering Blog (Mark Harman, September 2025). arXiv:2501.12862 abstract directly verified: 9,095 mutants, 571 tests, 10,795 Android Kotlin classes, 7 platforms, same statistics.

Source: Meta Engineering Blog + arXiv preprint abstract

Source Actually Supports Claim: YES — arXiv abstract directly read and all statistics match verbatim.

Classification: FACT — with implementation-specific scoping caveat

Severity: LOW

Notes:
- The "73% acceptance" and "36% privacy-relevant" statistics are from a specific Android Kotlin privacy testing context at Meta, not a general benchmark for all LLM-assisted mutation testing.
- Research files explicitly qualify this as implementation-specific and scope it correctly.
- Confidence is appropriately rated MEDIUM in research/03-evidence.md Evidence 12.

---

## Claim 8
Claim: "Five barriers to industrial adoption" of mutation testing.

Location: research/03-evidence.md — Evidence 9; research/05-report.md — Finding 7

Evidence Provided: Meta Engineering Blog (Mark Harman, September 2025). "Traditional mutation testing generates a very large number of mutants..."

Source: Meta Engineering Blog

Source Actually Supports Claim: PARTIAL

Classification: INTERPRETATION

Severity: MEDIUM

Notes:
- The "five barriers" (scalability, unrealistic mutants, equivalent mutants, computational cost, overstretching) is Meta/Harman's specific framing, not a universally established canonical list.
- Research/05-report.md notes "The 'five barriers' framing is Harman/Meta's specific articulation, not a universal taxonomy." This is appropriate disclosure.
- Research files classify confidence correctly as MEDIUM and scope it to one source.

---

## Claim 9
Claim: Subsumed mutants do not contribute to coverage metrics.

Location: research/03-evidence.md — Evidence 14; research/05-report.md — Finding 10

Evidence Provided: Wikipedia article, Mutation testing — subsumed mutants section.

Source: Wikipedia

Source Actually Supports Claim: YES
Wikipedia text directly verified; exact definition matches the claim verbatim.

Classification: FACT

Severity: LOW

Notes: Well-supported by direct Wikipedia source reading.

---

## Claim 10
Claim: Go has two community-driven mutation testing tools (go-mutesting and gremlins) but neither matches PIT or Stryker in maturity.

Location: research/03-evidence.md — Evidence 13; research/05-report.md — Finding 8

Evidence Provided: GitHub URLs for both tools; PIT FAQ confirming Java/Kotlin support; no Go in PIT/Stryker.

Source: GitHub (go-mutesting, gremlins), PIT FAQ, Stryker docs

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Tools are real and accessible. The maturity comparison is a reasonable inference from ecosystem evidence, not a precisely measured claim.

---

## Claim 11
Claim: No industry-standard mutation score threshold exists.

Location: research/05-report.md — Finding 11

Evidence Provided: Absence of threshold in PIT FAQ and Stryker docs; stated as "no consensus found."

Source: PIT FAQ, Stryker docs (absence of evidence)

Source Actually Supports Claim: PARTIAL

Classification: INTERPRETATION

Severity: LOW

Notes:
- Correct methodology: absence of threshold in both primary tool sources is evidence of no consensus standard from those sources. However, "no universally accepted threshold exists" is a broader claim. Academic literature (e.g., Jia & Harman survey) was not opened. The claim is reasonable but technically based on limited source inspection. Appropriately hedged in research text: "Sources consulted did not establish a consensus standard."
