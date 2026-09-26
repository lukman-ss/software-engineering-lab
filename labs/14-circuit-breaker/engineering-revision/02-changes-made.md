## Revision 1

Audit Issue: Lack of dedicated test asserting consecutive failure reset in CLOSED state
Severity: LOW
Files Changed: `internal/circuitbreaker/circuit_breaker_test.go`
Action: Added `TestSuccessInClosedResetsFailures` asserting that non-failing execution in CLOSED state clears the failure accumulator so subsequent failures below threshold do not trigger state transition to OPEN.
Verification: `go test -race ./...` passed.
Status: RESOLVED

## Revision 2

Audit Issue: Lack of test verifying canary concurrency throttling in HALF_OPEN state
Severity: LOW
Files Changed: `internal/circuitbreaker/circuit_breaker_test.go`
Action: Added `TestHalfOpenThrottlesExcessCalls` proving that calls exceeding `HalfOpenMaxCalls` while a canary probe is in-flight immediately fail-fast with `ErrCircuitOpen`.
Verification: `go test -race ./...` passed.
Status: RESOLVED
