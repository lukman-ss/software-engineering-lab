# Test Audit

Target Lab: `labs/31-oauth2-and-oidc`

## Test Coverage Summary

File: `tests/oauth_test.go`
Total Test Functions: 13

### Covered Paths & Categories

1. **PKCE Validation**:
   - `TestPKCE_S256_Valid`: Verifies `S256` verifier generation and verification pass.
   - `TestPKCE_Plain_Method`: Verifies `plain` method support.
   - `TestPKCE_InvalidMethod`: Verifies rejection of unknown PKCE methods.
   - `TestPKCE_Mismatch`: Verifies rejection when code verifier does not match challenge.

2. **OIDC ID Token Validation**:
   - `TestOIDC_IDToken_Valid`: Validates happy path ID Token parsing and claims.
   - `TestOIDC_IDToken_TamperedSignature`: Verifies rejection of tampered JWT signatures.
   - `TestOIDC_IDToken_Expired`: Verifies rejection of expired ID Tokens.
   - `TestOIDC_IDToken_MismatchClaims`: Verifies rejection on `iss`, `aud`, or `nonce` mismatches.
   - `TestOIDC_MalformedJWT`: Verifies handling of non-base64 and malformed 2-part tokens.

3. **OAuth 2.0 Authorization Server & Flows**:
   - `TestOAuth2_FullFlowAndPKCEInterception`: Verifies full Auth Code grant, PKCE interception defense failure on wrong verifier, legitimate exchange, and auth code single-use constraint.
   - `TestOAuth2_RefreshTokenRotation_AndReplayDetection`: Verifies token rotation, replay detection on consumed token, and family revocation blocking subsequent active session refreshes.
   - `TestOAuth2_NegativePaths`: Tests unregistered clients, invalid redirect URIs, empty challenges, invalid exchange codes, mismatched client IDs, missing/excess scopes, and nonexistent refresh tokens.

4. **Concurrency & Thread Safety**:
   - `TestOAuth2_ConcurrencyAndRace`: Launches 20 concurrent goroutines performing authorization request, code exchange, access token validation, and token refresh under Go race detector (`go test -race ./...`).

## Execution Command Outputs

```bash
go test ./...
# Output: ok labs/31-oauth2-and-oidc/tests 0.123s

go test -race ./...
# Output: ok labs/31-oauth2-and-oidc/tests 1.149s
```

## Test Audit Verdict

PASS — Test suite provides comprehensive coverage for happy paths, negative error cases, security attack scenarios (interception & replay), malformed inputs, and concurrency race safety.
