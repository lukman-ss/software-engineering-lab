# Test Audit

Test File: tests/sharding_test.go
Coverage:

- Happy path:
- TestPartitionPruning: inserts records into partitions, verifies pruning on QueryRange (PartitionsScanned=1).
- TestClusterScatterGatherAndGSI: inserts records, verifies GetByShardKey, GetByEmailUsingGSI, and ScatterGatherBroadcast success.
- TestIDGenerators: verifies UUIDv7 generation, time extraction, sequential IDs.

- Failure path:
- TestClusterScatterGatherAndGSI: uses pre-canceled context to confirm ScatterGatherBroadcastWithContext returns zero responses (ctx cancellation).

- Edge cases:
- TestRoutingAndConsistentHashRelocation: boundary conditions for move ratios (Modulo >=65%, Consistent <=40%/>5%).
- TestIDGenerators: UUID lexicographic ordering check.

- Recovery/Rollback: NONE
- Concurrency:
- TestConcurrentClusterAccess: 200 concurrent goroutines, verifies race-free access and total record count.

- Negative cases: NONE
- Tests: PASS (go test ./...)
- Race detector: PASS (go test -race ./...)
- Gaps: no explicit negative tests for invalid shard key removal, no tests for RebalanceData, no benchmark tests.
