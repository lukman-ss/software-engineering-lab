# Engineering Audit Plan

Target Lab: labs/36-cors-and-csrf
Implementation Files:
- `internal/cors/middleware.go`
- `internal/csrf/token.go`
- `internal/csrf/middleware.go`
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
1. CORS middleware enforces origin validation and credential rules (spec compliance).
2. CORS does NOT stop cross-origin ambient cookie request execution on state-changing POST requests.
3. HMAC-SHA256 session-bound CSRF token middleware effectively blocks cross-origin requests lacking valid tokens.
4. Sec-Fetch-Site metadata checks reject cross-site state-changing requests.
5. All tests compile and pass without race conditions (`go test -race ./...`).
6. Demo output is authentic and deterministic matching claims.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent account state modification.
- Mismatch between README documentation/claims and actual implementation behavior.
- Missing negative test cases for token tampering or expired sessions.
