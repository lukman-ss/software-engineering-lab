# Contradictions Audit: Optimistic vs Pessimistic Locking

## Contradiction 1: Meaning and Semantics of "Repeatable Read" Across Databases

Statement A:
PostgreSQL documentation states that `REPEATABLE READ` implements Snapshot Isolation (SI), preventing both non-repeatable reads and phantom reads without gap locks, but raising serialization failures on concurrent row modification.

Location:
`research/04-contradictions.md`, Contradiction 1; PostgreSQL Docs 13.2

Statement B:
MySQL/InnoDB documentation states that `REPEATABLE READ` (its default level) uses consistent read snapshots for plain SELECTs and 2PL with gap locks / next-key locks for locking reads.

Location:
`research/04-contradictions.md`, Contradiction 1; MySQL Docs 15.7.2.1

Type:
SOURCE_CONFLICT / IMPLEMENTATION_DIFFERENCE

Impact:
Applications written expecting PostgreSQL snapshot isolation error semantics will behave differently under MySQL gap-locking blocking behavior.

Assessment:
Properly identified and resolved as an implementation-specific divergence. Research report correctly warns developers that isolation levels are not uniform across database engines.

---

## Contradiction 2: READ UNCOMMITTED Implementation

Statement A:
PostgreSQL treats `READ UNCOMMITTED` as identical to `READ COMMITTED` because MVCC does not allow dirty reads.

Location:
`research/04-contradictions.md`, Contradiction 2; PostgreSQL Docs 13.2

Statement B:
ANSI SQL standard and traditional lock-based engines allow dirty reads under `READ UNCOMMITTED`.

Location:
`research/04-contradictions.md`, Contradiction 2; Wikipedia Concurrency Control

Type:
SOURCE_CONFLICT / IMPLEMENTATION_DIFFERENCE

Impact:
Developer performance optimizations relying on dirty reads in other systems provide zero benefit in PostgreSQL.

Assessment:
Accurately documented and resolved as an engine-specific limitation of PostgreSQL MVCC.

---

## Summary Assessment of Contradictions
No internal contradictions exist between the research plan, evidence, and report. All identified cross-vendor discrepancies reflect genuine differences in database architectures (MVCC snapshot isolation vs 2PL gap locks) and are documented transparently.
