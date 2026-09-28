# Code Audit

## Finding 1

Location: `pkg/pkce/pkce.go:23-45`
Claimed Behavior: Generate PKCE verifier and S256/plain challenge per RFC 7636.
Observed Implementation: Verifiers are 32 cryptographically random bytes base64raw-url encoded (43 chars). `ComputeChallenge` correctly handles S256 (SHA-256 base64url) and plain.
Assessment: PASS
Severity: LOW
Notes: Complies fully with RFC 7636 specs.

## Finding 2

Location: `pkg/oidc/oidc.go:40-116`
Claimed Behavior: HMAC-SHA256 JWT ID token signing and verification with strict claims validation.
Observed Implementation: Signs header and claims with HMAC-SHA256. Validates signature via constant-time `hmac.Equal`, checks issuer, audience, expiration, future issued-at (+300s skew), and nonce.
Assessment: PASS
Severity: LOW
Notes: Timing attack safe using `hmac.Equal`.

## Finding 3

Location: `pkg/server/server.go:132-217`
Claimed Behavior: Single-use Auth Code exchange with PKCE validation and token generation.
Observed Implementation: Server locks state, checks code expiry, checks `ac.Used`, validates PKCE verifier, marks `ac.Used = true`, and issues access, refresh, and OIDC ID tokens.
Assessment: PASS
Severity: LOW
Notes: Replay of auth code is blocked.

## Finding 4

Location: `pkg/server/server.go:219-286`
Claimed Behavior: Refresh token rotation with family revocation on replay detection.
Observed Implementation: Checks if token exists, if family is revoked, or if specific token was already revoked. On reuse of revoked token, marks `revokedFams[familyID] = true` and rejects request. Active tokens in the same family subsequently fail.
Assessment: PASS
Severity: LOW
Notes: Fully aligns with RFC 9700 Section 4.14.

## Finding 5

Location: `pkg/server/server.go:63-349`
Claimed Behavior: Concurrency safety across server operations.
Observed Implementation: All stateful operations (`Authorize`, `ExchangeCode`, `Refresh`, `ValidateAccessToken`, `RegisterClient`) lock `s.mu`.
Assessment: PASS
Severity: LOW
Notes: No race conditions detected under `go test -race`.
