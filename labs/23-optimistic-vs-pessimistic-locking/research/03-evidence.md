# Evidence: Optimistic vs Pessimistic Locking

Research date: 2026-09-26

---

## Evidence 01

Claim: The "lost update" problem occurs when a second transaction overwrites a first transaction's update without seeing it, causing the first value to be lost.

Evidence: "The lost update problem: A second transaction writes a second value of a data-item (datum) on top of a first value written by a first concurrent transaction, and the first value is lost to other transactions running concurrently which need, by their precedence, to read the first value."

Source: Wikipedia - Concurrency Control (citing Bernstein et al. 1987; Weikum and Vossen 2001)

URL: https://en.wikipedia.org/wiki/Concurrency_control

Confidence: HIGH

Corroborated By: Source 10 (Oracle 19c Concepts - Data Concurrency) documents identical pattern in Table 10-2 (Banda salary 7000 overwritten to 6300); Source 2 (PostgreSQL Docs 13.2) describes lost-update risk under Read Committed

Notes: Foundational definition. Stock example (10 -> 7 vs 6 -> final 6 instead of correct 3) from topic specification is a direct instantiation of this definition. Direct fetch of Wikipedia article confirms verbatim definition.

---

## Evidence 02

Claim: Under READ COMMITTED isolation, the classic lost-update scenario is reproducible: Session 1 updates a row, Session 2 reads the stale pre-commit value, Session 1 commits, Session 2 overwrites with stale-based value.

Evidence: Oracle documents: Session 1 updates Banda salary to 7000 (holds TX row lock). Session 2 reads salary as 6200 (read-consistent snapshot before commit). Session 1 commits. Session 2 updates to 6300 using old value. Result: Session 1's update to 7000 is lost. Recommended prevention: `UPDATE employees SET salary = 7000 WHERE last_name = 'Banda' AND salary = 6200` (include original value in WHERE clause). This appears in Table 10-2 "Conflicting Writes and Lost Updates in a READ COMMITTED Transaction".

Source: Oracle Database Concepts 19c - 9 Data Concurrency and Consistency

URL: https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html

Confidence: HIGH

Corroborated By: Source 2 (PostgreSQL transaction-iso docs describe Read Committed re-evaluation allowing second updater to overwrite); Source 3 (Wikipedia lost-update definition)

Notes: Oracle's example uses salary; topic spec uses stock. Same anomaly class. Direct fetch verified Table 10-2 content.

---

## Evidence 03

Claim: Pessimistic locking via SELECT ... FOR UPDATE causes retrieved rows to be locked as though for update; other transactions attempting UPDATE/DELETE/SELECT FOR UPDATE on same rows are blocked until the holding transaction ends.

Evidence: "FOR UPDATE causes the rows retrieved by the SELECT statement to be locked as though for update. This prevents them from being locked, modified or deleted by other transactions until the current transaction ends. That is, other transactions that attempt UPDATE, DELETE, SELECT FOR UPDATE, SELECT FOR NO KEY UPDATE, SELECT FOR SHARE or SELECT FOR KEY SHARE of these rows will be blocked until the current transaction ends"

Source: PostgreSQL 18 Documentation - 13.3. Explicit Locking, Section 13.3.2 Row-Level Locks

URL: https://www.postgresql.org/docs/current/explicit-locking.html

Confidence: HIGH (Direct fetch verified, 2026-09-26)

Corroborated By: Source 7 (MySQL InnoDB Locking Reads: SELECT ... FOR UPDATE locks rows and index entries, blocked until commit/rollback); Source 10 (Oracle: TX row lock acquired by SELECT ... FOR UPDATE, held until commit/rollback)

Notes: Directly opened and verified. Covers FOR UPDATE; doc also defines FOR NO KEY UPDATE, FOR SHARE, FOR KEY SHARE with weaker blocking semantics and compatibility matrix Table 13.3.

---

## Evidence 04

Claim: MySQL/InnoDB supports pessimistic locking via SELECT ... FOR SHARE (shared lock) and SELECT ... FOR UPDATE (exclusive lock), released only at transaction commit/rollback, with non-blocking variants NOWAIT and SKIP LOCKED.

Evidence: "InnoDB supports two types of locking reads: SELECT ... FOR SHARE sets a shared mode lock; SELECT ... FOR UPDATE locks the rows and associated index entries, same as UPDATE. All locks are released when the transaction is committed or rolled back. NOWAIT and SKIP LOCKED options are also available."

Source: MySQL 8.0 Reference Manual - 15.7.2.4 Locking Reads

URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html (mirror: https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-locking-reads.html - direct MySQL domain returned 403 on fetch but content confirmed via Oracle CDN mirror and subagent extraction)

Confidence: MEDIUM (URL content confirmed via alternate official Oracle CDN mirror hosting identical MySQL documentation; direct fetch via Oracle CDN mirror 2026-09-26)

Corroborated By: Source 1 (PostgreSQL equivalent FOR UPDATE semantics); Source 10 (Oracle TX row lock equivalent)

Notes: Direct fetch via Oracle CDN mirror confirmed verbatim MySQL locking reads documentation. InnoDB-specific gap/next-key lock behavior not present in PostgreSQL; see Evidence 10.

---

## Evidence 05

Claim: Pessimistic locking introduces deadlock risk; database systems auto-detect deadlocks and abort one transaction, requiring application retry logic. Consistent lock ordering mitigates risk.

Evidence: PostgreSQL docs: "PostgreSQL automatically detects deadlock situations and resolves them by aborting one of the transactions involved, allowing the other(s) to complete. (Exactly which transaction will be aborted is difficult to predict and should not be relied upon.)" + example of two transactions updating accounts 11111/22222 in opposite order causing deadlock. "The best defense against deadlocks is generally to avoid them by being certain that all applications using a database acquire locks on multiple objects in a consistent order."

Source: PostgreSQL 18 Documentation - 13.3. Explicit Locking, Section 13.3.4 Deadlocks

URL: https://www.postgresql.org/docs/current/explicit-locking.html

Confidence: HIGH

Corroborated By: Source 3 (Wikipedia: "Most non-optimistic mechanisms (with blocking) are prone to deadlocks which are resolved by an intentional abort"); Source 13 (DynamoDB transaction conflict handling as analogous distributed case)

Notes: Directly verified. Row-level deadlocks can occur even without explicit locking (two concurrent UPDATEs in opposite row order).

---

## Evidence 06

Claim: Holding a pessimistic-lock transaction open for a long time (e.g., waiting for user input or external HTTP call) is harmful: conflicting transactions wait indefinitely and throughput drops.

Evidence: "So long as no deadlock situation is detected, a transaction seeking either a table-level or row-level lock will wait indefinitely for conflicting locks to be released. This means it is a bad idea for applications to hold transactions open for long periods of time (e.g., while waiting for user input)."

Source: PostgreSQL 18 Documentation - 13.3. Explicit Locking, Section 13.3.4 Deadlocks

URL: https://www.postgresql.org/docs/current/explicit-locking.html

Confidence: HIGH

Corroborated By: Topic specification's anti-pattern example (BEGIN -> Lock -> Call Payment Gateway 15s -> Update -> COMMIT) aligns with doc warning; Source 5 (Pessimistic Offline Lock) notes reduced concurrency as core cost

Notes: Directly supports topic spec claim "Semakin lama memegang lock: Concurrency turun, Wait time naik, Deadlock risk naik" and "Transaction harus pendek. Jangan lakukan Call Payment Gateway di dalam transaction yang memegang lock."

---

## Evidence 07

Claim: Optimistic locking prevents conflicts by validating at commit time that data has not changed since it was read (e.g., version check in WHERE clause); if no row is affected (0 rows), a conflict is detected.

Evidence: "Optimistic Offline Lock solves this problem by validating that the changes about to be committed by one session don't conflict with the changes of another session. A successful pre-commit validation is, in a sense, obtaining a lock indicating it's okay to go ahead with the changes to the record data. So long as the validation and the updates occur within a single system transaction the business transaction will display consistency."

Source: Martin Fowler - Optimistic Offline Lock (Patterns of Enterprise Application Architecture)

URL: https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html

Confidence: HIGH

Corroborated By: Source 3 (Wikipedia optimistic category: "only check for violations at each transaction's commit; if violations detected, transaction is aborted and restarted"); Source 8/Hibernate docs: version check increment or OptimisticLockException; Source 10 (Oracle WHERE salary=6200 original-value pattern = optimistic check); Baeldung (Source 12): "Before the transaction wants to make an update, it checks the version property again. If changed, OptimisticLockException is thrown."

Notes: Topic spec pattern `UPDATE products SET stock=7, version=6 WHERE id=10 AND version=5` with 0-rows-affected detection is exactly this mechanism. Page fetched and verified.

---

## Evidence 08

Claim: Optimistic locking assumes conflict is rare and allows concurrent work; pessimistic locking assumes conflict is frequent and serializes access by acquiring locks early.

Evidence (optimistic): "Whereas Pessimistic Offline Lock assumes that the chance of session conflict is high and therefore limits the system's concurrency, Optimistic Offline Lock assumes that the chance of conflict is low. The expectation that session conflict isn't likely allows multiple users to work with the same data at the same time." Evidence (pessimistic): "Pessimistic Offline Lock prevents conflicts between concurrent business transactions by allowing only one business transaction at a time to access data." and "The first approach to try is Optimistic Offline Lock. However, if several people access the same data within a business transaction, one commits easily but others fail... If this happens a lot on lengthy business transactions the system will soon become very unpopular. Pessimistic Offline Lock prevents conflicts by avoiding them altogether."

Source: Martin Fowler - Optimistic Offline Lock + Pessimistic Offline Lock

URL: https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html , https://martinfowler.com/eaaCatalog/pessimisticOfflineLock.html

Confidence: HIGH

Corroborated By: Source 3 (Wikipedia: "This approach [optimistic] is very efficient when few transactions are aborted" vs pessimistic "Blocking operations is typically involved with performance reduction"); Baeldung Source 12 (optimistic suitable when many more reads than writes)

Notes: Directly supports topic spec "Pessimistic = Prevent conflict, Optimistic = Detect conflict" and use-case split (optimistic for Edit Profile/CRM/CMS/master data; pessimistic for wallet/balance/stock/booking).

---

## Evidence 09

Claim: Optimistic and pessimistic are the two primary categories of concurrency control; there is also a semi-optimistic category that mixes them.

Evidence: "The main categories of concurrency control mechanisms are: Optimistic - Allow transactions to proceed without blocking any of their (read, write) operations ... and only check for violations at each transaction's commit. ... Pessimistic - Block an operation of a transaction, if it may cause violation of the rules, until the possibility of violation disappears. ... Semi-optimistic - Responds pessimistically or optimistically depending on the type of violation and how quickly it can be detected."

Source: Wikipedia - Concurrency Control, Section Categories

URL: https://en.wikipedia.org/wiki/Concurrency_control

Confidence: HIGH

Corroborated By: Source 4 and 5 (Fowler's offline-lock pair maps directly to this taxonomy)

Notes: Academic framing. Topic spec's binary pessimistic/optimistic split is a simplification; semi-optimistic exists but not needed for this lab's scope.

---

## Evidence 10

Claim: Atomic database operations (e.g., UPDATE SET stock = stock - 3 WHERE id=? AND stock >= 3 and checking affected_rows) eliminate the read-modify-write race window without explicit locking.

Evidence: Direct verification from PostgreSQL 13.4.2 (Enforcing Consistency with Explicit Blocking Locks): "When non-serializable writes are possible, to ensure the current validity of a row and protect it against concurrent updates one must use SELECT FOR UPDATE, SELECT FOR SHARE, or an appropriate LOCK TABLE statement." Critically, PostgreSQL adds: "SELECT FOR UPDATE does not ensure that a concurrent transaction will not update or delete a selected row. To do that in PostgreSQL you must actually update the row, even if no values need to be changed. SELECT FOR UPDATE temporarily blocks other transactions from acquiring the same lock or executing an UPDATE or DELETE which would affect the locked row, but once the transaction holding this lock commits or rolls back, a blocked transaction will proceed with the conflicting operation unless an actual UPDATE of the row was performed while the lock was held."

This confirms: (a) an actual UPDATE is the authoritative conflict-resolution action in PostgreSQL; (b) a single conditional UPDATE (`SET stock = stock - 3 WHERE stock >= 3`) both performs the modification and is atomic at statement level, removing the read-modify-write window entirely. Supporting component evidence: (e) MySQL: under READ COMMITTED "InnoDB holds locks only for rows that it updates or deletes" — single conditional UPDATE holds lock only for matched row. (f) Oracle: WHERE-clause guard pattern (AND salary=6200) is the documented prevention. (g) ACID atomicity: "A SQL statement is an atomic unit; on failure, only its effects are rolled back."

Source: PostgreSQL 18 Documentation - 13.4 Data Consistency Checks at the Application Level (13.4.2); MySQL innodb-transaction-isolation-levels; Oracle data-concurrency + transactions ACID section

URL: https://www.postgresql.org/docs/current/applevel-consistency.html ; https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-transaction-isolation-levels.html ; https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html ; https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/transactions.html

Confidence: HIGH (upgraded from MEDIUM, 2026-09-26 — PostgreSQL 13.4.2 directly fetched and confirms actual-UPDATE requirement; component guarantees verified across MySQL, Oracle)

Corroborated By: Topic specification's own claim aligns with standard SQL single-statement atomicity; MySQL counter example (SELECT ... FOR UPDATE then UPDATE counter = counter + 1) in 17.7.2.4 documents the pessimistic variant of the same problem the atomic form solves.

Notes: The exact `SET stock = stock - N WHERE stock >= N` recipe is not quoted verbatim in a single vendor doc, but the mechanism (statement atomicity + conditional WHERE + affected_rows check) is fully documented; PostgreSQL 13.4.2 directly states that an actual UPDATE, not merely a lock, is what prevents concurrent modification. Upgrade justified per evidence rules (authoritative source + corroborating sources).

---

## Evidence 11

Claim: Default transaction isolation levels differ across databases and affect whether lost updates can occur without explicit locking.

Evidence: PostgreSQL docs: default is Read Committed (and Read Uncommitted behaves identically to Read Committed due to MVCC). MySQL docs: default is Repeatable Read (unlike others) and Serializable implicitly converts plain SELECT to SELECT FOR SHARE. Oracle docs: default is Read Committed; Oracle does not implement Read Uncommitted or Repeatable Read as separate levels (only Read Committed, Serializable, Read Only).

Source: PostgreSQL 18 Docs 13.2 Transaction Isolation; MySQL 8.0 Docs 15.7.2.1 Transaction Isolation Levels; Oracle 19c Concepts 9 Data Concurrency

URL: https://www.postgresql.org/docs/current/transaction-iso.html ; https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-transaction-isolation-levels.html ; https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html

Confidence: HIGH (for PostgreSQL and Oracle - directly verified; MEDIUM for MySQL via mirror)

Corroborated By: Source 3 (Wikipedia notes isolation-level anomalies critique via Berenson et al. 1995, confirming levels are not uniformly implemented)

Notes: Directly supports topic spec statement "Hasilnya bergantung pada query dan isolation level yang digunakan" and "Menganggap database transaction otomatis mencegah semua lost update — Tidak selalu."

---

## Evidence 12

Claim: Wrapping code in a database transaction alone does not automatically prevent lost updates; it depends on the queries and isolation level used. Repeatable Read in PostgreSQL uses Snapshot Isolation and can throw serialization errors rather than silently losing updates.

Evidence: PostgreSQL: "Within a REPEATABLE READ or SERIALIZABLE transaction, however, an error will be thrown if a row to be locked has changed since the transaction started." and Read Committed re-evaluates WHERE after concurrent commit. MySQL: READ COMMITTED performs "semi-consistent read" to reduce deadlocks but may still allow overwrites if app does not account for it. Oracle: explicit lost-update table showing committed transaction's effect lost despite transaction boundaries.

Source: PostgreSQL 18 Docs 13.3.2 + 13.2; MySQL 8.0 Docs 15.7.2.1; Oracle 19c 9 Data Concurrency Table 10-2

URL: https://www.postgresql.org/docs/current/explicit-locking.html ; https://www.postgresql.org/docs/current/transaction-iso.html ; https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html

Confidence: HIGH

Corroborated By: Source 3 (Wikipedia: isolation is the goal of concurrency control, but levels relax serializability for performance)

Notes: Supports topic spec "Kesalahan Umum #1: menganggap database transaction otomatis mencegah semua lost update. Tidak selalu."

---

## Evidence 13

Claim: Common pessimistic-locking mistakes include acquiring locks in inconsistent order, not acquiring the most restrictive lock first, and exhausting shared-memory lock limits.

Evidence: PostgreSQL Section 13.3.4: deadlock example with opposite lock order; "One should also ensure that the first lock acquired on an object in a transaction is the most restrictive mode that will be needed." PostgreSQL Section 13.3.5: shared-memory limit via max_locks_per_transaction x max_connections, "Care must be taken not to exhaust this memory or the server will be unable to grant any locks at all." Advisory lock LIMIT pitfall: `SELECT pg_advisory_lock ... LIMIT 100 -- danger!`

Source: PostgreSQL 18 Documentation - 13.3 Explicit Locking

URL: https://www.postgresql.org/docs/current/explicit-locking.html

Confidence: HIGH

Corroborated By: General distributed-systems literature on lock ordering (not separately sourced here)

Notes: Supports topic spec guidance on deadlock risk and lock-holding duration.

---

## Evidence 14

Claim: ACID atomicity guarantees that a single SQL statement is an atomic unit; on failure only its effects are rolled back, and committed transaction effects persist through crashes via non-volatile logging.

Evidence: "Atomicity - Either the effects of all or none of its operations remain ("all or nothing" semantics) when a transaction is completed" and "A SQL statement is an atomic unit; on failure, only its effects are rolled back" and "Durability - Effects of successful (committed) transactions must persist through crashes (typically by recording the transaction's effects and its commit event in a non-volatile memory)."

Source: Wikipedia Concurrency Control (ACID section) and Oracle 19c Concepts 10 Transactions

URL: https://en.wikipedia.org/wiki/Concurrency_control ; https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/transactions.html

Confidence: HIGH

Corroborated By: PostgreSQL MVCC + durability documentation (Source 2)

Notes: Supports distinction between single-statement atomicity (sufficient for decrement pattern) vs multi-statement transaction atomicity.

---

## Evidence 15

Claim: For low-conflict, read-heavy workloads or operations spanning multiple requests (offline/business transactions), optimistic locking is more suitable; pessimistic locking can cause widespread wasted work when many users contend.

Evidence: Baeldung: "This mechanism is suitable for applications that do many more reads than updates or deletes. It's also useful in situations where entities must be detached for some time and locks cannot be held." Fowler Optimistic: "The expectation that session conflict isn't likely allows multiple users to work with the same data at the same time." Fowler Pessimistic: "If several people access the same data within a business transaction, one will commit easily but the others will conflict and fail. Since the conflict is only detected at the end, the victims will do all the transaction work only to find at the last minute that the whole thing will fail."

Source: Baeldung - Optimistic Locking in JPA + Martin Fowler Optimistic/Pessimistic Offline Lock

URL: https://www.baeldung.com/jpa-optimistic-locking ; https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html ; https://martinfowler.com/eaaCatalog/pessimisticOfflineLock.html

Confidence: MEDIUM (Baeldung date not verified; Fowler 2003 remains canonical)

Corroborated By: Topic spec use-case lists (optimistic: Edit Profile/CRM/CMS/Master Data/Dokumen; pessimistic: Saldo/Stock terbatas/Booking kursi)

Notes: Topic spec trade-off table aligns with sources.

---

## Evidence 16

Claim: After detecting an optimistic-lock conflict (0 rows affected / OptimisticLockException), the application must handle it explicitly: reload, recalculate, retry, or return 409 Conflict to the client.

Evidence: Fowler Optimistic: conflict detection followed by rolling back the transaction. Hibernate/JPA: "If they are not from the same version, Hibernate will throw either OptimisticEntityLockException or StaleObjectStateException." Baeldung: "After that, we can retry updating the data." Topic spec example shows `0 rows affected` -> reload -> recalculate -> retry or 409.

Source: https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html ; https://docs.jboss.org/hibernate/orm/6.3/userguide/html_single/Hibernate_User_Guide.html#locking-optimistic ; https://www.baeldung.com/jpa-optimistic-locking

Confidence: HIGH

Corroborated By: PostgreSQL docs on REPEATABLE READ error requiring retry (Source 2)

Notes: Supports topic spec "Keempat: optimistic locking mendeteksi conflict tetapi aplikasi tidak menangani conflict — 0 rows affected tidak boleh dianggap sukses."

---
