# Engineering Audit Verdict

Target Lab: labs/33-read-replicas-and-replication-lag
Audit Date: 2026-09-28

## Summary

Code Files Reviewed: 3 (`internal/cluster/cluster.go`, `internal/router/router.go`, `cmd/demo/main.go`)
Tests Reviewed: 1 (`tests/replication_test.go`, 6 test functions)
Commands Executed: `go build ./...`, `go test -v -count=1 ./...`, `go test -race -count=1 ./...`, `go run ./cmd/demo`, `go vet ./...`
Failures: 0
Warnings: 10 gap items (1 BROKEN_IMPLEMENTATION MEDIUM, 1 UNHANDLED_ERROR LOW, 1 MISSING_EDGE_CASE LOW, 7 MISSING_TEST LOW/MEDIUM)

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS
Demo: PASS
Research Alignment: PASS
Documentation Accuracy: WARNING

## Blocking Issues
1. None. No HIGH or CRITICAL issues found.

## Non-Blocking Issues
1. BROKEN_IMPLEMENTATION (MEDIUM): `Cluster.Write` async branch silently drops WAL entries when replica channel buffer full (`cluster.go:220-225`). Permanent replica divergence possible under burst load, no error surfaced.
2. UNHANDLED_ERROR (LOW): `ReadWithToken` discards `WaitForLSN` error (`router.go:107`); caller cannot distinguish wait-success from timeout-fallback.
3. MISSING_EDGE_CASE (LOW): `ReadWithToken` with pre-cancelled context; leaked waiter goroutine on timeout (`Node.WaitForLSN`, `cluster.go:79-96`).
4. MISSING_TEST (MEDIUM): `ReadWithToken` timeout→fallback branch untested; concurrent test asserts no races but not correctness.
5. DOC_CODE_MISMATCH (LOW): design doc test name `TestReadYourOwnWrites_LSNToken` vs actual `TestReadWithToken_LSN`.

## Required Revisions
1. Fix or explicitly bound async WAL send: block on channel send or return error on full buffer instead of silent drop.
2. Document or fix `WaitForLSN` goroutine leak on context timeout/cancel.
3. Add tests: token timeout→fallback, token fast path, zero-replica fallback, write-after-close error, dynamic mode switch.
4. Strengthen concurrency test with value assertions or add LSN-ordering check.
5. Fix design doc test name to match actual test.

## Final Status

APPROVED_WITH_WARNINGS
