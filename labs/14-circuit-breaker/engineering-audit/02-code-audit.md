# Engineering Code Audit: Circuit Breaker

## Finding 1

Location: `internal/circuitbreaker/circuit_breaker.go:141-158`
Claimed Behavior: Circuit transitions from CLOSED to OPEN when consecutive failures reach `FailureThreshold`.
Observed Implementation: `onFailureLocked` checks generation match, increments `b.failures`, and when `b.failures >= b.cfg.FailureThreshold`, sets `b.state = Open`, updates `b.openedAt = now`, and increments `b.generation++`.
Assessment: PASS
Severity: LOW
Notes: Correctly tracks failures and resets or increments generation on transition.

## Finding 2

Location: `internal/circuitbreaker/circuit_breaker.go:95-98`
Claimed Behavior: When OPEN, calls fail fast immediately returning `ErrCircuitOpen` with zero downstream calls.
Observed Implementation: Mutex acquired, `advanceLocked` checks if cooldown elapsed. If still OPEN, returns `ErrCircuitOpen` immediately after mutex unlock without invoking `fn()`.
Assessment: PASS
Severity: LOW
Notes: Sub-microsecond fail-fast observed in execution and tests. Downstream function is not called.

## Finding 3

Location: `internal/circuitbreaker/circuit_breaker.go:83-89`
Claimed Behavior: After cooldown duration (`OpenTimeout`), circuit transitions to HALF-OPEN.
Observed Implementation: State advancement is evaluated lazily in `advanceLocked` invoked by both `State()` and `Execute()`. If `now.Sub(b.openedAt) >= b.cfg.OpenTimeout`, transitions to `HalfOpen`, resets `halfOpenIn = 0`, and bumps `generation++`.
Assessment: PASS
Severity: LOW
Notes: Lazy evaluation avoids unnecessary timer goroutines. Thread-safe under `b.mu`.

## Finding 4

Location: `internal/circuitbreaker/circuit_breaker.go:99-106`
Claimed Behavior: In HALF-OPEN state, limit concurrent probes up to `HalfOpenMaxCalls`; excess calls fail fast.
Observed Implementation: While in `HalfOpen`, if `halfOpenIn >= cfg.HalfOpenMaxCalls`, rejects additional concurrent requests with `ErrCircuitOpen`. Otherwise increments `halfOpenIn` and executes probe.
Assessment: PASS
Severity: LOW
Notes: Successfully prevents probe storms against recovering downstream services.

## Finding 5

Location: `internal/circuitbreaker/circuit_breaker.go:94,127-130,141-144`
Claimed Behavior: In-flight calls from earlier states/generations must not corrupt newly transitioned states.
Observed Implementation: `gen := b.generation` captured under lock before execution. In deferred completion handlers (`onSuccessLocked` / `onFailureLocked`), if `b.generation != gen`, the result is discarded.
Assessment: PASS
Severity: LOW
Notes: Effectively neutralizes race condition where slow requests from CLOSED or OPEN state resolve after breaker has transitioned to HALF-OPEN or CLOSED.

## Finding 6

Location: `internal/circuitbreaker/circuit_breaker.go:109-125`
Claimed Behavior: Panics during user execution must not leave breaker corrupted or lock held.
Observed Implementation: Mutex is unlocked during `fn()` execution. Deferred function runs with boolean flag `panicked := true` flipped to `false` only on clean return. If `fn()` panics, defer re-locks mutex and treats panic as failure (`onFailureLocked`), allowing panic to re-propagate cleanly without stranding state counters.
Assessment: PASS
Severity: LOW
Notes: Verified by `TestPanicInHalfOpenCleansUpState`.

## Finding 7

Location: `internal/circuitbreaker/circuit_breaker.go:63-74`
Claimed Behavior: Invalid or zero configuration values must fallback to safe defaults.
Observed Implementation: `New()` validates `FailureThreshold <= 0`, `OpenTimeout <= 0`, `HalfOpenMaxCalls <= 0` and sets defaults (3, 300ms, 1 respectively).
Assessment: PASS
Severity: LOW
Notes: Robust guard against uninitialized struct config.

## Finding 8

Location: `internal/payment/fake_server.go:34-45` & `internal/payment/client.go:27-37`
Claimed Behavior: Simulated HTTP client and server handle timeout, 500 error, and healthy 200 response states.
Observed Implementation: `FakeServer` uses `httptest.Server` with atomic mode flag (`ModeHealthy`, `ModeSlow`, `ModeDown`). `Client` uses standard `http.Client` with explicit timeouts.
Assessment: PASS
Severity: LOW
Notes: Clean test doubles; no external dependencies or mock frameworks required.
