---
title: "OAuth 2.0 & OpenID Connect: PKCE Protection & Token Rotation Replay Detection"
published: false
description: "Master OAuth 2.0 and OIDC in Go: distinguish authorization from authentication, implement PKCE S256, and detect refresh token replay with family revocation."
tags: go, security, oauth2, architecture
---

# OAuth 2.0 & OpenID Connect: PKCE Protection & Token Rotation Replay Detection

Confusing access credentials with identity assertions remains one of distributed systems' most catastrophic security design flaws. Lab 31 constructs an end-to-end identity and authorization engine in pure Go—combining Authorization Code Flow, RFC 7636 PKCE S256 binding, multi-step cryptographic ID Token verification, and single-use refresh token rotation backed by instantaneous token family revocation.

> OAuth grants delegated authority; OIDC asserts verified identity; PKCE binds the exchange; family rotation neutralizes replay.

## Real-World Vulnerability: Fragile Web Identity & Sessions

Treating access tokens as caller identity invites confused deputy exploits because access tokens lack client audience binding. Furthermore, public clients expose authorization codes during browser redirects. Without cryptographic proof-of-possession, attackers intercepting codes redeem valid tokens. Finally, long-lived static refresh tokens leak persistent API access through XSS or compromised network logs without server awareness.

## Why Naive Approaches & Implicit Flow Fail

Naive implementations parse JWT base64 payloads without cryptographic signature checks or claim bounds. Attackers trivially forge administrative identities. Concurrently, legacy Implicit Flow exposes tokens in URL fragments and browser histories. RFC 9700 formally deprecates Implicit Flow, mandating Authorization Code Grant with PKCE to prevent credential exposure in transit.

```go
func HandleAuthNaive(rawJWT string) (*User, error) {
    // FATAL VULNERABILITY: Decoding payload without HMAC/RSA verification!
    parts := strings.Split(rawJWT, ".")
    claimsBytes, _ := base64.RawURLEncoding.DecodeString(parts[1])
    var claims NaiveClaims
    json.Unmarshal(claimsBytes, &claims)
    // Skips iss, aud, exp, and nonce verification!
    return &User{ID: claims.Subject}, nil
}
```

## Mental Model: Dual-Token Architecture & Family Lineage

OAuth 2.0 issues opaque access tokens for API authorization. OIDC issues cryptographically signed ID Tokens asserting authenticated identity across six mandatory checkpoints: `iss`, `sub`, `aud`, `exp`, `iat`, and `nonce`. Concurrently, `RefreshTokenMeta` maintains session lineage through a unique `FamilyID` and `Revoked` flag, enforcing strict single-use consumption semantics across session renewals.

```go
type IDTokenClaims struct {
    Issuer     string `json:"iss"`   // Trusted Identity Provider
    Subject    string `json:"sub"`   // Unique user identity
    Audience   string `json:"aud"`   // Target Client ID
    Expiration int64  `json:"exp"`   // Expiration timestamp
    IssuedAt   int64  `json:"iat"`   // Issuance window (+/- 5m)
    Nonce      string `json:"nonce"` // Anti-replay challenge
}
type RefreshTokenMeta struct {
    FamilyID  string    // Session family lineage
    Revoked   bool      // Single-use consumption flag
    ExpiresAt time.Time // Absolute session expiry
}
```

## End-to-End Pipeline: PKCE, Token Minting, & Rotation

The end-to-end pipeline coordinates Client, Authorization Server, and Resource Server. The client derives `code_challenge = S256(verifier)` and injects an unguessable nonce. The server validates parameters and returns a short-lived authorization code. The client exchanges the code with the raw `code_verifier`. Upon cryptographic match, the server issues tokens and tracks family rotation.

```
1. Client.BuildAuthorizationRequest() -> Compute S256(verifier) & random nonce.
2. Client -> AS: GET /authorize?response_type=code&code_challenge=...&nonce=...
3. AS.Authorize() -> Validate redirect URI, persist challenge, issue AuthCode (TTL 5m).
4. Client -> AS: POST /token grant_type=authorization_code&code=...&code_verifier=...
5. AS.ExchangeCode() -> Verify SHA256(code_verifier) == challenge; mark code Used=true.
6. AS Mints: AccessToken (opaque), IDToken (HS256 JWT), RefreshToken (FamilyID).
7. Client.ParseAndVerifyIDToken() -> Validate HMAC signature, iss, aud, exp, nonce.
8. Client.RefreshTokens() -> AS rotates: mark old RT Revoked=true, issue new pair in Family.
```

## Core Implementation: Refresh Token Rotation & Family Revocation

The core session defense resides in `pkg/server/server.go`. Whenever a client requests a token refresh, the server inspects the presented token's `Revoked` state. Presenting an already-consumed token indicates token theft and active replay. The server instantly blacklists the entire `FamilyID`, invalidating all active sessions across that lineage.

```go
if s.revokedFams[meta.FamilyID] {
    return nil, ErrTokenReplayDetected // Family already blacklisted
}
if meta.Revoked {
    // REPLAY DETECTED: Stolen token reused after legitimate rotation!
    s.revokedFams[meta.FamilyID] = true // Revoke all sessions in family
    return nil, ErrTokenReplayDetected
}
meta.Revoked = true // Consume token: enforce single-use semantics
// Mint new refresh token retaining the identical FamilyID...
```

## Failure Path: Interception Rejection & Replay Defense

Concrete failure execution logs prove attack mitigation. An attacker intercepting an authorization code attempts redemption with a counterfeit verifier, but the server terminates the exchange. Furthermore, when an attacker replays a compromised refresh token, the server catches the consumed credential, revokes the token family, and terminates the victim's active session to force re-authentication.

```sh
[Step 6] Test Interception Attack (Attempt code exchange with wrong verifier)...
Interception Attack Blocked! Error: pkce verification failed: code_verifier does not match code_challenge
[Step 7] Test Refresh Token Rotation...
Rotated New Access Token : at_8c8facd2a7dce0b35b74243d16172970a30c4e94f05853bd
Rotated New Refresh Token: rt_a0aafd72efbf3ffb2f3cfafeacf71cf4d4e6503170fee590
[Step 8] Test Refresh Token Replay Detection (Attacker reuses stolen Refresh Token)...
Refresh Token Replay Detected! Error: refresh token reuse detected: entire token family revoked
Re-testing active refresh token after family revocation...
Active session revoked due to family revocation! Error: refresh token reuse detected: entire token family revoked
```

## Mechanisms vs. Symptoms: Evaluating Auth Solutions

Evaluating security mechanisms against real-world attack vectors demonstrates why partial solutions fail. Raw OAuth lacks caller identity verification. Deprecated Implicit Flow leaks tokens in URLs. Only combining strict OIDC validation, PKCE S256 binding, and single-use Refresh Token Rotation with Family Revocation produces a resilient security posture.

| Mechanism | Primary Role | Failure Mode / Attack Vector | Production Resilience |
|---|---|---|---|
| Pure OAuth 2.0 | API Authorization | Token misused as identity; no client audience check | Vulnerable (Confused Deputy) |
| Implicit Flow | Legacy Client Auth | Credentials leak in browser history & referers | Broken (RFC 9700 Deprecated) |
| OIDC ID Token | Caller Identity | Skipping iss, aud, exp, or nonce checks | Secure (Strict validation) |
| PKCE (RFC 7636) | Code Binding | Permitting 'plain' transformation method | Robust (Mandatory S256) |
| Refresh Rotation | Session Continuity | Static token replay without family tracking | Robust (With Family Revocation) |

## Three Critical Identity Architecture Decisions

First, mandate PKCE S256 unconditionally across all clients, rejecting plain transformations to preserve SHA-256 entropy. Second, execute multi-step cryptographic verification on ID Tokens, pinning issuer, audience, and nonce constraints. Third, enforce single-use refresh token semantics backed by FamilyID lineage tracking to trigger instantaneous session severance upon replay.

## Execution Evidence: Complete OIDC, PKCE, & Replay Defense

Terminal execution traces from `cmd/demo/main.go` confirm flawless lifecycle operations. Steps 1 through 5 showcase PKCE challenge creation, code issuance, legitimate token redemption, ID Token validation, and authorized resource dispatch. Steps 6 through 8 confirm active attack neutralization: PKCE drops intercepted codes and family revocation drops replayed tokens.

```sh
=== LAB 31: OAuth 2.0 & OpenID Connect (OIDC) Demo ===
[Step 1] Initiate Authorization Request with PKCE & OIDC scope...
 Generated PKCE code_challenge (S256): dbNVgXttw5HU018kaBwk3JusvwqMVanN4c4ghBsH88w
 Generated OIDC nonce: 2152fa0c9254b322d11b27dba96b99de
[Step 2] Authorization Server issues Authorization Code...
 Issued Code: 5fabbdcdb6bbd582c3811481395a5161 (expires: 15:31:47)
[Step 3] Client exchanges Code + code_verifier for Tokens...
 Access Token : at_4e8dd945eaf2aed814aaf9717a7ee6bab364db0914917bb2
 Refresh Token: rt_c624781adb52c6effa11d85bbe24fccf7f73eb20feeaf928
 ID Token     : eyJhbGciOiJIUzI1NiIsInR5cCI6Ik...
[Step 4] Validate OIDC ID Token claims... Subject: user_alice_99 | Nonce Match: SUCCESS
[Step 5] Access Resource Server using Access Token... Resource verified. Subject: user_alice_99
[Step 6] Test Interception Attack... Interception Attack Blocked! Error: pkce verification failed
[Step 7] Test Refresh Token Rotation... Rotated New RT: rt_a0aafd72efbf3ffb2f3cfafeacf71cf4...
[Step 8] Test Refresh Token Replay Detection... Replay Detected! Error: entire token family revoked
=== Demo Completed Successfully ===
```

## Scope Boundaries: What Is Not Proven

Laboratory claims reflect controlled in-memory conditions. State storage uses local maps rather than distributed databases (Redis/Postgres). ID Token signatures employ symmetric HMAC-SHA256 rather than asymmetric RS256/ES256 backed by public JWKS discovery endpoints. PKCE string comparisons do not use constant-time operations, and the implementation omits HTML login forms, MFA, consent screens, and transport-level sender-constraining (DPoP/mTLS).

## Architectural Trade-Offs: Session Security vs. Server Complexity

Enforcing PKCE S256 and token rotation introduces engineering trade-offs. Eliminating static tokens neutralizes long-term session hijacking but requires the Authorization Server to persist stateful FamilyID lineage and revocation sets. Furthermore, transient network partitions during token rotation can trigger false-positive replays when legitimate clients retry interrupted requests.

## Concurrency Validation: 10 Goroutines Contending for Refresh Token

`TestOAuth2_ConcurrentRefreshReplay` validates state consistency under concurrent race conditions. Ten goroutines simultaneously submit the identical refresh token. Because `as.Refresh()` serializes execution under `sync.Mutex`, exactly one goroutine succeeds in rotating the credential, while at least nine callers are rejected with replay errors. Go's race detector reports zero data races across 20 workers.

```go
func TestOAuth2_ConcurrentRefreshReplay(t *testing.T) {
    stolenRT := resp.RefreshToken
    var wg sync.WaitGroup; errCount := 0; var mu sync.Mutex
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            _, err := as.Refresh(stolenRT, "client_conc_ref")
            if err != nil { mu.Lock(); errCount++; mu.Unlock() }
        }()
    }
    wg.Wait()
    // Assert: errCount >= 9 due to atomic rotation and family revocation
}
```

## Before vs. After: Evolution of Web Session Security

Legacy architectures relied on Implicit Flow and stored static refresh tokens in browser `localStorage`, granting attackers persistent account control upon credential theft. Modern RFC 9700 architectures bind authorization codes cryptographically via PKCE S256, verify ID Token claims rigorously, and rotate refresh tokens continuously, turning replayed tokens into immediate intrusion alarms.

## Test 1: Cryptographic Validation of PKCE S256 Challenge

Unit tests verify RFC 7636 compliance in `pkg/pkce/pkce.go`. `GeneratePKCEPair("S256")` derives a 43-character base64url verifier with 256 bits of cryptographic entropy. Tests verify that `BASE64URL(SHA256(verifier))` matches the issued challenge. Verifiers shorter than 43 characters or exceeding 128 characters trigger strict bounds rejection via `ErrInvalidVerifierLength`.

```go
func TestPKCE_S256_Valid(t *testing.T) {
    pair, err := pkce.GeneratePKCEPair("S256")
    // Assert: len(pair.CodeVerifier) >= 43 && len(pair.CodeVerifier) <= 128
    err = pkce.Verify(pair.CodeVerifier, pair.CodeChallenge, "S256")
    // Assert: err == nil (SHA-256 base64url matching)
}
func TestPKCE_VerifierLength_Bounds(t *testing.T) {
    // Verifier < 43 characters  -> Reject ErrInvalidVerifierLength
    // Verifier > 128 characters -> Reject ErrInvalidVerifierLength
}
```

## Test 2: Multi-Step ID Token Cryptographic Verification

Unit tests prove ID Token tamper-resistance. In `TestOIDC_IDToken_TamperedSignature`, mutating a single character of the signature causes `hmac.Equal` to fail with `ErrSignatureInvalid`. Claim evaluation tests confirm that mismatched issuers, unexpected audiences, expired timestamps, or invalid nonces trigger immediate rejection without parsing untrusted claims.

```go
// TestOIDC_IDToken_TamperedSignature: Signature modification
tampered := parts[0] + "." + parts[1] + ".invalidSig=="
_, err = oidc.ParseAndVerifyIDToken(tampered, secret, iss, aud, "", now)
// Assert: err == ErrSignatureInvalid

// TestOIDC_IDToken_MismatchClaims: Claim evaluations
// 1. Wrong issuer     -> Reject ErrIssuerMismatch
// 2. Wrong audience   -> Reject ErrAudienceMismatch
// 3. Wrong nonce      -> Reject ErrNonceMismatch
// 4. Past expiration  -> Reject ErrTokenExpired
```

## Test 3: Authorization Code Interception Defense

`TestOAuth2_FullFlowAndPKCEInterception` proves that intercepted authorization codes cannot be redeemed without the local `code_verifier`. An attacker presenting a forged verifier is rejected by PKCE challenge verification. The authentic client redeems the code successfully. Subsequent attempts to reuse the redeemed code fail due to strict single-use code semantics.

```go
// 1. Legitimate client initiates flow and obtains code
ch, _ := cli.BuildAuthorizationRequest("openid profile")
ac, _ := as.Authorize("client_1", "https://app.com/cb", "openid profile", "user_sub_1", ch, "S256", cli.Nonce)

// 2. Attacker intercepts code, attempts exchange with guessed verifier
_, err = as.ExchangeCode(ac.Code, "client_1", "https://app.com/cb", "attacker_wrong_verifier_12345678901234567890123")
// Assert: err != nil (PKCE verification failed: challenge mismatch)

// 3. Legitimate client exchanges with authentic verifier -> SUCCESS
resp, err := cli.Exchange(ac.Code)

// 4. Replay attempt using same authorization code -> REJECTED (Used = true)
```

## Industry Case Study: RFC 9700 Security Consensus

RFC 9700 (OAuth 2.0 Security Best Current Practice) deprecates the Implicit Grant due to token leaks in browser histories and referer headers. Modern consensus mandates Authorization Code Flow with PKCE for all client types, exact redirect URI matching, and migration toward sender-constrained tokens like DPoP (RFC 9449) to eliminate token exfiltration across untrusted networks.

## Technical Analysis: Anatomy of Token Family Revocation

Per RFC 9700 Section 4.14, Token Family Revocation turns consumed credentials into active tripwires. Initial token exchange assigns a unique `FamilyID`. Each rotation invalidates the active token and mints a successor within the same family. Replaying an invalidated token indicates session compromise, prompting the server to revoke the entire family lineage immediately.

## Production Implementation Checklist: OAuth 2.0 & OIDC

Production deployment demands strict adherence to security invariants:
- [ ] Enforce mandatory `code_challenge` and `code_challenge_method=S256` on `/authorize`.
- [ ] Generate `code_verifier` with 32 bytes of cryptographic entropy (43–128 base64url characters).
- [ ] Limit authorization code TTL to 5 minutes under strict single-use consumption.
- [ ] Verify ID Token signatures with constant-time equality and validate `iss`, `aud`, `exp`, and `nonce`.
- [ ] Enforce single-use refresh token rotation and track `FamilyID` lineage.
- [ ] Quarantine tokens from browser `localStorage`; store session tokens in `httpOnly` `SameSite` cookies.
- [ ] Disable deprecated Implicit Grant and Resource Owner Password Credentials (ROPC).
- [ ] Guard token endpoints against concurrent race conditions using atomic transactions.

## Anti-Patterns: Common Identity Design Pitfalls

Pervasive anti-patterns include using access tokens for authentication, omitting audience or nonce claims validation, and storing tokens in browser `localStorage`. To reproduce vulnerabilities locally in the lab: disable `ac.Used` in `server.go:141` to permit code reuse, omit `meta.Revoked` to allow undetected refresh replay, or bypass signature validation in `oidc.go` to observe forgeable identities.

## Independent Verification: Testing & Running the Lab

Verify the laboratory deterministically using standard Go tooling without external dependencies:

```sh
cd labs/31-oauth2-and-oidc
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

The test suite runs 17/17 tests with 100% pass rate in ~0.32s. The race detector confirms thread safety under 20 concurrent workers in ~1.37s.

## Architectural Principles: Distributed Identity Design

Three tenets govern secure distributed identity design:
1. **Authorization is not Authentication:** OAuth grants scoped access to APIs; OIDC asserts verified caller identity.
2. **Bind Credentials Cryptographically:** Never trust credentials traversing front-channels without PKCE S256 and anti-replay nonces.
3. **Assume Token Exfiltration:** Treat refresh tokens as single-use consumables; leverage family tracking to turn replayed tokens into automated session circuit breakers.

## Source

[Software Engineering Lab 31 - OAuth 2.0 & OpenID Connect (OIDC)](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/31-oauth2-and-oidc), snapshot commit `dec7514f`. Verified pure Go implementation covering `pkg/pkce`, `pkg/oidc`, `pkg/server`, `pkg/client`, and `tests/oauth_test.go`.
