# Claim Audit

## Claim 1: Asynchronous Replication Default

**Claim**: Streaming replication in PostgreSQL and default replication in MySQL/MongoDB are asynchronous, meaning commits complete locally before changes are confirmed on replicas, creating a window for stale reads.  
**Location**: `03-evidence.md` (Evidence 1), `05-report.md` (Finding 1)  
**Evidence Provided**: PostgreSQL docs §26.2.5, Vitess documentation, Azure PG documentation.  
**Source**: Sources 1, 3, 7  
**Source Actually Supports Claim**: YES  
**Classification**: FACT  
**Severity**: LOW  
**Notes**: Supported across database engines and cloud offerings.

---

## Claim 2: Replication Lag Magnitude

**Claim**: Replication lag can range from sub-second under normal conditions to seconds, minutes, or hours under heavy write load or cross-region network delays.  
**Location**: `03-evidence.md` (Evidence 2), `05-report.md` (Finding 1)  
**Evidence Provided**: Azure docs specifically mention seconds to minutes/hours; AWS RDS lag metric in seconds; PostgreSQL docs note "typically under one second".  
**Source**: Sources 1, 2, 3  
**Source Actually Supports Claim**: YES  
**Classification**: FACT  
**Severity**: LOW  
**Notes**: Demonstrates nuance between optimal local network streaming and overloaded cloud replica behavior.

---

## Claim 3: Synchronous Durability vs Visibility

**Claim**: Synchronous replication modes (e.g. MySQL semi-sync, PostgreSQL `remote_write`) guarantee that replicas have received the WAL/relay log, but do not guarantee visibility/read freshness unless applied (`remote_apply`).  
**Location**: `03-evidence.md` (Evidence 3, Evidence 6), `05-report.md` (Executive Summary 2, Finding 5)  
**Evidence Provided**: PostgreSQL documentation §26.2.8, Vitess documentation.  
**Source**: Sources 1, 6, 7  
**Source Actually Supports Claim**: YES  
**Classification**: FACT  
**Severity**: LOW  
**Notes**: Critical distinction accurately captured. Confirms semi-sync is not sufficient for read-your-writes.

---

## Claim 4: Follower Reads Violate Linearizability

**Claim**: Reading from an asynchronously updated replica without synchronization violates linearizability, as clients may observe stale data after a successful write.  
**Location**: `03-evidence.md` (Evidence 11), `05-report.md` (Finding 2)  
**Evidence Provided**: Martin Kleppmann (2015) analysis of linearizability and follower reads.  
**Source**: Source 9  
**Source Actually Supports Claim**: YES  
**Classification**: FACT  
**Severity**: LOW  
**Notes**: Accurate distributed systems theoretical framing.

---

## Claim 5: Native Causal Consistency in MongoDB Sessions

**Claim**: MongoDB provides native causal consistency in client sessions, guaranteeing read-your-writes, monotonic reads, monotonic writes, and writes-follow-reads via cluster time tracking.  
**Location**: `03-evidence.md` (Evidence 4), `05-report.md` (Executive Summary 3)  
**Evidence Provided**: MongoDB Read Isolation, Consistency, and Recency manual.  
**Source**: Source 4  
**Source Actually Supports Claim**: YES  
**Classification**: FACT  
**Severity**: LOW  
**Notes**: Correctly identified as an engine-native feature in MongoDB that is absent from relational databases like PostgreSQL.

---

## Claim 6: Shared Storage Architecture in Aurora

**Claim**: Amazon Aurora reader replicas connect to the same distributed shared storage volume as the primary, eliminating classic log shipping at the storage layer while supporting up to 15 read replicas.  
**Location**: `03-evidence.md` (Evidence 7)  
**Evidence Provided**: AWS Aurora DB Clusters overview documentation.  
**Source**: Source 5  
**Source Actually Supports Claim**: YES  
**Classification**: FACT  
**Severity**: LOW  
**Notes**: Nuance preserved in `06-open-questions.md` regarding in-memory buffer cache invalidation lag.

---

## Claim 7: Application-Level Routing via Middleware

**Claim**: Middleware and ORMs like GORM DBResolver can automatically route SELECT statements to read replicas while routing mutations to the primary, with manual override support.  
**Location**: `03-evidence.md` (Evidence 5), `05-report.md` (Finding 4)  
**Evidence Provided**: GORM DBResolver documentation and Vitess proxy docs.  
**Source**: Sources 7, 8  
**Source Actually Supports Claim**: YES  
**Classification**: IMPLEMENTATION-SPECIFIC  
**Severity**: LOW  
**Notes**: Accurately presented as an application/middleware pattern rather than a universal database feature.

---

## Claim 8: Fixed 5-Second Sticky Routing Window

**Claim**: Subsequent reads within 5 seconds after a write should be routed to the primary or stickied before falling back to read replicas.  
**Location**: `04-contradictions.md` (§Uncertainty Areas), `05-report.md` (Conclusion), `06-open-questions.md` (Question 1)  
**Evidence Provided**: Lab specification and tutorial conventions; acknowledged that no official database documentation prescribes 5 seconds.  
**Source**: Uncited industry rule of thumb / lab requirement.  
**Source Actually Supports Claim**: PARTIAL  
**Classification**: HEURISTIC  
**Severity**: MEDIUM  
**Notes**: The research correctly flags this in `04-contradictions.md` and `06-open-questions.md` as an arbitrary heuristic that depends on actual replication lag in production. Because the report candidly identifies it as a heuristic rather than universal fact, it does not deceive readers.
