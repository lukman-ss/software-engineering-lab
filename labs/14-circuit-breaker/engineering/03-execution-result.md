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
go test -v ./...
```
Result:
```text
?   	circuitbreaker/cmd/demo	[no test files]
?   	circuitbreaker/internal/checkout	[no test files]
=== RUN   TestCircuitBreaker
=== RUN   TestCircuitBreaker/1._initial_state_is_CLOSED
=== RUN   TestCircuitBreaker/2._successful_calls_stay_CLOSED
=== RUN   TestCircuitBreaker/3._failures_below_threshold_stay_CLOSED
=== RUN   TestCircuitBreaker/4._threshold_reached_changes_state_to_OPEN
=== RUN   TestCircuitBreaker/5._OPEN_calls_fail_fast
=== RUN   TestCircuitBreaker/6._OPEN_calls_do_not_execute_downstream_function
=== RUN   TestCircuitBreaker/7._cooldown_moves_breaker_toward_HALF_OPEN_behavior
=== RUN   TestCircuitBreaker/8._successful_HALF_OPEN_probe_closes_circuit
=== RUN   TestCircuitBreaker/9._failed_HALF_OPEN_probe_opens_circuit_again
=== RUN   TestCircuitBreaker/10._circuit_recovers_after_dependency_becomes_healthy
=== RUN   TestCircuitBreaker/11._concurrency_and_race_safety
=== RUN   TestCircuitBreaker/12._HalfOpenMaxCalls_>_1_limits_probes_and_requires_consecutive_successes
=== RUN   TestCircuitBreaker/13._panic_in_fn_during_HALF_OPEN_sets_state_back_to_OPEN
=== RUN   TestCircuitBreaker/14._default_config_values_applied
=== RUN   TestCircuitBreaker/15._panic_in_fn_during_CLOSED_records_failure
=== RUN   TestCircuitBreaker/16._success_in_CLOSED_resets_consecutive_failure_count
--- PASS: TestCircuitBreaker (0.00s)
    --- PASS: TestCircuitBreaker/1._initial_state_is_CLOSED (0.00s)
    --- PASS: TestCircuitBreaker/2._successful_calls_stay_CLOSED (0.00s)
    --- PASS: TestCircuitBreaker/3._failures_below_threshold_stay_CLOSED (0.00s)
    --- PASS: TestCircuitBreaker/4._threshold_reached_changes_state_to_OPEN (0.00s)
    --- PASS: TestCircuitBreaker/5._OPEN_calls_fail_fast (0.00s)
    --- PASS: TestCircuitBreaker/6._OPEN_calls_do_not_execute_downstream_function (0.00s)
    --- PASS: TestCircuitBreaker/7._cooldown_moves_breaker_toward_HALF_OPEN_behavior (0.00s)
    --- PASS: TestCircuitBreaker/8._successful_HALF_OPEN_probe_closes_circuit (0.00s)
    --- PASS: TestCircuitBreaker/9._failed_HALF_OPEN_probe_opens_circuit_again (0.00s)
    --- PASS: TestCircuitBreaker/10._circuit_recovers_after_dependency_becomes_healthy (0.00s)
    --- PASS: TestCircuitBreaker/11._concurrency_and_race_safety (0.00s)
    --- PASS: TestCircuitBreaker/12._HalfOpenMaxCalls_>_1_limits_probes_and_requires_consecutive_successes (0.00s)
    --- PASS: TestCircuitBreaker/13._panic_in_fn_during_HALF_OPEN_sets_state_back_to_OPEN (0.00s)
    --- PASS: TestCircuitBreaker/14._default_config_values_applied (0.00s)
    --- PASS: TestCircuitBreaker/15._panic_in_fn_during_CLOSED_records_failure (0.00s)
    --- PASS: TestCircuitBreaker/16._success_in_CLOSED_resets_consecutive_failure_count (0.00s)
PASS
ok  	circuitbreaker/internal/circuitbreaker	0.180s
?   	circuitbreaker/internal/payment	[no test files]
=== RUN   TestCircuitBreakerIntegration
=== RUN   TestCircuitBreakerIntegration/downstream_fails,_CB_trips_open
=== RUN   TestCircuitBreakerIntegration/cooldown_and_recovery
--- PASS: TestCircuitBreakerIntegration (0.15s)
    --- PASS: TestCircuitBreakerIntegration/downstream_fails,_CB_trips_open (0.00s)
    --- PASS: TestCircuitBreakerIntegration/cooldown_and_recovery (0.15s)
PASS
ok  	circuitbreaker/tests	0.156s
```

## Race Detector
Command:
```bash
go test -race -v ./...
```
Result:
```text
?   	circuitbreaker/cmd/demo	[no test files]
?   	circuitbreaker/internal/checkout	[no test files]
=== RUN   TestCircuitBreaker
=== RUN   TestCircuitBreaker/1._initial_state_is_CLOSED
=== RUN   TestCircuitBreaker/2._successful_calls_stay_CLOSED
=== RUN   TestCircuitBreaker/3._failures_below_threshold_stay_CLOSED
=== RUN   TestCircuitBreaker/4._threshold_reached_changes_state_to_OPEN
=== RUN   TestCircuitBreaker/5._OPEN_calls_fail_fast
=== RUN   TestCircuitBreaker/6._OPEN_calls_do_not_execute_downstream_function
=== RUN   TestCircuitBreaker/7._cooldown_moves_breaker_toward_HALF_OPEN_behavior
=== RUN   TestCircuitBreaker/8._successful_HALF_OPEN_probe_closes_circuit
=== RUN   TestCircuitBreaker/9._failed_HALF_OPEN_probe_opens_circuit_again
=== RUN   TestCircuitBreaker/10._circuit_recovers_after_dependency_becomes_healthy
=== RUN   TestCircuitBreaker/11._concurrency_and_race_safety
=== RUN   TestCircuitBreaker/12._HalfOpenMaxCalls_>_1_limits_probes_and_requires_consecutive_successes
=== RUN   TestCircuitBreaker/13._panic_in_fn_during_HALF_OPEN_sets_state_back_to_OPEN
=== RUN   TestCircuitBreaker/14._default_config_values_applied
=== RUN   TestCircuitBreaker/15._panic_in_fn_during_CLOSED_records_failure
=== RUN   TestCircuitBreaker/16._success_in_CLOSED_resets_consecutive_failure_count
--- PASS: TestCircuitBreaker (0.00s)
PASS
ok  	circuitbreaker/internal/circuitbreaker	1.161s
?   	circuitbreaker/internal/payment	[no test files]
=== RUN   TestCircuitBreakerIntegration
=== RUN   TestCircuitBreakerIntegration/downstream_fails,_CB_trips_open
=== RUN   TestCircuitBreakerIntegration/cooldown_and_recovery
--- PASS: TestCircuitBreakerIntegration (0.15s)
PASS
ok  	circuitbreaker/tests	1.277s
```

## Demo
Command:
```bash
go run ./cmd/demo
```
Result:
```text
==================================================
  LAB 14: CIRCUIT BREAKER PATTERN DEMONSTRATION   
==================================================

=== SCENARIO 1: WITHOUT CIRCUIT BREAKER (SLOW DEPENDENCY) ===
request=1 result=timeout/error duration=102ms
request=2 result=timeout/error duration=102ms
request=3 result=timeout/error duration=102ms
downstream_calls=3 (all requests blocked and hit downstream)

=== SCENARIO 2: WITH CIRCUIT BREAKER (FAIL-FAST ON DOWN DEPENDENCY) ===
request=1 result=payment_error              state=CLOSED    duration=553.875µs
request=2 result=payment_error              state=CLOSED    duration=168.625µs
request=3 result=payment_error              state=OPEN      duration=166.625µs
request=4 result=circuit_open (fail-fast)   state=OPEN      duration=708ns
request=5 result=circuit_open (fail-fast)   state=OPEN      duration=500ns
request=6 result=circuit_open (fail-fast)   state=OPEN      duration=375ns
downstream_calls=3 (downstream calls stopped once OPEN)

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
probe result: err=true, state after failed probe=OPEN
sending next request while re-opened...
next request result: err=checkout payment failed (with CB): circuit breaker is open, state=OPEN

==================================================
  DEMO COMPLETED SUCCESSFULLY                     
==================================================
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
