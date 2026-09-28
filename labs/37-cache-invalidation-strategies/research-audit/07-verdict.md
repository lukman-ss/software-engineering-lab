# Audit Verdict

Target Lab: `labs/37-cache-invalidation-strategies`

Audit Date: 2026-09-28

Auditor Note: PIPELINE OVERRIDE applied. Code/implementation audit excluded from this stage. Code Correctness and Tests quality gates rated NOT_APPLICABLE.

## Summary

Major Claims Reviewed: 14
Sources Reviewed: 11
Unsupported Claims: 1 (Claim 6 — XFetch optimality partially supported; primary paper proofs NOT VERIFIED)
Contradictions: 5 (C1 critical formula error; C2–C5 resolved tensions)
Code Issues: NOT APPLICABLE (pipeline override)
Test Failures: NOT APPLICABLE (pipeline override)
Research Gaps: 6 (1 blocking, 5 non-blocking)

## Quality Gates

Source Integrity:
WARNING

Rationale: 3 of 11 sources pass with warnings or fail:
- Source 04 and Wikipedia sources (04–07) are Tier-3; acceptable for context but not sole authority for formal claims.
- Source 08/09 (Vattani et al. PDF) reachable but unextractable; optimality proof NOT VERIFIED.
- Source 10 (Redis/AWS docs) FAIL — all attempted URLs returned 404/403.

Claim Support:
WARNING

Rationale: The majority of major claims are well-supported (Cache-Aside mechanics, singleflight semantics, RFC 5861 SWR, stampede definition). The single HIGH-severity gap is the XFetch optimality claim resting on secondary sources (Wikipedia + DOI metadata) with the primary PDF unverifiable. The formula sign error (Contradiction 1) represents an identified-and-documented issue rather than a silent unsupported claim.

Internal Consistency:
PASS

Rationale: Research Agent self-identified all material contradictions (C1–C6) and documented them explicitly. There are no hidden contradictions not addressed in `04-contradictions.md`. The research report is internally consistent with its evidence file.

Code Correctness:
NOT_APPLICABLE

(Pipeline override — code audit excluded)

Tests:
NOT_APPLICABLE

(Pipeline override — test execution excluded)

Documentation Accuracy:
PASS

Rationale: Research report (`05-report.md`) accurately reflects evidence (`03-evidence.md`) and contradictions (`04-contradictions.md`). No research file misrepresents its sources. Confidence ratings (HIGH/MEDIUM/LOW) match actual evidence quality. Synthetic/unverified numbers are explicitly labeled.

## Blocking Issues

1. **Gap 2 — XFetch formula sign error**: The lab specification formula `Δ·β·ln(rand()) > TTL_remaining` is ALWAYS FALSE for `rand() ∈ (0,1)` because `ln(rand())` is negative. Any implementation using this formula verbatim will silently disable probabilistic early expiration. The correct form is `-Δ·β·ln(rand()) > TTL_remaining`. Research correctly identified this error in `04-contradictions.md` C1, but the lab specification formula has NOT been corrected. The research must explicitly communicate the corrected formula to the engineering stage to prevent implementation failure.

## Non-Blocking Issues

1. **Gap 1 — XFetch optimality NOT VERIFIED from primary PDF**: Vattani et al. PVLDB 2015 PDF parse failed. Optimality claim accepted on bibliographic authority + Wikipedia summary. This does not invalidate the algorithm description but weakens the "proven optimal" assertion.

2. **Gap 3 — Redis official docs unreachable**: No Redis-specific pattern guidance verified. Redis command behavior (`SET NX PX`, `TTL key`) not verified from redis.io. Claims use Microsoft Azure/general documentation as proxies.

3. **Gap 4 — Synthetic benchmark numbers**: 10,000 RPS, 500 goroutines, P99 figures are lab exercise parameters, not production measurements. Research correctly labels them synthetic.

4. **Gap 5 — RFC 5861 status terminology**: RFC 5861 is Informational, not a formal IETF Standard. Research correctly identifies this but downstream consumers should not call it an "Internet Standard".

5. **Gap 6 — Jitter-sufficiency claim is inferential**: The claim that jitter alone fails to prevent single hot-key stampede is a sound inference but lacks a direct primary source quote.

## Required Revisions

1. **CRITICAL**: In research report (`05-report.md`) Finding 6, explicitly replace the lab formula with the mathematically correct form: `-Δ·β·ln(rand()) > TTL_remaining` (equivalently `Δ·β·(-ln(rand())) > TTL_remaining`). Include a warning block: "WARNING: The original lab specification formula omits the negation sign. Do not implement `Δ·β·ln(rand()) > TTL_remaining`; it is always false."

2. **MEDIUM**: In `05-report.md` Finding 6, add an explicit disclaimer: "The XFetch optimality claim is accepted on bibliographic authority alone. Primary paper mathematical proofs and experimental benchmarks were not independently extracted due to PDF parsing failure."

3. **LOW**: Update `02-sources.md` Source 10 to note that Redis documentation URLs were unreachable and indicate alternative sources or current URL patterns that would be needed to verify Redis-specific claims.

## Final Status

NEEDS_REVISION
