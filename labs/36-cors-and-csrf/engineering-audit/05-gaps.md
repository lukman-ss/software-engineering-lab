# Gap Analysis

## Identified Gaps

### GAP-01

Type: MISSING_TEST
Severity: LOW
Location: `internal/cors/middleware.go:44-58`
Description: There is no unit test verifying that a non-preflight (e.g., POST) request from a disallowed origin still reaches the backend handler (i.e., returns 200 OK and mutates state). The integration test `TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution` covers this path, but no unit-level test in `internal/cors/middleware_test.go` explicitly tests that disallowed cross-origin simple requests pass through to `next`.
Impact: Low. Integration tests cover this critical claim. Not a blocking issue.

---

### GAP-02

Type: MISSING_EDGE_CASE
Severity: LOW
Location: `internal/csrf/token.go:62-109`
Description: The token format `sessionID:timestamp:nonce:hmac` uses `:` as a delimiter. If `sessionID` contains a `:` character, `strings.Split(string(decoded), ":")` would produce more than 4 parts and `len(parts) != 4` check would reject valid tokens. There is no test or input validation rejecting session IDs containing colons.
Impact: Low. Session IDs in this lab are hardcoded (`session-victim-secret`, `session-abc`) and do not contain colons. Would require enforcement if this code is extended beyond lab scope.

---

### GAP-03

Type: MISSING_TEST
Severity: LOW
Location: `internal/bank/app.go:101-146`
Description: `HandleTransferVulnerable` and `HandleTransferProtected` do not have explicit tests for `amount <= 0`, `recipient not found`, and `insufficient balance` error paths.
Impact: Low. Error paths return 400 Bad Request. Not critical for security claims; the lab focuses on CSRF demonstrations.

---

## No Critical Gaps Found

The following gap types were checked and **not found**:

- `BROKEN_IMPLEMENTATION`: All implementations behave as claimed.
- `DOC_CODE_MISMATCH`: No mismatches found between README, engineering notes, and code.
- `RACE_CONDITION`: `-race` detector clean; no data races found.
- `UNHANDLED_ERROR`: Error paths are propagated and return appropriate HTTP status codes.
- `IMPLEMENTATION_OVERCLAIM`: All claims are supported by test evidence and execution results.
- `RESEARCH_MISMATCH`: Research findings align with implementation behavior.
- `FAKE_DEMO`: Demo output verified as real — produced from actual code execution with account balance state changes.
- `FAKE_BENCHMARK`: No benchmark claims exist in this lab.
- `UNVERIFIED_RESULT`: All stated results are verified by test suite execution.
