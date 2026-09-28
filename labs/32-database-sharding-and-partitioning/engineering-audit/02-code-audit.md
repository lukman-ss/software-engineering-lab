# Code Audit

Target Lab: `labs/32-database-sharding-and-partitioning`

## Finding 1

Location: `internal/partitioning/table.go:73-84`
Claimed Behavior: Inserts records into designated partitions based on `CreatedAt` falling into `[Start, End)` interval. Returns error if no partition matches.
Observed Implementation: RLock held while iterating partitions; each partition has its own mutex lock upon append. Range boundary checks properly evaluate `[Start, End)`.
Assessment: PASS
Severity: LOW
Notes: Granular locking avoids contention between partitions during insertions.

## Finding 2

Location: `internal/partitioning/table.go:93-119`
Claimed Behavior: Range query prunes irrelevant partitions; only scans partitions overlapping query interval `[start, end)`.
Observed Implementation: Overlap logic `start.Before(p.Range.End) && end.After(p.Range.Start)` correctly identifies candidate partitions. Inside candidate partitions, individual records are filtered by `[start, end)`. Metrics `PartitionsScanned` and `TotalPartitions` reflect pruning efficiency.
Assessment: PASS
Severity: LOW
Notes: Correct mathematical formulation of interval overlap for half-open intervals.

## Finding 3

Location: `internal/sharding/sharding.go:30-77`
Claimed Behavior: ModuloRouter routes keys using standard FNV-1a hash modulo node count.
Observed Implementation: Thread-safe read/write via RWMutex. Correctly computes `hashKey(key) % N`. Handles empty shard list cleanly by returning `ErrShardNotFound`.
Assessment: PASS
Severity: LOW
Notes: Correctly models naive hash sharding.

## Finding 4

Location: `internal/sharding/sharding.go:84-177`
Claimed Behavior: ConsistentHashRouter assigns virtual nodes on a 64-bit ring, sorting nodes and binary searching with `sort.Search` wrapping at 0.
Observed Implementation: Virtual nodes generated as `shardID#index` and hashed with 64-bit FNV-1a. Ring is sorted upon shard addition. Ring deletion filters out removed shard virtual nodes. `sort.Search` wraps around ring when index equals ring length.
Assessment: PASS
Severity: LOW
Notes: Standard ring topology with configurable virtual nodes.

## Finding 5

Location: `internal/sharding/sharding.go:343-403`
Claimed Behavior: `ScatterGatherBroadcastWithContext` executes parallel shard queries across all shards, respecting `context.Context` cancellation and timeouts.
Observed Implementation: Buffered channel `ch := make(chan shardResult, len(shards))` prevents goroutine blocking on send. Each shard query checks `ctx.Done()` before execution and inside row iteration. `wg.Wait()` ensures all goroutines complete before channel close. Aggregator sums responded shards and matching records.
Assessment: PASS
Severity: LOW
Notes: Concurrency pattern is robust; no goroutine leakage or channel deadlocks.

## Finding 6

Location: `internal/sharding/sharding.go:228-251`
Claimed Behavior: GlobalSecondaryIndex maps secondary attributes (email) directly to shard key for direct routing without broadcast.
Observed Implementation: RWMutex protects `map[string]string`. `Cluster.Insert` indexes email on insert, and `Cluster.GetByEmailUsingGSI` resolves shard key before calling `GetByShardKey`.
Assessment: PASS
Severity: LOW
Notes: Effectively demonstrates lookup table/vindex pattern used in distributed database middlewares.

## Finding 7

Location: `internal/idgen/idgen.go:11-39`
Claimed Behavior: UUIDv7 generates monotonically time-ordered 128-bit identifiers compliant with RFC 9562 (version 7, variant 10).
Observed Implementation: First 48 bits encode Unix millisecond epoch timestamp. Version bits `0x70` masked at byte 6. Variant bits `0x80` masked at byte 8. Cryptographic randomness populates remaining bits. Formatted as standard UUID hyphenated string.
Assessment: PASS
Severity: LOW
Notes: Compliant with RFC 9562 specification.

## Finding 8

Location: `internal/idgen/idgen.go:41-73`
Claimed Behavior: SequenceBlockAllocator reserves blocks of sequential IDs, amortizing round-trips to central sequence generator.
Observed Implementation: Protected by Mutex. Requests batch from `fetcher` when `current >= max`. Generates contiguous sequential IDs without duplicates.
Assessment: PASS
Severity: LOW
Notes: Direct implementation of Vitess sequence block allocation pattern.
