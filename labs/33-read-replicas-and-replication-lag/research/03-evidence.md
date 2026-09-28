# Evidence

## Evidence 1: PostgreSQL Streaming Replication is Asynchronous by Default

**Claim**: PostgreSQL streaming replication is asynchronous by default; only WAL shipped after commit creates a data loss window, and lag is monitored via LSN comparison.

**Evidence**: PostgreSQL docs: "Streaming replication is asynchronous by default... there is a small delay between committing a transaction in the primary and the changes becoming visible in the standby." Lag monitoring: "You can calculate this lag by comparing the current WAL write location on the primary with the last WAL location received by the standby. These locations can be retrieved using `pg_current_wal_lsn` on the primary and `pg_last_wal_receive_lsn` on the standby."

**Source**: PostgreSQL 18 Documentation, §26.2.5, §26.2.5.2
**URL**: https://www.postgresql.org/docs/current/warm-standby.html
**Published**: 2026 (PostgreSQL 18)
**Confidence**: HIGH
**Corroborated By**: Source 7 (Vitess docs also states MySQL default replication is asynchronous); Source 3 (Azure docs state replicas "are updated asynchronously")

**Notes**: Multiple independent vendors converge on async as default. PostgreSQL documents the mechanism (WAL shipping after commit) and monitoring functions (LSN diff).

---

## Evidence 2: Replication Lag Ranges from Seconds to Hours in Cloud Replicas

**Claim**: Cloud read replica replication lag typically ranges from seconds to minutes, but can extend to hours in heavy write or high-latency scenarios.

**Evidence**: Azure docs: "delay, which typically ranges from a few seconds to minutes, and in some heavy workload or high-latency scenarios, this delay could extend to hours." Also: "persistent heavy write-intensive primary workloads, the replication lag can continue to grow and might only be able to catch up with the primary."

**Source**: Azure Database for PostgreSQL Flexible Server — Read Replicas
**URL**: https://learn.microsoft.com/en-us/azure/postgresql/read-replica/concepts-read-replicas
**Published**: 2026-07-13
**Confidence**: HIGH
**Corroborated By**: Source 2 (AWS RDS Read Replica Lag metric "in seconds"); Source 1 (PostgreSQL docs: "typically under one second assuming the standby is powerful enough")

**Notes**: AWS RDS reports lag in seconds for healthy replicas; Azure acknowledges degraded scenarios reach hours. PostgreSQL reports streaming lag "typically under one second." Lag magnitude depends on workload and network conditions.

---

## Evidence 3: Synchronous Replication Trade-off — Latency vs Durability

**Claim**: Enabling synchronous replication in PostgreSQL forces commit to wait for standby confirmation, increasing response time by at least the round-trip time.

**Evidence**: PostgreSQL docs: "each commit of a write transaction will wait until confirmation is received that the commit has been written to the write-ahead log on disk of both the primary and standby server." Also: "The minimum wait time is the round-trip time between primary and standby."

**Source**: PostgreSQL 18 Documentation, §26.2.8
**URL**: https://www.postgresql.org/docs/current/warm-standby.html
**Published**: 2026 (PostgreSQL 18)
**Confidence**: HIGH
**Corroborated By**: Source 7 (Vitess docs on semi-sync: "semi-sync guarantees that at least one replica has the transaction in its relay log, but it has not necessarily been applied yet")

**Notes**: PostgreSQL provides per-transaction `synchronous_commit` control (local, off, remote_write, remote_apply) for graduated consistency vs performance.

---

## Evidence 4: MongoDB Causal Consistency Sessions Provide Read-Your-Writes

**Claim**: MongoDB causally consistent client sessions guarantee read-your-writes, monotonic reads, monotonic writes, and writes-follow-reads across replica set members.

**Evidence**: MongoDB docs list causal consistency guarantees: "Read your writes: Read operations reflect the results of write operations that precede them." "Monotonic reads: Read operations do not return results that correspond to an earlier state of the data than a preceding read operation." These hold across "all members of the MongoDB deployment."

**Source**: MongoDB — Read Isolation, Consistency, and Recency
**URL**: https://www.mongodb.com/docs/manual/core/read-isolation-consistency-recency/
**Published**: Current (accessed 2026)
**Confidence**: HIGH
**Corroborated By**: Source 10 (Terry et al. — session guarantees for weakly consistent replicated data)

**Notes**: MongoDB implements session guarantees via `clusterTime`/`operationTime` tracking; requires `"majority"` read/write concern.

---

## Evidence 5: Application-Level Read/Write Splitting via Middleware

**Claim**: ORM/proxy middleware can automatically route reads to replicas and writes to primary based on SQL statement type.

**Evidence**: GORM DBResolver: "For `Query`, `Row` callback, will use `replicas` unless `Write` mode specified... For `Raw` callback, statements are considered read-only and will use `replicas` if the SQL starts with `SELECT`." Manual override: `db.Clauses(dbresolver.Write)` forces primary.

**Source**: GORM DBResolver
**URL**: https://gorm.io/docs/dbresolver.html
**Published**: 2026-08-04
**Confidence**: MEDIUM
**Corroborated By**: Source 7 (Vitess — query routing at proxy level)

**Notes**: Read/write split is implemented in application layer (GORM) or proxy layer (Vitess, PgBouncer with routing rules). This enables primary routing for critical reads.

---

## Evidence 6: Semi-Sync Does Not Guarantee Up-to-Date Reads

**Claim**: MySQL semi-synchronous replication ensures a transaction is in at least one replica's relay log but not necessarily applied; only primary reads guarantee freshness.

**Evidence**: Vitess docs: "semi-sync guarantees that at least one replica has the transaction in its relay log, but it has not necessarily been applied yet. The only way Vitess guarantees a fully up-to-date read is to send the request to the primary."

**Source**: Vitess Docs — Replication
**URL**: https://vitess.io/docs/24.0/reference/features/mysql-replication/
**Published**: 2025-10-29
**Confidence**: HIGH
**Corroborated By**: Source 1 (PostgreSQL `synchronous_commit=remote_apply` needed for visibility); Source 3 (Azure: "aren't intended for synchronous replication scenarios")

**Notes**: Distinguishes "durability" (WAL written) from "visibility" (replay applied). `remote_apply` / read-after-write on replica requires apply, not just receive.

---

## Evidence 7: Aurora Read Replicas Share Same Storage Volume

**Claim**: Aurora Replicas connect to the same distributed storage volume as the primary, enabling read scaling with minimal replication lag.

**Evidence**: AWS docs: "Aurora Replica (reader DB instance) – Connects to the same storage volume as the primary DB instance but supports only read operations. Each Aurora DB cluster can have up to 15 Aurora Replicas."

**Source**: Amazon Aurora — DB Clusters Overview
**URL**: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/Aurora.Overview.html
**Published**: Current (accessed 2026)
**Confidence**: HIGH
**Corroborated By**: Source 2 (AWS RDS read replicas for non-Aurora engines use async replication with lag metrics)

**Notes**: Aurora's storage-sharing architecture reduces replication lag vs. classic log-shipping but compute/apply lag still exists for read-after-write.

---

## Evidence 8: Secondary Reads Are Eventually Consistent by Default

**Claim**: Reads from read-only secondaries/replicas may be stale due to asynchronous replication, reflecting eventual consistency.

**Evidence**: MongoDB docs: "reads directed to secondaries through a non-primary read preference may return data that lags behind the primary, since secondaries replicate asynchronously... This reflects the eventual consistency model MongoDB uses." PostgreSQL hot standby docs: "data on the standby is eventually consistent with the primary."

**Source**: MongoDB docs; PostgreSQL hot-standby.html
**URLs**:
- https://www.mongodb.com/docs/manual/core/read-isolation-consistency-recency/
- https://www.postgresql.org/docs/current/hot-standby.html
**Published**: 2026
**Confidence**: HIGH
**Corroborated By**: Source 9 (Kleppmann: async replication + follower reads = non-linearizable)

**Notes**: All async-replicated databases accept eventual consistency for replica reads; strong consistency requires primary or synchronous apply.

---

## Evidence 9: Critical Reads Should Route to Primary

**Claim**: Applications requiring up-to-date reads must route queries to the primary/writer, as replicas cannot guarantee freshness under async replication.

**Evidence**: Vitess docs: "The only way Vitess guarantees a fully up-to-date read is to send the request to the primary." AWS RDS: "the data on the read replica might be stale because the source DB instance is unavailable." Azure: "use this feature for workloads that can accommodate this delay."

**Sources**: Vitess, AWS RDS, Azure PostgreSQL
**URLs**:
- https://vitess.io/docs/24.0/reference/features/mysql-replication/
- https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_ReadRepl.html
- https://learn.microsoft.com/en-us/azure/postgresql/read-replica/concepts-read-replicas
**Confidence**: HIGH
**Corroborated By**: All three sources converge on: replicas = eventual, primary = fresh.

**Notes**: Routing logic (sticky session, LSN/timestamp check, read-after-write token) must be implemented at application/proxy layer.

---

## Evidence 10: Replication Lag Monitoring via LSN/Bytes/Seconds

**Claim**: Databases expose replication lag via LSN differences (PostgreSQL), bytes (Azure `physical_replication_delay_in_bytes`), or seconds (Azure `physical_replication_delay_in_seconds`, AWS Read Replica Lag).

**Evidence**: PostgreSQL: `pg_current_wal_lsn` vs `pg_last_wal_receive_lsn` → `pg_wal_lsn_diff()`. Azure metrics: `Max Physical Replication Lag` (bytes), `Read Replica Lag` (seconds). AWS RDS: "Read Replica Lag" metric.

**Sources**:
- PostgreSQL: https://www.postgresql.org/docs/current/functions-admin.html
- Azure: https://learn.microsoft.com/en-us/azure/postgresql/read-replica/concepts-read-replicas
- AWS: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_ReadRepl.html
**Confidence**: HIGH
**Corroborated By**: Multiple vendors provide native lag monitoring.

**Notes**: LSN-based lag is precise (byte-level); seconds-based lag is easier for alerting but less precise.

---

## Evidence 11: Linearizability Violation with Lagging Replica

**Claim**: A read from a lagging replica after a completed write violates linearizability (read-your-writes not guaranteed without coordination).

**Evidence**: Kleppmann: "If operation B started after operation A successfully completed, then operation B must see the system in the same state as it was on completion of operation A, or a newer state." Also: "If you allow the application to make reads from a follower, and the replication is asynchronous... then a follower may be a little behind the leader when you read from it. In this case, your reads will not be linearizable."

**Source**: Martin Kleppmann — Please stop calling databases CP or AP
**URL**: https://martin.kleppmann.com/2015/05/11/please-stop-calling-databases-cp-or-ap.html
**Published**: 2015-05-11
**Confidence**: HIGH
**Corroborated By**: Source 10 (Terry et al. session guarantees); Source 1 (PostgreSQL eventually consistent standby)

**Notes**: Read-your-writes is a weaker session guarantee than full linearizability; achievable via session tokens/LSN tracking without global coordination.
