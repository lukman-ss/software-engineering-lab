# Documentation vs Code Comparison

## README vs Implementation

### README Claim 1: `internal/cors` is spec-compliant CORS middleware with preflight, allowed origins, method/header safelists, and credential checks.

Code verification:
- `middleware.go:44-92`: Handles preflight `OPTIONS`, allowed origins, methods, headers, and `AllowCredentials`. Reflects origin explicitly when credentials used. ✓

Result: PASS

---

### README Claim 2: `internal/csrf` provides HMAC-SHA256 signed session-bound tokens, Fetch Metadata, and custom header middleware.

Code verification:
- `token.go:44-58`: HMAC-SHA256 signed, session+nonce+timestamp bound. ✓
- `middleware.go:56-70`: `Sec-Fetch-Site` inspection implemented. ✓
- `middleware.go:73-86`: `RequireCustomHeaderMiddleware` implemented. ✓

Result: PASS

---

### README Claim 3: `internal/bank` simulates cookie-authenticated balance inquiries, vulnerable transfer, and protected transfer.

Code verification:
- `app.go:55-73`: Cookie-based session authentication. ✓
- `app.go:90-99`: Balance query handler. ✓
- `app.go:102-146`: Vulnerable transfer with no CSRF check. ✓
- `app.go:149-152`: Protected transfer delegates to same logic, CSRF enforcement done by middleware wrapper. ✓

Result: PASS

---

### README Claim 4: `tests` covers cross-origin requests, preflight, race safety, and attack mitigation.

Code verification:
- `integration_test.go` contains: CORS no-block of execution, CSRF token block, legitimate flow, SecFetchSite, concurrency, custom header, header-based token, cross-session token rejection. ✓

Result: PASS

---

### README Command: `go test -v ./...` and `go test -race ./...` and `go run ./cmd/demo`

Execution results:
- `go test -v ./...`: PASS (15/15 tests pass).
- `go test -race ./...`: PASS (no races).
- `go run ./cmd/demo`: PASS (real output produced, balances updated correctly).

Result: PASS

---

## Engineering Notes vs Code

### Implementation Notes Claim: Zero third-party dependencies.

Verified: `go.mod` only contains `module labs/36-cors-and-csrf` and `go 1.22.0`. No external dependencies. ✓

Result: PASS

---

### Implementation Notes Claim: Signed token encodes `sessionID:timestamp:nonce:hmac`.

Verified: `token.go:52-58` exactly builds `payload = sessionID + ":" + ts + ":" + nonceStr`, signs it, and base64-encodes the full `payload:sig`. ✓

Result: PASS

---

### Implementation Notes Claim: Thread-safe in-memory maps with `sync.RWMutex`.

Verified: `app.go:21` declares `mu sync.RWMutex`, used at `app.go:45-63`, `122-123`. ✓

Result: PASS

---

## Mismatches Found

None identified. All README claims, engineering notes, and research-level assertions correspond directly to observable, testable behavior in the implementation.
