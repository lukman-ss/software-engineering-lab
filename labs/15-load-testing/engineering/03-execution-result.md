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
PASS
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest	0.546s
?   	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/server	[no test files]
=== RUN   TestLoadTest_SmokeVsStress
--- PASS: TestLoadTest_SmokeVsStress (1.05s)
=== RUN   TestServer_MethodNotAllowed
--- PASS: TestServer_MethodNotAllowed (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/tests	1.600s
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
Total Requests:  182
Success:         182
Errors:          0
RPS:             90.96
Average:         21.873222ms
P50:             21.727792ms
P95:             22.674291ms
P99:             29.268042ms

--- Running Stress Test (50 VUs) ---
Total Requests:  470
Success:         470
Errors:          0
RPS:             234.92
Average:         200.484306ms
P50:             209.993833ms
P95:             212.389875ms
P99:             214.703708ms
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
