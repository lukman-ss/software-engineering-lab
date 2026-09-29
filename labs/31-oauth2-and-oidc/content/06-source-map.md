# Source Map

## OAuth 2.0 vs OIDC Fundamentals

**Research:**
- `research/05-report.md` — Finding 1 (OAuth bukan autentikasi), Finding 2 (OIDC identity layer), Finding 3 (ID Token validation)
- `research/03-evidence.md` — RFC 6749, OIDC Core Sec 1, RFC 8725

**Implementation:**
- `pkg/oidc/oidc.go` — ID Token signing dan validation
- `pkg/server/server.go` — Authorization Server logic

**Tests:**
- `tests/oauth_test.go` — `TestOIDC_IDToken_Valid`, `TestOIDC_IDToken_TamperedSignature`, `TestOIDC_IDToken_Expired`, `TestOIDC_IDToken_MismatchClaims`, `TestOIDC_MalformedJWT`, `TestOIDC_IDToken_IssuedInFuture`

---

## PKCE (Proof Key for Code Exchange)

**Research:**
- `research/05-report.md` — Finding 4 (PKCE prevents interception)
- `research/03-evidence.md` — RFC 7636 Sec 4.1/4.2/4.6

**Implementation:**
- `pkg/pkce/pkce.go` — Generator dan verifier PKCE (S256 dan plain)
- `pkg/server/server.go` — `Authorize()` menyimpan challenge, `ExchangeCode()` memverifikasi

**Tests:**
- `tests/oauth_test.go` — `TestPKCE_S256_Valid`, `TestPKCE_InvalidMethod`, `TestPKCE_Mismatch`, `TestPKCE_Plain_Method`, `TestPKCE_VerifierLength_Bounds`

---

## Authorization Code Flow

**Research:**
- `research/05-report.md` — Finding 5 (Implicit deprecated)
- `research/03-evidence.md` — RFC 9700, OAuth 2.1

**Implementation:**
- `pkg/server/server.go` — `Authorize()`, `ExchangeCode()`
- `pkg/client/client.go` — `BuildAuthorizationRequest()`, `Exchange()`
- `cmd/demo/main.go` — Step 1-3 demo

**Tests:**
- `tests/oauth_test.go` — `TestOAuth2_FullFlowAndPKCEInterception`, `TestOAuth2_NegativePaths`, `TestOAuth2_ExpiredAuthCode_And_ExpiredTokens`

---

## ID Token Validation

**Research:**
- `research/05-report.md` — Finding 3 (13 steps validation)
- `research/03-evidence.md` — OIDC Core Sec 3.1.3.7, RFC 8725

**Implementation:**
- `pkg/oidc/oidc.go` — `ParseAndVerifyIDToken()` dengan 6+ validasi
- `pkg/client/client.go` — `Exchange()` memanggil validasi

**Tests:**
- `tests/oauth_test.go` — `TestOIDC_IDToken_Valid`, `TestOIDC_IDToken_TamperedSignature`, `TestOIDC_IDToken_Expired`, `TestOIDC_IDToken_MismatchClaims`, `TestOIDC_MalformedJWT`, `TestOIDC_IDToken_IssuedInFuture`

---

## Refresh Token Rotation & Replay Detection

**Research:**
- `research/05-report.md` — Finding 7 (Rotation wajib)
- `research/03-evidence.md` — RFC 9700 Section 4.14

**Implementation:**
- `pkg/server/server.go` — `Refresh()` dengan rotation dan family revocation
- `pkg/client/client.go` — `RefreshTokens()`

**Tests:**
- `tests/oauth_test.go` — `TestOAuth2_RefreshTokenRotation_AndReplayDetection`, `TestOAuth2_ConcurrentRefreshReplay`

---

## Concurrency & Race Safety

**Research:**
- (Not a separate research item; covered by engineering quality gates)

**Implementation:**
- `pkg/server/server.go` — Mutex-protected maps
- `pkg/client/client.go` — Thread-safe client operations

**Tests:**
- `tests/oauth_test.go` — `TestOAuth2_ConcurrentRefreshReplay`, `TestOAuth2_ConcurrencyAndRace`

---

## Demo Execution

**Implementation:**
- `cmd/demo/main.go` — 8-step walkthrough

**Verified Behavior:**
- `engineering/03-execution-result.md` — Demo output (all steps pass)

---

## Audits

**Research Audit:**
- `research-audit/07-verdict.md` — APPROVED
- `research-audit/02-source-audit.md`
- `research-audit/03-claim-audit.md`
- `research-audit/04-contradictions.md`

**Engineering Audit:**
- `engineering-audit/06-verdict.md` — APPROVED
- `engineering-audit/02-code-audit.md`
- `engineering-audit/03-test-audit.md`
- `engineering-audit/04-docs-vs-code.md`
- `engineering-audit/05-gaps.md`

---

## Documentation

**Design & Implementation:**
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`

**Research:**
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`
