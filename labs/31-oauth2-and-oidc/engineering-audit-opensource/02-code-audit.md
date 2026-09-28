## Finding 1

Location: pkg/pkce/pkce.go:63-71
Claimed Behavior: PKCE verification fails when verifier does not match challenge
Observed Implementation: Verify computes expected challenge and returns ErrChallengeMismatch
Assessment: PASS
Severity: LOW
Notes: Covers interception mitigation as claimed

## Finding 2

Location: pkg/oidc/oidc.go:64-62
Claimed Behavior: ID token signed with HMAC‑SHA256, claims validated (iss, aud, exp, nonce)
Observed Implementation: SignIDToken creates header, claims JSON, HMAC‑SHA256 signature; ParseAndVerifyIDToken checks malformed, signature, issuer, audience, expiration, issued‑in‑future, nonce
Assessment: PASS
Severity: MEDIUM
Notes: Expiration check uses <= now, correct; nonce optional but validated when provided

## Finding 3

Location: pkg/server/server.go:167-187 (refresh token rotation)
Claimed Behavior: Refresh token rotation invalidates old token, issues new token with same family ID
Observed Implementation: On Refresh, meta.Revoked = true, new refresh token created with same FamilyID
Assessment: PASS
Severity: LOW
Notes: Demonstrates rotation as claimed

## Finding 4

Location: pkg/server/server.go:219-286 (replay detection)
Claimed Behavior: Replay of consumed refresh token triggers family revocation and error
Observed Implementation: Checks meta.Revoked and s.revokedFams[FamilyID]; on replay, revokes family and returns ErrTokenReplayDetected
Assessment: PASS
Severity: MEDIUM
Notes: Works as shown in demo step 8

## Finding 5

Location: pkg/server/server.go:63-72 (mutexes on maps)
Claimed Behavior: Server is safe for concurrent use
Observed Implementation: All methods lock s.mu before accessing shared maps
Assessment: PASS
Severity: LOW
Notes: Concurrency test passes; race detector shows no races

## Finding 6

Location: pkg/server/server.go:132-153 (code exchange)
Claimed Behavior: Authorization code one‑time use enforced
Observed Implementation: ac.Used set true after successful exchange; second use returns ErrCodeAlreadyUsed
Assessment: PASS
Severity: LOW
Notes: Covered by test