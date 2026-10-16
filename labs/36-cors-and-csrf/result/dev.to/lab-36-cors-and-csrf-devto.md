---
title: "CORS & CSRF: Securing Backend from Cross-Origin Attacks"
published: false
description: "Demystify CORS vs CSRF in Go: prove why CORS is not a backend firewall, and implement multi-layer anti-CSRF defense using signed tokens and Sec-Fetch-Site."
tags: go, security, webdev, architecture
---

# CORS & CSRF: Securing Backend from Cross-Origin Attacks

Assuming CORS protects your backend from malicious write mutations is one of the most widespread security misconceptions. Lab 36 builds a spec-compliant CORS middleware and an HMAC-SHA256 signed CSRF protection engine in pure Go to expose why CORS fails against CSRF and how defense-in-depth mitigates cross-origin exploits.

> SOP restricts cross-origin read; CSRF exploits cross-origin write. CORS relaxes read, it never defends write.

## Fatal Misconception: CORS Is Not a Firewall

A common fallacy assumes CORS acts as a backend firewall. In reality, CORS only governs whether client JavaScript may inspect response bodies. Backend state mutation completes before browser CORS checks occur.

Engineering teams frequently assume that omitting an attacker's origin from their CORS whitelist prevents malicious write operations. Under this assumption, when an attacker triggers an unauthorized request from `https://evil-attacker.com`, they expect the server to drop the connection. In practice, standard simple HTTP POST requests execute their state changes inside the database completely before the browser evaluates response headers. The browser's CORS enforcement merely prevents attacker scripts from inspecting the response payload, leaving unauthorized balance transfers or email mutations intact.

## Core Mental Model: SOP vs CORS vs CSRF

Web security rests on three pillars: SOP blocks cross-origin reads; CORS selectively relaxes read restrictions; CSRF exploits unrestricted cross-origin writes via ambient credentials.

The Same-Origin Policy (SOP) prohibits scripts on one origin from reading confidential responses returned by another origin. However, SOP historically permits outbound write requests, such as HTML form submissions, to allow inter-site navigation. Cross-Origin Resource Sharing (CORS) serves as a selective opt-in protocol, enabling servers to tell browsers which origins are permitted to read sensitive API responses. Conversely, Cross-Site Request Forgery (CSRF) exploits this exact asymmetry: an adversary does not need to read the server's reply to steal assets or mutate records. Defending against CSRF requires enforcing server-side write authorization rather than relying on browser-side read relaxations.

## Lab 36 Banking Architecture Topology

Lab 36 hosts a banking service comparing vulnerable and protected endpoints. The vulnerable endpoint relies solely on CORS, whereas protected endpoints compose HMAC token verification, Fetch Metadata, and Custom Header middleware.

The banking server is implemented cleanly in Go using only standard library packages (`net/http`, `crypto/hmac`, `crypto/subtle`, `sync`). The service routes incoming traffic across distinct endpoints to benchmark defensive behaviors against identical cross-origin attacks:

- `GET  /api/balance`: Inquires account balances using authenticated session cookies.
- `GET  /api/csrf-token`: Issues HMAC-SHA256 anti-CSRF tokens cryptographically bound to the active session.
- `POST /api/transfer/vulnerable`: A vulnerable balance transfer route wrapped exclusively with standard CORS middleware.
- `POST /api/transfer/protected`: A secured transfer route requiring both valid session cookies and HMAC token verification.
- `POST /api/transfer/fetch-metadata`: A secured transfer route guarded by browser-level `Sec-Fetch-Site` inspection.
- `POST /api/transfer/custom-header`: A secured transfer route requiring the `X-Requested-With` header.

## Simple Request vs Preflight OPTIONS

The WHATWG Fetch standard bifurcates traffic into simple and preflighted requests. Simple requests fire immediately without permission checks, making them the primary vehicle for CSRF exploits. Non-simple requests mandate OPTIONS preflight before sending actual payloads.

Browsers categorize an HTTP request as a "simple request" if it uses methods `GET`, `HEAD`, or `POST`, avoids custom headers, and specifies a standard encoding such as `application/x-www-form-urlencoded`, `multipart/form-data`, or `text/plain`. For these requests, browsers transmit payloads immediately without requesting permission from the server. Preflight requests (`OPTIONS`) only trigger when clients use non-simple methods like `PUT`, `DELETE`, or `PATCH`, custom HTTP headers, or `application/json` content types. Attackers weaponize simple requests using hidden HTML forms because the backend processes the mutation before any CORS headers can be inspected.

## Anatomy of Cross-Site Request Forgery

CSRF exploits the server's ambient trust in the user's browser session. When a hidden form submits cross-site, the browser attaches authenticated session cookies automatically. The backend executes the transfer before CORS can block response reading.

The sequence of exploitation follows five discrete stages:

1. A victim authenticates with their banking portal, receiving an HTTP session cookie stored in the browser cookie jar.
2. The victim visits an untrusted site (`https://evil-attacker.com`) while their banking session remains active.
3. The malicious page executes an automated, hidden HTML `<form>` submission targeting `https://bank.com/api/transfer/vulnerable`.
4. The user's browser automatically attaches ambient session cookies to the outbound cross-origin POST request.
5. The banking server validates the ambient session cookie and executes the financial transfer immediately.

Because the attacker's script is blind to the response due to SOP and CORS, the exploit succeeds anyway: the victim's account balance has already been deducted.

## Empirical Proof: CORS Allows Backend Execution

Integration tests confirm that cross-origin simple POST requests mutate backend state regardless of CORS configuration. The vulnerable transfer endpoint returns 200 OK and deducts $300 from the victim despite omitting Access-Control-Allow-Origin headers.

The test `TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution` in `tests/integration_test.go` provides empirical verification:

```go
req := httptest.NewRequest(http.MethodPost, "/api/transfer/vulnerable", strings.NewReader(form.Encode()))
req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
req.Header.Set("Origin", "https://evil-attacker.com")
req.AddCookie(&http.Cookie{Name: "session_id", Value: "session-victim-secret"})

rec := httptest.NewRecorder()
app.ServeHTTP(rec, req)

// Response Code: 200 OK!
// Header Access-Control-Allow-Origin: Omitted (denied for evil-attacker.com)
// Victim Balance: Deducted from $1000 down to $700!
```

The server returns HTTP 200 OK because the business handler completes successfully. The response lacks `Access-Control-Allow-Origin`, causing the browser console to report a CORS violation, yet the backend balance mutation remains executed in database memory.

## W3C Spec-Compliant CORS Middleware

The CORS middleware adheres strictly to W3C Fetch rules. Disallowed origins receiving preflight OPTIONS are rejected with 403 Forbidden. However, non-preflight simple requests pass through to subsequent handlers without CORS headers.

The implementation in `internal/cors/middleware.go` separates preflight validation from actual request dispatch:

```go
if !m.isOriginAllowed(origin) {
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusForbidden)
        return
    }
    next.ServeHTTP(w, r) // Non-preflight requests pass straight to business handlers!
    return
}
if m.config.AllowCredentials {
    w.Header().Set("Access-Control-Allow-Origin", origin)
    w.Header().Set("Access-Control-Allow-Credentials", "true")
}
```

When an origin fails the whitelist, preflight `OPTIONS` requests abort immediately with HTTP 403 Forbidden. However, standard simple requests proceed to `next.ServeHTTP(w, r)` without receiving CORS response headers, faithfully demonstrating spec compliance.

## Credential Rules & Wildcard Origins

Web standards forbid pairing Access-Control-Allow-Origin: * with Access-Control-Allow-Credentials: true. If credentials are included, the server must explicitly whitelist and echo the validated origin, accompanied by a Vary: Origin header.

The W3C Fetch specification explicitly forbids servers from emitting `Access-Control-Allow-Origin: *` when `Access-Control-Allow-Credentials` is set to `true`. If a backend attempts this combination, compliant browsers immediately reject the response and withhold data from JavaScript clients to prevent leaking private credentials across arbitrary domains. When credentialed requests are supported, backends must inspect incoming origins against an explicit whitelist, dynamically reflect the authorized origin, and append a `Vary: Origin` header to protect intermediary caches and CDNs from cache poisoning.

## Multi-Layer Anti-CSRF Defense Architecture

Robust anti-CSRF defense avoids single points of failure. The defense-in-depth model composes HMAC-signed double-submit tokens, Fetch Metadata verification, API custom header enforcement, and strict SameSite cookie policies.

Relying on a single security mechanism leaves systems vulnerable to perimeter bypasses or configuration drift. Lab 36 structures anti-CSRF architecture into four distinct, complementary layers:

1. **Layer 1: Signed CSRF Tokens**: Stateless HMAC-SHA256 tokens cryptographically bound to user session IDs, generation timestamps, and random nonces.
2. **Layer 2: Fetch Metadata**: Browser-asserted `Sec-Fetch-Site` header checks that drop unauthorized cross-site mutations before executing application logic.
3. **Layer 3: Custom API Headers**: Enforcement of headers like `X-Requested-With`, which native HTML forms cannot transmit.
4. **Layer 4: SameSite Cookie Policies**: Configuration of `SameSite=Lax` or `SameSite=Strict` cookie attributes to restrict ambient credential delivery during cross-site requests.

## Signed Double-Submit CSRF Token (HMAC)

The Signed Double-Submit pattern delivers stateless CSRF defense without server memory overhead. Tokens combine Session ID, timestamp, and random nonce signed via HMAC-SHA256, protecting against cookie injection and subdomain compromise.

Naive Double-Submit cookie implementations are susceptible to cookie injection attacks from compromised sibling subdomains or man-in-the-middle networks. Signed Double-Submit tokens solve this vulnerability by calculating a cryptographic signature over the token payload using a secret key held exclusively on the server:

```text
Payload:   sessionID : timestamp : nonce
Signature: HMAC-SHA256(secret_key, payload)
Wire:      base64(payload : base64(signature))
```

This architecture is completely stateless, eliminating the operational overhead of synchronizing CSRF tokens across distributed Redis caches or database tables.

## Token Generation & Validation in Go

The ValidateToken method validates session binding, timestamp expiration, and cryptographic authenticity. Leveraging crypto/subtle.ConstantTimeCompare eliminates timing side-channel leakage during signature verification.

The verification logic in `internal/csrf/token.go` enforces strict validation invariants:

```go
parts := strings.Split(string(decoded), ":")
sessionID, tsStr, nonceStr, sigB64 := parts[0], parts[1], parts[2], parts[3]
if sessionID != expectedSessionID {
    return ErrInvalidToken
}

ts, err := strconv.ParseInt(tsStr, 10, 64)
if err != nil || time.Since(time.Unix(ts, 0)) > tm.ttl {
    return ErrExpiredToken
}

expectedSig := tm.sign(payload)
if subtle.ConstantTimeCompare(actualSig, expectedSig) != 1 {
    return ErrInvalidToken
}
```

The method decomposes the incoming token, checks session equivalence, evaluates TTL expiration against current server time, and verifies the HMAC signature using `subtle.ConstantTimeCompare` to neutralize timing side-channel attacks.

## Anti-CSRF Middleware in HTTP Handler

The anti-CSRF middleware wraps sensitive HTTP routes. It passes safe idempotent methods transparently while interrogating state-mutating requests. It extracts tokens from either custom headers or form payloads and rejects mismatches with 403 Forbidden.

The middleware implementation in `internal/csrf/middleware.go` guards against unauthorized state mutations:

```go
if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
    next.ServeHTTP(w, r)
    return
}
cookie, err := r.Cookie(m.cookieName)
if err != nil || cookie.Value == "" {
    http.Error(w, "Unauthorized: missing session cookie", http.StatusUnauthorized)
    return
}
token := r.Header.Get("X-CSRF-Token")
if token == "" {
    token = r.PostFormValue("csrf_token")
}
if err := m.tokenManager.ValidateToken(token, cookie.Value); err != nil {
    http.Error(w, "Forbidden: CSRF validation failed", http.StatusForbidden)
    return
}
next.ServeHTTP(w, r)
```

Safe idempotent methods (`GET`, `HEAD`, `OPTIONS`) proceed without interruption. State-changing requests extract the session cookie and CSRF token, rejecting missing, invalid, or expired tokens with HTTP 403 Forbidden.

## Fetch Metadata: Leveraging Sec-Fetch-Site

Fetch Metadata provides browser-asserted contextual headers that JavaScript cannot tamper with. A Sec-Fetch-Site value of cross-site on state-changing methods allows backend middleware to drop malicious requests instantly.

The W3C Fetch Metadata specification introduces security headers attached directly by the browser network stack. Because JavaScript execution cannot overwrite or forge these headers, servers can inspect `Sec-Fetch-Site` to evaluate request origins:

- `same-origin`: Initiated from the exact same protocol, domain, and port (Safe).
- `same-site`: Initiated from a related subdomain within the same registrable domain.
- `cross-site`: Initiated from an external, untrusted third-party site (High Risk).
- `none`: Direct user navigation, such as typing a URL into the browser address bar.

```go
fetchSite := r.Header.Get("Sec-Fetch-Site")
if fetchSite == "cross-site" {
    if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
        http.Error(w, "Forbidden: Cross-site request rejected via Sec-Fetch-Site", http.StatusForbidden)
        return
    }
}
```

State-changing requests originating from `cross-site` contexts are immediately rejected with HTTP 403 Forbidden.

## API Custom Header Enforcement

Requiring custom headers like X-Requested-With completely neutralizes classic HTML form CSRF. Standard HTML form tags cannot attach custom HTTP headers, immediately breaking exploit attempts against API endpoints.

Standard HTML `<form>` elements cannot transmit custom HTTP headers without programmatic JavaScript interception. By mandating custom headers such as `X-Requested-With: XMLHttpRequest` on API mutation endpoints, servers ensure that standard browser form-submission exploits fail immediately:

```go
func RequireCustomHeaderMiddleware(headerKey, expectedValue string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            val := r.Header.Get(headerKey)
            if val == "" || !strings.EqualFold(val, expectedValue) {
                http.Error(w, "Forbidden: Missing required custom header", http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

Attempting to attach custom headers via JavaScript in cross-origin environments triggers a CORS preflight `OPTIONS` request, blocking unauthorized access before the payload is sent.

## Cookie SameSite Attributes: Lax vs Strict

SameSite cookie attributes govern ambient credential transmission. SameSite=Lax withholds cookies on cross-origin POST requests while transmitting them on top-level GET links. Consequently, GET endpoints must remain strictly read-only and idempotent.

Modern browser releases treat `SameSite=Lax` as the default cookie configuration. Under `SameSite=Lax`, browsers withhold session cookies when transmitting cross-site POST requests, neutralizing standard form CSRF attacks. However, browsers still attach cookies during top-level GET navigations, such as when users click external hyperlinks. For higher-security contexts, `SameSite=Strict` completely blocks ambient cookie attachment across all cross-site requests, though it requires users to re-authenticate when following external links. Because `SameSite=Lax` permits credential attachment on GET requests, application architectures must guarantee that `GET` endpoints remain strictly idempotent and read-only.

## Demo Execution: Direct Terminal Evidence

The demonstration in cmd/demo/main.go provides conclusive terminal evidence. A cross-origin exploit against the vulnerable route executes successfully, while the protected route rejects unauthenticated mutations with HTTP 403 Forbidden.

Running `go run ./cmd/demo` produces verifiable trace output illustrating attacks against both vulnerable and secured endpoints:

```text
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

The output confirms the security boundary: the vulnerable endpoint processes the attacker's transfer despite missing CORS headers, whereas the protected endpoint halts unauthenticated requests and admits authorized tokens.

## Test Suite Verification & Race Safety

All 13 unit and integration test cases pass cleanly under the Go race detector. Tests validate origin whitelisting, cryptographic signatures, expiration timeouts, Sec-Fetch-Site drops, cross-session reuse barriers, and thread concurrency.

The test suite exercises all security boundaries deterministically:

```text
=== RUN   TestBankTransfer_VulnerableEndpoint
--- PASS: TestBankTransfer_VulnerableEndpoint (0.00s)
=== RUN   TestCORS_Preflight_Success
--- PASS: TestCORS_Preflight_Success (0.00s)
=== RUN   TestCORS_Credentials_With_Wildcard_DisallowedInSpec
--- PASS: TestCORS_Credentials_With_Wildcard_DisallowedInSpec (0.00s)
=== RUN   TestTokenManager_GenerateAndValidate
--- PASS: TestTokenManager_GenerateAndValidate (0.00s)
=== RUN   TestTokenManager_ExpiredToken
--- PASS: TestTokenManager_ExpiredToken (0.02s)
=== RUN   TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution
--- PASS: TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution (0.00s)
=== RUN   TestIntegration_CSRF_Token_Prevents_Attack
--- PASS: TestIntegration_CSRF_Token_Prevents_Attack (0.00s)
=== RUN   TestIntegration_SecFetchSite_Protection
--- PASS: TestIntegration_SecFetchSite_Protection (0.00s)
=== RUN   TestIntegration_Concurrency_RaceCondition
--- PASS: TestIntegration_Concurrency_RaceCondition (0.00s)
PASS
```

Running `go test -race ./...` verifies thread safety across 20 concurrent workers generating and validating tokens simultaneously without data races.

## Comparative Matrix: 5 Security Dimensions

The comparative matrix clarifies defense boundaries. CORS is ineffective against form CSRF as it targets read access. Signed CSRF tokens remain the sole mechanism providing 100% server-side autonomy independent of browser capabilities.

Security Mechanism | Protection Target | Enforcement Layer | Blocks Form CSRF? | Browser Dependency
---|---|---|---|---
CORS Middleware | Response Read Access | Browser & API | NO | Full (W3C Fetch Spec)
Signed CSRF Token | State Write Mutation | Backend Server | YES | None (Universal)
Fetch Metadata | State Write Mutation | Gateway / App | YES | Requires Modern Browser
Custom Header API | State Write Mutation | Gateway / App | YES | None (API Clients)
Cookie SameSite | Ambient Credentials | Browser Engine | YES (Lax/Strict) | Requires Modern Browser

While Fetch Metadata and SameSite cookies provide perimeter defense, cryptographic tokens deliver complete server-side autonomy independent of browser implementation variations.

## Hidden Threat: XSS Defeats Anti-CSRF

All anti-CSRF mechanisms collapse if an application contains Cross-Site Scripting vulnerabilities. An injected script executes within the trusted origin context, allowing it to read CSRF tokens legally under SOP and authorize malicious transactions.

Cross-Site Scripting (XSS) completely bypasses CSRF defenses. If an adversary injects JavaScript into the target domain through a stored or reflected XSS flaw, that script executes directly within the trusted origin context. Under the Same-Origin Policy, the injected script is permitted to invoke `GET /api/csrf-token`, read the valid token from the DOM or JSON response, and dispatch an authenticated `POST` request with ambient cookies attached. Anti-CSRF architectures must always be paired with strict Content Security Policies (CSP), HTML output encoding, and input validation to eliminate XSS entry points.

## Lab 36 Model Boundaries & Constraints

This lab simplifies production infrastructure for clarity. HMAC secret keys reside in memory and tokens parse via colon delimiters. Production deployments demand managed key vaults, automatic rotation, and structured token payloads.

To isolate protocol mechanics, Lab 36 incorporates deliberate educational simplifications:

1. **In-Memory Secret Key**: The HMAC signing key is stored as a static byte slice in memory. Production environments require hardware security modules (HSM), AWS Secrets Manager, or HashiCorp Vault with regular key rotation.
2. **String Delimiter Parsing**: Tokens use colon delimiters (`:`) to separate fields. If a session ID contains colons, parsing will fail. Production systems should use structured serialization formats such as JSON Web Signatures (JWS) or Protocol Buffers.
3. **Browser Scope Boundaries**: SOP, CORS, and CSRF govern browser runtime environments. Non-browser HTTP clients like cURL, backend scripts, and microservice calls do not enforce these policies.

## Security Audit Checklist: CORS & CSRF

Apply this checklist during backend security audits. Guarantee that GET routes remain strictly read-only, prohibit wildcard origins with credentials, mandate session-bound HMAC tokens for mutations, and configure SameSite cookies.

Engineering teams should verify the following security controls before shipping web backends:

- [x] Verify that all `GET`, `HEAD`, and `OPTIONS` routes are strictly read-only and free of state mutations.
- [x] Ensure CORS whitelists specify explicit trusted domains without wildcard reflections.
- [x] Confirm `AllowCredentials: true` is never combined with `Access-Control-Allow-Origin: *`.
- [x] Guard all state-changing endpoints (`POST`, `PUT`, `DELETE`) with Anti-CSRF Token middleware.
- [x] Verify CSRF tokens are signed via HMAC-SHA256, bound to user session IDs, and carry valid TTL limits.
- [x] Use `crypto/subtle.ConstantTimeCompare` for all signature verifications.
- [x] Set authentication cookies with `HttpOnly`, `Secure`, and `SameSite=Lax` or `SameSite=Strict`.
- [x] Inspect `Sec-Fetch-Site` and enforce custom headers (`X-Requested-With`) at API gateways.
- [x] Implement robust Content Security Policies (CSP) to mitigate Cross-Site Scripting (XSS) bypasses.

## Three Core Cross-Origin Security Principles

Securing backend systems demands conceptual clarity. Never rely on CORS to protect database integrity. Cryptographically bind anti-CSRF tokens to active sessions, and compose multi-layer defenses to neutralize cross-origin attack vectors.

Three foundational principles guide secure cross-origin system design:

1. **Never Trust Client Boundaries for Server Integrity**: CORS is a browser read-compliance mechanism, not a database firewall. Server-side authorization and state-mutation boundaries must be enforced independently on the backend.
2. **Bind Tokens Cryptographically to Sessions**: Anti-CSRF tokens must be bound to active session identifiers using HMAC signatures to prevent token forgery and cross-session reuse attacks.
3. **Defense in Depth Wins**: Combine signed tokens, SameSite cookies, Fetch Metadata, and custom API headers to ensure resilient security even if a single defensive perimeter fails.

## Source

[Software Engineering Lab 36 - CORS & CSRF](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/36-cors-and-csrf), snapshot commit `7ced81ed`. Source: README, implementation (`internal/cors`, `internal/csrf`, `internal/bank`), tests (`tests/integration_test.go`), and execution records. Code excerpts belong in their original package context.
