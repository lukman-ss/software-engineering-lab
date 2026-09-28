# Code Audit Findings

## Finding 1

Location: `pkg/pkce/pkce.go:23-72`
Claimed Behavior: RFC 7636 PKCE code verifier and challenge generation and verification with S256 (mandatory length check 43-128 chars) and plain support.
Observed Implementation: Verifier generation uses `crypto/rand` 32 bytes (43 base64url characters). `ComputeChallenge` checks verifier length `[43, 128]` and produces Base64URL-encoded SHA-256 digest for `S256` or identity for `plain`. `Verify` validates against challenge.
Assessment: PASS
Severity: LOW
Notes: Compliant with RFC 7636 and RFC 9700.

## Finding 2

Location: `pkg/oidc/oidc.go:40-116`
Claimed Behavior: RFC 7519 / OIDC Core 1.0 ID token signing and validation with standard claims (`iss`, `sub`, `aud`, `exp`, `iat`, `nonce`).
Observed Implementation: HMAC-SHA256 signature calculation and constant-time verification (`hmac.Equal`). Claims parsing validates issuer, audience, expiration time, future skew bounds, and nonce matching.
Assessment: PASS
Severity: LOW
Notes: Implementation correctly parses standard compact JWT format (`header.payload.signature`).

## Finding 3

Location: `pkg/server/server.go:63-286`
Claimed Behavior: Authorization code exchange requires PKCE verification, marks codes as used, supports OIDC ID token issuance when `openid` scope requested.
Observed Implementation: `Authorize` enforces mandatory `code_challenge` and valid method. `ExchangeCode` validates grant expiration, prevents code reuse (`ErrCodeAlreadyUsed`), validates redirect URI and client ID, verifies PKCE challenge, issues access token, generates refresh token with `FamilyID`, and issues signed ID token if scope contains `openid`.
Assessment: PASS
Severity: LOW
Notes: Correct synchronization with mutex lock covering entire state mutation.

## Finding 4

Location: `pkg/server/server.go:219-286`
Claimed Behavior: Refresh Token Rotation with token family tracking (RFC 9700 Section 4.14). Replay of a consumed refresh token revokes entire family.
Observed Implementation: When a token is refreshed, old token is marked `Revoked = true`. If a revoked token is presented, `revokedFams[meta.FamilyID]` is set to `true` and `ErrTokenReplayDetected` is returned. Subsequent attempts to use active tokens from that family are blocked.
Assessment: PASS
Severity: LOW
Notes: Fully implements single-use rotation and family revocation.

## Finding 5

Location: `pkg/server/server.go:288-321`
Claimed Behavior: Access token validation and scope enforcement.
Observed Implementation: Validates existence and expiration of access tokens under mutex lock. Helper `containsScope` checks all requested scopes against granted space-delimited scopes.
Assessment: PASS
Severity: LOW
Notes: Scope parsing handles whitespace delimiters accurately.
