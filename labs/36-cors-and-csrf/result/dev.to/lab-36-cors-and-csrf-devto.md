---
title: "CORS & CSRF: Securing Backend from Cross-Origin Attacks: Demystifying CORS and Implementing Server-Side Multi-Layer Anti-CSRF Defense"
published: false
description: "Demystify CORS vs CSRF in Go: prove why CORS is not a backend firewall, and implement multi-layer anti-CSRF defense using signed tokens and Sec-Fetch-Site."
tags: go, security, webdev, architecture
---

# CORS & CSRF: Securing Backend from Cross-Origin Attacks: Demystifying CORS and Implementing Server-Side Multi-Layer Anti-CSRF Defense

Assuming CORS protects your backend from malicious write mutations is one of the most widespread security misconceptions. Lab 36 builds a spec-compliant CORS middleware and an HMAC-SHA256 signed CSRF protection engine in pure Go to expose why CORS fails against CSRF and how defense-in-depth mitigates cross-origin exploits.

> CORS governs read permissions in JavaScript, not execution permissions on backend servers.

## Fatal Misconception: CORS Is Not a Firewall

A common fallacy assumes CORS acts as a backend firewall. In reality, CORS only governs whether client JavaScript may inspect response bodies. Backend state mutation completes before browser CORS checks occur.

## Core Mental Model: SOP vs CORS vs CSRF

Web security rests on three pillars: SOP blocks cross-origin reads; CORS selectively relaxes read restrictions; CSRF exploits unrestricted cross-origin writes via ambient credentials.

## Lab 36 Banking Architecture Topology

Lab 36 hosts a banking service comparing vulnerable and protected endpoints. The vulnerable endpoint relies solely on CORS, whereas protected endpoints compose HMAC token verification, Fetch Metadata, and Custom Header middleware.

## Simple Request vs Preflight OPTIONS

The WHATWG Fetch standard bifurcates traffic into simple and preflighted requests. Simple requests fire immediately without permission checks, making them the primary vehicle for CSRF exploits. Non-simple requests mandate OPTIONS preflight before sending actual payloads.

## Anatomy of Cross-Site Request Forgery

CSRF exploits the server's ambient trust in the user's browser session. When a hidden form submits cross-site, the browser attaches authenticated session cookies automatically. The backend executes the transfer before CORS can block response reading.

## Empirical Proof: CORS Allows Backend Execution

Integration tests confirm that cross-origin simple POST requests mutate backend state regardless of CORS configuration. The vulnerable transfer endpoint returns 200 OK and deducts $300 from the victim despite omitting Access-Control-Allow-Origin headers.

```go
req := httptest.NewRequest(POST, "/api/transfer/vulnerable", body)
req.Header.Set("Origin", "https://evil-attacker.com")
req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

app.ServeHTTP(rec, req)
// Response: 200 OK!
// ACAO Header: Kosong (ditolak untuk evil-attacker.com)
// Saldo Korban: Turun dari $1000 menjadi $700!
```

## W3C Spec-Compliant CORS Middleware

The CORS middleware adheres strictly to W3C Fetch rules. Disallowed origins receiving preflight OPTIONS are rejected with 403 Forbidden. However, non-preflight simple requests pass through to subsequent handlers without CORS headers.

```go
if !m.isOriginAllowed(origin) {
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusForbidden); return
    }
    next.ServeHTTP(w, r) // Non-preflight tetap diteruskan!
    return
}
if m.config.AllowCredentials {
    w.Header().Set("Access-Control-Allow-Origin", origin)
    w.Header().Set("Access-Control-Allow-Credentials", "true")
}
```

## Credential Rules & Wildcard Origins

Web standards forbid pairing Access-Control-Allow-Origin: * with Access-Control-Allow-Credentials: true. If credentials are included, the server must explicitly whitelist and echo the validated origin, accompanied by a Vary: Origin header.

## Multi-Layer Anti-CSRF Defense Architecture

Robust anti-CSRF defense avoids single points of failure. The defense-in-depth model composes HMAC-signed double-submit tokens, Fetch Metadata verification, API custom header enforcement, and strict SameSite cookie policies.

## Signed Double-Submit CSRF Token (HMAC)

The Signed Double-Submit pattern delivers stateless CSRF defense without server memory overhead. Tokens combine Session ID, timestamp, and random nonce signed via HMAC-SHA256, protecting against cookie injection and subdomain compromise.

## Token Generation & Validation in Go

The ValidateToken method validates session binding, timestamp expiration, and cryptographic authenticity. Leveraging crypto/subtle.ConstantTimeCompare eliminates timing side-channel leakage during signature verification.

```go
parts := strings.Split(string(decoded), ":")
sessionID, tsStr, nonceStr, sigB64 := parts[0], parts[1], parts[2], parts[3]
if sessionID != expectedSessionID { return ErrInvalidToken }

ts, err := strconv.ParseInt(tsStr, 10, 64)
if err != nil || time.Since(time.Unix(ts, 0)) > tm.ttl {
    return ErrExpiredToken
}
expectedSig := hmac.New(sha256.New, tm.secret).Sum(payload)
if subtle.ConstantTimeCompare(actualSig, expectedSig) != 1 {
    return ErrInvalidToken
}
```

## Anti-CSRF Middleware in HTTP Handler

The anti-CSRF middleware wraps sensitive HTTP routes. It passes safe idempotent methods transparently while interrogating state-mutating requests. It extracts tokens from either custom headers or form payloads and rejects mismatches with 403 Forbidden.

```go
if r.Method == GET || r.Method == HEAD || r.Method == OPTIONS {
    next.ServeHTTP(w, r); return
}
cookie, err := r.Cookie(m.cookieName)
if err != nil || cookie.Value == "" {
    http.Error(w, "Unauthorized: missing cookie", 401); return
}
token := r.Header.Get("X-CSRF-Token")
if token == "" { token = r.PostFormValue("csrf_token") }

if err := m.tokenManager.ValidateToken(token, cookie.Value); err != nil {
    http.Error(w, "Forbidden: CSRF validation failed", 403); return
}
```

## Fetch Metadata: Leveraging Sec-Fetch-Site

Fetch Metadata provides browser-asserted contextual headers that JavaScript cannot tamper with. A Sec-Fetch-Site value of cross-site on state-changing methods allows backend middleware to drop malicious requests instantly.

```go
fetchSite := r.Header.Get("Sec-Fetch-Site")
if fetchSite == "cross-site" {
    if r.Method != GET && r.Method != HEAD && r.Method != OPTIONS {
        http.Error(w, "Forbidden: Cross-site rejected", 403)
        return
    }
}
```

## API Custom Header Enforcement

Requiring custom headers like X-Requested-With completely neutralizes classic HTML form CSRF. Standard HTML form tags cannot attach custom HTTP headers, immediately breaking exploit attempts against API endpoints.

```go
func RequireCustomHeaderMiddleware(headerKey, expectedValue string) {
    val := r.Header.Get(headerKey)
    if val == "" || !strings.EqualFold(val, expectedValue) {
        http.Error(w, "Forbidden: Missing custom header", 403); return
    }
}
```

## Cookie SameSite Attributes: Lax vs Strict

SameSite cookie attributes govern ambient credential transmission. SameSite=Lax withholds cookies on cross-origin POST requests while transmitting them on top-level GET links. Consequently, GET endpoints must remain strictly read-only and idempotent.

## Demo Execution: Direct Terminal Evidence

The demonstration in cmd/demo/main.go provides conclusive terminal evidence. A cross-origin exploit against the vulnerable route executes successfully, while the protected route rejects unauthenticated mutations with HTTP 403 Forbidden.

```go
[1] Initial State: Victim: $1000 | Attacker: $50
[2] DEMO: Attack on Vulnerable Endpoint (Origin: https://evil.com)
    HTTP Response Code: 200 | CORS Header: ""
    [RESULT] Transfer EXECUTED! Victim: $600 | Attacker: $450
[3] DEMO: Attack on Protected Endpoint (No CSRF Token)
    HTTP Response Code: 403
    [RESULT] Transfer BLOCKED! Victim: $600 | Attacker: $450
[4] DEMO: Legitimate Flow with Valid CSRF Token
    Obtained Signed Token: c2Vzc2lvbi12aWN0aW0tc2Vjc...
    [RESULT] Transfer Succeeded! Victim: $500 | Attacker: $550
```

## Test Suite Verification & Race Safety

All 13 unit and integration test cases pass cleanly under the Go race detector. Tests validate origin whitelisting, cryptographic signatures, expiration timeouts, Sec-Fetch-Site drops, cross-session reuse barriers, and thread concurrency.

```sh
PASS: TestBankTransfer_VulnerableEndpoint       (0.00s)
PASS: TestCORS_Preflight_Success               (0.00s)
PASS: TestCORS_Credentials_With_Wildcard       (0.00s)
PASS: TestTokenManager_GenerateAndValidate     (0.00s)
PASS: TestTokenManager_ExpiredToken            (0.02s)
PASS: TestIntegration_CORS_Does_Not_Prevent    (0.00s)
PASS: TestIntegration_CSRF_Token_Prevents      (0.00s)
PASS: TestIntegration_SecFetchSite_Protection  (0.00s)
PASS: TestIntegration_Concurrency_RaceCondition (0.00s)
PASS: go test -race ./... (Race-Clean Verified)
```

## Comparative Matrix: 5 Security Dimensions

The comparative matrix clarifies defense boundaries. CORS is ineffective against form CSRF as it targets read access. Signed CSRF tokens remain the sole mechanism providing 100% server-side autonomy independent of browser capabilities.

Mekanisme Keamanan | Target Proteksi | Lokasi Penegakan | Tangkal Form CSRF? | Ketergantungan Browser
|---|---|---|---|---|
CORS Middleware    | Relaksasi Read  | Browser & API    | TIDAK              | Penuh (Fetch Spec)
Signed CSRF Token  | Mutasi Write    | Backend Server   | YA                 | Nol (Bekerja di semua)
Fetch Metadata     | Mutasi Write    | Gateway / App    | YA                 | Butuh browser modern
Custom Header API  | Mutasi Write    | Gateway / App    | YA                 | Nol (API Client)
Cookie SameSite    | Ambient Cookie  | Browser Engine   | YA (Lax/Strict)    | Butuh browser modern

## Hidden Threat: XSS Defeats Anti-CSRF

All anti-CSRF mechanisms collapse if an application contains Cross-Site Scripting vulnerabilities. An injected script executes within the trusted origin context, allowing it to read CSRF tokens legally under SOP and authorize malicious transactions.

## Lab 36 Model Boundaries & Constraints

This lab simplifies production infrastructure for clarity. HMAC secret keys reside in memory and tokens parse via colon delimiters. Production deployments demand managed key vaults, automatic rotation, and structured token payloads.

## Security Audit Checklist: CORS & CSRF

Apply this checklist during backend security audits. Guarantee that GET routes remain strictly read-only, prohibit wildcard origins with credentials, mandate session-bound HMAC tokens for mutations, and configure SameSite cookies.

## Three Core Cross-Origin Security Principles

Securing backend systems demands conceptual clarity. Never rely on CORS to protect database integrity. Cryptographically bind anti-CSRF tokens to active sessions, and compose multi-layer defenses to neutralize cross-origin attack vectors.

## Source

[Software Engineering Lab 36 - CORS & CSRF: Securing Backend from Cross-Origin Attacks](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/36-cors-and-csrf), snapshot commit `7ced81ed`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
