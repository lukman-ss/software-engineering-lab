# Engineering Revision Plan

Target Lab: labs/15-load-testing
Previous Verdict: APPROVED_WITH_WARNINGS (OpenSource Audit: 0 Critical, 0 High, 1 Medium, 5 Low)

## Blocking Issues
None.

## Non-Blocking Issues
1. (MEDIUM) MISSING_TEST: Server connection pool concurrency limit (`MaxDBConnections`) not directly asserted.
2. (LOW) MISSING_TEST: Unit test for RPS computation against known duration and sample count missing.
3. (LOW) MISSING_TEST: Percentile single-sample (`len == 1`) edge-case calculation missing.
4. (LOW) MISSING_TEST: Invariant `SuccessCount + ErrorCount == TotalRequests` not asserted end-to-end.
5. (LOW) DOC_CODE_MISMATCH: `engineering/01-design.md` described LoadTester with "iterations" instead of bounded duration.
6. (LOW) DOC_CODE_MISMATCH: `engineering/01-design.md` described MetricsAggregator as a single mutexed collector rather than per-VU buffers aggregated via pure function.
7. (LOW) DOC_CODE_MISMATCH: `engineering/02-implementation-notes.md` stated mock wait duration was fixed at 20ms, whereas 20ms is demo-specific and configurable in server.

## Files To Change
- `internal/server/server.go`: Expose `ActiveConnections() int` helper for testing current semaphore usage.
- `internal/loadtest/metrics_test.go`: Add `TestCalculateMetrics_SingleSample` and `TestCalculateMetrics_RPS`.
- `tests/loadtest_test.go`: Add `TestServer_MaxDBConnectionsBound` and `TestLoadTest_SuccessAndErrorInvariant`.
- `engineering/01-design.md`: Clarify runner execution duration and per-VU aggregation model.
- `engineering/02-implementation-notes.md`: Clarify configurable `DBQueryDuration` vs demo default.

## Tests To Add/Modify
- `TestCalculateMetrics_SingleSample` (Unit)
- `TestCalculateMetrics_RPS` (Unit)
- `TestServer_MaxDBConnectionsBound` (Integration/Concurrency)
- `TestLoadTest_SuccessAndErrorInvariant` (Integration)

## Validation Commands
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
