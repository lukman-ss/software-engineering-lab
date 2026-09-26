# Engineering Revision Plan

Target Lab: labs/14-circuit-breaker
Previous Verdict: NEEDS_REVISION

## Blocking Issues
- **Trailing In-Flight Request Corruption (HIGH)**: Late responses from a previous generation/state corrupt `Open` cooldown timer or prematurely fail `HalfOpen`.
- **Panic Induced Stuck State (HIGH)**: If `fn()` panics during `HalfOpen`, probe counter `b.halfOpenIn` is never decremented, causing permanent lockup in `HalfOpen`.

## Non-Blocking Issues
- **Missing Timeout Test (MEDIUM)**: Test suite does not test that `ModeSlow` (timeouts) trips breaker.
- **Missing Interleaved Concurrency Test (MEDIUM)**: Concurrency test does not verify in-flight requests crossing state transitions.
- **Demo Output Formatting (LOW)**: Output formatting newline in `fake_server.go` and README mismatch; Scenario 1 printing `CLOSED` without breaker.

## Files To Change
- `labs/14-circuit-breaker/internal/circuitbreaker/circuit_breaker.go`: Add generation tracking to isolate state transitions; handle panics with `defer` to safely clear probe counters.
- `labs/14-circuit-breaker/internal/checkout/service.go`: Fix `CheckoutWithoutBreaker` state reporting (use state string/indicator or distinct representation).
- `labs/14-circuit-breaker/internal/payment/fake_server.go`: Strip trailing newline or format cleaner errors without confusing formatting.
- `labs/14-circuit-breaker/README.md`: Align demo output.

## Tests To Add/Modify
- `labs/14-circuit-breaker/internal/circuitbreaker/circuit_breaker_test.go`:
  - Test trailing request from previous generation does not reset cooldown or trip HalfOpen.
  - Test panic in `fn()` during `HalfOpen` does not permanently lock breaker.
  - Add interleaved concurrent transition test.
- `labs/14-circuit-breaker/tests/integration_test.go`:
  - Add test verifying `ModeSlow` trips the breaker via timeout.

## Validation Commands
```bash
cd labs/14-circuit-breaker
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
