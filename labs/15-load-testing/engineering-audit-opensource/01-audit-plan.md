# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files: 
- cmd/demo/main.go
- internal/server/server.go
- internal/loadtest/runner.go
- internal/loadtest/metrics.go
Tests:
- tests/loadtest_test.go
- internal/loadtest/metrics_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: engineering/01-design.md, engineering/02-implementation-notes.md
Main Claims To Verify:
1. Load testing reveals boundaries, saturation, and degradation patterns
2. Average response time conceals tail latency spikes; percentiles (P95, P99) are necessary
3. Incremental test stages (smoke vs stress) differentiate baseline performance from resource exhaustion
4. Downstream resource saturation causes non-linear latency degradation for tail requests
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in concurrent metrics collection
- Incorrect percentile calculations
- Demo not showing expected smoke vs stress differences
- Tests not validating core behavior claims