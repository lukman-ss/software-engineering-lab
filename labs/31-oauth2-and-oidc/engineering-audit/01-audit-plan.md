# Engineering Audit Plan

Target Lab: `labs/31-oauth2-and-oidc`
Implementation Files:
- `pkg/pkce/pkce.go`
- `pkg/oidc/oidc.go`
- `pkg/server/server.go`
- `pkg/client/client.go`
- `go.mod`
Tests:
- `tests/oauth_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`
- `README.md`

Main Claims To Verify:
1. Pure Go standard library implementation without external third-party dependencies.
2. Authorization Code Grant flow with mandatory PKCE (RFC 7636 / RFC 9700) using `S256` method, successfully protecting against code interception.
3. Separation of OAuth 2.0 (Access Token for authorization/resource access) and OIDC (ID Token JWT for authentication identity).
4. OIDC ID Token signing and claims validation (`iss`, `sub`, `aud`, `exp`, `iat`, `nonce`) per OIDC Core 1.0 Section 3.1.3.7.
5. Refresh Token Rotation with token family tracking and automatic family revocation on replay detection (RFC 9700 Section 4.14).
6. Race safety under concurrent client requests and refresh operations.
7. README, engineering docs, and demo output match actual code and test execution.

Commands To Run:
```bash
go test ./...
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Thread safety and state mutation leaks in concurrent authorization / token issuance / token refresh.
- Insufficient cryptographic or claims verification checks in OIDC parser (e.g. signature bypass, missing audience/issuer check, expired tokens).
- PKCE challenge generation or verification mismatches (e.g., base64 encoding errors, padding discrepancies).
- Discrepancies between documented behaviors in README/engineering notes and actual codebase implementation.
