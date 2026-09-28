# Docs vs Code Audit

## Comparison Matrix

| Claim / Specification | Documentation Source | Implementation / Test | Verdict |
| :--- | :--- | :--- | :--- |
| Asynchronous vs Synchronous Replication | README.md, 01-design.md | `cluster.AsyncReplication`, `cluster.SyncReplication`, `TestSynchronousReplication_Freshness` | MATCH |
| Stale Read Anomaly Demonstration | README.md, 01-design.md | `cmd/demo/main.go`, `TestNaiveReplicationLag_StaleRead` | MATCH |
| Time-Based Sticky Session Routing | README.md, 01-design.md | `router.ReadWithStickySession`, `TestStickySessionRouting` | MATCH |
| Causal Token / LSN Wait Routing | README.md, 01-design.md | `router.ReadWithToken`, `TestReadWithToken_LSN` | MATCH |
| Lag-Aware Dynamic Replica Routing | README.md, 01-design.md | `router.ReadLagAware`, `TestReplicaLagThreshold_Fallback` | MATCH |
| Execution Output in Engineering Notes | 03-execution-result.md | Matches live execution runs verbatim | MATCH |

## Identified Discrepancies

- None detected.
- Code identifiers, API signatures, package structure, and CLI output match documentation accurately.
