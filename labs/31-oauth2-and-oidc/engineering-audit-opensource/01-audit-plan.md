# Engineering Audit Plan

Target Lab: `labs/31-oauth2-and-oidc`
Implementation Files:
- `pkg/pkce/pkce.go` — PKCE code verifier/challenge generation & verification
- `pkg/oidc/oidc.go` — ID Token signing (HS256) & claims validation
- `pkg/server/server.go` — In-memory OAuth2 Authorization Server (Auth Code + PKCE + OIDC + Refresh Token Rotation)
- `pkg/client/client.go` — Client orchestrator (request build, exchange, refresh)
- `cmd/demo/main.go` — Executable demonstration

Tests:
- `tests/oauth_test.go`

Executable/Demo:
- `go run ./cmd/demo`

Approved Research Inputs:
- `research/01-plan.md`, `research/02-sources.md`, `research/03-evidence.md`, `research/04-contradictions.md`, `research/05-report.md`, `research/06-open-questions.md`
- Key RFCs referenced: RFC 6749, RFC 7519, RFC 7636, RFC 8252, RFC 9700, OpenID Connect Core 1.0

Main Claims To Verify:
1. OAuth2 Authorization Code Flow + PKCE (S256) works end-to-end.
2. PKCE intercepts code-injection attacks (wrong verifier rejected).
3. ID Token (JWT) signed & verified with iss/aud/exp/iat/nonce claims.
4. Authorization code single-use (redemption blocked on reuse).
5. Refresh Token Rotation: new token issued, old one becomes single-use.
6. Refresh Token Replay Detection: family revocation on reuse of consumed token.
7. Access token scope validation (granted vs required).
8. Concurrency-safe state transitions under `-race`.
9. Demo output reflects real, reproducible behavior.
10. README matches code & tests.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- PKCE `plain` method still accepted (RFC 9700 deprecates it for native apps).
- ID Token signed HS256 (reference impl with shared key) — non-standard alg vs OIDC typical RS256.
- Error wrapping uses `%v` for inner errors (not unwrappable).
- Coverage gaps in boundary times (auth-code expiry, refresh expiry, skew).
