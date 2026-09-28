# Revision Plan

Target Lab: labs/37-cache-invalidation-strategies

Previous Audit Status: NEEDS_REVISION

## Blocking Issues

1. **Gap 2 — XFetch formula sign error (CRITICAL)**
   - Lab formula `Δ·β·ln(rand()) > TTL_remaining` is ALWAYS FALSE for `rand() ∈ (0,1)` because `ln(rand())` is negative
   - Must correct in research report with explicit warning

## Non-Blocking Issues

1. **Gap 1 — XFetch optimality NOT VERIFIED from primary PDF (HIGH)**
   - Primary paper proofs not independently extracted
   - Need explicit disclaimer in Finding 6

2. **Gap 3 — Redis official docs unreachable (MEDIUM)**
   - Source 10 URLs returned 404/403
   - Update source documentation

3. **Gap 4 — Synthetic benchmark numbers (LOW)**
   - Already correctly labeled in research

4. **Gap 5 — RFC 5861 status terminology (LOW)**
   - Already correctly called "Informational" in report

5. **Gap 6 — Jitter-sufficiency claim is inferential (MEDIUM)**
   - Need to label as inferential deduction, not primary source claim

## Files To Modify

- research/05-report.md (Finding 6, Finding 8)
- research/02-sources.md (Source 10)
- research/03-evidence.md (Evidence 16)

## Verification Plan

- Verify source corrections against audit findings
- Confirm formula correction is mathematically correct
- Ensure disclaimers are prominent and accurate
- No code changes required (pipeline override)
- Mark as READY_FOR_RESEARCH_REAUDIT upon completion