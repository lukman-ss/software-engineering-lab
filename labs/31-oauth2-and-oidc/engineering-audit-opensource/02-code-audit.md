# Code Audit

Target Lab: labs/31-oauth2-and-oidc

## Finding 1 — Scope enforcement missing in ValidateAccessToken

Location: pkg/server/server.go:268-278
Claimed Behavior: Design doc states "Resource server accepts Access Token and checks scopes." README states Access Tokens represent "delegated authorization grants to protected resources (scopes...)". Function signature includes `requiredScope string` parameter, implying scope enforcement.
Observed Implementation: `ValidateAccessToken` accepts `requiredScope` but never references it. Only existence of the token is checked; scope is ignored entirely.
Assessment: FAIL
Severity: MEDIUM
Notes: DOC_CODE_MISMATCH + IMPLEMENTATION_OVERCLAIM. The test suite does not exercise scope rejection (TestOAuth2_ConcurrencyAndRace calls ValidateAccessToken with "openid" but never asserts that a wrong scope is rejected). Scope is a core OAuth2 concept; unproven enforcement is a real gap.

## Finding 2 — Unchecked crypto/rand.Read return values

Location: pkg/client/client.go:44,47; pkg/server/server.go:153,160,164,240,244,254
Claimed Behavior: Robust random generation for codes, tokens, state, nonce.
Observed Implementation: `rand.Read(...)` return values (n, err) are silently discarded in 8 places.
Assessment: WARNING
Severity: LOW
Notes: crypto/rand.Read returns an error only under catastrophic system entropy failure (and may panic). Not a functional bug today, but error-swallowing is a code smell. Should be `if _, err := rand.Read(...); err != nil { return ... }`.

## Finding 3 — ErrIssuedInFuture defined but untested

Location: pkg/oidc/oidc.go:107, pkg/oidc/oidc.go:21
Claimed Behavior: ID tokens issued in the future (> 5 min leeway) must be rejected.
Observed Implementation: Check exists at line 107. No test exercises it.
Assessment: WARNING
Severity: LOW
Notes: MISSING_EDGE_CASE. The 5-minute leeway is generous; a token issued 6 minutes in the future would be rejected, but this path is unverified.

## Finding 4 — Design doc references non-existent test file names

Location: engineering/01-design.md:80-85
Claimed Behavior: Test files: pkce_test.go, oidc_test.go, server_test.go, race_test.go.
Observed Implementation: All tests live in a single file: tests/oauth_test.go.
Assessment: WARNING
Severity: LOW
Notes: DOC_CODE_MISMATCH (minor). Design doc file names do not match actual layout. Does not affect correctness.

## Finding 5 — Missing edge-case tests (negative paths)

Location: tests/oauth_test.go
Claimed Behavior: Tests cover happy path, failure path, edge cases, transitions, recovery, rollback, concurrency.
Observed Implementation: 10 tests total. Missing: malformed JWT (wrong part count), unauthorized client on refresh, expired refresh token, expired auth code, missing code_challenge, invalid code_challenge_method, unauthorized client/redirect on Authorize, RefreshTokens before any exchange, Refresh with ErrRefreshTokenNotFound, ValidateAccessToken scope rejection.
Assessment: WARNING
Severity: MEDIUM
Notes: MISSING_EDGE_CASE / MISSING_TEST. Happy paths and the four demonstrated attack scenarios are covered, but many negative paths are unverified. The suite is stronger than it appears but incomplete relative to the design's stated test strategy.

## Finding 6 — Short-circuit protects nil access in ExchangeCode

Location: pkg/server/server.go:133
Claimed Behavior: Safe handling of unknown codes.
Observed Implementation: `if !exists || time.Now().After(ac.ExpiresAt)` — Go short-circuit evaluation guarantees `ac.ExpiresAt` is never dereferenced when `ac` is nil. Correct.
Assessment: PASS
Severity: LOW
Notes: Verified correct. No action needed.

## Finding 7 — Concurrency safety

Location: pkg/server/server.go (all exported methods)
Claimed Behavior: Thread-safe in-memory store under concurrent load.
Observed Implementation: All state-mutating methods (RegisterClient, Authorize, ExchangeCode, Refresh, ValidateAccessToken) acquire `s.mu` via `s.mu.Lock()` / `defer s.mu.Unlock()`. No lock-free reads of shared maps.
Assessment: PASS
Severity: LOW
Notes: Verified with `go test -race ./...` — PASS, no data races detected across 20 concurrent workers.

## Finding 8 — Token family revocation logic

Location: pkg/server/server.go:206-266
Claimed Behavior: Refresh token rotation with family revocation on replay (RFC 9700 Section 4.14).
Observed Implementation: Correct ordering — family revoked BEFORE checking `meta.Revoked`, revoked flag set on consumed token, new rotated token issued in same family. Replay of consumed token triggers family-wide revocation.
Assessment: PASS
Severity: LOW
Notes: Logic verified against test expectations. Demo confirms behavior.