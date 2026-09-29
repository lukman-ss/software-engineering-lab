# Engineering Code Audit

## Finding 1

Location: `internal/cors/middleware.go:44-92`
Claimed Behavior: CORS middleware should process requests without aborting non-preflight requests when Origin is disallowed, matching browser behavior where simple requests reach the server.
Observed Implementation: For non-preflight requests with disallowed Origin, `next.ServeHTTP(w, r)` is executed without setting CORS response headers. Preflight `OPTIONS` requests from disallowed origins return `403 Forbidden`.
Assessment: PASS
Severity: LOW
Notes: Correctly demonstrates that CORS header absence does not block backend execution of simple POST requests.

## Finding 2

Location: `internal/cors/middleware.go:62-72`
Claimed Behavior: When `AllowCredentials: true`, wildcard `*` is illegal and origin must be explicitly reflected.
Observed Implementation: When `m.config.AllowCredentials` is true, `Access-Control-Allow-Origin` is explicitly set to `origin`. When false and wildcard origin configured, it sets `*`.
Assessment: PASS
Severity: LOW
Notes: Strict compliance with W3C / Fetch Standard CORS credential rules.

## Finding 3

Location: `internal/csrf/token.go:44-109`
Claimed Behavior: Secure, session-bound HMAC-SHA256 signed CSRF token generation and constant-time validation.
Observed Implementation: Tokens encode `sessionID:timestamp:nonce:hmac`. Validation decodes base64, verifies exact session ID, checks timestamp against TTL, and uses `subtle.ConstantTimeCompare` for HMAC comparison.
Assessment: PASS
Severity: LOW
Notes: Uses constant-time comparison to prevent timing attacks. Nonce ensures token freshness.

## Finding 4

Location: `internal/bank/app.go:122-145`
Claimed Behavior: Thread-safe in-memory state updates during fund transfer operations.
Observed Implementation: Uses `b.mu.Lock()` and `defer b.mu.Unlock()` around sender and recipient balance mutations. Reads return deep copies of `Account` struct under `RLock()`.
Assessment: PASS
Severity: LOW
Notes: Clean thread-safety primitives without data race potential.

## Finding 5

Location: `internal/csrf/middleware.go:56-70`
Claimed Behavior: Fetch Metadata middleware inspects `Sec-Fetch-Site` and rejects cross-site state-changing requests.
Observed Implementation: If `Sec-Fetch-Site` is `cross-site`, non-safe methods (POST/PUT/DELETE/etc.) return `403 Forbidden`. Safe methods (GET/HEAD/OPTIONS) are allowed through.
Assessment: PASS
Severity: LOW
Notes: Accurately reflects W3C Fetch Metadata specification guidelines for browser-initiated requests.
