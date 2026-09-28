# Code Audit

## Finding 1

Location: `internal/cors/middleware.go:44-93`
Claimed Behavior: Spec-compliant CORS middleware handling preflight requests and header decoration.
Observed Implementation: Handlers check origin against allowed list, set `Access-Control-Allow-Origin` dynamically or wildcard, handle preflight `OPTIONS` returning 204 No Content, enforce non-wildcard when `AllowCredentials` is true.
Assessment: PASS
Severity: LOW
Notes: Correctly demonstrates CORS as browser-side read relaxation rather than backend write prevention.

## Finding 2

Location: `internal/csrf/token.go:44-110`
Claimed Behavior: Cryptographically secure HMAC-SHA256 signed session-bound CSRF token generation and validation.
Observed Implementation: Constructs payload with sessionID, timestamp, and random nonce; generates HMAC-SHA256 signature; validates using `subtle.ConstantTimeCompare` and checks TTL expiry.
Assessment: PASS
Severity: LOW
Notes: Implements constant-time signature comparison to prevent timing attacks.

## Finding 3

Location: `internal/bank/app.go:122-145`
Claimed Behavior: Thread-safe bank account balance mutations.
Observed Implementation: Guarded by `sync.RWMutex` (`b.mu.Lock()` and `defer b.mu.Unlock()`). Account state mutations are completely protected from concurrent data races.
Assessment: PASS
Severity: LOW
Notes: Tested under Go race detector with zero race warnings.
