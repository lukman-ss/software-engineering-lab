# Engineering Audit Plan

Target Lab: `labs/36-cors-and-csrf`
Implementation Files:
- `internal/cors/middleware.go`
- `internal/csrf/token.go`
- `internal/csrf/middleware.go`
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
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
Main Claims To Verify:
1. CORS middleware does not act as CSRF protection; simple cross-origin POST requests still execute state mutation if unprotected by anti-CSRF measures.
2. Anti-CSRF signed tokens (HMAC-SHA256 bound to session) reject untrusted or missing tokens across origins and sessions.
3. Modern browser defenses (`Sec-Fetch-Site`) and custom headers (`X-Requested-With`) block cross-site unauthorized requests.
4. Concurrency and race safety across account balances, session lookups, and token validations.
5. All tests compile cleanly and pass with `go test -race ./...`.
6. Demo executes deterministically without mocked/fake outputs.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent token generation or account transfers.
- Incomplete CORS header parsing or improper wildcard handling with credentials.
- Discrepancies between README instructions and codebase reality.
