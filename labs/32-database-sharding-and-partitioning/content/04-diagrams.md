# Diagrams

## 1. Logical Partitioning vs Physical Sharding Architecture

Diagram perbandingan arsitektural antara Logical Partitioning pada single engine dan Physical Sharding lintas beberapa node independen:

```text
========================================================================================
1. LOGICAL PARTITIONING (Single Database Instance)
========================================================================================
                              [ Client / Application ]
                                         │
                                         ▼ SQL Query (WHERE created_at BETWEEN ...)
                            ┌────────────────────────┐
                            │ Database Query Planner │
                            └───────────┬────────────┘
                                        │ (Partition Pruning)
                ┌───────────────────────┼───────────────────────┐
                ▼                       ▼                       ▼
     ┌─────────────────────┐ ┌─────────────────────┐ ┌─────────────────────┐
     │ Partition: p2026_01 │ │ Partition: p2026_02 │ │ Partition: p2026_03 │
     │  [ Jan 01 - Feb 01 ]│ │  [ Feb 01 - Mar 01 ]│ │  [ Mar 01 - Apr 01 ]│
     │      (PRUNED)       │ │     (SCANNED: 1)    │ │      (PRUNED)       │
     └─────────────────────┘ └─────────────────────┘ └─────────────────────┘
     ───────────────────────────────────────────────────────────────────────
     Physical Host: Single Server (Shared Memory, Shared Disk, Single Engine)


========================================================================================
2. PHYSICAL DATABASE SHARDING (Distributed Multi-Instance Cluster)
========================================================================================
                              [ Client / Application ]
                                         │
                                         ▼
                            ┌────────────────────────┐
                            │  Cluster Shard Router  │
                            │ (Consistent Hash Ring) │
                            └───────────┬────────────┘
                                        │
                ┌───────────────────────┼───────────────────────┐
                │ (Route: tenant-A)     │ (Route: tenant-B)     │ (Route: tenant-C)
                ▼                       ▼                       ▼
     ┌─────────────────────┐ ┌─────────────────────┐ ┌─────────────────────┐
     │   Shard Node 0      │ │   Shard Node 1      │ │   Shard Node 2      │
     │   (Host A / Port 1) │ │   (Host B / Port 2) │ │   (Host C / Port 3) │
     │  Independent Engine │ │  Independent Engine │ │  Independent Engine │
     │     Isolated DB     │ │     Isolated DB     │ │     Isolated DB     │
     └─────────────────────┘ └─────────────────────┘ └─────────────────────┘
```

---

## 2. Hash Modulo vs Consistent Hash Ring Resizing (Data Migration)

Diagram visualisasi pergerakan data saat cluster bertambah dari 4 node menjadi 5 node:

```text
========================================================================================
A. HASH MODULO (N % M): Scale-out 4 Shards -> 5 Shards
========================================================================================
Keys: [k0, k1, k2, k3, k4, k5, k6, k7, k8, k9, ...]
Before (M = 4): k_i % 4
After  (M = 5): k_i % 5

Result: ~79.84% of all keys change destinations!
[k0] -> Shard 0 -> Shard 0 (Stay)
[k1] -> Shard 1 -> Shard 1 (Stay)
[k2] -> Shard 2 -> Shard 2 (Stay)
[k3] -> Shard 3 -> Shard 3 (Stay)
[k4] -> Shard 0 -> Shard 4 (MOVED)
[k5] -> Shard 1 -> Shard 0 (MOVED)
[k6] -> Shard 2 -> Shard 1 (MOVED)
[k7] -> Shard 3 -> Shard 2 (MOVED)
... (Massive cluster-wide data shuffle)


========================================================================================
B. CONSISTENT HASH RING WITH VIRTUAL NODES (Scale-out 4 Shards -> 5 Shards)
========================================================================================
Ring Space: [0 ................................................. 2^64 - 1]

                      (vnode: s1#12)
                            o
                (vnode: s0#4) \   o (vnode: s2#88)
                             \ /
                     --- CIRCULAR RING ---
                            / \
              (NEW: s4#01) o   \ o (vnode: s3#55)
                                o
                           (vnode: s0#99)

Result: Adding Node 4 (s4) only intercepts keys immediately preceding its vnodes.
Theoretical Migration: ~1/(N+1) = ~20.00%
Actual Lab Result: ~12.00% to 16.00% keys moved.
Remaining ~84% to 88% keys stay on their original shards!
```

---

## 3. Querying Non-Sharded Attributes: Scatter-Gather vs Global Secondary Index

Diagram alur eksekusi saat mencari record berdasarkan non-sharding key (`email`):

```text
========================================================================================
OPTION 1: SCATTER-GATHER (Full Broadcast)
Query: SELECT * WHERE email = 'carol@example.com' (No ShardKey provided)
========================================================================================
                       [ Client / Application ]
                                  │
                                  ▼
                      ┌───────────────────────┐
                      │ Cluster Router Engine │
                      └───────────┬───────────┘
                                  │ Broadcast in Parallel (Goroutines + sync.WaitGroup)
                ┌─────────────────┼─────────────────┬─────────────────┐
                ▼                 ▼                 ▼                 ▼
         ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐
         │   Shard 0   │   │   Shard 1   │   │   Shard 2   │   │   Shard 3   │
         │  (Scan DB)  │   │  (Scan DB)  │   │  (Scan DB)  │   │  (Scan DB)  │
         │  Matched: 0 │   │  Matched: 0 │   │  Matched: 1 │   │  Matched: 0 │
         └──────┬──────┘   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘
                └─────────────────┼─────────────────┴─────────────────┘
                                  ▼ Aggregate Results
                      ┌───────────────────────┐
                      │ Combined Result: 1 rec│ (High network & CPU overhead on all nodes)
                      └───────────────────────┘


========================================================================================
OPTION 2: GLOBAL SECONDARY INDEX (Lookup Vindex Point-Lookup)
Query: SELECT * WHERE email = 'carol@example.com'
========================================================================================
                       [ Client / Application ]
                                  │
                                  ▼
                      ┌───────────────────────┐
                      │ 1. GSI Lookup Engine  │
                      │ email -> ShardKey     │ (In-memory map / Lookup Vindex Table)
                      └───────────┬───────────┘
                                  │ Resolved: ShardKey = "tenant-C"
                                  ▼
                      ┌───────────────────────┐
                      │ 2. Consistent Router  │
                      │ ShardKey -> Shard 2   │
                      └───────────┬───────────┘
                                  │ 3. Direct Point-Lookup (Targeted single node)
                                  ▼
                           ┌─────────────┐
                           │   Shard 2   │
                           │  (Fetch ID) │
                           └─────────────┘
                      (Only 1 Shard contacted; 0 broadcast overhead!)
```

---

## 4. Distributed ID Generation Mechanics

Diagram perbandingan struktur UUIDv7 (RFC 9562) vs Vitess-style Sequence Block Allocation:

```text
========================================================================================
A. RFC 9562 UUIDv7 (128-bit Time-Ordered Layout)
========================================================================================
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           unix_ts_ms                          | (Bits 0-31: MSB timestamp)
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          unix_ts_ms           |  ver (0111)   |    rand_a     | (Bits 32-47: LSB timestamp +
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+  Bits 48-51: Ver 7 + rand)
|var(10)|                             rand_b                    | (Bits 64-65: Variant RFC 4122
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+  Bits 66-95: Random bits)
|                            rand_b                             | (Bits 96-127: Random bits)
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
Properties: Globally Unique, Zero Network Coordination, Lexicographically Time-Ordered.


========================================================================================
B. VITESS-STYLE SEQUENCE BLOCK ALLOCATION
========================================================================================
                    ┌────────────────────────────┐
                    │ Central Sequence Generator │ (Cursor = 1000)
                    └─────────────┬──────────────┘
                                  │
          ┌───────────────────────┴───────────────────────┐
          │ AllocateBlock(Size=500)                       │ AllocateBlock(Size=500)
          ▼                                               ▼
┌───────────────────────────┐                   ┌───────────────────────────┐
│       App Worker A        │                   │       App Worker B        │
│ Range Allocated: [1 - 500]│                   │ Range: [501 - 1000]       │
│ Local Cursor: 1, 2, 3...  │                   │ Local Cursor: 501, 502... │
│ (Zero DB calls for 500 IDs│                   │ (Zero DB calls for 500 IDs│
└───────────────────────────┘                   └───────────────────────────┘
```
