# Claim Audit: Optimistic vs Pessimistic Locking

## Claim 1
Claim: Two concurrent transactions performing uncoordinated read-modify-write under READ COMMITTED isolation can cause the first transaction's write to be silently lost.  
Location: `research/05-report.md:21-34` (Finding 1); `research/03-evidence.md:7-40` (Evidence 01, 02)  
Evidence Provided: Wikipedia Concurrency Control, Oracle 19c Concepts (Table 10-2 Banda salary case), PostgreSQL 18 Docs 13.2.  
Source: Sources 2, 3, 10  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Classic database anomaly with direct documentation in academic and vendor materials.  

---

## Claim 2
Claim: `SELECT ... FOR UPDATE` acquires row-level exclusive locks held until transaction commit/rollback, blocking concurrent modifications and locking reads while allowing plain `SELECT`.  
Location: `research/05-report.md:35-47` (Finding 2); `research/03-evidence.md:43-76` (Evidence 03, 04)  
Evidence Provided: PostgreSQL 13.3.2, MySQL InnoDB Locking Reads 15.7.2.4, Oracle Concepts 9.  
Source: Sources 1, 7, 10  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Verified across PostgreSQL, MySQL, and Oracle documentation.  

---

## Claim 3
Claim: Pessimistic locking reduces system concurrency and introduces deadlock risks; transactions holding locks across slow/external calls (e.g., payment gateways) degrade throughput and heighten lock wait times.  
Location: `research/05-report.md:48-59` (Finding 3); `research/03-evidence.md:79-113` (Evidence 05, 06)  
Evidence Provided: PostgreSQL 13.3.4 (Deadlocks & long-held locks warning), Fowler Pessimistic Offline Lock.  
Source: Sources 1, 3, 5  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Vendor docs explicitly advise against keeping open transactions waiting on external inputs.  

---

## Claim 4
Claim: Optimistic locking validates at commit time (e.g., `UPDATE ... WHERE id = ? AND version = ?`) and requires explicit handling of zero affected rows (reload/recalculate/retry or 409 Conflict).  
Location: `research/05-report.md:60-74` (Finding 4); `research/03-evidence.md:115-132, 278-290` (Evidence 07, 16)  
Evidence Provided: Fowler Optimistic Offline Lock, Oracle WHERE-guard recommendation, Hibernate User Guide.  
Source: Sources 4, 10, 12, 14  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Pattern and conflict handling rules are solidly supported.  

---

## Claim 5
Claim: Pessimistic locking prevents conflicts (ideal for high-contention or financial/stock invariants), whereas optimistic locking detects conflicts (ideal for low-contention, read-heavy workflows).  
Location: `research/05-report.md:75-88` (Finding 5); `research/03-evidence.md:133-166, 260-275` (Evidence 08, 09, 15)  
Evidence Provided: Fowler PoEAA, Wikipedia Concurrency Control.  
Source: Sources 3, 4, 5, 12  
Source Actually Supports Claim: YES  
Classification: INTERPRETATION  
Severity: LOW  
Notes: Standard software architecture heuristics supported by canonical literature.  

---

## Claim 6
Claim: Atomic single-statement updates (`UPDATE ... SET stock = stock - N WHERE stock >= N` + checking `affected_rows == 1`) eliminate the read-modify-write race window without application-level or explicit row locking.  
Location: `research/05-report.md:89-101` (Finding 6); `research/03-evidence.md:169-185` (Evidence 10)  
Evidence Provided: Oracle ACID Statement-level atomicity, PostgreSQL/MySQL statement execution semantics.  
Source: Sources 1, 9, 10, 11  
Source Actually Supports Claim: PARTIAL  
Classification: IMPLEMENTATION-SPECIFIC  
Severity: MEDIUM  
Notes: The SQL statement atomicity guarantee is verified; however, the exact vendor recommendation and dialect nuance for composite business checks are synthesized rather than quoted directly from vendor docs. The research honestly flagged this as MEDIUM confidence.  

---

## Claim 7
Claim: Wrapping code in a plain database transaction does not automatically prevent lost updates without explicit query locking or isolation level escalation.  
Location: `research/05-report.md:102-115` (Finding 7); `research/03-evidence.md:205-222` (Evidence 12)  
Evidence Provided: PostgreSQL 13.2 / 13.3.2, Oracle lost-update example Table 10-2, MySQL semi-consistent reads.  
Source: Sources 1, 2, 9, 10  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Common developer misconception refuted by official database isolation docs.  

---

## Claim 8
Claim: Default isolation levels differ among database engines (PostgreSQL and Oracle default to READ COMMITTED; MySQL InnoDB defaults to REPEATABLE READ; PostgreSQL treats READ UNCOMMITTED as READ COMMITTED).  
Location: `research/05-report.md:116-128` (Finding 8); `research/03-evidence.md:187-203` (Evidence 11)  
Evidence Provided: PostgreSQL 13.2 docs, MySQL 15.7.2.1 docs, Oracle Concepts 9.  
Source: Sources 2, 9, 10  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Confirmed against primary vendor reference manuals.  
