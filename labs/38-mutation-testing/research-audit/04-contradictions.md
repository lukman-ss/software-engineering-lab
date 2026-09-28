# Contradictions Audit

## Summary

No material contradictions were found. The research's own contradictions analysis file (`04-contradictions.md`) is sound, complete, and correctly resolves three apparent tensions.

---

## Auditor Review of Research Contradictions File

### Tension 1: "Mutation testing requires a test to already exist" vs. "LLMs can generate tests from mutants"

RESEARCH ASSESSMENT: Not a contradiction; LLM systems extend workflow rather than replace it.

AUDITOR VERIFICATION: CONFIRMED. The Meta blog post clarifies that ACH still requires existing test infrastructure as an oracle. No contradiction.

### Tension 2: "Equivalent mutants mathematically undecidable" vs. "0.95/0.96 precision/recall"

RESEARCH ASSESSMENT: Theory vs. practice; heuristics can perform well despite theoretical undecidability.

AUDITOR VERIFICATION: CONFIRMED. Theoretical undecidability (Rice's Theorem domain) does not preclude high-precision approximations.

### Tension 3: Computational cost concerns (PIT FAQ) vs. tool claims of speed (Stryker)

RESEARCH ASSESSMENT: Relative improvement over older systems; incremental targeting mitigates this.

AUDITOR VERIFICATION: CONFIRMED. PIT FAQ's exact words: "PIT is fast compared to other mutation testing systems, but that can still mean that things will take a while." Stryker's "easy to use and fast to run" refers to its relative improvement over previous-generation tools.

---

## Additional Contradictions Checked by Auditor

### Contradiction A: Martin Fowler DRAFT citation used as confirmed reference

Statement A:
"Martin Fowler confirms: 'Each run makes a small modification...'" — `research/05-report.md:29`

Statement B:
"Page carries 'This is a draft entry' notice" — `research/02-sources.md:29`

Type:
INTERNAL

Impact:
MEDIUM. The research simultaneously acknowledges the DRAFT status (correctly) in `02-sources.md`, `05-report.md` Areas of Disagreement, and `06-open-questions.md`, but continues to use Fowler's text as supporting evidence throughout `03-evidence.md` and `05-report.md` without flagging the draft caveat in-line.

Assessment:
MINOR INCONSISTENCY. The research is self-aware but does not consistently embed the DRAFT caveat in each citation site. Should not be misconstrued as fabricated — the content is real and the draft concern is acknowledged.

### Contradiction B: Stryker evidence described as from `/docs/` vs. `/docs/General/example/`

Statement A:
Evidence for 100% coverage vs 60% mutation score attributed to Stryker source listed as Source 7 (stryker-mutator.io/docs/General/example/)

Statement B:
Several claims in `03-evidence.md` attribute the "sandwich/paste" analogy to `https://stryker-mutator.io/docs/` (Source 6), whereas it appears on the main docs introduction page.

Type:
INTERNAL

Impact:
LOW. Both Source 6 (introduction) and Source 7 (example page) were reviewed and are consistent. The analogy appears in the introduction.

Assessment:
MINOR. No incorrect attribution; no wrong claim.
