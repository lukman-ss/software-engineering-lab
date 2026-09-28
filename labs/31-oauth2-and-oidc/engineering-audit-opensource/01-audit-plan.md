# Engineering Audit Plan

Target Lab: labs/31-oauth2-and-oidc
Implementation Files: pkg/**.go, cmd/demo/main.go
Tests: tests/**/*.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: (pipeline override – not audited)
Main Claims To Verify:
- PKCE S256 verification prevents code interception
- OIDC ID token signature and claim validation
- Refresh token rotation and family revocation on replay
- Concurrency safety of in‑memory server
Commands To Run:
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
Primary Risks:
- Race conditions under high concurrency
- Token family revocation correctness
- Missing expiry tests for access tokens
