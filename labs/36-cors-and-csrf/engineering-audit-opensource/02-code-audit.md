# Code Audit

## Finding 1

Location: `internal/cors/middleware.go:62-72`
Claimed Behavior: Reflect explicit Origin when `AllowCredentials` is enabled, adhering strictly to CORS specification (wildcard `*` prohibited with credentials).
Observed Implementation: When `m.config.AllowCredentials` is true, `Access-Control-Allow-Origin` is explicitly populated with the validated request origin.
Assessment: PASS
Severity: LOW
Notes: Complies with W3C Fetch standard for credentials mode.

## Finding 2

Location: `internal/cors/middleware.go:46-59`
Claimed Behavior: Disallowed cross-origin preflight (`OPTIONS`) requests are rejected with 403 Forbidden, while standard requests proceed without CORS headers (demonstrating CORS does not block backend execution).
Observed Implementation: Preflight requests matching forbidden origins return 403. Non-preflight requests pass through to `next.ServeHTTP(w, r)` without adding CORS headers, allowing the browser to enforce read restrictions while backend execution occurs.
Assessment: PASS
Severity: LOW
Notes: Correctly demonstrates the core myth-busting premise: backend state mutation completes even when CORS headers are withheld.

## Finding 3

Location: `internal/csrf/token.go:44-59`
Claimed Behavior: Generates HMAC-SHA256 signed CSRF tokens containing `sessionID:timestamp:nonce:signature`.
Observed Implementation: `GenerateToken` constructs payload, computes HMAC-SHA256 using server secret, and base64url encodes the output.
Assessment: PASS
Severity: LOW
Notes: Clean cryptographic token implementation.

## Finding 4

Location: `internal/csrf/token.go:62-110`
Claimed Behavior: Validates signature using constant-time comparison, validates session binding, and checks TTL.
Observed Implementation: Decodes base64url token, verifies session match, validates timestamp expiration, and performs `subtle.ConstantTimeCompare` against recomputed HMAC signature.
Assessment: PASS
Severity: LOW
Notes: Timing-attack resilient validation using constant-time comparison.

## Finding 5

Location: `internal/bank/app.go:122-139`
Claimed Behavior: Thread-safe state mutation for bank account balance transfers.
Observed Implementation: Encloses account balance checks and mutations in `b.mu.Lock()` and `b.mu.Unlock()`. `GetAccount` uses `b.mu.RLock()` and returns deep copies.
Assessment: PASS
Severity: LOW
Notes: Data race safety verified across concurrent requests.
