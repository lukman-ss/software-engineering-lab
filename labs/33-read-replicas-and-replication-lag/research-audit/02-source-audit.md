# Source Audit

## Source 1

**Claimed Title**: PostgreSQL 18 Documentation: Log-Shipping Standby Servers (Streaming Replication) & Hot Standby  
**Claimed Publisher**: PostgreSQL Global Development Group  
**URL**: https://www.postgresql.org/docs/current/warm-standby.html  

**Reachable**: YES  
**Source Type**: PRIMARY  
**Relevant**: YES  
**Supports Claimed Topic**: YES  

**Problems**:
- None. URL is live and content directly supports streaming replication async default, LSN monitoring (`pg_current_wal_lsn`, `pg_last_wal_receive_lsn`), and `synchronous_commit` parameters (including `remote_apply`).

**Assessment**: PASS

---

## Source 2

**Claimed Title**: Working with DB Instance Read Replicas  
**Claimed Publisher**: Amazon Web Services  
**URL**: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_ReadRepl.html  

**Reachable**: YES  
**Source Type**: PRIMARY  
**Relevant**: YES  
**Supports Claimed Topic**: YES  

**Problems**:
- None. Standard official AWS RDS documentation for MySQL/PostgreSQL/MariaDB read replicas.

**Assessment**: PASS

---

## Source 3

**Claimed Title**: Read replicas in Azure Database for PostgreSQL Flexible Server  
**Claimed Publisher**: Microsoft Learn  
**URL**: https://learn.microsoft.com/en-us/azure/postgresql/read-replica/concepts-read-replicas  

**Reachable**: YES  
**Source Type**: PRIMARY  
**Relevant**: YES  
**Supports Claimed Topic**: YES  

**Problems**:
- None. Accurately cited for lag metrics (`physical_replication_delay_in_seconds`, `physical_replication_delay_in_bytes`) and lag duration warnings (seconds to minutes, hours under heavy load).

**Assessment**: PASS

---

## Source 4

**Claimed Title**: Read Isolation, Consistency, and Recency  
**Claimed Publisher**: MongoDB Inc.  
**URL**: https://www.mongodb.com/docs/manual/core/read-isolation-consistency-recency/  

**Reachable**: YES  
**Source Type**: PRIMARY  
**Relevant**: YES  
**Supports Claimed Topic**: YES  

**Problems**:
- None. Officially details MongoDB causal consistency sessions (`read-your-writes`, `monotonic reads`, `monotonic writes`, `writes-follow-reads`).

**Assessment**: PASS

---

## Source 5

**Claimed Title**: Amazon Aurora DB Clusters  
**Claimed Publisher**: Amazon Web Services  
**URL**: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/Aurora.Overview.html  

**Reachable**: YES  
**Source Type**: PRIMARY  
**Relevant**: YES  
**Supports Claimed Topic**: YES  

**Problems**:
- None. Documents Aurora's shared storage volume architecture with up to 15 reader replicas.

**Assessment**: PASS

---

## Source 6

**Claimed Title**: PostgreSQL 18: 19.6. Replication (Configuration Parameters)  
**Claimed Publisher**: PostgreSQL Global Development Group  
**URL**: https://www.postgresql.org/docs/current/runtime-config-replication.html  

**Reachable**: YES  
**Source Type**: PRIMARY  
**Relevant**: YES  
**Supports Claimed Topic**: YES  

**Problems**:
- None. Documents GUC parameters `synchronous_standby_names`, `synchronous_commit`, `hot_standby_feedback`, `max_standby_streaming_delay`.

**Assessment**: PASS

---

## Source 7

**Claimed Title**: Vitess Docs: Replication  
**Claimed Publisher**: Vitess (CNCF / Open Source)  
**URL**: https://vitess.io/docs/24.0/reference/features/mysql-replication/  

**Reachable**: YES  
**Source Type**: PRIMARY  
**Relevant**: YES  
**Supports Claimed Topic**: YES  

**Problems**:
- None. Correctly cited for semi-synchronous MySQL replication semantics (relay log receipt != applied) and primary query routing.

**Assessment**: PASS

---

## Source 8

**Claimed Title**: DBResolver (GORM)  
**Claimed Publisher**: GORM (Go ORM)  
**URL**: https://gorm.io/docs/dbresolver.html  

**Reachable**: YES  
**Source Type**: SECONDARY  
**Relevant**: YES  
**Supports Claimed Topic**: YES  

**Problems**:
- None. Demonstrates application ORM read/write splitting and manual `dbresolver.Write` overrides.

**Assessment**: PASS

---

## Source 9

**Claimed Title**: Please stop calling databases CP or AP  
**Claimed Publisher**: Martin Kleppmann (Cambridge University)  
**URL**: https://martin.kleppmann.com/2015/05/11/please-stop-calling-databases-cp-or-ap.html  

**Reachable**: YES  
**Source Type**: SECONDARY  
**Relevant**: YES  
**Supports Claimed Topic**: YES  

**Problems**:
- None. High-relevance foundational blog post by prominent distributed systems researcher explaining linearizability violations on asynchronous followers.

**Assessment**: PASS

---

## Source 10

**Claimed Title**: Replicated Data Consistency Explained Through Baseball (MSR-TR-2011-137)  
**Claimed Publisher**: Microsoft Research (Douglas B. Terry et al.)  
**URL**: https://www.microsoft.com/en-us/research/wp-content/uploads/2011/10/ConsistencyAndBaseballReport.pdf  

**Reachable**: YES  
**Source Type**: PRIMARY  
**Relevant**: YES  
**Supports Claimed Topic**: YES  

**Problems**:
- None. Canonical research paper defining session guarantees (read-your-writes, monotonic reads, monotonic writes, writes-follow-reads).

**Assessment**: PASS
