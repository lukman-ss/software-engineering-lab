# Gap Analysis

Target Lab: labs/15-load-testing

## DOC_CODE_MISMATCH
- **Description:** `01-design.md` claims the failure scenario involves "timeouts for the 95th percentile". However, `cmd/demo/main.go` runs for 2 seconds while the HTTP client timeout in `internal/loadtest/runner.go` is 5 seconds. Timeouts do not occur (Errors: 0 in stress test).
- **Severity:** LOW
- **Assessment:** The primary mechanism (latency degradation due to saturation) is accurately simulated and proven. The missing timeouts do not invalidate the lab's core lesson on percentiles vs averages.

## All Other Quality Gates
- **Race Conditions:** None. Concurrency model is robust.
- **Test Coverage:** High. Explicit coverage of latency invariants.
- **Implementation Alignment:** Hand-rolled exact percentile sorting correctly implemented and explicitly traded-off.
- **Code Correctness:** Channel-based semaphore successfully mimics pool exhaustion without CPU spinning.