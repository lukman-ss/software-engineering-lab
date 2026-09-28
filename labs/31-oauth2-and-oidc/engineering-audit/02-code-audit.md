# Code Audit

## Finding 1

Location: `pkg/pkce/pkce.go:47-72`
Claimed Behavior: PKCE S256 challenge generation and constant-time/safe verification with strict length validation (43..128 characters).
Observed Implementation: Correctly validates verifier character length range, generates SHA256 digest, base64-URL-encodes without padding (`base64.RawURLEncoding`), and verifies matches.
Assessment: PASS
Severity: LOW
Notes: Compliant with RFC 7636.

## Finding 2

Location: `pkg/oidc/oidc.go:40-116`
Claimed Behavior: OIDC ID Token signing and claim verification (`iss`, `sub`, `aud`, `exp`, `nonce`) using HMAC-SHA256 (`HS256`).
Observed Implementation: Properly signs JWT formatted strings (`header.claims.signature`), uses `hmac.Equal` for timing-attack-safe signature verification, verifies issuer, audience, expiration time, future skew bounds, and nonce integrity.
Assessment: PASS
Severity: LOW
Notes: Claims validation adheres to OpenID Connect Core 1.0 specifications.

## Finding 3

Location: `pkg/server/server.go:88-204`
Claimed Behavior: Authorization code issuance and single-use redemption with PKCE enforcement and OIDC ID Token generation when `openid` scope is present.
Observed Implementation: Authorization codes are stored with metadata, checked for expiry (5 mins), single-use flag checked and flipped atomically under mutex lock `s.mu`, PKCE verified before token issuance, and `openid` scope triggers ID Token construction and signing.
Assessment: PASS
Severity: LOW
Notes: State transitions are correctly locked and guarded against replay.

## Finding 4

Location: `pkg/server/server.go:206-266`
Claimed Behavior: Refresh Token Rotation with token family tracking and automatic family-wide revocation upon detecting replay of consumed tokens (RFC 9700 §4.14).
Observed Implementation: Each refresh token is linked to a `FamilyID`. When a refresh token marked `Revoked: true` is presented for refresh, `s.revokedFams[meta.FamilyID] = true` is set, blocking any subsequent usage of active tokens from the same lineage family.
Assessment: PASS
Severity: LOW
Notes: Replay detection logic conforms to modern OAuth 2.0 Security Best Current Practice.

## Finding 5

Location: `pkg/server/server.go:57-67`, `pkg/server/server.go:268-278`
Claimed Behavior: Concurrency safety across all state mutations (authorization codes, tokens, refresh metadata, family revocation lists).
Observed Implementation: All server endpoints (`RegisterClient`, `Authorize`, `ExchangeCode`, `Refresh`, `ValidateAccessToken`) synchronize state access with `s.mu.Lock()` and `s.mu.Unlock()`.
Assessment: PASS
Severity: LOW
Notes: Thread-safe in-memory authorization server implementation.
