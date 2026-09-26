# Code Audit

## Finding 1

Location: internal/circuitbreaker/circuit_breaker.go:73-155 (Execute)
Claimed Behavior: CLOSED counts failures, successes reset count, trips to OPEN at `failures >= FailureThreshold`. OPEN fail-fasts with ErrCircuitOpen. HALF_OPEN limits probes to HalfOpenMaxCalls; success -> CLOSED; failure -> OPEN.
Observed Implementation: `Execute` locks mutex, runs `checkStateTransitionLocked` (OPEN -> HALF_OPEN after OpenTimeout), then switches on state. CLOSED branch increments `failureCount` on error and trips to OPEN when `>= FailureThreshold` (line 145-148); resets `failureCount = 0` on success (line 152). HALF_OPEN branch increments `halfOpenCalls`, allows `fn()` only when `halfOpenCalls < HalfOpenMaxCalls` else returns ErrCircuitOpen (lines 83-86); on success increments `consecutiveSuccesses` and moves to CLOSED when `>= HalfOpenMaxCalls` (lines 114-120); on error/probe-failure moves to OPEN (lines 106-111). OPEN returns ErrCircuitOpen immediately (lines 78-80).
Assessment: PASS
Severity: LOW
Notes: Matches research state machine exactly. Mutex held only outside downstream `fn()` execution (correct — avoids blocking on network).

## Finding 2

Location: internal/circuitbreaker/circuit_breaker.go:64-71 (checkStateTransitionLocked)
Claimed Behavior: OPEN transitions to HALF_OPEN after OpenTimeout cooldown elapses (time-triggered, not probe-triggered).
Observed Implementation: Checks `cb.state == StateOpen && cb.now().Sub(cb.lastStateChange) >= cb.config.OpenTimeout`, sets HALF_OPEN, resets probe counters. Called under lock from both `State()` and `Execute()`.
Assessment: PASS
Severity: LOW
Notes: Correct time-triggered transition. `now` injectable for deterministic tests.

## Finding 3

Location: internal/circuitbreaker/circuit_breaker.go:90-137 (panic handling)
Claimed Behavior: Panics propagate to caller; breaker state restored safely (HALF_OPEN panic -> OPEN; CLOSED panic -> increment failure, trip to OPEN if threshold reached).
Observed Implementation: HALF_OPEN branch has defer-recover that sets state=OPEN, resets counters, then re-panics (lines 90-100). CLOSED branch defer-recover increments failureCount and trips to OPEN if threshold met, then re-panics (lines 126-137). No double-unlock: recover acquires its own lock; original unlock already done.
Assessment: PASS
Severity: LOW
Notes: Panic propagation preserved; state machine consistency maintained under panic.

## Finding 4

Location: internal/circuitbreaker/circuit_breaker.go:38-55 (New)
Claimed Behavior: Sensible defaults when config values are zero/negative.
Observed Implementation: FailureThreshold<=0 -> 3; OpenTimeout<=0 -> 5s; HalfOpenMaxCalls<=0 -> 1.
Assessment: PASS
Severity: LOW
Notes: Defensive defaults. Covered by test 14.

## Finding 5

Location: internal/circuitbreaker/circuit_breaker.go:11 (State type) / var ErrCircuitOpen
Claimed Behavior: Distinct error ErrCircuitOpen returned for fail-fast OPEN state.
Observed Implementation: `ErrCircuitOpen = errors.New("circuit breaker is open")`; returned from StateOpen and probe-limit-exceeded HALF_OPEN paths.
Assessment: PASS
Severity: LOW
Notes: Errors.Is-compatible (value sentinel); tested in tests 5, 6, 12.

## Finding 6

Location: internal/payment/fake_server.go / client.go / checkout/service.go
Claimed Behavior: Checkout service guards payment via circuit breaker; nil CB bypasses breaker (no-CB path).
Observed Implementation: `Service.Checkout` wraps `paymentClient.ProcessPayment(ctx)` in `cb.Execute` when cb != nil (lines 24-31); direct call otherwise (lines 34-37). Payment client uses http.Client.Timeout; fake server supports HEALTHY/SLOW/DOWN modes and tracks request count.
Assessment: PASS
Severity: LOW
Notes: Clean composition. FakeServer.RequestCount enables verifying "zero downstream calls when OPEN".

## Finding 7

Location: cmd/demo/main.go:64-77, 118 (output formatting)
Claimed Behavior: Demo prints state and duration per request and downstream call count, demonstrating fail-fast and recovery.
Observed Implementation: Iterates requests, prints result/state/duration; prints `downstream_calls` from fakeServer.RequestCount. Scenario 4 prints probe error truthiness via `%v` with `err != nil` boolean (line 118) — matches README "err=true".
Assessment: PASS
Severity: LOW
Notes: Demo mirrors README structure exactly. Timing values are illustrative (see Gaps).

## Finding 8

Location: internal/circuitbreaker/circuit_breaker.go — unused fields
Claimed Behavior: (none claimed)
Observed Implementation: `consecutiveSuccesses` and `halfOpenCalls` are tracked; `consecutiveSuccesses` used only in HALF_OPEN success path; no dead fields affecting correctness.
Assessment: PASS
Severity: LOW
Notes: Minor: `consecutiveSuccesses` is redundant with `halfOpenCalls >= HalfOpenMaxCalls` logic but harmless and clarifies intent.
