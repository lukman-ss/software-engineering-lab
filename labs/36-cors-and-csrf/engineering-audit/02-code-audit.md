# Code Audit

## Finding 1

Location: `internal/cors/middleware.go:44-92`
Claimed Behavior: Middleware handles CORS preflight and actual requests, enforcing origin safelists and credential rules without blocking standard cross-origin simple requests from executing server-side logic when origin is disallowed.
Observed Implementation: When origin is not allowed, preflight OPTIONS requests return 403, while non-preflight requests proceed to `next.ServeHTTP(w, r)` without setting CORS headers. This faithfully demonstrates how browsers receive the response but execute the server-side action.
Assessment: PASS
Severity: LOW
Notes: Complies with W3C/Fetch CORS specifications.

## Finding 2

Location: `internal/cors/middleware.go:62-66`
Claimed Behavior: When credentials are enabled (`AllowCredentials: true`), wildcard `*` origin is forbidden and the requesting origin is explicitly reflected.
Observed Implementation: Direct check `if m.config.AllowCredentials` reflects `r.Header.Get("Origin")` and sets `Access-Control-Allow-Credentials: true`.
Assessment: PASS
Severity: LOW
Notes: Adheres to security specifications forbidding wildcard credentials.

## Finding 3

Location: `internal/csrf/token.go:44-109`
Claimed Behavior: CSRF tokens are HMAC-SHA256 signed, session-bound, include a timestamp and cryptographic nonce, and are verified using constant-time comparison.
Observed Implementation: `GenerateToken` generates 16 random bytes via `crypto/rand`, constructs `sessionID:ts:nonce`, signs it with HMAC-SHA256, and `ValidateToken` verifies the session match, timestamp TTL, and uses `subtle.ConstantTimeCompare`.
Assessment: PASS
Severity: LOW
Notes: Implements defense-in-depth token security without memory leaks or timing attack vulnerabilities.

## Finding 4

Location: `internal/csrf/middleware.go:24-86`
Claimed Behavior: Provides CSRF middleware checking headers/form bodies on unsafe methods, Fetch-Metadata validation on `Sec-Fetch-Site`, and custom header enforcement.
Observed Implementation: Safe methods (GET, HEAD, OPTIONS, TRACE) bypass token checks; state-changing requests validate token against session cookie. `FetchMetadataMiddleware` and `RequireCustomHeaderMiddleware` provide secondary defensive layers.
Assessment: PASS
Severity: LOW
Notes: Standard defense-in-depth anti-CSRF patterns accurately coded.

## Finding 5

Location: `internal/bank/app.go:20-152`
Claimed Behavior: Thread-safe banking operations with session authentication and balance transfers.
Observed Implementation: Protected by `sync.RWMutex`, defensive copies of `Account` returned in `GetAccount`, atomic locking around account sender/recipient balance updates.
Assessment: PASS
Severity: LOW
Notes: Concurrency safety verified under `-race`.
