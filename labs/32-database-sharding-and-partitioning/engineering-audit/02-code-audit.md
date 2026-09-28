# Code Audit

Target Lab: `labs/32-database-sharding-and-partitioning`

## Finding 1

Location: `internal/partitioning/table.go:73-84`, `internal/partitioning/table.go:93-119`
Claimed Behavior: Logical partition pruning matches only partitions whose range overlaps `[start, end)`.
Observed Implementation: Range overlap condition `start.Before(p.Range.End) && end.After(p.Range.Start)` correctly prunes non-overlapping partitions. Row filtering inside partition verifies `CreatedAt` within boundaries.
Assessment: PASS
Severity: LOW
Notes: Clean synchronization using `sync.RWMutex` on Table and Partition levels.

## Finding 2

Location: `internal/sharding/sharding.go:42-52`, `internal/sharding/sharding.go:143-159`
Claimed Behavior: Deterministic routing via ModuloRouter and ConsistentHashRouter with virtual nodes.
Observed Implementation: ModuloRouter computes `hashKey(key) % N`. ConsistentHashRouter hashes `shardID#i`, sorts ring by hash, and uses `sort.Search` for binary search on ring lookup with wrap-around to ring start.
Assessment: PASS
Severity: LOW
Notes: Standard library `hash/fnv` used without third-party dependencies.

## Finding 3

Location: `internal/sharding/sharding.go:343-403`
Claimed Behavior: Scatter-Gather queries broadcast to all physical shards concurrently, supporting context timeout and cancellation.
Observed Implementation: Spawns goroutine per shard, uses buffered channel sized to `len(shards)`, checks `ctx.Done()`, uses `sync.WaitGroup` to wait, and aggregates responses.
Assessment: PASS
Severity: LOW
Notes: No channel leak or goroutine leak detected; context cancellation handled before and during record iteration.

## Finding 4

Location: `internal/sharding/sharding.go:228-251`, `internal/sharding/sharding.go:322-329`
Claimed Behavior: Global Secondary Index (Lookup Vindex) provides point lookup by non-shard key (e.g. Email) without broadcasting.
Observed Implementation: GSI stores mapping `secondaryKey -> shardKey`. Lookup retrieves `shardKey` and executes point lookup directly against mapped shard.
Assessment: PASS
Severity: LOW
Notes: Concurrently protected with `sync.RWMutex`.

## Finding 5

Location: `internal/idgen/idgen.go:12-38`, `internal/idgen/idgen.go:41-73`
Claimed Behavior: RFC 9562 UUIDv7 generator and Vitess-style sequence block allocator.
Observed Implementation: UUIDv7 embeds 48-bit millisecond timestamp, 4-bit version `0111`, and 2-bit variant `10`. SequenceBlockAllocator fetches blocks atomically under lock to avoid duplicate IDs across nodes.
Assessment: PASS
Severity: LOW
Notes: Verified time extraction utility correctly parses the 48-bit millisecond timestamp.
