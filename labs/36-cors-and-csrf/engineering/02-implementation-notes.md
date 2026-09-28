# Implementation Notes

## Files Added

- `go.mod`: Module definition for `labs/36-cors-and-csrf` using Go 1.22+.
- `internal/cors/middleware.go`: CORS middleware handling simple requests, non-simple preflight `OPTIONS`, allowed origins/headers/methods, and credential security restrictions.
- `internal/cors/middleware_test.go`: Unit tests for origin checks, wildcard credentials handling, and preflight behavior.
- `internal/csrf/token.go`: HMAC-SHA256 signed CSRF token generator and validator bound to session IDs.
- `internal/csrf/middleware.go`: Anti-CSRF verification middleware, Fetch Metadata (`Sec-Fetch-Site`) middleware, and API custom header middleware.
- `internal/csrf/token_test.go`: Unit tests for token generation, expiration, tampering detection, and middleware filtering.
- `internal/bank/app.go`: In-memory bank domain with session authentication, balance query, vulnerable transfer, and protected transfer handlers.
- `internal/bank/app_test.go`: Unit test confirming vulnerable bank transfer execution upon forged request.
- `tests/integration_test.go`: End-to-end integration tests confirming:
  - CORS does NOT protect server-side state from mutating cross-origin simple requests.
  - CSRF tokens reliably block unauthorized state mutations.
  - Legitimate clients successfully query tokens and complete transactions.
  - Modern browser `Sec-Fetch-Site: cross-site` protection works.
  - Concurrent token requests pass Go race detector without data races.
- `cmd/demo/main.go`: Interactive runnable CLI demonstration comparing vulnerable and protected flows.

## Core Design Decisions

1. **Pure Go Standard Library**:
   - Uses Go 1.22+ standard HTTP routing (`net/http.ServeMux`) and standard cryptographic libraries (`crypto/hmac`, `crypto/sha256`, `crypto/rand`, `crypto/subtle`).
   - Zero third-party dependencies.
2. **Signed Token Pattern**:
   - Generates tokens encoding `sessionID:timestamp:nonce:hmac` to ensure stateless verification while strictly binding the CSRF token to the authenticated session.
3. **Spec-Compliant CORS Behavior**:
   - Enforces the Fetch Standard rule where `Access-Control-Allow-Origin: *` cannot be combined with `Access-Control-Allow-Credentials: true`. Origin is explicitly reflected only when in allowed list.

## Implementation-Specific Choices

- Bank state and session storage are kept in thread-safe in-memory maps using `sync.RWMutex` to facilitate fast, isolated, and repeatable test execution.
- Token TTL is configurable with a default of 1 hour.

## Known Limitations

- In-memory state does not persist across application restarts.
- Cookie transport in tests simulates HTTP requests with `Cookie` headers rather than headless browser rendering engines.

## Trade-offs

- Choosing pure standard library requires writing minimal CORS and CSRF token primitives rather than pulling in external middleware packages, ensuring zero dependency bloat and transparent verification.

## What Is Demonstrated

- The fundamental difference between SOP/CORS (browser-enforced response read policy) and CSRF (browser-executed ambient-credential write vulnerability).
- Proof that strict CORS configuration fails to prevent server-side state mutation if simple POST requests are processed without anti-CSRF measures.
- Multi-layer CSRF defenses: HMAC-signed tokens, Fetch Metadata (`Sec-Fetch-Site`), and API custom headers.

## What Is Not Demonstrated

- Browser UI rendering or DOM manipulation.
- Network-level packet inspection.
