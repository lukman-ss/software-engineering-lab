# Claim Audit

## Claim 1

Claim: Code coverage metrics measure execution rather than assertion/fault detection capability.

Location: `research/05-report.md` (Finding 1)

Evidence Provided: Quoted statements from PIT homepage ("Traditional test coverage... measures only which code is executed... does not check that your tests are actually able to detect faults") and Stryker ("sandwich covered with paste" analogy).

Source: Source 4 (PIT) and Source 6 (Stryker).

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Directly supported by official documentation of leading mutation testing tools.

---

## Claim 2

Claim: Mutation Score formula is (Killed Mutants / Total Mutants) × 100%.

Location: `research/05-report.md` (Finding 2)

Evidence Provided: Wikipedia definition and PIT documentation ("gauged from the percentage of mutations killed").

Source: Source 2 (Wikipedia) and Source 5 (PIT FAQ).

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Standard mathematical expression of mutation coverage percentage.

---

## Claim 3

Claim: Mutation testing rests on the Competent Programmer Hypothesis and the Coupling Effect.

Location: `research/05-report.md` (Finding 4)

Evidence Provided: Wikipedia definitions citing DeMillo et al. (1978) and Offutt (1992).

Source: Source 2 (Wikipedia).

Source Actually Supports Claim: YES

Classification: FACT / THEORY

Severity: LOW

Notes: Universally acknowledged theoretical foundations of mutation testing.

---

## Claim 4

Claim: The RIP model specifies three conditions (Reach, Infect, Propagate) required for strong mutation testing to kill a mutant.

Location: `research/05-report.md` (Finding 5)

Evidence Provided: Wikipedia summary citing Offutt & Untch (2000).

Source: Source 2 (Wikipedia).

Source Actually Supports Claim: YES

Classification: FACT / MODEL

Severity: LOW

Notes: Standard model for test excitation and assertion propagation.

---

## Claim 5

Claim: Detecting equivalent mutants is mathematically undecidable.

Location: `research/05-report.md` (Finding 6)

Evidence Provided: Wikipedia citations and Meta ACH paper statement ("Determining whether a mutant is equivalent or not is known to be mathematically undecidable").

Source: Source 2 (Wikipedia) and Source 12 (arXiv:2501.12862).

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Fundamental theoretical limit acknowledged across literature.

---

## Claim 6

Claim: Meta's ACH system uses LLMs to achieve 73% engineer acceptance, 36% privacy relevance, and 0.95/0.96 precision/recall for equivalent mutant detection with preprocessing.

Location: `research/05-report.md` (Finding 7)

Evidence Provided: Meta Engineering blog post (Sep 2025) and arXiv preprint abstract (Foster et al., Jan 2025).

Source: Source 8 (Meta Blog) and Source 12 (arXiv:2501.12862).

Source Actually Supports Claim: YES

Classification: IMPLEMENTATION-SPECIFIC / EXPERIMENTAL

Severity: MEDIUM

Notes: Trial statistics are accurate to the source, but restricted to Meta's specific Android Kotlin codebase and privacy test-a-thons. Appropriately contextualized in research revision.

---

## Claim 7

Claim: Go has community mutation testing tools (`go-mutesting`, `gremlins`), but tooling maturity trails JVM (PIT) and JS/TS (Stryker) ecosystems.

Location: `research/05-report.md` (Finding 8)

Evidence Provided: GitHub repositories for `go-mutesting` and `gremlins`, alongside PIT FAQ and Stryker supported language documentation.

Source: Source 5 (PIT FAQ), Source 6 (Stryker), Source 10 (`go-mutesting`), Source 11 (`gremlins`).

Source Actually Supports Claim: YES

Classification: FACT / EVALUATION

Severity: LOW

Notes: Accurately reflects current open-source tooling landscape for Go.

---

## Claim 8

Claim: No universal industry-standard mutation score threshold exists (e.g., 80% vs 90%).

Location: `research/05-report.md` (Finding 11)

Evidence Provided: Review of PIT and Stryker documentation revealing no prescribed target threshold.

Source: Source 5 (PIT FAQ) and Source 6 (Stryker docs).

Source Actually Supports Claim: YES

Classification: FACT / OBSERVATION

Severity: LOW

Notes: Honest identification of an unstandardized metric threshold in industry practice.
