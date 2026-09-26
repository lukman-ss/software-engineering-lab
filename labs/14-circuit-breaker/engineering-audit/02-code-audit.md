# Code Audit

## Finding 1

Location: `internal/circuitbreaker/circuit_breaker.go:88` (`onFailureLocked`)
Claimed Behavior: Breaker transitions cleanly between states and honors the cooldown timer.
Observed Implementation: `onFailureLocked` does not check if the request that failed was initiated in the current state. If a long-running request starts during `Closed`, takes 500ms, and fails, but the breaker tripped to `Open` at 100ms, the trailing request calls `onFailureLocked(now)` while the state is `Open` or `HalfOpen`. If `Open`, it overwrites `b.openedAt = now`, completely resetting the cooldown timer. If `HalfOpen`, it forces the state to `Open`, treating the old trailing request as a failed probe.
Assessment: FAIL
Severity: HIGH
Notes: Requests must carry a "generation" or the state machine must ignore failures from requests that started before the last state transition.

## Finding 2

Location: `internal/circuitbreaker/circuit_breaker.go:49` (`Execute`)
Claimed Behavior: Circuit breaker handles state robustly and recovers via probe mechanism.
Observed Implementation: When state is `HalfOpen`, `b.halfOpenIn` is incremented. `onSuccessLocked` or `onFailureLocked` are called after `fn()` returns to reset `b.halfOpenIn`. However, if `fn()` panics, neither function is called. The mutex is unlocked properly, but `b.halfOpenIn` remains permanently incremented. Any future request in `HalfOpen` sees `b.halfOpenIn >= b.cfg.HalfOpenMaxCalls` and returns `ErrCircuitOpen`. The breaker is permanently stuck in `HalfOpen`.
Assessment: FAIL
Severity: HIGH
Notes: A `defer` should be used to catch panics and properly fail the request or clear the probe counter.

## Finding 3

Location: `internal/checkout/service.go:28`
Claimed Behavior: Scenario 1 (Without Breaker) accurately represents an unprotected call.
Observed Implementation: `CheckoutWithoutBreaker` returns a `Result` struct where the `State` field is zero-valued. When printed in `Result.String()`, state 0 resolves to `"CLOSED"`, making the demo output print `state=CLOSED` even though no circuit breaker is present.
Assessment: WARNING
Severity: LOW
Notes: Minor display issue that could slightly confuse readers.

## Finding 4

Location: `internal/payment/fake_server.go:34`
Claimed Behavior: Fake server simulates failures effectively.
Observed Implementation: `http.Error` writes the error string followed by a newline `\n`. This newline propagates into the error returned by the client, causing multi-line formatting in logs.
Assessment: WARNING
Severity: LOW
Notes: Causes documentation mismatch in README demo output.