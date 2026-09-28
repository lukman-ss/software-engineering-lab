# Research Plan: Read Replicas & Replication Lag

## Research Topic

Read Replicas & Replication Lag — Scaling read-heavy workloads while maintaining data consistency and avoiding stale reads.

## Objective

To investigate read replica architectures for handling read-heavy workloads, focusing on:

1. The problem of replication lag in asynchronous replication
2. Consistency models: eventual vs. causal vs. linearizable
3. Strategies for handling read-your-own-writes and monotonic reads
4. Database-specific replication implementations (PostgreSQL, MySQL, MongoDB, Aurora)
5. Application-level routing patterns (sticky routing, session-based consistency)

## Research Questions

1. How does asynchronous replication cause replication lag, and what are typical lag ranges (ms to hours)?
2. What consistency models do major databases provide by default vs. with configuration?
3. How do read-your-own-writes and monotonic read guarantees work at database and application levels?
4. What are the trade-offs between synchronous vs. asynchronous replication for read replicas?
5. What application patterns can mitigate replication lag without sacrificing scalability?

## Search Strategy

1. **Primary sources**: Database official documentation (PostgreSQL, MySQL, MongoDB, Aurora, RDS)
2. **Secondary sources**: Industry best practices, replication pattern papers
3. **Case studies**: Real-world replication lag issues and solutions

## Expected Primary Sources

- PostgreSQL documentation: streaming replication, hot standby, async sync modes
- MySQL/Oracle: InnoDB replication, GTIDs, semi-sync
- MongoDB: causal consistency sessions, read concerns
- AWS RDS/Aurora: read replica implementation
- Microsoft Research: Replicated Data Consistency Explained Through Baseball

## Risks / Unknowns

- Evidence on real-world lag distributions for various workloads
- Limited evidence on application-level routing effectiveness in production
- Unclear best practices for critical read threshold (when to route to primary)
- Contextual differences between OLTP vs. analytics workloads
