# Engineering Audit Plan

Target Lab: labs/31-oauth2-and-oidc
Implementation Files:
- pkg/pkce/pkce.go
- pkg/oidc/oidc.go
- pkg/server/server.go
- pkg/client/client.go

Tests:
- tests/oauth_test.go (17 test cases)

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/01-plan.md
- research/02-sources.md
- research/03-evidence.md
- research/04-contradictions.md
- research/05-report.md
- research/06-open-questions.md
- research-revision/01-revision-plan.md
- research-revision/02-changes-made.md
- research-revision/03-revision-result.md

Main Claims To Verify:
1. PKCE execution prevents code interception attack using S256 and plain code challenge methods.
2. OIDC ID Token signature validation (HMAC-SHA256) and claim checks (`iss`, `sub`, `aud`, `exp`, `iat`, `nonce`).
3. Refresh Token Rotation with lineage tracking and family-wide revocation on replay attempt.
4. Thread safety of Authorization Server operations under concurrent requests.

Commands To Run:
```bash
go test ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Unhandled race conditions under concurrent code redemption or refresh token reuse.
- Inaccurate token claim expiration or validation boundary checks.
- Discrepancy between README documentation claims and implementation behavior.
