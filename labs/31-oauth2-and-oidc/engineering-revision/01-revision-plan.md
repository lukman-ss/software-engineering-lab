# Engineering Revision Plan

Target Lab: labs/31-oauth2-and-oidc
Previous Verdict: APPROVED_WITH_WARNINGS

## Blocking Issues
None.

## Non-Blocking Issues
1. `ValidateAccessToken` ignored `requiredScope` parameter, leaving claimed scope validation unenforced (`pkg/server/server.go`).
2. Access tokens lacked expiration timestamps despite token responses returning `ExpiresIn: 3600`.
3. Unhandled error returns from `crypto/rand.Read` across `pkg/server/server.go` and `pkg/client/client.go`.
4. Missing test coverage for negative branches (malformed JWT, plain PKCE, scope rejection, bad client/redirect/challenge, invalid refresh).

## Files To Change
- `pkg/server/server.go`: Add `AccessTokenMeta` with `ExpiresAt` and `Scope`; implement scope checking and expiration validation in `ValidateAccessToken`; check errors on `rand.Read`.
- `pkg/client/client.go`: Check errors on `rand.Read` in `BuildAuthorizationRequest`.
- `tests/oauth_test.go`: Add `TestPKCE_Plain_Method`, `TestOIDC_MalformedJWT`, and `TestOAuth2_NegativePaths`.

## Tests To Add/Modify
- `TestPKCE_Plain_Method`
- `TestOIDC_MalformedJWT`
- `TestOAuth2_NegativePaths`

## Validation Commands
```bash
go test ./...
go test -race ./...
go run ./cmd/demo
```
