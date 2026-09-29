# Research Gap Analysis

## Summary of Gaps

This audit identified 4 research gaps across the target lab's research files. None are classified as CRITICAL. One is classified as HIGH (source availability failure). The remaining are classified as MEDIUM or LOW.

---

## Gap 1

Type: WEAK_SOURCE / BROKEN_URL

Severity: HIGH

Location: `research/02-sources.md` (Source 3) and `research/03-evidence.md` (Evidence 1, 3)

Problem:
Martin Fowler's bliki URL `https://martinfowler.com/bliki/MutationTesting.html` returned HTTP 404 during this audit. The Research Agent acknowledged it was a "pre-publication draft" carrying a "This is a draft entry" notice. A draft URL that is now unreachable should not be cited as an active Tier 1 source without verification of its current availability or retrieval via archive.org / alternate authoritative source.

Required Revision:
Either remove Source 3 and rely on the other primary sources (PIT, Stryker, Wikipedia) that fully support the claims, or update the source entry with an archived URL (`web.archive.org`) or alternative citation from Fowler's published works (e.g., Refactoring 2nd ed. or testing guides).

Can Be Approved Without Fix: NO (in a full production pipeline; acceptable with WARNING if non-blocking because core claims have independent multi-source backing).

---

## Gap 2

Type: UNVERIFIED_CLAIM / NUMERIC_RECOMMENDATIONS

Severity: MEDIUM

Location: `research/05-report.md` (Finding 11) and `research/06-open-questions.md` (Weak Evidence 2)

Problem:
Finding 11 mentions specific threshold percentages: "proposed values (80%, 85%, 90%) vary by project and context." The Research Agent flagged this in `06-open-questions.md` as "NOT VERIFIED," but included the specific numbers parenthetically in the main report without a source. While the primary finding ("no standard threshold exists") is well-supported, mentioning specific numeric ranges without attribution introduces arbitrary numbers.

Required Revision:
Remove specific numbers `(80%, 85%, 90%)` from Finding 11 or explicitly label them as illustrative community heuristics with a citation or explicit disclaimer.

Can Be Approved Without Fix: YES (the primary conclusion that no standard exists is sound; numbers are illustrative).

---

## Gap 3

Type: MISSING_SOURCE / SECONDARY_ATTRIBUTION

Severity: LOW

Location: `research/02-sources.md` (Source 1, Source 9)

Problem:
The foundational 1978 DeMillo, Lipton, and Sayward paper and the 2009 Jia & Harman survey are cited based exclusively on Wikipedia's reference list. The Research Agent explicitly declared "Paper text NOT accessed directly" and "NOT VERIFIED for direct quotations." While this disclosure is honest, relying on secondary citations for historical paper details is an attribution weakness.

Required Revision:
For higher-rigor publication, access the papers via open-access repositories (e.g., IEEE Xplore, Semantic Scholar open access, author personal copies) to verify key quotations directly.

Can Be Approved Without Fix: YES (the Research Agent properly disclosed the secondary attribution; the bibliographic metadata is accurate).

---

## Gap 4

Type: SCOPE_ERROR / GENERALIZATION

Severity: LOW

Location: `research/03-evidence.md` (Evidence 8); `research/05-report.md` (Executive Summary)

Problem:
Claim 8 states: "All major mutation testing tools follow the same pattern: generate mutants → run tests → classify as killed/survived." While true for standard mutation testing engines, modern LLM-augmented tools (such as ACH) invert or extend this pattern by using mutants to *generate* tests rather than just running existing ones.

Required Revision:
Qualify the statement to clarify that "all major *evaluation-focused* mutation testing tools" follow this pattern, while newer *generative* approaches (like ACH) extend it.

Can Be Approved Without Fix: YES (context throughout report makes the distinction clear).
