# Code Audit

File: internal/circuitbreaker/circuit_breaker.go

## Finding 1
Location: Execute + onSuccessLocked/onFailureLocked generation check (circuit_breaker.go:91-157)
Claimed Behavior: In-flight requests from a prior generation do not corrupt newer state transitions.
Observed Implementation: Execute captures `gen := b.generation` after acquiring lock and re-checking state; onSuccessLocked/onFailureLocked skip mutations when `b.generation != gen`. advanceLocked increments generation when transitioning Open->HalfOpen or when HalfOpen/Open trip.
Assessment: PASS
Severity: LOW
Notes: Correctly prevents stale trailing requests from flipping state. Covered by TestTrailingInFlightRequestDoesNotCorruptNewState.

## Finding 2
Location: Execute Open/HalfOpen path (circuit_breaker.go:95-106)
Claimed Behavior: OPEN fails immediately with ErrCircuitOpen without invoking fn; HALF_OPEN allows up to HalfOpenMaxCalls, rest fail fast.
Observed Implementation: In OPEN, returns ErrCircuitOpen before unlocking, fn never called. In HALF_OPEN, increments halfOpenIn when below limit; returns ErrCircuitOpen when at limit.
Assessment: PASS
Severity: LOW
Notes: Matches README "HALF_OPEN ... allows limited probe calls". The halfOpenIn counter is reset on transition to HalfOpen via advanceLocked and on generation bump.

## Finding 3
Location: onFailureLocked transition (circuit_breaker.go:141-157)
Claimed Behavior: When failures reach FailureThreshold, breaker trips to OPEN and records openedAt.
Observed Implementation: failures++ then `if b.failures >= b.cfg.FailureThreshold { state = Open; openedAt = now; generation++ }`.
Assessment: PASS
Severity: LOW
Notes: Threshold semantics `>=` consistent with README "When failures >= FailureThreshold". DefaultConfig FailureThreshold=3 matches README scenario 2 (3 failures then OPEN on 3rd).

## Finding 4
Location: advanceLocked (circuit_breaker.go:83-89)
Claimed Behavior: OPEN -> HALF_OPEN after OpenTimeout elapses.
Observed Implementation: Checks `b.state == Open && now.Sub(b.openedAt) >= b.cfg.OpenTimeout`; sets HalfOpen, resets halfOpenIn, bumps generation.
Assessment: PASS
Severity: LOW
Notes: `>=` semantics allow the first call after cooldown to become the probe.

## Finding 5
Location: onSuccessLocked HalfOpen branch (circuit_breaker.go:127-139)
Claimed Behavior: Successful probe in HALF_OPEN -> CLOSED, clears failures and generation bump to invalidate stale ops.
Observed Implementation: Sets state=Closed, failures=0, halfOpenIn=0, generation++.
Assessment: PASS
Severity: LOW
Notes: Correct recovery path.

## Finding 6
Location: onFailureLocked HalfOpen branch (circuit_breaker.go:145-150)
Claimed Behavior: Failed probe in HALF_OPEN -> OPEN with new openedAt (restart cooldown).
Observed Implementation: Sets state=Open, openedAt=now, halfOpenIn=0, generation++.
Assessment: PASS
Severity: LOW
Notes: Cooldown restarts on failed probe (README scenario 4).

## Finding 7
Location: State() (circuit_breaker.go:76-81)
Claimed Behavior: State() reflects elapsed-time transitions (OPEN->HALF_OPEN) without a successful call.
Observed Implementation: Calls advanceLocked(time.Now()) under lock before returning.
Assessment: PASS
Severity: LOW
Notes: Necessary so Read-only observers see time-based transition; State is used by tests and demo to assert/recover states.

## Finding 8
Location: New config guard (circuit_breaker.go:63-74)
Claimed Behavior: Non-positive config values fall back to safe defaults.
Observed Implementation: FailureThreshold<=0 -> 3; OpenTimeout<=0 -> 300ms; HalfOpenMaxCalls<=0 -> 1.
Assessment: PASS
Severity: LOW
Notes: Defensive defaults match DefaultConfig.

## Finding 9
Location: Panic handling (circuit_breaker.go:109-124)
Claimed Behavior: Panic in fn treats as failure and releases lock safely.
Observed Implementation: `panicked:=true`; deferred re-lock; on panic path calls onFailureLocked; panic propagates after unlock (defer recovers implicitly via lock release). The `return err` path is skipped on panic since panic unwinds; deferred fn runs, locks, records failure, unlocks; panic continues propagating outside.
Assessment: PASS
Severity: LOW
Notes: Lock hygiene correct (unlock -> fn() -> lock). Covered by TestPanicInHalfOpenCleansUpState. Panic in HalfOpen correctly re-trips to OPEN.

## Finding 10
Location: checkout/service.go CheckoutWithBreaker (service.go:34-42)
Claimed Behavior: Wraps payment call in breaker; reports downstream OK / failure / state.
Observed Implementation: Calls breaker.Execute(ProcessPayment); on err returns Result{Err,Duration,State,HasBreaker:true}; on success returns Result{Duration,State,DownstreamOK:true,HasBreaker:true}.
Assessment: PASS
Severity: LOW
Notes: State snapshot taken after Execute settles state, so reported State matches post-call state (matches README scenario 2 where request=3 shows state=OPEN after trip).

## Finding 11
Location: checkout/service.go CheckoutWithoutBreaker (service.go:44-49)
Claimed Behavior: Bypasses breaker; HasBreaker=false so String() omits state field.
Observed Implementation: HasBreaker: false; no breaker.Execute. String() emits `stateStr=""` when !HasBreaker.
Assessment: PASS
Severity: LOW
Notes: Matches README scenario 1 (no state field present).

## Finding 12
Location: internal/payment/fake_server.go requestCount
Claimed Behavior: Each /pay hit increments atomic counter, enabling downstream_call assertions.
Observed Implementation: fs.requestCount.Add(1) at handler entry.
Assessment: PASS
Severity: LOW
Notes: Used by demo and integration tests to prove fail-fast means zero extra downstream calls.

## Finding 13
Location: internal/circuitbreaker State enum / String() (circuit_breaker.go:11-37)
Claimed Behavior: States Closed/Closed/HalfOpen map to "CLOSED"/"OPEN"/"HALF_OPEN"; exported StateClosed/StateOpen/StateHalfOpen aliases for test compatibility.
Observed Implementation: String() returns exact README casing; aliases defined.
Assessment: PASS
Severity: LOW
Notes: README scenario 4 uses `state=OPEN`/`state=HALF_OPEN`/`state=CLOSED` which match.

## Finding 14
Location: circuit_breaker.go exported constants naming
Claimed Behavior: Integration tests use StateClosed/StateOpen/StateHalfOpen.
Observed Implementation: const block exports these aliases; integration_test.go uses circuitbreaker.StateOpen etc.
Assessment: PASS
Severity: LOW
Notes: Minor API duplication but harmless and documented as "for integration test compatibility".

## Finding 15
Location: README vs code semantics (README:17 "Successes reset failure count")
Claimed Behavior: A success in CLOSED resets failures to 0.
Observed Implementation: onSuccessLocked in Closed branch sets b.failures=0.
Assessment: PASS
Severity: LOW
Notes: Covered by TestSuccessInClosedResetsFailures.
