## Code Audit

## Finding 1: Summary uses write Lock instead of RLock
Location: internal/metrics/tracker.go:117-128 (Summary method)
Claimed Behavior: Summary is a read-only operation that should allow concurrent readers.
Observed Implementation: Summary uses w.mu.Lock() (write lock) which serializes all Summary calls (and Record calls), reducing read parallelism.
Assessment: WARNING
Severity: LOW
Notes: Correctness not affected (mutual exclusion still holds), but performance suboptimal under high-concurrency read workloads. Readers block readers unnecessarily.

## Finding 2: Eviction uses event timestamp, not wall clock
Location: internal/metrics/tracker.go:46 (Record calls evictStaleLocked(e.Timestamp)); line 106-115 (evictStaleLocked)
Claimed Behavior: WindowTracker should evict buckets older than windowSize from current time.
Observed Implementation: evictStaleLocked is called with e.Timestamp (the event’s timestamp). If an out-of-order event arrives with an old timestamp, cutoff = e.Timestamp.Add(-windowSize) may be too recent, causing truly stale buckets (based on wall clock) to NOT be evicted.
Assessment: WARNING
Severity: MEDIUM
Notes: Edge case: if real-time events are recorded then delayed old events arrive, old buckets may accumulate beyond the intended window. Not triggered in demo/tests (events near-monotonic). Could cause memory leak in pathological cases.

## Finding 3: Unused BurnRateRule fields
Location: internal/alerting/engine.go:17-24 (BurnRateRule struct)
Claimed Behavior: BurnRateRule configures alert thresholds for multi-window policy.
Observed Implementation: Fields LongWindow, ShortWindow, BudgetConsumedPct are defined but never used. Check() uses the shortTracker and longTracker passed to AlertEngine, ignoring the rule’s window durations. BudgetConsumedPct unused.
Assessment: WARNING
Severity: LOW
Notes: Dead code / configuration that does nothing. Rule names in demo embed window semantics (e.g., "6.0x - 5% in 6h") but windows come from tracker construction, not the rule. Minor unnecessary complexity.

## Finding 4: Test coverage not 100% (design claim)
Location: engineering/01-design.md line 21 ("100% test coverage on core math and sliding window calculations")
Claimed Behavior: Success criterion states 100% test coverage.
Observed Implementation: Actual statement coverage (via go test -coverpkg=./internal/... ./tests/) is 95.3%. Uncovered branches:
  - CalculateBurnRate: total == 0 branch; allowedErrorRate <= 0 branch.
  - NewWindowTracker: bucketSize <= 0 branch; numBuckets < 1 branch.
  - Record: small fraction (~3.1%) (likely out-of-order insertion edge path).
Assessment: WARNING
Severity: MEDIUM
Notes: Design overstates test coverage. Core behavior (SLI, error budget, burn rate, alerting, CanDeploy) is tested. Missing tests are edge-case argument validation and zero-input handling. Not a correctness bug but violates the stated success criterion.

## Finding 5: Design over-describes latency histogram
Location: engineering/01-design.md line 26 ("Sliding-window event recorder (histogram latency buckets & success counts)")
Claimed Behavior: Metrics recorder maintains latency histograms.
Observed Implementation: internal/metrics/tracker.go stores Duration per Event but does not bucket/aggregate durations into histogram. Only good/bad counts per time bucket are aggregated.
Assessment: WARNING
Severity: LOW
Notes: Minor inaccuracy in design notes; implementation correctly computes SLI via good/total using Duration threshold (isGoodEvent). No functional impact; README description is accurate ("recording requests and measuring good vs. total events").

## Finding 6: Missing Phase 4 in recorded execution result
Location: engineering/03-execution-result.md lines 59-65 (demo output)
Claimed Behavior: engineering/03-execution-result.md records actual demo output from `go run ./cmd/demo`.
Observed Implementation: Recorded demo output ends after [PHASE 3] and jumps to DEMO COMPLETE. Actual demo output (verified by `go run ./cmd/demo`) includes [PHASE 4] Endpoint Criticality Comparison showing Reports SLO 95% vs Payment 99.9%.
Assessment: WARNING
Severity: MEDIUM
Notes: Recorded execution result is stale/incomplete (omits Phase 4). Actual demo is correct and runs fully. Documentation mismatch between claimed execution output and real demo.

## Finding 7: Summary() evicts stale buckets before computing totals
Location: internal/metrics/tracker.go:117-128 (Summary method)
Claimed Behavior: Summary should return stats for the window at the given time.
Observed Implementation: Summary(now) calls evictStaleLocked(now) before iterating buckets, ensuring the window is cut to [now - windowSize, now] at the moment of the call.
Assessment: PASS
Severity: -
Notes: Correct behavior. Eviction uses the Summary timestamp (wall clock) not event timestamp, so window is accurate for the query time.

## Finding 8: CalculateBurnRate handles zero total and zero allowedErrorRate
Location: internal/alerting/engine.go:51-61 (CalculateBurnRate)
Claimed Behavior: Returns 0 when no traffic or when SLO=1 (no allowed error).
Observed Implementation: Correctly returns 0.0 for total == 0 and allowedErrorRate <= 0.
Assessment: PASS
Severity: -
Notes: Covered by inspection; missing tests (see Finding 4). Logic correct.