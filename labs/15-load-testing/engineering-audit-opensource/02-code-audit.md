# Code Audit

Lab: labs/15-load-testing

## Scope
Reviewed files: `internal/server/server.go`, `internal/loadtest/runner.go`, `internal/loadtest/metrics.go`, `cmd/demo/main.go`. Excluded: research/ and content/ (per pipeline override).

## Methodology
Each finding records the file/line, the claimed behavior, the observed implementation, and an assessment. Severity levels: LOW, MEDIUM, HIGH, CRITICAL. Claims are considered "verified" only when code + execution prove them.

---

## Finding 1

Location: internal/server/server.go:60-66
Claimed Behavior: Server simulates a database connection pool with a hard concurrency limit via a buffered-channel semaphore. Slots are released on request completion, including when the client aborts (context cancellation).
Observed Implementation: A `chan struct{}` of capacity `MaxDBConnections` is filled on handler entry (`s.semaphore <- struct{}{}`) and drained via `defer func(){ <-s.semaphore }()`. The defer is registered only after successful acquisition, so an aborted pre-acquire (context done in first select) returns without touching the semaphore. On abort after acquisition (timer select, line 77-81), the deferred drain fires.
Assessment: PASS
Severity: HIGH (resource correctness is core)
Notes: No semaphore leak on any return path. `activeReq` is managed with `atomic.AddInt64` and is consistent.

---

## Finding 2

Location: internal/server/server.go:69-73
Claimed Behavior: Under saturation the server models non-linear tail latency: 10% of requests queued past capacity take ~25x longer.
Observed Implementation: When `activeReq > MaxDBConnections`, 10% of iterations multiply `DBQueryDuration` by 25. Uses `math/rand.Float32()` (Go 1.20+ auto-seeded globally; go.mod pins 1.22).
Assessment: PASS
Severity: LOW
Notes: Randomness is acceptable for a mock; non-determinism is expected and acknowledged in docs. The multiplier is bounded and deterministic in structure.

---

## Finding 3

Location: internal/loadtest/runner.go:44-102
Claimed Behavior: Load generator runs N virtual users concurrently for a bounded duration, collecting latencies error-free (no data races) via per-VU slices.
Observed Implementation: Each goroutine writes exclusively to its own index `results[vuID]`; no shared mutable state during the run. Aggregation happens after `wg.Wait()`. Timeout uses an inner `context.WithTimeout`; context-cancellation errors are filtered by `ctx.Err() == nil` checks (lines 73, 86, 114) so cancelled/dialed requests are not double-counted.
Assessment: PASS
Severity: HIGH
Notes: `go test -race` passes (verified, see execution). The `io.Discard` body drain and `resp.Body.Close()` are always reached. The custom `http.Transport` (lines 31-34) correctly prevents client-side pooling from masking server saturation.

---

## Finding 4

Location: internal/loadtest/runner.go:94-98
Claimed Behavior: Only HTTP 2xx/3xx (status < 400) responses contribute latencies; >=400 responses are treated as errors.
Observed Implementation: `if resp.StatusCode >= 400 { errs++ } else { lats = append(...) }`.
Assessment: PASS
Severity: MEDIUM
Notes: This is a deliberate scoping decision: latency is measured only on "successful" round-trips. The implementation matches the documented "latency = request-to-response time on success." Acceptable and consistent with design.

---

## Finding 5

Location: internal/loadtest/metrics.go:23-65
Claimed Behavior: Latencies are sorted and percentiles (P50, P90, P95, P99) computed accurately; sum/count yields a correct Avg; Min/Max taken from the sorted extremes.
Observed Implementation: `CalculateMetrics` copies the input slice, sorts ascending, computes `idx = int((len-1) * pct/100)` (nearest-rank method). Avg via integer division of `sum` by count. Min = sorted[0], Max = sorted[len-1].
Assessment: PASS
Severity: MEDIUM
Notes: The percentile method is "nearest rank" (a valid, common definition) and is internally validated by `TestCalculateMetrics` and `TestCalculateMetrics_Invariants` which assert ordering invariants Min <= P50 <= P90 <= P95 <= P99 <= Max. The `ponytail:` comment correctly flags the scaling ceiling (sorting not suitable for millions of samples) and upgrade path. Matches design doc's stated "exact sorting" choice.

---

## Finding 6

Location: cmd/demo/main.go:16-59
Claimed Behavior: Demo runs a Smoke scenario (2 VUs) then a Stress scenario (50 VUs) against a server with 5 DB connections and 20ms query time, printing comparative metric tables.
Observed Implementation: Exactly as claimed. Uses `httptest.NewServer`, defers `ts.Close()`. Configs match the design doc. `printResults` formats the full Result struct.
Assessment: PASS
Severity: LOW
Notes: Demo is self-contained (no external services). Verified by actual execution (stress P95/P99 >> smoke P95/P99, invariant holds across runs).

---

## Finding 7

Location: internal/loadtest/metrics.go:36-38; cmd/demo/main.go
Claimed Behavior: RPS computed as total requests divided by wall-clock duration.
Observed Implementation: `res.RPS = float64(total) / totalDuration.Seconds()` where `total = len(latencies) + errors`.
Assessment: PASS
Severity: LOW
Notes: `totalDuration` in the runner is the actual wall-clock span from `start` to after `wg.Wait()`, so RPS reflects observed throughput. Consistent with the demo output.

---

## Finding 8

Location: internal/server/server.go:51-88
Claimed Behavior: Method validation rejects non-POST with 405; context cancellation aborts cleanly.
Observed Implementation: `r.Method != POST` returns 405 (line 52-55). Both the semaphore-acquire select and the timer select also honor `r.Context().Done()`.
Assessment: PASS
Severity: LOW
Notes: `TestServer_MethodNotAllowed` and `TestServer_ContextCanceled` confirm both behaviors.

---

## Summary

No CRITICAL or HIGH defects found. The semaphore lifecycle, per-VU isolation, percentile math, and context handling all behave as designed. Two LOW-accuracy terminology mismatches in docs (see 04-docs-vs-code) but the code itself is correct.