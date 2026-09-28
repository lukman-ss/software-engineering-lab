# Docs vs Code Audit

## README vs Implementation Comparison

1. **Architecture Claim**:
   - `internal/cors`: Spec-compliant CORS middleware (`OPTIONS` preflight, allowed origins, method/header safelists, credential checks).
   - Code reality: Implemented exactly in `internal/cors/middleware.go`. Matches.
   - `internal/csrf`: Anti-CSRF mechanisms including HMAC-SHA256 signed session-bound tokens, Fetch Metadata (`Sec-Fetch-Site`), and API custom header middleware.
   - Code reality: Implemented in `internal/csrf/token.go` and `internal/csrf/middleware.go`. Matches.
   - `internal/bank`: Bank application service simulating cookie-authenticated balance inquiries, vulnerable transfer endpoints, and protected transfer endpoints.
   - Code reality: Implemented in `internal/bank/app.go`. Matches.
   - `cmd/demo`: Runnable CLI program showcasing attacks against vulnerable vs. protected configurations.
   - Code reality: Implemented in `cmd/demo/main.go`. Output reproduces identical results to documentation notes.

2. **Test Instructions**:
   - Instructions specify `go test -v ./...`, `go test -race ./...`, and `go run ./cmd/demo`.
   - All documented commands run with exit code 0 and exact expected outputs.

## Mismatch Inventory
- `DOC_CODE_MISMATCH`: None observed.
- `TEST_CLAIM_MISMATCH`: None observed.
- `RESEARCH_IMPLEMENTATION_MISMATCH`: None observed.
