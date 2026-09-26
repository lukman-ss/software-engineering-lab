# Engineering Audit Plan

Target Lab: labs/15-load-testing

Implementation Files:
- internal/loadtest/runner.go (VU orchestrator, 115 lines)
- internal/loadtest/metrics.go (percentile calculator, 73 lines)
- internal/server/server.go (constrained mock server, 89 lines)
- cmd/demo/main.go (smoke vs stress demo, 72 lines)

Tests:
- internal/loadtest/metrics_test.go (3 unit tests)
- tests/loadtest_test.go (5 integration tests)

Executable/Demo:
- cmd/demo (httptest-backed smoke 2 VUs + stress 50 VUs, 2s each)

Approved Research Inputs: OUT OF SCOPE per pipeline override (implementation + tests only; engineering/ used as claim source, not audited as research)

Main Claims To Verify:
1. Runner spawns N VUs concurrently, aggregates latencies without lock contention, counts HTTP >=400 and transport errors, respects Duration timeout
2. Metrics computes Min/Max/Avg/P50/P90/P95/P99 accurately via exact sort
3. Server semaphore (MaxDBConnections) causes queueing + tail spike under stress (P95/P99 >> smoke)
4. Smoke run: 0 errors, P95 ~= Avg; Stress run: P95 > Avg and P95 >> smoke P95
5. No data races; demo output real; README matches code

Commands To Run:
- go build ./...
- go vet ./...
- go test -v -count=1 ./...
- go test -race -count=1 ./...
- go run ./cmd/demo

Primary Risks:
- Timing-sensitive assertion (P95 > Avg, stress P95 > smoke P95) flaky on loaded CI
- 10% rand-based spike (math/rand) nondeterministic by design
- Percentile uses floor-index nearest-rank, not interpolated — approximation risk
- Zero/negative Duration edge untested (immediate cancel path)
