# Open Questions: Optimistic vs Pessimistic Locking

## Unanswered Questions

### 1. MySQL/InnoDB Specific Behavior
**Question**: How exactly does MySQL/InnoDB REPEATABLE READ handle lost updates compared to PostgreSQL's Snapshot Isolation?
- Known: MySQL uses next-key locks (gap + record locks) in REPEATABLE READ, which prevents phantom reads but behaves differently from PostgreSQL's MVCC-based snapshot.
- Need: Official MySQL documentation verification (was inaccessible during research).
- Impact: Applications porting between PostgreSQL and MySQL may encounter different concurrency behavior at the same named isolation level.

### 2. SQL Server Snapshot Isolation vs Serializable
**Question**: Does SQL Server's SNAPSHOT isolation level (enabled by `ALLOW_SNAPSHOT_ISOLATION`) provide the same lost-update prevention as PostgreSQL's REPEATABLE READ?
- EF Core docs mention: "SQL Server snapshot isolation level, as well as by the PostgreSQL repeatable reads isolation level" both implement a form of optimistic locking via serialization errors.
- Need: Direct verification from SQL Server documentation.

### 3. Version Counter Overflow in Optimistic Locking
**Question**: What is the practical impact of integer version counter overflow in high-throughput systems?
- An int32 version column overflows at ~2.1B updates.
- Need: Real-world data or best practices for handling (e.g., use BIGINT, timestamp, or GUID).
- Not addressed in primary sources consulted.

### 4. NOWAIT / SKIP LOCKED Behavior Under Contention
**Question**: How do `FOR UPDATE NOWAIT` and `FOR UPDATE SKIP LOCKED` behave differently from blocking `FOR UPDATE` in high-contention scenarios?
- PostgreSQL docs mention NOWAIT/SKIP LOCKED as non-blocking variants.
- Need: Detailed behavior and use cases (e.g., work queue patterns with SKIP LOCKED).
- Not fully explored in this research.

### 5. Advisory Locks vs Row Locks for Application-Level Coordination
**Question**: When should PostgreSQL advisory locks be used instead of row-level locks?
- Docs: "Advisory locks can be useful for locking strategies that are an awkward fit for the MVCC model."
- Need: Concrete examples where advisory locks outperform row locks.

### 6. Merge Strategies for Optimistic Conflict Resolution
**Question**: What are the recommended merge algorithms for different data types when optimistic locking detects a conflict?
- EF Core shows manual merge with "TODO: decide which value should be written to database."
- Need: Patterns for last-writer-wins, field-level merge, operational transform, CRDTs.
- Not covered in database docs; this is application-level design.

### 7. Optimistic Locking with Composite Entities
**Question**: How to implement optimistic locking when a business entity spans multiple tables/rows?
- Version column on one row doesn't protect related rows.
- Need: Patterns for multi-table optimistic locking (e.g., root entity version, separate lock table).

### 8. Distributed Lock Necessity Threshold
**Question**: At what scale/architecture does a distributed lock (Redis, etcd) become necessary vs database-level locking?
- Topic spec says: "Using distributed lock for problems that can be solved by database is an anti-pattern."
- Need: Concrete criteria (multi-database, cross-service, single-database with connection pooling limits).

## Weak Evidence Areas

### 1. Quantitative Performance Comparisons
- No benchmark data found in official docs comparing pessimistic vs optimistic vs atomic operations under load.
- All sources qualitative; no throughput/latency numbers.

### 2. Deadlock Probability in Practice
- PostgreSQL docs state deadlocks "typically low" likelihood but can occur.
- No data on how lock ordering, transaction duration, or isolation level affects probability.

### 3. Optimistic Locking Retry Storm Behavior
- Under high contention, optimistic locking can cause retry storms (thundering herd).
- No official guidance on backoff strategies, max retries, or circuit breaker patterns.

## Claims Needing Deeper Research

| Claim | Current Evidence | Needed |
|-------|------------------|--------|
| "Atomic UPDATE always safe" | Strong for single row; depends on WHERE clause design | Verify under SERIALIZABLE with predicate locks |
| "Optimistic better for low conflict" | Theoretical (Kung & Robinson 1981) | Empirical threshold definitions |
| "Pessimistic causes deadlocks" | Documented possibility | Frequency under realistic workloads |
| "MVCC prevents dirty reads" | Both PG & Oracle confirm | Verify edge cases (deferred inserts, etc.) |

## Possible Next Research Directions

1. **MySQL/InnoDB Deep Dive**: Fetch official docs on locking reads, gap locks, next-key locks, and isolation level behavior.

2. **ORM Comparison**: Compare Laravel `lockForUpdate()`, Hibernate `@Version`, SQLAlchemy `with_for_update()` — verify they generate correct SQL and handle edge cases.

3. **Distributed Concurrency**: Research when database locks aren't enough (sharding, microservices, saga patterns).

4. **Real-World Case Studies**: Analyze concurrency bugs from production incidents (e.g., GitHub, Shopify, Uber engineering blogs).

5. **Benchmark Suite**: Design reproducible benchmarks for:
   - Atomic UPDATE vs SELECT+UPDATE
   - Optimistic (version) vs Pessimistic (FOR UPDATE) vs SERIALIZABLE
   - Varying conflict rates (1%, 10%, 50%)
   - Varying transaction durations

6. **Advanced Patterns**: Research:
   - Semantic locking (lock by business key, not row ID)
   - Escrow locking (for aggregate constraints like "sum of balances")
   - Conflict-free replicated data types (CRDTs) for eventually consistent systems