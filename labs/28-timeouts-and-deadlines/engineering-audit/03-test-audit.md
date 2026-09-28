# Test Audit

Target Lab: `labs/28-timeouts-and-deadlines`

## Test Execution Results

```text
go test ./...
?       timeouts-and-deadlines/cmd/demo [no test files]
ok      timeouts-and-deadlines/internal/circuit 0.400s
ok      timeouts-and-deadlines/internal/deadline        0.378s
ok      timeouts-and-deadlines/internal/idempotency     0.457s
ok      timeouts-and-deadlines/internal/retry   0.365s
ok      timeouts-and-deadlines/tests    0.347s

go test -race ./...
ok      timeouts-and-deadlines/internal/circuit 1.163s
ok      timeouts-and-deadlines/internal/deadline        1.141s
ok      timeouts-and-deadlines/internal/idempotency     1.222s
ok      timeouts-and-deadlines/internal/retry   1.129s
ok      timeouts-and-deadlines/tests    1.112s
```

## Coverage by Domain

### 1. `internal/deadline`
- Happy path: Fast function execution within budget (`TestExecuteWithBudget_Success`).
- Failure path: Function exceeding budget returns `context.DeadlineExceeded` (`TestExecuteWithBudget_Timeout`).
- Context cancellation propagation: Pre-canceled and canceled parent context propagates immediately (`TestExecuteWithBudget_ParentCancel`).

### 2. `internal/retry`
- Happy path: First attempt success without delay (`TestRetrier_SuccessFirstAttempt`).
- Failure path: Exhaustion of `MaxAttempts` returning combined errors (`TestRetrier_MaxRetriesExceeded`).
- Recovery path: Success on subsequent attempt (`TestRetrier_EventualSuccess`).
- Backoff bounds: Verifies jittered delays are within theoretical `[0, min(max, base*2^(i-1))]` bounds (`TestRetrier_CalculateBackoffBounds`).
- Context cancellation: Immediate abort upon context cancellation (`TestRetrier_ContextCancellation`).

### 3. `internal/circuit`
- Happy path: Calls allowed in `CLOSED` state (`TestBreaker_InitialClosed`).
- Transitions: Threshold failure transitions to `OPEN` (`TestBreaker_TripToOpen`).
- Rejection: Fast rejection with `ErrCircuitOpen` while `OPEN` (`TestBreaker_RejectWhenOpen`).
- Half-Open Cooldown: Cooldown expiration allows trial call (`TestBreaker_HalfOpenTransition`).
- Recovery: Consecutive successes in `HALF_OPEN` restore `CLOSED` state (`TestBreaker_HalfOpenToClosed`).
- Regression: Failure in `HALF_OPEN` immediately trips back to `OPEN` (`TestBreaker_HalfOpenFailureTripsToOpen`).
- Concurrency: Parallel calls with race detector (`TestBreaker_ConcurrentAccess`).

### 4. `internal/idempotency`
- Basic operations: Store and retrieve response key (`TestStore_SetAndGet`).
- Missing key: Returns not found (`TestStore_GetNotFound`).
- Expiration / TTL: Records expire past configured TTL (`TestStore_TTLExpiration`).
- Concurrency: Parallel writes and reads tested under `-race` (`TestStore_ConcurrentAccess`).

### 5. `tests/integration_test.go`
- End-to-end composite pipeline: Deadline budget + Exponential backoff + Circuit breaker + Idempotency deduplication interacting together (`TestEndToEnd_ResiliencePipeline`).

Assessment: PASS
All critical failure modes, edge cases, state transitions, and concurrency race conditions are covered with deterministic assertions.
