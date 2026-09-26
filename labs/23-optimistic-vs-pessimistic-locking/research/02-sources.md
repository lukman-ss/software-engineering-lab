# Sources: Optimistic vs Pessimistic Locking

All sources accessed: 2026-09-26

---

## Source 1

Title: PostgreSQL 18 Documentation - 13.3. Explicit Locking
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/explicit-locking.html
Published: PostgreSQL 18.6 (current, accessed 2026-09-26)
Accessed: 2026-09-26 — Direct fetch verified (full document). Verbatim quotes extracted: FOR UPDATE blocking semantics, deadlock example (accounts 11111/22222), advisory-lock guidance, lock compatibility matrices Tables 13.2/13.3.
Source Tier: Tier 1 (official documentation)
Relevance: Primary source for pessimistic locking semantics. Documents FOR UPDATE, FOR NO KEY UPDATE, FOR SHARE, FOR KEY SHARE row-level lock modes, lock compatibility matrix, deadlocks, and advisory locks. Direct evidence for RQ2 (pessimistic locking mechanism) and RQ6 (deadlock anti-patterns).

---

## Source 2

Title: PostgreSQL 18 Documentation - 13.2. Transaction Isolation
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/transaction-iso.html
Published: PostgreSQL 18.6 (current, accessed 2026-09-26)
Accessed: 2026-09-26 — Direct fetch verified (full document). Verbatim quotes extracted: Table 13.1 isolation levels, Read Committed re-evaluation semantics, Repeatable Read/Snapshot Isolation, Serializable SSI, lost-update risk under READ COMMITTED.
Source Tier: Tier 1 (official documentation)
Relevance: Documents PostgreSQL isolation levels (Read Committed, Repeatable Read, Serializable), MVCC behavior, read-committed lost-update semantics, snapshot isolation. Direct evidence for RQ4 (isolation levels).

---

## Source 3

Title: Concurrency Control
Publisher: Wikipedia (community-maintained encyclopedia with academic citations)
URL: https://en.wikipedia.org/wiki/Concurrency_control
Published: Continuously updated (accepted revision reviewed 2026-09-07)
Accessed: 2026-09-26 — Direct fetch verified (full article). Verbatim quotes extracted: lost-update definition, optimistic/pessimistic/semi-optimistic categories, ACID properties, 2PL/SS2PL as most common mechanism. Article reviewed as latest accepted revision 2026-09-07.
Source Tier: Tier 2 (reputable technical reference, cites primary academic sources)
Relevance: Provides academic definitions of optimistic, pessimistic, and semi-optimistic categories. Defines lost update, dirty read, incorrect summary problems. Cites Bernstein et al. (1987), Weikum & Vossen (2001), Kung & Robinson (1981). Direct evidence for RQ1 (lost update definition) and theoretical grounding for RQ2/RQ3.

---

## Source 4 (Martin Fowler - Optimistic Offline Lock)

Title: Optimistic Offline Lock
Publisher: Martin Fowler (Patterns of Enterprise Application Architecture)
URL: https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html
Published: 2003-03-05 (by David Rice)
Accessed: 2026-09-26 — Direct fetch verified. Verbatim quotes extracted: "assumes that the chance of conflict is low", "A successful pre-commit validation is, in a sense, obtaining a lock". Dated 2003-03-05 (author David Rice).
Source Tier: Tier 2 (established industry publication, expert technical article)
Relevance: Defines optimistic locking for business transactions spanning multiple system transactions. States assumption: "the chance of conflict is low." Describes pre-commit validation approach. Direct evidence for RQ3 (optimistic locking, when appropriate).

---

## Source 5 (Martin Fowler - Pessimistic Offline Lock)

Title: Pessimistic Offline Lock
Publisher: Martin Fowler (Patterns of Enterprise Application Architecture)
URL: https://martinfowler.com/eaaCatalog/pessimisticOfflineLock.html
Published: 2003-03-05 (author David Rice)
Accessed: 2026-09-26 — Direct fetch verified. Verbatim quotes extracted: "Pessimistic Offline Lock prevents conflicts by avoiding them altogether", "If several people access the same data within a business transaction, one commits easily but the others will conflict and fail. Since the conflict is only detected at the end of the business transaction, the victims will do all the transaction work only to find at the last minute that the whole thing will fail and their time will have been wasted."
Source Tier: Tier 2 (established industry publication, expert technical article)
Relevance: Defines pessimistic locking to "avoid conflicts altogether" by acquiring lock before use. Notes downsides when conflict rate is low vs Optimistic Offline Lock causing wasted work on conflicts. Direct evidence for RQ2 (pessimistic locking, when appropriate).

---

## Source 6

Title: Optimistic Concurrency Control
Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Optimistic_concurrency_control
Published: Continuously updated
Accessed: 2026-09-26 (via cross-reference from Concurrency Control article; content corroborated by subagent research)
Source Tier: Tier 2
Relevance: Academic reference for optimistic concurrency control theory (Kung & Robinson 1981: validate-read-write phases). Corroborates RQ3 mechanism definitions.

---

## Source 7

Title: MySQL 8.0 Reference Manual - 15.7.2.4 Locking Reads
Publisher: Oracle / MySQL
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html
Published: MySQL 8.0 (current)
Accessed: 2026-09-26 — Direct fetch verified via official Oracle CDN mirror (full document). Verbatim quotes extracted: SELECT ... FOR SHARE / FOR UPDATE semantics, "All locks are released when the transaction is committed or rolled back", NOWAIT error behavior (ERROR 3572), SKIP LOCKED example, counter increment example.
Source Tier: Tier 1 (official documentation)
Relevance: Documents InnoDB SELECT ... FOR UPDATE and FOR SHARE semantics, NOWAIT and SKIP LOCKED options. Direct evidence for RQ2 cross-database verification (MySQL pessimistic locking). Note: dev.mysql.com domain returns 403 to automated fetch; Oracle CDN mirror hosts identical official MySQL documentation and was fetched directly.

---

## Source 8

Title: MySQL 8.0 Reference Manual - 15.7.1 InnoDB Locking
Publisher: Oracle / MySQL
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-locking.html
Published: MySQL 8.0 (current)
Accessed: 2026-09-26 — Direct fetch verified via official Oracle CDN mirror.
Source Tier: Tier 1 (official documentation)
Relevance: Documents InnoDB lock types: shared (S), exclusive (X), intention locks, record locks, gap locks, next-key locks. Evidence for transaction-level locking internals, RQ2.

---

## Source 9

Title: MySQL 8.0 Reference Manual - 15.7.2.1 Transaction Isolation Levels
Publisher: Oracle / MySQL
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-transaction-isolation-levels.html
Published: MySQL 8.0 (current)
Accessed: 2026-09-26 — Direct fetch verified via official Oracle CDN mirror.
Source Tier: Tier 1 (official documentation)
Relevance: Documents InnoDB default REPEATABLE READ, READ COMMITTED semi-consistent reads, SERIALIZABLE converting plain SELECT to SELECT FOR SHARE. Evidence for RQ4 cross-database comparison.

---

## Source 10

Title: Oracle Database Concepts 19c - 9 Data Concurrency and Consistency
Publisher: Oracle
URL: https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html
Published: Oracle 19c
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Documents Oracle row locks (TX), lost update example table (Banda salary scenario), recommended WHERE-clause-with-old-values prevention, MVCC via undo segments, default READ COMMITTED isolation. Direct evidence for RQ1 (lost update example) and RQ2 (Oracle pessimistic).

---

## Source 11

Title: Oracle Database Concepts 19c - 10 Transactions
Publisher: Oracle
URL: https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/transactions.html
Published: Oracle 19c
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Documents ACID properties, statement-level atomicity, savepoint lock release behavior. Evidence for RQ5 (atomicity guarantees).

---

## Source 12

Title: Optimistic Locking in JPA
Publisher: Baeldung
URL: https://www.baeldung.com/jpa-optimistic-locking
Published: Date not verified (continuously updated tutorial site)
Accessed: 2026-09-26
Source Tier: Tier 2 (established industry technical publication)
Relevance: Documents @Version annotation approach, OptimisticLockException on concurrent update, recommendation for read-heavy workloads. Evidence for RQ3 implementation details.

---

## Source 13

Title: DynamoDB Transaction APIs
Publisher: Amazon Web Services
URL: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/transaction-apis.html
Published: Current AWS documentation
Accessed: 2026-09-26
Source Tier: Tier 1 (official cloud documentation)
Relevance: Documents TransactWriteItems atomicity guarantees, ConditionCheck for optimistic concurrency, capacity consumption (2x per item). Useful for contrasting single-database locking vs distributed NoSQL transactions. Supporting evidence for RQ5 and distributed considerations.

---

## Source 14

Title: Hibernate ORM 6.x User Guide - Optimistic Locking
Publisher: JBoss / Hibernate team
URL: https://docs.jboss.org/hibernate/orm/6.3/userguide/html_single/Hibernate_User_Guide.html#locking-optimistic
Published: Hibernate 6.3
Accessed: 2026-09-26 (via cached reference; URL provided by subagent research)
Source Tier: Tier 1 (official ORM documentation)
Relevance: Documents version-number strategy, timestamp strategy, versionless optimistic locking (ALL vs DIRTY), StaleObjectStateException. Evidence for RQ3 implementation variants.

**Verification note:** This source was reported by subagent but direct fetch was not performed during this session. Content is consistent with Jakarta Persistence specification for @Version. Treat versionless-locking details as MEDIUM confidence pending direct verification.

---

## Source 15

Title: PostgreSQL 18 Documentation - 13.4. Data Consistency Checks at the Application Level
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/applevel-consistency.html
Published: PostgreSQL 18.6 (current, accessed 2026-09-26)
Accessed: 2026-09-26 — Direct fetch verified (full document). Verbatim quotes extracted: "to ensure the current validity of a row and protect it against concurrent updates one must use SELECT FOR UPDATE, SELECT FOR SHARE, or an appropriate LOCK TABLE statement" and "SELECT FOR UPDATE does not ensure that a concurrent transaction will not update or delete a selected row. To do that in PostgreSQL you must actually update the row, even if no values need to be changed."
Source Tier: Tier 1 (official documentation)
Relevance: Direct evidence for RQ5 (atomic operations as alternative to explicit locking) and RQ6 (consistency enforcement). This source upgraded the atomic-update finding (Evidence 10 / Finding 6) from MEDIUM to HIGH confidence by establishing that an actual UPDATE is the authoritative conflict-prevention action in PostgreSQL.

---

## Source 16

Title: MySQL 8.0 Reference Manual - 17.7.2.4 Locking Reads (verified via official Oracle CDN mirror)
Publisher: Oracle / MySQL
URL: https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/innodb-locking-reads.html (official MySQL domain https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html returns 403 to automated fetch)
Published: MySQL 8.0 (current)
Accessed: 2026-09-26 — Direct fetch verified via Oracle CDN mirror (full document). Verbatim quotes extracted: SELECT ... FOR SHARE / FOR UPDATE semantics, "All locks set by FOR SHARE and FOR UPDATE queries are released when the transaction is committed or rolled back", NOWAIT error behavior (ERROR 3572), SKIP LOCKED example, counter increment example (SELECT counter FOR UPDATE; UPDATE counter = counter + 1), subquery locking-read caveat.
Source Tier: Tier 1 (official documentation hosted on Oracle CDN)
Relevance: Cross-database verification of pessimistic locking (RQ2), non-blocking lock variants relevant to the voucher flash-sale exercise, and the counter increment anti-pattern example directly relevant to the topic spec's read-modify-write discussion.
