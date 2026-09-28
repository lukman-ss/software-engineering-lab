# Code Audit

Target Lab: `labs/31-oauth2-and-oidc`

## Finding 1: Standard Library Zero-Dependency Implementation

Location: `go.mod`, `pkg/pkce/pkce.go`, `pkg/oidc/oidc.go`, `pkg/server/server.go`, `pkg/client/client.go`
Claimed Behavior: Pure Go standard library implementation without third-party external dependencies.
Observed Implementation: `go.mod` specifies `module labs/31-oauth2-and-oidc` with Go 1.22.0 and no external `require` directives. Imports use standard packages (`crypto/hmac`, `crypto/sha256`, `crypto/rand`, `encoding/base64`, `encoding/json`, `sync`, `time`, `errors`, `fmt`).
Assessment: PASS
Severity: LOW
Notes: Clean, minimal, zero-dependency design.

## Finding 2: PKCE Challenge Generation & Verification (RFC 7636)

Location: `pkg/pkce/pkce.go`
Claimed Behavior: Generates base64url-encoded code verifiers and S256 code challenges; rejects mismatched verifiers and invalid methods.
Observed Implementation:
- `GeneratePKCEPair` validates method (`S256` or `plain`), reads 32 random bytes, encodes with `base64.RawURLEncoding` (yielding 43 chars), and computes challenge.
- `ComputeChallenge` checks verifier length (43..128) and calculates standard SHA-256 raw URL base64 digest for `S256`.
- `Verify` re-computes challenge and compares string equality.
Assessment: PASS
Severity: LOW
Notes: Correct standard implementation of RFC 7636 PKCE validation.

## Finding 3: OIDC ID Token JWT Generation and Claim Validation

Location: `pkg/oidc/oidc.go`
Claimed Behavior: Signs ID Tokens using HMAC-SHA256 (HS256) and verifies claims (`iss`, `aud`, `exp`, `nonce`, future `iat`).
Observed Implementation:
- `SignIDToken` constructs standard 3-part JWT (`header.claims.signature`) with `base64.RawURLEncoding`.
- `ParseAndVerifyIDToken` verifies HMAC-SHA256 signature using `hmac.Equal` (constant-time protection), validates `iss` equality, `aud` equality, expiration (`claims.Expiration <= unixNow`), future `iat` check (`unixNow+300`), and optional `nonce` equality.
Assessment: PASS
Severity: LOW
Notes: Secure signature comparison using `hmac.Equal` avoids timing attack vectors.

## Finding 4: Authorization Server State and Concurrency Safety

Location: `pkg/server/server.go`
Claimed Behavior: State transitions for Auth Codes, Access Tokens, and Refresh Token family lineages are protected by mutex locking.
Observed Implementation:
- `AuthorizationServer` uses `s.mu.Lock()` and `defer s.mu.Unlock()` across all exported methods: `RegisterClient`, `Authorize`, `ExchangeCode`, `Refresh`, and `ValidateAccessToken`.
- `AuthCode` validation marks `ac.Used = true` while under mutex lock.
- `Refresh` checks `revokedFams`, marks token `Revoked = true`, and updates maps while under mutex lock.
Assessment: PASS
Severity: LOW
Notes: No race conditions found under Go race detector.

## Finding 5: Refresh Token Rotation & Family Revocation (RFC 9700)

Location: `pkg/server/server.go:219-286`
Claimed Behavior: Single-use refresh token rotation issuing new tokens in the same family, with replay detection revoking the entire token family.
Observed Implementation:
- When a refresh token is presented, `Refresh` checks if `revokedFams[meta.FamilyID]` is true, returning `ErrTokenReplayDetected`.
- If `meta.Revoked` is true (indicating replay of a consumed refresh token), `s.revokedFams[meta.FamilyID] = true` is set, blocking all future tokens in that family.
- On valid refresh, original token is marked `meta.Revoked = true`, and a new refresh token is stored with identical `FamilyID`.
Assessment: PASS
Severity: LOW
Notes: Implementation matches RFC 9700 Section 4.14 specs.
