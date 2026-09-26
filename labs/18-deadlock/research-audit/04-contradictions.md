# Contradiction Audit

## Contradiction 1

Statement A: PostgreSQL uses periodic interval checking for deadlocks (`deadlock_timeout`).
Location: 05-report.md (Finding 4)

Statement B: MySQL/InnoDB uses active wait-for graph cycle detection dynamically.
Location: 04-contradictions.md

Type: SOURCE_CONFLICT
Impact: Not material for Postgres, but highlights difference between RDBMS implementations.
Assessment: The research notes the MySQL documentation was inaccessible (403), appropriately marking claims about MySQL as UNVERIFIED. No contradictory behavior documented internally for PostgreSQL.

## Contradiction 2

Statement A: Explicit locking and updates can lead to deadlocks (40P01).
Location: 05-report.md (Finding 1)

Statement B: Repeatable Read/Serializable isolation can lead to serialization failures (40001).
Location: 05-report.md (Finding 10)

Type: INTERNAL
Impact: Low
Assessment: These are complementary, not contradictory. The research appropriately clarifies the difference between deadlock errors and serialization failures.

No material contradictions found.
