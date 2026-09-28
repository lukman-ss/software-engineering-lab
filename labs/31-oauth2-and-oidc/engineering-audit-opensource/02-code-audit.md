# Code Audit — labs/31-oauth2-and-oidc

## Finding 1

Location: `pkg/pkce/pkce.go:23-72`
Claimed Behavior: PKCE S256/plain generation and verification per RFC 7636; mandatory challenge.
Observed Implementation: `GeneratePKCEPair` rejects non-S256/plain; 32 random bytes → 43-char base64url verifier; `ComputeChallenge` enforces 43–128 length, SHA256+RawURLEncoding for S256, identity for plain; `Verify` recomputes and compares.
Assessment: PASS
Severity: LOW
Notes: String `!=` compare in `Verify` is not constant-time; acceptable here since challenge is not a long-lived secret and HMAC path uses `hmac.Equal`. No correctness impact.

## Finding 2

Location: `pkg/oidc/oidc.go:40-116`
Claimed Behavior: ID Token HS256 signing + strict claims validation (iss/aud/exp/iat/nonce).
Observed Implementation: `SignIDToken` JSON→base64url→HMAC-SHA256; `ParseAndVerifyIDToken` enforces 3-part JWT, base64 decode, `hmac.Equal` signature check, JSON unmarshal, iss/aud equality, `exp <= now` reject, `iat > now+300` reject, nonce check when expected non-empty.
Assessment: PASS
Severity: LOW
Notes: `alg` header not explicitly validated (always HS256 assumed). Safe in this closed lab because any header mutation invalidates HMAC; no `none` bypass possible. Documented HS256-for-RS256 substitution in engineering notes is correctly scoped.

## Finding 3

Location: `pkg/server/server.go:92-130` (`Authorize`)
Claimed Behavior: Auth request validates client, redirect, mandatory PKCE.
Observed Implementation: Mutex-protected; `Clients[clientID]` + redirect exact match else `ErrUnauthorizedClient`; empty challenge → `ErrInvalidRequest`; method must be S256/plain; 16-byte random code, 5-min expiry, `Used=false`.
Assessment: PASS
Severity: LOW
Notes: Correct state transition (create unconsumed code). Failure handling and error propagation clean.

## Finding 4

Location: `pkg/server/server.go:132-217` (`ExchangeCode`)
Claimed Behavior: One-time code, PKCE proof-of-possession, token issuance.
Observed Implementation: Checks existence/expiry → `ErrInvalidGrant`; `Used` → `ErrCodeAlreadyUsed`; client/redirect match → `ErrUnauthorizedClient`; `pkce.Verify` failure → `ErrInvalidPKCE`; sets `Used=true` before issuance; issues 24-byte access + refresh tokens, family ID, 1h/30d expiry; mints ID Token only when `openid` scope present with correct iss/aud/exp/iat/auth_time/nonce.
Assessment: PASS
Severity: LOW
Notes: Mark-used-before-issue prevents double-redeem on partial failure. No timeout/recovery needed beyond expiry. Cleanup: expired codes/tokens never purged (unbounded map growth) — acceptable for in-memory lab, noted as limitation.

## Finding 5

Location: `pkg/server/server.go:219-286` (`Refresh`)
Claimed Behavior: Single-use rotation; replay revokes entire family (RFC 9700 §4.14).
Observed Implementation: Not-found → `ErrRefreshTokenNotFound`; `revokedFams[family]` → `ErrTokenReplayDetected`; `meta.Revoked` → sets `revokedFams[family]=true` + `ErrTokenReplayDetected`; then client-match, expiry checks; marks old revoked, issues new token in same family + new access token.
Assessment: PASS
Severity: LOW
Notes: Check order (family-revoked before client-match) leaks no secret. Behavior proven by test + demo. Concurrency safe under single `sync.Mutex`.

## Finding 6

Location: `pkg/server/server.go:288-306`, `308-349` (`ValidateAccessToken`, scope helpers)
Claimed Behavior: Bearer token + scope enforcement; separation of Access Token (authorization) vs ID Token (authentication).
Observed Implementation: Existence + expiry check; `containsScope` requires all requested scopes present; custom `splitSpaces` handles space/tab.
Assessment: PASS
Severity: LOW
Notes: No complexity issue. Error strings generic, no secret leak.

## Finding 7

Location: `pkg/server/server.go:63-90` + all methods (concurrency)
Claimed Behavior: Thread-safe in-memory AS.
Observed Implementation: Single `sync.Mutex` guards `Clients/authCodes/tokens/refreshMeta/revokedFams` on every read/write path (`RegisterClient/Authorize/ExchangeCode/Refresh/ValidateAccessToken`).
Assessment: PASS
Severity: LOW
Notes: `go test -race` passed; 20-worker concurrency test passed. No race observed. Coarse lock is correct for lab scale.

## Finding 8

Location: `pkg/client/client.go:36-86`
Claimed Behavior: Correct flow parameters, verifier storage, ID Token validation.
Observed Implementation: `BuildAuthorizationRequest` generates S256 pair, stores verifier, random state+nonce, returns challenge; `Exchange` calls `ExchangeCode` then verifies ID Token against `Server.Issuer/ClientID/Nonce`; `RefreshTokens` updates stored tokens.
Assessment: PASS
Severity: LOW
Notes: Client holds `SigningKey` symmetric to server — matches documented lab simplification. No state/nonce mismatch handling beyond server binding; adequate.

## Finding 9

Location: `cmd/demo/main.go:1-112`
Claimed Behavior: End-to-end walkthrough of legitimate flows + attack defenses.
Observed Implementation: Real calls to `Authorize/Exchange/ValidateAccessToken/Refresh`; interception with wrong verifier must fail; rotation then replay must fail; post-revocation active token must fail; `os.Exit(1)` on unexpected success.
Assessment: PASS
Severity: LOW
Notes: Demo output verified real by re-execution (see 03-test-audit.md). No fabricated output, no benchmark claims.
