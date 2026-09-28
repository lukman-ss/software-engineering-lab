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
- `research/05-report.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. PKCE computation and verification (`S256` SHA-256 + Base64URL and `plain` methods with 43-128 length validation).
2. OIDC ID Token signing and claims validation (HMAC-SHA256, `iss`, `sub`, `aud`, `exp`, `iat`, `nonce`).
3. Authorization Server flow handling (Authorization Code grant, single-use auth code redemption, PKCE binding).
4. Refresh Token Rotation with lineage/family tracking and entire family revocation upon reuse/replay detection (RFC 9700 §4.14).
5. Thread safety and concurrency correctness across state maps under concurrent access.
6. Negative scenarios and failure paths (tampered JWT, expired token, mismatched claims, bad redirect URI, replay attempts).

Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`

Primary Risks:
- Data race conditions in map operations inside `AuthorizationServer`.
- Flawed refresh token family revocation semantics allowing revoked tokens to refresh or leaking access.
- Incomplete validation of OIDC claims or PKCE challenges.
