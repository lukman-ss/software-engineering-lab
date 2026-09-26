# Gap Analysis

## BROKEN_IMPLEMENTATION: Trailing In-Flight Requests
- Type: BROKEN_IMPLEMENTATION
- Location: `internal/circuitbreaker/circuit_breaker.go` (in `onFailureLocked`)
- Description: Requests initiated during `Closed` state that finish with an error after the breaker has transitioned to `Open` or `HalfOpen` will corrupt the state machine. If `Open`, they reset the `openedAt` cooldown timer. If `HalfOpen`, they act as a false failed probe, prematurely tripping the circuit back to `Open` without a real probe failure.

## BROKEN_IMPLEMENTATION: Panic in Downstream Call
- Type: BROKEN_IMPLEMENTATION
- Location: `internal/circuitbreaker/circuit_breaker.go` (`Execute`)
- Description: If the provided `fn()` panics during the `HalfOpen` state, `b.halfOpenIn` is never decremented or cleared. The breaker gets permanently stuck in `HalfOpen` and rejects all future traffic.

## MISSING_TEST: Timeout / Slow Dependency
- Type: MISSING_TEST
- Location: `tests/integration_test.go`
- Description: The lab's primary motivation is cascading failures caused by slow dependencies, but there are no tests verifying that `ModeSlow` successfully trips the circuit breaker.

## MISSING_TEST: Interleaved State Transitions
- Type: MISSING_TEST
- Location: `internal/circuitbreaker/circuit_breaker_test.go`
- Description: Tests only cover synchronous or fast concurrent failures. There are no tests for slow in-flight requests returning after a state change.

## DOC_CODE_MISMATCH: Demo Output Formatting
- Type: DOC_CODE_MISMATCH
- Location: `README.md`
- Description: The Expected Behavior block omits the response body and newline character present in the actual `cmd/demo` output.