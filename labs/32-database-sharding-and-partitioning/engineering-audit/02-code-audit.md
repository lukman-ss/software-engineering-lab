# Code Audit

## Finding 1

Location: `internal/sharding/sharding.go:337-382` (`ScatterGatherBroadcast`)
Claimed Behavior: Scatter-gather broadcast across all shards with parallel execution.
Observed Implementation: Launches goroutines per shard, collects results into a buffered channel of size `len(shards)`, and waits with `sync.WaitGroup`. However, no timeout or context cancellation is supported; if a shard operation hangs, the query blocks indefinitely.
Assessment: WARNING
Severity: LOW
Notes: Acceptable for an in-memory educational lab without network latency, but standard production scatter-gather requires context timeout / deadline propagation.

## Finding 2

Location: `internal/sharding/sharding.go:279-298` (`Cluster.Insert` and `GlobalSecondaryIndex.Index`)
Claimed Behavior: Sharded record insert with Global Secondary Index maintenance.
Observed Implementation: Insert acquires shard lock, then writes to GSI map under separate lock. The operation is not atomic across shard and GSI (no 2PC / transaction).
Assessment: PASS
Severity: LOW
Notes: Accurately scoped and documented in `engineering/02-implementation-notes.md` under Known Limitations (no 2PC / cross-shard transactions).

## Finding 3

Location: `internal/sharding/sharding.go:83-176` (`ConsistentHashRouter`)
Claimed Behavior: Ring-based consistent hashing with virtual nodes and binary search routing.
Observed Implementation: Uses `fnv.New64a` hash of `shardID#vnodeIndex`, maintains sorted ring slice, and performs `sort.Search` for clockwise lookup. Dynamic `AddShard` and `RemoveShard` safely lock `sync.RWMutex`.
Assessment: PASS
Severity: LOW
Notes: Clean, idiomatic, thread-safe Go implementation.

## Finding 4

Location: `internal/partitioning/table.go:93-119` (`Table.QueryRange`)
Claimed Behavior: Logical range partitioning with partition pruning.
Observed Implementation: Scans only partitions whose range overlaps `[start, end)` (`start.Before(p.Range.End) && end.After(p.Range.Start)`), correctly counting `PartitionsScanned` vs `TotalPartitions`.
Assessment: PASS
Severity: LOW
Notes: Accurately simulates storage-engine level partition pruning.

## Finding 5

Location: `internal/idgen/idgen.go:13-39` (`NewUUIDv7`)
Claimed Behavior: RFC 9562 compliant time-ordered UUIDv7.
Observed Implementation: Encodes 48-bit millisecond timestamp in high bits, sets 4-bit version `0111`, and sets 2-bit variant `10`.
Assessment: PASS
Severity: LOW
Notes: Produces lexicographically sortable identifiers consistent with RFC 9562 specifications.

## Finding 6

Location: `internal/idgen/idgen.go:42-73` (`SequenceBlockAllocator`)
Claimed Behavior: Vitess-style sequence block allocation avoiding single-point DB sequence contention.
Observed Implementation: Atomically allocates sequential IDs locally until `max` block limit, refreshing from central coordinator via `fetcher` callback.
Assessment: PASS
Severity: LOW
Notes: Implementation correctly models chunked sequence allocation with mutex synchronization.
