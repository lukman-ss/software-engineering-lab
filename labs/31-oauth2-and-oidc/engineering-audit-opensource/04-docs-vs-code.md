# Docs vs Code Audit

## README claims vs implementation

| README Claim | Code Reality | Status |
|---|---|---|
| `pkg/pkce`: PKCE code verifier and S256 challenge | `pkce.go` generates S256 OR plain | MATCH (note: plain extra) |
| `pkg/oidc`: ID Token creation, HMAC-SHA256 signature verification, iss/aud/exp/nonce | `oidc.go` HS256, validates iss/aud/exp/iat/nonce | MATCH |
| `pkg/server`: Auth Code grant, PKCE, OIDC ID Tokens, Refresh Token Rotation | `server.go` implements all | MATCH |
| `pkg/client`: Client helper orchestrating requests and token parsing | `client.go` orchestrates build/exchange/refresh + ID token verify | MATCH |
| `cmd/demo`: Executable walkthrough of legitimate flows and attack defenses | `main.go` exactly does this | MATCH |
| README says `PKCE (RFC 7636 / RFC 9700)`: S256 | code accepts `plain` too | DOC_CODE_MISMATCH (minor) |
| README emphasizes HS256 only | matches code | MATCH |

## Run instructions vs actual

| README command | Actually works |
|---|---|
| `go test -v ./...` | PASS |
| `go test -race ./...` | PASS (clean race report) |
| `go run ./cmd/demo` | PASS, exits 0, prints expected steps 1-8 + "Demo Completed Successfully" |

## Research alignment
Key RFC claims (RFC 6749, RFC 7519, RFC 7636, RFC 8252, RFC 9700, OIDC Core) match the implemented flows. No RESEARCH_IMPLEMENTATION_MISMATCH.

## Mismatches found
1. DOC_CODE_MISMATCH: README implies S256-only PKCE; server also accepts deprecated `plain`. See research/04-contradictions.md note on RFC 9700 deprecation.
2. No TEST_CLAIM_MISMATCH — all documented test commands succeed.
3. No README "features" that code omits.

Verdict: documentation is faithful except the minor `plain`-method note (LOW).
