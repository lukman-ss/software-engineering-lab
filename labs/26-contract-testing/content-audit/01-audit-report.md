# Content Audit Report — Lab 26 Contract Testing

Auditor: Kiro (Technical Content Auditor)  
Date: 2026-09-27  
Target: labs/26-contract-testing/content/

---

## Summary

| Quality Gate | Status | Notes |
|--------------|--------|-------|
| Research Alignment | PASS | Content aligns with approved research |
| Core Accuracy | PARTIAL | CDC concepts accurate; header validation overstated |
| Completeness | WARNING | Missing clarification on verification scope |
| Clarity | WARNING | Headers field described but not validated |

---

## Findings

### Finding 1: Headers Validation Overstatement (HIGH)

**Location:** 
- `02-master-draft.md:75` — "Verifier verifies status, header, and body"
- `04-diagrams.md:22` — Diagram implies header comparison in verifier block
- `03-code-snippets.md:131` — "Notes exists... not part of contract" (conflicts with header validation claim)

**Evidence:** `internal/contract/verifier.go:57-115`
- Lines 70-71: Request headers set, sent
- **Lines 87-92**: Only checks `resp.StatusCode` and `body`; **NO header validation loop**
- `ResponseDefinition.Headers` field defined (line 29) but never read in verification

**Impact:** Consumers/readers may expect Content-Type, X-Custom-Header, etc. to be verified; they are not.

**Source:** engineering-audit-opensource/06-verdict.md GAP-01 confirms this as HIGH severity implementation gap.

### Finding 2: V2 Endpoint Documentation Mismatch (MEDIUM)

**Location:**
- Content describes ProviderDual with `/v2` endpoint (Snippet 4, line 147-203)
- Tests only verify V1 path (`TestProviderDual_ContractVerification_Success` uses `/v1`)

**Evidence:**
- `tests/contract_test.go:74-94` — Dual provider test calls `/v1/orders/ORD-123` only
- `internal/provider/server.go:112-128` — V2 endpoint exists but has no contract interaction, no test, no demo

**Impact:** Content advertises V2 safe evolution pattern; code routes V2 but provides no verification feedback loop for consumers migrating to V2.

**Source:** engineering-audit-opensource/06-verdict.md GAP-02 confirms as HIGH severity gap.

### Finding 3: Verification Error Order Non-Deterministic (LOW)

**Location:** `04-diagrams.md:40-45` shows sequential error output

**Evidence:** `verifier.go:131` uses `range` on map (non-deterministic order)

**Impact:** Demo transcripts will show diffs in varying order; test only asserts `>=3` errors, not specific categories.

**Source:** engineering-audit-opensource/06-verdict.md GAP-06 documents this limitation.

---

## Discrepancy Matrix

| Content Statement | Code Reality | Severity |
|-------------------|--------------|----------|
| "Verifier checks headers" | No header checking | HIGH |
| "V2 endpoint verified" | V2 endpoint exists, unverified | HIGH |
| "ProviderState in contract" | Present in JSON, unused by verifier | LOW |
| Request headers sent | True (lines 70-72) | None |
| Response body verified | True | None |
| Status code verified | True | None |
| json.Number type checking | True | None |
| Recursive field diffing | True | None |

---

## Positive Findings

### Correct Content

- CDC definition aligns with research (Finding 2 in research/05-report.md)
- Breaking change taxonomy (enum casing, field rename, type mutation) matches code
- Type precision (json.Number vs string) correctly documented
- Expand/contract pattern accurately described
- Over-specification anti-pattern correctly identified

---

## Recommendations

1. **Add verification scope disclaimer**: State explicitly "The current implementation verifies HTTP status code and response body. Response header validation is a planned enhancement (see engineering-audit-opensource GAP-01)."

2. **Clarify dual-provider status**: Note "V2 endpoint exists in ProviderDual but has no associated consumer contract or verification test. See engineering-audit-opensource GAP-02."

3. **Add non-determinism disclaimer**: "Verifier error ordering depends on map iteration; order may vary between runs."

---

## Dependencies

Content should be cross-referenced with:

- `engineering-audit-opensource/06-verdict.md` — identifies gaps (GAP-01, GAP-02, GAP-06)
- `engineering/02-implementation-notes.md:29-33` — known limitations section
- `engineering-audit/06-verdict.md` — shows APPROVED status for main audit track

---

**Auditor Note:** Content demonstrates strong technical accuracy for CDC concepts. The implementation gaps (headers, V2 verification) represent simplifications typical of educational labs. Content should transparently frame these as "current implementation scope" vs "full production capability."