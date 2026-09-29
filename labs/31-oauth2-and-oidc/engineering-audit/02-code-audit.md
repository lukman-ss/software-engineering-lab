# Code Audit

## Finding 1

Location: pkg/pkce/pkce.go:63-71 (`Verify`)
Claimed Behavior: Constant-time comparison to prevent timing attacks on PKCE challenge verification.
Observed Implementation: Uses `expected != challenge` (string equality, NOT `subtle.ConstantTimeCompare`). For PKCE, the verifier is public during exchange; timing attack risk on code_challenge comparison is marginal in server-side usage but non-RFC-compliant for security-conscious implementations.
Assessment: WARNING
Severity: LOW
Notes: RFC 7636 does not mandate constant-time comparison. In practice the server holds the challenge server-side; attacker cannot probe timing. Acceptable for a reference lab.

---

## Finding 2

Location: pkg/server/server.go:229-238 (`Refresh`)
Claimed Behavior: Replayed (already-consumed) refresh token triggers full family revocation.
Observed Implementation: Two-phase check:
1. Line 229-231: if `revokedFams[meta.FamilyID]` → return `ErrTokenReplayDetected`
2. Line 234-238: if `meta.Revoked` → set `revokedFams[meta.FamilyID] = true` → return `ErrTokenReplayDetected`

Order is correct: a token that was revoked via phase 2 sets the family flag; subsequent calls on any token in the same family hit phase 1. Logic is sound.
Assessment: PASS
Severity: N/A

---

## Finding 3

Location: pkg/server/server.go:64 (`AuthorizationServer`)
Claimed Behavior: Concurrency-safe via `sync.Mutex`.
Observed Implementation: Single `mu sync.Mutex` guards all map mutations: `Authorize`, `ExchangeCode`, `Refresh`, `ValidateAccessToken`, `RegisterClient` all call `s.mu.Lock()/defer s.mu.Unlock()`. No map accessed outside lock.
Assessment: PASS
Severity: N/A

---

## Finding 4

Location: pkg/server/server.go:137 (`ExchangeCode`)
Claimed Behavior: Expired authorization codes rejected.
Observed Implementation: `time.Now().After(ac.ExpiresAt)` check combined with existence check. However: `!exists || time.Now().After(ac.ExpiresAt)` — if code does not exist, short-circuits before expiry check. If code exists but expired, correctly rejected.
Assessment: PASS
Severity: N/A

---

## Finding 5

Location: pkg/server/server.go:141-143 (`ExchangeCode`)
Claimed Behavior: Auth codes are single-use (replay prevented).
Observed Implementation: `ac.Used` boolean checked; set to `true` after first successful exchange. Replay of same code returns `ErrCodeAlreadyUsed`. Correct.
Assessment: PASS
Severity: N/A

---

## Finding 6

Location: pkg/oidc/oidc.go:107-109
Claimed Behavior: Tokens issued in the future are rejected.
Observed Implementation: `claims.IssuedAt > unixNow+300` — allows 5-minute clock skew. Correctly rejects iat significantly in the future.
Assessment: PASS
Severity: N/A
Notes: 300s clock skew allowance is conventional and matches OpenID Connect spec.

---

## Finding 7

Location: pkg/oidc/oidc.go:103-105
Claimed Behavior: Expired tokens rejected.
Observed Implementation: `claims.Expiration <= unixNow` — uses `<=` meaning a token expiring exactly at `now` is rejected. Strict; acceptable.
Assessment: PASS
Severity: N/A

---

## Finding 8

Location: pkg/oidc/oidc.go:80 (`ParseAndVerifyIDToken`)
Claimed Behavior: Tampered signatures rejected using `hmac.Equal`.
Observed Implementation: `hmac.Equal(sig, expectedSig)` — constant-time comparison. Correct and secure.
Assessment: PASS
Severity: N/A

---

## Finding 9

Location: pkg/server/server.go:101-107 (`Authorize`)
Claimed Behavior: Empty or invalid `code_challenge_method` rejected.
Observed Implementation: Empty challenge returns `ErrInvalidRequest`. Invalid method (`!= "S256" && != "plain"`) returns `ErrInvalidRequest`. Correct.
Assessment: PASS
Severity: N/A

---

## Finding 10

Location: pkg/client/client.go:36-56 (`BuildAuthorizationRequest`)
Claimed Behavior: Client generates PKCE pair and nonce for each authorization request.
Observed Implementation: Generates S256 pair, 16-byte random state, 16-byte random nonce. State is stored but never verified server-side (server does not receive/validate state parameter). State CSRF protection is incomplete.
Assessment: WARNING
Severity: LOW
Notes: `state` parameter is generated but never passed to or validated by the server. In a real implementation, state must be validated on redirect to prevent CSRF. This is an in-memory simulation; no HTTP redirect occurs, so the omission is architectural but acceptable for a lab demo. README does not claim CSRF state validation.

---

## Finding 11

Location: pkg/server/server.go (no token cleanup / expiry eviction)
Claimed Behavior: N/A (no claim made about cleanup)
Observed Implementation: `tokens` and `refreshMeta` maps grow unbounded. No eviction of expired tokens. Acceptable for a lab; would be a memory leak in production.
Assessment: WARNING
Severity: LOW
Notes: Not a correctness issue for lab scope.

---

## Finding 12

Location: pkg/pkce/pkce.go:29-33 (`GeneratePKCEPair`)
Claimed Behavior: Generates verifier of valid length (43-128 chars).
Observed Implementation: 32 random bytes → base64url (no padding) → exactly 43 characters. Length is always valid. `ComputeChallenge` length guard (line 48) will always pass for generated pairs.
Assessment: PASS
Severity: N/A

---

## Finding 13

Location: cmd/demo/main.go:51
Claimed Behavior: Demo shows truncated ID Token.
Observed Implementation: `resp.IDToken[:30]` — hardcoded slice. If ID token were shorter than 30 chars, this would panic. In practice the generated JWT is always much longer. Not a real risk.
Assessment: WARNING
Severity: LOW
Notes: Defensive slice would be `resp.IDToken[:min(30, len(resp.IDToken))]`. Minor.
