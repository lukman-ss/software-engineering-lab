# Code Audit

## Finding 1

Location: cmd/demo/main.go:30-58 (smoke 2 VUs vs stress 50 VUs)
Claimed Behavior: Smoke test runs with low VUs below the connection-pool limit (no queuing); stress test runs with 10x the DB connection limit, causing queue saturation and tail-latency growth.
Observed Implementation: `smokeCfg.VUs = 2`, `stressCfg.VUs = 50`, server `MaxDBConnections = 5`. Stress VUs (50) greatly exceed the 5-slot semaphore, so requests queue. Smoke VUs (2) stay under 5, so no queuing. Demo output reproduced: Smoke P95 ~22ms, Stress P95 ~710ms, Stress P99 ~1.15s.
Assessment: PASS
Severity: LOW
Notes: Matches research Finding 1 (smoke vs stress) and Finding 4 (downstream saturation bottleneck). Invariant reproduced in live run.

## Finding 2

Location: internal/loadtest/metrics.go:67-73
Claimed Behavior: Loads report P50, P90, P95, P99 percentiles in addition to Min/Max/Avg.
Observed Implementation: `percentile()` uses `idx = int(float64(len(sorted)-1) * (pct/100.0))` (lower/nearest-rank with floor). All four percentiles computed and exposed on `Result`. `metrics_test.go` asserts exact values for 100-sample 1..100ms series (P50=50ms, P95=95ms, P99=99ms) — all pass.
Assessment: PASS
Severity: LOW
Notes: Deterministic lower-interpolation percentile (valid, not linear-interpolation like k6). Self-consistent with its own tests; acceptable for lab scale. Document the definition if aligning with k6's p(N) is required.

## Finding 3

Location: internal/loadtest/runner.go:91-98
Claimed Behavior: Latency is recorded per request.
Observed Implementation: Latency (`time.Since(reqStart)`) is appended to `lats` ONLY in the success branch (`else`); on HTTP>=400 or dial transport errors the request duration is discarded. `latencies` is therefore strictly successful-request latencies. `TotalRequests = len(latencies) + errors`.
Assessment: WARNING
Severity: MEDIUM
Notes: Percentile response-time metrics (research Finding 2) normally include all requests (errors still take time). Excluding error durations biases P50/P95/P99 downward for error-heavy runs. In the demo all responses are 201, so Smoke/Stress percentiles are unaffected; the gap only matters when errors carry meaningful latency. Acceptable scoping but should be explicit.

## Finding 4

Location: internal/loadtest/runner.go:48-114 (Run)
Claimed Behavior: Virtual users execute concurrently and aggregate safely; no race conditions.
Observed Implementation: Per-VU `results[vuID]` slice is written by exactly one goroutine (indexed by VU id), read only after `wg.Wait()`. No shared mutable state during the run. `go test -race ./...` is clean for both packages with tests.
Assessment: PASS
Severity: LOW
Notes: Confirmed by -race. Per-goroutine buffers avoid lock contention as documented in engineering/02-implementation-notes.md.

## Finding 5

Location: internal/server/server.go:56-93 (handleBooking)
Claimed Behavior: Server enforces a bounded DB connection pool; requests saturate and queue; context cancellation aborts in-flight acquires.
Observed Implementation: Semaphore `chan struct{}` of size `MaxDBConnections`. Acquire via `select { case s.semaphore <- struct{}{}: case <-r.Context().Done(): return }` with `defer func() { <-s.semaphore }()`. `TestServer_ContextCanceled` (pre-canceled context) — request aborts, no 201. `TestServer_MaxDBConnectionsBound` confirms peak ≤ max. `r.Context().Done()` also checked around the timer, aborting mid-query on cancellation.
Assessment: PASS
Severity: LOW
Notes: Acquire/release pairing is correct; nothing leaked on cancel. Verified by test + race.

## Finding 6

Location: internal/server/server.go:74-78
Claimed Behavior: Tail latency under stress is strictly a function of queuing time behind the connection pool. (mirrors engineering/02-implementation-notes.md "Tail latency is strictly a function of queuing time")
Observed Implementation: When `activeReq > MaxDBConnections`, 10% of requests multiply `DBQueryDuration` by 25 (→500ms vs 20ms baseline). This random slowdown is INDEPENDENT of semaphore queuing and is an additional tail-latency amplifier.
Assessment: WARNING
Severity: LOW
Notes: Engineering note is slightly inaccurate. Tail latency is queuing OR random-slowdown, not strictly queuing. Not a bug — the slowdown reinforces the queuing demonstration and the Smoke-vs-Stress invariant still holds — but it is a secondary, undocumented-in-code amplifier. Low risk; clarify the note.

## Finding 7

Location: cmd/demo/main.go:71-72 (printResults)
Claimed Behavior: Demo prints Total Requests, Success, Errors, RPS, Average, P50, P90, P95, P99.
Observed Implementation: `printResults` prints TotalRequests, SuccessCount, ErrorCount, RPS, AvgLatency, P50, P95, P99. P90 is computed in `Result` but NOT printed.
Assessment: WARNING
Severity: LOW
Notes: P90 available in Result struct but omitted from demo table. No correctness impact; cosmetic mismatch between documented metric set and printed set.

## Finding 8

Location: internal/server/server.go:5
Claimed Behavior: Standard library only, no third-party dependencies.
Observed Implementation: Imports: encoding/json, math/rand, net/http, sync/atomic, time. `go.mod` has zero require directives. `go build`/`go vet` clean.
Assessment: PASS
Severity: LOW
Notes: Matches engineering/design "stdlib only" decision.

## Finding 9

Location: runner.go (gofmt -l) / trailing whitespace
Claimed Behavior: Code is clean / formatted.
Observed Implementation: `gofmt -l .` reports `internal/loadtest/runner.go` (trailing tab on the blank line at 93). `go vet` and build are clean.
Assessment: WARNING
Severity: LOW
Notes: Cosmetic; does not affect correctness or tests. Fix: remove trailing whitespace.

## Finding 10 (concurrency safety summary)

Location: server.go activeReq; runner.go per-VU slices; server semaphore channel.
Claimed Behavior: All shared state is concurrency-safe.
Observed Implementation: `activeReq` is `atomic.Int64` (Add/Load). Server semaphore is a single buffered channel (no per-element state). Runner avoids shared counters by partitioning per VU. `-race` clean.
Assessment: PASS
Severity: LOW
Notes: No race conditions detected across the full suite.
