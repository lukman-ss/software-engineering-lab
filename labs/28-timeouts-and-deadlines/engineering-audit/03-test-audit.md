# Test Audit

Target Lab: labs/28-timeouts-and-deadlines

## Test Coverage Summary

1. `internal/deadline/deadline_test.go`:
   - `TestExecuteWithBudget_Success`: Verifies quick worker finishes with nil error. (PASS)
   - `TestExecuteWithBudget_Timeout`: Verifies long worker triggers context deadline exceeded. (PASS)
   - `TestExecuteWithBudget_ParentTimeoutInherited`: Verifies shorter parent timeout takes precedence over budget. (PASS)

2. `internal/retry/retry_test.go`:
   - `TestRetrier_SuccessOnFirstTry`: Verifies no retry when first attempt succeeds. (PASS)
   - `TestRetrier_RetryUntilSuccess`: Verifies multiple attempts until transient failure resolves. (PASS)
   - `TestRetrier_ExceedMaxAttempts`: Verifies terminal failure when max attempts reached. (PASS)
   - `TestRetrier_ContextCanceled`: Verifies retry loop halts immediately on context cancellation. (PASS)
   - `TestRetrier_JitterBoundsAndZeroConfig`: Verifies backoff duration remains within `[0, temp]` bounds and zero-value config defaults. (PASS)

3. `internal/circuit/circuit_test.go`:
   - `TestCircuitBreaker_StateTransitions`: Verifies CLOSED -> OPEN -> HALF_OPEN -> CLOSED cycle. (PASS)
   - `TestCircuitBreaker_HalfOpenFailureTripsOpen`: Verifies single failure in HALF_OPEN trips immediately back to OPEN. (PASS)
   - `TestCircuitBreaker_DefaultZeroConfig`: Verifies default threshold/cooldown fallbacks. (PASS)

4. `internal/idempotency/idempotency_test.go`:
   - `TestStore_GetSet`: Verifies record write and read. (PASS)
   - `TestStore_ConcurrentAccess`: Verifies concurrent read/write safety under race detector. (PASS)
   - `TestStore_LazyEvictionOnGet`: Verifies expired record deletion upon Get. (PASS)

5. `tests/integration_test.go`:
   - `TestIntegration_RetryWithCircuitBreaker`: Verifies retrier encountering failing service trips the circuit breaker to OPEN. (PASS)
   - `TestIntegration_IdempotentRetry`: Verifies retried duplicate executions return cached response without duplicate side-effects. (PASS)

## Execution Output

### `go test -v ./...`
```text
?   	timeouts-and-deadlines/cmd/demo	[no test files]
=== RUN   TestCircuitBreaker_StateTransitions
--- PASS: TestCircuitBreaker_StateTransitions (0.06s)
=== RUN   TestCircuitBreaker_HalfOpenFailureTripsOpen
--- PASS: TestCircuitBreaker_HalfOpenFailureTripsOpen (0.04s)
=== RUN   TestCircuitBreaker_DefaultZeroConfig
--- PASS: TestCircuitBreaker_DefaultZeroConfig (0.00s)
PASS
ok  	timeouts-and-deadlines/internal/circuit	0.185s
=== RUN   TestExecuteWithBudget_Success
--- PASS: TestExecuteWithBudget_Success (0.00s)
=== RUN   TestExecuteWithBudget_Timeout
--- PASS: TestExecuteWithBudget_Timeout (0.02s)
=== RUN   TestExecuteWithBudget_ParentTimeoutInherited
--- PASS: TestExecuteWithBudget_ParentTimeoutInherited (0.02s)
PASS
ok  	timeouts-and-deadlines/internal/deadline	0.129s
=== RUN   TestStore_GetSet
--- PASS: TestStore_GetSet (0.12s)
=== RUN   TestStore_ConcurrentAccess
--- PASS: TestStore_ConcurrentAccess (0.00s)
=== RUN   TestStore_LazyEvictionOnGet
--- PASS: TestStore_LazyEvictionOnGet (0.03s)
PASS
ok  	timeouts-and-deadlines/internal/idempotency	0.237s
=== RUN   TestRetrier_SuccessOnFirstTry
--- PASS: TestRetrier_SuccessOnFirstTry (0.00s)
=== RUN   TestRetrier_RetryUntilSuccess
--- PASS: TestRetrier_RetryUntilSuccess (0.00s)
=== RUN   TestRetrier_ExceedMaxAttempts
--- PASS: TestRetrier_ExceedMaxAttempts (0.00s)
=== RUN   TestRetrier_ContextCanceled
--- PASS: TestRetrier_ContextCanceled (0.02s)
=== RUN   TestRetrier_JitterBoundsAndZeroConfig
--- PASS: TestRetrier_JitterBoundsAndZeroConfig (0.00s)
PASS
ok  	timeouts-and-deadlines/internal/retry	0.109s
=== RUN   TestIntegration_RetryWithCircuitBreaker
--- PASS: TestIntegration_RetryWithCircuitBreaker (0.02s)
=== RUN   TestIntegration_IdempotentRetry
--- PASS: TestIntegration_IdempotentRetry (0.00s)
PASS
ok  	timeouts-and-deadlines/tests	0.101s
```

### `go test -race ./...`
```text
ok  	timeouts-and-deadlines/internal/circuit	0.181s
ok  	timeouts-and-deadlines/internal/deadline	0.129s
ok  	timeouts-and-deadlines/internal/idempotency	0.230s
ok  	timeouts-and-deadlines/internal/retry	0.106s
ok  	timeouts-and-deadlines/tests	0.101s
```

### `go run ./cmd/demo`
```text
=== Timeouts & Deadlines Laboratory Demo ===

--- Demo 1: Context Deadline & Budget Propagation ---
Deadline propagation result: context deadline exceeded

--- Demo 2: Exponential Backoff with Full Jitter ---
  Attempt #1 executed
  Attempt #2 executed
  Attempt #3 executed
Retry execution result: <nil> (total attempts: 3)

--- Demo 3: Circuit Breaker State Transitions ---
Initial state: CLOSED
State after 2 failures: OPEN
Execution attempt while OPEN: circuit breaker is open
State after cooldown: HALF_OPEN
Execution attempt in HALF_OPEN (success): <nil> -> New State: CLOSED

--- Demo 4: Idempotent Request Retry Protection ---
First request execution: Charged $100 successfully
Retried request execution: Charged $100 successfully (DEDUPLICATED)

=== Demo Completed Successfully ===
```
