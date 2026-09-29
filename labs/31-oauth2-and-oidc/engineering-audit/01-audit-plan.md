# Engineering Audit Plan

Target Lab: labs/31-oauth2-and-oidc

Implementation Files:
- pkg/pkce/pkce.go
- pkg/oidc/oidc.go
- pkg/server/server.go
- pkg/client/client.go
- cmd/demo/main.go

Tests:
- tests/oauth_test.go (17 test functions)

Executable/Demo:
- cmd/demo/main.go (go run ./cmd/demo)

Approved Research Inputs:
- research/ (OAuth 2.0, OIDC, PKCE RFC 7636/9700, Refresh Token Rotation RFC 9700 §4.14)
- research-audit/07-verdict.md
- engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md
- engineering-revision/

Main Claims To Verify:
1. PKCE S256 challenge generation and verification (RFC 7636)
2. PKCE plain method support
3. PKCE verifier length bounds (43–128)
4. Authorization code single-use enforcement
5. Authorization code expiry enforcement
6. OIDC ID Token signing (HMAC-SHA256) and verification
7. ID Token claims validation: iss, aud, exp, iat, nonce
8. Refresh Token Rotation (single-use, family lineage)
9. Refresh Token replay detection + full family revocation
10. Concurrency safety (sync.Mutex on AuthorizationServer)
11. Demo walkthrough is real (not fabricated)

Commands To Run:
```bash
cd labs/31-oauth2-and-oidc
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Race conditions in concurrent refresh rotation
- PKCE verification using non-constant-time comparison
- Family revocation ordering issue (revoked check before family check)
- Demo output fabricated (must verify against real execution)
