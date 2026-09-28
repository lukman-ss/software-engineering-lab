# Code Audit

Target Lab: labs/31-oauth2-and-oidc

## Finding 1

Location: `pkg/pkce/pkce.go:23-71`
Claimed Behavior: RFC 7636 PKCE pair generation and verification using S256 and plain methods.
Observed Implementation: Verifier generated via `crypto/rand` (32 bytes base64url encoded -> 43 characters). Challenge correctly hashes with SHA256 and encodes to base64url without padding. Method validation and length checks (43 to 128 characters) enforced.
Assessment: PASS
Severity: LOW
Notes: Compliant with RFC 7636.

## Finding 2

Location: `pkg/oidc/oidc.go:40-116`
Claimed Behavior: ID Token signing and claims validation (iss, aud, exp, iat, nonce) with HMAC-SHA256.
Observed Implementation: Standard HS256 JWT generation with base64url encoding. Verification parses parts, recomputes HMAC signature, uses `hmac.Equal` to prevent timing attacks, and validates issuer, audience, expiration, clock skew for iat, and nonce.
Assessment: PASS
Severity: LOW
Notes: Correctly handles tamper detection and claim mismatches.

## Finding 3

Location: `pkg/server/server.go:63-74, 92-217`
Claimed Behavior: Authorization Server enforces client registration, mandatory PKCE, one-time auth code usage, and scopes.
Observed Implementation: `s.mu.Lock()` protects all server state operations. Auth codes expire after 5 minutes and cannot be reused (`ac.Used` check). Code exchange verifies PKCE challenge with supplied verifier before minting access, refresh, and ID tokens.
Assessment: PASS
Severity: LOW
Notes: Thread-safe in-memory implementation.

## Finding 4

Location: `pkg/server/server.go:219-286`
Claimed Behavior: Refresh Token Rotation with token family tracking and reuse detection (RFC 9700 Section 4.14).
Observed Implementation: Each initial refresh token is tagged with a unique `FamilyID`. When a refresh token is presented:
1. Checks if `s.revokedFams[meta.FamilyID]` is true. If so, rejects.
2. Checks if `meta.Revoked` is true (reuse of consumed token). If so, marks `s.revokedFams[meta.FamilyID] = true` and rejects.
3. If valid, marks `meta.Revoked = true`, issues a new refresh token sharing the same `FamilyID`, and issues a new access token.
Assessment: PASS
Severity: LOW
Notes: Properly models family revocation upon replay attack detection.

## Finding 5

Location: `pkg/server/server.go:288-322`
Claimed Behavior: Resource server token validation and space-delimited scope checks.
Observed Implementation: Validates existence, expiry, and required scope subsets using custom space-tokenization helper.
Assessment: PASS
Severity: LOW
Notes: Scope parsing handles whitespace and multiple scopes properly.
