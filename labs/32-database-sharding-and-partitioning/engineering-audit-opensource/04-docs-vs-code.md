## Docs vs Code

- **README.md** lists components: Partitioning, Sharding routers, GSI, ID generation. Matches code packages `internal/partitioning`, `internal/sharding`, `internal/idgen`.
- **Engineering Design** (engineering/01-design.md) claims: range pruning, hotspot demo, resharding, scatter/gather, GSI, ID gen. All demonstrated in demo and tests.
- Claims
  1. Partitioning with pruning.
  - Verified: `Table.QueryRange` implements pruning; demo prints partition scan numbers.
  2. Sharding key hot‑spot.
  - Verified: `demoShardingHotspots` shows monotonic key all writes to one shard, high‑cardinality key distributes uniformly.
    (Bar chart visual matches claim).
  3. Routing data movement.
  - Verified: `TestRoutingAndConsistentHashRelocation` logs moved keys; demo prints migration percentages aligning with claim.
  - Claim says "Hash modulo moves ~79.84% vs Consistent ~12%" – demo shows these values.
  4. Scatter‑gather vs GSI.
  - Verified: `TestClusterScatterGatherAndGSI` asserts broadcast count 3, GSI direct lookup; demo prints timings.
  - Documentation states "Scatter‑gather broadcast across all shards when sharding key is unknown" – code matches.
  - Documentation states "Global Secondary Index avoids broadcast" – code matches.
  - Execution result prints broadcasted nodes numbers consistent with code.
  5. Distributed ID generation.
  - Verified: `TestIDGenerators` checks UUIDv7 format and ordering; demo prints a UUID and block IDs.

No mismatches found between README/Design and actual implementation or demo output.

Potential minor doc gaps:
- README's "Quick Start" uses `go test -v ./...` but test output uses default without `-v`; still works.
- Design mentions "Two‑phase commit (2PC) not implemented" – not a claim, just limitation.
- Docs don't mention that `RebalanceData` is not demonstrated; but not claimed.

Overall documentation accurately reflects code and behavior.
