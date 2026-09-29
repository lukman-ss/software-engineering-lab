# Engineering Audit Verdict

Target Lab: `labs/36-cors-and-csrf`
Audit Date: 2026-09-29

## Summary

Code Files Reviewed: 5
- `internal/cors/middleware.go`
- `internal/csrf/middleware.go`
- `internal/csrf/token.go`
- `internal/bank/app.go`
- `cmd/demo/main.go`

Tests Reviewed: 4
- `internal/cors/middleware_test.go`
- `internal/csrf/token_test.go`
- `internal/bank/app_test.go`
- `tests/integration_test.go`

Commands Executed:
- `rtk go test -v ./...` (14 passed)
- `rtk go test -race ./...` (14 passed, 0 races)
- `rtk go run ./cmd/demo` (success, clean output matching claims)

Failures: 0
Warnings: 4 (1 Medium, 3 Low)

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS
Demo: PASS
Research Alignment: PASS
Documentation Accuracy: PASS

## Blocking Issues

None.

## Non-Blocking Issues

1. **GAP-001 (MEDIUM)**: Missing test for `RequireCustomHeaderMiddleware` at `/api/transfer/custom-header`.
2. **GAP-002 (LOW)**: `strings.Split` in token validation breaks if `sessionID` contains `:` delimiters.
3. **GAP-003 (LOW)**: No test verifying negative/zero amount validation in `HandleTransferVulnerable`.
4. **GAP-004 (LOW)**: Disallowed origin simple request is not tested at the unit level in `internal/cors/middleware_test.go` (covered transitively in integration test).

## Required Revisions

1. Add a test case for `RequireCustomHeaderMiddleware` in `internal/csrf/token_test.go` or `tests/integration_test.go` verifying that requests lacking `X-Requested-With` are rejected with 403 and requests containing the header pass.
2. Consider switching `strings.Split` to `strings.SplitN(..., 4)` or JSON payload in `TokenManager.ValidateToken` to prevent potential colon collisions in session identifiers.

## Final Status

APPROVED_WITH_WARNINGS
