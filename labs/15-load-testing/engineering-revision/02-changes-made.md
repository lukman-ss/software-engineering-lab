# Changes Made

Target Lab: labs/15-load-testing

## Revision 1

Audit Issue: Gap 1 (MEDIUM) — Missing direct test for `MaxDBConnections` concurrency bound enforcement.
Severity: MEDIUM
Files Changed:
- `internal/server/server.go`
- `tests/loadtest_test.go`
Action: Added `ActiveConnections()` to `Server` and created `TestServer_MaxDBConnectionsBound` to poll and assert peak concurrency <= `MaxDBConnections` under heavy concurrent traffic.
Verification: `go test -race ./tests` passed.
Status: RESOLVED

---

## Revision 2

Audit Issue: Gap 2 & 3 (LOW) — Missing unit tests for RPS calculation and single-sample (`len == 1`) percentile edge case.
Severity: LOW
Files Changed:
- `internal/loadtest/metrics_test.go`
Action: Added `TestCalculateMetrics_SingleSample` and `TestCalculateMetrics_RPS`.
Verification: `go test -v ./internal/loadtest` passed.
Status: RESOLVED

---

## Revision 3

Audit Issue: Gap 4 (LOW) — Missing end-to-end assertion for `SuccessCount + ErrorCount == TotalRequests`.
Severity: LOW
Files Changed:
- `tests/loadtest_test.go`
Action: Added `TestLoadTest_SuccessAndErrorInvariant` asserting sum equality on multi-VU test run.
Verification: `go test -v ./tests` passed.
Status: RESOLVED

---

## Revision 4

Audit Issue: Gap 5, 6, 7 (LOW) — Documentation terminology mismatches in `01-design.md` and `02-implementation-notes.md`.
Severity: LOW
Files Changed:
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
Action:
- Updated `01-design.md` to describe bounded duration and per-VU buffer aggregation.
- Updated `02-implementation-notes.md` to clarify configurable `DBQueryDuration` vs demo 20ms parameter.
Verification: Manual review of markdown files against implementation.
Status: RESOLVED
