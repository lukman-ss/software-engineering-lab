# Research Plan: Optimistic vs Pessimistic Locking

## Research Topic

Optimistic vs Pessimistic Locking — preventing lost updates caused by concurrent requests to the same data in database systems. The topic covers concurrency control mechanisms, specifically pessimistic locking (which prevents conflicts by blocking) and optimistic locking (which detects conflicts at commit time), as well as atomic database operations as an alternative to read-modify-write patterns.

## Objective

Investigate the theoretical foundations, practical implementations, and trade-offs of optimistic vs pessimistic locking approaches in database systems. The research will:

1. Define and verify the "lost update" problem with concrete examples
2. Document how pessimistic locking works (SELECT FOR UPDATE, shared/exclusive locks)
3. Document how optimistic locking works (version/timestamp-based conflict detection)
4. Examine the role of database isolation levels in preventing lost updates
5. Document atomic update operations as an alternative approach
6. Identify common mistakes and anti-patterns
7. Collect database-specific behaviors (PostgreSQL, MySQL/InnoDB, Oracle, SQL Server)
8. Document performance trade-offs: concurrency, throughput, deadlock risk

## Research Questions

### RQ1: What is the "lost update" problem and how does it manifest?
- Definition of lost update in database concurrency control
- Concrete example with stock/quantity data
- The role of isolation levels in lost update occurrence
- How different databases handle lost updates under default isolation

### RQ2: How does pessimistic locking work and when is it appropriate?
- Mechanism: blocking concurrent access using SELECT FOR UPDATE
- Lock modes: shared (S) and exclusive (X) locks, row-level vs table-level
- Database implementations: PostgreSQL FOR UPDATE, MySQL InnoDB locking reads, Oracle row locks
- Trade-offs: reduced concurrency, deadlock risk, lock ordering requirements
- Appropriate use cases: high-conflict scenarios, financial/balances, inventory allocation

### RQ3: How does optimistic locking work and when is it appropriate?
- Mechanism: version/timestamp checking at commit time
- The UPDATE WHERE id=? AND version=? pattern
- Affected rows check as conflict detection
- Application-level vs database-level implementation
- Trade-offs: overhead of version column, retry logic requirements, 0-rows-affected handling
- Appropriate use cases: low-conflict, read-heavy, offline/slow business transactions

### RQ4: How do database isolation levels affect these locking strategies?
- ANSI SQL isolation levels: READ UNCOMMITTED, READ COMMITTED, REPEATABLE READ, SERIALIZABLE
- Database-specific default isolation levels (PostgreSQL, MySQL, Oracle, SQL Server)
- How each level handles lost updates, dirty reads, non-repeatable reads, phantom reads
- Snapshot isolation as an alternative approach

### RQ5: How do atomic operations provide an alternative to explicit locking?
- The UPDATE ... WHERE condition >= value pattern
- Affected rows check as success/failure signal
- When atomic operations are sufficient vs when locking is needed
- Application-level vs database-level atomicity guarantees

### RQ6: What are common mistakes and anti-patterns in concurrency control?
- Assuming default transactions prevent all lost updates
- Using distributed locks (Redis) when database can solve the problem
- Holding pessimistic locks too long (network calls in transaction)
- Optimistic locking conflict not handled in application
- Read-modify-write pattern without concurrency consideration

## Search Strategy

### Primary Sources (Tier 1)
- PostgreSQL official documentation (https://www.postgresql.org/docs/)
- MySQL official documentation (https://dev.mysql.com/doc/)
- Oracle Database documentation (https://docs.oracle.com/en/database/)
- Microsoft SQL Server documentation (https://learn.microsoft.com/)
- SQL standards and academic references

### Secondary Sources (Tier 2)
- Martin Fowler's enterprise integration patterns (https://martinfowler.com/)
- Wikipedia articles on concurrency control, lost update
- Industry technical blogs (Baeldung, etc.)
- ORM documentation (Hibernate/JPA, Laravel)

### Search Approach
1. Access official documentation for each major database system
2. Verify key claims across multiple database vendors
3. Cross-check academic definitions (Wikipedia, textbooks)
4. Extract code examples and implementation details
5. Document database-specific behaviors and edge cases

## Expected Primary Sources

| Source | Topic | URL |
|--------|-------|-----|
| PostgreSQL 18 Docs - Explicit Locking | Pessimistic locking, FOR UPDATE, deadlocks | https://www.postgresql.org/docs/current/explicit-locking.html |
| PostgreSQL 18 Docs - Transaction Isolation | Isolation levels, Read Committed | https://www.postgresql.org/docs/current/transaction-iso.html |
| MySQL 8.0 Reference - Locking Reads | SELECT FOR UPDATE, FOR SHARE | https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html |
| MySQL 8.0 Reference - Lock Types | InnoDB lock types | https://dev.mysql.com/doc/refman/8.0/en/innodb-locking.html |
| Oracle 19c Concepts - Data Concurrency | Row-level locking, lost update | https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html |
| Oracle 19c Concepts - Transactions | ACID, atomicity | https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/transactions.html |
| Wikipedia - Concurrency Control | Definitions, academic references | https://en.wikipedia.org/wiki/Concurrency_control |
| Wikipedia - Optimistic Concurrency Control | OCC mechanism | https://en.wikipedia.org/wiki/Optimistic_concurrency_control |
| Martin Fowler - Optimistic Offline Lock | Application-level OCC | https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html |
| Martin Fowler - Pessimistic Offline Lock | Application-level PCC | https://martinfowler.com/eaaCatalog/pessimisticOfflineLock.html |

## Risks / Unknowns

1. **Database-specific behavior**: Different databases may have different default isolation levels and locking semantics. Need to verify for PostgreSQL, MySQL, Oracle, SQL Server separately.

2. **Isolation level anomalies**: The ANSI SQL isolation level definitions do not fully capture all possible anomalies (per Berenson et al. 1995). PostgreSQL's implementation of READ UNCOMMITTED behaves identically to READ COMMITTED (no true dirty read support).

3. **Snapshot isolation**: Not part of the ANSI standard levels but widely used. Behavior may vary between implementations.

4. **Lost update under REPEATABLE READ**: In PostgreSQL, REPEATABLE READ uses snapshot isolation which prevents lost updates; in MySQL (InnoDB) under REPEATABLE READ, the behavior may differ due to gap locking and the "2PL" approach.

5. **ORM-level locking**: Different ORMs (Laravel, Hibernate, SQLAlchemy) may implement optimistic locking differently. Need to focus on database-level mechanisms rather than ORM abstractions.

6. **Distributed locking**: The topic specification mentions distributed locks (Redis) as an anti-pattern. This research phase should document when distributed locks are appropriate vs database locks.

7. **Atomic operations**: The `UPDATE ... SET x = x - N WHERE id = ? AND x >= N` pattern effectiveness depends on isolation level and index behavior. Need to verify under which conditions this is truly safe.

8. **Version counter overflow**: Optimistic locking with integer version counters can overflow. This is a practical concern not always documented in theoretical sources.

9. **NOWAIT / SKIP LOCKED**: Modern databases support non-blocking lock acquisition variants. Behavior under contention differs from traditional blocking locks.
