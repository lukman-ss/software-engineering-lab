# Code Audit

## Finding 1

Location: internal/metrics/tracker.go:46-104 (Record), 117-127 (Summary)
Claimed Behavior: Thread-safe sliding-window bucket aggregation with eviction.
Observed Implementation: sync.RWMutex used (Lock on Record, Lock on Summary). Buckets kept sorted by StartTime; out-of-order events inserted in order. Stale buckets evicted by cutoff = now - windowSize.
Assessment: PASS
Severity: LOW
Notes: Race detector clean. Summary acquires write lock; could use RLock for read-only path but not a correctness issue.

## Finding 2

Location: internal/metrics/tracker.go:67-91 (out-of-order insertion)
Claimed Behavior: Handles out-of-order timestamps by inserting into correct bucket position.
Observed Implementation: Loop scans existing buckets; inserts at sorted position or updates existing bucket. Verified by TestOutOfOrderTimestamps.
Assessment: PASS
Severity: LOW
Notes: Insertion into middle of slice is O(n); acceptable for demo scale.

## Finding 3

Location: internal/slo/evaluator.go:41-71 (Evaluate)
Claimed Behavior: SLI = good/total, error budget = (1 - targetSLO) * total, consumed = bad count, remaining = budget - consumed.
Observed Implementation: Count-based budget. SLI rounded to 4 decimals, budget to 2 decimals. CanDeploy false when remaining <= 0.
Assessment: WARNING
Severity: MEDIUM
Notes: Design doc says "Error Budget = 1 - SLO" (ratio), but implementation uses count-based budget (ratio * total). Valid interpretation but differs from standard SRE time-based error budget. LatencyThreshold in Config is never used by evaluator — dead field (tracked via isGood callback instead).

## Finding 4

Location: internal/alerting/engine.go:71-75 (Check)
Claimed Behavior: Multi-window burn-rate alerting.
Observed Implementation: Trigger requires BOTH shortBurn >= factor AND longBurn >= factor.
Assessment: WARNING
Severity: MEDIUM
Notes: Standard SRE multi-window burn-rate uses OR (either window exceeding threshold fires). AND logic is more conservative and can miss fast-burn incidents where long window is still clean. Single BurnRateFactor per rule — no separate short/long factors.

## Finding 5

Location: cmd/demo/main.go (demo output)
Claimed Behavior: Demo shows baseline, incident, alert, endpoint comparison.
Observed Implementation: Demo output matches engineering/03-execution-result.md exactly (verified by execution).
Assessment: PASS
Severity: LOW
Notes: No fake demo. Output is real.