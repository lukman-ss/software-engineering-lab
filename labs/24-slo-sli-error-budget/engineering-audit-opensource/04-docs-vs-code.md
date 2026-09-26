# Docs vs Code — labs/24-slo-sli-error-budget

Method: compared README.md, engineering/*.md (design, implementation-notes, execution-result) against actual source, tests, and a fresh execution of `go run ./cmd/demo`.

## 1. README vs Code — structure & commands
| README Claim | Code Reality | Match |
|---|---|---|
| internal/metrics: sliding-window time-bucketed tracker | tracker.go: WindowTracker + Bucket | PASS |
| internal/slo: SLI/SLI ratio + error budget + release freeze | evaluator.go Evaluate returns Status with CanDeploy | PASS |
| internal/alerting: multi-window multi-burn-rate alerting | engine.go Check | WARNING (see #4) |
| cmd/demo: baseline, incident budget burn, alert trigger | main.go 4 phases | PASS (plus extra Phase 4 not in README prose) |
| tests/: unit + concurrency + thread-safety | tests/slo_test.go 6 tests, -race clean | PASS |
| `go test ./...` / `go run ./cmd/demo` | identical to README | PASS |

## 2. engineering/01-design.md vs code
- DOC: "internal/metrics: Sliding-window event recorder (histogram latency buckets & success counts)." -> CODE: NO latency histograms. metrics tracks per-bucket counts (Total/Good/Bad); latency only appears as an Event field compared by isGood closure. `histogram latency buckets` is design aspiration not implemented.
  - Classification: RESEARCH/DESIGN_IMPLEMENTATION_MISMATCH (LOW)
- DOC: "BurnRateAlertEngine: Monitors short and long windows for fast/slow burn threshold breaches." -> CODE: single BurnRateFactor per rule, shared engine-level short/long trackers; LongWindow/ShortWindow/BudgetConsumedPct unused.
  - Classification: DOC_CODE_MISMATCH (LOW — see code Finding 1)
- DOC success criteria: "100% test coverage on core math and sliding window calculations." -> MEASURED: 94.1% of ./internal statements.
  - Classification: TEST_CLAIM_MISMATCH (LOW)
- DOC success criteria: "Demonstration output shows ... recovery" -> CODE: demo has NO recovery phase after the simulated incident (budget stays exhausted, no restore).
  - Classification: DOC_CODE_MISMATCH (LOW)
- DOC: "in-memory ring/time-bucketed window tracking to simulate 30-day/rolling windows in compressed real-time" -> CODE: window30d=30*time.Minute (30 min, compressed). Intentional per "Implementation Decisions"; variable name misleading.
  - Classification: WARNING (naming) / not a mismatch on intent

## 3. engineering/03-execution-result.md vs actual run — FAKE/INCOMPLETE DEMO RECORD
The recorded execution-result.md asserts this demo output:
```
[PHASE 3] Checking Multi-Window Burn Rate Alerts...
>>> ALERT TRIGGERED: [TICKET] ... (Threshold: 6.00x)
   <then DEMO COMPLETE, no Phase 4>
```
The ACTUAL `go run ./cmd/demo` output (verified this run) ALSO prints:
```
[PHASE 4] Endpoint Criticality Comparison (Payment 99.9% vs Reports 95.0%)...
Reports Target SLO: 95.0% | Current SLI: 90.0% | Budget Remaining: -5.00
Payment CanDeploy: false | Reports CanDeploy: false (Reports has wider 5% error tolerance)
```
The demo has FOUR phases but the recorded doc captured only three and dropped Phase 4 entirely. The recorded output is also missing the Phase-2 numbers that ARE present, so it is not a full capture.
  - Classification: DOC_CODE_MISMATCH / incomplete demo record (MEDIUM). Not fabricated values (Phase 3 matches byte-for-byte and math is correct), but the record is stale relative to the final source (Phase 4 was added after the doc was written, OR the capture was truncated).
- execution-result.md test list also INCOMPLETE: lists only TestMetricsWindowTracker, TestSLOEvaluator, TestAlertEngineBurnRate, TestConcurrencyMetrics. Actual `go test -v` shows 6 tests; TestOutOfOrderTimestamps and TestEvaluatorZeroTraffic are omitted from the record.
  - Classification: DOC_CODE_MISMATCH (LOW)

## 4. README vs demo output (no claim numbers to contradict)
README makes no numeric claims about demo output; it only promises "error budget depletion during an incident, and burn rate alert triggering." Demo delivers: Phase 2 budget -8.90 CanDeploy=false; Phase 3 TICKET alert. PASS.

## 5. Mismatch taxonomy (summary)
```text
DOC_CODE_MISMATCH   design.md: "histogram latency buckets" (not in code)           -> LOW
DOC_CODE_MISMATCH   design.md: "100% test coverage" (measured 94.1%)                -> LOW
DOC_CODE_MISMATCH   design.md: demo shows "recovery" (no recovery phase)            -> LOW
DOC_CODE_MISMATCH   design.md/engine: BurnRateEngine per-rule windows unused       -> LOW
TEST_CLAIM_MISMATCH execution-result.md test list omits 2 of 6 tests               -> LOW
DOC_CODE_MISMATCH   execution-result.md demo record omits Phase 4 output           -> MEDIUM
```
No RESEARCH_IMPLEMENTATION_MISMATCH / TEST_CLAIM_MISMATCH on core math (SLI, error budget, burn rate all verified to reproduce). No FAKE_DEMO (numbers are real and reproducible). No FAKE_BENCHMARK (no benchmarks claimed). No UNVERIFIED_RESULT (tests and demo both re-run and match).
