# Code Audit

Target Lab: `labs/32-database-sharding-and-partitioning`

## Finding 1

Location: `internal/partitioning/table.go:73-84`
Claimed Behavior: Inserts records into the partition whose time range [Start, End) encompasses `Record.CreatedAt`.
Observed Implementation: Evaluates each partition's bounds `CreatedAt >= Start && CreatedAt < End`. Successfully inserts into matching partition or returns `ErrNoMatchingPartition`. Table holds RLock while checking, partition takes Lock when appending.
Assessment: PASS
Severity: LOW
Notes: Linear scan over partitions is efficient for typical partition counts (e.g. dozens of monthly or quarterly partitions).

## Finding 2

Location: `internal/partitioning/table.go:93-119`
Claimed Behavior: Partition pruning scans only partitions overlapping query interval `[start, end)`.
Observed Implementation: Overlap check `start.Before(p.Range.End) && end.After(p.Range.Start)` correctly prunes non-overlapping partitions. Returns total partitions count and partitions scanned count for auditing pruning effectiveness.
Assessment: PASS
Severity: LOW
Notes: Both partition count and scanned counts are reported accurately.

## Finding 3

Location: `internal/sharding/sharding.go:30-77`
Claimed Behavior: ModuloRouter distributes keys via `hash(key) % N` and dynamically tracks shard additions/removals.
Observed Implementation: Uses FNV-64a hash of key modulo shard slice length. Protected by `sync.RWMutex`.
Assessment: PASS
Severity: LOW
Notes: `GetShard` returns `ErrShardNotFound` if shard list is empty.

## Finding 4

Location: `internal/sharding/sharding.go:84-172`
Claimed Behavior: ConsistentHashRouter assigns virtual nodes on ring, uses binary search ring lookup, and minimizes key relocation upon shard addition/removal.
Observed Implementation: Configurable virtual nodes (`vnodeCount`), sorted ring slice of `(hash, shardID)`, binary search via `sort.Search`. Ring wraps around to index 0 when search reaches end of ring. Add/Remove shard rebuilds ring appropriately.
Assessment: PASS
Severity: LOW
Notes: Ring updates are synchronized via `sync.RWMutex`.

## Finding 5

Location: `internal/sharding/sharding.go:273-330`
Claimed Behavior: Cluster manages shard nodes, routes inserts, registers secondary keys to GSI, and supports direct shard key or GSI point lookups.
Observed Implementation: `Insert` routes by `rec.ShardKey` and indexes `rec.Email` into `GlobalSecondaryIndex`. `GetByShardKey` and `GetByEmailUsingGSI` execute without cluster broadcast.
Assessment: PASS
Severity: LOW
Notes: Consistent with Vitess lookup vindex design.

## Finding 6

Location: `internal/sharding/sharding.go:343-403`
Claimed Behavior: ScatterGather queries all shards concurrently with context cancellation/timeout handling.
Observed Implementation: Spawns a goroutine per shard with `sync.WaitGroup`, monitors `ctx.Done()`, buffers results into channel of size `len(shards)`, counts completed shard responses and aggregated records.
Assessment: PASS
Severity: LOW
Notes: Gracefully handles canceled contexts without channel blocking or goroutine leaks.

## Finding 7

Location: `internal/idgen/idgen.go:13-39`
Claimed Behavior: Generates RFC 9562 UUIDv7 IDs with 48-bit millisecond timestamp and version/variant bits.
Observed Implementation: Places millisecond Unix time in high 48 bits, reads cryptographic random bytes into remaining bytes, bitmasks version 7 (`0x70`) and variant 2 (`0x80`), formats standard 8-4-4-4-12 hex string.
Assessment: PASS
Severity: LOW
Notes: Lexicographical order correlates with generation timestamp.

## Finding 8

Location: `internal/idgen/idgen.go:42-73`
Claimed Behavior: SequenceBlockAllocator allocates sequential auto-increment IDs in chunks from a coordinator.
Observed Implementation: Fetches `blockSize` IDs atomically when local range exhausted (`current >= max`), dispensing sequentially under mutex lock.
Assessment: PASS
Severity: LOW
Notes: Matches Vitess sequence table block allocation behavior.
