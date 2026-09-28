# Engineering Audit Verdict

Target Lab: labs/31-oauth2-and-oidc
Audit Date: 2026-09-28

## Summary

Code Files Reviewed: 5 (`pkg/pkce/pkce.go`, `pkg/oidc/oidc.go`, `pkg/server/server.go`, `pkg/client/client.go`, `cmd/demo/main.go`)
Tests Reviewed: 1 (`tests/oauth_test.go`, 10 tests)
Commands Executed: 5 (`go build ./...`, `go test -v -count=1 ./...`, `go test -race -count=1 ./...`, `go run ./cmd/demo`, `go vet ./...`)
Failures: 0 command failures
Warnings: 2 MEDIUM, 3 LOW

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS
Demo: PASS
Research Alignment: NOT_APPLICABLE (skipped per pipeline override — implementation/tests only)
Documentation Accuracy: WARNING

## Blocking Issues
None. No HIGH/CRITICAL issues. Core claimed behaviors (PKCE S256 interception defense, ID Token sign/verify with iss/aud/exp/nonce, auth-code single-use, refresh rotation with family revocation, concurrency safety) are implemented and proven by passing tests + live demo.

## Non-Blocking Issues
1. [MEDIUM] `ValidateAccessToken` ignores `requiredScope` — scope enforcement claimed in design/README but unimplemented (`pkg/server/server.go:268-278`). DOC_CODE_MISMATCH / IMPLEMENTATION_OVERCLAIM.
2. [MEDIUM] Missing negative-path tests: malformed JWT, expired code, bad client/redirect, missing/invalid challenge, refresh expiry/not-found/wrong-client, plain PKCE, scope rejection (`tests/oauth_test.go`). MISSING_TEST.
3. [LOW] `ExpiresIn: 3600` returned but access tokens stored with no expiry; `ValidateAccessToken` never expires them (`pkg/server/server.go:176-182,255-257`). IMPLEMENTATION_OVERCLAIM.
4. [LOW] `crypto/rand.Read` return values discarded in 8 places (`pkg/client/client.go:44,47`, `pkg/server/server.go:153,160,164,240,244,254`). UNHANDLED_ERROR.
5. [LOW] Design doc test filenames (`pkce_test.go`, `oidc_test.go`, `server_test.go`, `race_test.go`) do not match actual `tests/oauth_test.go` (`engineering/01-design.md:80-85`). DOC_CODE_MISMATCH.

## Required Revisions
1. (Recommended, non-blocking) Enforce `requiredScope` in `ValidateAccessToken` or remove the parameter and correct README/design claims.
2. (Recommended, non-blocking) Add expiry timestamp to access tokens and enforce it, or document that access tokens are non-expiring in this lab.
3. (Recommended, non-blocking) Add negative-path tests for the branches listed in `05-gaps.md`; check `rand.Read` errors.

## Final Status

APPROVED_WITH_WARNINGS
