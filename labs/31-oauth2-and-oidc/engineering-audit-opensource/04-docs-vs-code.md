# Documentation vs Code Audit

## 1. Documentation Review

Document: `README.md`
Claims Checked:
1. **OAuth 2.0 (Authorization)**: Access Tokens represent delegated authorization grants to protected resource servers.
   - Verified in `pkg/server/server.go:288-306` (`ValidateAccessToken` checking scope and expiration).
2. **OIDC (Authentication)**: ID Tokens (signed JWTs) provide claims (`iss`, `sub`, `aud`, `exp`, `nonce`).
   - Verified in `pkg/oidc/oidc.go:29-38` and `pkg/server/server.go:198-215`.
3. **PKCE (RFC 7636 / RFC 9700)**: S256 challenge generation and code interception mitigation.
   - Verified in `pkg/pkce/pkce.go` and `tests/oauth_test.go:TestOAuth2_FullFlowAndPKCEInterception`.
4. **Refresh Token Rotation (RFC 9700 Section 4.14)**: Single-use refresh tokens with lineage tracking and token family revocation upon reuse.
   - Verified in `pkg/server/server.go:219-286` and `tests/oauth_test.go:TestOAuth2_RefreshTokenRotation_AndReplayDetection`.

## 2. Discrepancy Analysis

- `DOC_CODE_MISMATCH`: None found.
- `TEST_CLAIM_MISMATCH`: None found.
- `RESEARCH_IMPLEMENTATION_MISMATCH`: None found.

## 3. Demo Alignment

- Executable path: `cmd/demo/main.go`
- Demo steps match all architectural components described in the README (Steps 1–8 covering PKCE setup, Auth Code issuance, Token exchange, OIDC ID Token validation, Resource access, PKCE attack blocking, Refresh Token rotation, and Replay detection family revocation).
