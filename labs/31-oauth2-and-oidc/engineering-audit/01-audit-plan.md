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
- research/05-report.md
- engineering/01-design.md
- engineering/02-implementation-notes.md

Main Claims To Verify:
1. PKCE (RFC 7636/9700) implementation using S256 method prevents code interception attacks.
2. OIDC ID Token issuing and cryptographic HMAC-SHA256 signature and standard claim (`iss`, `sub`, `aud`, `exp`, `nonce`) validation.
3. OAuth 2.0 Authorization Code Grant with single-use code redemption.
4. Refresh Token Rotation with lineage/family tracking and immediate family revocation upon reuse detection.
5. Thread safety and concurrency handling across in-memory server state.

Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo

Primary Risks:
- Race conditions during concurrent token exchanges or refresh requests.
- Incorrect claim checking logic or loose signature validation in OIDC package.
- Inconsistent scope validation or memory leaks in token maps.
