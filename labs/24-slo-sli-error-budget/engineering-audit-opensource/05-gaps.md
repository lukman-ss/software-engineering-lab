# Gap Analysis

Target Lab: `labs/24-slo-sli-error-budget`

## Gaps

### GAP-1 — IMPLEMENTATION_OVERCLAIM (MEDIUM)

Per-rule multi-window configurability implied by `BurnRateRule{LongWindow, ShortWindow,
BudgetConsumedPct}` is unimplemented; engine uses two fixed trackers and threshold-only gating.
Core burn math and dual-window gating are proven; configurability is not.
Fix: either wire per-rule windows + budget-consumed guard into `Check`, or remove the dead fields
and document the two-tracker design as the intended simplification.

### GAP-2 — MISSING_EDGE_CASE (MEDIUM)

Demo short/long windows contain identical events (all traffic within ~52s ⊂ both 5-min and 60-min
windows), so the demo never discriminates windows; only the unit-test transient case does.
Fix: stagger demo incident traffic so the short window burns while the long window stays clean
(or vice versa), showing each rule respond independently.

### GAP-3 — DOC_CODE_MISMATCH (MEDIUM)

Design doc claims "histogram latency buckets" and "100% test coverage"; neither holds
(no latency distribution in `Bucket`; uncovered branches, no coverage artifact).
Fix: reword to "good/bad event counts with caller-supplied latency predicate" and replace "100%"
with the actual measured coverage number.

### GAP-4 — DOC_CODE_MISMATCH (LOW)

Design doc promises a demo "recovery" phase; demo ends with both services frozen.
Fix: add a recovery phase (traffic heals, budget restored, alerts clear) or drop the claim.

### GAP-5 — MISSING_TEST (LOW)

No 100%-error test despite being listed in the design Test Strategy; no direct test of
`CalculateBurnRate` zero guards.
Fix: add `TestAllBadEvents` (SLI 0, budget deeply negative, Page fires) and a zero-total burn-rate case.

### GAP-6 — MISSING_EDGE_CASE (LOW)

Concurrency test asserts totals, not the deterministic 1800/200 good/bad split; nil `isGoodEvent`
and out-of-range `TargetUptime` are unguarded.
Fix: assert exact split; add constructor guards or document preconditions.

## Explicitly NOT gaps (checked, clear)

- FAKE_DEMO / FAKE_BENCHMARK / UNVERIFIED_RESULT: demo re-ran byte-identical numerically; no
  benchmarks claimed; execution record accurate.
- RACE_CONDITION: `-race` clean over same-bucket contention test.
- BROKEN_IMPLEMENTATION / UNHANDLED_ERROR: no failing behavior found; error paths (zero traffic,
  zero total, degenerate SLO) all guarded.
- RESEARCH_MISMATCH: out of scope per pipeline override.
