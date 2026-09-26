# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files:
- labs/15-load-testing/internal/server/server.go
- labs/15-load-testing/internal/loadtest/runner.go
- labs/15-load-testing/internal/loadtest/metrics.go
- labs/15-load-testing/cmd/demo/main.go
Tests:
- labs/15-load-testing/internal/loadtest/metrics_test.go
- labs/15-load-testing/tests/loadtest_test.go
Executable/Demo: labs/15-load-testing/cmd/demo/main.go
Approved Research Inputs: 
- labs/15-load-testing/engineering/01-design.md
- labs/15-load-testing/engineering/03-execution-result.md
Main Claims To Verify:
1. Code compiles successfully
2. All tests pass (unit and integration)
3. No race conditions detected with -race flag
4. Demo output shows Smoke Test vs Stress Test latency divergence (P95/P99)
5. Load test harness correctly calculates percentiles (P50, P95, P99)
6. Server simulates connection pool exhaustion correctly via semaphore
7. README matches implementation (structure, running instructions)
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Inaccurate percentile calculation due to sorting approach (acceptable for scale <1M samples)
- Context cancellation handling in load test runner may undercount errors
- Server's activeReq increment/decrement not atomic in all code paths? Actually it is atomic via AddInt64.
- Potential resource leak if context canceled while holding semaphore? The defer func() { <-s.semaphore }() is after the semaphore acquire but before the timer; if context.Done() fires in the select, we return without releasing the semaphore. This is a bug.