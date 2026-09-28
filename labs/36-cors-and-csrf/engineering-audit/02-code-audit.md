# Code Audit

## Finding 1

Location: internal/bank/app.go:21-38, 122-146
Claimed Behavior: Thread-safe bank account state mutation and lookup.
Observed Implementation: Uses `sync.RWMutex` to guard account maps (`b.mu.Lock()` on transfers, `b.mu.RLock()` on lookups/authentication).
Assessment: PASS
Severity: LOW
Notes: Correct synchronization primitives in place; concurrency tests verify race-free execution.

## Finding 2

Location: internal/cors/middleware.go:62-72
Claimed Behavior: Fetch standard compliance for credentials vs wildcard origin.
Observed Implementation: When `AllowCredentials` is true, explicitly sets `Access-Control-Allow-Origin` to the verified request origin; does not emit wildcard `*`.
Assessment: PASS
Severity: LOW
Notes: Fully aligns with CORS/Fetch specification rules.

## Finding 3

Location: internal/csrf/token.go:44-109
Claimed Behavior: Cryptographically secure HMAC-SHA256 session-bound CSRF token generation and constant-time verification with expiration.
Observed Implementation: Uses `crypto/hmac`, `crypto/sha256`, `crypto/rand`, and `crypto/subtle.ConstantTimeCompare`. Checks session ID equivalence, TTL expiration, and HMAC integrity.
Assessment: PASS
Severity: LOW
Notes: No token forgery or replay across sessions is possible.

## Finding 4

Location: internal/csrf/middleware.go:25-53
Claimed Behavior: Middleware enforces anti-CSRF token verification on state-changing HTTP methods while allowing safe read methods.
Observed Implementation: Skips GET, HEAD, OPTIONS, TRACE. For mutating requests, requires matching session cookie and valid token passed via `X-CSRF-Token` header or form field `csrf_token`.
Assessment: PASS
Severity: LOW
Notes: Correct defense-in-depth design.
