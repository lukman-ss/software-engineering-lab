# Docs vs Code

Source: README.md, engineering/01-design.md, engineering/02-implementation-notes.md, pkg/*, cmd/demo/main.go, tests/oauth_test.go.

## Match Analysis

| Artifact | Claim | Status |
|---|---|---|
| README structure (`pkg/pkce`, `pkg/oidc`, `pkg/server`, `pkg/client`, `cmd/demo`, `tests`) | Matches | PASS |
| README run commands (`go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`) | Verified real | PASS |
| Design "Resource server accepts Access Token and checks scopes" | `ValidateAccessToken` ignores `requiredScope` | FAIL (DOC_CODE_MISMATCH) |
| Design test strategy (`pkce_test.go`, `oidc_test.go`, `server_test.go`, `race_test.go`) | All tests in `tests/oauth_test.go` | WARNING (DOC_CODE_MISMATCH) |
| Implementation notes "Bearer Access Token validation" | Only token existence checked — scope & expiry not enforced | WARNING |
| "ID Token verification (`iss`, `aud`, `exp`, `nonce`)" | All four enforced in `ParseAndVerifyIDToken` | PASS |
| "Refresh Token Rotation with family revocation" | Implemented & verified | PASS |
| "100% Go standard library, no external deps" | go.mod has no requires | PASS |

## Mismatches

1. **DOC_CODE_MISMATCH**: Scope checking explicitly described as Resource Server behavior. Parameter exists in signature but is never used.
2. **DOC_CODE_MISMATCH**: Test file names in design doc do not match actual file layout.
3. **IMPLEMENTATION_OVERCLAIM**: `ExpiresIn` returned to client (3600s) but access tokens never expire server‑side — `s.tokens` map has no expiry tracking; `ValidateAccessToken` only checks existence.
