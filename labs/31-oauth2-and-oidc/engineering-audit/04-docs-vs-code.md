# Docs vs Code Audit

Target Lab: `labs/31-oauth2-and-oidc`

## Comparison Matrix

| Item | Documented Claim | Code / Test Reality | Match Status |
| --- | --- | --- | --- |
| Zero-dependency | Pure Go standard library | `go.mod` has 0 third-party dependencies | PASS |
| PKCE Method | `S256` default, RFC 7636 | Implemented in `pkg/pkce`, tested in `tests/` | PASS |
| ID Token Verification | HMAC-SHA256 JWT, claims verification (`iss`, `sub`, `aud`, `exp`, `nonce`) | Implemented in `pkg/oidc/oidc.go`, tested in `tests/` | PASS |
| Refresh Token Rotation | Token family lineage, reuse detection revokes entire family | Implemented in `pkg/server/server.go:220-286`, tested in `tests/` | PASS |
| Commands in README | `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | All 3 commands execute without error | PASS |
| Demo Output | Shows step 1 through step 8 walkthrough | Real terminal execution matches demo logs exactly | PASS |

## Identified Mismatches

None.
- No `DOC_CODE_MISMATCH`
- No `TEST_CLAIM_MISMATCH`
- No `RESEARCH_IMPLEMENTATION_MISMATCH`
