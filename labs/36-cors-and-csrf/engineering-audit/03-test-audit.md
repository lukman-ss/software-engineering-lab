# Test Audit

## Summary

Test packages audited:
- `internal/cors/middleware_test.go` — 4 tests
- `internal/csrf/token_test.go` — 4 tests (inc. middleware unit)
- `internal/bank/app_test.go` — 1 test
- `tests/integration_test.go` — 5 integration tests

All tests passed with `-count=1` and `-race`:

```
ok  labs/36-cors-and-csrf/internal/bank   0.243s
ok  labs/36-cors-and-csrf/internal/cors   0.243s
ok  labs/36-cors-and-csrf/internal/csrf   0.266s
ok  labs/36-cors-and-csrf/tests           0.244s
```

Race detector:
```
ok  labs/36-cors-and-csrf/internal/bank   1.125s
ok  labs/36-cors-and-csrf/internal/cors   1.121s
ok  labs/36-cors-and-csrf/internal/csrf   1.145s
ok  labs/36-cors-and-csrf/tests           1.137s
```

No race conditions detected.

## Coverage Assessment

### Happy path
- PASS: `TestIntegration_Legitimate_Flow_With_CSRF_Token` — full flow: fetch token, submit protected transfer, balance verified.
- PASS: `TestCORS_Preflight_Success` — allowed origin preflight returns 204 with correct headers.
- PASS: `TestBankTransfer_VulnerableEndpoint` — vulnerable endpoint executes without CSRF protection.

### Failure / attack path
- PASS: `TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution` — disallowed cross-origin request still executes on the vulnerable endpoint.
- PASS: `TestIntegration_CSRF_Token_Prevents_Attack` — protected endpoint returns 403 without CSRF token.
- PASS: `TestCSRFMiddleware_FormPost` — both valid token (200) and missing token (403) tested in same test.

### Edge cases
- PASS: `TestCORS_NoOrigin` — non-CORS requests pass through without headers.
- PASS: `TestCORS_DisallowedOrigin_Preflight` — disallowed preflight returns 403.
- PASS: `TestCORS_Credentials_With_Wildcard_DisallowedInSpec` — wildcard wildcard not reflected when credentials enabled.
- PASS: `TestTokenManager_ExpiredToken` — expired token returns `ErrExpiredToken`.
- PASS: `TestTokenManager_GenerateAndValidate` — session mismatch and tampered token rejected.

### Concurrency
- PASS: `TestIntegration_Concurrency_RaceCondition` — 20 goroutines concurrently request CSRF tokens; no races with `-race`.

## Gaps

### Missing: CSRF token tested only over form field; header submission path not integration-tested
`ValidateToken` reads from `r.Header.Get("X-CSRF-Token")` OR `r.PostFormValue("csrf_token")`. The header path is not exercised in integration tests; only form post is tested. This is minor because the header path is indirectly exercised in unit tests.

### Missing: `RequireCustomHeaderMiddleware` integration test
`TestIntegration` has no test exercising `/api/transfer/custom-header` endpoint from the perspective of an attacker missing the header. Only `FetchMetadataMiddleware` has an integration test (`TestIntegration_SecFetchSite_Protection`). The custom header middleware is covered by its own setup in `setupBankApp` but lacks an integration-level attack/block test.

### Missing: Same-origin Sec-Fetch-Site positive test (integration level)
`TestFetchMetadataMiddleware` in the unit package covers same-origin but the integration test does not repeat this path.

### Missing: Token-reuse / token-binding not tested
No test verifies that a token generated for `session-A` is rejected when submitted with `session-B` cookie in the integration tier (this is tested in unit tests but not in integration).
