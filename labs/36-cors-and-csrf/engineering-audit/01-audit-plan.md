# Engineering Audit Plan

Target Lab: `labs/36-cors-and-csrf`
Implementation Files:
- `internal/cors/middleware.go`
- `internal/csrf/middleware.go`
- `internal/csrf/token.go`
- `internal/bank/app.go`

Tests:
- `internal/cors/middleware_test.go`
- `internal/csrf/token_test.go`
- `internal/bank/app_test.go`
- `tests/integration_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`

Main Claims To Verify:
1. CORS rejection of origin does not prevent backend execution of state-mutating simple requests (`POST` form).
2. `Access-Control-Allow-Origin: *` cannot be combined with `Access-Control-Allow-Credentials: true` (Fetch spec compliance).
3. Preflight `OPTIONS` behaves correctly for allowed vs disallowed origins.
4. Anti-CSRF signed tokens (HMAC-SHA256, session-bound, timestamped, nonce) prevent forged submissions while allowing legitimate requests.
5. Cross-session token reuse is blocked.
6. Modern defenses (`Sec-Fetch-Site: cross-site` rejection and custom header checks) protect against unauthorized cross-site requests.
7. Concurrency safety: Token generation and in-memory account access are race-free.
8. Demo and integration tests produce real, verifiable output matching documentation.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions in bank account updates or token validation.
- Incomplete error propagation or missing validation branches in token parser.
- Mismatch between README documentation and executable behavior.
- Insecure HMAC verification (timing attacks).
