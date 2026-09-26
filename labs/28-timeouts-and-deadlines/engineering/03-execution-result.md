# Execution Result

## Build
Command: `go build ./...`
Result: SUCCESS

## Tests
Command: `go test ./...`
Result:
```text
?   	timeouts-and-deadlines/cmd/demo	[no test files]
ok  	timeouts-and-deadlines/internal/circuit	0.469s
ok  	timeouts-and-deadlines/internal/deadline	0.433s
ok  	timeouts-and-deadlines/internal/idempotency	0.498s
ok  	timeouts-and-deadlines/internal/retry	0.367s
ok  	timeouts-and-deadlines/tests	0.372s
```

## Race Detector
Command: `go test -race ./...`
Result:
```text
?   	timeouts-and-deadlines/cmd/demo	[no test files]
ok  	timeouts-and-deadlines/internal/circuit	1.439s
ok  	timeouts-and-deadlines/internal/deadline	1.490s
ok  	timeouts-and-deadlines/internal/idempotency	1.475s
ok  	timeouts-and-deadlines/internal/retry	1.381s
ok  	timeouts-and-deadlines/tests	1.454s
```

## Demo
Command: `go run ./cmd/demo`
Result:
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

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
