# Test Audit

## Coverage Analysis

### 1. Happy Path Coverage
- PKCE S256 generation and verification (`TestPKCE_S256_Valid`).
- Plain PKCE generation and verification (`TestPKCE_Plain_Method`).
- OIDC ID token creation and claim parsing (`TestOIDC_IDToken_Valid`).
- Full OAuth 2.0 + OIDC code exchange flow (`TestOAuth2_FullFlowAndPKCEInterception`).
- Refresh token rotation lifecycle (`TestOAuth2_RefreshTokenRotation_AndReplayDetection`).

### 2. Failure & Negative Path Coverage
- PKCE invalid method (`TestPKCE_InvalidMethod`).
- PKCE mismatch rejection (`TestPKCE_Mismatch`).
- Tampered JWT signature rejection (`TestOIDC_IDToken_TamperedSignature`).
- Expired ID token rejection (`TestOIDC_IDToken_Expired`).
- ID token issuer/audience/nonce mismatch rejection (`TestOIDC_IDToken_MismatchClaims`).
- Malformed JWT inputs (`TestOIDC_MalformedJWT`).
- OAuth authorization code interception prevention (`TestOAuth2_FullFlowAndPKCEInterception`).
- Authorization code reuse / replay prevention (`TestOAuth2_FullFlowAndPKCEInterception`).
- Refresh token reuse triggering entire token family revocation (`TestOAuth2_RefreshTokenRotation_AndReplayDetection`).
- Authorization failures (unregistered client, wrong redirect URI, missing challenge, invalid method) (`TestOAuth2_NegativePaths`).
- Token exchange failures (invalid code, mismatched client, mismatched redirect URI) (`TestOAuth2_NegativePaths`).
- Access token validation failures (scope mismatch, non-existent token) (`TestOAuth2_NegativePaths`).
- Refresh failures (non-existent token, wrong client) (`TestOAuth2_NegativePaths`).

### 3. Concurrency & Race Detector Coverage
- Concurrent multi-worker authorization, exchange, validation, and refresh test (`TestOAuth2_ConcurrencyAndRace`) runs 20 parallel routines without data races.

## Actual Test Execution Result

```text
=== RUN   TestPKCE_S256_Valid
--- PASS: TestPKCE_S256_Valid (0.00s)
=== RUN   TestPKCE_InvalidMethod
--- PASS: TestPKCE_InvalidMethod (0.00s)
=== RUN   TestPKCE_Mismatch
--- PASS: TestPKCE_Mismatch (0.00s)
=== RUN   TestOIDC_IDToken_Valid
--- PASS: TestOIDC_IDToken_Valid (0.00s)
=== RUN   TestOIDC_IDToken_TamperedSignature
--- PASS: TestOIDC_IDToken_TamperedSignature (0.00s)
=== RUN   TestOIDC_IDToken_Expired
--- PASS: TestOIDC_IDToken_Expired (0.00s)
=== RUN   TestOIDC_IDToken_MismatchClaims
--- PASS: TestOIDC_IDToken_MismatchClaims (0.00s)
=== RUN   TestOAuth2_FullFlowAndPKCEInterception
--- PASS: TestOAuth2_FullFlowAndPKCEInterception (0.00s)
=== RUN   TestOAuth2_RefreshTokenRotation_AndReplayDetection
--- PASS: TestOAuth2_RefreshTokenRotation_AndReplayDetection (0.00s)
=== RUN   TestPKCE_Plain_Method
--- PASS: TestPKCE_Plain_Method (0.00s)
=== RUN   TestOIDC_MalformedJWT
--- PASS: TestOIDC_MalformedJWT (0.00s)
=== RUN   TestOAuth2_NegativePaths
--- PASS: TestOAuth2_NegativePaths (0.00s)
=== RUN   TestOAuth2_ConcurrencyAndRace
--- PASS: TestOAuth2_ConcurrencyAndRace (0.00s)
PASS
ok  	labs/31-oauth2-and-oidc/tests	1.116s
```

All 13 test suites pass cleanly under race detection.
