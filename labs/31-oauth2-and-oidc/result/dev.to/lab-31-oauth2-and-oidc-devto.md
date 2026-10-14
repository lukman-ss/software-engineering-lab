---
title: "OAuth 2.0 & OpenID Connect (OIDC): Authorization, Authentication, and Modern Token Protection"
published: false
description: "Explore OAuth 2.0 & OpenID Connect (OIDC) in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# OAuth 2.0 & OpenID Connect (OIDC): Authorization, Authentication, and Modern Token Protection

A production system must remain resilient under load. Lab 31 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> OAuth 2.0 is for authorization (what you can access), OIDC is for authentication (who you are).

## Auth vs Authz: Kebingungan yang Berbahaya

OAuth 2.0 sering disalahpahami sebagai protokol autentikasi Access Token menjawab 'apa yang boleh diakses', bukan 'siapa penggunanya' OIDC menambahkan layer identitas di atas OAuth 2.0 untuk menjawab 'siapa user' Menyimpan bearer token di localStorage → XSS mengeksploitasi Menggunakan implicit flow → token tereksposisi di URL fragment Refresh token statis → attacker mempertahankan akses selamanya RFC 9700: Wajib Authorization Code Flow + PKCE untuk browser client OIDC Core 1.0: ID Token JWT dengan signature & claims wajib Implicit flow dihapus dari OAuth 2.1

## Empat Mekanisme Inti

Client generate code_verifier (32 byte random) code_challenge = SHA-256(code_verifier) Mencegah authorization code interception Client → AS: response_type=code, code_challenge, nonce AS → Client: authorization code (one-time, expires 5 min) Client → AS: exchange code + code_verifier for tokens JWT signed dengan HMAC-SHA256 Claims wajib: iss, sub, aud, exp, iat, nonce Client verifikasi signature, issuer, audience, expiry, nonce Setiap refresh menghasilkan token baru Token lama di-revoke Jika token lama digunakan ulang → entire token family revoked

## Authorization Request dengan PKCE

response_type=code client_id=spa-client-123 redirect_uri=https://app.example.com/callback scope=openid profile email state=<random> → anti-CSRF code_challenge=<challenge> → PKCE S256 code_challenge_method=S256 nonce=<random> → anti-replay ID Token 1. Client generates PKCE pair + nonce + state 2. Client redirects to AS with all parameters 3. AS stores challenge & nonce bound to code 4. AS redirects back with authorization code

## Token Exchange: Verifikasi PKCE

grant_type=authorization_code code=<authorization_code> redirect_uri=https://app.example.com/callback client_id=spa-client-123 code_verifier=dbNVgXttw5HU018kaBwk3JusvwqMVanN4c4ghBsH88w 1. Check code is unused & not expired 2. Compute SHA-256(code_verifier) 3. Compare with stored code_challenge 4. If match → issue tokens, mark code used Attacker intercepts authorization code Attacker does NOT have code_verifier Exchange with wrong verifier → verification fails Error: pkce verification failed: code_verifier does not match code_challenge

## ID Token: JWT Structure

header.claims.signature (3 base64url parts) Header: {"alg":"HS256","typ":"JWT"} Claims: {"iss":"...","sub":"...","aud":"...","exp":...,"iat":...,"nonce":"..."} 1. Encode header & claims to base64url 2. unsignedToken = header + '.' + claims 3. signature = HMAC-SHA256(unsignedToken, shared_secret) 4. final = unsignedToken + '.' + signature 1. Decode & verify HMAC-SHA256 signature (constant-time) 2. iss exact match → https://auth.example.com 3. aud match → spa-client-123 4. exp > now, iat < now+300s 5. nonce match → anti-replay

## Validation Implementation

```go
parts = rawJWT.split('.')
if len(parts) != 3: return ErrMalformedJWT
```

```go
unsigned = parts[0] + '.' + parts[1]
expected_sig = HMAC-SHA256(unsigned, secret)
if not hmac.Equal(sig, expected_sig): return ErrSignatureInvalid
```

```go
if claims.Issuer != expectedIssuer: return ErrIssuerMismatch
if claims.Audience != expectedAudience: return ErrAudienceMismatch
if claims.Expiration <= now: return ErrTokenExpired
if claims.IssuedAt > now+300: return ErrIssuedInFuture
if claims.Nonce != expectedNonce: return ErrNonceMismatch
```

## Token Response

Bearer token opaque string (e.g., at_4e8dd945eaf2aed814aaf9717a7ee6bab364db0914917bb2) Lifetime: 1 hour Scope-limited, used to access protected resources NOT for identity assertion Signed JWT with claims (iss, sub, aud, exp, iat, nonce) Lifetime: 1 hour Client verifies signature, issuer, audience, expiry, nonce Used to assert user identity Long-lived token (e.g., 30 days) Used to obtain new access tokens Single-use: each refresh invalidates old token Family tracking: replay revokes entire family

## Resource Server Verification

1. Extract bearer token from Authorization header 2. Validate token exists in AS token store 3. Check token scope matches requested resource 4. Return protected resource with subject from token Resource Server verified Access Token. Subject: user_alice_99 Scope check: profile → allowed Response: {"id":99,"name":"Alice","email":"alice@example.com"}

```go
GET /api/profile HTTP/1.1
Host: api.example.com
Authorization: Bearer at_4e8dd945eaf2aed814aaf9717a7ee6bab364db0914917bb2
```

## Refresh Token Rotation

1. Client sends current refresh_token 2. AS marks token.Revoked = true 3. AS generates new refresh_token with same FamilyID 4. AS generates new access_token 5. AS returns both tokens Rotated New Access Token : at_8c8facd2a7dce0b35b74243d16172970a30c4e94f05853bd Rotated New Refresh Token: rt_a0aafd72efbf3ffb2f3cfafeacf71cf4d4e6503170fee590 FamilyID preserved: rt_c624... and rt_a0aa... belong to same family Attacker reuses old refresh_token (already revoked) AS checks: meta.Revoked = true AS sets s.revokedFams[FamilyID] = true AS returns: ErrTokenReplayDetected

## Rotation Implementation

```go
meta, exists := s.refreshMeta[refreshToken]
if !exists: return ErrRefreshTokenNotFound
if s.revokedFams[meta.FamilyID]: return ErrTokenReplayDetected
if meta.Revoked: s.revokedFams[meta.FamilyID] = true; return ErrTokenReplayDetected
```

```go
meta.Revoked = true
meta.UsedAt = time.Now()
```

```go
newRtBytes := make([]byte, 24)
rand.Read(newRtBytes)
newRefreshToken := 'rt_' + hex.EncodeToString(newRtBytes)
s.refreshMeta[newRefreshToken] = &RefreshTokenMeta{FamilyID: meta.FamilyID, Revoked: false}
```

## Scenario 1: Authorization Code Interception

Attacker intercepts authorization code from redirect URL Attacker attempts token exchange with intercepted code Client code + attacker's random verifier AS computes SHA-256(attacker_verifier) Stored challenge: SHA-256(original_verifier) Mismatch: challenge does not match Error: pkce verification failed: code_verifier does not match code_challenge Token exchange rejected Attacker cannot obtain tokens Legitimate client with correct verifier succeeds

## Scenario 2: Forged / Tampered ID Token

Attacker intercepts ID Token Attacker modifies payload (e.g., change sub from user_alice to admin) Attacker generates new signature with known key Client verifies HMAC-SHA256 signature Computed signature ≠ embedded signature Return: ErrSignatureInvalid Token rejected before claims validation Attacker forges entirely new JWT without signing key HMAC-SHA256 verification fails immediately Return: ErrSignatureInvalid No claims validation needed — signature invalid

## Scenario 3: Refresh Token Replay

Attacker steals refresh token (rt_c624781adb52c6effa11d85bbe24fccf7f73eb20feeaf928) Attacker uses stolen token after legitimate refresh Client refreshes token → new refresh token issued, old revoked Attacker sends revoked refresh token to AS AS checks: meta.Revoked = true AS sets s.revokedFams[meta.FamilyID] = true AS returns: ErrTokenReplayDetected All active tokens in family immediately revoked Active session: Error: refresh token reuse detected: entire token family revoked Client must re-authenticate (force login) Attack mitigated — attacker loses access

## Scenario 4: Insecure Token Storage

Store access token in localStorage Application vulnerable to XSS attack Malicious script reads localStorage Attacker steals access token XSS payload: fetch('/api/steal', {method:'POST', body: localStorage.getItem('access_token')}) Access token sent to attacker server Attacker uses token to access protected resources No PKCE or refresh token protection — token is bearer httpOnly Secure SameSite cookie (BFF pattern) In-memory storage (react state, service variable) DPoP or mTLS sender-constrained tokens Never localStorage for bearer tokens

## Component Structure

pkg/pkce/pkce.go: PKCE code verifier + S256 challenge generation/verification pkg/oidc/oidc.go: ID Token JWT signing + verification + claims validation pkg/server/server.go: In-memory Authorization Server + rotation + replay detection pkg/client/client.go: Client helper for orchestration cmd/demo/main.go: End-to-end walkthrough demo tests/oauth_test.go: 17 unit tests + race detector go.mod: Go 1.22+ (standard library only) No external dependencies crypto/rand: 256-bit entropy for verifiers crypto/hmac, crypto/sha256: HMAC-SHA256 JWT signature encoding/json: JWT claims marshaling encoding/base64: base64url encoding

## Test Suite Coverage

TestPKCE_S256_Valid: verifier 43-128 chars, challenge valid TestPKCE_InvalidMethod: reject non-S256/plain TestPKCE_Mismatch: wrong verifier rejected TestPKCE_VerifierLength_Bounds: reject <43 or >128 TestOIDC_IDToken_Valid: signature valid, claims match TestOIDC_IDToken_TamperedSignature: modified signature rejected TestOIDC_IDToken_Expired: exp in past rejected TestOIDC_IDToken_MismatchClaims: iss/aud/nonce mismatch rejected TestOIDC_IDToken_IssuedInFuture: iat > now+300 rejected TestOIDC_MalformedJWT: wrong structure rejected TestOAuth2_FullFlowAndPKCEInterception: valid flow, interception blocked TestOAuth2_RefreshTokenRotation_AndReplayDetection: rotation works, replay revoked TestOAuth2_NegativePaths: unregistered client, bad redirect, invalid code TestOAuth2_ExpiredAuthCode_And_ExpiredTokens: expired resources rejected TestOAuth2_ConcurrentRefreshReplay: 10 goroutines, 9+ fail TestOAuth2_ConcurrencyAndRace: 20 workers, no race detected

## Test Execution Results

Command: go test -race ./... Result: ok  	labs/31-oauth2-and-oidc/tests	1.366s No data races detected Thread-safe operations verified

```go
=== RUN   TestPKCE_S256_Valid
--- PASS: TestPKCE_S256_Valid (0.00s)
=== RUN   TestOIDC_IDToken_Valid
--- PASS: TestOIDC_IDToken_Valid (0.00s)
=== RUN   TestOAuth2_FullFlowAndPKCEInterception
--- PASS: TestOAuth2_FullFlowAndPKCEInterception (0.00s)
=== RUN   TestOAuth2_RefreshTokenRotation_AndReplayDetection
--- PASS: TestOAuth2_RefreshTokenRotation_AndReplayDetection (0.00s)
PASS — 17/17 tests
```

```go
=== LAB 31: OAuth 2.0 & OpenID Connect (OIDC) Demo ===
[Step 6] Interception Attack Blocked! Error: pkce verification failed
[Step 8] Refresh Token Replay Detected! Error: entire token family revoked
=== Demo Completed Successfully ===
```

## Security Checklist

[ ] PKCE S256 wajib untuk semua browser client [ ] Implicit flow & resource owner password flow dihindari [ ] ID Token diverifikasi: signature, iss, aud, exp, iat, nonce [ ] Authorization code one-time use, expires 5 minutes [ ] Refresh token rotation setiap penggunaan [ ] Replay detection → family revocation [ ] Access token tidak digunakan untuk authentication [ ] Token storage: httpOnly cookie atau in-memory [ ] Concurrency thread-safe (mutex, race detector) RFC 6749: OAuth 2.0 Authorization Framework RFC 7636: Proof Key for Code Exchange RFC 7519: JSON Web Token RFC 8725: OAuth 2.0 Security Best Current Practice RFC 9700: OAuth 2.0 Security Best Current Practice (2025) OIDC Core 1.0: Authentication layer Use JWKS endpoint for asymmetric signing (RS256/ES256) Implement constant-time comparison for PKCE verification Use DPoP or mTLS for sender-constrained tokens Graceful eviction for expired codes/tokens Validate state parameter for anti-CSRF

## 6 Common Mistakes

Using access token for authentication Access token is NOT identity — it's authorization grant Solution: Use ID Token for user identity Not validating ID Token signature Anyone can forge JWT without signature verification Solution: Verify HMAC-SHA256 signature before claims Ignoring issuer/audience checks ID Token could be intended for another client Solution: Exact match issuer, audience exact or array match Storing bearer token in localStorage XSS attacks can read localStorage and steal token Solution: httpOnly Secure SameSite cookie or in-memory Using implicit flow Token exposed in URL fragment — browser history, logs, referrer Solution: Authorization Code Flow + PKCE (RFC 9700) Static refresh tokens Same token reused indefinitely allows permanent access Solution: Single-use rotation with family tracking

## Case Study: Attack Flow

Generated PKCE code_challenge (S256): dbNVgXttw5HU018kaBwk3JusvwqMVanN4c4ghBsH88w Generated OIDC nonce: 2152fa0c9254b322d11b27dba96b99de Access Token: at_4e8dd945eaf2aed814aaf9717a7ee6bab364db0914917bb2 ID Token Subject: user_alice_99 ID Token Issuer: https://auth.example.com ID Token Audience: spa-client-123 ID Token Nonce Match: SUCCESS Interception Attack Blocked! Error: pkce verification failed Rotated New Refresh Token: rt_a0aafd72efbf3ffb2f3cfafeacf71cf4d4e6503170fee590 Refresh Token Replay Detected! Error: entire token family revoked Active session revoked due to family revocation! Error: refresh token reuse detected

## Lab Implementation Boundaries

Store menggunakan in-memory map (not real database) No network latency or real HTTP transport Tokens represented as strings, not cryptographically secure random HMAC-SHA256 symmetric signing for demo simplicity Production should use asymmetric (RS256/ES256) with JWKS No key rotation mechanism demonstrated Demo does not validate state parameter (client-side only) Production must validate state for anti-CSRF No PKCE constant-time comparison (audit: LOW severity)

## Three Core Principles

OAuth 2.0 = Authorization (access token, scopes) OIDC = Authentication (ID Token JWT, identity claims) Never use access token as identity assertion Access Token: Validate existence + scope in AS token store ID Token: Verify HMAC-SHA256 signature + all claims (iss, aud, exp, iat, nonce) Refresh Token: Track usage, rotate, detect replay PKCE prevents code interception ID Token nonce prevents replay Refresh token rotation prevents long-term token reuse Family revocation limits blast radius

## Source

[Software Engineering Lab 31 - OAuth 2.0 & OpenID Connect (OIDC)](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/31-oauth2-and-oidc), snapshot commit `cf0cd907`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
