## Revision 1

Audit Issue: `ValidateAccessToken` ignored `requiredScope`
Severity: MEDIUM
Files Changed: `pkg/server/server.go`
Action: Added scope matching check in `ValidateAccessToken` with `containsScope` helper.
Verification: Verified with `TestOAuth2_NegativePaths` scope checks.
Status: RESOLVED

Audit Issue: Access tokens stored with no expiration
Severity: LOW
Files Changed: `pkg/server/server.go`
Action: Created `AccessTokenMeta` struct holding `Subject`, `Scope`, and `ExpiresAt` (1 hour expiration). Enforced expiry check during `ValidateAccessToken`.
Verification: Verified with `go test ./...`.
Status: RESOLVED

Audit Issue: Unhandled `rand.Read` error returns
Severity: LOW
Files Changed: `pkg/server/server.go`, `pkg/client/client.go`
Action: Handled error return values for all `rand.Read` calls.
Verification: Verified with `go test ./...`.
Status: RESOLVED

Audit Issue: Missing negative path test coverage
Severity: MEDIUM
Files Changed: `tests/oauth_test.go`
Action: Added `TestPKCE_Plain_Method`, `TestOIDC_MalformedJWT`, and `TestOAuth2_NegativePaths`.
Verification: `go test -v ./...` passed with 13 total test cases.
Status: RESOLVED
