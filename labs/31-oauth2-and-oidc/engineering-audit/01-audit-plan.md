# Engineering Audit Plan

Target Lab: `labs/31-oauth2-and-oidc`
Implementation Files:
- `pkg/pkce/pkce.go`
- `pkg/oidc/oidc.go`
- `pkg/server/server.go`
- `pkg/client/client.go`
Tests:
- `tests/oauth_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/` (RFC 6749, RFC 7636, OpenID Connect Core 1.0, RFC 9700 OAuth 2.0 Security BCP)
- `research-audit/` (All gates passed)
Main Claims To Verify:
1. PKCE (`S256` and `plain`) correctly protects against authorization code interception attacks.
2. OIDC ID Token signing (HMAC-SHA256) and strict claims verification (`iss`, `aud`, `exp`, `iat`, `nonce`).
3. Authorization code single-use semantics and client/redirect URI binding.
4. Refresh Token Rotation with lineage/family tracking and immediate family invalidation on replay detection.
5. Thread-safety across concurrent token issuance, authorization, validation, and refresh operations.
Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent refresh token rotation.
- Missing edge case validations in PKCE/JWT claim checks.
- Documentation divergence between README and implementation.
