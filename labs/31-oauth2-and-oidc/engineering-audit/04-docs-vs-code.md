# Docs vs Code Audit

## Comparison Matrix

| Component / Claim | README & Design Claim | Code Implementation | Test Verification | Verdict |
|---|---|---|---|---|
| OAuth 2.0 Auth Code Grant | Server issues Auth Code bound to Client & Redirect URI | `pkg/server/server.go:92-130` | `TestOAuth2_NegativePaths` & `TestOAuth2_FullFlowAndPKCEInterception` | MATCH |
| PKCE S256 & Plain | Mandates PKCE, validates verifier vs challenge | `pkg/pkce/pkce.go` | `TestPKCE_S256_Valid`, `TestPKCE_Plain_Method`, `TestPKCE_Mismatch` | MATCH |
| OIDC ID Token Verification | Signs JWT with HMAC-SHA256, verifies claims (`iss`, `sub`, `aud`, `exp`, `nonce`) | `pkg/oidc/oidc.go` | `TestOIDC_IDToken_Valid`, `TestOIDC_IDToken_Expired`, `TestOIDC_IDToken_MismatchClaims`, `TestOIDC_IDToken_TamperedSignature` | MATCH |
| Refresh Token Rotation | Single-use refresh token, lineage tracking | `pkg/server/server.go:219-286` | `TestOAuth2_RefreshTokenRotation_AndReplayDetection` | MATCH |
| Replay Detection Revocation | Replaying consumed refresh token revokes entire family | `pkg/server/server.go:234-238` | `TestOAuth2_RefreshTokenRotation_AndReplayDetection` | MATCH |
| Scope Enforcement | Resource server checks granted scopes | `pkg/server/server.go:288-321` | `TestOAuth2_NegativePaths` | MATCH |
| Demo Execution | `go run ./cmd/demo` executes end-to-end walkthrough | `cmd/demo/main.go` | Executed live with 0 exit code | MATCH |

## Discrepancies Found
- None. Documentation, code signatures, and runnable demo commands correspond exactly with no discrepancies.
