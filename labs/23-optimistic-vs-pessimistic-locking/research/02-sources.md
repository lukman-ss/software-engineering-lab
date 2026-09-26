# Sources

## Source 1
Title: PostgreSQL Documentation: Explicit Locking (Section 13.3)
Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/explicit-locking.html
Published: N/A (Current as of PostgreSQL 18)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Official specification for `SELECT FOR UPDATE`, `SELECT FOR SHARE`, `FOR UPDATE`, `FOR SHARE`, `FOR KEY SHARE`, row-level lock modes, conflict behavior, and deadlock handling.
Notes: Covers all locking read variants; includes Table 13.2 (table-level) and Table 13.3 (row-level) lock compatibility matrices.

## Source 2
Title: PostgreSQL Documentation: Transaction Isolation (Section 13.2)
Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/transaction-iso.html
Published: N/A (Current as of PostgreSQL 18)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Defines READ COMMITTED, REPEATABLE READ, SERIALIZABLE isolation levels and their interaction with MVCC. Critical for understanding how optimistic vs pessimistic locking behaves under different isolation levels.
Notes: PostgreSQL's Repeatable Read implements Snapshot Isolation. Serialization anomaly detection described.

## Source 3
Title: PostgreSQL Documentation: Data Consistency Checks at the Application Level (Section 13.4)
Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/applevel-consistency.html
Published: N/A (Current as of PostgreSQL 18)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Practical guidance on when to use explicit blocking locks (`SELECT FOR UPDATE`, `SELECT FOR SHARE`) vs serializable transactions for consistency checks. Documents `SELECT FOR UPDATE` behavior nuances (temporarily blocks, requires actual UPDATE to enforce).
Notes: Documents the read/write conflict issue with serializable transactions and explicit locking patterns.

## Source 4
Title: PostgreSQL Documentation: MVCC Introduction (Section 13.1)
Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/mvcc-intro.html
Published: N/A (Current as of PostgreSQL 18)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Foundational explanation of Multiversion Concurrency Control (MVCC) and its relationship to optimistic locking. Explains how MVCC provides non-blocking reads while allowing writers to proceed.
Notes: Establishes why PostgreSQL's Repeatable Read/SERIALIZABLE can serve as optimistic concurrency control.

## Source 5
Title: Microsoft SQL Server Documentation: Transaction Locking and Row Versioning Guide
Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-transaction-locking-and-row-versioning-guide
Published: 2026-04-22 (Updated)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Authoritative guide to SQL Server's optimistic concurrency via row versioning (READ COMMITTED SNAPSHOT, SNAPSHOT isolation) and pessimistic locking via table hints (`UPDLOCK`, `UPDATENAME`, etc.). Describes deadlock avoidance patterns.
Notes: Documents SQL Server's "Optimized locking" (2023+) and row-versioning isolation levels.

## Source 6
Title: Oracle Database SQL Language Reference: SELECT Statement
Publisher: Oracle Corporation
URL: https://docs.oracle.com/en/database/oracle/oracle-database/23/sql/SELECT.html
Published: N/A (Current as of Oracle Database 23c)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Official documentation of `FOR UPDATE` clause syntax and semantics for Oracle's pessimistic locking. Includes `NOWAIT`, `WAIT`, `SKIP LOCKED` options.
Notes: Documents that Oracle uses `ORA_ROWSCN` pseudo-column for optimistic locking via `SELECT ... FOR UPDATE` alternatives.

## Source 7
Title: Hibernate ORM User Guide: Locking
Publisher: Red Hat Hibernate Team
URL: https://docs.jboss.org/hibernate/orm/6.6/userguide/html_single/Hibernate_User_Guide.html#locking
Published: 2024 (Hibernate ORM 6.6)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Official documentation of Hibernate's `@OptimisticLock`, `@OptimisticLocking`, `@Version`, and `@PessimisticLock` annotations. Describes Entity-level and Property-level optimistic locking.
Notes: Documents JPA locking modes and Hibernate-specific extensions.

## Source 8
Title: "Concurrency Control and Recovery in Database Systems" by Philip A. Bernstein, Vassos Hadzilacos, Nathan Goodman
Publisher: Addison-Wesley
Edition: 1st Edition (1987)
ISBN: 978-0201107159
URL: https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/ccontrol.pdf
Accessed: 2026-09-25
Source Tier: Tier 2 (Academic)
Relevance: Foundational textbook defining optimistic vs pessimistic concurrency control theory, timestamp ordering, validation protocols, and distributed transaction models.
Notes: Defines OCC validation phase and 2PL. Freely available PDF from Microsoft Research. Does not cover modern MVCC variants post-1987.

## Source 9
Title: "Database System Concepts" by Abraham Silberschatz, Henry F. Korth, S. Sudarshan
Publisher: McGraw-Hill Education
Edition: 7th Edition (2019)
ISBN: 978-0078022159
Accessed: 2026-09-25
Source Tier: Tier 2 (Academic)
Relevance: Standard database textbook covering locking protocols, Two-Phase Locking (2PL), strict 2PL, and serializability theory. Chapter 15 covers concurrency control in detail.
Notes: Defines shared/exclusive locks, lock compatibility, deadlock detection, and recovery. Edition 2019 includes modern MVCC discussion.

## Source 10
Title: MySQL 8.0 Reference Manual: InnoDB Locking Reads
Publisher: Oracle Corporation (MySQL)
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html
Published: MySQL 8.0 (Current)
Accessed: 2026-09-25 — NOT VERIFIED (Technical Difficulties error page returned)
Source Tier: Tier 1
Relevance: Official documentation of MySQL/InnoDB `SELECT ... FOR UPDATE`, `SELECT ... LOCK IN SHARE MODE`, `NOWAIT`, `SKIP LOCKED` syntax and behavior.
Notes: URL returned "Technical Difficulties" page on 2026-09-25 (verified across 8.0 and 5.7). MySQL locking behavior in this research inferred from PostgreSQL/SQL Server/Oracle patterns and community knowledge. Requires re-verification when MySQL docs accessible.

## Source 11
Title: Laravel Database Documentation: Locking
Publisher: Laravel
URL: https://laravel.com/docs/11.x/queries#locking-rows
Published: Laravel 11.x (2024)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Official documentation of Laravel's `lockForUpdate()`, `sharedLock()`, `lockForUpdateNowait()` query builder methods. Documents Eloquent ORM interpretation of SQL locking.
Notes: Laravel translates to native SQL locking per database driver (PostgreSQL, MySQL, SQL Server).

## Source 12
Title: Entity Framework Core Documentation: Concurrency
Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/ef/core/performance/efficient-query-patterns/concurrent-updates
Published: N/A (Current)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Official Microsoft guidance on EF Core optimistic concurrency using `[Timestamp]`/`rowversion` columns and `ConcurrencyCheck` attributes. Documents retry logic patterns.
Notes: EF Core translates optimistic concurrency to SQL `WHERE` clauses comparing original values.