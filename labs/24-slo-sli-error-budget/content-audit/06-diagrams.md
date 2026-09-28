# Diagrams Verification

## Lab
`labs/24-slo-sli-error-budget`

## Audit Scope
Verification of `content/04-diagrams.md` against engineering implementation and tests.

## Verification Results

### D1 — High-Level Architecture ✅ ACCURATE
- **Components**: cmd/demo → SLO Evaluator, Short Window Tracker, Long Window Tracker → AlertEngine
- **Data Flow**: Demo records to all trackers; Evaluator uses Summary; AlertEngine reads both trackers
- **Implementation Match**: 
  - `cmd/demo/main.go:28-30`: Creates three WindowTracker instances
  - `internal/slo/evaluator.go:41`: Uses single tracker
  - `internal/alerting/engine.go:26-30`: Holds shortTracker and longTracker

### D2 — WindowTracker Bucket Lifecycle ✅ ACCURATE
- **Flow**: evictStaleLocked → Truncate timestamp → isGoodEvent → bucket lookup/insert
- **Three insertion paths correctly shown**:
  1. Last bucket matches → increment (fast path)
  2. Earlier timestamp → search sorted slice → increment or insert
  3. New timestamp → append
- **Out-of-order handling**: Search for matching bucket start; if gap, insert in order
- **Eviction**: Buckets older than `now - windowSize` removed
- **Implementation Match**: `internal/metrics/tracker.go:46-104` (Record) and `106-115` (evictStaleLocked)

### D3 — Error Budget & Burn Rate Relationship ✅ ACCURATE
- **SLI Calculation**: `good/total = 1090/1100 = 99.09%` ✅
- **Error Budget**: `totalBudget = (1-SLO) × total = 0.1% × 1100 = 1.1` ✅
- **Budget Remaining**: `1.1 - 10 = -8.9` (exhausted) ✅
- **CanDeploy**: `false` when budget ≤ 0 ✅
- **Burn Rate**: `actualErrorRate / allowedErrorRate = (10/1100) / 0.001 = 9.09x` ✅
- **AlertEngine Multi-Window**: Requires BOTH shortBurn AND longBurn ≥ threshold ✅
- **Implementation Match**: `internal/slo/evaluator.go:44-57`, `internal/alerting/engine.go:51-88`

### D4 — Multi-Window Alert: True Positive vs False Positive ❌ CONTAINS ERROR

| Column | Diagram Shows | Correct? | Issue |
|--------|---------------|----------|-------|
| TRUE POSITIVE (alert fires) | Short: 98 OK + 2 ERR (2% → 20x). Long: 9800 OK + 200 ERR (2% → 20x). Both ≥ 6.0 → TICKET ✓ | ✅ | Correct |
| FALSE POSITIVE (no alert) | Short: 90 OK + 10 ERR (10% → 100x). Long: 90 OK + 10 ERR (10% → 100x). "Short ≥ 6.0 BUT Long < 6.0 → NO ALERT" | ❌ | **INTERNAL CONTRADICTION** — If long window is 100x, it IS ≥ 6.0. This would be a TRUE POSITIVE, not false positive. |
| TRUE NEGATIVE (no alert) | Short: 98 OK + 2 ERR (2% → 20x). Long: 9999 OK + 1 ERR (0.01% → 0.1x). "Short ≥ 6.0 BUT Long < 6.0 → NO ALERT" | ✅ | Correct (this IS the false positive prevention case) |

**Root Cause**: The FALSE POSITIVE column appears to be a copy-paste of the TRUE POSITIVE column data with incorrect labeling. The scenario showing "Short ≥ 6.0 BUT Long < 6.0 → NO ALERT" is actually the TRUE NEGATIVE column's scenario (transient spike in short window only).

**Correct FALSE POSITIVE Prevention Scenario** (from `tests/slo_test.go:130-151`):
- **Short window**: 90 good + 10 bad → 10% error rate → 100x burn rate
- **Long window**: 9999 good + 1 bad → 0.01% error rate → 0.1x burn rate
- **Result**: Short ≥ 14.4x (YES), Long ≥ 14.4x (NO) → **NO ALERT** ✓

### D5 — Test Coverage Map ✅ ACCURATE
| Test | Target Feature | Verified By |
|------|----------------|-------------|
| TestMetricsWindowTracker | Bucket agg + eviction | `tests/slo_test.go:13-59` |
| TestSLOEvaluator | SLI & Error Budget policy | `tests/slo_test.go:61-91` |
| TestAlertEngineBurnRate | Multi-window burn rate (positive + negative) | `tests/slo_test.go:93-152` |
| TestOutOfOrderTimestamps | Bucket ordering + eviction | `tests/slo_test.go:154-177` |
| TestEvaluatorZeroTraffic | SLI=1.0, CanDeploy=true on zero traffic | `tests/slo_test.go:179-198` |
| TestConcurrencyMetrics | 20 goroutines × 100 requests, race-free | `tests/slo_test.go:200-236` |

### D6 — Demo Flow State Transition ✅ ACCURATE
| Phase | State | Demo Output |
|-------|-------|-------------|
| 1 | Baseline: 1000 good, 0 bad | SLI=100%, Budget=+1.00, CanDeploy=true |
| 2 | Incident: 1100 total (1000+90 good, 10 bad) | SLI=99.09%, Budget=-8.90, CanDeploy=false |
| 3 | Alert Check | TICKET triggered (9.09x ≥ 6.0x) |
| 4 | Criticality Comparison | Payment 99.9% (budget -8.90), Reports 95.0% (budget -5.00), both CanDeploy=false |

Matches `engineering/03-execution-result.md:48-73` exactly.

## Diagram Accuracy Summary

| Diagram | Status | Notes |
|---------|--------|-------|
| D1 | ✅ PASS | Architecture correctly mapped |
| D2 | ✅ PASS | Bucket lifecycle matches implementation |
| D3 | ✅ PASS | Formulas match code and demo output |
| D4 | ❌ FAIL | FALSE POSITIVE column has internal contradiction; should show transient spike scenario |
| D5 | ✅ PASS | Test coverage correctly mapped |
| D6 | ✅ PASS | Demo phases match execution results |

## Required Correction for D4

**Current FALSE POSITIVE column** (incorrect):
```
Short Window: 90 OK + 10 ERR → 10% error → 100x burn
Long Window: 90 OK + 10 ERR → 10% error → 100x burn
Short ≥ 6.0 BUT Long < 6.0 → NO ALERT ✓
```

**Should be** (transient spike scenario):
```
Short Window: 90 OK + 10 ERR → 10% error → 100x burn
Long Window: 9999 OK + 1 ERR → 0.01% error → 0.1x burn
Short ≥ 14.4x AND Long ≥ 14.4x? NO → NO ALERT ✓ (transient spike filtered)
```

Or relabel existing TRUE NEGATIVE column as "FALSE POSITIVE PREVENTION" and ensure no duplicate data columns exist.