# Documentation vs Code

## DOC_CODE_MISMATCH: Demo Output Formatting
Location: `README.md` Expected Behavior vs `cmd/demo/main.go`
Claim:
```text
request=1 result=err=payment failed: status 500 duration=366.75µs state=CLOSED
```
Observation:
Actual output includes a newline and body content:
```text
request=1 result=err=payment failed: status 500 body internal payment server failure
 duration=870.667µs state=CLOSED
```
The README has manually stripped the newline and response body, which misrepresents the raw output of the demo.

## TEST_CLAIM_MISMATCH: Timeout Protection
Location: `README.md` vs `tests/integration_test.go`
Claim: "Checkout service worker threads block on slow HTTP responses... trips open to block calls"
Observation: The test suite never asserts that a slow HTTP response actually trips the breaker. Only `ModeDown` (500 internal server error) is tested in `integration_test.go`.

## RESEARCH_IMPLEMENTATION_MISMATCH
Observation: The implementation aligns well with the research architecture (CLOSED -> OPEN -> HALF_OPEN). The design note appropriately scoped down the sliding window error rate into a consecutive failure counter, so there is no conflict.