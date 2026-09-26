# Docs vs Code

## Claims

1.  **Docs Claim**: Average response time conceals tail latency spikes; percentiles are necessary.
    *   **Code**: Implementation tracks P50, P90, P95, P99 alongside Avg. Demo output visually confirms average is much lower than P95 under stress.
    *   **Assessment**: Match.

2.  **Docs Claim**: Server implements a mock database connection pool using a semaphore.
    *   **Code**: `server.go` line 62 (`s.semaphore <- struct{}{}`) controls concurrency.
    *   **Assessment**: Match.

3.  **Docs Claim**: Demo contrasts Smoke test with Stress test metrics.
    *   **Code**: `cmd/demo/main.go` runs Smoke (2 VUs) then Stress (50 VUs), outputting tab-writer tables.
    *   **Assessment**: Match.

## Findings
No mismatches. Code and tests perfectly fulfill the engineering design document and README.
