# Changes Made

## Revision 1

### Audit Issue: CRITICAL — Gap 2, XFetch formula sign error

**Files Changed:**
- research/05-report.md (Finding 6)

**Action:**
- Added explicit `WARNING` block to Finding 6 documenting the lab formula sign error
- Stated clearly: `Δ·β·ln(rand()) > TTL_remaining` is ALWAYS FALSE for `rand() ∈ (0,1)`
- Provided corrected formula: `-Δ·β·ln(rand()) > TTL_remaining` (and equivalent form)
- Included original notation from Wikipedia for reference
- Explicitly directed engineers NOT to implement the unnegated formula

**Verification:**
- Formula derivation verified: `log(rand(0,1))` is negative for rand∈(0,1); `-delta*beta*log(rand)` yields positive offset
- Lab qualitative description ("higher probability closer to expiry / longer query time") is consistent with corrected formula
- C1 contradiction in 04-contradictions.md already documents this; finding now carries explicit engineer-facing warning

**Status:** RESOLVED

---

### Audit Issue: HIGH — Gap 1, XFetch optimality NOT VERIFIED from primary PDF

**Files Changed:**
- research/05-report.md (Finding 6, Evidence paragraph)

**Action:**
- Expanded Evidence paragraph to explicitly state: "primary paper mathematical proofs and experimental benchmarks were not independently extracted due to PDF parsing failure"
- Both PDF URLs (cseweb.ucsd.edu and vldb.org) returned binary streams that could not be decompressed
- Theorem statements and quantitative benchmarks are NOW EXPLICITLY MARKED NOT VERIFIED
- Optimality claim remains accepted on bibliographic authority alone, with this limitation clearly documented

**Verification:**
- Consistent with Evidence 11 which already documents PDF fetch failure
- No new claims introduced; existing MEDIUM confidence retained

**Status:** RESOLVED

---

### Audit Issue: MEDIUM — Gap 6, Jitter-sufficiency claim is inferential

**Files Changed:**
- research/05-report.md (Finding 8)
- research/03-evidence.md (Evidence 16)

**Action:**
- Added `NOTE` to Finding 8 labeling the single-key insufficiency claim as "inferential deduction based on single-key contention mechanics, not a direct quote from a primary caching source"
- Updated Evidence 16 similarly with explicit NOTE
- Changed Confidence label from "MEDIUM for sufficiency claim" to "MEDIUM (inferential) for single-key insufficiency claim"
- Clarified in Evidence 16 notes that the deduction is logically sound but should be labeled as interpretation

**Verification:**
- No primary source states "jitter alone fails to prevent stampede on a single hot key" verbatim
- Wikipedia Thundering herd §Mitigation describes jitter for retry backoff desynchronization but does not address single-hot-key cache stampede specifically
- Deduction is correct; labeling updated to reflect evidence quality accurately

**Status:** RESOLVED

---

### Audit Issue: LOW — Gap 3, Redis official docs unreachable

**Files Changed:**
- research/02-sources.md (Source 10)

**Action:**
- Added `REVISER UPDATE (2026-09-28)` block to Source 10
- Documented that redis.io site structure has moved from `/docs/latest/develop/` to different paths
- Noted alternative source URLs needed for re-verification
- Stated that until re-verified, all Redis-specific claims must be treated as NOT VERIFIED
- Recommended Microsoft Learn (Source 01) as primary vendor guidance for Redis-tier claims

**Verification:**
- Evidence 19 already documents HTTP 404/403 on attempted Redis URLs
- No Redis-official quotations were used as primary evidence for any major claim
- Microsoft Learn provides adequate general cache-aside guidance

**Status:** RESOLVED

---

## Unresolved Items

None. All CRITICAL, HIGH, MEDIUM, and LOW issues from the audit have been addressed.
