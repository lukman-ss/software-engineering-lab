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
- research/01-plan.md
- research/02-sources.md
- research/03-evidence.md
- research/04-contradictions.md
- research/05-report.md
- research/06-open-questions.md

Main Claims To Verify:
1. PKCE S256 verification prevents authorization code interception attacks (RFC 7636 / RFC 9700).
2. OIDC ID Token claims (`iss`, `sub`, `aud`, `exp`, `nonce`) are validated correctly using HMAC-SHA256 signatures.
3. Refresh token rotation issues new access/refresh tokens and invalidates consumed refresh tokens.
4. Refresh token replay causes complete family revocation (RFC 9700 Section 4.14).
5. State transitions and single-use authorization code constraints are enforced.
6. Execution and race detector passes cleanly without data races.

Commands To Run:
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Data race conditions in concurrent authorization, code exchange, token rotation, or token validation requests.
- Incorrect implementation of Refresh Token Family Revocation or PKCE code verification.
- Discrepancies between README / documentation claims and actual code behavior.
