# Sources

## Source 1: PostgreSQL 18 — Streaming Replication & Hot Standby

**Title**: PostgreSQL 18 Documentation: Log-Shipping Standby Servers (Streaming Replication) & Hot Standby
**Publisher**: PostgreSQL Global Development Group
**URL**: https://www.postgresql.org/docs/current/warm-standby.html
**Published**: Current (PostgreSQL 18, September 2026)
**Accessed**: 2026-09-28
**Source Tier**: 1 (Official database documentation)
**Relevance**: Defines streaming replication as asynchronous by default, documents LSN-based lag monitoring, synchronous replication modes (`synchronous_commit`), and hot standby read-only behavior with eventual consistency.

## Source 2: Amazon RDS — Working with DB Instance Read Replicas

**Title**: Working with DB Instance Read Replicas
**Publisher**: Amazon Web Services
**URL**: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_ReadRepl.html
**Published**: Current (accessed 2026)
**Accessed**: 2026-09-28
**Source Tier**: 1 (Cloud vendor official guidance)
**Relevance**: Describes read replicas as read-only copies using asynchronous replication; documents use cases (scaling, reporting, DR), monitoring via Read Replica Lag metric (seconds), cross-region lag, and promotion behavior.

## Source 3: Azure Database for PostgreSQL Flexible Server — Read Replicas

**Title**: Read replicas in Azure Database for PostgreSQL Flexible Server
**Publisher**: Microsoft Learn
**URL**: https://learn.microsoft.com/en-us/azure/postgresql/read-replica/concepts-read-replicas
**Published**: 2026-07-13 (updated 2026-07-16)
**Accessed**: 2026-09-28
**Source Tier**: 1 (Cloud vendor official guidance)
**Relevance**: Documents asynchronous physical replication with replication slots, lag metrics (`physical_replication_delay_in_bytes`, `physical_replication_delay_in_seconds`), replica states, and caveat that replicas are optimized for near real-time with lag from seconds to minutes or hours under heavy load.

## Source 4: MongoDB — Read Isolation, Consistency, and Recency

**Title**: Read Isolation, Consistency, and Recency
**Publisher**: MongoDB Inc.
**URL**: https://www.mongodb.com/docs/manual/core/read-isolation-consistency-recency/
**Published**: Current (accessed 2026)
**Accessed**: 2026-09-28
**Source Tier**: 1 (Official database documentation)
**Relevance**: Documents causal consistency via client sessions guaranteeing read-your-writes, monotonic reads, monotonic writes, and writes-follow-reads. Also documents secondary stale reads and event-consistent model.

## Source 5: Amazon Aurora — DB Clusters Overview

**Title**: Amazon Aurora DB Clusters
**Publisher**: Amazon Web Services
**URL**: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/Aurora.Overview.html
**Published**: Current (accessed 2026)
**Accessed**: 2026-09-28
**Source Tier**: 1 (Cloud vendor official guidance)
**Relevance**: Describes Aurora architecture separating compute from storage; primary (writer) instances and up to 15 Aurora Replicas (reader) sharing same storage volume, enabling read scaling with low replication lag.

## Source 6: PostgreSQL — Configuration Parameters for Replication

**Title**: PostgreSQL 18: 19.6. Replication (Configuration Parameters)
**Publisher**: PostgreSQL Global Development Group
**URL**: https://www.postgresql.org/docs/current/runtime-config-replication.html
**Published**: Current (PostgreSQL 18, September 2026)
**Accessed**: 2026-09-28
**Source Tier**: 1 (Official database documentation)
**Relevance**: Documents `synchronous_standby_names`, `synchronous_commit`, `hot_standby`, `hot_standby_feedback`, `max_standby_streaming_delay`, `wal_receiver_status_interval`, `recovery_min_apply_delay`.

## Source 7: Vitess — MySQL Replication Modes

**Title**: Vitess Docs: Replication
**Publisher**: Vitess (YouTube/Google open-source project)
**URL**: https://vitess.io/docs/24.0/reference/features/mysql-replication/
**Published**: 2025-10-29 (last updated)
**Accessed**: 2026-09-28
**Source Tier**: 1 (Open-source middleware documentation)
**Relevance**: Documents MySQL replication default async, semi-sync option, GTID-based replication, and caveat that semi-sync does not guarantee up-to-date reads—only that a replica has the transaction in its relay log.

## Source 8: GORM DBResolver — Read/Write Splitting

**Title**: DBResolver (GORM)
**Publisher**: GORM (Go ORM)
**URL**: https://gorm.io/docs/dbresolver.html
**Published**: 2026-08-04 (last updated)
**Accessed**: 2026-09-28
**Source Tier**: 2 (Established open-source project documentation)
**Relevance**: Demonstrates application-level read/write splitting via middleware: automatic routing of SELECT to replicas, INSERT/UPDATE/DELETE to primary, manual override via `dbresolver.Write` clause.

## Source 9: Martin Kleppmann — Please Stop Calling Databases CP or AP

**Title**: Please stop calling databases CP or AP
**Publisher**: Martin Kleppmann (Cambridge University)
**URL**: https://martin.kleppmann.com/2015/05/11/please-stop-calling-databases-cp-or-ap.html
**Published**: 2015-05-11
**Accessed**: 2026-09-28
**Source Tier**: 2 (Academic practitioner)
**Relevance**: Explains why async-replicated systems with read-only secondaries are neither CAP-consistent nor CAP-available; defines linearizability; documents why followers returning stale data violate linearizability.

## Source 10: Terry et al. — Replicated Data Consistency Explained Through Baseball

**Title**: Replicated Data Consistency Explained Through Baseball (MSR-TR-2011-137)
**Publisher**: Microsoft Research (Douglas B. Terry et al.)
**URL**: https://www.microsoft.com/en-us/research/wp-content/uploads/2011/10/ConsistencyAndBaseballReport.pdf
**Published**: October 2011
**Accessed**: 2026-09-28
**Source Tier**: 1 (Academic research paper)
**Relevance**: Defines session guarantees including read-your-writes, monotonic reads, monotonic writes, and writes-follow-reads; the theoretical basis for read-your-own-writes consistency.
