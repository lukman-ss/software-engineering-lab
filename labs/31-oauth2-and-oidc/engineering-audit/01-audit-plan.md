# Engineering Audit Plan

Target Lab: labs/31-oauth2-and-oidc
Implementation Files:
- pkg/pkce/pkce.go
- pkg/oidc/oidc.go
- pkg/server/server.go
- pkg/client/client.go
Tests:
- tests/oauth_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- RFC 6749 (OAuth 2.0)
- RFC 7636 (PKCE)
- RFC 9700 (OAuth 2.0 Security Best Current Practice - Section 4.14 Refresh Token Rotation)
- OpenID Connect Core 1.0 (ID Token claims and validation)
Main Claims To Verify:
1. PKCE (S256 and plain methods) protects against code interception attacks.
2. OIDC ID Tokens are valid HMAC-SHA256 signed JWTs with correct claim verification (`iss`, `sub`, `aud`, `exp`, `nonce`).
3. Authorization Server enforces single-use auth codes and PKCE verification.
4. Refresh Token Rotation (RFC 9700) assigns family IDs, detects reuse of revoked refresh tokens, and invalidates the entire token family.
5. Concurrency safety across authorization, code exchange, access token validation, and token refresh.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go test -count=1 ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent token exchanges or refresh token rotation.
- Missing validations on ID token claims or PKCE parameters.
- Mismatch between README documentation/claims and actual Go implementation.
