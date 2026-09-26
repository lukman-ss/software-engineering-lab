# Docs vs Code Audit

Target Lab: labs/15-load-testing

## Comparison Summary

| Item | Documentation (`README.md`, `01-design.md`) | Code Implementation (`internal/`, `cmd/`) | Audit Result |
| :--- | :--- | :--- | :--- |
| **Demo Command** | `go run ./cmd/demo` | `cmd/demo/main.go` exists and executes cleanly | MATCH |
| **Test Commands** | `go test -v ./...`, `go test -race ./...` | All packages pass standard and race test execution | MATCH |
| **Structure Claim** | Mentions `cmd/demo`, `internal/server`, `internal/loadtest`, `tests`, `engineering/` | All listed directories exist and contain specified logic | MATCH |
| **Resource Bottleneck** | Max DB connection pool constraint with queuing latency degradation | `internal/server/server.go` implements channel semaphore pool | MATCH |
| **Metrics Calculated** | Min, Max, Avg, P50, P90, P95, P99, RPS | `internal/loadtest/metrics.go` calculates all 8 metrics | MATCH |

## Identified Mismatches
None.

## Verification of Demo Output
Execution of `go run ./cmd/demo` produced actual output matching claims in `03-execution-result.md`:
- Smoke test: 2 VUs, 0 errors, P95 ≈ 21.4ms
- Stress test: 50 VUs against 5 DB connection slots, 0 errors, P95 ≈ 770.7ms (demonstrating queue tail latency expansion)
- No hardcoded or fake benchmark values found.
