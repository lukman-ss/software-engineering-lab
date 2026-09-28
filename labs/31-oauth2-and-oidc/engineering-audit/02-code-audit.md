# Code Audit: labs/31-oauth2-and-oidc

## Finding 1

Location: `pkg/pkce/pkce.go:47-61`
Claimed Behavior: PKCE RFC 7636 / RFC 9700 S256 verification and character constraint enforcement (43-128 chars).
Observed Implementation: Verifier length is bounded within `[43, 128]`. SHA-256 is computed and raw URL-safe base64 encoded.
Assessment: PASS
Severity: LOW
Notes: Correctly rejects non-compliant verifiers and invalid methods.

## Finding 2

Location: `pkg/oidc/oidc.go:64-116`
Claimed Behavior: OIDC Core 1.0 ID Token claims validation including cryptographic signature, issuer, audience, expiration, clock skew, and nonce.
Observed Implementation: Checks header/payload/signature format, verifies HMAC-SHA256 with constant-time `hmac.Equal`, checks `claims.Issuer == expectedIssuer`, `claims.Audience == expectedAudience`, `claims.Expiration <= unixNow`, future `iat` tolerance (+300s), and `nonce` match.
Assessment: PASS
Severity: LOW
Notes: Matches OIDC specifications for symmetrically-signed ID tokens.

## Finding 3

Location: `pkg/server/server.go:132-217`
Claimed Behavior: Auth Code exchange enforces single-use redemption, client binding, redirect URI verification, PKCE verification, and conditional ID token issuance.
Observed Implementation: Single lock `s.mu.Lock()` guards map access. Re-used authorization codes return `ErrCodeAlreadyUsed`. Verifier is checked with `pkce.Verify`. ID Token is only added if `openid` scope is present.
Assessment: PASS
Severity: LOW
Notes: Implemented without deadlocks or race windows.

## Finding 4

Location: `pkg/server/server.go:219-286`
Claimed Behavior: Refresh Token Rotation with token family tracking and revocation on reuse (RFC 9700 §4.14).
Observed Implementation: Each lineage has `FamilyID`. Consumed refresh tokens have `meta.Revoked = true`. If a revoked token is presented, `s.revokedFams[meta.FamilyID] = true` and `ErrTokenReplayDetected` is returned. Any subsequent attempt using any token in `s.revokedFams` returns `ErrTokenReplayDetected`.
Assessment: PASS
Severity: LOW
Notes: Accurately implements RFC 9700 §4.14 family invalidation semantics.

## Finding 5

Location: `pkg/server/server.go:63-73`
Claimed Behavior: Concurrency safe state mutations across all endpoint handlers.
Observed Implementation: All exported stateful operations (`RegisterClient`, `Authorize`, `ExchangeCode`, `Refresh`, `ValidateAccessToken`) acquire `s.mu.Lock()`.
Assessment: PASS
Severity: LOW
Notes: No unsynchronized memory reads or writes detected.
