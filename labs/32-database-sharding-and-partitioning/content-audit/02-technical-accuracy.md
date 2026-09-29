# Technical Accuracy & Claim Verification

Audit against:
- `research/05-report.md`
- `engineering/03-execution-result.md`
- `internal/` codebase

## Evaluations

### 1. Logical Partitioning vs Physical Sharding Mental Model
- Draft: Logical table partitioning operates within a single engine/host; physical sharding distributes across independent database nodes.
- Status: **VERIFIED**. Correctly maps concepts without conflating pruning with distributed horizontal scaling.

### 2. Monotonic Key Write Hotspot
- Draft: Monotonic key (`date-string`, sequential ID) sends 100% writes to a single shard; high-cardinality non-monotonic keys (`user_id`) balance writes across shards.
- Evidence: Matches `cmd/demo/main.go` and `engineering/03-execution-result.md` (Scenario A: 1000 records on shard-2; Scenario B: 400, 400, 200 distributed across nodes).
- Status: **VERIFIED**.

### 3. Hash Modulo vs Consistent Hashing Rebalancing
- Draft: Hash Modulo moves ~79.84% keys upon scale-out 4 -> 5 nodes. Consistent Hashing moves 12% to 16% keys.
- Theoretical validation: Formula $N/(N+1)$ for $N=4$ gives $80\%$, consistent with $79.84\%$ measured in demo. Minimal relocation $1/(N+1) = 20\%$, actual in demo is $12.00\%$ and test range is $5\%-40\%$.
- Status: **VERIFIED**.

### 4. Non-Shard Key Queries (Scatter-Gather vs GSI)
- Draft: Scatter-gather broadcasts to 4/4 nodes; GSI maps `email -> shard_key` for 1 direct point lookup.
- Evidence: Matches execution logs in `engineering/03-execution-result.md` (Scatter-Gather: 4/4 shards, GSI: 1 shard).
- Status: **VERIFIED**.

### 5. Distributed Unique ID Generators
- Draft: RFC 9562 UUIDv7 provides 128-bit time-ordered IDs without network coordination; Sequence Block Allocator simulates Vitess Sequences with local allocation from batches.
- Status: **VERIFIED**. Accurately describes standard specifications and gaps on node crash trade-offs.

### 6. Cross-Shard Transactions & TwoPC
- Draft: Notes that TwoPC provides atomicity but not traditional full ACID isolation (risk of fractured reads), citing Vitess documentation.
- Status: **VERIFIED**. Factually sound and backed by Vitess official architectural papers.
