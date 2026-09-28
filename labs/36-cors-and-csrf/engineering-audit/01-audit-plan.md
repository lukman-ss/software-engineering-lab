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
- research/01-plan.md
- research/02-sources.md
- research/03-evidence.md
- research/04-contradictions.md
- research/05-report.md
- research/06-open-questions.md
Main Claims To Verify:
1. CORS response read blocking does NOT prevent server-side state execution during cross-origin state-changing requests.
2. Wildcard Access-Control-Allow-Origin with credentials active is spec-disallowed and correctly handled.
3. Synchronizer / HMAC-signed anti-CSRF token middleware blocks forged ambient-credential requests.
4. Fetch Metadata (`Sec-Fetch-Site`) blocks cross-site mutating requests.
5. Implementation runs clean under Go race detector.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent token generation/validation or bank transfer execution.
- Discrepancy between demo output, integration test expectations, and README documentation.
- Faulty/missing test cases for edge cases or negative failure paths.
