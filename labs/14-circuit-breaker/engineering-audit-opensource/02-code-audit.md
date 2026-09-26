## Finding 1

Location: `internal/circuitbreaker/circuit_breaker.go:75-80` (State method)
Claimed Behavior: `State()` returns current state without side effects.
Observed Implementation: `State()` locks, calls `advanceLocked(time.Now())`, which may mutate state (OPEN → HALF_OPEN) based on elapsed time, then unlocks and returns state.
Assessment: WARNING (behavioral side effect in getter)
Severity: MEDIUM
Notes: A getter that mutates internal state violates principle of least surprise. Although protected by mutex and conceptually correct (state should advance on query), callers may expect `State()` to be a pure read. The pattern is acceptable if documented, but it is not. In typical circuit breaker implementations, state transition on query is standard (see Akka, Hystrix), so this is idiomatic. Still flagged as MEDIUM because it violates pure-getter expectation; however, it is correct for the domain.

## Finding 2

Location: `internal/circuitbreaker/circuit_breaker.go:82-86` (advanceLocked)
Claimed Behavior: After `OpenTimeout` elapses, state advances from OPEN to HALF_OPEN exactly once per cooldown period.
Observed Implementation: Condition checks `b.state == Open && now.Sub(b.openedAt) >= b.cfg.OpenTimeout`, then sets `b.state = HalfOpen` and resets `b.halfOpenIn = 0`. No guard against repeated calls advancing state multiple times if `State()` called rapidly while still in cooldown window.
Assessment: PASS
Severity: LOW
Notes: `b.state == Open` guard ensures transition only happens once per cooldown because once state becomes HalfOpen, guard fails. Subsequent calls while still HalfOpen do nothing. Correct.

## Finding 3

Location: `internal/circuitbreaker/circuit_breaker.go:89-116` (Execute method)
Claimed Behavior: `Execute` atomically decides whether to allow a call, increments halfOpenIn if needed, calls `fn`, then updates state based on result.
Observed Implementation:
- Lock, advanceLocked, switch on state:
  - OPEN: unlock, return ErrCircuitOpen.
  - HALF_OPEN: if `b.halfOpenIn >= b.cfg.HalfOpenMaxCalls`, unlock, return ErrCircuitOpen; else increment.
  - CLOSED: fallthrough.
- Unlock before calling `fn()` (critical: allows concurrency).
- Lock again on return, defer unlock.
- On success: call `onSuccessLocked`.
- On failure: call `onFailureLocked(now)`.
Observed call to `fn()` occurs WITHOUT lock held. This is correct to avoid blocking downstream on breaker lock.
Assessment: PASS
Severity: LOW
Notes: Textbook "check-then-act" pattern with lock released for the duration of the call. Potential race: between unlocking before `fn` and re-locking after, state could change (e.g., another goroutine trips breaker OPEN). This is by design; circuit breaker permits racing calls to see latest state on re-entry. Implementation correctly re-locks and re-evaluates via `advanceLocked` at top of second critical section.

## Finding 4

Location: `internal/circuitbreaker/circuit_breaker.go:118-126` (onSuccessLocked)
Claimed Behavior: On success in HALF_OPEN, reset to CLOSED, zero failures and halfOpenIn.
Observed Implementation: If `b.state == HalfOpen`, set state to Closed, zero failures and halfOpenIn, return. Else (CLOSED/OPEN) just zero failures.
Assessment: PASS
Severity: LOW
Notes: Correct and matches spec. Note: success in OPEN state is impossible because Execute short-circuits with ErrCircuitOpen before calling `fn`. So onSuccessLocked only sees CLOSED or HALF_OPEN.

## Finding 5

Location: `internal/circuitbreaker/circuit_breaker.go:128-139` (onFailureLocked)
Claimed Behavior: On failure in HALF_OPEN, re-open immediately; in CLOSED, increment failures and open if threshold reached.
Observed Implementation: If HalfOpen: set state Open, record openedAt, zero halfOpenIn, return. Else (CLOSED/OPEN): increment failures; if >= FailureThreshold, set Open and record openedAt.
Assessment: PASS
Severity: LOW
Notes: Correct. Note: failures in OPEN state are ignored (breaker already Open) – acceptable because OPEN state already failing fast.

## Finding 6

Location: `internal/circuitbreaker/circuit_breaker.go:53-60` (Breaker struct)
Claimed Behavior: Fields protected by mutex; failures, state, openedAt, halfOpenIn correctly typed.
Observed Implementation:
- `mu sync.Mutex` guards all mutable state.
- `cfg Config` (immutable after construction).
- `state State`.
- `failures int`.
- `openedAt time.Time`.
- `halfOpenIn int` (calls allowed in HALF_OPEN).
Assessment: PASS
Severity: LOW
Notes: All fields that change on transitions are under `mu`. Config is set at construction and never mutated (New copies and applies defaults). No data races.

## Finding 7

Location: `internal/circuitbreaker/circuit_breaker.go:62-73` (New)
Claimed Behavior: Constructor applies sensible defaults for zero values.
Observed Implementation: If `FailureThreshold <= 0` → 3; if `OpenTimeout <= 0` → 300ms; if `HalfOpenMaxCalls <= 0` → 1.
Assessment: PASS
Severity: LOW
Notes: Defensive copying prevents panics due to misconfiguration. Good.

## Finding 8

Location: `internal/checkout/service.go:32-41` (CheckoutWithBreaker)
Claimed Behavior: Service method wraps payment call in breaker.Execute, returns Result with Err, Duration, DownstreamOK, State.
Observed Implementation:
- Records start time.
- Calls `s.breaker.Execute(func() error { return s.payment.ProcessPayment(context.Background()) })`.
- On error: returns Result{Err: err, Duration: time.Since(start), State: s.breaker.State()}.
- On success: returns Result{Duration: ..., State: ..., DownstreamOK: true}.
Assessment: PASS
Severity: LOW
Notes: Correct. Note: DownstreamOK only set on success (consistent with demo). State() called after lock release in Execute – may race with state transitions, but that's fine; Result.State is best-effort.

## Finding 9

Location: `internal/payment/client.go:28-45` (ProcessPayment)
Claimed Behavior: Simple HTTP POST with context and timeout; maps non-200 to error.
Observed Implementation: Creates POST request, does client.Do, reads body on non-200, wraps ErrPaymentFailed with status and body.
Assessment: PASS
Severity: LOW
Notes: Proper use of context, timeout, and error wrapping. No issues.

## Finding 10

Location: `internal/payment/fake_server.go:18-71` (FakeServer)
Claimed Behavior: httptest server with HEALTHY/SLOW/DOWN modes, request counting, SetMode atomic.
Observed Implementation:
- `mode atomic.Value` for safe concurrent SetMode.
- `requestCount atomic.Int64` for lock-free increments.
- Handler reads mode via `fs.mode.Load().(ServerMode)`.
- Slow mode `time.Sleep(fs.slowDelay)`.
Assessment: PASS
Severity: LOW
Notes: Clean, race-free test double. Good.

## Finding 11

Location: `tests/integration_test.go:14-69` (TestCircuitBreakerIntegration)
Claimed Behavior: Integration test verifies breaker trips open after FailureThreshold downstream failures, then cooldown to HALF_OPEN, then probe success leads to CLOSED.
Observed Implementation:
- FakeServer with 250ms slow delay (not used in test, ModeDown used).
- Client timeout 50ms to exaggerate slowness (though ModeDown immediate).
- Breaker Config{FailureThreshold: 2, OpenTimeout: 100ms, HalfOpenMaxCalls: 1}.
- Subtest "downstream fails, CB trips open": two failing calls, verify OPEN, third call fails fast, request count == 2.
- Subtest "cooldown and recovery": sleep 150ms (> OpenTimeout), set server ModeHealthy, verify HALF_OPEN, call succeeds, verify CLOSED.
Assessment: PASS
Severity: LOW
Notes: Times chosen provide margin; test is deterministic and correct.

## Finding 12

Location: `internal/circuitbreaker/circuit_breaker_test.go:163-183` (TestConcurrentAccess)
Claimed Behavior: Concurrent Execute and State calls under load do not corrupt state or panic.
Observed Implementation: 50 goroutines, each calls Execute (with 1/3 failure probability) and State, waits for WG, then asserts final state is one of {Open, Closed, HalfOpen}.
Assessment: PASS
Severity: LOW
Notes: Basic liveness safety check; does not verify exact counts or fairness, but sufficient to catch gross races or panics. Passes with race detector.