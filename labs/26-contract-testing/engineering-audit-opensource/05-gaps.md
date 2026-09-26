# Gaps

Target Lab: labs/26-contract-testing

## Identified Gaps

GAP-01 [MISSING_TEST]
internal/contract/verifier.go diffValues slice recursion
Claimed Behavior: general JSON comparison
Observed Implementation: no `[]interface{}` recursion; would panic on array diff
Severity: LOW
Status: Known, scoped, accepted (revision plan documents; no refactor needed, single-resource contract uses no arrays).

GAP-02 [MISSING_TEST]
internal/contract/verifier.go response headers
Claimed Behavior: engineering 02-implementation-notes item 3 / 01-design components mentions "headers" in verification
Observed Implementation: response headers in contract ignored
Severity: LOW
Status: Not core claim (status + body subset). Unenforced.

GAP-03 [MISSING_TEST]
internal/contract/verifier.go decoder error path / malformed body
Claimed Behavior: handles non-JSON provider responses
Observed Implementation: error path at verifier.go:98 exists (continue) but no test exercises malformed body or non-200 status
Severity: LOW
Status: Branch exists, untested.

GAP-04 [MISSING_TEST]
internal/provider/server.go V2 endpoint
Claimed Behavior: engineering 01-design:18 "V1 contract passes against /v1, while V2 contract passes against /v2"
Observed Implementation: V2 handler present; no V2 contract; no test hits /v2/orders/
Severity: LOW
Status: V2 code is trivial GET returning V2 DTO; claim partially unproven. V1 preservation proven.

GAP-05 [MISSING_TEST]
negative status / wrong path
Observed Implementation: no test asserting non-200 or missing 404 path triggers verifier failure
Severity: LOW
Status: Would exercise verifier status check (verifier.go:81), untested.

GAP-06 [UNHANDLED_ERROR]
internal/contract/verifier.go:51-55 client timeout
Claimed Behavior: executes contract verification in CI gate
Observed Implementation: http.Client{} no Timeout -> indefinite hang on dead provider
Severity: LOW
Status: httptest in lab is reliable; production gate would need timeout. Out of scope.

## Summary By Severity

CRITICAL: 0
HIGH: 0
MEDIUM: 0
LOW: 6 (all scoped; core behavior proven; no blockers)
