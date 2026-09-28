# Code Audit

## Finding 1

Location: `internal/cors/middleware.go:44-92`
Claimed Behavior: Spec-compliant CORS middleware handles origins, preflight (OPTIONS), credentials constraints, and headers.
Observed Implementation: Properly evaluates Origin headers, distinguishes preflight vs simple requests, restricts wildcards with credentials, sets Vary header, and returns appropriate status codes (204 for OPTIONS preflight, 403 for disallowed preflight origins).
Assessment: PASS
Severity: LOW
Notes: Implementation adheres to W3C/WHATWG CORS specification logic.

## Finding 2

Location: `internal/csrf/token.go:44-110`
Claimed Behavior: Signed HMAC-SHA256 session-bound tokens with TTL expiry and constant-time signature comparison.
Observed Implementation: Payload formatted as `sessionID:timestamp:nonce`, signed via HMAC-SHA256, validated using `subtle.ConstantTimeCompare`, and enforces expiry against configured TTL.
Assessment: PASS
Severity: LOW
Notes: Cryptographically sound construction preventing tampering and replay beyond expiry window.

## Finding 3

Location: `internal/csrf/middleware.go:24-86`
Claimed Behavior: Defense-in-depth CSRF protection via token middleware, `Sec-Fetch-Site` inspection, and custom header enforcement.
Observed Implementation: Excludes safe methods (GET, HEAD, OPTIONS, TRACE), extracts session cookie, checks token from header and form body, and provides standalone Fetch Metadata & Custom Header middlewares.
Assessment: PASS
Severity: LOW
Notes: Covers both traditional Synchronizer Token Pattern and modern Fetch Metadata defenses.

## Finding 4

Location: `internal/bank/app.go:20-152`
Claimed Behavior: Thread-safe bank state management simulating victim and attacker accounts with cookie authentication and transfer endpoints.
Observed Implementation: Uses `sync.RWMutex` to protect internal maps (`accounts`, `sessions`) and balance mutations. Provides vulnerable and protected endpoints.
Assessment: PASS
Severity: LOW
Notes: State manipulation is thread-safe and isolated.
