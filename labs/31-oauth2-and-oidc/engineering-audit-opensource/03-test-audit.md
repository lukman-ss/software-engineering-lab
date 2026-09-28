# Test Audit

## Coverage summary

Tests in `tests/oauth_test.go`: 12 functions.

### Happy path
- `TestOAuth2_FullFlowAndPKCEInterception` — full Auth Code + PKCE + OIDC + token issuance (PASS)
- `TestPKCE_S256_Valid` — valid S256 pair (PASS)
- `TestPKCE_Plain_Method` — plain method (PASS)
- `TestOIDC_IDToken_Valid` — round-trip sign/verify (PASS)
- `TestOAuth2_RefreshTokenRotation_AndReplayDetection` — normal + rotation + replay (PASS)
- `TestOAuth2_ConcurrencyAndRace` — 20 goroutines full flow (PASS) under `-race`

### Failure path
- `TestPKCE_InvalidMethod` — bad method (PASS)
- `TestPKCE_Mismatch` — wrong verifier (PASS)
- `TestOAuth2_NegativePaths` — bad client, bad redirect, empty challenge, invalid method,
  nonexistent code, mismatched client/redirect, wrong scope, invalid token, wrong-client refresh,
  nonexistent refresh (PASS)
- `TestOIDC_IDToken_TamperedSignature` (PASS)
- `TestOIDC_IDToken_Expired` (PASS)
- `TestOIDC_IDToken_MismatchClaims` — issuer/audience/nonce mismatch (PASS)
- `TestOIDC_MalformedJWT` — 2-part token, non-base64 (PASS)

### Attack path
- Interception (wrong verifier → PKCE failure) — covered in FullFlow test (PASS)
- Replay of consumed refresh token → family revocation (PASS)

### Execution result

```
$ go test -v ./...
=== RUN   TestPKCE_S256_Valid --- PASS
=== RUN   TestPKCE_InvalidMethod --- PASS
=== RUN   TestPKCE_Mismatch --- PASS
=== RUN   TestOIDC_IDToken_Valid --- PASS
=== RUN   TestOIDC_IDToken_TamperedSignature --- PASS
=== RUN   TestOIDC_IDToken_Expired --- PASS
=== RUN   TestOIDC_IDToken_MismatchClaims --- PASS
=== RUN   TestOAuth2_FullFlowAndPKCEInterception --- PASS
=== RUN   TestOAuth2_RefreshTokenRotation_AndReplayDetection --- PASS
=== RUN   TestPKCE_Plain_Method --- PASS
=== RUN   TestOIDC_MalformedJWT --- PASS
=== RUN   TestOAuth2_NegativePaths --- PASS
=== RUN   TestOAuth2_ConcurrencyAndRace --- PASS
PASS
ok labs/31-oauth2-and-oidc/tests

$ go test -race ./...
ok labs/31-oauth2-and-oidc/tests
```

## Gaps (MISSING_TEST)

### G1 — PKCE boundary lengths
Location: `pkg/pkce/pkce.go:48-49`
Claimed: verifier 43–128 chars.
Tests: only checks `len(verifier) < 43 || > 128` via a happy path at exactly 43. No test for 128-char upper bound or 42/129 rejection.
Severity: LOW

### G2 — Auth code expiry
Location: `pkg/server/server.go:137`
Claimed: expired codes rejected with `ErrInvalidGrant`.
Tests: no test forces `time.Now().After(ac.ExpiresAt)`. Only negative path for nonexistent code.
Severity: LOW

### G3 — Refresh token expiry
Location: `pkg/server/server.go:244`
Claimed: expired refresh returns `ErrInvalidGrant`.
Tests: never exercised.
Severity: LOW

### G4 — iat skew boundary (300s)
Location: `pkg/oidc/oidc.go:107`
Claimed: token issued within 5min future accepted.
Tests: no edge test at `+299s`/`+301s`.
Severity: LOW

### G5 — PKCE `plain` downgrade
Location: `pkg/pkce/pkce.go:57`, `pkg/server/server.go:105`
Claimed (README): S256 emphasis.
Tests: `plain` is exercised as a *happy* path, never as a *downgrade-risk* negative case (interception with plain is moot since challenge==verifier, but the server still accepting plain is the gap).
Severity: LOW

## Assessment

Test suite is solid, covers happy/failure/attack paths and concurrency. All 12 tests pass with and without `-race`.
