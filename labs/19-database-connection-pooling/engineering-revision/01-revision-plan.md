# Engineering Revision Plan

Target Lab: labs/19-database-connection-pooling
Previous Verdict: NEEDS_REVISION

## Blocking Issues
1. Data Race in `MockDriver.Open` (`internal/pool/mockdb.go:41`) where `d.activeConns` is read non-atomically.

## Non-Blocking Issues
1. Missing Deadlock Test: Pool-locking deadlock scenario (Finding 11) is not tested.
2. Implementation Overclaim: Engineering doc claims "throughput degradation or max connection enforcement" while implementation only enforces max connections.

## Files To Change
- `internal/pool/mockdb.go`
- `engineering/01-design.md`

## Tests To Add/Modify
- `tests/pool_test.go`: Add `TestPoolLockingDeadlock` to test the pool-locking deadlock edge case.

## Validation Commands
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`