# Source Audit: Optimistic vs Pessimistic Locking

## Source 1

Claimed Title: 13.3. Explicit Locking
Claimed Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/explicit-locking.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Accurately describes `FOR UPDATE`, `FOR SHARE`, row locks, deadlocks, and advisory locks.

Assessment:
PASS

---

## Source 2

Claimed Title: 13.2. Transaction Isolation
Claimed Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/transaction-iso.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Documents Read Committed behavior, Snapshot Isolation under Repeatable Read ("ERROR: could not serialize access due to concurrent update"), and SSI under Serializable.

Assessment:
PASS

---

## Source 3

Claimed Title: Data Concurrency and Consistency (Chapter 10)
Claimed Publisher: Oracle Corporation
URL: https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Explicitly details multiversion read consistency, TX row locks, Table 10-2 lost update under READ COMMITTED, and ORA-08177 under SERIALIZABLE.

Assessment:
PASS

---

## Source 4

Claimed Title: Concurrency Control
Claimed Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Concurrency_control

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Community reference, but adequately cites Bernstein et al. 1987 and Weikum & Vossen 2001 for theoretical classification (Optimistic, Pessimistic, Semi-optimistic, 2PL, MVCC).

Assessment:
PASS

---

## Source 5

Claimed Title: Optimistic Concurrency Control
Claimed Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Optimistic_concurrency_control

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Accurately details Kung & Robinson (1981), OCC phases (Begin, Modify, Validate, Commit/Rollback), and ecosystem examples.

Assessment:
PASS

---

## Source 6

Claimed Title: Optimistic Offline Lock
Claimed Publisher: Martin Fowler (Patterns of Enterprise Application Architecture)
URL: https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html

Reachable:
YES

Source Type:
SECONDARY (Authoritative Industry Pattern Catalog)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Web page is a pattern summary directing to Chapter 16 of PoEAA book, but core definition and assumption (conflict is unlikely) are fully verified.

Assessment:
PASS

---

## Source 7

Claimed Title: Pessimistic Offline Lock
Claimed Publisher: Martin Fowler (Patterns of Enterprise Application Architecture)
URL: https://martinfowler.com/eaaCatalog/pessimisticOfflineLock.html

Reachable:
YES

Source Type:
SECONDARY (Authoritative Industry Pattern Catalog)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Web page is a pattern summary directing to Chapter 16 of PoEAA book, but core definition and trade-offs are fully verified.

Assessment:
PASS

---

## Source 8

Claimed Title: Handling Concurrency Conflicts - EF Core
Claimed Publisher: Microsoft Learn
URL: https://learn.microsoft.com/en-us/ef/core/saving/concurrency

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Explicitly details `UPDATE ... WHERE Version = @p2`, `DbUpdateConcurrencyException`, conflict resolution, and comparison with transaction isolation levels.

Assessment:
PASS

---

## Source 9

Claimed Title: Write-write conflict (Lost Update)
Claimed Publisher: Wikipedia
URL: https://en.wikipedia.org/wiki/Write%E2%80%93write_conflict

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Article is flagged on Wikipedia as needing additional citations, but its formal definition matches Berenson et al. (1995) and Stearns & Rosenkrantz (1981).

Assessment:
PASS
