# Gap Analysis

Target Lab: `labs/24-slo-sli-error-budget`

## Identified Gaps

### GAP-01: Out-of-Order Timestamp Handling
- Type: `MISSING_EDGE_CASE`
- Severity: LOW
- Location: `internal/metrics/tracker.go:56-75`
- Description: `WindowTracker.Record` assumes incoming events are monotonically increasing in timestamp. If out-of-order events arrive, buckets are appended in out-of-order sequence, which may affect prefix-slice eviction in `evictStaleLocked`.
- Recommended Action: In a future iteration, insert buckets in sorted order or discard events older than current window cutoff.

### GAP-02: Missing Multi-Window Negative Alert Test
- Type: `MISSING_TEST`
- Severity: LOW
- Location: `tests/slo_test.go:93`
- Description: Test suite verifies that an alert triggers when both short and long windows exceed the burn rate threshold, but does not include a test asserting that an alert is suppressed when only the short window or only the long window exceeds the threshold.
- Recommended Action: Add unit test verifying suppression of alert when `shortBurn >= threshold && longBurn < threshold`.

### GAP-03: Endpoint Criticality Demonstration Omitted
- Type: `DOC_CODE_MISMATCH`
- Severity: LOW
- Location: `engineering/01-design.md:10`
- Description: Design doc lists endpoint criticality bucketing (Payment 99.9% vs Reports 95%) as a concept to prove, but `cmd/demo/main.go` only instantiates a single evaluator for Payment Service.
- Recommended Action: Update design doc or add the secondary evaluator in demo to illustrate contrasting policies between endpoints.
