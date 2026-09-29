## Revision 1

Audit Issue: GAP-01 — Custom header defense (`RequireCustomHeaderMiddleware`) lacked integration test.
Severity: LOW
Files Changed: `tests/integration_test.go`
Action: Added `TestIntegration_CustomHeader_Protection` verifying missing header returns 403 (balance unchanged) and valid `X-Requested-With: XMLHttpRequest` returns 200 (balance updated).
Verification: `go test -v -run TestIntegration_CustomHeader_Protection ./tests`
Status: RESOLVED

## Revision 2

Audit Issue: GAP-02 — Header-based CSRF token submission (`X-CSRF-Token`) missing from integration tier.
Severity: LOW
Files Changed: `tests/integration_test.go`
Action: Added `TestIntegration_CSRF_Token_In_Header` verifying successful token validation and state transfer when token is supplied via request header.
Verification: `go test -v -run TestIntegration_CSRF_Token_In_Header ./tests`
Status: RESOLVED

## Revision 3

Audit Issue: GAP-03 — Cross-session CSRF token reuse not tested at integration level.
Severity: LOW
Files Changed: `tests/integration_test.go`
Action: Added `TestIntegration_CrossSession_Token_Reuse_Rejected` verifying token issued to victim session is rejected with 403 when presented alongside attacker session cookie.
Verification: `go test -v -run TestIntegration_CrossSession_Token_Reuse_Rejected ./tests`
Status: RESOLVED

## Revision 4

Audit Issue: GAP-04 — `Sec-Fetch-Site: same-origin` positive flow not tested at integration level.
Severity: LOW
Files Changed: `tests/integration_test.go`
Action: Added `TestIntegration_SecFetchSite_SameOrigin_Allowed` verifying requests with `Sec-Fetch-Site: same-origin` successfully execute transfer and return 200.
Verification: `go test -v -run TestIntegration_SecFetchSite_SameOrigin_Allowed ./tests`
Status: RESOLVED
