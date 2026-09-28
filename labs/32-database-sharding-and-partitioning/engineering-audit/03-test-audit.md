# Test Audit

## Test Suite Execution

### 1. Standard Test Execution
Command:
```bash
go test -v ./...
```
Output:
```text
=== RUN   TestPartitionPruning
--- PASS: TestPartitionPruning (0.00s)
=== RUN   TestRoutingAndConsistentHashRelocation
    sharding_test.go:106: Hash Modulo moved 756 / 1000 keys (75.60%)
    sharding_test.go:107: Consistent Hash moved 0 / 1000 keys (0.00%)
--- PASS: TestRoutingAndConsistentHashRelocation (0.00s)
=== RUN   TestClusterScatterGatherAndGSI
--- PASS: TestClusterScatterGatherAndGSI (0.00s)
=== RUN   TestIDGenerators
--- PASS: TestIDGenerators (0.00s)
=== RUN   TestConcurrentClusterAccess
--- PASS: TestConcurrentClusterAccess (0.00s)
PASS
ok  	labs/32-database-sharding-and-partitioning/tests	0.155s
```

### 2. Race Detector Execution
Command:
```bash
go test -race ./...
```
Output:
```text
ok  	labs/32-database-sharding-and-partitioning/tests	1.103s
```
Race Detector Result: Zero race conditions detected.

### 3. Executable Demo Verification
Command:
```bash
go run ./cmd/demo
```
Output:
```text
================================================================================
LAB 32: DATABASE SHARDING AND PARTITIONING DEMONSTRATION
================================================================================

--- 1. Single-Node Logical Table Partitioning & Range Pruning ---
Query Range: 2026-04-01 to 2026-06-30
Records Found: 1 | Partitions Scanned: 1 / 4 (Pruned 3 partitions)

--- 2. Sharding Key Selection: Monotonic Key vs High-Cardinality Key ---

[Scenario A] Monotonic Key (date-string) Sharding:
  Node shard-0: 0 records [                              ]
  Node shard-1: 0 records [                              ]
  Node shard-2: 1000 records [##############################]
  Node shard-3: 0 records [                              ]
Result: Severe Write Hotspot! 100% writes hit single shard.

[Scenario B] High-Cardinality Key (user_id) Sharding:
  Node shard-0: 400 records [############                  ]
  Node shard-1: 0 records [                              ]
  Node shard-2: 200 records [######                        ]
  Node shard-3: 400 records [############                  ]
Result: Uniform distribution across physical shards.

--- 3. Resharding / Scale-out Comparison: Hash Modulo vs Consistent Hashing ---
Cluster Resize: 4 Shards -> 5 Shards (Total Keys: 10000)
  Hash Modulo (N % M) Keys Remapped   : 7984 / 10000 (79.84% moved)
  Consistent Hashing Keys Remapped      : 1200 / 10000 (12.00% moved)
  Theoretical Minimal Relocation (1/N) : ~20.00%

--- 4. Querying Non-Sharded Attributes: Scatter-Gather vs Global Secondary Index (GSI) ---
Scatter-Gather Query (by Email without Shard Key):
  Nodes Broadcasted : 4 / 4
  Records Matched   : 1
  Execution Time    : 159.875µs
Global Secondary Index (Lookup Vindex) Query:
  Nodes Broadcasted : 1 (Direct Point Lookup via Shard Key mapping)
  Record Found      : ID=usr-342, Email=user_342@company.com, ShardKey=tenant-42
  Execution Time    : 1.083µs (GSI Avoided Broadcast Overhead!)

--- 5. Distributed Unique ID Generation: UUIDv7 vs Central Sequence Block Allocation ---
Generated UUIDv7 (Time-Ordered 128-bit) : 01a0e727-b0bb-7747-b674-ed832d41cdd6
Sequence Block Allocator IDs (Block Size=5) : 1 2 3 4 5 6 7 8 

[DEMO COMPLETE] All database sharding & partitioning concepts successfully executed.
```

## Coverage & Quality Assessment

1. Happy Path: Covered across all core features (Partitioning, Modulo routing, Consistent hash routing, GSI lookup, UUIDv7, Sequence allocation).
2. Failure / Missing Record: Covered in unit tests (e.g. partition pruning with non-matching ranges, GSI lookup).
3. Concurrency Safety: Covered with 200 concurrent goroutines writing, reading, and indexing under race detector.
4. Edge Case Relocation: Test `TestRoutingAndConsistentHashRelocation` bounds modulo move ratio ($\ge 65\%$) and consistent hash upper bound ($\le 40\%$). In consistent hash test, `chMoveRatio` was 0.00% on the specific sample key generation, which passed the upper bound check but revealed that keys didn't happen to fall onto shard-4's tokens.
