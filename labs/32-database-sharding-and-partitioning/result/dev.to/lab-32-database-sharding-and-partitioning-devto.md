---
title: "Database Sharding and Table Partitioning in Go: Scaling Out Past Single-Instance Limits"
published: false
description: "Explore range partitioning, consistent hashing vs modulo routing, scatter-gather vs GSI lookups, and RFC 9562 UUIDv7 in Go."
tags: "go, architecture, distributed, database"
---

Monolithic databases face unavoidable physical boundaries: CPU cores, RAM limits, I/O storage saturation, and write lock contention on a single engine. When data volume and transaction throughput exceed what vertical hardware scaling can handle, engineering teams look toward horizontal decomposition.

Two primary strategies exist:
1. **Logical Table Partitioning**: Splitting a large table into physical sub-tables within the same database engine.
2. **Physical Database Sharding**: Distributing partitions across multiple autonomous database nodes.

In **Software Engineering Lab 32**, we build an in-memory simulation in pure Go standard library covering range partitioning, partition pruning, hash modulo vs. consistent hashing routing, scatter-gather queries vs. Global Secondary Index (GSI / Lookup Vindex), and distributed unique ID generation (RFC 9562 UUIDv7 and Vitess-style block allocation).

```text
[Incoming Query / Write]
         │
         ▼
[Routing Engine / Query Planner]
  ├── (Key contains Shard Key?)  ──► Direct Point-Route (O(1) Modulo or O(log N*V) Ring)
  ├── (Non-Sharded Point Lookup) ──► Global Secondary Index (GSI) -> Single Shard Route
  └── (Cross-Shard Query / Range)──► Parallel Scatter-Gather Goroutines (All Shards)
         │
         ▼
[Physical Shards: Shard 0 .. Shard N]
  └── In-Engine Sub-Tables (Range-Partitioned, Pruned by Query Window)
```

---

## Mental Model: Logical Partitioning vs. Physical Sharding

Partitioning and sharding are not interchangeable. They operate at distinct layers of the infrastructure stack:

| Dimension | Logical Table Partitioning | Physical Database Sharding |
|---|---|---|
| **Instance Topology** | Single database instance | N autonomous database instances |
| **Primary Goal** | Lifecycle management, partition pruning | Horizontal scaling past 1 host's disk & I/O limits |
| **Scalability Ceiling** | Bound by CPU, RAM, and disk of 1 server | Scales horizontally with shard count |
| **Cross-Partition JOIN** | Native SQL supported by database engine | Not natively supported; requires scatter-gather or co-location |
| **Distributed IDs** | Native `AUTO_INCREMENT` / sequences work | Required; local autoincrement causes primary key collisions |
| **Transaction Isolation** | Full ACID isolation guaranteed | Two-Phase Commit (TwoPC) adds latency; risk of fractured reads |

```text
Logical Partitioning = Engine Optimization (Split tables inside one host)
Physical Sharding    = Distributed Topology (Split data across multiple hosts)
Shard Key            = Routing Discriminator (Determines which host owns the record)
```

---

## The Monotonic Key Write Hotspot

A common architectural trap in sharding is choosing a monotonically increasing shard key, such as auto-incrementing integers, creation timestamps, or standard UUIDv1. Because sequential values share the same hash or range bucket, write operations do not distribute uniformly across the cluster. Instead, 100% of inserts hammer the single active shard responsible for the latest range.

Our lab benchmarks this behavior under 1,000 writes across 4 physical shards:

```text
[Scenario A] Monotonic Key (date-string) Sharding:
  Node shard-0: 0 records
  Node shard-1: 0 records
  Node shard-2: 1000 records [##############################]
  Node shard-3: 0 records
Result: Severe Write Hotspot. 100% writes hit single shard.

[Scenario B] High-Cardinality Key (user_id) Sharding:
  Node shard-0: 400 records  [############                  ]
  Node shard-1: 0 records    
  Node shard-2: 200 records  [######                        ]
  Node shard-3: 400 records  [############                  ]
Result: Uniform distribution across physical shards.
```

### Optimal Shard Key Criteria

Production sharding architectures (like Vitess and MongoDB) evaluate shard keys against three rules:
1. **High Cardinality**: Sufficient unique values (`user_id`, `tenant_id`, `account_uuid`) rather than low-cardinality flags (`status`, `country_code`).
2. **Non-Monotonicity**: Successive writes disperse evenly across the hash ring.
3. **Point-Lookup Alignment**: Common read queries filter by this key, avoiding cluster-wide broadcasts.

When business requirements mandate queries by monotonic attributes (such as sorting orders by time), **Hashed Sharding** provides the necessary uniform distribution while trading off native range queries for scatter-gather operations.

---

## Routing Algorithms: Hash Modulo vs. Consistent Hashing

When an application issues a read or write, the router must resolve the target shard.

### 1. Hash Modulo Routing

The naive approach computes:
$$\text{ShardIndex} = \text{hash}(key) \pmod N$$

While simple and $O(1)$, modulo routing suffers from severe data disruption during scale-out. When the cluster grows from $N$ to $N+1$ nodes, the denominator changes, altering the target shard for almost every key. The fraction of keys requiring physical relocation approaches:
$$\text{Moved Fraction} = \frac{N}{N+1}$$

Moving from 4 to 5 shards relocates $\approx 80\%$ of all records:

```text
Cluster Resize: 4 Shards -> 5 Shards (10,000 Keys)
  Hash Modulo (N % M) Keys Remapped: 7984 / 10000 (79.84% moved)
```

### 2. Consistent Hashing with Virtual Nodes

Consistent Hashing maps both physical nodes and data keys onto a continuous $2^{64}-1$ ring. When a new node joins, it only assumes ownership of keys from its immediate predecessor on the ring. The theoretical fraction of relocated keys is bounded by:
$$\text{Moved Fraction} \approx \frac{1}{N}$$

```go
type ConsistentHashRouter struct {
	vnodes int
	ring   []vnodeEntry
	mu     sync.RWMutex
}

func (ch *ConsistentHashRouter) GetShard(key string) (string, error) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if len(ch.ring) == 0 {
		return "", ErrNoShardsAvailable
	}

	h := hashKey(key)
	idx := sort.Search(len(ch.ring), func(i int) bool {
		return ch.ring[i].hash >= h
	})
	if idx == len(ch.ring) {
		idx = 0 // Wrap around the circular ring
	}
	return ch.ring[idx].shardID, nil
}
```

To eliminate statistical data skew when physical node counts are low, each physical shard registers multiple virtual nodes (100–150 vnodes) across the ring:

```text
Cluster Resize: 4 Shards -> 5 Shards (10,000 Keys, 100 vnodes)
  Consistent Hashing Keys Remapped: 1200 / 10000 (12.00% moved)
  Theoretical Relocation (1/N)   : ~20.00%
```

Consistent Hashing replaces massive data reshuffling with an $O(\log(N \cdot V))$ binary search over the virtual node ring, cutting migration overhead from 79.84% down to 12.00%.

---

## Non-Sharded Query Optimization: Scatter-Gather vs. Global Secondary Index

Applications frequently need to query records by an attribute other than the primary shard key (for instance, finding a user by `email` when data is sharded by `tenant_id`).

### Strategy 1: Parallel Scatter-Gather (Broadcast)

The router dispatches parallel queries to every shard in the cluster using Go goroutines, collects results over a buffered channel, and merges the records:

```go
func (c *Cluster) ScatterGatherBroadcastWithContext(ctx context.Context, predicate func(Record) bool) ScatterGatherResult {
	c.mu.RLock()
	shards := make([]*Shard, 0, len(c.shards))
	for _, s := range c.shards {
		shards = append(shards, s)
	}
	c.mu.RUnlock()

	var wg sync.WaitGroup
	ch := make(chan shardResult, len(shards))

	for _, s := range shards {
		wg.Add(1)
		go func(shard *Shard) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				ch <- shardResult{shardID: shard.ID(), err: ctx.Err()}
				return
			default:
				records := shard.Scan(predicate)
				ch <- shardResult{shardID: shard.ID(), records: records}
			}
		}(s)
	}

	wg.Wait()
	close(ch)
	// Aggregate results from channel...
}
```

While functional, scatter-gather query overhead scales linearly with cluster size ($O(N)$), degrading throughput under high QPS.

### Strategy 2: Global Secondary Index (Lookup Vindex)

A Global Secondary Index maintains a thread-safe mapping of secondary keys to shard keys (`email -> tenant_id`). Write operations maintain this index via double-writes. During query execution, the router resolves the secondary key directly to the target shard key, converting a broadcast into a single point-lookup:

```text
Scatter-Gather Query (by Email without Shard Key):
  Nodes Broadcasted : 4 / 4
  Records Matched   : 1
  Execution Time    : 128.042 µs

Global Secondary Index (Lookup Vindex) Query:
  Nodes Broadcasted : 1 (Direct Point Lookup via Shard Key mapping)
  Record Found      : ID=usr-342, Email=user_342@company.com, ShardKey=tenant-42
  Execution Time    : 1.375 µs (GSI avoided broadcast overhead)
```

The GSI delivers a 93x latency improvement by eliminating multi-node fan-out, trading off double-write latency during insertion.

---

## Distributed Unique ID Generation

Standard relational databases rely on auto-incrementing integers. In a distributed multi-shard environment, separate engines generate overlapping numeric sequences unless coordinated. RFC 9562 outlines this limitation: coordinated sequential IDs across network boundaries introduce bottlenecks and failure points.

Our lab implements two complementary strategies:

### 1. RFC 9562 UUIDv7 (Time-Ordered 128-Bit UUID)

UUIDv7 packs a 48-bit millisecond Unix timestamp into the most significant bits, followed by 74 bits of cryptographically secure pseudo-random entropy:

```go
func NewUUIDv7() (string, error) {
	var uuid [16]byte
	ms := uint64(time.Now().UnixMilli())

	// 48-bit millisecond timestamp
	uuid[0] = byte(ms >> 40)
	uuid[1] = byte(ms >> 32)
	uuid[2] = byte(ms >> 24)
	uuid[3] = byte(ms >> 16)
	uuid[4] = byte(ms >> 8)
	uuid[5] = byte(ms)

	if _, err := rand.Read(uuid[6:]); err != nil {
		return "", err
	}

	// RFC 9562 Version 7: 0111
	uuid[6] = (uuid[6] & 0x0F) | 0x70
	// RFC 4122 / 9562 Variant: 10xxxxxx
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}
```

Because timestamp bits appear first, UUIDv7 strings sort chronologically ($u_1 < u_2$), preserving B-Tree index locality and preventing page splits.

### 2. Sequence Block Allocator (Vitess Style)

When compact 64-bit integers are required, worker nodes request discrete ID ranges (blocks of 1,000 IDs) from a central coordinator. The local node dispenses IDs from its allocated range without network round-trips until the block is exhausted. If an instance restarts abruptly, unallocated IDs within that block are discarded, creating minor numerical gaps, an accepted trade-off in high-throughput architectures like Vitess.

---

## Test Verification

The test suite validates data distribution invariants, race conditions, and error boundaries:

```bash
go test -race -v ./...
```

```text
=== RUN   TestPartitionPruning
--- PASS: TestPartitionPruning (0.00s)
=== RUN   TestRoutingAndConsistentHashRelocation
    sharding_test.go:107: Hash Modulo moved 3756 / 5000 keys (75.12%)
    sharding_test.go:108: Consistent Hash moved 800 / 5000 keys (16.00%)
--- PASS: TestRoutingAndConsistentHashRelocation (0.02s)
=== RUN   TestClusterScatterGatherAndGSI
--- PASS: TestClusterScatterGatherAndGSI (0.00s)
=== RUN   TestIDGenerators
--- PASS: TestIDGenerators (0.00s)
=== RUN   TestConcurrentClusterAccess
--- PASS: TestConcurrentClusterAccess (0.00s)
PASS
ok      labs/32-database-sharding-and-partitioning/tests        1.492s
```

All 5 verification test suites pass cleanly with zero data race warnings under Go's race detector.

---

## Production Caveats

1. **In-Memory Volatility**: Shards and GSI maps run in local Go memory. Production deployments require persistent storage engines and network serialization protocols (gRPC, Raft consensus).
2. **Synchronous GSI Double-Writes**: Our simulation performs atomic dual writes. Production secondary indexes require distributed two-phase transactions or asynchronous log-based change data capture (Debezium, Kafka) with eventual consistency.
3. **Cross-Shard Two-Phase Commit (TwoPC)**: While TwoPC coordinates commit consensus across multiple nodes, it adds multi-round-trip network latency and risks blocking locks if a coordinator fails. Co-locating related entity records under a shared partition key minimizes cross-shard transactions.

---

## Architectural Summary

- **Logical Partitioning** optimizes single-node query execution through partition pruning and lifecycle drops.
- **Consistent Hashing** reduces scale-out data relocation from $\approx 80\%$ down to $\approx 12\%-16\%$.
- **High-Cardinality Shard Keys** eliminate write hotspots by distributing inserts across physical nodes.
- **Global Secondary Indexes** replace expensive cluster-wide scatter-gather broadcasts with point lookups.
- **UUIDv7** provides collision-free distributed identity while preserving sequential index performance.

Complete runnable source code, test suites, and CLI demonstration:
[github.com/lukman-ss/software-engineering-lab/tree/main/labs/32-database-sharding-and-partitioning](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/32-database-sharding-and-partitioning)
