# Code Audit

## Finding 1
Location: internal/server/server.go:70
Claimed Behavior: Implementation notes state "Wait durations in server mock are fixed (20ms), making the tail latency strictly a function of queuing time when VUs exceed max DB connections."
Observed Implementation: The server has a 10% chance of applying 25x query duration (500ms) when activeReq exceeds MaxDBConnections. This random slow-query behavior is not mentioned in the implementation notes.
Assessment: WARNING
Severity: LOW
Notes: The random slow query simulates realistic DB-side degradation under load, which is reasonable behavior. However, the implementation notes explicitly claim tail latency is "strictly a function of queuing time," which is inaccurate. This is a DOC_CODE_MISMATCH between engineering documentation and actual code behavior.

## Finding 2
Location: internal/server/server.go:69
Claimed Behavior: Server simulates DB connection pool saturation via semaphore
Observed Implementation: activeReq is incremented before acquiring the semaphore, then used to check if the pool is saturated. The check `activeReq > MaxDBConnections` correctly identifies when requests are queued.
Assessment: PASS
Severity: LOW
Notes: The semaphore pattern correctly limits concurrent DB access. The activeReq counter provides an additional signal for over-capacity detection.

## Finding 3
Location: internal/loadtest/runner.go:53,68
Claimed Behavior: Load tester spawns VUs and tracks metrics without race conditions
Observed Implementation: Each goroutine writes results to a pre-allocated slice at its unique index (`results[vuID]`). No mutex is needed since each index is written by exactly one goroutine.
Assessment: PASS
Severity: LOW
Notes: The per-VU result pattern is a clean concurrency design. Aggregation happens after wg.Wait(), eliminating all data races.

## Finding 4
Location: internal/loadtest/metrics.go:44
Claimed Behavior: "Thread-safe latency collector" (per design doc line 33)
Observed Implementation: CalculateMetrics receives a fully aggregated slice. It is called from Run() after wg.Wait() completes, so no concurrent access occurs. The "thread-safe" description is technically accurate due to timing (called after all goroutines finish) but misleading - there is no locking because it's unnecessary.
Assessment: PASS
Severity: LOW
Notes: Not a defect. The metric computation is called only after all goroutines complete. The "thread-safe" label in design docs is not wrong, just unnecessary. No fix needed.

## Finding 5
Location: internal/loadtest/runner.go:76-77,86-87
Claimed Behavior: Context cancellation errors are not counted as errors; only genuine failures are counted
Observed Implementation: Both `http.NewRequestWithContext` and `r.client.Do` check `ctx.Err() == nil` before incrementing error count. This prevents context cancellations from inflating error counts.
Assessment: PASS
Severity: LOW
Notes: Correct implementation of context-aware error handling for load testing.

## Finding 6
Location: internal/server/server.go:62-65,77-81
Claimed Behavior: Server handles request cancellation gracefully during semaphore acquisition and query processing
Observed Implementation: Two select blocks with `<-r.Context().Done()` allow the handler to abort if the client disconnects. The deferred semaphore release in `<-s.semaphore` ensures no resource leak.
Assessment: PASS
Severity: LOW
Notes: Proper context propagation and resource cleanup. The TestServer_ContextCanceled test validates this behavior.

## Finding 7
Location: internal/loadtest/runner.go:31-34
Claimed Behavior: Client-side connection pooling does not mask server bottlenecks
Observed Implementation: Custom http.Transport with MaxIdleConns=1000 and MaxIdleConnsPerHost=1000 prevents client-side connection limits from being the bottleneck.
Assessment: PASS
Severity: LOW
Notes: Good implementation decision documented in implementation notes (line 24). Correct approach for isolating server-side bottlenecks.

## Finding 8
Location: internal/loadtest/metrics.go:71
Claimed Behavior: Percentile calculation accurately identifies P50, P90, P95, P99
Observed Implementation: Uses nearest-rank method: `idx := int(float64(len(sorted)-1) * (pct / 100.0))`. For 100 samples, this yields P50=50ms, P95=95ms, P99=99ms, matching the test expectations.
Assessment: PASS
Severity: LOW
Notes: The percentile calculation is mathematically correct for the nearest-rank method. Tests in metrics_test.go validate the computation.

## Finding 9
Location: cmd/demo/main.go
Claimed Behavior: Demo runs smoke test (low VUs, no queuing) then stress test (high VUs, queuing expected), and prints comparative metrics
Observed Implementation: Smoke test uses 2 VUs (below 5 DB connections); stress test uses 50 VUs (10x DB max). Both run for 2 seconds against /booking endpoint.
Assessment: PASS
Severity: LOW
Notes: Demo output confirms expected behavior: smoke P95=21.6ms (close to average), stress P95=1105ms (far exceeds average=612ms), demonstrating the tail latency masking effect.

## Finding 10
Location: internal/loadtest/runner.go:55-56
Claimed Behavior: Load test runs for exactly the specified duration
Observed Implementation: context.WithTimeout sets the deadline. VUs check ctx.Done() in their loop and exit cleanly.
Assessment: PASS
Severity: LOW
Notes: The duration is approximate due to loop timing, but this is standard for load testing tools. No issues.

## Finding 11
Location: internal/server/server.go:5 (import)
Claimed Behavior: math/rand used for random slow query simulation
Observed Implementation: math/rand.Float32() is safe for concurrent use per Go documentation. No manual seeding needed in Go 1.22+ (auto-seeded).
Assessment: PASS
Severity: LOW
Notes: No concurrency safety issue with rand usage.
