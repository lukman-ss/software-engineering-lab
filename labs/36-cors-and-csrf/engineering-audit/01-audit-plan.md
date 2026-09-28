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
- `research-audit/07-verdict.md` (APPROVED)
Main Claims To Verify:
1. CORS middleware enforces origin checks and preflight headers without blocking backend mutation on unauthorized simple requests.
2. Wildcard `*` cannot be paired with `AllowCredentials: true` per spec.
3. Synchronizer / HMAC signed double-submit CSRF token validates session-bound nonces and enforces expiry.
4. Mutating state-changing requests without CSRF token are blocked with HTTP 403.
5. Fetch Metadata (`Sec-Fetch-Site: cross-site`) blocks cross-site mutating requests.
6. Bank state is race-safe under concurrent requests (`go test -race ./...`).
Commands To Run:
- `go test ./...`
- `go test -count=1 -race -v ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent token generation or account mutation.
- Spec deviation in CORS preflight handling or credential reflection.
- Documentation mismatch regarding attack simulations.
