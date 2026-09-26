# Contradictions: Optimistic vs Pessimistic Locking

Research date: 2026-09-26

---

## Contradiction 1: Repeatable Read Implementation Differs Across Databases

**SOURCE A:** PostgreSQL Documentation - Transaction Isolation (13.2)
> "Repeatable Read uses Snapshot Isolation... Serializable adds predicate locking (Serializable Snapshot Isolation) on top of Snapshot Isolation, detecting serialization anomalies."

**SOURCE B:** MySQL Documentation - Transaction Isolation Levels (15.7.2.1)
> "REPEATABLE READ (default): Consistent reads use snapshot from first read. Gap locks / next-key locks prevent phantom rows for locking reads."

**SOURCE C:** ANSI SQL Isolation Level Definition (Wikipedia Concurrency Control)
> Repeatable Read: "Not possible" for non-repeatable read, "Possible" for phantom read

**ASSESSMENT:** The term "Repeatable Read" means different things across implementations. PostgreSQL implements it as true Snapshot Isolation (SI) which prevents both non-repeatable reads AND write skews, but implements SERIALIZABLE as a separate level (SSI). MySQL uses a form of Strict 2PL with gap/next-key locks that prevent phantoms in "REPEATABLE READ" but the ANSI standard does not require this. Oracle implements neither as distinct levels. This creates confusion: code written for PostgreSQL's REPEATABLE READ is not equivalent to MySQL's REPEATABLE READ. The semantic gap means developers cannot assume "REPEATABLE READ prevents all anomalies" across all databases.

---

## Contradiction 2: READ UNCOMMITTED Behavior Varies

**SOURCE A:** PostgreSQL Documentation
> "Read Uncommitted behaves identically to Read Committed (prevents dirty reads via MVCC)."

**SOURCE B:** Other traditional databases (inferred from Wikipedia)
> READ UNCOMMITTED typically allows dirty reads (reading uncommitted data).

**ASSESSMENT:** PostgreSQL's MVCC architecture does not support true dirty reads, making READ UNCOMMITTED a no-op that behaves like READ COMMITTED. Other database systems (SQL Server, Oracle in some modes) can support dirty reads. This is a technical limitation of MVCC, not a design choice, causing applications that rely on READ UNCOMMITTED for performance in other databases to get no benefit in PostgreSQL.

---

## Contradiction 3: How Lost Update Manifests Under Different Isolation Levels

**SOURCE A:** PostgreSQL Transaction Isolation
Under READ COMMITTED: re-reads WHERE clause after concurrent commit can silently overwrite

**SOURCE B:** Oracle Lost Update Example (Table 10-2)
Under READ COMMITTED (default): Lost update occurs when Session 2 uses stale snapshot to update

**SOURCE C:** MySQL Transaction Isolation Levels under READ COMMITTED
> "Semi-consistent reads: For UPDATE, InnoDB returns latest committed version to evaluate WHERE condition"

**ASSESSMENT:** All three databases can exhibit lost update under READ COMMITTED, but the mechanism differs:
- PostgreSQL: each statement re-evaluates with latest committed data
- Oracle: read-consistent snapshot is taken at statement start, but the lost update anomaly occurs when overwriting with old value
- MySQL: uses semi-consistent reads to reduce deadlocks, still vulnerable if app doesn't check affected_rows

All three require explicit locking or WHERE-guard patterns. This is NOT a contradiction in outcomes but reveals different concurrency models.

---

## Contradiction 4: Transaction-Level Lock Release Semantics

**SOURCE A:** PostgreSQL Explicit Locking
> "Once acquired, a lock is normally held until the end of the transaction. But if a lock is acquired after establishing a savepoint, the lock is released immediately if the savepoint is rolled back to."

**SOURCE B:** MySQL InnoDB Documentation
> "All locks are released when the transaction is committed or rolled back."

**ASSESSMENT:** Both agree locks release at transaction end, but PostgreSQL has additional semantics for savepoint rollback releasing locks. This matters for long-running transactions with error handling. MySQL's lock duration across savepoints is less documented in the sources fetched.

---

## No Material Contradictions Discovered

The following areas show expected variation but not contradictions:

### 4.1 Oracle vs PostgreSQL Isolation Levels
- Oracle: READ COMMITTED, SERIALIZABLE, READ ONLY (no READ UNCOMMITTED, no REPEATABLE READ)
- PostgreSQL: READ COMMITTED, REPEATABLE READ, SERIALIZABLE (no READ UNCOMMITTED effect)

This is a difference in available options, not contradictory guidance. The core concurrency control principles remain the same.

### 4.2 Advisory Lock Variants
- PostgreSQL: session-level vs transaction-level advisory locks with different lifetime semantics
- Other databases: similar concepts (SQL Server sp_getapplock, MySQL GET_LOCK) but different APIs

Same pattern, different implementations.

---

## Resolution for Research Report

The identified contradictions are implementation-specific semantic differences under the same concurrency control model. They should be reported as:
1. "Repeatable Read" is not consistently defined across databases
2. PostgreSQL's MVCC does not support dirty reads (READ UNCOMMITTED is a no-op)
3. Each database's locking/release semantics have nuances relevant for edge cases
4. Lost update prevention requires explicit action (locking or WHERE-guard) in all systems under their default isolation levels