# Engineering Audit Plan

Target Lab: labs/31-oauth2-and-oidc
Implementation Files:
- pkg/pkce/pkce.go
- pkg/oidc/oidc.go
- pkg/server/server.go
- pkg/client/client.go
- cmd/demo/main.go
Tests:
- tests/oauth_test.go
Executable/Demo:
- go run ./cmd/demo
Approved Research Inputs:
- research/01-plan.md through research/06-open-questions.md (Audit Status: APPROVED)
Main Claims To Verify:
1. OAuth 2.0 Authorization (Access Tokens) vs OIDC Authentication (ID Tokens)
2. PKCE S256 challenge/verification prevents code interception
3. Strict ID Token signature and claims validation (iss, aud, exp, iat, nonce)
4. Refresh token rotation with single-use and family revocation on reuse
5. Thread-safe in-memory storage with mutex protection
6. Demo demonstrates happy path and attack scenarios
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in shared maps (authCodes, tokens, refreshMeta, revokedFams)
- Incorrect error handling (e.g., token expiration vs revocation)
- PKCE verifier length validation
- Refresh token family revocation logic correctness