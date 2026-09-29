# Claim Audit

Target Lab: labs/38-mutation-testing
Audit Scope: Research Files Only (PIPELINE OVERRIDE)

---

## Claim 1

Claim: Traditional code coverage (line, statement, branch) measures only which code is executed, not whether tests can detect faults.

Location: research/05-report.md — Finding 1; research/03-evidence.md — Evidence 4

Evidence Provided: PIT documentation verbatim quote: "Traditional test coverage (i.e line, statement, branch, etc.) measures only which code is executed by your tests. It does not check that your tests are actually able to detect faults in the executed code."

Source: PIT (https://pitest.org/) — verified reachable. Stryker (https://stryker-mutator.io/docs/) — verified reachable.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Verified directly from PIT homepage text. Stryker provides the "chocolate paste" analogy corroborating the same point. Both sources are official tool documentation. This is the central thesis of the lab and is well-supported.

---

## Claim 2

Claim: Mutation Score = (Killed Mutants / Total Mutants) × 100%. A high score indicates robust tests; a low score reveals weak assertions.

Location: research/05-report.md — Finding 2; research/03-evidence.md — Evidence 3

Evidence Provided: Wikipedia: "The value of a test suite is measured by the percentage of mutants that it kills." PIT FAQ: "The quality of your tests can be gauged from the percentage of mutations killed."

Source: Wikipedia (https://en.wikipedia.org/wiki/Mutation_testing) — accessible; PIT (https://pitest.org/) — verified.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Formula is standard. Wikipedia states it as a ratio. No authoritative source disputes this calculation. Martin Fowler DRAFT entry cited in-line with appropriate qualifier.

---

## Claim 3

Claim: Mutation operators include statement deletion, arithmetic operator replacement, relational operator replacement, Boolean replacement, and statement insertion.

Location: research/05-report.md — Finding 3; research/03-evidence.md — Evidence 5

Evidence Provided: Wikipedia lists operators (citing Hamimoune & Falah, 2016). Stryker output shows "BinaryOperator" and "RemoveConditionals" mutators. PIT FAQ lists DEFAULTS, STRONGER, ALL mutator groups.

Source: Wikipedia, Stryker docs (verified), PIT FAQ.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: The categorization is encyclopedic and well-established across independent sources.

---

## Claim 4

Claim: Mutation testing is based on two hypotheses: Competent Programmer Hypothesis and Coupling Effect.

Location: research/05-report.md — Finding 4; research/03-evidence.md — Evidence 2

Evidence Provided: Wikipedia verbatim: "The first is the competent programmer hypothesis. This hypothesis states that competent programmers write programs that are close to being correct. The second hypothesis is called the coupling effect. The coupling effect asserts that simple faults can cascade or couple to form other emergent faults."

Source: Wikipedia (citing DeMillo et al. 1978; Offutt 1992; Acree et al. 1979) — reachable.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: DeMillo et al. 1978 original paper was not directly read, but the Wikipedia page cites this claim accurately with established academic citations. Offutt 1992 (coupling effect formalization) and Acree et al. (1979) are secondary-attributed. This limitation is acknowledged in the research.

---

## Claim 5

Claim: The RIP model (Reach, Infect, Propagate) defines the three conditions required for a test to kill a mutant.

Location: research/05-report.md — Finding 5; research/03-evidence.md — Evidence 7

Evidence Provided: Wikipedia verbatim on RIP model (citing Offutt & Untch, 2000 — Mutation 2000, with link to cs.gmu.edu paper).

Source: Wikipedia (https://en.wikipedia.org/wiki/Mutation_testing)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Standard terminology in mutation testing literature. The primary paper (Offutt & Untch 2000) was not directly read, but the GMU URL is cited in Wikipedia, providing credible academic context. The term is well-established in the field.

---

## Claim 6

Claim: Equivalent mutant detection is mathematically undecidable and is the primary obstacle to industrial-scale mutation testing.

Location: research/05-report.md — Finding 6; research/03-evidence.md — Evidence 6

Evidence Provided: Wikipedia (citing sources [18], [19]). Meta Engineering Blog: "Determining whether a mutant is equivalent or not is known to be mathematically undecidable."

Source: Wikipedia + Meta Engineering Blog (https://engineering.fb.com/2025/09/30/...) — verified reachable.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Undecidability claim is mathematically well-established (reduction to the Halting Problem) and is confirmed by two independent sources: Wikipedia's reference list and the Meta Engineering blog by a researcher specializing in the field.

---

## Claim 7

Claim: Meta ACH system achieves 0.79/0.47 precision/recall for equivalent mutant detection, rising to 0.95/0.96 with simple static analysis preprocessing. Trial: 9,095 mutants, 571 tests, 10,795 classes, 7 platforms, 73% acceptance, 36% privacy relevance.

Location: research/05-report.md — Finding 7; research/03-evidence.md — Evidence 12

Evidence Provided: arXiv abstract (https://arxiv.org/abs/2501.12862) verified; contains all claimed numeric values verbatim.

Source: arXiv:2501.12862 (verified reachable) + Meta Engineering Blog (verified reachable, Sep 30, 2025).

Source Actually Supports Claim: YES

Classification: FACT (empirical; from specific trial context)

Severity: LOW

Notes: All numeric figures in the abstract were verified directly. Important: the claim notes these are Kotlin/privacy/Meta-specific results. The research correctly qualifies the confidence as MEDIUM (single-vendor trial, industry preprint). The arXiv submission was to FSE 2025 Industry Track — peer review status at time of writing is unclear but the preprint is publicly available. The research appropriately caveats this.

---

## Claim 8

Claim: Mutation testing adoption has been blocked by five major barriers: scalability, unrealistic mutants, equivalent mutants, computational cost, and overstretching.

Location: research/05-report.md — Finding 7; research/03-evidence.md — Evidence 9

Evidence Provided: Meta Engineering Blog (Mark Harman) — verified and detailed five distinct barriers with explanatory subsections.

Source: Meta Engineering Blog (verified reachable).

Source Actually Supports Claim: YES (Claim correctly attributed as Harman/Meta's framing, not a universal academic taxonomy)

Classification: INTERPRETATION (Meta/Harman's framing, not universal academic classification)

Severity: LOW

Notes: The revision (Revision 4 in `research-revision/02-changes-made.md`) addressed prior overgeneralization risk by explicitly attributing this as Harman/Meta's articulation. The report currently reads "According to Harman (Meta, 2025)..." framing applied. No independent academic source was cross-checked for this exact five-part taxonomy; however, the barriers themselves are individually supported by multiple sources (PIT FAQ on cost; Wikipedia on equivalent mutants; etc.).

---

## Claim 9

Claim: Go mutation testing tools exist (go-mutesting, gremlins) but lack feature parity with PIT and Stryker.

Location: research/05-report.md — Finding 8; research/03-evidence.md — Evidence 13

Evidence Provided: GitHub repositories cited; PIT FAQ explicitly states Java/Kotlin support only; Stryker supports JS/TS, C#, Scala (no Go).

Source: go-mutesting GitHub, gremlins GitHub/gremlins.dev, PIT FAQ, Stryker docs.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: The claim about maturity gap is reasonable and appropriately framed as the researcher's assessment based on what features are documented vs. not documented for each tool.

---

## Claim 10

Claim: Subsumed mutants are mutants at the same source code location as another mutant; they do not contribute to coverage metrics.

Location: research/05-report.md — Finding 10; research/03-evidence.md — Evidence 14; research/02-sources.md — Source 13

Evidence Provided: Wikipedia verbatim on subsumed mutants.

Source: Wikipedia (https://en.wikipedia.org/wiki/Mutation_testing)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: A narrow but well-supported definition from Wikipedia. It is correctly framed as an advanced refinement topic that does not undermine core concepts.

---

## Claim 11

Claim: There is no universally accepted mutation score threshold; proposed values (80%, 85%, 90%) vary by project.

Location: research/05-report.md — Finding 11

Evidence Provided: PIT docs and Stryker docs — neither provides a recommended minimum threshold.

Source: PIT FAQ, Stryker docs.

Source Actually Supports Claim: PARTIAL

Classification: FACT (negative evidence; absence of standard threshold noted)

Severity: LOW

Notes: The claim is supported by the absence of threshold guidance in major tool documentation. However, the numeric proposals (80%, 85%, 90%) are stated as "variously proposed elsewhere but NOT VERIFIED here" — meaning the research acknowledges it cannot cite where those numbers come from. This is acceptable: the point is that no consensus exists, not that the numbers are authoritative. The research correctly labels them as unverified, which is accurate.

---

## Claim 12

Claim: Mutation testing was first proposed by Richard Lipton in 1971 and formally published by DeMillo, Lipton, and Sayward in 1978.

Location: research/03-evidence.md — Evidence 1

Evidence Provided: Wikipedia citing [1] (the DeMillo 1978 paper). Martin Fowler bliki (draft).

Source: Wikipedia + Martin Fowler (draft).

Source Actually Supports Claim: YES (via Wikipedia summary of the paper)

Classification: FACT (via secondary source; original paper not read)

Severity: LOW

Notes: Historical attribution is widely agreed upon across sources including Wikipedia and Fowler. The primary paper was not directly accessed, but the bibliographic reference from Wikipedia (IEEE Computer, 11(4):34-41, April 1978) is clearly identified. This is the only primary source not directly inspectable, but the limitation is explicitly disclosed.

---

## Claim 13

Claim: Mutation testing is white-box testing intended to help develop effective regression tests.

Location: research/05-report.md — Executive Summary; research/03-evidence.md — Evidence 10

Evidence Provided: Wikipedia verbatim: "Mutation testing is a form of white-box testing. Its purpose is to help the tester develop effective regression tests..."

Source: Wikipedia (citing Ostrand 2002; Misra 2003)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Standard definition. Corroborated by all tool documentation and Fowler's bliki.
