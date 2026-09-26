# Engineering Revision Plan

Target Lab: `labs/24-slo-sli-error-budget`
Previous Verdict: APPROVED (with non-blocking warnings)

## Blocking Issues

None.

## Non-Blocking Issues

1. **GAP-01**: `WindowTracker.Record` bucket append assumption on timestamp monotonicity could misorder buckets if out-of-order timestamps arrive.
2. **GAP-02**: Missing negative assertion test for multi-window burn rate alert suppression when short window fires but long window does not.
3. **GAP-03**: Endpoint criticality bucketing referenced in design doc was omitted from the interactive demonstration `cmd/demo/main.go`.
4. **Edge Case**: Missing test coverage for zero-traffic evaluation in `SLOEvaluator`.

## Files To Change

- `internal/metrics/tracker.go`: Implement sorted in-place insertion / bucket merging for out-of-order events.
- `cmd/demo/main.go`: Add Phase 4 demonstrating endpoint criticality comparison (stricter Payment SLO 99.9% vs relaxed Reports SLO 95.0%).
- `tests/slo_test.go`: Add negative multi-window alert test, out-of-order timestamp insertion test, and zero-traffic evaluator test.

## Tests To Add/Modify

- `TestAlertEngineBurnRate`: Add negative test case asserting no alerts trigger when only short window exceeds threshold.
- `TestOutOfOrderTimestamps`: Test out-of-order timestamp insertion and bucket eviction.
- `TestEvaluatorZeroTraffic`: Test zero traffic SLI calculation and deployment allowance.

## Validation Commands

```bash
go test -v ./tests
go test -race ./tests
go run ./cmd/demo
```
