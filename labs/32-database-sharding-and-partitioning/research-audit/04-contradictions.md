# Contradiction Audit

## Evaluation of Research Internal Consistency & Source Alignment

No material contradictions found.

### Evaluation of Candidate Nuances

1. **Cross-shard JOIN Capability vs Claim of "No Native JOIN"**
   - Statement A: High-level architectural descriptions frequently state "sharding breaks relational JOINs."
   - Statement B: Vitess Gen4 planner and MongoDB 4.2+ support distributed JOINs / distributed multi-document transactions.
   - Audit Result: Not a contradiction. The research explicitly resolves this: while distributed query engines *can* execute scatter-gather JOINs, they incur substantial network and latency penalties and lack full ACID isolation across shards. Co-location remains the recommended design pattern.

2. **Automatic vs Operator-Triggered Resharding**
   - Statement A: MongoDB uses a background balancer for continuous chunk migration.
   - Statement B: Vitess uses operator-triggered workflows (`vtctldclient Reshard`) with brief read-only cutover.
   - Audit Result: Implementation variance, not a contradiction. Both adhere to the underlying distributed systems necessity of minimal downtime rebalancing.

3. **Consistent Hashing Formula: $1/n$ vs $n/m$**
   - Statement A: Karger et al. 1997 demonstrates adding the $n$-th server redistributes on average $1/n$ of the keys.
   - Statement B: Wikipedia states resizing remaps $n/m$ keys (where $n$ is keys and $m$ is slots).
   - Audit Result: Mathematically identical under standard notation ($K/N$ items per node where $K$ is total keys and $N$ is total nodes). When adding 1 node to $N$ existing nodes, the fraction of relocated keys is $1/(N+1) \approx 1/N$.

All evaluated sources and internal research documents maintain high consistency.
