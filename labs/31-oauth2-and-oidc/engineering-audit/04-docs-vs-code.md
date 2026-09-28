# Docs vs Code Audit

## Documentation Review

### Files Examined
1. `README.md`
2. `engineering/01-design.md`
3. `engineering/02-implementation-notes.md`
4. `engineering/03-execution-result.md`
5. Implementation packages (`pkg/pkce`, `pkg/oidc`, `pkg/server`, `pkg/client`)
6. Executable Demo (`cmd/demo/main.go`)

### Comparison Analysis

1. **OAuth 2.0 vs OIDC Separation**:
   - Docs Claim: OAuth 2.0 provides delegated access tokens for resource servers, while OIDC delivers ID Tokens (JWT) asserting identity claims.
   - Code Verification: `pkg/server/server.go` issues `AccessToken` and optionally `IDToken` if `openid` scope is present. `pkg/oidc` handles ID token validation.

2. **PKCE Protection**:
   - Docs Claim: Uses RFC 7636 / RFC 9700 S256 challenge creation and verification to protect authorization codes from interception.
   - Code Verification: Implemented in `pkg/pkce/pkce.go` and enforced in `server.ExchangeCode`.

3. **Refresh Token Rotation & Family Revocation**:
   - Docs Claim: Complies with RFC 9700 §4.14; detects reuse and revokes entire token family.
   - Code Verification: `server.Refresh` checks `meta.Revoked` and marks `revokedFams[meta.FamilyID] = true`.

4. **Demo Walkthrough**:
   - Docs Claim: `cmd/demo` executes full 8-step flow demonstrating legitimate issuance, ID token verification, resource access, interception defense, token rotation, and replay detection.
   - Execution Verification: Demo output matches documented execution trace exactly.

### Discrepancies Found
- None. All claimed features, mechanisms, and architectural designs in `README.md` and `engineering/` accurately reflect the implementation.
