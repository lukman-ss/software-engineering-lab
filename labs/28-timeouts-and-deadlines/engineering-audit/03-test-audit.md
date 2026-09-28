# Test Audit

Target Lab: labs/28-timeouts-and-deadlines

## Test Execution Results

Command: `go test -count=1 -v ./...`
Output:
```text
=== RUN   TestCircuitBreaker_StateTransitions
--- PASS: TestCircuitBreaker_StateTransitions (0.06s)
=== RUN   TestCircuitBreaker_HalfOpenFailureTripsOpen
--- PASS: TestCircuitBreaker_HalfOpenFailureTripsOpen (0.04s)
=== RUN   TestCircuitBreaker_DefaultZeroConfig
--- PASS: TestCircuitBreaker_DefaultZeroConfig (0.00s)
PASS
ok  	timeouts-and-deadlines/internal/circuit	0.183s
=== RUN   TestExecuteWithBudget_Success
--- PASS: TestExecuteWithBudget_Success (0.00s)
=== RUN   TestExecuteWithBudget_Timeout
--- PASS: TestExecuteWithBudget_Timeout (0.02s)
=== RUN   TestExecuteWithBudget_ParentTimeoutInherited
--- PASS: TestExecuteWithBudget_ParentTimeoutInherited (0.02s)
PASS
ok  	timeouts-and-deadlines/internal/deadline	0.123s
=== RUN   TestStore_GetSet
--- PASS: TestStore_GetSet (0.12s)
=== RUN   TestStore_ConcurrentAccess
--- PASS: TestStore_ConcurrentAccess (0.00s)
=== RUN   TestStore_LazyEvictionOnGet
--- PASS: TestStore_LazyEvictionOnGet (0.03s)
PASS
ok  	timeouts-and-deadlines/internal/idempotency	0.233s
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
ok  	timeouts-and-deadlines/internal/retry	0.100s
=== RUN   TestIntegration_RetryWithCircuitBreaker
--- PASS: TestIntegration_RetryWithCircuitBreaker (0.02s)
=== RUN   TestIntegration_IdempotentRetry
--- PASS: TestIntegration_IdempotentRetry (0.00s)
PASS
ok  	timeouts-and-deadlines/tests	0.095s
```

Command: `go test -count=1 -race ./...`
Output:
```text
ok  	timeouts-and-deadlines/internal/circuit	1.211s
ok  	timeouts-and-deadlines/internal/deadline	1.395s
ok  	timeouts-and-deadlines/internal/idempotency	1.261s
ok  	timeouts-and-deadlines/internal/retry	1.137s
ok  	timeouts-and-deadlines/tests	1.128s
```

## Test Coverage Evaluation

- Happy path: Covered across all packages and integration tests.
- Failure path: Covered for retry exhaustion, breaker open trips, and deadline expirations.
- Concurrency & Race conditions: Explicitly tested (`TestStore_ConcurrentAccess`) and verified clean under `-race`.
- Edge cases: Default zero config fallbacks and jitter bounding tested.
- Integration: End-to-end integration tests verify retry + circuit breaker and retry + idempotency deduplication.
