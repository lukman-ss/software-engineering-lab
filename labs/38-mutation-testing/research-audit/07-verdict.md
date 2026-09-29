# Audit Verdict

Target Lab: labs/38-mutation-testing
Audit Date: 2026-09-29

---

## Summary

Major Claims Reviewed: 13
Sources Reviewed: 13 (7 directly verified via URL; 2 unread academic papers via secondary citation; 2 community GitHub projects; 2 Wikipedia pages)
Unsupported Claims: 0
Contradictions: 0 material contradictions
Code Issues: NOT_APPLICABLE (Research-only phase)
Test Failures: NOT_APPLICABLE (Research-only phase)
Research Gaps: 5 (all LOW or MEDIUM; all non-blocking)

---

## Quality Gates

Source Integrity:
PASS
— All 9 directly-accessible URLs verified and confirmed reachable. Two paywalled/academic papers explicitly disclosed as secondary-attributed only. No fabricated sources or URLs.

Claim Support:
PASS
— All 13 audited claims are supported by accessible sources or explicitly disclosed as unverifiable secondary citations. Confidence levels are accurately assigned. No CRITICAL or HIGH unsupported claims.

Internal Consistency:
PASS
— No material contradictions between research files, evidence log, sources list, and open questions document. Theoretical tension points (undecidability vs. heuristic approximation; LLM test generation vs. "requires existing tests") are clearly resolved and contextualized.

Code Correctness:
NOT_APPLICABLE
— Pipeline override: research-only audit stage.

Tests:
NOT_APPLICABLE
— Pipeline override: research-only audit stage.

Documentation Accuracy:
PASS
— Research revision cycle resolved all prior audit flags: Martin Fowler draft properly qualified in-line, Meta ACH stats attributed to arXiv-verifiable abstract, "five barriers" framing attributed to Harman/Meta explicitly, Go tooling gap filled with community sources, subsumed mutants documented, and threshold ambiguity explicitly flagged as unresolved in the literature.

---

## Blocking Issues

None.

---

## Non-Blocking Issues

1. **Foundational papers unread directly**: DeMillo et al. (1978) and Jia & Harman (2009) cited via Wikipedia secondary references. No direct text verification of the original papers. Transparently disclosed in source notes with "NOT VERIFIED for direct quotations."

2. **Draft source**: Martin Fowler bliki entry carries "This is a draft entry" banner. Properly qualified in all inline citations.

3. **Single-vendor trial data for ACH**: Meta ACH performance statistics (73% acceptance, 0.95/0.96 precision/recall) originate from a single internal Meta trial on Kotlin Android privacy testing. The arXiv preprint abstract corroborates the blog post claims. Findings are appropriately scoped as implementation-specific.

4. **FSE 2025 peer review status**: The arXiv preprint (2501.12862) was submitted to FSE 2025 Industry Track. Full peer review outcome is unknown at time of research and is noted in the limitations section.

---

## Required Revisions

None. Research is approved as-is.

---

## Final Status

APPROVED_WITH_WARNINGS

The research for labs/38-mutation-testing is technically sound and ready to serve as the foundation for lab implementation or publication. All major claims are either directly verified against reachable primary/official sources, or transparently disclosed as secondary-attributed with appropriate confidence downgrades. The source inventory is comprehensive for the scope of the topic. Internal consistency is maintained across all research documents. Research gaps are explicitly catalogued in `06-open-questions.md` and were addressed in the prior revision cycle (`research-revision/`).

Warnings are limited to inherent limitations of the research pipeline:
- Two key historical papers remain unread directly (paywalled academic literature).
- One industry source is still in editorial draft stage.
- One industry trial is from a single vendor and is scoped to a specific platform and domain.

None of these warnings undermine the core claims, which are independently corroborated by official tool documentation (PIT, Stryker), Wikipedia, and a publicly accessible arXiv preprint.
