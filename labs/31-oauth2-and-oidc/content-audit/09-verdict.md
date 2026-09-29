# Content Audit Verdict

**Target Lab:** labs/31-oauth2-and-oidc  
**Audit Scope:** Content files only (01-content-brief.md through 06-source-map.md)  
**Audit Date:** 2026-09-29

## Summary

Content files reviewed: 6
- `01-content-brief.md`
- `02-master-draft.md`
- `03-code-snippets.md`
- `04-diagrams.md`
- `05-key-takeaways.md`
- `06-source-map.md`

Reference materials cross-checked:
- `research/05-report.md` (Research Report - APPROVED)
- `engineering/01-design.md` (Design)
- `engineering/02-implementation-notes.md` (Implementation Notes)
- `engineering/03-execution-result.md` (Execution Results - 17/17 tests PASS)
- `engineering-audit/02-code-audit.md` (Code Audit)
- `engineering-audit/04-docs-vs-code.md` (Docs vs Code)
- `engineering-audit/05-gaps.md` (Gap Analysis)
- `pkg/pkce/pkce.go`, `pkg/oidc/oidc.go`, `pkg/server/server.go`, `pkg/client/client.go` (Source code)
- `tests/oauth_test.go` (Test suite - 17 tests)

---

## Detailed Findings

### Findings Requiring Revision

1. **HIGH - State Parameter CSRF Protection Claim Mismatch** (`content/02-master-draft.md` lines 277, 341)

   The content checklist claims "State parameter divalidasi (untuk anti-CSRF)" but the implementation does NOT validate the state parameter:
   - `pkg/client/client.go:47` generates `c.State` but never sends it to `server.Authorize()`
   - `pkg/server/server.go` has no state parameter in any function signature
   - Engineering Audit Gap 3 (LOW): "Client `State` parameter generated but not passed to server"
   - Code Audit Finding 10: `State is stored but never verified server-side`

   This is misleading and overstates the lab's security posture.

2. **MEDIUM - ID Token Validation Steps Claim Overstated** (`content/02-master-draft.md` lines 71-78)

   The research report (Finding 3) cites OIDC Core 3.1.3.7 with 13 validation steps. The implementation (`pkg/oidc/oidc.go:64-116`) performs only 7 checks and uses simplified HS256 symmetric signing rather than JWKS. The content should explicitly state the scope of validation implemented versus the full OIDC spec, rather than implying full OIDC Core compliance.

### Warnings (Accurately Documented)

3. **LOW - PKCE Non-Constant-Time Comparison** (`content/02-master-draft.md` line 266, `content/03-code-snippets.md` line 72)

   The content correctly documents this as a LOW severity issue per the engineering audit.

4. **LOW - In-Memory Storage** (`content/01-content-brief.md` line 38)

   Properly noted as a lab limitation.

5. **LOW - Symmetric HMAC-SHA256 vs RSA/JWKS** (`content/01-content-brief.md` lines 39-40)

   Properly documented as a simplification choice.

### Verified Accurate

- Test count (17) and all test scenario descriptions match `tests/oauth_test.go`
- Demo walkthrough matches `engineering/03-execution-result.md` output exactly
- RFC references are correct (RFC 6749, 7636, 7519, 8725, 9700, OIDC Core)
- PKCE S256 challenge/verifier mechanism accurately described
- Refresh Token Rotation and Family Revocation accurately explained
- Code snippets have correct line references and match actual source code
- Key Takeaways accurately summarize the lab's scope

---

## Verdict

### APPROVED_WITH_WARNINGS

The content is **substantially accurate** and effectively conveys OAuth 2.0/OIDC concepts with correct technical depth. The code snippets match actual implementation, test documentation aligns with the test suite, and demo output is verified against live execution.

**Required revision before final publication:**
1. Remove or qualify the claim "State parameter divalidasi (untuk anti-CSRF)" with an explicit note that the lab omits server-side state validation (architectural demo limitation).
2. Clarify that ID Token validation implements a 7-step subset of OIDC Core 3.1.3.7 (HMAC-SHA256, not JWKS; no alg pinning; no at_hash/c_hash), explicitly marking what is in scope for the lab.

These are not hallucinated facts or platform biases, but rather omissions of precision that could mislead readers about the lab's security posture.
