# Test Audit

## Overview

Test coverage spans unit tests for each subsystem (`internal/cors`, `internal/csrf`, `internal/bank`) and end-to-end integration tests (`tests/integration_test.go`).

## Execution Results

Command:
```bash
go test -count=1 -race -v ./...
```

Output:
```text
?   	labs/36-cors-and-csrf/cmd/demo	[no test files]
=== RUN   TestBankTransfer_VulnerableEndpoint
--- PASS: TestBankTransfer_VulnerableEndpoint (0.00s)
PASS
ok  	labs/36-cors-and-csrf/internal/bank	1.205s
=== RUN   TestCORS_NoOrigin
--- PASS: TestCORS_NoOrigin (0.00s)
=== RUN   TestCORS_DisallowedOrigin_Preflight
--- PASS: TestCORS_DisallowedOrigin_Preflight (0.00s)
=== RUN   TestCORS_Preflight_Success
--- PASS: TestCORS_Preflight_Success (0.00s)
=== RUN   TestCORS_Credentials_With_Wildcard_DisallowedInSpec
--- PASS: TestCORS_Credentials_With_Wildcard_DisallowedInSpec (0.00s)
PASS
ok  	labs/36-cors-and-csrf/internal/cors	1.199s
=== RUN   TestTokenManager_GenerateAndValidate
--- PASS: TestTokenManager_GenerateAndValidate (0.00s)
=== RUN   TestTokenManager_ExpiredToken
--- PASS: TestTokenManager_ExpiredToken (0.02s)
=== RUN   TestCSRFMiddleware_FormPost
--- PASS: TestCSRFMiddleware_FormPost (0.00s)
=== RUN   TestFetchMetadataMiddleware
--- PASS: TestFetchMetadataMiddleware (0.00s)
PASS
ok  	labs/36-cors-and-csrf/internal/csrf	1.232s
=== RUN   TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution
--- PASS: TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution (0.00s)
=== RUN   TestIntegration_CSRF_Token_Prevents_Attack
--- PASS: TestIntegration_CSRF_Token_Prevents_Attack (0.00s)
=== RUN   TestIntegration_Legitimate_Flow_With_CSRF_Token
--- PASS: TestIntegration_Legitimate_Flow_With_CSRF_Token (0.00s)
=== RUN   TestIntegration_SecFetchSite_Protection
--- PASS: TestIntegration_SecFetchSite_Protection (0.00s)
=== RUN   TestIntegration_Concurrency_RaceCondition
--- PASS: TestIntegration_Concurrency_RaceCondition (0.00s)
PASS
ok  	labs/36-cors-and-csrf/tests	1.202s
```

## Coverage Analysis
- Happy Path: Verified for legitimate transfers with signed token (`TestIntegration_Legitimate_Flow_With_CSRF_Token`).
- Failure / Attack Path: Verified that cross-origin POST executes on vulnerable endpoint (`TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution`) and is blocked on protected endpoint (`TestIntegration_CSRF_Token_Prevents_Attack`).
- Preflight & CORS Spec: Verified allowed vs disallowed origin preflight responses (`TestCORS_DisallowedOrigin_Preflight`, `TestCORS_Preflight_Success`, `TestCORS_Credentials_With_Wildcard_DisallowedInSpec`).
- Expiration & Tampering: Verified in `TestTokenManager_ExpiredToken` and `TestTokenManager_GenerateAndValidate`.
- Concurrency & Race Detector: Verified 20 concurrent workers generating tokens and hitting endpoints under `-race` flag without data races.
