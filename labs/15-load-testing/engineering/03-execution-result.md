# Execution Result

## Build
Command:
```bash
go build ./...
```
Result:
```text
(success, no output)
```

## Tests
Command:
```bash
go test -v ./...
```
Result:
```text
?   	github.com/lukman/software-engineering-lab/labs/15-load-testing/cmd/demo	[no test files]
=== RUN   TestCalculateMetrics
--- PASS: TestCalculateMetrics (0.00s)
=== RUN   TestCalculateMetrics_Empty
--- PASS: TestCalculateMetrics_Empty (0.00s)
=== RUN   TestCalculateMetrics_Invariants
--- PASS: TestCalculateMetrics_Invariants (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest	0.347s
?   	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/server	[no test files]
=== RUN   TestLoadTest_SmokeVsStress
--- PASS: TestLoadTest_SmokeVsStress (1.36s)
=== RUN   TestLoadTest_ErrorCount
--- PASS: TestLoadTest_ErrorCount (0.10s)
=== RUN   TestServer_MethodNotAllowed
--- PASS: TestServer_MethodNotAllowed (0.00s)
=== RUN   TestLoadTest_DialError
--- PASS: TestLoadTest_DialError (0.10s)
=== RUN   TestServer_ContextCanceled
--- PASS: TestServer_ContextCanceled (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/tests	1.657s
```

## Race Detector
Command:
```bash
go test -race ./...
```
Result:
```text
?   	github.com/lukman/software-engineering-lab/labs/15-load-testing/cmd/demo	[no test files]
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest	1.385s
?   	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/server	[no test files]
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/tests	2.448s
```

## Demo
Command:
```bash
go run ./cmd/demo
```
Result:
```text
Starting Load Test Demo (Booking Bengkel)
Server simulated DB connections: 5
Simulated DB query duration: 20ms

--- Running Smoke Test (2 VUs) ---
Total Requests:  188
Success:         188
Errors:          0
RPS:             93.96
Average:         21.250847ms
P50:             21.2175ms
P95:             21.372291ms
P99:             22.289875ms

--- Running Stress Test (50 VUs) ---
Total Requests:  120
Success:         120
Errors:          0
RPS:             59.92
Average:         739.091366ms
P50:             669.568542ms
P95:             1.35666975s
P99:             1.587917041s
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
