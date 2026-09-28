# Report

## Research Question

**How can an application reliably handle replication lag when using read replicas for scaling, while guaranteeing that a user sees the result of their own write (read‑your‑own‑writes) without sacrificing read‑scale?**

## Executive Summary

1. **Async replication is the default** across PostgreSQL, MySQL, MongoDB, and cloud RDS/Aurora. Replicas receive WAL/logs *after* the primary commits, so reads from replicas are **eventually consistent** (Sources 1, 3, 7).
2. **Strong consistency (linearizability) requires reads to go to the primary** or to a replica that has *applied* the transaction (`remote_apply` in PostgreSQL, Aurora reader with minimal compute lag). Semi‑sync guarantees durability, **not visibility** (Sources 6, 7).
3. **Session‑level guarantees** (read‑your‑writes, monotonic reads) are explicitly provided by MongoDB (Source 4) and formalized in academic work (Source 10). PostgreSQL/MySQL lack built‑in session guarantees; they must be implemented by the application (Sources 1, 6, 7).
4. **Lag monitoring** is native: LSN diff (PostgreSQL), replica‑lag seconds/bytes (Azure, AWS), GTID lag (MySQL). These metrics enable *dynamic routing* based on observed lag (Sources 1, 2, 3, 10).
5. **Application‑level patterns** – sticky routing, latency‑based routing, or explicit read‑after‑write tokens – are recommended. Middleware (GORM DBResolver, Vitess proxy) already supports automatic read/write split with manual overrides (Sources 8, 7).
6. **Trade‑off**: Enabling synchronous replication eliminates lag but adds write latency equal to the network RTT (Source 3). Semi‑sync reduces data‑loss risk without guaranteeing freshness (Source 7).

## Findings

### Finding 1 – Async Replication Causes Variable Lag

- **Claim**: Replication lag can be sub‑second, a few seconds, or hours under heavy load (Sources 1, 3).
- **Evidence**: PostgreSQL streaming lag "typically under one second"; Azure mentions "seconds to minutes, possibly hours"; AWS provides a numeric `Read Replica Lag` metric.
- **Confidence**: HIGH.

### Finding 2 – Guarantees of Read‑Your‑Own‑Writes

- **Claim**: Only primary reads guarantee that a user sees their own write; otherwise use *session guarantees*.
- **Evidence**: Kleppmann’s linearizability argument (Source 9) and MongoDB causal sessions (Source 4) show that without coordination a replica read can be stale. PostgreSQL `synchronous_commit=remote_apply` provides the same guarantee.
- **Confidence**: HIGH.

### Finding 3 – Session Guarantees via LSN/ClusterTime

- **Claim**: Applications can achieve read‑your‑writes by tracking the primary’s commit LSN (PostgreSQL) or `clusterTime` (MongoDB) and routing subsequent reads to a replica that has *caught up*.
- **Evidence**: PostgreSQL function `pg_last_wal_receive_lsn`, Azure replica‑lag metrics, and MongoDB session timestamps (Sources 1, 3, 4).
- **Confidence**: MEDIUM (requires implementation).

### Finding 4 – Middleware Supports Automatic Splitting

- **Claim**: ORM/proxy layers can automatically route reads to replicas and writes to primary, with manual overrides for critical reads.
- **Evidence**: GORM DBResolver (Source 8) and Vitess query router (Source 7) provide read/write split and allow forced primary reads (`dbresolver.Write`).
- **Confidence**: MEDIUM (depends on correct configuration).

### Finding 5 – Synchronous Replication Eliminates Lag but Increases Latency

- **Claim**: Enabling `synchronous_commit=remote_apply` (PostgreSQL) or semi‑sync (MySQL) makes writes wait for replica acknowledgment, increasing round‑trip latency.
- **Evidence**: PostgreSQL docs (Source 3) and Vitess semi‑sync description (Source 7).
- **Confidence**: HIGH.

## Areas of Agreement

- Replicas are **read‑only** and **asynchronously updated** (PostgreSQL, MySQL, AWS RDS, Azure).
- **Primary reads guarantee freshness**; replica reads are eventually consistent.
- **Lag monitoring is essential** and supported natively across all platforms.
- **Application‑level routing** is the canonical solution for read‑your‑own‑writes.

## Areas of Disagreement (None Significant)

All sources align on the core model of async replication and the need for explicit routing for strong consistency. Differences are only in terminology ("semi‑sync" vs. "synchronous_commit=remote_apply") and granularity of configuration.

## Limitations

- No quantitative threshold for acceptable lag; must be determined per workload.
- Lack of built‑in session consistency in PostgreSQL/MySQL means extra engineering effort.
- Cloud vendor docs (AWS/Azure) do not prescribe exact routing algorithms; they only expose metrics.

## Conclusion

To guarantee **read‑your‑own‑writes** while still scaling reads, an application should:
1. **Write to primary** (default).
2. **Record the commit LSN / clusterTime** after each write.
3. **Route subsequent reads** either:
   - Directly to **primary** for critical operations, **or**
   - Use **sticky routing** with a timeout (e.g., 5 s) that checks the replica’s reported lag against the recorded LSN. If the replica is caught up, serve the read; otherwise fallback to primary.
4. **Monitor lag** continuously and alert when lag exceeds a configured SLA.
5. **Optionally enable synchronous replication** (`remote_apply` or semi‑sync) for workloads where write latency can be tolerated.

By combining database‑provided lag metrics, session tokens (LSN/clusterTime), and middleware routing, developers can achieve both high read throughput and the read‑your‑own‑writes guarantee.
