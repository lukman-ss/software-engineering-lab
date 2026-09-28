# Engineering Audit Plan

Target Lab: labs/31-oauth2-and-oidc
Implementation Files:
- pkg/pkce/pkce.go
- pkg/oidc/oidc.go
- pkg/server/server.go
- pkg/client/client.go
- cmd/demo/main.go
Tests:
- tests/oauth_test.go (13 tests)
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs: (pipeline override — research/content not audited in this stage)
Main Claims To Verify:
1. Authorization Code flow with mandatory PKCE (S256) blocks code interception
2. Auth code single-use enforcement
3. OIDC ID Token (HS256 JWT) signature + claims (iss/aud/exp/iat/nonce) validation
4. Access Token scope enforcement (authorization vs authentication separation)
5. Refresh token rotation with family revocation on replay
6. Concurrency safety under race detector
7. Demo output is real and reproducible
Commands To Run:
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Token family revocation bypass
- Race on shared maps
- Claim validation bypass (iss/aud/exp/nonce)
- Fake demo / unverified results
