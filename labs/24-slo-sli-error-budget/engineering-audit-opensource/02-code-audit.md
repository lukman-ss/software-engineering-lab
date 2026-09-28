# Code Audit

Target Lab: `labs/24-slo-sli-error-budget`
Note: read-only audit; no code modified.

## Finding 1

Location: `internal/metrics/tracker.go:46-64` (`Record`, `Summary`, `evictStaleLocked`)
Claimed Behavior: sliding-window bucketed aggregation of good vs total events with stale eviction.
Observed Implementation: `sync.RWMutex`-guarded bucket slice; `Record` and `Summary` both take the
write lock and evict buckets with `StartTime.Before(now - windowSize)` before aggregating.
Assessment: PASS
Severity: LOW
Notes: `Summary` uses `mu.Lock` instead of `RLock` (conservative, still correct). Eviction boundary
is strict-`Before`, so a bucket exactly on the cutoff is retained — deterministic, documented in test.

## Finding 2

Location: `internal/metrics/tracker.go:66-92` (out-of-order insert path)
Claimed Behavior: tracker tolerates out-of-order / late events.
Observed Implementation: linear scan finds a matching bucket or sorted-insert position; covered by
`TestOutOfOrderTimestamps` including partial-eviction of the earlier bucket.
Assessment: PASS
Severity: LOW
Notes: O(n) insert is fine at this scale (bounded by window/bucket count, e.g. ~180 buckets in demo).

## Finding 3

Location: `internal/metrics/tracker.go:22-28,30-44` + `internal/slo/evaluator.go:10-14`
Claimed Behavior: robust construction.
Observed Implementation: `NewWindowTracker` guards `bucketSize <= 0` but does not guard a nil
`isGoodEvent` (nil would panic in `Record`). `NewEvaluator` does not validate `TargetUptime` range.
`Config.LatencyThreshold` is stored but never read — latency classification lives in the caller's
`isGood` closure (see `cmd/demo/main.go:20-22`).
Assessment: WARNING
Severity: LOW
Notes: No crash path is reachable from tests or demo; robustness-only gap, not a correctness failure.

## Finding 4

Location: `internal/slo/evaluator.go:41-70`
Claimed Behavior: SLI = good/total; budget = allowed-rate*total - bad; freeze when exhausted.
Observed Implementation: matches claim exactly. Zero-traffic yields SLI 1.0 + `CanDeploy=true`
(`total > 0` guard). Over-consumption goes negative (demo shows -8.90), correctly signalling debt.
Rounding (SLI 4dp, budget 2dp) is cosmetic only.
Assessment: PASS
Severity: LOW
Notes: `TestSLOEvaluator` proves the exact boundary (99 good + 1 bad deploys; 2nd bad freezes).

## Finding 5

Location: `internal/alerting/engine.go:51-89`
Claimed Behavior: multi-window multi-burn-rate alerting.
Observed Implementation: burn-rate math (`actual/allowed`) is correct with zero guards
(`total == 0`, `allowedErrorRate <= 0`). `Check` fires only when BOTH short and long burn exceed the
rule threshold — correct anti-flap semantics, proven by the transient-spike negative test.
BUT: `BurnRateRule.LongWindow`, `ShortWindow`, and `BudgetConsumedPct` (lines 21-24) are never read.
Window geometry comes solely from the two trackers passed to `NewAlertEngine`; there is no
per-rule window or budget-consumed guard as in the Google MWMBR formulation.
Assessment: WARNING
Severity: MEDIUM
Notes: Core threshold logic is proven; per-rule configurability implied by the struct is unimplemented.

## Finding 6

Location: `cmd/demo/main.go:24-30,106-115`
Claimed Behavior: demo exercises distinct short vs long burn windows.
Observed Implementation: all ~1100 events land within ~52s, inside BOTH the 5-min and 60-min windows,
so `ShortBurn == LongBurn == 9.09x`. The demo never separates the windows; only the unit test's
transient-spike case (separate trackers) proves true multi-window discrimination. Consequence visible
in output: only the 6x TICKET rule fires, the 14.4x PAGE rule does not (9.09 < 14.4).
Assessment: WARNING
Severity: MEDIUM
Notes: Output is honest (no faked alert); the demo's multi-window aspect is degenerate, not fabricated.

## Finding 7

Location: `cmd/demo/main.go` (whole) vs `engineering/03-execution-result.md:45-73`
Claimed Behavior: recorded demo output.
Observed Implementation: re-ran `go run ./cmd/demo` — output matches the recorded transcript
line-for-line numerically (1000/1000/0 → SLI 100%, budget 1.00; 1100/1090/10 → SLI 99.09%,
budget -8.90; TICKET 9.09x/9.09x thr 6.00x; Reports 90.0%, budget -5.00).
Assessment: PASS
Severity: LOW
Notes: No FAKE_DEMO. Phase 4 narrative ("wider 5% tolerance") is confusing since both services end
frozen (`CanDeploy=false`), but the numbers are real.

## Finding 8

Location: concurrency paths (`tracker.go:46-48,117-119`; `tests/slo_test.go:200-236`)
Claimed Behavior: thread-safe concurrent recording.
Observed Implementation: every mutation/aggregation path holds the mutex; `go test -race ./...`
passes; 20 goroutines x 100 records lose no events (total == 2000, good+bad == total).
Assessment: PASS
Severity: LOW
Notes: Test asserts the sum invariant, not exact 1800/200 split — sufficient for safety, see test audit.
