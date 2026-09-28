# Code Snippets Verification

## Lab
`labs/24-slo-sli-error-budget`

## Audit Scope
Verification of `content/03-code-snippets.md` against actual source code.

## Verification Results

### Verbatim Extract Check
All 13 snippets (1-13) in `content/03-code-snippets.md` are **verbatim extracts** from the approved engineering implementation. The file states: "Semua snippet diambil verbatim dari implementasi yang telah di-audit dan disetujui. Tidak ada kode sintetis." This claim is **VERIFIED** — every snippet matches the source code exactly.

| Snippet | Source File | Lines | Verbatim? |
|---------|------------|-------|----------|
| 1 | internal/metrics/tracker.go | 8-20 | ✅ YES |
| 2 | internal/metrics/tracker.go | 22-44 | ✅ YES |
| 3 | internal/metrics/tracker.go | 46-104 | ✅ YES |
| 4 | internal/metrics/tracker.go | 106-128 | ✅ YES |
| 5 | internal/slo/evaluator.go | 10-32 | ✅ YES |
| 6 | internal/slo/evaluator.go | 41-71 | ✅ YES |
| 7 | internal/alerting/engine.go | 9-49 | ✅ YES |
| 8 | internal/alerting/engine.go | 51-61 | ✅ YES |
| 9 | internal/alerting/engine.go | 63-89 | ✅ YES |
| 10 | cmd/demo/main.go | (config section) | ✅ YES |
| 11 | cmd/demo/main.go | 74-104 | ✅ YES |
| 12 | tests/slo_test.go | 130-151 | ✅ YES |
| 13 | tests/slo_test.go | 200-236 | ✅ YES |

## Critical Observations

### ✅ Correctly Documented Limitations
1. **Snippet 5 / LatencyThreshold**: Correctly notes `Config.LatencyThreshold` is not read by Evaluator; latency judged by caller's `isGood` closure (line 198). Matches GAP-2 from engineering-audit-opensource/05-gaps.md.

2. **Snippet 7 / BurnRateRule fields**: Correctly notes `LongWindow`, `ShortWindow`, `BurnConsumedPct` fields in `BurnRateRule` are ignored; all rules share construction-time trackers (line 285). Matches GAP-1 from engineering-audit-opensource/05-gaps.md.

3. **Snippet 6 / Evaluate logic**: Correctly shows SLI=good/total with zero-traffic fallback to 1.0, budget calculation, and CanDeploy logic (lines 241-242 in content).

### ❌ Inaccuracy: Snippet 11 Burn Rate Claim
| Aspect | Content Claims | Actual Implementation |
|--------|----------------|----------------------|
| Incident batch error rate | 10% (10 errors / 100 incident requests) | Correct ✓ |
| Cumulative error rate | Not explicitly stated | 10/1100 = 0.91% |
| Burn rate calculation | `// 10% / 0.1% = 100x` | `(0.91% / 0.1%) = 9.09x` |
| Demo output | "burn rate = 100x" | Output shows "ShortBurn: 9.09x | LongBurn: 9.09x" |

**Root Cause**: Snippet 11 records the same events to `sloTracker`, `shortTracker`, AND `longTracker` in Phase 1 (baseline 1000 good) and Phase 2 (90 good + 10 bad). The cumulative state across all trackers is 1100 total with 10 bad. The content's explanation computes burn rate using only the Phase 2 batch (10/100 = 10%) rather than the cumulative window state.

**Code Evidence**:
- `cmd/demo/main.go:63-65`: Phase 1 events RECORDED to all three trackers
- `cmd/demo/main.go:84-98`: Phase 2 events RECORDED to all three trackers
- `cmd/demo/main.go:106-113`: AlertEngine.Check computes from tracker.Summary which includes ALL prior events
- `engineering/03-execution-result.md:62`: Output shows `ShortBurn: 9.09x | LongBurn: 9.09x`

### ✅ Correctly Verified: Snippet 12 Transient Spike Test
Content states: "Short window: 100x burn rate. Long window: 10,000 requests, 1 error = 0.01% error rate (0.1x burn rate < 14.4x threshold). Alert TIDAK terpicu."

**Verified**:
- `tests/slo_test.go:137-139`: 90 good to short tracker
- `tests/slo_test.go:140-142`: 10 bad to short tracker (10% error → 100x burn)
- `tests/slo_test.go:143-145`: 9999 good to long tracker
- `tests/slo_test.go:146`: 1 bad to long tracker (0.01% error → 0.1x burn)
- `tests/slo_test.go:148-151`: Asserts 0 alerts triggered

## Explanation Text Accuracy

| Snippet | Explanation Accuracy | Notes |
|---------|---------------------|-------|
| 1 | ✅ | Event/Bucket structure correctly described |
| 2 | ✅ | WindowTracker thread-safety and isGoodEvent correctly described |
| 3 | ✅ | Three insertion paths correctly described |
| 4 | ✅ | Eviction and Summary logic correctly described |
| 5 | ✅ | LatencyThreshold note correct; Status fields accurate |
| 6 | ✅ | SLI computation, budget, and CanDeploy logic correct |
| 8 | ✅ | Burn rate formula and edge cases (total==0, SLO==1.0) correct |
| 9 | ✅ | Multi-window AND-gating logic correct |
| 10 | ✅ | Time compression and 3-tracker setup correct |
| 11 | ❌ | Burn rate 100x claim incorrect; actual is 9.09x |
| 12 | ✅ | Transient spike filtering correctly described |
| 13 | ✅ | Concurrency scenario (20×100 = 2000 total) correct |

## Overall Assessment

Code Snippets Verification: **PASS with 1 HIGH severity correction**

The content is verbatim from source code (as claimed) with accurate explanations of all code behaviors and limitations. The ONLY error is the burn rate calculation in Snippet 11's explanatory text, where the isolated incident batch error rate (10%) is used instead of the cumulative window error rate (0.91%).