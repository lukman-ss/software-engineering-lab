## Finding 1

Location: `internal/loadtest/runner.go`
Claimed Behavior: Thread-safe per-VU aggregation without mutexes.
Observed Implementation: Uses pre-allocated slice `results = make([]vuResult, VUs)`, writes to disjoint indices `results[vuID]`. Safe.
Assessment: PASS
Severity: LOW
Notes: Perfect concurrent design for isolation.

## Finding 2

Location: `internal/server/server.go`
Claimed Behavior: Simulates connection pool bottleneck and non-linear degradation.
Observed Implementation: Bounded channel `semaphore` creates queuing. `atomic.AddInt64` measures total active (queued + executing). If `activeReq > max`, 10% of requests sleep 25x longer.
Assessment: PASS
Severity: LOW
Notes: Accurate simulation of saturated resource cascading failure.

## Finding 3

Location: `internal/loadtest/metrics.go`
Claimed Behavior: Accurately computes percentiles (P50, P95, P99).
Observed Implementation: Nearest-rank exact sorting `idx := int(float64(len(sorted)-1) * (pct / 100.0))`.
Assessment: PASS
Severity: LOW
Notes: Adequate for test limits. Memory scales linearly but explicitly noted in trade-offs.