## Finding 1

Location: internal/server/server.go:61-66, 69-73
Claimed Behavior: Server simulates database connection pool with bounded semaphore; when active requests exceed MaxDBConnections, there's a 10% chance to increase query duration 25x to simulate degradation under load.
Observed Implementation: 
- Uses buffered channel as semaphore for DB connection slots (line 41)
- Tracks active requests with atomic counter (line 19)
- When semaphore acquire succeeds, checks if activeReq > MaxDBConnections (line 69)
- If true and random < 0.10, multiplies DBQueryDuration by 25 (lines 70-72)
- Otherwise uses normal duration
Assessment: PASS
Severity: 
Notes: Implementation correctly models connection pool exhaustion. When all slots are full (activeReq >= MaxDBConnections), new requests still get a slot but experience increased latency 10% of the time. This creates the tail latency under stress.

## Finding 2

Location: internal/loadtest/runner.go:44-112
Claimed Behavior: Load generator spawns VUs as goroutines, each making HTTP requests until context cancellation, collecting latencies and errors, then aggregating results.
Observed Implementation:
- Uses WaitGroup for goroutine management (line 45)
- Per-VU result slices to avoid lock contention (lines 48-53)
- Context with timeout for test duration (line 55)
- Each VU loops until ctx.Done() (lines 65-98)
- Proper error handling: counts errors only when context not cancelled (lines 84-86)
- Status code >= 400 counted as error (lines 92-94)
- Latency measured from request start to response body read (line 95)
Assessment: PASS
Severity: 
Notes: Well-implemented concurrent load generator with proper context handling and error classification. Per-VU slicing reduces lock contention.

## Finding 3

Location: internal/loadtest/metrics.go:23-65
Claimed Behavior: CalculateMetrics computes total requests, success/error counts, RPS, min/max/avg latency, and percentiles (P50, P90, P95, P99).
Observed Implementation:
- Sorts latency slice to compute percentiles (lines 45-49)
- Uses standard percentile formula: idx = int(float64(len-1) * pct/100) (lines 67-72)
- Computes average by summing and dividing (lines 51-55)
- Min/max from sorted array (lines 56-57)
- RPS = totalRequests / duration.Seconds() (lines 36-38)
Assessment: PASS
Severity: 
Notes: Correct implementation of statistical calculations. Note the "ponytail" comment acknowledging exact sort is fine for test scale but would need histogram for larger scale.

## Finding 4

Location: internal/loadtest/runner.go:31-41
Claimed Behavior: Custom http.Transport with large idle conn limits prevents load generator from becoming bottleneck.
Observed Implementation:
- Sets MaxIdleConns and MaxIdleConnsPerHost to 1000 (lines 32-33)
- Uses 5 second timeout (line 39)
Assessment: PASS
Severity: 
Notes: Properly configures HTTP client to avoid self-throttling, ensuring the server is the bottleneck.

## Finding 5

Location: internal/loadtest/runner.go:80-96
Claimed Behavior: Latency measurement captures full round-trip time including server processing.
Observed Implementation:
- Records time before Do() call (line 81)
- Measures time.Since() after response body read and close (line 95)
- Actually: reqStart is set before Do(), latency calculated after resp.Body.Close()
Assessment: PASS
Severity: 
Notes: Latency measurement includes network round trip, server processing, and response transfer - correct for end-to-end measurement.

## Finding 6

Location: internal/server/server.go:57-58
Claimed Behavior: Active request counter tracks concurrent requests for degradation logic.
Observed Implementation:
- atomic.AddInt64(&s.activeReq, 1) at start of handler (line 57)
- defer atomic.AddInt64(&s.activeReq, -1) to decrement (line 58)
Assessment: PASS
Severity: 
Notes: Correct use of atomic for concurrent increment/decrement. However, note that activeReq includes requests waiting for semaphore AND those being processed. This is correct for the degradation logic which triggers when activeReq > MaxDBConnections (meaning some requests are waiting).

## Finding 7

Location: internal/loadtest/runner.go:71-74
Claimed Behavior: HTTP request creation error handling.
Observed Implementation:
- If http.NewRequestWithContext fails, increments error count and continues (lines 72-74)
Assessment: WARNING
Severity: LOW
Notes: Should distinguish between context cancellation (ctx.Err() != nil) and other errors. Currently counts all request creation errors as test errors, which could be misleading if context is cancelled. However, in practice request creation rarely fails unless severely misconfigured.

## Finding 8

Location: internal/loadtest/runner.go:82-88
Claimed Behavior: HTTP client Do() error handling.
Observed Implementation:
- If err != nil from r.client.Do(), checks if ctx.Err() == nil before counting as error (lines 83-86)
Assessment: PASS
Severity: 
Notes: Correctly excludes context cancellation errors from error count.

## Finding 9

Location: internal/loadtest/runner.go:92-94
Claimed Behavior: HTTP status codes >= 400 are counted as errors.
Observed Implementation:
- if resp.StatusCode >= 400 { errs++ } (lines 92-94)
Assessment: PASS
Severity: 
Notes: Standard practice to count 4xx and 5xx as errors.

## Finding 10

Location: internal/loadtest/runner.go:105-110
Claimed Behavior: After WaitGroup, aggregates all VU latencies and errors into final result.
Observed Implementation:
- Loops through results, appends latencies, sums errors (lines 107-110)
- Calls CalculateMetrics with aggregated data (line 112)
Assessment: PASS
Severity: 
Notes: Correct aggregation of per-VU results.

## Finding 11

Location: cmd/demo/main.go
Claimed Behavior: Demo runs smoke test (2 VUs) then stress test (50 VUs) against server with 5 max DB connections, 20ms query duration, 2-second duration each.
Observed Implementation:
- Smoke: 2 VUs < 5 max connections → no queuing expected
- Stress: 50 VUs >> 5 max connections → significant queuing expected
- Prints comparative tables showing latency increase
Assessment: PASS
Severity: 
Notes: Demo correctly configured to show the contrast between baseline and saturated conditions.