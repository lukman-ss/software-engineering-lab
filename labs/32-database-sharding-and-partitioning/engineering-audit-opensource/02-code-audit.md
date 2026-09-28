# Code Audit

## Finding 1
Location: internal/sharding/sharding.go:163-172 (GetShards)
Claimed Behavior: ConsistentHashRouter.GetShards returns shard list deterministically
Observed Implementation: Uses random map iteration order then sorts; output deterministic.
Assessment: PASS
Severity: LOW
Notes: Minor inefficiency rebuilding slice each call.

## Finding 2
Location: internal/sharding/sharding.go:416-438 (RebalanceData)
Claimed Behavior: Resharding moves records to correct shard per router
Observed Implementation: Reads all records, resets shards, re-inserts by router.GetShard (ignores error). No concurrent writes safety; assumes router not changed concurrently.
Assessment: PASS
Severity: MEDIUM
Notes: Router GetShard error is discarded, causing potential nil target if router empty; in single-threaded rebalance safe.

## Finding 3
Location: internal/idgen/idgen.go:90-113 (ExtractTimeFromUUIDv7)
Claimed Behavior: Extracts timestamp from UUIDv7.
Observed Implementation: Parses 12 hex chars, but returns empty time on success in non-error branch (`return time.Time{}, nil`). Bug — success path returns zero value.
Assessment: FAIL
Severity: HIGH
Notes: Function never returns extracted time; always returns epoch zero on success. Not used in demo/tests but is public API.

## Finding 4
Location: internal/sharding/sharding.go:337-403 (ScatterGatherBroadcastWithContext)
Claimed Behavior: Concurrent parallel query with context cancellation.
Observed Implementation: Checks ctx.Done() in default-select loop; pre-canceled context yields ShardResponded=0, matching test.
Assessment: PASS
Severity: LOW

## Finding 5
Location: internal/sharding/sharding.go:212 (AllRecords)
Claimed Behavior: Thread-safe snapshot of shard data.
Observed Implementation: Returns slice copy under RLock.
Assessment: PASS

## Finding 6
Location: internal/partitioning/table.go:92-119 (QueryRange)
Claimed Behavior: Partition pruning scans overlapping partitions only.
Observed Implementation: Overlap condition correct (`start < p.End && end > p.Start`); row-level filter applied.
Assessment: PASS

## Finding 7
Location: internal/sharding/sharding.go:295-297 (Insert)
Claimed Behavior: GSI updated when email present.
Observed Implementation: GSI index stores email -> shardKey (not shardID), consistent with GetByEmailUsingGSI usage of shardKey in GetByShardKey, which re-routes key.
Assessment: PASS
