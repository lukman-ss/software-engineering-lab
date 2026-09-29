# Source Map

## Perbedaan Arsitektur: Logical Partitioning vs Physical Sharding

Research:
`research/05-report.md` → Finding 1: Partitioning vs Sharding Architecture

Implementation:
- `internal/partitioning/table.go` — LogicalTable partitioned by time range
- `internal/sharding/sharding.go` — Shard, Cluster, Router (physical sharding)

Tests:
- `tests/sharding_test.go` → `TestPartitionPruning` (logical partitioning)
- `tests/sharding_test.go` → `TestClusterScatterGatherAndGSI` (physical cluster)

---

## Sharding Key Selection & Monotonic Write Hotspot

Research:
`research/05-report.md` → Finding 2: Sharding Key Selection is Critical for Write Distribution
`research/05-report.md` → Finding 6: Hashed Sharding Mitigates Monotonic-Key Hotspots

Implementation:
- `cmd/demo/main.go` — Section 2: Monotonic Key vs High-Cardinality Key demonstration
- `internal/sharding/sharding.go` → `ModuloRouter.GetShard()`, `ConsistentHashRouter.GetShard()`

Tests:
- `tests/sharding_test.go` → `TestRoutingAndConsistentHashRelocation`

Engineering:
- `engineering/01-design.md` → Failure Scenario: Monotonic write hotspot
- `engineering/03-execution-result.md` → Demo Section 2 output (100% hotspot vs uniform)

---

## Routing Algorithms: Hash Modulo vs Consistent Hashing

Research:
`research/05-report.md` → Finding 3: Consistent Hashing Minimizes Data Movement During Resizing

Implementation:
- `internal/sharding/sharding.go` → `ModuloRouter`, `ConsistentHashRouter`, `vnode`, `hashKey()`
- `internal/sharding/sharding.go` → `AddShard()`, `RemoveShard()`, `GetShard()` (binary search ring lookup)

Tests:
- `tests/sharding_test.go` → `TestRoutingAndConsistentHashRelocation`
  - Hash Modulo: verifikasi >= 65% key moved (actual: 75.12%)
  - Consistent Hash: verifikasi 5%–40% key moved (actual: 16.00%)

Engineering:
- `engineering/01-design.md` → Expected Behavior, Components
- `engineering/02-implementation-notes.md` → Implementation-Specific Choices, Trade-offs
- `engineering/03-execution-result.md` → Demo Section 3: Resharding comparison (79.84% vs 12.00%)

---

## Queries Without Shard Key: Scatter-Gather dan Global Secondary Index

Research:
`research/05-report.md` → Finding 5: Handling Queries Without the Sharding Key

Implementation:
- `internal/sharding/sharding.go` → `ScatterGatherBroadcast()`, `ScatterGatherBroadcastWithContext()`
- `internal/sharding/sharding.go` → `GlobalSecondaryIndex`, `GetByEmailUsingGSI()`
- `internal/sharding/sharding.go` → `Cluster.Insert()` (double-write to GSI on insert with email)

Tests:
- `tests/sharding_test.go` → `TestClusterScatterGatherAndGSI`
  - ScatterGather broadcast ke 3 shard
  - GSI point-lookup untuk `bob@example.com`
  - Context cancellation test (0 shard responded on canceled context)

Engineering:
- `engineering/03-execution-result.md` → Demo Section 4: Scatter-Gather 92µs vs GSI 1µs

---

## Distributed Unique ID Generation

Research:
`research/05-report.md` → Finding 4: Core Complexities Introduced by Sharding (RFC 9562, Vitess Sequences)
`research/02-sources.md` → Source 4: RFC 9562 (IETF Standard)

Implementation:
- `internal/idgen/idgen.go` → `NewUUIDv7()` (RFC 9562 compliant: 48-bit timestamp + version 7 + variant bits + random)
- `internal/idgen/idgen.go` → `SequenceBlockAllocator`, `MemoryCentralSequence`
- `internal/idgen/idgen.go` → `ExtractTimeFromUUIDv7()`

Tests:
- `tests/sharding_test.go` → `TestIDGenerators`
  - UUIDv7 time-ordering: $u_1 < u_2$ setelah 2ms sleep
  - Timestamp extraction dari UUIDv7
  - Sequential IDs dari SequenceBlockAllocator dengan 25 consecutive calls, block size 10

Engineering:
- `engineering/03-execution-result.md` → Demo Section 5: UUIDv7 output & Sequence Block IDs

---

## Resharding / Scale-Out Behavior

Research:
`research/05-report.md` → Finding 7: Resharding Reorganizes Data with Brief Downtime

Implementation:
- `internal/sharding/sharding.go` → `Cluster.AddShardNode()`, `Cluster.RebalanceData()`

Tests:
- `tests/sharding_test.go` → `TestRoutingAndConsistentHashRelocation` (key migration measurement before/after AddShard)

Engineering:
- `engineering/03-execution-result.md` → Demo Section 3: Cluster Resize 4→5 Shards

---

## Concurrent Safety

Implementation:
- `internal/partitioning/table.go` → `sync.RWMutex` per Partition and Table
- `internal/sharding/sharding.go` → `sync.RWMutex` per Shard, Cluster, GlobalSecondaryIndex, Router
- `internal/idgen/idgen.go` → `sync.Mutex` per SequenceBlockAllocator and MemoryCentralSequence

Tests:
- `tests/sharding_test.go` → `TestConcurrentClusterAccess` — 200 goroutines concurrent Insert/Get/GSI, race detector pass

Engineering:
- `engineering/03-execution-result.md` → `go test -race -v ./...` PASS
- `engineering-audit/03-test-audit.md` — Race Detector: PASS
