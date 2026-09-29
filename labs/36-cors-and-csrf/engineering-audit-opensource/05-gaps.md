# Gap Analysis

## GAP-001

Type: MISSING_TEST
Severity: MEDIUM
Location: `internal/csrf/middleware.go:73-85` — `RequireCustomHeaderMiddleware`
Description: `RequireCustomHeaderMiddleware` is fully implemented and wired in `tests/integration_test.go:48` as the `/api/transfer/custom-header` route, but no test case actually sends requests to that route to verify the middleware behavior. Neither a request with the correct `X-Requested-With: XMLHttpRequest` header nor a request missing the header is exercised.
Impact: Custom header defense pathway is unvalidated by any test. A regression in this middleware would not be detected.

---

## GAP-002

Type: MISSING_EDGE_CASE
Severity: LOW
Location: `internal/csrf/token.go:72-74`
Description: Token validation splits on `:` without bounds restriction (`strings.Split`, not `strings.SplitN`). If a `sessionID` contains a `:` character, the token is malformed and fails validation with `ErrInvalidToken`. No test exists for session IDs containing colons.
Impact: Latent structural incompatibility if session ID format changes to include colons. Not exploitable in current codebase.

---

## GAP-003

Type: MISSING_EDGE_CASE
Severity: LOW
Location: `internal/bank/app.go:117-120` — transfer amount validation
Description: `HandleTransferVulnerable` returns 400 Bad Request for `amount <= 0`, but no test validates this constraint. The negative-amount path through the bank's transfer handler is exercised by neither unit nor integration tests.
Impact: Minor test completeness issue. Implementation is correct.

---

## GAP-004

Type: MISSING_TEST
Severity: LOW
Location: `internal/cors/middleware_test.go` — disallowed origin simple request (non-OPTIONS)
Description: The CORS tests cover: no origin, disallowed origin preflight (OPTIONS), allowed origin preflight, and credentials/wildcard. There is no test for a disallowed-origin non-preflight simple request (e.g., GET or POST from an unknown origin). The code handles it correctly (passes to next handler, no CORS headers), but this behavior is untested directly in the cors package unit tests.
Impact: Low — covered transitively by `TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution`, which exercises the same path end-to-end.

---

## Summary Table

| Gap ID | Type | Severity | Blocking |
| ------ | ---- | -------- | -------- |
| GAP-001 | MISSING_TEST | MEDIUM | No |
| GAP-002 | MISSING_EDGE_CASE | LOW | No |
| GAP-003 | MISSING_EDGE_CASE | LOW | No |
| GAP-004 | MISSING_TEST | LOW | No |

No HIGH or CRITICAL issues identified.
