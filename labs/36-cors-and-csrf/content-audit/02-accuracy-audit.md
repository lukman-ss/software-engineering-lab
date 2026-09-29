# Content Audit Accuracy Matrix

Target Lab: `labs/36-cors-and-csrf`
Audit Date: 2026-09-29

## Cross-Check vs Research

| Claim in Content | Research Source (`research/05-report.md`) | Status | Notes |
|---|---|---|---|
| CORS is not a backend firewall, but SOP relaxation for read access | Finding 1, Finding 2 | PASS | Accurately explains browser enforcement vs server execution |
| Simple POST requests bypass preflight `OPTIONS` | Finding 3 | PASS | WHATWG Fetch simple request conditions matched |
| `Access-Control-Allow-Origin: *` does not protect against CSRF | Finding 2 | PASS | Correctly highlights that server write already occurred before browser inspects response |
| Credentialed requests prohibit `Access-Control-Allow-Origin: *` | Finding 2 | PASS | Enforced in code and explained according to Fetch standard |
| Signed Double-Submit CSRF Token (HMAC-SHA256 bound to Session ID) | Finding 5 | PASS | Defeats cookie injection and subdomain takeover vulnerabilities |
| Constant-time comparison prevents timing attacks | Finding 5 | PASS | Documented and verified via `subtle.ConstantTimeCompare` |
| Fetch Metadata (`Sec-Fetch-Site`) defends against cross-site requests | Finding 6 | PASS | Correctly attributes browser-set immutable header |
| Custom header requirement blocks standard HTML form CSRF | Finding 3 | PASS | Accurately identifies browser limitation on custom headers for forms |
| `SameSite=Lax` sends cookie on top-level GET navigation | Finding 4 | PASS | Explicit warning provided: GET requests must remain safe/idempotent |
| XSS bypasses all CSRF mitigations on same origin | Limitations | PASS | Explicitly noted in warnings, checklist, and key takeaways |

## Cross-Check vs Engineering

| Claim / Code in Content | Actual Code / Test File | Status | Notes |
|---|---|---|---|
| CORS middleware passes non-preflight disallowed origins to next handler | `internal/cors/middleware.go:52-58` | PASS | Verified in integration test `TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution` |
| Preflight `OPTIONS` rejected with 403 Forbidden for disallowed origins | `internal/cors/middleware.go:49-51` | PASS | Verified in unit test `TestCORS_DisallowedOrigin_Preflight` |
| Token format: `base64(sessionID:timestamp:nonce:signature)` | `internal/csrf/token.go:44-59` | PASS | Verbatim match with implementation |
| Token validation checks TTL, session ID matching, and HMAC | `internal/csrf/token.go:62-110` | PASS | Verbatim match with implementation |
| Anti-CSRF Middleware skips safe methods (`GET`, `HEAD`, `OPTIONS`, `TRACE`) | `internal/csrf/middleware.go:26-30` | PASS | Verbatim match with implementation |
| Fetch Metadata Middleware checks `Sec-Fetch-Site == "cross-site"` | `internal/csrf/middleware.go:56-70` | PASS | Verbatim match with implementation |
| Custom Header Middleware checks presence and value of custom header | `internal/csrf/middleware.go:73-86` | PASS | Verbatim match with implementation |
| Vulnerable endpoint allows cross-origin transfer: victim $1000 -> $700 | `internal/bank/app.go`, `cmd/demo/main.go`, `tests/integration_test.go` | PASS | Verbatim match with demo output and integration tests |
| Protected transfer endpoint rejects attacker without token (403) | `tests/integration_test.go:89-115` | PASS | Verbatim match |
| Token delimiter `:` parsing caveat (GAP-02) disclosed | `content/01-content-brief.md`, `content/02-master-draft.md:332-333` | PASS | Explicitly documented in warnings and production considerations |
| Modular middleware routing structure accurately portrayed | `content/02-master-draft.md:94-100, 123-136`, `content/04-diagrams.md:141-145` | PASS | Accurately notes that Fetch Metadata & Custom Header are tested on separate endpoints |
