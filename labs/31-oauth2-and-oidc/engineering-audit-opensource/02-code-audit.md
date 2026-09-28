# Code Audit

## Finding 1 — PKCE verifier/challenge correctness (S256)

Location: `pkg/pkce/pkce.go:23-71`
Claimed Behavior: Generates a 43–128 char verifier; S256 base64url challenge = BASE64URL(SHA256(verifier)); rejects mismatch.
Observed Implementation: `crypto/rand` 32 bytes → 43-char base64url verifier; SHA256 then RawURLEncoding for challenge; length check 43–128; `Verify` recomputes and compares exact strings.
Assessment: PASS
Severity: LOW
Notes: Correct per RFC 7636 §4.1. Only concern: 32-byte input → 43-char verifier is exactly at the lower bound (acceptable).

## Finding 2 — PKCE `plain` method accepted

Location: `pkg/pkce/pkce.go:24`, `pkg/server/server.go:105`
Claimed Behavior: README states S256 PKCE. Server accepts both `plain` and `S256`.
Observed Implementation: `GeneratePKCEPair("plain")` returns challenge == verifier. Server `Authorize` accepts `plain`.
Assessment: WARNING
Severity: LOW
Notes: RFC 9700 & RFC 8252 deprecate `plain` for native/public clients; S256 is strongly recommended. Server allowing `plain` is a downgrade surface. By-design in the package but contradicts README emphasis on S256.

## Finding 3 — ID Token HMAC-SHA256 signing & verification

Location: `pkg/oidc/oidc.go:40-116`
Claimed Behavior: HS256 JWT signed with shared secret; verification validates signature (hmac.Equal constant-time), iss, aud, exp, iat skew (300s), nonce.
Observed Implementation: Header `{"alg":"HS256","typ":"JWT"}`; HMAC-SHA256 over `b64(header).b64(claims)`; constant-time sig compare; claims decoded & checked.
Assessment: PASS
Severity: LOW
Notes: OIDC production typically uses RS256/JWKS. HS256 acceptable for a reference impl w/ shared signing key, but should be documented as dev-only. `aud` modeled as single string (OIDC allows array) — scoped simplification. No `jti`/`nbf`/`azp` validation.

## Finding 4 — Authorization code single-use & expiry

Location: `pkg/server/server.go:132-217`
Claimed Behavior: Code is `Used=true` after exchange; second exchange returns `ErrCodeAlreadyUsed`; expired codes rejected.
Observed Implementation: `ac.Used` checked before PKCE verification; set true before issuing tokens; expiry checked via `time.Now().After(ac.ExpiresAt)`.
Assessment: PASS
Severity: LOW
Notes: Correct ordering — expiry checked before `Used` flag. Attacker replay of stale-but-unexpired code is blocked.

## Finding 5 — Refresh Token Rotation with family revocation

Location: `pkg/server/server.go:219-286`
Claimed Behavior: Consumed refresh token marked `Revoked=true`; new token issued in same `FamilyID`; replay of consumed token sets `revokedFams[familyID]=true` and returns `ErrTokenReplayDetected`; any later use of family token returns family-revoked error.
Observed Implementation: Matches claimed flow exactly. Family revoked check precedes per-token `Revoked` check, so all family members are invalidated atomically.
Assessment: PASS
Severity: LOW
Notes: Correct per RFC 9700 §4.14 reuse detection & family revocation.

## Finding 6 — Access token scope validation

Location: `pkg/server/server.go:288-321`
Claimed Behavior: `ValidateAccessToken` verifies existence, expiry, and that required scope subset of granted.
Observed Implementation: `containsScope` splits on whitespace; requires all required scopes present.
Assessment: PASS
Severity: LOW
Notes: Correct subset semantics for space-delimited scopes.

## Finding 7 — Concurrency safety

Location: `pkg/server/server.go:63-84`, all handlers use `s.mu sync.Mutex` + Lock/Unlock
Claimed Behavior: Thread-safe under `-race`.
Observed Implementation: Every public method (`Authorize`, `ExchangeCode`, `Refresh`, `ValidateAccessToken`, `RegisterClient`) holds `s.mu`. Token strings generated with `crypto/rand` (non-blocking, no shared state).
Assessment: PASS (under -race)
Severity: LOW
Notes: Demo Step 8 relies on the rotated token carrying the same FamilyID so family revocation propagates — verified by race test + demo.

## Finding 8 — Error wrapping uses `%v` (not unwrappable)

Location: `pkg/server/server.go:150`, `pkg/oidc/oidc.go:95,99,112`
Claimed Behavior: Errors carry sentinel via `%w`.
Observed Implementation: e.g. `fmt.Errorf("%w: %v", ErrInvalidPKCE, err)` — sentinel wraps, but inner error uses `%v`.
Assessment: WARNING
Severity: LOW
Notes: Inner cause is not unwrappable. Sentinel errors ARE chainable; only the nested detail is lost. Tests check via string `Contains`, so they pass. Minor hygiene issue.

## Finding 9 — Demo produces real output

Location: `cmd/demo/main.go`
Claimed Behavior: Walkthrough of legitimate flows + attack defenses.
Observed Implementation: Demo exits non-zero on failure, prints actual token substrings (not hardcoded placeholders — tokens are random per run but structure is deterministic).
Assessment: PASS
Severity: LOW
Notes: See `02-code-audit.md#step8` capture. Output is genuinely generated, not fabricated.

## Finding 10 — Client nonce/state stored unguarded

Location: `pkg/client/client.go:19-24`
Claimed Behavior: Single-use client; holds transient `Verifier`, `State`, `Nonce`.
Observed Implementation: Plain fields, no mutex. Not used concurrently in tests/demo (each goroutine builds its own `Client`).
Assessment: PASS (not in concurrent use)
Severity: LOW
Notes: No action needed; each worker instantiates a fresh `Client`.
