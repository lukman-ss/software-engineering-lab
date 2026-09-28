# Code Audit

## Finding 1 — Lock-order inversion between Insert and RebalanceData

Location: internal/sharding/sharding.go:280 Insert (router lock -> c.mu RLock), internal/sharding/sharding.go:416 RebalanceData (c.mu Lock -> router lock)
Claimed Behavior: Cluster supports concurrent insert and rebalance operations.
Observed Implementation: Insert acquires the router lock first then the cluster RWMutex (RLock). RebalanceData acquires the cluster RWMutex (Lock) first then the router lock inside the same critical section. These two lock orderings are inverted. If an Insert and RebalanceData run concurrently, each goroutine can hold one lock and wait for the other -> permanent deadlock.
Assessment: FAIL
Severity: HIGH
Notes: Not detected by the race detector (it is a logical deadlock, not a data race). No test exercises Insert concurrently with RebalanceData, so the hazard is unproven and latent.

## Finding 2 — RebalanceData nil-map panic on missing target shard

Location: internal/sharding/sharding.go:432-434
Claimed Behavior: RebalanceData migrates all records back onto correct shards after a scale event.
Observed Implementation: After resetting every shard, the loop calls `targetShard, _ := c.router.GetShard(...)` and then `c.shards[targetShard].Put(rec)`. The error from GetShard is discarded and the map lookup assumes targetShard exists in c.shards. If a shard was removed from the cluster but not the router (or GetShard returns ErrShardNotFound), `c.shards[targetShard]` is nil and `.Put` panics on a nil-pointer dereference.
Assessment: FAIL
Severity: HIGH
Notes: Latent — RebalanceData is never invoked by the demo or tests. The function is effectively dead code with a panic surface.

## Finding 3 — ExtractTimeFromUUIDv7 silently returns zero time

Location: internal/idgen/idgen.go:90-113
Claimed Behavior: ExtractTimeFromUUIDv7 extracts a millisecond timestamp from a UUIDv7 string.
Observed Implementation: The primary path uses `fmt.Sscan(clean[:12], "%12x", &high)` to parse 12 hex chars into an [8]byte. On any Sscanf success the function returns `time.Time{}, nil` (zero value, no error); the manual fallback loop only runs when Sscanf fails outright. For a well-formed UUIDv7 the primary parse "succeeds" and the correct timestamp is never returned.
Assessment: FAIL
Severity: MEDIUM
Notes: Dead code — not called by demo or tests. Still a broken implementation sitting in the module.

## Finding 4 — Table.Insert holds table reader lock while mutating partition

Location: internal/partitioning/table.go:73-83
Claimed Behavior: Thread-safe concurrent inserts.
Observed Implementation: Insert acquires `t.mu.RLock` and then calls `p.Insert(rec)` which acquires the partition's `p.mu.Lock`. Multiple concurrent inserts can therefore proceed in parallel and are only serialized at the partition level — correct and safe, lock ordering is consistent (table -> partition) across Insert and QueryRange.
Assessment: PASS
Severity: LOW
Notes: Acceptable design; the RLock on the table is held only for the partition-selection loop, partitions have their own mutexes.

## Finding 5 — Scatter-gather bounded goroutine fan-out

Location: internal/sharding/sharding.go:359-384
Claimed Behavior: Parallel broadcast across shards with context cancellation.
Observed Implementation: A goroutine is spawned per shard (fan-out bounded by shard count). Each goroutine re-checks the context during iteration and reports errors via a buffered channel closed by `wg.Wait`. Results are aggregated after all goroutines complete. Correct.
Assessment: PASS
Severity: LOW
Notes: Shard count is small and config-controlled; goroutine count is bounded.

## Finding 6 — ConsistentHashRouter ring rebuild on AddShard

Location: internal/sharding/sharding.go:115-122
Claimed Behavior: AddShard inserts a shard onto the ring.
Observed Implementation: Appends vnodeCount entries and re-sorts the entire ring on every AddShard call — O(V log V) per add. Correct but inefficient; with many scale operations this is a bottleneck.
Assessment: WARNING (performance, not correctness)
Severity: LOW
Notes: Acceptable for the demo scale (4-5 shards).

## Finding 7 — GetShardCounts returns consistent snapshot

Location: internal/sharding/sharding.go:405-414
Claimed Behavior: Report per-shard record counts.
Observed Implementation: Acquires cluster RLock then calls `s.Count()` (which takes the per-shard RLock). Lock ordering cluster -> shard is consistent with Insert. Safe.
Assessment: PASS
Severity: LOW

## Finding 8 — UUIDv7 generation

Location: internal/idgen/idgen.go:13-39
Claimed Behavior: RFC 9562 time-ordered UUIDv7.
Observed Implementation: 48-bit big-endian millisecond timestamp in bytes 0-5, version nibble (0x70) in high nibble of byte 6, variant (10) in byte 8, remaining bytes random. Format string produces correct 8-4-4-4-12 layout. Test verifies u1 < u2 lexicographically.
Assessment: PASS
Severity: LOW

## Finding 9 — SequenceBlockAllocator

Location: internal/idgen/idgen.go:57-73
Claimed Behavior: Block allocation of monotonic IDs.
Observed Implementation: Returns base+1..base+blockSize, then fetches a new block via fetcher when exhausted. MemoryCentralSequence.AllocateBlock returns cursor+1 and advances. Test confirms monotonic sequential output across block boundary.
Assessment: PASS
Severity: LOW
Notes: No overflow guard, but acceptable for in-memory demo.
