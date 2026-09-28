# Test Audit

## Test Suite Execution Results

### 1. `go test -v ./...`
```text
=== RUN   TestPartitionPruning
--- PASS: TestPartitionPruning (0.00s)
=== RUN   TestRoutingAndConsistentHashRelocation
    sharding_test.go:107: Hash Modulo moved 3719 / 5000 keys (74.38%)
    sharding_test.go:108: Consistent Hash moved 1162 / 5000 keys (23.24%)
--- PASS: TestRoutingAndConsistentHashRelocation (0.01s)
=== RUN   TestClusterScatterGatherAndGSI
--- PASS: TestClusterScatterGatherAndGSI (0.00s)
=== RUN   TestIDGenerators
--- PASS: TestIDGenerators (0.00s)
=== RUN   TestConcurrentClusterAccess
--- PASS: TestConcurrentClusterAccess (0.01s)
PASS
ok  	labs/32-database-sharding-and-partitioning/tests	0.279s
```

### 2. `go test -race ./...`
```text
PASS
ok  	labs/32-database-sharding-and-partitioning/tests	1.121s
```

### 3. `go run ./cmd/demo`
Successfully executed all 5 demo sections:
1. Single-Node Logical Table Partitioning & Range Pruning (1 partition scanned out of 4, 3 pruned).
2. Sharding Key Selection (Monotonic key showing 100% write hotspot on 1 shard vs high-cardinality key evenly distributed).
3. Resharding / Scale-out Comparison (Hash modulo remapped 79.84% keys vs consistent hashing 12.00%).
4. Querying Non-Sharded Attributes (Scatter-Gather broadcasting to 4/4 shards vs GSI point lookup querying 1 shard).
5. Distributed Unique ID Generation (UUIDv7 time-ordered validation and Sequence Block Allocator).

## Coverage & Gap Review

- Happy path: Tested across all components.
- Failure path: Unmatched partition returns `ErrNoMatchingPartition`, missing shard key / email tested, context cancellation in scatter-gather tested.
- Edge cases: Wrap-around on consistent hash ring verified, sequence block boundary crossing verified.
- Concurrency: `TestConcurrentClusterAccess` exercises parallel inserts, point lookups, and GSI lookups across 200 concurrent goroutines with `-race` enabled.
