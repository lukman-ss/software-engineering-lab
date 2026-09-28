# Contradictions

## No Material Contradictions Discovered

All sources converge on the core model: async replication → lag → eventual consistency; strong consistency → primary reads or sync apply.

## Areas of Divergence (Not Contradictions)

### Synchronous Replication Semantics

- **Source 1 (PostgreSQL)**: `synchronous_commit=remote_write` waits for standby to receive + write to OS; `remote_apply` waits for replay (visibility). Default is `on` = local commit wait only.
- **Source 7 (Vitess)**: MySQL semi-sync "guarantees that at least one replica has the transaction in its relay log" but "has not necessarily been applied yet."

**Assessment**: Same underlying distinction — WAL receipt vs. WAL apply. PostgreSQL offers granular control; MySQL semi-sync provides durability guarantee but not visibility. No contradiction.

### Lag Ranges

- **Source 1 (PostgreSQL streaming)**: "typically under one second assuming the standby is powerful enough"
- **Source 3 (Azure)**: "from a few seconds to minutes... could extend to hours" in heavy workloads
- **Source 2 (AWS RDS)**: "Read Replica Lag" metric in seconds (typical healthy case)

**Assessment**: Different workload profiles. Low-latency streaming can achieve <1s; heavy write or geo-distant replicas degrade to seconds/hours. Sources are complementary, not contradictory.

### Read Replica Use Cases

- **Source 2 (AWS RDS)**: Recommends read replicas for "scaling beyond compute or I/O capacity," "business reporting," "disaster recovery"
- **Source 3 (Azure)**: Warns "aren't intended for synchronous replication scenarios requiring up-to-the-minute data accuracy"

**Assessment**: Consistent. Read replicas are for read scaling where staleness is acceptable; not for strong consistency. Different vendors emphasize different use cases.

## Uncertainty Areas Requiring Judgment

### Optimal Replication Lag Threshold for Read-Your-Own-Writes

- **No source provides**: A single threshold value for "acceptable" lag before routing reads to primary
- **Terry et al. (10)**: Defines session guarantees qualitatively, not with ms thresholds
- **Cloud vendors**: Report lag but don't define acceptable bounds

**Assessment**: Lag threshold is application-specific. Financial systems may require <100ms; analytics tolerate seconds. Engineering judgment required.

### Sticky Session Timeout Duration

- **Lab spec**: Suggests "5 seconds" for sticky routing after write
- **No authoritative source**: Validates this specific duration

**Assessment**: Timeout depends on observed replication lag in production; 5s is a reasonable starting point for same-region deployments but must be tuned per workload.

### Causal Consistency Implementation Across Databases

- **Source 4 (MongoDB)**: Client sessions with `clusterTime`/`operationTime` tracking
- **Source 1 (PostgreSQL)**: No built-in session-based consistency; requires application-level LSN tracking
- **Source 6 (MySQL)**: GTID-based lag monitoring; no built-in session consistency

**Assessment**: Not a contradiction — different databases offer different levels of native support for session consistency. MongoDB has the most mature causal consistency support; PostgreSQL/MySQL require application-level implementation.
