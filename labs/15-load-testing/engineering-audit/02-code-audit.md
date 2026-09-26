# Code Audit

Target Lab: labs/15-load-testing

## Finding 1

Location: `internal/loadtest/runner.go:48-53`, `runner.go:107-113`
Claimed Behavior: Thread-safe, low-contention aggregation of concurrent request results across virtual users (VUs).
Observed Implementation: Each VU worker writes solely to its allocated `results[vuID]` slice without sharing mutexes or atomic variables during iteration. Results are combined after `wg.Wait()`.
Assessment: PASS
Severity: LOW
Notes: Pattern eliminates lock contention during load generation.

## Finding 2

Location: `internal/server/server.go:60-66`
Claimed Behavior: Simulated connection pool bounding concurrency and handling client context cancellations without leaking semaphore tokens.
Observed Implementation: Buffered channel semaphore of size `cfg.MaxDBConnections`. When acquiring slot, `select` checks `s.semaphore <- struct{}{}` and `<-r.Context().Done()`. Release is deferred right after acquisition (`defer func() { <-s.semaphore }()`).
Assessment: PASS
Severity: LOW
Notes: Properly avoids channel token leak if context is canceled before acquiring semaphore.

## Finding 3

Location: `internal/server/server.go:68-81`
Claimed Behavior: Tail latency degradation under overload with context cancellation support during DB query sleep.
Observed Implementation: Uses `time.NewTimer` with `defer t.Stop()` and `select` listening to both `t.C` and `r.Context().Done()`. Under overload (`activeReq > MaxDBConnections`), random 10% tail penalty simulates slow queries.
Assessment: PASS
Severity: LOW
Notes: Timer cleanup prevents timer leaks on canceled requests.

## Finding 4

Location: `internal/loadtest/metrics.go:23-42`, `metrics.go:67-73`
Claimed Behavior: Accurate calculation of percentiles and summary statistics, safe on empty slices.
Observed Implementation:
- Handled empty latencies and zero duration gracefully without dividing by zero.
- Percentile index uses standard nearest rank: `idx := int(float64(len(sorted)-1) * (pct / 100.0))`.
- Copies slice before sorting to prevent mutating input slice.
Assessment: PASS
Severity: LOW
Notes: Math and bounds checks verified.

## Finding 5

Location: `internal/loadtest/runner.go:85-89`
Claimed Behavior: Don't miscount expected test cancellations at end of duration as server errors.
Observed Implementation: Inspects `ctx.Err() == nil` before incrementing `errs` on request failures.
Assessment: PASS
Severity: LOW
Notes: Distinguishes natural timeout cutoff from genuine server/network errors.
