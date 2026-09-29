# Engineering Audit Plan

Target Lab: labs/36-cors-and-csrf
Implementation Files:
- `internal/cors/middleware.go`
- `internal/csrf/middleware.go`
- `internal/csrf/token.go`
- `internal/bank/app.go`
- `cmd/demo/main.go`

Tests:
- `internal/cors/middleware_test.go`
- `internal/csrf/token_test.go`
- `internal/bank/app_test.go`
- `tests/integration_test.go`

Executable/Demo: `cmd/demo/main.go`

Approved Research Inputs:
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

Main Claims To Verify:
1. CORS middleware correctly handles preflight requests, allowed origins, custom/exposed headers, max age, and credential reflection rules.
2. CORS does NOT stop cross-origin state-changing POST requests (CSRF execution) on simple requests even when origins are unlisted; browser sends request and server executes before CORS check blocks reading response.
3. HMAC-SHA256 signed session-bound CSRF token validation effectively prevents cross-origin requests lacking a valid token.
4. Defense-in-depth mechanisms (`Sec-Fetch-Site` header check, custom header requirements like `X-Requested-With`) block unauthorized cross-site POSTs.
5. All code compiles, race detector passes, and interactive demo runs cleanly producing accurate outputs.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Thread safety in `BankServer` during concurrent token generation/account balance modifications.
- Spec non-compliance in CORS headers when credentials are true vs false.
- Discrepancies between demo execution output and documented README instructions/claims.
