# Execution Result

## Build
Command:
```bash
go build ./...
```
Result:
```text
Exit Code: 0 (Successful compilation across all packages)
```

## Tests
Command:
```bash
go test -v ./...
```
Result:
```text
=== RUN   TestFaultInjector
--- PASS: TestFaultInjector (0.01s)
=== RUN   TestCircuitBreakerStateTransitions
--- PASS: TestCircuitBreakerStateTransitions (0.06s)
=== RUN   TestCircuitBreakerGracefulDegradation
--- PASS: TestCircuitBreakerGracefulDegradation (0.00s)
=== RUN   TestExperimentAutoAbortOnSteadyStateViolation
--- PASS: TestExperimentAutoAbortOnSteadyStateViolation (0.01s)
=== RUN   TestConcurrencyAndRace
--- PASS: TestConcurrencyAndRace (0.00s)
PASS
ok  	labs/34-chaos-engineering/tests	0.184s
```

## Race Detector
Command:
```bash
go test -race ./...
```
Result:
```text
ok  	labs/34-chaos-engineering/tests	1.481s
```

## Demo
Command:
```bash
go run ./cmd/demo
```
Result:
```text
=== Chaos Engineering & Fault Injection Demo ===

[1] Baseline Steady-State Traffic (Normal Operation)...
Req #1 -> Payment Success (Primary) | CB State: CLOSED
Req #2 -> Payment Success (Primary) | CB State: CLOSED
Req #3 -> Payment Success (Primary) | CB State: CLOSED
Req #4 -> Payment Success (Primary) | CB State: CLOSED
Req #5 -> Payment Success (Primary) | CB State: CLOSED

[2] Running Chaos Experiment: Injecting Downstream Failures with Resilient Client...
Req #6 -> Payment Queued (Fallback) | CB State: CLOSED
Req #7 -> Payment Queued (Fallback) | CB State: CLOSED
Req #8 -> Payment Queued (Fallback) | CB State: OPEN
Req #9 -> Payment Queued (Fallback) | CB State: OPEN
Req #10 -> Payment Queued (Fallback) | CB State: OPEN
Req #11 -> Payment Queued (Fallback) | CB State: OPEN
Req #12 -> Payment Queued (Fallback) | CB State: OPEN
Req #13 -> Payment Queued (Fallback) | CB State: OPEN
Req #14 -> Payment Queued (Fallback) | CB State: OPEN
Req #15 -> Payment Queued (Fallback) | CB State: OPEN

Experiment State: COMPLETED
Steady-State Metrics with Fallback: Total=15, Success=15, Failed=0, ErrorRate=0.00%

[3] Running Second Chaos Experiment: Unmitigated Downstream Fault (Triggers Auto-Abort)...
Raw Call #1 Failed: chaos: injected fault failure
Raw Call #2 Failed: chaos: injected fault failure
Experiment State: ABORTED
Abort Reason: Steady state breached: error rate 33.33% exceeded threshold
Fault Injector Active: false (Neutralized on Abort)

[4] Traffic After Experiment Ends / Recovery...
Req #16 -> Payment Success (Primary) | CB State: CLOSED
Req #17 -> Payment Success (Primary) | CB State: CLOSED
Req #18 -> Payment Success (Primary) | CB State: CLOSED
Req #19 -> Payment Success (Primary) | CB State: CLOSED
Req #20 -> Payment Success (Primary) | CB State: CLOSED

=== Demo Complete ===
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
