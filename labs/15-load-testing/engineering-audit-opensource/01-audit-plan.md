# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files:
- cmd/demo/main.go
- internal/server/server.go
- internal/loadtest/metrics.go
- internal/loadtest/runner.go
Tests:
- tests/loadtest_test.go
- internal/loadtest/metrics_test.go
Executable/Demo: go run ./cmd/demo
Approved Research Inputs: research/05-report.md (Research Report)
Main Claims To Verify:
1. Load testing reveals boundaries, saturation, and degradation patterns.
2. Average response time conceals tail latency spikes; percentiles (P95, P99) are necessary to uncover degradation.
3. Incremental test stages (smoke vs stress) differentiate baseline performance from resource exhaustion.
4. Downstream resource saturation (e.g., database connection pool exhaustion) causes non-linear latency degradation for tail requests.
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in concurrent load test implementation
- Incorrect percentile calculations
- Demo not matching documented expected behavior
- Tests not covering edge cases (error conditions, context cancellation, etc.)