# Gap Analysis

Target Lab: labs/31-oauth2-and-oidc
Gap types used: MISSING_TEST, DOC_CODE_MISMATCH, IMPLEMENTATION_OVERCLAIM, UNHANDLED_ERROR.

## Gap List

1. MISSING_TEST: No test for `ValidateAccessToken` scope enforcement (parameter ignored).
   - File: pkg/server/server.go
   - Lines: 268-278

2. MISSING_TEST: No test for malformed JWT (wrong number of parts) in `ParseAndVerifyIDToken`.
   - File: pkg/oidc/oidc.go
   - Lines: 65-67

3. MISSING_TEST: No test for `ErrIssuedInFuture` (token issued too far in future).
   - File: pkg/oidc/oidc.go
   - Lines: 107

4. MISSING_TEST: No test for unauthorized client on Refresh (clientID mismatch).
   - File: pkg/server/server.go
   - Lines: 227-229

5. MISSING_TEST: No test for expired refresh token.
   - File: pkg/server/server.go
   - Lines: 231-233

6. MISSING_TEST: No test for `Refresh` returning `ErrRefreshTokenNotFound`.
   - File: pkg/server/server.go
   - Lines: 210-213

7. MISSING_TEST: No test for `Authorize` with missing `code_challenge`.
   - File: pkg/server/server.go
   - Lines: 97-99

8. MISSING_TEST: No test for `Authorize` with invalid `code_challenge_method`.
   - File: pkg/server/server.go
   - Lines: 101-103

9. MISSING_TEST: No test for `Authorize` with unauthorized client/redirectURI.
   - File: pkg/server/server.go
   - Lines: 92-95

10. MISSING_TEST: No test for `ExchangeCode` with expired auth code.
    - File: pkg/server/server.go
    - Lines: 133-135

11. MISSING_TEST: No test for `ExchangeCode` with wrong clientID/redirectURI.
    - File: pkg/server/server.go
    - Lines: 141-143

12. MISSING_TEST: No test for `BuildAuthorizationRequest` generating distinct state/nonce each call.
    - File: pkg/client/client.go
    - Lines: 44-49

13. MISSING_TEST: No test for PKCE "plain" method generation/verification.
    - File: pkg/pkce/pkce.go
    - Lines: 24-25, 56-57

14. MISSING_TEST: No test for `RefreshTokens` when client has never exchanged (empty refresh token).
    - File: pkg/client/client.go
    - Lines: 74-75

15. DOC_CODE_MISMATCH: `ValidateAccessToken` claims scope enforcement via `requiredScope` param but ignores it.
    - File: pkg/server/server.go
    - Lines: 268-278

16. DOC_CODE_MISMATCH: Design doc references test files `pkce_test.go`, `oidc_test.go`, `server_test.go`, `race_test.go`; actual file is `tests/oauth_test.go`.
    - File: engineering/01-design.md
    - Lines: 80-85

17. IMPLEMENTATION_OVERCLAIM: `ExpiresIn: 3600` returned in TokenResponse, but server stores access tokens with no expiry; `ValidateAccessToken` cannot expire them.
    - File: pkg/server/server.go
    - Lines: 176-182 (AccessToken creation), 255-257 (Refresh token creation)

18. UNHANDLED_ERROR: Ignored return value from `crypto/rand.Read` (8 instances).
    - Files: pkg/client/client.go (44,47); pkg/server/server.go (153,160,164,240,244,254)
    - Error: silent discard of `(n, err)` from `rand.Read()`