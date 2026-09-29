# Test Audit

## Test Suite Overview

Test file: `tests/oauth_test.go`
Execution Command: `go test -v -race -count=1 ./tests/...`
Result: 17 passed, 0 failed, 0 race conditions detected.

## Test Coverage Analysis

### 1. Happy Path
- `TestPKCE_S256_Valid`: Verifies PKCE S256 generation and verification.
- `TestPKCE_Plain_Method`: Verifies PKCE plain method compatibility.
- `TestOIDC_IDToken_Valid`: Verifies ID token signing and verification with all standard claims.
- `TestOAuth2_FullFlowAndPKCEInterception`: Verifies full end-to-end OAuth + OIDC flow from code generation to token exchange and access token validation.

### 2. Failure & Negative Paths
- `TestPKCE_InvalidMethod`: Verifies rejection of unsupported challenge methods.
- `TestPKCE_Mismatch`: Verifies rejection when challenge does not match verifier.
- `TestPKCE_VerifierLength_Bounds`: Verifies bounds checking for verifier length < 43 and > 128 characters.
- `TestOIDC_IDToken_TamperedSignature`: Verifies rejection of tampered HMAC signatures.
- `TestOIDC_IDToken_Expired`: Verifies rejection of expired ID tokens.
- `TestOIDC_IDToken_MismatchClaims`: Verifies rejection on mismatch of issuer, audience, and nonce claims.
- `TestOIDC_IDToken_IssuedInFuture`: Verifies rejection of tokens with `iat` in the future.
- `TestOIDC_MalformedJWT`: Verifies rejection of malformed tokens (invalid split count, non-base64).
- `TestOAuth2_NegativePaths`: Exhaustively tests invalid client ID, invalid redirect URI, missing challenge, invalid challenge method, non-existent code, client ID/redirect URI mismatch during exchange, token scope rejection, and unknown refresh token.
- `TestOAuth2_ExpiredAuthCode_And_ExpiredTokens`: Verifies rejection when authorization code has expired.

### 3. Attack Scenarios & Security Mitigations
- `TestOAuth2_FullFlowAndPKCEInterception`: Interception attack simulation where code is exchanged with an illegitimate verifier. Asserts PKCE failure and single-use code redemption protection.
- `TestOAuth2_RefreshTokenRotation_AndReplayDetection`: Replay attack simulation where a stolen refresh token is reused, asserting detection and subsequent family revocation of the active session.

### 4. Concurrency & Race Safety
- `TestOAuth2_ConcurrentRefreshReplay`: 10 concurrent goroutines racing to reuse a refresh token; asserts that at least 9 requests fail due to rotation and family revocation.
- `TestOAuth2_ConcurrencyAndRace`: 20 concurrent goroutines performing concurrent authorize, code exchange, token validation, and token refresh under race detector without data races.

## Execution Output

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
=== RUN   TestOIDC_IDToken_IssuedInFuture
--- PASS: TestOIDC_IDToken_IssuedInFuture (0.00s)
=== RUN   TestPKCE_VerifierLength_Bounds
--- PASS: TestPKCE_VerifierLength_Bounds (0.00s)
=== RUN   TestOAuth2_ExpiredAuthCode_And_ExpiredTokens
--- PASS: TestOAuth2_ExpiredAuthCode_And_ExpiredTokens (0.00s)
=== RUN   TestOAuth2_ConcurrentRefreshReplay
--- PASS: TestOAuth2_ConcurrentRefreshReplay (0.00s)
=== RUN   TestOAuth2_ConcurrencyAndRace
--- PASS: TestOAuth2_ConcurrencyAndRace (0.00s)
PASS
ok  	labs/31-oauth2-and-oidc/tests	1.119s
```
