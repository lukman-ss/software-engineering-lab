## Gaps

1. MISSING_TEST — No negative test for `partitioning.Table.Insert` with out-of-range timestamp returning `ErrNoMatchingPartition`.
   Severity: MEDIUM — documents failure path but never proven.

2. MISSING_TEST — No negative tests for `Cluster.GetByShardKey` record-not-found, `GetByEmailUsingGSI` email-miss, `Router.GetShard` on empty router (`ErrShardNotFound`), or `SequenceBlockAllocator` fetcher error propagation.
   Severity: MEDIUM — error propagation paths unproven.

3. MISSING_TEST — `Cluster.RebalanceData` (destructive reset + re-route) claimed in design as "resharding migration" but never invoked in tests; only router-level remap ratio is measured, not actual data migration correctness.
   Severity: MEDIUM — core resharding behavior partially unproven at cluster level.

4. MISSING_EDGE_CASE — Partition boundary semantics (Start inclusive / End exclusive) and `DropPartition` of non-existent name (returns false) not tested.
   Severity: LOW.

5. MISSING_EDGE_CASE — `ConsistentHashRouter` duplicate `AddShard`, `RemoveShard` of missing ID, zero/negative vnodeCount defaulting, and `ModuloRouter.RemoveShard` not tested.
   Severity: LOW.

6. UNVERIFIED_RESULT — Demo timing comparison Scatter-Gather vs GSI (`~117µs` vs `~667ns`) is wall-clock on in-memory goroutines, not a real network/latency benchmark; presented only as illustrative overhead, not a throughput benchmark. No fake benchmark, but readers could misread as performance claim.
   Severity: LOW — demo text says "Avoided Broadcast Overhead", which is directionally true; no fabricated numbers beyond run-local timing.
