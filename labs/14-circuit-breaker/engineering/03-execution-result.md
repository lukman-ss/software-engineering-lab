# Execution Result

## Build
Command:
```bash
go build ./cmd/demo
```
Result:
PASS (exited with 0)

## Tests
Command:
```bash
go test -count=1 -v ./...
```
Result:
```text
?   	circuitbreaker/cmd/demo	[no test files]
?   	circuitbreaker/internal/checkout	[no test files]
=== RUN   TestInitialStateIsClosed
--- PASS: TestInitialStateIsClosed (0.00s)
=== RUN   TestSuccessfulCallsStayClosed
--- PASS: TestSuccessfulCallsStayClosed (0.00s)
=== RUN   TestFailuresBelowThresholdStayClosed
--- PASS: TestFailuresBelowThresholdStayClosed (0.00s)
=== RUN   TestThresholdReachedOpens
--- PASS: TestThresholdReachedOpens (0.00s)
=== RUN   TestOpenFailsFast
--- PASS: TestOpenFailsFast (0.00s)
=== RUN   TestOpenDoesNotCallDownstream
--- PASS: TestOpenDoesNotCallDownstream (0.00s)
=== RUN   TestCooldownMovesToHalfOpenBehavior
--- PASS: TestCooldownMovesToHalfOpenBehavior (0.04s)
=== RUN   TestSuccessfulHalfOpenProbeCloses
--- PASS: TestSuccessfulHalfOpenProbeCloses (0.04s)
=== RUN   TestFailedHalfOpenProbeReopens
--- PASS: TestFailedHalfOpenProbeReopens (0.04s)
=== RUN   TestRecoveryAfterDependencyHealthy
--- PASS: TestRecoveryAfterDependencyHealthy (0.04s)
=== RUN   TestConcurrentAccess
--- PASS: TestConcurrentAccess (0.00s)
PASS
ok  	circuitbreaker/internal/circuitbreaker	0.258s
?   	circuitbreaker/internal/payment	[no test files]
=== RUN   TestCircuitBreakerIntegration
=== RUN   TestCircuitBreakerIntegration/downstream_fails,_CB_trips_open
=== RUN   TestCircuitBreakerIntegration/cooldown_and_recovery
--- PASS: TestCircuitBreakerIntegration (0.15s)
    --- PASS: TestCircuitBreakerIntegration/downstream_fails,_CB_trips_open (0.00s)
    --- PASS: TestCircuitBreakerIntegration/cooldown_and_recovery (0.15s)
PASS
ok  	circuitbreaker/tests	0.289s
```

## Race Detector
Command:
```bash
go test -race -count=1 -v ./...
```
Result:
```text
?   	circuitbreaker/cmd/demo	[no test files]
?   	circuitbreaker/internal/checkout	[no test files]
=== RUN   TestInitialStateIsClosed
--- PASS: TestInitialStateIsClosed (0.00s)
=== RUN   TestSuccessfulCallsStayClosed
--- PASS: TestSuccessfulCallsStayClosed (0.00s)
=== RUN   TestFailuresBelowThresholdStayClosed
--- PASS: TestFailuresBelowThresholdStayClosed (0.00s)
=== RUN   TestThresholdReachedOpens
--- PASS: TestThresholdReachedOpens (0.00s)
=== RUN   TestOpenFailsFast
--- PASS: TestOpenFailsFast (0.00s)
=== RUN   TestOpenDoesNotCallDownstream
--- PASS: TestOpenDoesNotCallDownstream (0.00s)
=== RUN   TestCooldownMovesToHalfOpenBehavior
--- PASS: TestCooldownMovesToHalfOpenBehavior (0.04s)
=== RUN   TestSuccessfulHalfOpenProbeCloses
--- PASS: TestSuccessfulHalfOpenProbeCloses (0.04s)
=== RUN   TestFailedHalfOpenProbeReopens
--- PASS: TestFailedHalfOpenProbeReopens (0.04s)
=== RUN   TestRecoveryAfterDependencyHealthy
--- PASS: TestRecoveryAfterDependencyHealthy (0.04s)
=== RUN   TestConcurrentAccess
--- PASS: TestConcurrentAccess (0.00s)
PASS
ok  	circuitbreaker/internal/circuitbreaker	1.259s
?   	circuitbreaker/internal/payment	[no test files]
=== RUN   TestCircuitBreakerIntegration
=== RUN   TestCircuitBreakerIntegration/downstream_fails,_CB_trips_open
=== RUN   TestCircuitBreakerIntegration/cooldown_and_recovery
--- PASS: TestCircuitBreakerIntegration (0.15s)
    --- PASS: TestCircuitBreakerIntegration/downstream_fails,_CB_trips_open (0.00s)
    --- PASS: TestCircuitBreakerIntegration/cooldown_and_recovery (0.15s)
PASS
ok  	circuitbreaker/tests	1.282s
```

## Demo
Command:
```bash
go run ./cmd/demo
```
Result:
```text
=== SCENARIO 1: WITHOUT CIRCUIT BREAKER (SLOW DEPENDENCY) ===
request=1 result=err=payment request error: Post "http://127.0.0.1:52472/pay": context deadline exceeded (Client.Timeout exceeded while awaiting headers) duration=101.99425ms state=CLOSED
request=2 result=err=payment request error: Post "http://127.0.0.1:52472/pay": context deadline exceeded (Client.Timeout exceeded while awaiting headers) duration=101.168541ms state=CLOSED
request=3 result=err=payment request error: Post "http://127.0.0.1:52472/pay": context deadline exceeded (Client.Timeout exceeded while awaiting headers) duration=100.652833ms state=CLOSED

downstream_calls=3

=== SCENARIO 2: WITH CIRCUIT BREAKER (FAIL-FAST ON DOWN DEPENDENCY) ===
request=1 result=err=payment failed: status 500 body internal payment server failure
 duration=993.958µs state=CLOSED
request=2 result=err=payment failed: status 500 body internal payment server failure
 duration=166.875µs state=CLOSED
request=3 result=err=payment failed: status 500 body internal payment server failure
 duration=159.167µs state=OPEN
request=4 result=err=circuit breaker is open duration=83ns state=OPEN
request=5 result=err=circuit breaker is open duration=83ns state=OPEN
request=6 result=err=circuit breaker is open duration=84ns state=OPEN

downstream_calls=3

=== SCENARIO 3: RECOVERY (HALF_OPEN -> CLOSED) ===
initial state=OPEN
waiting for cooldown (300ms)...
payment server recovered to HEALTHY. Current CB state=HALF_OPEN
sending probe request...
probe result: err=<nil>, state after probe=CLOSED
sending subsequent regular request...
subsequent result: err=<nil>, state=CLOSED

=== SCENARIO 4: FAILED RECOVERY (HALF_OPEN -> OPEN AGAIN) ===
circuit forced back to: OPEN
waiting for cooldown (300ms)...
dependency still DOWN. Current CB state=HALF_OPEN
sending probe request...
probe result: err=payment failed: status 500 body internal payment server failure
, state after failed probe=OPEN
sending next request while re-opened...
next request result: err=circuit breaker is open, state=OPEN
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
