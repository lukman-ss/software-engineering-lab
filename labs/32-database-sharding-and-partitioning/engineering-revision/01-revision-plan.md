# Engineering Revision Plan

Target Lab: `labs/32-database-sharding-and-partitioning`
Previous Verdict: APPROVED

## Blocking Issues
None.

## Non-Blocking Issues
1. **UNHANDLED_ERROR (LOW)**: `ScatterGatherBroadcast` in `internal/sharding/sharding.go` lacked `context.Context` timeout / cancellation support.
2. **MISSING_TEST (LOW)**: `TestRoutingAndConsistentHashRelocation` in `tests/sharding_test.go` used key generator format yielding 0% relocation in unit test due to small sample size / key prefix alignment.

## Files To Change
- `internal/sharding/sharding.go`
- `tests/sharding_test.go`

## Tests To Add/Modify
- Update `TestRoutingAndConsistentHashRelocation` in `tests/sharding_test.go` to use larger sample size (5,000 keys) and prefix matching dataset to prove positive relocation (> 5% and <= 40%).
- Add context cancellation test for `ScatterGatherBroadcastWithContext` in `tests/sharding_test.go`.

## Validation Commands
```bash
cd labs/32-database-sharding-and-partitioning
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
