# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files: internal/server/server.go, internal/loadtest/runner.go, internal/loadtest/metrics.go
Tests: tests/loadtest_test.go, internal/loadtest/metrics_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: research/05-report.md
Main Claims To Verify:
1. Load testing reveals system boundaries and saturation.
2. P95/P99 percentiles expose latency outliers masked by averages.
3. Constrained resources (connection pool) cause non-linear tail latency degradation.
Commands To Run:
- `go test -count=1 -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions in load generator.
- Inaccurate percentile calculation.
- Tests passing by coincidence without asserting degradation.
