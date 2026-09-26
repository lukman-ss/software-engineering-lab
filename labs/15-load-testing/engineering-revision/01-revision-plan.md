# Engineering Revision Plan

Target Lab: labs/15-load-testing
Previous Verdict: NEEDS_REVISION

## Blocking Issues

1. **RESEARCH_MISMATCH**: In `internal/server/server.go`, uniform query latency (20ms) under a closed-loop virtual user harness caused queuing delay to distribute uniformly, leading to Average latency (~200ms) matching P95 latency (~211ms). This failed to prove the core research finding that average response times conceal tail latency spikes.

## Non-Blocking Issues

1. **TEST_ASSERTION_DEFICIENCY**: `TestLoadTest_SmokeVsStress` did not assert that `stressRes.P95Latency` diverged from `stressRes.AvgLatency`.
2. **MISSING_INVARIANT_TEST**: No test verified that percentile calculation maintained mathematical ordering invariants (Min <= P50 <= P90 <= P95 <= P99 <= Max) on arbitrary or non-monotonic samples.
3. **DOC_CODE_MISMATCH**: `engineering/03-execution-result.md` contained runtime numbers with uniform latencies that contradicted the claimed masking effect.

## Files To Change

- `internal/server/server.go`: Introduce intermittent tail contention (10% chance of 25x duration increase when saturated) to simulate long-tail bottleneck behavior.
- `internal/loadtest/metrics_test.go`: Add `TestCalculateMetrics_Invariants` to verify percentile monotonicity under irregular distributions.
- `tests/loadtest_test.go`: Update `TestLoadTest_SmokeVsStress` to verify that `stressRes.P95Latency > stressRes.AvgLatency`.
- `engineering/03-execution-result.md`: Update benchmark execution logs to reflect actual tail divergence.

## Tests To Add/Modify

- `TestCalculateMetrics_Invariants`: New test validating Min <= P50 <= P90 <= P95 <= P99 <= Max on non-monotonic inputs.
- `TestLoadTest_SmokeVsStress`: Strengthened with assertion `stressRes.P95Latency > stressRes.AvgLatency`.

## Validation Commands

```bash
cd labs/15-load-testing
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
