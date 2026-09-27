# Research Report: Optimistic vs Pessimistic Locking

## Research Question

How do optimistic and pessimistic locking mechanisms prevent lost updates in concurrent database access, what are their trade-offs, and what alternative approaches exist?

## Executive Summary

This research confirms that the **lost update (write-write conflict)** is a fundamental concurrency anomaly occurring when two transactions read the same data, independently modify it, and the second commit silently overwrites the first. Neither PostgreSQL's nor Oracle's default **READ COMMITTED** isolation level prevents this — both documentation explicitly show lost update scenarios under default isolation.

Three primary strategies exist:
1. **Pessimistic Locking** — Block concurrent access using explicit locks (`SELECT FOR UPDATE`), preventing conflicts at the cost of reduced concurrency and deadlock risk.
2. **Optimistic Locking** — Allow concurrent access, detect conflicts at commit time via version/timestamp columns, and retry on failure.
3. **Atomic Operations** — Use single-statement database operations (`UPDATE ... SET x = x - N WHERE condition`) that eliminate the read-modify-write race window entirely.

The choice depends on conflict frequency, transaction duration, and whether the operation can be expressed atomically. Senior engineers should first consider atomic operations, then optimistic locking for low-conflict scenarios, then pessimistic locking for high-conflict/critical sections, and avoid distributed locks when database-level solutions suffice.

## Findings

### Finding 1: Lost Update is Not Prevented by Default Isolation Levels

**Claim**: The default READ COMMITTED isolation level in both PostgreSQL and Oracle does not prevent lost updates. Applications must explicitly implement concurrency control.

**Evidence**:
- PostgreSQL docs state that under Read Committed, "it is possible for an updating command to see an inconsistent snapshot" because each command gets a new snapshot.
- Oracle docs Table 10-2 explicitly demonstrates a lost update scenario in a READ COMMITTED transaction and states: "Devising a strategy to handle lost updates is an important part of application development."

**Sources**: PostgreSQL 13.2 Transaction Isolation; Oracle 19c Data Concurrency and Consistency

**Confidence**: HIGH

### Finding 2: Pessimistic Locking Prevents Conflicts by Blocking

**Claim**: Pessimistic locking uses explicit row-level locks (`SELECT FOR UPDATE`) to prevent other transactions from modifying the same row until the current transaction commits or rolls back.

**Evidence**:
- PostgreSQL: "FOR UPDATE causes the rows retrieved by the SELECT statement to be locked as though for update. This prevents them from being locked, modified or deleted by other transactions until the current transaction ends."
- Oracle: "A row lock, also called a TX lock, is a lock on a single row of table. A transaction acquires a row lock for each row modified by an INSERT, UPDATE, DELETE, MERGE, or SELECT ... FOR UPDATE statement."
- Martin Fowler: "Pessimistic Offline Lock prevents conflicts by avoiding them altogether. It forces a business transaction to acquire a lock on a piece of data before it starts to use it."

**Trade-offs documented**:
- Reduced concurrency (other transactions wait)
- Deadlock risk (PostgreSQL: "The use of explicit locking can increase the likelihood of deadlocks")
- Must keep transactions short (no network calls while holding locks)

**Sources**: PostgreSQL 13.3 Explicit Locking; Oracle 19c Locking Mechanism; Martin Fowler Pessimistic Offline Lock

**Confidence**: HIGH

### Finding 3: Optimistic Locking Detects Conflicts at Commit Time

**Claim**: Optimistic concurrency control allows transactions to proceed without locks, validating at commit time that no other transaction modified the data. If a conflict is detected, the transaction aborts and can retry.

**Evidence**:
- Wikipedia: "OCC assumes that multiple transactions can frequently complete without interfering with each other... Before committing, each transaction verifies that no other transaction has modified the data it has read. If the check reveals conflicting modifications, the committing transaction rolls back and can be restarted."
- EF Core implementation: `UPDATE ... WHERE [PersonId] = @p1 AND [Version] = @p2` — if 0 rows affected, throws `DbUpdateConcurrencyException`.
- Martin Fowler: "Optimistic Offline Lock solves this problem by validating that the changes about to be committed by one session don't conflict with the changes of another session."

**Trade-offs documented**:
- Higher throughput when conflicts are rare
- Requires application-level retry logic
- Version column overhead (integer overflow possible)
- 0-rows-affected must not be treated as success

**Sources**: Wikipedia Optimistic Concurrency Control; EF Core Concurrency; Martin Fowler Optimistic Offline Lock; Kung & Robinson 1981

**Confidence**: HIGH

### Finding 4: Database Isolation Levels Provide Alternative Concurrency Control

**Claim**: Higher isolation levels (REPEATABLE READ, SERIALIZABLE) can prevent lost updates without explicit application-level locking, but with different mechanisms and trade-offs.

**Evidence**:
- PostgreSQL REPEATABLE READ: Uses Snapshot Isolation. If a transaction tries to modify a row changed by another committed transaction after the transaction started, it rolls back with "could not serialize access due to concurrent update."
- PostgreSQL SERIALIZABLE: Uses Serializable Snapshot Isolation with predicate locking to prevent serialization anomalies.
- Oracle SERIALIZABLE: Raises ORA-08177 when a serializable transaction tries to update a row changed by another transaction that committed after the serializable transaction began.
- EF Core docs note: "Repeatable reads guarantee that your transaction always sees the same data across queries inside the transaction, avoiding inconsistencies." But "this approach requires a transaction to span all the operations."

**Sources**: PostgreSQL 13.2 Transaction Isolation; Oracle 19c Transaction Isolation Levels; EF Core Concurrency

**Confidence**: HIGH

### Finding 5: Atomic Database Operations Eliminate Race Windows for Simple Cases

**Claim**: For operations like decrementing a counter or checking-and-updating a condition, a single atomic UPDATE statement with a WHERE clause can prevent lost updates without explicit locking.

**Evidence**:
- PostgreSQL docs show: `UPDATE accounts SET balance = balance + 100.00 WHERE acctnum = 12345;` — this acquires a row lock automatically and applies atomically.
- The pattern `UPDATE products SET stock = stock - 3 WHERE id = 10 AND stock >= 3;` followed by checking affected rows is a complete atomic operation.
- This is the "read-modify-write" anti-pattern solution: let the database do the read and write atomically.

**Sources**: PostgreSQL 13.2 (transaction isolation examples showing atomic updates); Universal SQL semantics

**Confidence**: HIGH

### Finding 6: Common Mistakes Are Well-Documented

**Claim**: Five common concurrency anti-patterns are identified across sources.

**Evidence**:
1. **Assuming transactions prevent lost updates**: Both PostgreSQL and Oracle docs show lost updates under default isolation.
2. **Using distributed locks (Redis) unnecessarily**: Topic specification and Fowler's patterns suggest using database locks first.
3. **Holding pessimistic locks too long**: PostgreSQL warns: "it is a bad idea for applications to hold transactions open for long periods of time (e.g., while waiting for user input)."
4. **Optimistic locking without conflict handling**: EF Core requires catching `DbUpdateConcurrencyException` and retrying or merging.
5. **Read-modify-write without concurrency consideration**: The `$product->stock -= 1; $product->save();` pattern is the classic race condition.

**Sources**: PostgreSQL Explicit Locking; Oracle Conflicting Writes; Martin Fowler Patterns; EF Core Concurrency; Topic Specification

**Confidence**: HIGH

### Finding 7: MVCC Is the Foundation Enabling Both Approaches

**Claim**: Both PostgreSQL and Oracle use Multiversion Concurrency Control (MVCC) as their underlying architecture, which provides read consistency without blocking readers and enables both optimistic and pessimistic approaches.

**Evidence**:
- PostgreSQL: MVCC is the core concurrency architecture. "PostgreSQL's Read Uncommitted mode behaves like Read Committed" because of MVCC.
- Oracle: "Oracle Database maintains multiversion read consistency... Readers and writers of data do not block one another."
- MVCC generates new versions on write, allowing reads to access historical versions.

**Sources**: PostgreSQL 13.1 Introduction; Oracle 19c Multiversion Read Consistency; Wikipedia Multiversion Concurrency Control

**Confidence**: HIGH

## Areas of Agreement

| Topic | PostgreSQL | Oracle | Wikipedia | Martin Fowler | EF Core |
|-------|------------|--------|-----------|---------------|---------|
| Lost update occurs under READ COMMITTED | ✓ | ✓ | ✓ | ✓ | ✓ |
| SELECT FOR UPDATE blocks writers | ✓ | ✓ | ✓ (2PL) | ✓ | — |
| Optimistic locking via version check | ✓ (implicit) | ✓ (implicit) | ✓ | ✓ | ✓ |
| Deadlock risk with pessimistic | ✓ | ✓ | ✓ | ✓ | — |
| Atomic UPDATE is atomic | ✓ | ✓ | ✓ | — | — |
| MVCC as foundation | ✓ | ✓ | ✓ | — | — |

## Areas of Disagreement / Implementation Differences

| Aspect | PostgreSQL | Oracle | SQL Server (via EF Core) |
|--------|------------|--------|--------------------------|
| Default isolation | READ COMMITTED | READ COMMITTED | READ COMMITTED |
| REPEATABLE READ implementation | Snapshot Isolation (no phantoms) | Not offered | 2PL with shared locks |
| SERIALIZABLE implementation | Serializable Snapshot Isolation | 2PL-style + ORA-08177 | Snapshot Isolation or 2PL |
| Deadlock handling | Auto-detect, abort one | Auto-detect | Auto-detect |
| Native optimistic support | Application-level | Application-level | rowversion / concurrency tokens |

## Limitations

1. **MySQL/InnoDB behavior not directly verified** — Official documentation was inaccessible (403). Known industry differences: MySQL REPEATABLE READ uses gap locking (next-key locks) which behaves differently from PostgreSQL's snapshot isolation at the same named level.

2. **SQL Server specifics** — Referenced only via EF Core documentation, not directly from SQL Server docs.

3. **ORM-level implementations** — Only EF Core was examined. Laravel, Hibernate/JPA, SQLAlchemy may have different defaults and behaviors.

4. **Performance benchmarks** — No quantitative performance data (throughput, latency) was collected from official sources.

5. **Distributed systems** — The research focused on single-database concurrency. Distributed locking (Redis, etcd, Zookeeper) and distributed transactions (2PC, saga) are out of scope.

## Conclusion

The lost update problem is real and not solved by default database isolation. Three practical strategies exist, ordered by preference for most applications:

1. **Atomic operations** — Use single-statement updates with conditional WHERE clauses. Eliminates the race window entirely. Best for counters, stock decrement, quota management.

2. **Optimistic locking** — Add a version/timestamp column. Check `UPDATE ... WHERE version = ?`. Handle 0-rows-affected with retry/merge logic. Best for read-heavy, low-conflict scenarios (CRUD, profiles, documents).

3. **Pessimistic locking** — `SELECT FOR UPDATE` within short transactions. Best for high-conflict, critical sections (wallet balance, seat reservation, inventory allocation).

Senior engineers should ask: "Can this operation be made atomic?" before reaching for locking. If locking is needed, match the strategy to conflict frequency and business transaction duration.