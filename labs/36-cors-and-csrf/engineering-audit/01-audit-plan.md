# Engineering Audit Plan

Target Lab: labs/36-cors-and-csrf
Implementation Files:
- internal/cors/middleware.go
- internal/csrf/token.go
- internal/csrf/middleware.go
- internal/bank/app.go
- cmd/demo/main.go

Tests:
- internal/cors/middleware_test.go
- internal/csrf/token_test.go
- internal/bank/app_test.go
- tests/integration_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/05-report.md
- research-audit/07-verdict.md
- engineering/01-design.md
- engineering/02-implementation-notes.md

Main Claims To Verify:
1. CORS does NOT protect against CSRF attacks; unauthorized cross-origin simple POST requests still execute on the server unless protected by anti-CSRF measures.
2. CORS middleware strictly follows the W3C/Fetch standard (preflight handling, origin reflection when credentials are true, no wildcard with credentials).
3. Anti-CSRF Token Manager uses HMAC-SHA256 session binding, constant-time comparison, and expiry checking.
4. Fetch-Metadata (`Sec-Fetch-Site`) middleware blocks cross-site state-changing requests while allowing same-origin and safe methods.
5. Custom header requirement blocks standard HTML form CSRF requests.
6. Concurrency safety: Bank state mutations, account reads, and token operations are thread-safe under `-race`.

Commands To Run:
- `go test -count=1 ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Thread-safety / race conditions during concurrent account balance transfers or reads.
- Discrepancies between demo execution output and documented examples.
- Overclaiming of CORS security guarantees.
