# Test Audit

## Test Suite Execution Results

Command:
```bash
go test -v -count=1 ./...
```
Output:
```text
=== RUN   TestPartitionPruning
--- PASS: TestPartitionPruning (0.00s)
=== RUN   TestRoutingAndConsistentHashRelocation
    sharding_test.go:107: Hash Modulo moved 3756 / 5000 keys (75.12%)
    sharding_test.go:108: Consistent Hash moved 800 / 5000 keys (16.00%)
--- PASS: TestRoutingAndConsistentHashRelocation (0.00s)
=== RUN   TestClusterScatterGatherAndGSI
--- PASS: TestClusterScatterGatherAndGSI (0.00s)
=== RUN   TestIDGenerators
--- PASS: TestIDGenerators (0.00s)
=== RUN   TestConcurrentClusterAccess
--- PASS: TestConcurrentClusterAccess (0.00s)
PASS
ok  	labs/32-database-sharding-and-partitioning/tests	0.088s
```

Command:
```bash
go test -race ./...
```
Output:
```text
ok  	labs/32-database-sharding-and-partitioning/tests	(cached)
```

## Coverage Verification

1. **Happy Path**:
   - `TestPartitionPruning`: Verifies insert and range query matching 2 records in Feb 2026.
   - `TestClusterScatterGatherAndGSI`: Verifies point lookup by ShardKey and GSI email lookup.
   - `TestIDGenerators`: Verifies UUIDv7 timestamp order and sequence block allocation continuity.

2. **Failure & Edge Cases**:
   - `TestPartitionPruning`: Verifies partition drop (`DropPartition`) and partition count updates.
   - `TestClusterScatterGatherAndGSI`: Verifies pre-canceled context handling in scatter-gather query (`ShardResponded == 0`).

3. **Routing & Scaling Behavior**:
   - `TestRoutingAndConsistentHashRelocation`: Compares 3-node to 4-node scale out for 5,000 keys. Modulo moves ~75.12% of keys while Consistent Hashing moves ~16.00% of keys.

4. **Concurrency & Race Conditions**:
   - `TestConcurrentClusterAccess`: 200 concurrent goroutines performing inserts and lookups concurrently under `-race`. Passed cleanly.
