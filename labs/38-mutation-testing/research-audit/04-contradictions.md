# Contradictions Analysis

Target Lab: `/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/38-mutation-testing`
Audit Date: 2026-09-29

---

## Contradiction Analysis Summary

No material contradictions were found across research files, nor between research files and README.md.

---

## Verification of Potential Tension Points

### Point 1: Research vs README Operator Scope
Statement A: `README.md` lists 4 implemented mutation operators: (1) Relational operator replacement, (2) Boolean flip, (3) Arithmetic operator replacement, (4) Boundary value shift.
Location: `README.md:48-53`
Statement B: `research/05-report.md` notes 5 mutation types for practical implementation including statement deletion.
Location: `research/05-report.md:45`
Type: CODE_DOC_MISMATCH (minor pedagogical nuance)
Impact: None. Research identifies the broad theoretical catalog (statement deletion is one of the classic operators), while the lab's specific Go implementation chooses a subset of 4 AST operators suited to its discount engine domain.
Assessment: PASS. Not a contradiction; README scopes the actual Go implementation, while research covers the broader literature.

---

### Point 2: Undecidability vs Empirical Detection Rates
Statement A: Semantic equivalence detection is mathematically undecidable in the general case.
Location: `research/05-report.md` Finding 6; `research/04-contradictions.md` Point 2
Statement B: Meta ACH reports 0.95 precision and 0.96 recall in equivalent mutant detection with preprocessing.
Location: `research/05-report.md` Finding 7; `research/03-evidence.md` Evidence 12
Type: INTERNAL / THEORETICAL_VS_EMPIRICAL
Impact: None.
Assessment: PASS. This tension is explicitly analyzed in `research/04-contradictions.md` Point 2. Undecidability is a general-case theoretical bound (Halting problem analogue), whereas empirical classification heuristics (LLM + static analysis) operate on specific code corpora without formal totality guarantees. Both coexist consistently.

---

### Point 3: Execution Models Across Tooling
Statement A: Evaluation-focused mutation testing executes mutants against existing test suites.
Location: `research/05-report.md` Finding 8; `research/03-evidence.md` Evidence 8
Statement B: Generative systems like Meta ACH use unkilled mutants to prompt LLMs to synthesize new tests.
Location: `research/05-report.md` Finding 7; `research/03-evidence.md` Evidence 12
Type: INTERNAL (workflow taxonomy)
Impact: None.
Assessment: PASS. Resolved in research revision by scoping traditional tools as "evaluation-focused" and distinguishing generative paradigms.

---

### Point 4: Martin Fowler Draft Entry Status
Statement A: Martin Fowler bliki URL marked as `REMOVED (UNREACHABLE)` in `research/02-sources.md`.
Location: `research/02-sources.md` Source 3
Statement B: Martin Fowler draft entry referenced in `research/05-report.md` Finding 2 Sources field.
Location: `research/05-report.md:31`
Type: INTERNAL (residual citation)
Impact: LOW. Finding 2 explicitly notes the URL was a pre-publication draft and that the claim is independently supported by Wikipedia and PIT. The revision documentation acknowledges this as non-blocking.
Assessment: WARNING (minor residual string in report, fully harmless given multi-source support).

---

## Conclusion

No material contradictions found. All investigated tension points are properly reconciled in the research text.
