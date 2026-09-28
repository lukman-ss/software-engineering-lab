# Engineering Audit Plan

Target Lab: labs/31-oauth2-and-oidc
Implementation Files: pkg/pkce/pkce.go, pkg/oidc/oidc.go, pkg/server/server.go, pkg/client/client.go, cmd/demo/main.go
Tests: tests/oauth_test.go
Executable/Demo: go run ./cmd/demo
Approved Research Inputs: research/05-report.md
Main Claims To Verify:
- PKCE S256 prevents authorization code interception
- ID Token validation (iss, aud, exp, nonce, signature)
- Refresh token rotation and family revocation on reuse
- Concurrency safety
Commands To Run:
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in token storage
- Incorrect PKCE verification
- ID token signature validation bypass
- Refresh token family revocation logic flaws