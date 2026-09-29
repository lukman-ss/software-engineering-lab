# Code Audit

## Finding 1

Location: `pkg/pkce/pkce.go:28-34`
Claimed Behavior: Generates random cryptographic verifier compliant with RFC 7636 (length 43 to 128 characters).
Observed Implementation: Uses `crypto/rand` reading 32 bytes and encodes using `base64.RawURLEncoding`, producing exactly 43 characters.
Assessment: PASS
Severity: LOW
Notes: Complies with RFC 7636 min length (43).

## Finding 2

Location: `pkg/pkce/pkce.go:48-61`
Claimed Behavior: Supports RFC 7636 challenge computation with length verification and methods (`S256`, `plain`).
Observed Implementation: Rejects verifiers < 43 or > 128 characters with `ErrInvalidVerifierLength`. Correctly computes SHA-256 base64url-unpadded hash for `S256` and verbatim value for `plain`.
Assessment: PASS
Severity: LOW
Notes: Correct standard implementation.

## Finding 3

Location: `pkg/oidc/oidc.go:64-115`
Claimed Behavior: Parse and cryptographically verify ID Token JWT signature, standard claims (`iss`, `aud`, `exp`, `iat`, `nonce`).
Observed Implementation: Uses constant-time comparison `hmac.Equal` for HS256 signature verification. Validates `exp` against current time, checks clock skew threshold on `iat` (`iat > unixNow + 300`), and enforces exact matching for `iss`, `aud`, and `nonce`.
Assessment: PASS
Severity: LOW
Notes: Clean cryptographic verification using stdlib primitives.

## Finding 4

Location: `pkg/server/server.go:63-72`, `pkg/server/server.go:93-130`, `pkg/server/server.go:133-217`, `pkg/server/server.go:220-286`
Claimed Behavior: Thread-safe authorization server handling state, auth code redemption, and refresh token rotation with token family revocation upon reuse.
Observed Implementation: All state mutating and lookup methods (`Authorize`, `ExchangeCode`, `Refresh`, `ValidateAccessToken`) synchronize using mutex `s.mu.Lock()`. Single-use flag `ac.Used` prevents authorization code replay. Refresh token rotation tracks `FamilyID` in `s.refreshMeta` and invalidates `s.revokedFams[meta.FamilyID]` on detecting a reused token.
Assessment: PASS
Severity: LOW
Notes: Concurrency safety verified with `-race` tests.
