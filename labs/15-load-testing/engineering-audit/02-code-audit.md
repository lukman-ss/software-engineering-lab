# Code Audit

## Finding 1

Location: `internal/server/server.go` (lines 59-62)
Claimed Behavior: Simulates connection pool bottleneck via semaphore.
Observed Implementation: Uses buffered channel to bound concurrent DB query sleeps.
Assessment: PASS
Severity: LOW
Notes: Functionally sound for simulating bounded concurrency. The channel capacity accurately restricts concurrent throughput and induces expected queuing behavior.

## Finding 2

Location: `internal/server/server.go` (line 60)
Claimed Behavior: Handle incoming booking requests.
Observed Implementation: Server blocks unconditionally on `s.semaphore <- struct{}{}`.
Assessment: WARNING
Severity: LOW
Notes: If the HTTP request is canceled by the client while waiting in the queue, the server will continue to hold the semaphore slot once acquired. Safe for this lab environment since test runs are short and timeouts are high, but suboptimal for production resilience.

## Finding 3

Location: `internal/loadtest/runner.go` (lines 48-52)
Claimed Behavior: Thread-safe metrics collection.
Observed Implementation: Pre-allocates `results` slice; each goroutine writes exclusively to its own index (`results[vuID]`). Merged sequentially after `wg.Wait()`.
Assessment: PASS
Severity: LOW
Notes: Avoids lock contention entirely. Clean, robust implementation of concurrent aggregation.

## Finding 4

Location: `internal/loadtest/runner.go` (line 88)
Claimed Behavior: Close HTTP response body to prevent leaks.
Observed Implementation: `_ = resp.Body.Close()`
Assessment: WARNING
Severity: LOW
Notes: Not reading the body before closing could prevent HTTP keep-alive connection reuse if the body is unread. However, since the simulated server returns a tiny JSON payload (read entirely into buffers by the TCP stack), `http.Transport` connection reuse is not materially affected in this lab scope.

## Finding 5

Location: `internal/loadtest/runner.go` (line 83)
Claimed Behavior: Filter out context cancellations at test completion.
Observed Implementation: `if ctx.Err() == nil { errs++ }`
Assessment: PASS
Severity: LOW
Notes: Excellent handling of shutdown dynamics. Prevents false positive errors when the test duration elapses and cleanly aborts in-flight requests.

## Finding 6

Location: `internal/loadtest/metrics.go` (lines 67-73)
Claimed Behavior: Calculates accurate latency percentiles.
Observed Implementation: Performs exact array sorting and retrieves nearest-rank index `idx := int(float64(len(sorted)-1) * (pct / 100.0))`.
Assessment: PASS
Severity: LOW
Notes: Standard rank-estimation approach. Appropriate for the dataset sizes used in the lab. `TotalRequests` and edge cases (empty slices) are handled safely.
