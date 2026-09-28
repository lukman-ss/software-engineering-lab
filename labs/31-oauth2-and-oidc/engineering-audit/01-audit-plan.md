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
- RFC 7636 (PKCE)
- RFC 9700 (OAuth 2.0 Security Best Current Practice)
- OpenID Connect Core 1.0

Main Claims To Verify:
1. PKCE (S256 and plain) code challenge computation and verification prevents code interception attacks.
2. OIDC ID Token signing (HMAC-SHA256) and validation (iss, aud, exp, iat, nonce).
3. Authorization Server enforces single-use Auth Codes, PKCE verification, client authentication, and scope validation.
4. Refresh Token Rotation with lineage tracking (family ID) revokes the entire family upon token reuse.
5. Concurrency safety across token generation, verification, and revocation under memory locking.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent authorization/exchange/refresh operations.
- State mutation flaws in token revocation and refresh token family revocation logic.
- Misalignment between README documentation and actual code API or CLI behavior.
