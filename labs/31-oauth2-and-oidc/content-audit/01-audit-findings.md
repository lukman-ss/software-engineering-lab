# Content Audit Findings

## Overview

**Target Lab:** labs/31-oauth2-and-oidc  
**Audit Scope:** Content files only (01-content-brief.md through 06-source-map.md)  
**Reference Materials:** Research report, Engineering design/implementation/test results, Engineering audit, Source code

## Issues Identified

### 1. HIGH - State Parameter CSRF Protection Claim Mismatch (DOC_CODE_MISMATCH)

**Location:** `content/02-master-draft.md` lines 48, 277, 341

**Claim in Content:**
- Line 48: "code_challenge=S256, nonce" - no mention of state
- Line 277: Checklist item "State parameter divalidasi (untuk anti-CSRF)"
- Line 341: Checklist item "State parameter divalidasi (untuk anti-CSRF)"

**Reality from Code Audit:**
- `pkg/client/client.go:47` generates `c.State` but it is **never sent to the server** during authorization request
- `pkg/server/server.go` does not receive or validate any state parameter
- Engineering Audit Gap 3: "Client `State` parameter generated but not passed to server"

**Impact:** The content claims state parameter validation for CSRF protection, but the implementation does not enforce CSRF mitigation. This is misleading for readers implementing the same pattern.

---

### 2. MEDIUM - ID Token Validation Steps Count Discrepancy

**Location:** `content/02-master-draft.md` lines 71-78, 227-248

**Research Claim (research/05-report.md Finding 3):**
- "ID Token wajib divalidasi multi-langkah (signature + iss + aud + exp + nonce)" - 13 steps per OIDC Core 3.1.3.7

**Implementation Reality:**
The lab ID Token validation (`pkg/oidc/oidc.go:64-116`) implements only **7 validation checks**:
1. JWT format (3 parts)
2. HMAC-SHA256 signature verification (`hmac.Equal`)
3. Issuer exact match
4. Audience exact match
5. Expiration check
6. Issued-at clock skew check (±5 min)
7. Nonce match

**Missing from Implementation (per full 13-step OIDC validation):**
- Optional decryption for encrypted ID Tokens (JWE)
- JWKS-based signature verification (lab uses symmetric HS256)
- `alg` claim pinning/checking
- `acr` (Authentication Context Class Reference)
- `auth_time` (authentication time) claims
- `at_hash` / `c_hash` (access/authorization code hash) validation
- `azp` (authorized party) validation

**Impact:** The content implies full OIDC Core 3.1.3.7 compliance. The lab intentionally uses a simplified HS256 model for educational purposes. This should be explicitly clarified rather than implied.

---

### 3. LOW - PKCE Verification Not Using Constant-Time Comparison (Already Documented as Warning)

**Location:** `content/02-master-draft.md` line 266, `content/03-code-snippets.md` line 72

The content **correctly documents** this as a LOW severity issue:
- Content-brief line 40: warns about "PKCE comparison uses string equality (not constant-time)"
- Master-draft line 266: mentions in Production Considerations to use `crypto/subtle.ConstantTimeCompare`
- Code-snippets line 72: explicitly notes "(Note: audit menandai penggunaan string equality bukan constant-time compare sebagai LOW severity issue.)"

**Status:** ACCURATELY DOCUMENTED

---

### 4. LOW - Audience Validation Discrepancy

**Location:** `content/03-code-snippets.md` line 174

**Claim:** "audience exact match"

**Implementation:** Uses exact string comparison (`claims.Audience != expectedAudience`)

**Research/OIDC Spec (OIDC Core 3.1.3.7):** Specifies `aud` must contain the client_id (for multiple audiences case)

**Impact:** Minor - exact match is more restrictive and safe. This is a valid implementation choice but differs from the spec's "contains" requirement.

---

### 5. ACCURATE REFERENCES - Verified

- **Test count (17 tests):** All content claims match actual test file (tests/oauth_test.go)
- **Test names:** All test scenarios listed in content match actual test functions
- **Demo output:** Content's demo walkthrough matches engineering/03-execution-result.md output exactly
- **RFC references:** RFC 6749, 7636, 7519, 8725, 9700, OIDC Core 1.0 correctly cited
- **Key takeaways:** Accurate summarization of lab's scope and limitations

---

## Summary of Findings

| Issue | Severity | Status |
|-------|----------|--------|
| State parameter CSRF claim mismatch | HIGH | Content claims validation but not implemented |
| ID Token validation steps count discrepancy | MEDIUM | Claims 13 steps but implements ~7 (intentional simplification) |
| PKCE constant-time comparison | LOW | Correctly flagged as warning in content |
| Audience validation exact vs contains | LOW | Implementation choice, not incorrect |
| All other content | N/A | Accurately reflects research, engineering, and tests |

---

## Content Quality Assessment

**Strengths:**
- Comprehensive coverage of OAuth 2.0 vs OIDC distinction
- Clear explanation of PKCE S256 mechanism
- Accurate code snippets with correct line references
- Proper documentation of security warnings (non-constant-time PKCE comparison)
- Good technical depth in ID Token validation explanation
- Accurate test documentation
- Verified demo output

**Weaknesses:**
- Overstates CSRF protection by claiming state validation (which isn't implemented)
- Overstates ID Token validation completeness (omits mention of simplified HS256 model)
- Checklist item implies full CSRF protection without caveat

---

## Recommendation

The content is **substantially accurate** but has:
1. One critical misrepresentation (State parameter validation)
2. One medium-issue omission (ID Token validation scope clarification)

These issues should be addressed in a revision before publication to prevent readers from implementing missing CSRF protections.