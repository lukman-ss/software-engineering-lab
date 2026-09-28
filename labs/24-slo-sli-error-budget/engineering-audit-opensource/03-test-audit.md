## Test Coverage Overview

- `TestMetricsWindowTracker`: verifies window tracking, good/bad counts, eviction.
- `TestSLOEvaluator`: validates SLI calculation, budget remaining logic, deployment gating.
- `TestAlertEngineBurnRate`: checks burn‑rate alert triggering and negative scenario.
- `TestOutOfOrderTimestamps`: ensures correct handling of out‑of‑order events.
- `TestEvaluatorZeroTraffic`: confirms zero‑traffic edge case.
- `TestConcurrencyMetrics`: stress‑tests thread‑safety with 20 goroutines recording 100 events each.

All tests compile and pass (`go test ./...`), race detector reports no data races.

**Missing Aspects**
- No performance benchmarks for large traffic volumes.
- No tests for `BurnRateRule.LongWindow`/`ShortWindow` fields (currently unused).
- No integration test exercising the demo binary end‑to‑end.
