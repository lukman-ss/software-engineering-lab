# Test Audit

## Coverage of Claims

### Partition Pruning
- Happy path: TestPartitionPruning inserts into Jan/Feb, queries Feb -> 2 records, 1 partition scanned.
- Edge case: none (no query with no overlap)
- Negative case: none (no test for timestamp before/after all partitions)
Assessment: PASS for pruning, missing edge-case coverage.

### ModuloRouter Relocation
- TestRoutingAndConsistentHashRelocation tests 3->4 node scale on 5000 keys.
- Observed: 3756/5000 keys moved (75.12%) matches expectation (~75%).
- No test for empty shard list, single shard, or removal scenario.
Assessment: PASS for baseline, missing edge cases.

### ConsistentHashRouter Relocation
- Same test: 800/5000 keys moved (16.00%) logged (actual run).
- Acceptance: 5-40% range per test -> PASS.
- No test for vnodeCount = 0, vnodeCount = 1, or removal edge.
Assessment: PASS for baseline, missing edge cases.

### Scatter‑gather & GSI
- TestClusterScatterGatherAndGSI inserts three records, verifies:
  - direct GetByShardKey works.
  - GSI lookup works.
  - scatter-gather by email works.
  - pre-canceled context yields 0 shard responses.
Assessment: PASS.

### ID Generation
- TestIDGenerators:
  - UUIDv7 generation length and ordering (u1 < u2).
  - SequenceBlockAllocator yields 25 sequential IDs.
Assessment: PASS.

### Concurrency
- TestConcurrentClusterAccess: 200 goroutines each insert + GetByShardKey + GetByEmailUsingGSI.
- Final total matches ops -> PASS.
- Gaps: no concurrent AddShardNode vs Insert, no concurrent RebalanceData vs Insert, no concurrent RemoveShard.
Assessment: PASS for exercised paths, missing concurrent-modification edge cases.

### Missing Test Classes
1. No direct unit tests for ModuloRouter or ConsistentHashRouter alone.
2. No test for empty shard slice (GetShard returns ErrShardNotFound).
3. No test for RemoveShard on live router.
4. No test for RebalanceData (dead code).
5. No test for Table.Insert with timestamp before/after all ranges.
6. No test for Table.QueryRange with zero-width window.
7. No negative test for malformed UUID input (ExtractTimeFromUUIDv7).
8. No test for context deadline exceeded in scatter-gather (only pre-cancel).
Assessment: MISSING_TEST and MISSING_EDGE_CASE gaps.

Test Quality
- Tests use actual implementation, no mocks; they exercise real concurrency and logic.
- Test logging shows actual numbers, not canned expectations (good).
- Test names are clear.

## Summary
Core behaviors (partition pruning, routing relocation, scatter-gather/GSI, ID generation) are unit‑tested and pass. Missing tests for edge cases, error paths, and concurrent modification leave some risk undetected. No test probes the deadlock hazard (Insert vs RebalanceData) or nil‑map panic in RebalanceData.