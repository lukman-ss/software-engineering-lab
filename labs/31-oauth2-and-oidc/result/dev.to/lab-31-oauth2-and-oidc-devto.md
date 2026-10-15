---
title: "OAuth 2.0 & OpenID Connect (OIDC): Authentication vs Authorization, PKCE Protection, and Token Rotation Replay Detection"
published: false
description: "Explore OAuth 2.0 & OpenID Connect (OIDC) in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# OAuth 2.0 & OpenID Connect (OIDC): Authentication vs Authorization, PKCE Protection, and Token Rotation Replay Detection

A production system must remain resilient under load. Lab 31 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Access tokens prove resource permission; only cryptographically signed ID tokens prove caller identity.

## Masalah Nyata: Rapuhnya Identitas & Sesi Web

In production distributed architectures, session breaches rarely stem from broken cryptographic algorithms, but rather from flawed authorization flows. Many engineering teams mistakenly use OAuth 2.0 access tokens to identify users, ignoring that access tokens are scoped for resource servers rather than clients. Furthermore, authorization codes traversing browser redirects on mobile or single-page applications remain highly vulnerable to interception unless anchored by cryptographic proof-of-possession like PKCE.

## Mengapa Pendekatan Naif & Implicit Flow Gagal?

The naive approach frequently decodes base64 JWT payloads directly without validating the cryptographic signature or mandatory claims. Such code allows an attacker to forge arbitrary user claims and instantly gain elevated privileges. Moreover, the legacy Implicit Flow has been formally deprecated by RFC 9700 because it leaks access tokens directly in URL hash fragments and browser histories, exposing credentials to catastrophic cross-site scripting (XSS) extraction.

```go
func HandleAuthNaive(rawJWT string) (*User, error) {
    // BAHAYA FATAL: Decode payload tanpa verifikasi signature HMAC/RSA!
    parts := strings.Split(rawJWT, ".")
    claimsBytes, _ := base64.RawURLEncoding.DecodeString(parts[1])
    var claims NaiveClaims
    json.Unmarshal(claimsBytes, &claims)
    // Tidak memeriksa iss, aud, exp, maupun nonce anti-replay!
    return &User{ID: claims.Subject}, nil
}
```

## Mental Model: Dual-Token Architecture & Family Lineage

The architectural mental model rests on a rigorous separation of concerns between authorization and authentication. ID Token claims establish verified caller identity bounded by six mandatory checkpoints: issuer (iss), subject (sub), audience (aud), expiration (exp), issuance window (iat), and anti-replay nonce. Concurrently, RefreshTokenMeta tracks session lineage via a unique FamilyID alongside a Revoked state flag to enforce strict single-use semantics across refresh token rotations.

```go
type IDTokenClaims struct {
    Issuer     string `json:"iss"`  // OIDC Provider sah
    Subject    string `json:"sub"`  // Identitas unik user
    Audience   string `json:"aud"`  // Client ID target
    Expiration int64  `json:"exp"`  // Waktu kedaluwarsa
    IssuedAt   int64  `json:"iat"`  // Waktu terbit (±5m)
    Nonce      string `json:"nonce"`// Anti-replay challenge
}
type RefreshTokenMeta struct {
    FamilyID  string    // Silsilah sesi tunggal
    Revoked   bool      // Status single-use consumption
    ExpiresAt time.Time // Kedaluwarsa absolut (30 hari)
}
```

## Pipeline End-to-End: PKCE, Token Minting, & Rotation

This end-to-end identity and authorization pipeline maps the sequential interactions between Client, Authorization Server, and Resource Server. During stages 1 through 3, the client establishes PKCE S256 parameters and a session nonce to obtain a short-lived authorization code. In stages 4 through 6, the client exchanges the code by supplying the unhashed code_verifier proof, prompting the server to issue access tokens, signed ID tokens, and a family-tracked refresh token. Subsequent stages verify claims and rotate tokens cleanly.

## Implementasi Inti: Refresh Token Rotation & Family Revocation

This core snippet forms the primary defense against refresh token replay attacks inside pkg/server/server.go. Whenever a client presents a refresh token, the server verifies whether the token was previously consumed via the meta.Revoked flag. If an already-revoked token is submitted, the server infers that the credential was stolen and is currently being replayed by an adversary. It immediately adds the meta.FamilyID to the revokedFams blacklist, terminating all associated active sessions.

```go
if s.revokedFams[meta.FamilyID] {
    return nil, ErrTokenReplayDetected // Keluarga sudah terkompromi
}
if meta.Revoked {
    // REPLAY TERDETEKSI: Penyerang menggunakan token yang sudah dirotasi!
    s.revokedFams[meta.FamilyID] = true // Cabut seluruh sesi aktif
    return nil, ErrTokenReplayDetected
}
meta.Revoked = true // Tandai token saat ini sebagai terpakai
// Terbitkan refresh token baru dengan FamilyID yang sama...
```

## Failure Path: Penolakan Intersepsi & Replay Serangan

This concrete failure path is extracted directly from the demo run in cmd/demo/main.go. In Step 6, an attacker who intercepted the authorization code attempts an exchange using a bogus code_verifier, but the server rejects the request due to PKCE challenge mismatch. In Step 8, the attacker replays the consumed refresh token from Step 7. The server instantly detects the replay, marks the entire token family as compromised, and invalidates all session tokens, forcing a clean re-authentication.

## Mekanisme vs Gejala: Evaluasi Solusi Autentikasi & Otorisasi

This table evaluates multiple authorization and authentication mechanisms against modern attack vectors. Relying solely on raw OAuth 2.0 for user sign-in creates severe vulnerabilities because access tokens lack caller audience constraints. The legacy Implicit Flow is formally banned by IETF RFC 9700 due to severe token leakage risks via browser history. Only an integrated combination of OIDC, PKCE S256, and Refresh Token Rotation with Token Family Revocation offers comprehensive defense.

Mekanisme | Peran Primer | Celah Kegagalan / Vektor Serangan | Ketahanan Produksi
|---|---|---|---|
OAuth 2.0 Murni | Otorisasi Akses API | Token ditukar sebagai identitas; no identity claims | Rentan (Confused Deputy)
Implicit Flow (Legacy) | Otorisasi SPA Langsung | Token bocor via URL history & referer header | Tidak Aman (RFC 9700 Deprecated)
OIDC ID Token | Autentikasi Pengguna | Skipping aud/iss/nonce verification | Aman jika divalidasi ketat
PKCE (RFC 7636) | Code Exchange Binding | Menggunakan metode 'plain' alih-alih S256 | Aman Mutlak (S256 Mandatory)
Refresh Rotation | Sesi Berkelanjutan | Static token replay tanpa pelacakan family | Aman Mutlak (dengan Family Revocation)

## 3 Keputusan Desain Kritis Arsitektur Identitas

These three architectural decisions ensure OAuth 2.0 and OIDC implementations satisfy rigorous IETF RFC 9700 standards. First, PKCE S256 is enforced unconditionally across all clients to ensure intercepted authorization codes cannot be redeemed. Second, ID Token validation must cryptographically verify digital signatures alongside strict issuer, audience, and nonce constraints. Third, refresh token rotation must maintain FamilyID lineage tracking to trigger immediate session revocation upon replay detection.

## Bukti Eksekusi: Siklus Lengkap OIDC, PKCE, & Replay Defense

The terminal execution trace verifies that the full authentication and authorization lifecycle functions flawlessly. Steps 1 through 5 demonstrate PKCE S256 challenge generation, code issuance, legitimate token exchange, ID Token claim verification, and authorized resource access. Steps 6 through 8 execute active attack simulations where code interception is neutralized by PKCE and stolen refresh token replays immediately trigger complete token family revocation.

## YANG TIDAK DIBUKTIKAN

The technical claims in this laboratory are strictly bounded by its self-contained test design. The authorization server manages state in-process without distributed database backends like Redis clusters. ID Token signing relies on symmetric HMAC-SHA256 rather than asymmetric public-key cryptography backed by dynamic JWKS endpoints. Production deployments must account for constant-time comparisons, database persistence, and advanced sender-constrained mechanisms like DPoP (RFC 9449).

## Trade-off Arsitektur: Keamanan Sesi vs Kompleksitas Server

Adopting modern OAuth 2.0 and OIDC architectures requires a calculated trade-off between security posture and operational complexity. Eliminating static tokens in favor of Refresh Token Rotation guarantees defense against persistent session hijacking, yet forces the server to maintain stateful FamilyID lineage metadata across distributed storage. Furthermore, transient network partitions during rotation exchanges can cause legitimate retries to be mistakenly flagged as replay attacks.

## Validasi Konkurensi: 10 Goroutine Berebut Refresh Token

Concurrency testing in TestOAuth2_ConcurrentRefreshReplay validates system resilience when a single stolen refresh token is replayed simultaneously by 10 goroutines. Because the Authorization Server protects the Refresh method via an exclusive sync.Mutex, state evaluations remain strictly linearizable. At most one request succeeds in rotating the credential, while the remaining nine immediately trigger ErrTokenReplayDetected, invalidating the entire family with zero data races flagged under -race.

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
    // Assert: errCount >= 9 akibat rotasi dan family revocation
}
```

## Before vs After: Evolusi Keamanan Sesi Web

This before-and-after comparison highlights the architectural paradigm shift in web session security over the past decade. Legacy systems relied on the Implicit Flow and stored static refresh tokens inside browser localStorage, inviting devastating session hijacking via cross-site scripting. Under modern RFC 9700 architectures, authorization codes are cryptographically bound by PKCE, ID Tokens undergo multi-step signature verification, and refresh tokens rotate continuously with automated replay detection.

## Test 1: Validasi Kriptografis Tantangan PKCE S256

The first unit test suite validates RFC 7636 compliance within pkg/pkce/pkce.go. GeneratePKCEPair demonstrably generates 256-bit cryptographic entropy yielding a 43-character base64url code_verifier. The test verifies that BASE64URL(SHA256(verifier)) matches the issued code_challenge. Furthermore, bounds testing ensures any verifier shorter than 43 characters or exceeding 128 characters is systematically rejected, preventing entropy degradation vulnerabilities.

```go
func TestPKCE_S256_Valid(t *testing.T) {
    pair, err := pkce.GeneratePKCEPair("S256")
    // Assert: len(pair.CodeVerifier) >= 43 && len(pair.CodeVerifier) <= 128
    err = pkce.Verify(pair.CodeVerifier, pair.CodeChallenge, "S256")
    // Assert: err == nil (SHA-256 base64url matching)
}
func TestPKCE_VerifierLength_Bounds(t *testing.T) {
    // Verifier < 43 karakter -> Tolak ErrInvalidVerifierLength
    // Verifier > 128 karakter -> Tolak ErrInvalidVerifierLength
}
```

## Test 2: Verifikasi Kriptografis & Claims ID Token

The OIDC test suite proves that ID Tokens are tamper-proof against unauthorized manipulation. In TestOIDC_IDToken_TamperedSignature, altering a single character of the JWT signature triggers instant rejection via hmac.Equal constant-time verification. Simultaneously, claim tests verify that any token presenting an unrecognized issuer, mismatched audience, expired timestamp, or deviated nonce is intercepted and rejected without exception.

```go
// TestOIDC_IDToken_TamperedSignature: Modifikasi signature
tampered := parts[0] + "." + parts[1] + ".invalidSig=="
_, err = oidc.ParseAndVerifyIDToken(tampered, secret, iss, aud, "", now)
// Assert: err == ErrSignatureInvalid
// TestOIDC_IDToken_MismatchClaims: Evaluasi claims
// 1. Wrong issuer -> Tolak ErrIssuerMismatch
// 2. Wrong audience -> Tolak ErrAudienceMismatch
// 3. Wrong nonce -> Tolak ErrNonceMismatch
// 4. Expiration lampau -> Tolak ErrTokenExpired
```

## Test 3: Pencegahan Serangan Intersepsi Authorization Code

The third test orchestrates an end-to-end interception attack simulation on public clients. When an adversary intercepts the authorization code and attempts redemption using a forged code_verifier, the server detects the cryptographic hash mismatch and aborts the exchange. The legitimate client then redeems the code using its authentic verifier. When anyone attempts to replay that redeemed code a second time, the server rejects it because authorization codes are strictly single-use.

```go
// 1. Client sah menginisiasi flow dan meminta authorization code
ch, _ := cli.BuildAuthorizationRequest("openid profile")
ac, _ := as.Authorize("client_1", "https://app.com/cb", "openid profile", "user_sub_1", ch, "S256", cli.Nonce)
// 2. Attacker menyadap ac.Code, mencoba exchange dengan verifier tebakan
_, err = as.ExchangeCode(ac.Code, "client_1", "https://app.com/cb", "attacker_wrong_verifier_12345678901234567890123")
// Assert: err != nil (PKCE verification failed: challenge mismatch)
// 3. Client sah menukar code dengan verifier asli -> SUKSES
resp, err := cli.Exchange(ac.Code) // Menerima Access, Refresh, & ID Token
// 4. Percobaan penggunaan ulang ac.Code yang sama -> DITOLAK (Used = true)
```

## Studi Kasus Industri: Konsensus Keamanan RFC 9700

IETF standardization reflects a decisive evolution in distributed web security. The legacy Implicit Flow originally designed for Single Page Applications has been excised from OAuth 2.1 and RFC 9700 because access tokens routinely leaked via browser histories and third-party analytics scripts. Modern best current practices mandate Authorization Code Flow backed by PKCE S256 across all client types, reinforced by emerging sender-constrained standards like DPoP to bind tokens to browser private keys.

## Analisis Teknis: Anatomi Token Family Revocation

The Token Family Revocation mechanism specified in RFC 9700 Section 4.14 converts passive refresh tokens into active intrusion alarms. When an adversary attempts to replay a stolen refresh token after the legitimate client has rotated it, the server recognizes the token was already consumed. Rather than merely rejecting the individual request, the server concludes that the session lineage has been breached. It revokes the entire FamilyID instantly, terminating unauthorized access and alerting the system.

## Checklist Implementasi Produksi OAuth 2.0 & OIDC

This checklist codifies nine critical requirements before deploying OAuth 2.0 and OIDC infrastructure to production enterprise environments. All authorization requests must mandate PKCE S256, and authorization codes must expire within five minutes under strict single-use rules. Client-side ID Token verification must evaluate signatures and all mandatory claims without shortcuts. Storing bearer tokens in localStorage must be eliminated in favor of hardened httpOnly SameSite cookie architectures.

## Anti-Pattern: Kesalahan Desain yang Kerap Terjadi

Ignoring fundamental identity security principles exposes enterprise user data to catastrophic compromises. The most pervasive anti-patterns include substituting access tokens for authentication and persisting long-lived credentials in browser localStorage, allowing XSS scripts to harvest sessions. Developers can reproduce these failure states locally by relaxing single-use validations in server.go and watching the test suite instantly flag security violations.

## Verifikasi Mandiri: Uji & Jalankan Lab

This laboratory is engineered for immediate, frictionless peer verification without external database infrastructure. All JWT HS256 signing and PKCE cryptographic logic is implemented using Go standard library primitives. Running go test -v executes 17 comprehensive test suites covering verifier bounds and token rotation. Executing go test -race proves thread safety under multi-worker contention, while go run demo demonstrates real attack mitigations.

```sh
cd labs/31-oauth2-and-oidc
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

## Prinsip Arsitektur: Desain Identitas Terdistribusi

These three architectural principles transcend specific frameworks and languages in modern software engineering. First, enforce an absolute separation between authorization and authentication to eliminate confused deputy vulnerabilities. Second, never trust credentials traversing public browsers without cryptographic proof-of-possession like PKCE and nonces. Third, adopt a Zero Trust mindset for user sessions: assume credentials leak, and leverage family-tracked token rotation as an automated defense mechanism.

## Source

[Software Engineering Lab 31 - OAuth 2.0 & OpenID Connect (OIDC)](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/31-oauth2-and-oidc), snapshot commit `dec7514f`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
