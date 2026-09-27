# Sources: Optimistic vs Pessimistic Locking

## Source 1

Title: 13.3. Explicit Locking
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/explicit-locking.html
Published: PostgreSQL 18 (current)
Accessed: 2026-09-27
Source Tier: Tier 1 (Official Documentation)
Relevance: Primary reference for pessimistic locking mechanisms in PostgreSQL — FOR UPDATE, FOR SHARE, row-level locks, table-level locks, deadlocks, advisory locks.

## Source 2

Title: 13.2. Transaction Isolation
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/transaction-iso.html
Published: PostgreSQL 18 (current)
Accessed: 2026-09-27
Source Tier: Tier 1 (Official Documentation)
Relevance: Official documentation on isolation levels (Read Committed, Repeatable Read, Serializable), how PostgreSQL implements them via MVCC and Snapshot Isolation, and how lost updates are handled.

## Source 3

Title: Data Concurrency and Consistency (Chapter 10)
Publisher: Oracle Corporation
URL: https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html
Published: Oracle Database 19c Concepts Guide
Accessed: 2026-09-27
Source Tier: Tier 1 (Official Documentation)
Relevance: Oracle's multiversion read consistency, row locking (TX locks), isolation levels (Read Committed, Serializable), lost update behavior under Read Committed, locking mechanisms.

## Source 4

Title: Concurrency Control
Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Concurrency_control
Published: Reviewed 7 September 2026
Accessed: 2026-09-27
Source Tier: Tier 2 (Academic/Reference)
Relevance: Comprehensive overview of concurrency control theory — optimistic vs pessimistic categories, methods (2PL, timestamp ordering, serialization graph checking), ACID, serializability. References Bernstein et al. 1987, Weikum & Vossen 2001.

## Source 5

Title: Optimistic Concurrency Control
Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Optimistic_concurrency_control
Published: Last edited 26 February 2026
Accessed: 2026-09-27
Source Tier: Tier 2 (Academic/Reference)
Relevance: Formal definition of OCC by Kung & Robinson (1981). Three phases (Begin, Modify, Validate, Commit/Rollback). Real-world usage: MediaWiki, Bugzilla, Ruby on Rails, Entity Framework, Redis WATCH, CouchDB, Elasticsearch, DynamoDB, Kubernetes, Firestore, Iceberg.

## Source 6

Title: Optimistic Offline Lock
Publisher: Martin Fowler (Patterns of Enterprise Application Architecture)
URL: https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html
Published: 05 March 2003
Accessed: 2026-09-27
Source Tier: Tier 2 (Authoritative Industry Publication)
Relevance: Application-level optimistic concurrency pattern. Validates changes at commit time, detects conflicts, allows rollback. Assumption: conflict is unlikely.

## Source 7

Title: Pessimistic Offline Lock
Publisher: Martin Fowler (Patterns of Enterprise Application Architecture)
URL: https://martinfowler.com/eaaCatalog/pessimisticOfflineLock.html
Published: 05 March 2003
Accessed: 2026-09-27
Source Tier: Tier 2 (Authoritative Industry Publication)
Relevance: Application-level pessimistic concurrency pattern. Forces acquisition of lock before data use. Prevents conflicts by limiting concurrency. Useful when business transactions are long and conflicts are frequent.

## Source 8

Title: Handling Concurrency Conflicts - EF Core
Publisher: Microsoft Learn
URL: https://learn.microsoft.com/en-us/ef/core/saving/concurrency
Published: Updated 2025-10-30
Accessed: 2026-09-27
Source Tier: Tier 1 (Official Documentation)
Relevance: Practical implementation of optimistic concurrency via concurrency tokens (Version/timestamp columns), DbUpdateConcurrencyException handling, comparison of optimistic tokens vs isolation-level-based concurrency control (Repeatable Read / Snapshot).

## Source 9

Title: Write-write conflict (Lost Update)
Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Write%E2%80%93write_conflict
Published: Last edited 25 December 2025
Accessed: 2026-09-27
Source Tier: Tier 2 (Academic/Reference)
Relevance: Formal definition of write-write conflict (lost update). Dirty writes. Reference to Stearns & Rosenkrantz 1981. Strict 2PL as solution with deadlock caveat.
