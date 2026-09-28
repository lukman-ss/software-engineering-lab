# Docs vs Code Audit

## Comparison Matrix

| Component / Claim | README Description | Code Implementation | Test Verification | Demo Execution | Status |
|---|---|---|---|---|---|
| OAuth 2.0 Access Token | Delegated authorization grant with scope checks | `pkg/server` (`ValidateAccessToken`) | `TestOAuth2_NegativePaths` | Step 5 in `cmd/demo` | MATCH |
| OIDC ID Token | HMAC-SHA256 JWT with `iss`, `sub`, `aud`, `exp`, `nonce` | `pkg/oidc` (`SignIDToken`, `ParseAndVerifyIDToken`) | `TestOIDC_IDToken_*` | Step 4 in `cmd/demo` | MATCH |
| PKCE Defense | RFC 7636 / RFC 9700 `S256` and `plain` methods | `pkg/pkce` (`ComputeChallenge`, `Verify`) | `TestPKCE_*`, `TestOAuth2_FullFlowAndPKCEInterception` | Step 1, 6 in `cmd/demo` | MATCH |
| Refresh Token Rotation | RFC 9700 Sec 4.14 family tracking & reuse revocation | `pkg/server` (`Refresh`) | `TestOAuth2_RefreshTokenRotation_AndReplayDetection` | Step 7, 8 in `cmd/demo` | MATCH |
| Running Commands | `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | Valid go module and package layout | Passes identically as described | Runs cleanly with exact step outputs | MATCH |

## Discrepancies Found

None. The documentation accurately reflects all implemented capabilities, parameter names, and verification workflows.
